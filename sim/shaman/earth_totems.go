package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const StrengthOfEarthTotemRanks = 5

var StrengthOfEarthTotemSpellId = [StrengthOfEarthTotemRanks + 1]int32{0, 8075, 8160, 8161, 10442, 25361}
var StrengthOfEarthTotemManaCost = [StrengthOfEarthTotemRanks + 1]float64{0, 25, 65, 125, 225, 275}
var StrengthOfEarthTotemLevel = [StrengthOfEarthTotemRanks + 1]int{0, 10, 24, 38, 52, 60}

func (shaman *Shaman) registerStrengthOfEarthTotemSpell() {
	shaman.StrengthOfEarthTotem = make([]*core.Spell, StrengthOfEarthTotemRanks+1)

	for rank := 1; rank <= StrengthOfEarthTotemRanks; rank++ {
		config := shaman.newStrengthOfEarthTotemSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.StrengthOfEarthTotem[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newStrengthOfEarthTotemSpellConfig(rank int) core.SpellConfig {
	spellId := StrengthOfEarthTotemSpellId[rank]
	manaCost := StrengthOfEarthTotemManaCost[rank]
	level := StrengthOfEarthTotemLevel[rank]

	duration := time.Second * 120

	spell := shaman.newTotemSpellConfig(manaCost, spellId)
	spell.RequiredLevel = level
	spell.Rank = rank
	spell.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		shaman.TotemExpirations[EarthTotem] = sim.CurrentTime + duration
		shaman.ActiveTotems[EarthTotem] = spell

		core.StrengthOfEarthTotemAura(&shaman.Unit, shaman.Level, 1+.075*float64(shaman.Talents.EnhancingTotems)).Activate(sim)
	}
	return spell
}

const StoneskinTotemRanks = 6

var StoneskinTotemSpellId = [StoneskinTotemRanks + 1]int32{0, 8071, 8154, 8155, 10406, 10407, 10408}
var StoneskinTotemManaCost = [StoneskinTotemRanks + 1]float64{0, 30, 60, 90, 115, 160, 210}
var StoneskinTotemLevel = [StoneskinTotemRanks + 1]int{0, 4, 14, 24, 34, 44, 54}

func (shaman *Shaman) registerStoneskinTotemSpell() {
	shaman.StoneskinTotem = make([]*core.Spell, StoneskinTotemRanks+1)

	for rank := 1; rank <= StoneskinTotemRanks; rank++ {
		config := shaman.newStoneskinTotemSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.StoneskinTotem[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newStoneskinTotemSpellConfig(rank int) core.SpellConfig {
	spellId := StoneskinTotemSpellId[rank]
	manaCost := StoneskinTotemManaCost[rank]
	level := StoneskinTotemLevel[rank]

	duration := time.Second * 120

	spell := shaman.newTotemSpellConfig(manaCost, spellId)
	spell.RequiredLevel = level
	spell.Rank = rank
	spell.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		shaman.TotemExpirations[EarthTotem] = sim.CurrentTime + duration
		shaman.ActiveTotems[EarthTotem] = spell
	}
	return spell
}

func (shaman *Shaman) registerTremorTotemSpell() {
	spellId := int32(8143)
	manaCost := float64(60)
	duration := time.Second * 120
	level := 18

	spell := shaman.newTotemSpellConfig(manaCost, spellId)
	spell.RequiredLevel = level
	spell.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		shaman.TotemExpirations[EarthTotem] = sim.CurrentTime + duration
		shaman.ActiveTotems[EarthTotem] = spell
	}
	shaman.TremorTotem = shaman.RegisterSpell(spell)
}
