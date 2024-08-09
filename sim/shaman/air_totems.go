package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
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
	spell.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		shaman.TotemExpirations[EarthTotem] = sim.CurrentTime + duration
		shaman.ActiveTotems[EarthTotem] = spell
	}
	return spell
}

const GraceOfAirTotemRanks = 3

var GraceOfAirTotemSpellId = [GraceOfAirTotemRanks + 1]int32{0, 8835, 10627, 25359}
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
	manaCost := GraceOfAirTotemManaCost[rank]
	level := GraceOfAirTotemLevel[rank]

	duration := time.Second * 120
	multiplier := []float64{1, 1.08, 1.15}[shaman.Talents.EnhancingTotems]

	hasFeralSpirit := shaman.HasRune(proto.ShamanRune_RuneCloakFeralSpirit)

	graceOfAirTotemAuras := make([]*core.Aura, core.TernaryInt32(hasFeralSpirit, 3, 1))
	graceOfAirTotemAuras[0] = core.GraceOfAirTotemAura(&shaman.Unit, shaman.Level, multiplier)
	if hasFeralSpirit {
		graceOfAirTotemAuras[1] = core.GraceOfAirTotemAura(&shaman.SpiritWolves.SpiritWolf1.Unit, shaman.Level, multiplier)
		graceOfAirTotemAuras[2] = core.GraceOfAirTotemAura(&shaman.SpiritWolves.SpiritWolf2.Unit, shaman.Level, multiplier)
	}

	spell := shaman.newTotemSpellConfig(manaCost, spellId)
	spell.RequiredLevel = level
	spell.Rank = rank
	spell.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		shaman.TotemExpirations[AirTotem] = sim.CurrentTime + duration
		shaman.ActiveTotems[AirTotem] = spell

		for _, aura := range graceOfAirTotemAuras {
			aura.Activate(sim)
		}
	}
	return spell
}

const WindwallTotemRanks = 3

var WindwallTotemSpellId = [WindwallTotemRanks + 1]int32{0, 15107, 15111, 15112}
var WindwallTotemManaCost = [WindwallTotemRanks + 1]float64{0, 115, 170, 225}
var WindwallTotemLevel = [WindwallTotemRanks + 1]int{0, 36, 46, 56}

func (shaman *Shaman) registerWindwallTotemSpell() {
	shaman.WindwallTotem = make([]*core.Spell, WindwallTotemRanks+1)

	for rank := 1; rank <= WindwallTotemRanks; rank++ {
		config := shaman.newWindwallTotemSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.WindwallTotem[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newWindwallTotemSpellConfig(rank int) core.SpellConfig {
	has6PEarthfuryResolve := shaman.HasSetBonus(ItemSetEarthfuryResolve, 6)

	spellId := WindwallTotemSpellId[rank]
	manaCost := WindwallTotemManaCost[rank]
	level := WindwallTotemLevel[rank]

	duration := time.Second * 120

	windwallTotemAura := core.ImprovedWindwallTotemAura(&shaman.Unit)

	spell := shaman.newTotemSpellConfig(manaCost, spellId)
	spell.RequiredLevel = level
	spell.Rank = rank
	spell.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		shaman.TotemExpirations[AirTotem] = sim.CurrentTime + duration
		shaman.ActiveTotems[AirTotem] = spell

		// We don't have a separation of Melee vs Ranged bonus physical damage at the moment
		// but it's probably fine because bosses generally don't have ranged physical attacks.
		// core.WindwallTotemAura(&shaman.Unit, shaman.Level, shaman.Talents.GuardianTotems).Activate(sim)
		if has6PEarthfuryResolve {
			windwallTotemAura.Activate(sim)
		}
	}
	return spell
}
