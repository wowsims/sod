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
	hasCoherenceRune := shaman.HasRune(proto.ShamanRune_RuneCloakCoherence)
	hasRiptideRune := shaman.HasRune(proto.ShamanRune_RuneBracersRiptide)

	spellId := ChainHealSpellId[rank]
	baseHealingMultiplier := 1 + shaman.purificationHealingModifier()
	baseHealingLow := ChainHealBaseHealing[rank][0] * baseHealingMultiplier
	baseHealingHigh := ChainHealBaseHealing[rank][1] * baseHealingMultiplier
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

	spell := core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellId},
		ClassSpellMask: ClassSpellMask_ShamanChainHeal,
		DefenseType:    core.DefenseTypeMagic,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          flags,

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

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			targets := sim.Environment.Raid.GetFirstNPlayersOrPets(targetCount)
			curTarget := targets[0]
			origMult := spell.GetDamageMultiplier()
			// TODO: This bounces to most hurt friendly...
			for hitIndex := 0; hitIndex < len(targets); hitIndex++ {
				originalDamageMultiplier := spell.GetDamageMultiplier()
				if hasRiptideRune && !isOverload && shaman.Riptide.Hot(curTarget).IsActive() {
					spell.ApplyMultiplicativeDamageBonus(1.25)
					shaman.Riptide.Hot(curTarget).Deactivate(sim)
				}
				spell.CalcAndDealHealing(sim, curTarget, sim.Roll(baseHealingLow, baseHealingHigh), spell.OutcomeHealingCrit)
				spell.SetMultiplicativeDamageBonus(originalDamageMultiplier)

				if !isOverload && shaman.procOverload(sim, "Chain Heal Overload", 1/3) {
					shaman.ChainHealOverload[rank].Cast(sim, target)
				}

				spell.ApplyMultiplicativeDamageBonus(bounceCoef)
				curTarget = targets[hitIndex]
			}
			spell.SetMultiplicativeDamageBonus(origMult)
		},
	}

	if isOverload {
		shaman.applyOverloadModifiers(&spell)
	}

	return spell
}
