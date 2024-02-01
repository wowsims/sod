package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const LightningBoltRanks = 10

// First entry is the base spell ID, second entry is the overload's spell ID
var LightningBoltSpellId = [LightningBoltRanks + 1][]int32{{0}, {403, 408439}, {529, 408440}, {548, 408441}, {915, 408442}, {943, 408443}, {6041, 408472}, {10391, 408473}, {10392, 408474}, {15207, 408475}, {15208, 408477}}
var LightningBoltBaseDamage = [LightningBoltRanks + 1][]float64{{0}, {15, 17}, {28, 33}, {48, 57}, {88, 100}, {131, 149}, {179, 202}, {235, 264}, {291, 326}, {357, 400}, {428, 477}}
var LightningBoltSpellCoef = [LightningBoltRanks + 1]float64{0, .1233, .314, .554, .857, .857, .857, .857, .857, .857, .857}
var LightningBoltCastTime = [LightningBoltRanks + 1]int32{0, 1500, 2000, 2500, 3000, 3000, 3000, 3000, 3000, 3000, 3000}
var LightningBoltManaCost = [LightningBoltRanks + 1]float64{0, 15, 30, 45, 75, 105, 135, 165, 195, 230, 265}
var LightningBoltLevel = [LightningBoltRanks + 1]int{0, 1, 8, 14, 20, 26, 32, 38, 44, 50, 56}

func (shaman *Shaman) registerLightningBoltSpell() {
	overloadRuneEquipped := shaman.HasRune(proto.ShamanRune_RuneChestOverload)

	for i := 1; i <= LightningBoltRanks; i++ {
		config := shaman.newLightningBoltSpellConfig(i, false)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.LightningBolt = shaman.RegisterSpell(config)

			if overloadRuneEquipped {
				shaman.LightningBoltOverload = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(i, true))
			}
		}
	}
}

func (shaman *Shaman) newLightningBoltSpellConfig(rank int, isOverload bool) core.SpellConfig {
	// First entry is the base spell ID, second entry is the overload's spell ID
	spellId := LightningBoltSpellId[rank][core.TernaryInt32(isOverload, 1, 0)]
	baseDamageLow := LightningBoltBaseDamage[rank][0]
	baseDamageHigh := LightningBoltBaseDamage[rank][1]
	spellCoeff := LightningBoltSpellCoef[rank]
	castTime := LightningBoltCastTime[rank]
	manaCost := LightningBoltManaCost[rank]
	level := LightningBoltLevel[rank]

	spell := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: spellId},
		manaCost,
		time.Millisecond*time.Duration(castTime),
		isOverload,
	)

	dmgBonus := shaman.electricSpellBonusDamage(spellCoeff)

	canOverload := !isOverload && shaman.OverloadAura != nil

	spell.RequiredLevel = level
	spell.Rank = rank
	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := dmgBonus + sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

		if canOverload && result.Landed() && sim.RandomFloat("LB Overload") < shaman.OverloadChance {
			shaman.LightningBoltOverload.Cast(sim, target)
		}

		spell.DealDamage(sim, result)
	}

	if isOverload {
		shaman.applyOverloadModifiers(&spell)
	}

	return spell
}
