package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const LightningBoltRanks = 10

var LightningBoltSpellId = [LightningBoltRanks + 1]int32{0, 403, 529, 548, 915, 943, 6041, 10391, 10392, 15207, 15208}
var LightningBoltBaseDamage = [LightningBoltRanks + 1][]float64{{0}, {15, 17}, {28, 33}, {48, 57}, {88, 100}, {131, 149}, {179, 202}, {230, 259}, {145, 163}, {347, 389}, {428, 477}}
var LightningBoltSpellCoef = [LightningBoltRanks + 1]float64{0, .1233, .314, .554, .857, .857, .857, .857, .857, .857, .857}
var LightningBoltCastTime = [LightningBoltRanks + 1]int32{0, 1500, 2000, 2500, 3000, 3000, 3000, 3000, 3000, 3000, 3000}
var LightningBoltManaCost = [LightningBoltRanks + 1]float64{0, 15, 30, 45, 75, 105, 135, 165, 195, 230, 265}
var LightningBoltLevel = [LightningBoltRanks + 1]int{0, 1, 8, 14, 20, 26, 32, 38, 44, 50, 56}

func (shaman *Shaman) registerLightningBoltSpell() {
	overloadRuneEquipped := shaman.HasRune(proto.ShamanRune_RuneChestOverload)

	shaman.LightningBolt = make([]*core.Spell, LightningBoltRanks+1)

	if overloadRuneEquipped {
		shaman.LightningBoltOverload = make([]*core.Spell, LightningBoltRanks+1)
	}

	for rank := 1; rank <= LightningBoltRanks; rank++ {
		config := shaman.newLightningBoltSpellConfig(rank, false)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.LightningBolt[rank] = shaman.RegisterSpell(config)

			if overloadRuneEquipped {
				shaman.LightningBoltOverload[rank] = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(rank, true))
			}
		}
	}
}

func (shaman *Shaman) newLightningBoltSpellConfig(rank int, isOverload bool) core.SpellConfig {
	spellId := LightningBoltSpellId[rank]
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
	spell.ClassSpellMask = ClassSpellMask_ShamanLightningBolt
	spell.MissileSpeed = 20
	spell.RequiredLevel = level
	spell.Rank = rank
	spell.BonusCoefficient = spellCoeff

	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			spell.DealDamage(sim, result)

			if !isOverload && shaman.procOverload(sim, "Lightning Bolt Overload", 1) {
				shaman.LightningBoltOverload[rank].Cast(sim, target)
			}
		})
	}

	if isOverload {
		shaman.applyOverloadModifiers(&spell)
	}

	return spell
}
