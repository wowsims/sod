package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ChainHealRanks = 3
const ChainHealTargetCount = int32(3)

var ChainHealSpellId = [ChainHealRanks + 1]int32{0, 1064, 10622, 10623}
var ChainHealBaseHealing = [ChainHealRanks + 1][]float64{{0}, {332, 381}, {416, 477}, {567, 646}}
var ChainHealSpellCoef = [ChainHealRanks + 1]float64{0, .714, .714, .714}
var ChainHealCastTime = [ChainHealRanks + 1]int32{0, 2500, 2500, 2500}
var ChainHealManaCost = [ChainHealRanks + 1]float64{0, 260, 315, 405}
var ChainHealLevel = [ChainHealRanks + 1]int{0, 40, 46, 54}

func (shaman *Shaman) registerChainHealSpell() {
	overloadRuneEquipped := shaman.HasRune(proto.ShamanRune_RuneChestOverload)

	shaman.ChainHeal = make([]*core.Spell, ChainHealRanks+1)

	if overloadRuneEquipped {
		shaman.ChainHealOverload = make([]*core.Spell, ChainHealRanks+1)
	}

	for rank := 1; rank <= ChainHealRanks; rank++ {
		config := shaman.newChainHealSpellConfig(rank, false)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.ChainHeal[rank] = shaman.RegisterSpell(config)

			if overloadRuneEquipped {
				shaman.ChainHealOverload[rank] = shaman.RegisterSpell(shaman.newChainHealSpellConfig(rank, true))
			}
		}
	}
}

func (shaman *Shaman) newChainHealSpellConfig(rank int, isOverload bool) core.SpellConfig {
	hasOverloadRune := shaman.HasRune(proto.ShamanRune_RuneChestOverload)
	hasCoherenceRune := shaman.HasRune(proto.ShamanRune_RuneCloakCoherence)

	spellId := ChainHealSpellId[rank]
	baseHealingLow := ChainHealBaseHealing[rank][0] * (1 + shaman.purificationHealingModifier())
	baseHealingHigh := ChainHealBaseHealing[rank][1] * (1 + shaman.purificationHealingModifier())
	spellCoeff := ChainHealSpellCoef[rank]
	castTime := ChainHealCastTime[rank]
	manaCost := ChainHealManaCost[rank]
	level := ChainHealLevel[rank]

	flags := core.SpellFlagHelpful
	if !isOverload {
		flags |= core.SpellFlagAPL
	}

	bounceCoef := .5 // 50% reduction per bounce
	targetCount := ChainHealTargetCount
	if hasCoherenceRune {
		bounceCoef = .65 // 35% reduction per bounce
		targetCount += 2
	}

	canOverload := !isOverload && hasOverloadRune

	spell := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellCode:   SpellCode_ShamanChainHeal,
		DefenseType: core.DefenseTypeMagic,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       flags,

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
		},

		BonusCritRating: float64(shaman.Talents.TidalMastery) * core.CritRatingPerCritChance,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			targets := sim.Environment.Raid.GetFirstNPlayersOrPets(targetCount)
			curTarget := targets[0]
			origMult := spell.DamageMultiplier
			// TODO: This bounces to most hurt friendly...
			for hitIndex := 0; hitIndex < len(targets); hitIndex++ {
				baseHealing := sim.Roll(baseHealingLow, baseHealingHigh)

				result := spell.CalcAndDealHealing(sim, curTarget, baseHealing, spell.OutcomeHealingCrit)

				if canOverload && result.Landed() && sim.RandomFloat("CH Overload") < ShamanOverloadChance {
					shaman.ChainHealOverload[rank].Cast(sim, target)
				}

				spell.DamageMultiplier *= bounceCoef
				curTarget = targets[hitIndex]
			}
			spell.DamageMultiplier = origMult
		},
	}

	if isOverload {
		shaman.applyOverloadModifiers(&spell)
	}

	return spell
}
