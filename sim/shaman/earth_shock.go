package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const EarthShockRanks = 7

// First entry is the base spell ID, second entry is the overload's spell ID
var EarthShockSpellId = [EarthShockRanks + 1]int32{0, 8042, 8044, 8045, 8046, 10412, 10413, 10414}
var EarthShockBaseDamage = [EarthShockRanks + 1][]float64{{0}, {19, 22}, {35, 38}, {65, 69}, {126, 134}, {235, 249}, {372, 394}, {517, 545}}
var EarthShockSpellCoef = [EarthShockRanks + 1]float64{0, .154, .212, .299, .386, .386, .386, .386}
var EarthShockManaCost = [EarthShockRanks + 1]float64{0, 30, 50, 85, 145, 240, 345, 450}
var EarthShockLevel = [EarthShockRanks + 1]int{0, 4, 8, 14, 24, 36, 48, 60}

func (shaman *Shaman) registerEarthShockSpell(shockTimer *core.Timer) {
	// Way of Earth gives earth shock a separate timer
	if shaman.HasRune(proto.ShamanRune_RuneLegsWayOfEarth) {
		shockTimer = shaman.NewTimer()
	}

	shaman.EarthShock = make([]*core.Spell, EarthShockRanks+1)

	for rank := 1; rank <= EarthShockRanks; rank++ {
		config := shaman.newEarthShockSpellConfig(rank, shockTimer)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.EarthShock[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newEarthShockSpellConfig(rank int, shockTimer *core.Timer) core.SpellConfig {
	spellId := EarthShockSpellId[rank]
	baseDamageLow := EarthShockBaseDamage[rank][0]
	baseDamageHigh := EarthShockBaseDamage[rank][1]
	spellCoeff := EarthShockSpellCoef[rank]
	manaCost := EarthShockManaCost[rank]
	level := EarthShockLevel[rank]

	spell := shaman.newShockSpellConfig(
		core.ActionID{SpellID: spellId},
		core.SpellSchoolNature,
		manaCost,
		shockTimer,
	)

	spell.RequiredLevel = level
	spell.Rank = rank

	spell.ThreatMultiplier = shaman.ShamanThreatMultiplier(2)

	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellPower()
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	}

	return spell
}
