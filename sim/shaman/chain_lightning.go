package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ChainLightningRanks = 4
const ChainLightningTargetCount = int32(3)

var ChainLightningSpellId = [ChainLightningRanks + 1]int32{0, 421, 930, 2860, 10605}
var ChainLightningBaseDamage = [ChainLightningRanks + 1][]float64{{0}, {200, 227}, {288, 323}, {383, 430}, {505, 564}}
var ChainLightningSpellCoef = [ChainLightningRanks + 1]float64{0, .714, .714, .714, .714}
var ChainLightningManaCost = [ChainLightningRanks + 1]float64{0, 280, 380, 490, 605}
var ChainLightningLevel = [ChainLightningRanks + 1]int{0, 32, 40, 48, 56}

func (shaman *Shaman) registerChainLightningSpell() {
	overloadRuneEquipped := shaman.HasRune(proto.ShamanRune_RuneChestOverload)

	shaman.ChainLightning = make([]*core.Spell, ChainLightningRanks+1)

	if overloadRuneEquipped {
		shaman.ChainLightningOverload = make([]*core.Spell, ChainLightningRanks+1)
	}

	cdTimer := shaman.NewTimer()

	for rank := 1; rank <= ChainLightningRanks; rank++ {
		config := shaman.newChainLightningSpellConfig(rank, cdTimer, false)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.ChainLightning[rank] = shaman.RegisterSpell(config)

			// TODO: Confirm how CL Overloads work in SoD
			if overloadRuneEquipped {
				shaman.ChainLightningOverload[rank] = shaman.RegisterSpell(shaman.newChainLightningSpellConfig(rank, cdTimer, true))
			}
		}
	}
}

func (shaman *Shaman) newChainLightningSpellConfig(rank int, cdTimer *core.Timer, isOverload bool) core.SpellConfig {
	hasOverloadRune := shaman.HasRune(proto.ShamanRune_RuneChestOverload)
	hasCoherenceRune := shaman.HasRune(proto.ShamanRune_RuneCloakCoherence)
	hasStormEarthAndFireRune := shaman.HasRune(proto.ShamanRune_RuneCloakStormEarthAndFire)

	spellId := ChainLightningSpellId[rank]
	baseDamageLow := ChainLightningBaseDamage[rank][0]
	baseDamageHigh := ChainLightningBaseDamage[rank][1]
	spellCoeff := ChainLightningSpellCoef[rank]
	manaCost := ChainLightningManaCost[rank]
	level := ChainLightningLevel[rank]

	cooldown := time.Second * 6
	if hasStormEarthAndFireRune {
		cooldown /= 2
	}
	castTime := time.Millisecond * 2500

	bounceCoef := .7 // 30% reduction per bounce
	targetCount := ChainLightningTargetCount
	if hasCoherenceRune {
		bounceCoef = .8 // 20% reduction per bounce
		targetCount += 2
	}

	canOverload := !isOverload && hasOverloadRune
	overloadChance := .1667

	spell := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: spellId},
		manaCost,
		castTime,
		isOverload,
	)

	spell.SpellCode = SpellCode_ShamanChainLightning
	spell.RequiredLevel = level
	spell.Rank = rank
	spell.BonusCoefficient = spellCoeff

	if !isOverload {
		spell.Cast.CD = core.Cooldown{
			Timer:    cdTimer,
			Duration: cooldown,
		}
	}

	results := make([]*core.SpellResult, min(targetCount, shaman.Env.GetNumTargets()))

	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		origMult := spell.DamageMultiplier
		for hitIndex := range results {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			results[hitIndex] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			target = sim.Environment.NextTargetUnit(target)
			spell.DamageMultiplier *= bounceCoef
		}

		for _, result := range results {
			spell.DealDamage(sim, result)

			if canOverload && sim.Proc(overloadChance, "CL Overload") {
				shaman.ChainLightningOverload[rank].Cast(sim, result.Target)
			}
		}

		spell.DamageMultiplier = origMult
	}

	if isOverload {
		shaman.applyOverloadModifiers(&spell)
	}

	return spell
}
