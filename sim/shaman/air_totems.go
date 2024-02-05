package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const WindfuryTotemRanks = 3

var WindfuryTotemSpellId = [WindfuryTotemRanks + 1]int32{0, 8512, 10613, 10614}
var WindfuryTotemBonusDamage = [WindfuryTotemRanks + 1]float64{0, 122, 229, 315}
var WindfuryTotemManaCost = [WindfuryTotemRanks + 1]float64{0, 115, 175, 250}
var WindfuryTotemLevel = [WindfuryTotemRanks + 1]int{0, 32, 42, 52}

func (shaman *Shaman) registerWindfuryTotemSpell() {
	shaman.WindfuryTotem = make([]*core.Spell, WindfuryTotemRanks+1)

	for rank := 1; rank <= WindfuryTotemRanks; rank++ {
		config := shaman.newWindfuryTotemSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.WindfuryTotem[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newWindfuryTotemSpellConfig(rank int) core.SpellConfig {
	spellId := WindfuryTotemSpellId[rank]
	// TODO: The sim won't respect the value of a totem dropped via the APL. It uses hard-coded values from buffs.go
	// bonusDamage := WindfuryTotemBonusDamage[rank]
	manaCost := WindfuryTotemManaCost[rank]
	level := WindfuryTotemLevel[rank]

	duration := time.Second * 120

	spell := shaman.newTotemSpellConfig(manaCost, spellId)
	spell.RequiredLevel = level
	spell.Rank = rank
	spell.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.TotemExpirations[AirTotem] = sim.CurrentTime + duration
	}
	return spell
}

const GraceOfAirTotemRanks = 3

var GraceOfAirTotemSpellId = [GraceOfAirTotemRanks + 1]int32{0, 8835, 10627, 25359}
var GraceOfAirTotemBonusAgi = [GraceOfAirTotemRanks + 1]float64{0, 43, 67, 77}
var GraceOfAirTotemManaCost = [GraceOfAirTotemRanks + 1]float64{0, 155, 250, 310}
var GraceOfAirTotemLevel = [GraceOfAirTotemRanks + 1]int{0, 42, 56, 60}

func (shaman *Shaman) registerGraceOfAirTotemSpell() {
	shaman.GraceOfAirTotem = make([]*core.Spell, GraceOfAirTotemRanks+1)

	for rank := 1; rank <= GraceOfAirTotemRanks; rank++ {
		config := shaman.newGraceOfAirTotemSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.GraceOfAirTotem[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newGraceOfAirTotemSpellConfig(rank int) core.SpellConfig {
	spellId := GraceOfAirTotemSpellId[rank]
	// TODO: The sim won't respect the value of a totem dropped via the APL. It uses hard-coded values from buffs.go
	// bonusDamage := GraceOfAirTotemBonusAgi[rank]
	manaCost := GraceOfAirTotemManaCost[rank]
	level := GraceOfAirTotemLevel[rank]

	duration := time.Second * 120

	spell := shaman.newTotemSpellConfig(manaCost, spellId)
	spell.RequiredLevel = level
	spell.Rank = rank
	spell.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.TotemExpirations[AirTotem] = sim.CurrentTime + duration
	}
	return spell
}
