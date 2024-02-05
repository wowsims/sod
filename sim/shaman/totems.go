package shaman

import (
	"github.com/wowsims/sod/sim/core"
)

func (shaman *Shaman) newTotemSpellConfig(flatCost float64, spellID int32) core.SpellConfig {
	return core.SpellConfig{
		ActionID: core.ActionID{SpellID: spellID},
		Flags:    SpellFlagTotem | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   flatCost,
			Multiplier: shaman.TotemManaMultiplier(),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
	}
}

// func (shaman *Shaman) registerCallOfTheElements() {
// 	airTotem := shaman.getAirTotemSpell(shaman.Totems.Air)
// 	earthTotem := shaman.getEarthTotemSpell(shaman.Totems.Earth)
// 	fireTotem := shaman.getFireTotemSpell(shaman.Totems.Fire)
// 	waterTotem := shaman.getWaterTotemSpell(shaman.Totems.Water)

// 	totalManaCost := 0.0
// 	if airTotem != nil {
// 		totalManaCost += airTotem.DefaultCast.Cost
// 	}
// 	if earthTotem != nil {
// 		totalManaCost += earthTotem.DefaultCast.Cost
// 	}
// 	if fireTotem != nil {
// 		totalManaCost += fireTotem.DefaultCast.Cost
// 	}
// 	if waterTotem != nil {
// 		totalManaCost += waterTotem.DefaultCast.Cost
// 	}

// 	shaman.RegisterSpell(core.SpellConfig{
// 		ActionID: core.ActionID{SpellID: 66842},
// 		Flags:    core.SpellFlagAPL,

// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: core.GCDDefault,
// 			},
// 		},
// 		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
// 			return shaman.CurrentMana() >= totalManaCost
// 		},

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			// Save GCD timer value, so we can safely reset it between each totem cast.
// 			nextGcdAt := shaman.GCD.ReadyAt()

// 			if airTotem != nil {
// 				shaman.GCD.Set(sim.CurrentTime)
// 				airTotem.Cast(sim, target)
// 			}
// 			if earthTotem != nil {
// 				shaman.GCD.Set(sim.CurrentTime)
// 				earthTotem.Cast(sim, target)
// 			}
// 			if fireTotem != nil {
// 				shaman.GCD.Set(sim.CurrentTime)
// 				fireTotem.Cast(sim, target)
// 			}
// 			if waterTotem != nil {
// 				shaman.GCD.Set(sim.CurrentTime)
// 				waterTotem.Cast(sim, target)
// 			}

// 			shaman.GCD.Set(nextGcdAt)
// 		},
// 	})
// }

// func (shaman *Shaman) getAirTotemSpell(totemType proto.AirTotem) *core.Spell {
// 	switch totemType {
// 	case proto.AirTotem_GraceOfAirTotem:
// 		return shaman.GraceOfAirTotem
// 	case proto.AirTotem_WindfuryTotem:
// 		return shaman.WindfuryTotem
// 	}
// 	return nil
// }

// func (shaman *Shaman) getEarthTotemSpell(totemType proto.EarthTotem) *core.Spell {
// 	switch totemType {
// 	case proto.EarthTotem_StrengthOfEarthTotem:
// 		return shaman.StrengthOfEarthTotem
// 	case proto.EarthTotem_TremorTotem:
// 		return shaman.TremorTotem
// 	case proto.EarthTotem_StoneskinTotem:
// 		return shaman.StoneskinTotem
// 	}
// 	return nil
// }

// func (shaman *Shaman) getFireTotemSpell(totemType proto.FireTotem) *core.Spell {
// 	switch totemType {
// 	case proto.FireTotem_TotemOfWrath:
// 		return shaman.TotemOfWrath
// 	case proto.FireTotem_SearingTotem:
// 		return shaman.SearingTotem
// 	case proto.FireTotem_MagmaTotem:
// 		return shaman.MagmaTotem
// 	case proto.FireTotem_FlametongueTotem:
// 		return shaman.FlametongueTotem
// 	}
// 	return nil
// }

// func (shaman *Shaman) getWaterTotemSpell(totemType proto.WaterTotem) *core.Spell {
// 	switch totemType {
// 	case proto.WaterTotem_ManaSpringTotem:
// 		return shaman.ManaSpringTotem
// 	case proto.WaterTotem_HealingStreamTotem:
// 		return shaman.HealingStreamTotem
// 	}
// 	return nil
// }
