package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const LesserHealingWaveRanks = 6

var LesserHealingWaveSpellId = [LesserHealingWaveRanks + 1]int32{0, 8004, 8008, 8010, 10466, 10467, 10468}
var LesserHealingWaveBaseHealing = [LesserHealingWaveRanks + 1][]float64{{0}, {170, 195}, {257, 292}, {349, 394}, {473, 529}, {649, 723}, {832, 928}}
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
	spellId := LesserHealingWaveSpellId[rank]
	baseHealingLow := LesserHealingWaveBaseHealing[rank][0]
	baseHealingHigh := LesserHealingWaveBaseHealing[rank][1]
	spellCoeff := LesserHealingWaveSpellCoef[rank]
	castTime := LesserHealingWaveCastTime[rank]
	manaCost := LesserHealingWaveManaCost[rank]
	level := LesserHealingWaveLevel[rank]

	spell := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellCode:   int32(SpellCode_LesserHealingWave),
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
			Multiplier: 1 *
				(1 - .01*float64(shaman.Talents.TidalFocus)),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * time.Duration(castTime),
			},
		},

		BonusCritRating: float64(shaman.Talents.TidalMastery) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(shaman.Talents.Purification)),
		CritMultiplier:   shaman.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healPower := spell.HealingPower(target)
			baseHealing := sim.Roll(baseHealingLow, baseHealingHigh) + spellCoeff*healPower
			result := spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			if result.Outcome.Matches(core.OutcomeCrit) {
				if shaman.HasRune(proto.ShamanRune_RuneFeetAncestralAwakening) {
					shaman.ancestralHealingAmount = result.Damage * AncestralAwakeningHealMultiplier

					// TODO: this should actually target the lowest health target in the raid.
					//  does it matter in a sim? We currently only simulate tanks taking damage (multiple tanks could be handled here though.)
					shaman.AncestralAwakening.Cast(sim, target)
				}
			}
		},
	}

	return spell
}
