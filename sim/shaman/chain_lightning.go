package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ChainLightningRanks = 4
const ChainLightningTargetCount = 3

// 30% reduction per bounce
const ChainLightningBounceCoeff = .7

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

	for rank := 1; rank <= ChainLightningRanks; rank++ {
		config := shaman.newChainLightningSpellConfig(rank, false)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.ChainLightning[rank] = shaman.RegisterSpell(config)

			// TODO: Confirm how CL Overloads work in SoD
			if overloadRuneEquipped {
				shaman.ChainLightningOverload[rank] = shaman.RegisterSpell(shaman.newChainLightningSpellConfig(rank, true))
			}
		}
	}
}

func (shaman *Shaman) newChainLightningSpellConfig(rank int, isOverload bool) core.SpellConfig {
	spellId := ChainLightningSpellId[rank]
	baseDamageLow := ChainLightningBaseDamage[rank][0]
	baseDamageHigh := ChainLightningBaseDamage[rank][1]
	spellCoeff := ChainLightningSpellCoef[rank]
	manaCost := ChainLightningManaCost[rank]
	level := ChainLightningLevel[rank]

	cooldown := time.Second * 6
	castTime := time.Millisecond * 2500

	bonusDamage := shaman.electricSpellBonusDamage(spellCoeff)
	canOverload := !isOverload && shaman.HasRune(proto.ShamanRune_RuneChestOverload)

	spell := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: spellId},
		manaCost,
		castTime,
		isOverload,
	)

	spell.SpellCode = SpellCode_ShamanChainLightning
	spell.RequiredLevel = level
	spell.Rank = rank

	if !isOverload {
		spell.Cast.CD = core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: cooldown,
		}
	}

	numHits := min(ChainLightningTargetCount, shaman.Env.GetNumTargets())

	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		curTarget := target
		bounceCoeff := 1.0
		for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
			baseDamage := bonusDamage + sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			baseDamage *= bounceCoeff
			result := spell.CalcAndDealDamage(sim, curTarget, baseDamage*shaman.ConcussionMultiplier(), spell.OutcomeMagicHitAndCrit)

			if canOverload && result.Landed() && sim.RandomFloat("CL Overload") <= ShamanOverloadChance {
				shaman.ChainLightningOverload[rank].Cast(sim, curTarget)
			}

			bounceCoeff *= ChainLightningBounceCoeff
			curTarget = sim.Environment.NextTargetUnit(curTarget)
		}
	}

	if isOverload {
		shaman.applyOverloadModifiers(&spell)
	}

	return spell
}
