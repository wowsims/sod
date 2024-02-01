package shaman

import (
	"github.com/wowsims/sod/sim/core"
)

const FrostShockRanks = 4

// First entry is the base spell ID, second entry is the overload's spell ID
var FrostShockSpellId = [FrostShockRanks + 1]int32{0, 8056, 8058, 10472, 10473}
var FrostShockBaseDamage = [FrostShockRanks + 1][]float64{{0}, {95, 101}, {215, 230}, {345, 366}, {492, 520}}
var FrostShockSpellCoef = [FrostShockRanks + 1]float64{0, .386, .386, .386, .386}
var FrostShockManaCost = [FrostShockRanks + 1]float64{0, 115, 225, 325, 430}
var FrostShockLevel = [FrostShockRanks + 1]int{0, 20, 34, 46, 58}

func (shaman *Shaman) registerFrostShockSpell(shockTimer *core.Timer) {
}

func (shaman *Shaman) newFrostShockSpellConfig(shockTimer *core.Timer, rank int) core.SpellConfig {
	spellId := FrostShockSpellId[rank]
	baseDamageLow := FrostShockBaseDamage[rank][0]
	baseDamageHigh := FrostShockBaseDamage[rank][1]
	spellCoeff := FrostShockSpellCoef[rank]
	manaCost := FrostShockManaCost[rank]
	level := FrostShockLevel[rank]

	spell := shaman.newShockSpellConfig(
		core.ActionID{SpellID: spellId},
		core.SpellSchoolFrost,
		manaCost,
		shockTimer,
	)

	spell.RequiredLevel = level
	spell.Rank = rank

	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellPower()
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	}

	return spell
}
