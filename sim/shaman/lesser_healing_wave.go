package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const LesserHealingWaveRanks = 6

var LesserHealingWaveSpellId = [LesserHealingWaveRanks + 1]int32{0, 8004, 8008, 8010, 10466, 10467, 10468}
var LesserHealingWaveBaseHealing = [LesserHealingWaveRanks + 1][]float64{{0}, {170, 195}, {257, 292}, {347, 391}, {473, 529}, {649, 723}, {832, 928}}
var LesserHealingWaveSpellCoef = [LesserHealingWaveRanks + 1]float64{0, .429, .429, .429, .429, .429, .429}
var LesserHealingWaveCastTime = [LesserHealingWaveRanks + 1]int32{0, 1500, 1500, 1500, 1500, 1500, 1500}
var LesserHealingWaveManaCost = [LesserHealingWaveRanks + 1]float64{0, 105, 145, 185, 235, 305, 380}
var LesserHealingWaveLevel = [LesserHealingWaveRanks + 1]int{0, 20, 28, 36, 44, 52, 60}

func (shaman *Shaman) registerLesserHealingWaveSpell() {
	shaman.LesserHealingWave = make([]*core.Spell, LesserHealingWaveRanks+1)

	for rank := 1; rank <= LesserHealingWaveRanks; rank++ {
		config := shaman.newLesserHealingWaveSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.LesserHealingWave[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newLesserHealingWaveSpellConfig(rank int) core.SpellConfig {
	hasMaelstromWeaponRune := shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon)

	spellId := LesserHealingWaveSpellId[rank]
	baseHealingMultiplier := 1 + shaman.purificationHealingModifier()
	baseHealingLow := LesserHealingWaveBaseHealing[rank][0] * baseHealingMultiplier
	baseHealingHigh := LesserHealingWaveBaseHealing[rank][1] * baseHealingMultiplier
	spellCoeff := LesserHealingWaveSpellCoef[rank]
	castTime := LesserHealingWaveCastTime[rank]
	manaCost := LesserHealingWaveManaCost[rank]
	level := LesserHealingWaveLevel[rank]

	switch shaman.Ranged().ID {
	case TotemOfTheStorm:
		baseHealingLow += 53
		baseHealingHigh += 53
	}

	spell := core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellId},
		ClassSpellMask: ClassSpellMask_ShamanLesserHealingWave,
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagHelpful | core.SpellFlagAPL,
		MetricSplits:   MaelstromWeaponSplits,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * time.Duration(castTime),
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				castTime := shaman.ApplyCastSpeedForSpell(cast.CastTime, spell)
				if hasMaelstromWeaponRune {
					stacks := shaman.MaelstromWeaponAura.GetStacks()
					spell.SetMetricsSplit(stacks)
					if stacks > 0 {
						return
					}
				}

				if castTime > 0 {
					shaman.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+castTime, false)
				}
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealHealing(sim, spell.Unit, sim.Roll(baseHealingLow, baseHealingHigh), spell.OutcomeHealingCrit)

			if result.Outcome.Matches(core.OutcomeCrit) {
				if shaman.HasRune(proto.ShamanRune_RuneFeetAncestralAwakening) {
					shaman.ancestralHealingAmount = result.Damage * AncestralAwakeningHealMultiplier

					// TODO: this should actually target the lowest health target in the raid.
					//  does it matter in a sim? We currently only simulate tanks taking damage (multiple tanks could be handled here though.)
					shaman.AncestralAwakening.Cast(sim, spell.Unit)
				}
			}
		},
	}

	return spell
}
