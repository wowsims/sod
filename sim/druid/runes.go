package druid

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (druid *Druid) ApplyRunes() {
	// Helm
	druid.applyGaleWinds()
	druid.applyGore()

	// Chest
	druid.applyFuryOfStormRage()
	druid.applyWildStrikes()

	// Bracers
	druid.applyElunesFires()

	// Hands
	druid.applyMangle()
	druid.registerSunfireSpell()

	// Belt
	druid.applyBerserk()
	druid.applyEclipse()

	// Legs
	druid.applyStarsurge()
	druid.applySavageRoar()

	// Feet
	druid.applyDreamstate()
	druid.applyKingOfTheJungle()
}

func (druid *Druid) applyGaleWinds() {
	if !druid.HasRune(proto.DruidRune_RuneHelmGaleWinds) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label:    "Gale Winds",
		ActionID: core.ActionID{SpellID: int32(proto.DruidRune_RuneHelmGaleWinds)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (druid *Druid) applyGore() {
	if !druid.HasRune(proto.DruidRune_RuneHelmGore) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label:    "Gore",
		ActionID: core.ActionID{SpellID: int32(proto.DruidRune_RuneHelmGore)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

const (
	Gore_BearResetProcChance = .15
	Gore_CatResetProcChance  = .05
)

// TODO: Bear spells not implemented: MangleBear, Swipe, Maul
func (druid *Druid) rollGoreBearReset(sim *core.Simulation) {
	if sim.RandomFloat("Gore (Bear)") < Gore_BearResetProcChance {
		druid.MangleBear.CD.Reset()
	}
}

func (druid *Druid) rollGoreCatReset(sim *core.Simulation) {
	if sim.RandomFloat("Gore (Cat)") < Gore_CatResetProcChance {
		druid.TigersFury.CD.Reset()
	}
}

func (druid *Druid) applyFuryOfStormRage() {
	if !druid.HasRune(proto.DruidRune_RuneChestFuryOfStormrage) {
		return
	}

	druid.FuryOfStormrageAura = druid.RegisterAura(core.Aura{
		Label:    "Fury Of Stormrage",
		ActionID: core.ActionID{SpellID: int32(proto.DruidRune_RuneChestFuryOfStormrage)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (druid *Druid) applyEclipse() {
	if !druid.HasRune(proto.DruidRune_RuneBeltEclipse) {
		return
	}

	solarCritBonus := 30.0
	lunarCastTimeReduction := time.Second * 1

	var affectedSolarSpells []*DruidSpell
	var affectedLunarSpells []*DruidSpell

	// Solar
	druid.SolarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:     "Solar Eclipse proc",
		Duration:  time.Second * 15,
		MaxStacks: 4,
		ActionID:  core.ActionID{SpellID: 408250},
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSolarSpells = core.FilterSlice(
				core.Flatten([][]*DruidSpell{druid.Wrath, {druid.Starsurge}}),
				func(spell *DruidSpell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSolarSpells, func(spell *DruidSpell) {
				spell.BonusCritRating += solarCritBonus
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSolarSpells, func(spell *DruidSpell) {
				spell.BonusCritRating -= solarCritBonus
			})
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode != SpellCode_DruidWrath && spell.SpellCode != SpellCode_DruidStarsurge {
				return
			}

			aura.RemoveStack(sim)
		},
	})

	// Lunar
	druid.LunarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:     "Lunar Eclipse proc",
		Duration:  time.Second * 15,
		MaxStacks: 4,
		ActionID:  core.ActionID{SpellID: 408255},
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedLunarSpells = core.FilterSlice(
				core.Flatten([][]*DruidSpell{druid.Starfire}),
				func(spell *DruidSpell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedLunarSpells, func(spell *DruidSpell) {
				spell.DefaultCast.CastTime -= lunarCastTimeReduction
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedLunarSpells, func(spell *DruidSpell) {
				spell.DefaultCast.CastTime += lunarCastTimeReduction
			})
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode != SpellCode_DruidStarfire {
				return
			}

			aura.RemoveStack(sim)
		},
	})

	druid.EclipseAura = druid.RegisterAura(core.Aura{
		Label:    "Eclipse",
		Duration: core.NeverExpires,
		ActionID: core.ActionID{SpellID: 408248},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !slices.Contains([]int32{SpellCode_DruidWrath, SpellCode_DruidStarfire, SpellCode_DruidStarsurge}, spell.SpellCode) || !result.Landed() {
				return
			}

			if spell.SpellCode == SpellCode_DruidWrath || spell.SpellCode == SpellCode_DruidStarsurge {
				druid.LunarEclipseProcAura.Activate(sim)
				// Solar gives 1 stack of lunar bonus
				druid.LunarEclipseProcAura.AddStack(sim)
			}

			if spell.SpellCode == SpellCode_DruidStarfire || spell.SpellCode == SpellCode_DruidStarsurge {
				druid.SolarEclipseProcAura.Activate(sim)
				// Lunar gives 2 staacks of solar bonus
				druid.SolarEclipseProcAura.AddStacks(sim, 2)
			}
		},
	})

}

func (druid *Druid) applyElunesFires() {
	if !druid.HasRune(proto.DruidRune_RuneBracersElunesFires) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label:    "Elune's Fires",
		ActionID: core.ActionID{SpellID: int32(proto.DruidRune_RuneBracersElunesFires)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

const (
	ElunesFires_MaxExtensions = 3

	ElunesFires_BonusMoonfireTicks = int32(2)
	ElunesFires_BonusSunfireTicks  = int32(1)
	ElunesFires_BonusRipTicks      = int32(1)

	ElunesFires_MaxBonusMoonfireTicks = ElunesFires_BonusMoonfireTicks * ElunesFires_MaxExtensions
	ElunesFires_MaxSunfireTicks       = SunfireTicks + ElunesFires_BonusSunfireTicks*ElunesFires_MaxExtensions
	ElunesFires_MaxRipTicks           = RipTicks + ElunesFires_BonusRipTicks*ElunesFires_MaxExtensions
)

func (druid *Druid) tryElunesFiresMoonfireExtension(sim *core.Simulation, unit *core.Unit) {
	for _, moonfire := range druid.Moonfire {
		if moonfire != nil {
			if dot := moonfire.Dot(unit); dot.IsActive() && dot.NumberOfTicks < MoonfireDotTicks[moonfire.Rank]+ElunesFires_MaxBonusMoonfireTicks {
				dot.NumberOfTicks += ElunesFires_BonusMoonfireTicks
				dot.RecomputeAuraDuration()
				dot.UpdateExpires(sim, dot.ExpiresAt()+time.Duration(ElunesFires_BonusMoonfireTicks)*dot.TickPeriod())
			}
		}
	}

	if dot := druid.Sunfire.Dot(unit); dot.IsActive() && dot.NumberOfTicks < ElunesFires_MaxSunfireTicks {
		dot.NumberOfTicks += ElunesFires_BonusSunfireTicks
		dot.RecomputeAuraDuration()
		dot.UpdateExpires(sim, dot.ExpiresAt()+time.Duration(ElunesFires_BonusSunfireTicks)*dot.TickPeriod())
	}
}

func (druid *Druid) tryElunesFiresSunfireExtension(sim *core.Simulation, unit *core.Unit) {
	if dot := druid.Sunfire.Dot(unit); dot.IsActive() && dot.NumberOfTicks < ElunesFires_MaxSunfireTicks {
		dot.NumberOfTicks += ElunesFires_BonusSunfireTicks
		dot.RecomputeAuraDuration()
		dot.UpdateExpires(sim, dot.ExpiresAt()+time.Duration(ElunesFires_BonusSunfireTicks)*dot.TickPeriod())
	}
}

func (druid *Druid) tryElunesFiresRipExtension(sim *core.Simulation, unit *core.Unit) {
	if dot := druid.Rip.Dot(unit); dot.IsActive() && dot.NumberOfTicks < ElunesFires_MaxRipTicks {
		dot.NumberOfTicks += ElunesFires_BonusRipTicks
		dot.RecomputeAuraDuration()
		dot.UpdateExpires(sim, dot.ExpiresAt()+time.Duration(ElunesFires_BonusRipTicks)*dot.TickPeriod())
	}
}

func (druid *Druid) applyMangle() {
	//druid.registerMangleBearSpell()
	druid.registerMangleCatSpell()
}

func (druid *Druid) applyWildStrikes() {
	if !druid.HasRune(proto.DruidRune_RuneChestWildStrikes) {
		return
	}

	druid.WildStrikesBuffAura = core.ApplyWildStrikes(druid.GetCharacter())
}

func (druid *Druid) applyKingOfTheJungle() {
	if druid.HasRune(proto.DruidRune_RuneFeetKingOfTheJungle) {
		druid.registerTigersFurySpellKotJ()
	}
}

func (druid *Druid) applyDreamstate() {
	if !druid.HasRune(proto.DruidRune_RuneFeetDreamstate) {
		return
	}

	manaRegenDuration := time.Second * 8

	manaRegenAura := druid.RegisterAura(core.Aura{
		Label:    "Dreamstate Mana Regen",
		ActionID: core.ActionID{SpellID: int32(proto.DruidRune_RuneFeetDreamstate)},
		Duration: manaRegenDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.SpiritRegenRateCasting += .5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.SpiritRegenRateCasting -= .5
		},
	})

	// Hidden aura
	druid.RegisterAura(core.Aura{
		Label:    "Dreamstate Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}

			manaRegenAura.Activate(sim)
			core.DreamstateAura(result.Target).Activate(sim)
		},
	})
}
