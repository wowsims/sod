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

	// Shoulder
	druid.applyShoulderRuneEffect()

	// Cloak
	druid.registerStarfallCD()
	druid.registerSwipeCatSpell()

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

func (druid *Druid) applyShoulderRuneEffect() {
	if druid.Equipment.Shoulders().Rune == int32(proto.DruidRune_DruidRuneNone) {
		return
	}

	switch druid.Equipment.Shoulders().Rune {
	// Balance
	case int32(proto.DruidRune_RuneShouldersLunatic):
		druid.applyT1Balance4PBonus()
	case int32(proto.DruidRune_RuneShouldersStarcaller):
		druid.applyT1Balance6PBonus()
	case int32(proto.DruidRune_RuneShouldersNight):
		druid.applyT2Balance2PBonus()
	case int32(proto.DruidRune_RuneShouldersKeepers):
		druid.applyT2Balance4PBonus()
	case int32(proto.DruidRune_RuneShouldersWrathful):
		druid.applyT2Balance6PBonus()
	case int32(proto.DruidRune_RuneShouldersCometcaller):
		druid.applyZGBalance3PBonus()
	case int32(proto.DruidRune_RuneShouldersForest):
		druid.applyZGBalance5PBonus()
	case int32(proto.DruidRune_RuneShouldersGraceful):
		druid.applyTAQBalance2PBonus()
	case int32(proto.DruidRune_RuneShouldersAstralAscendant):
		druid.applyTAQBalance4PBonus()

	// Cat
	case int32(proto.DruidRune_RuneShouldersIlluminator):
		druid.applyT1Feral2PBonus()
	case int32(proto.DruidRune_RuneShouldersPredatoryInstincts):
		druid.applyT1Feral4PBonus()
	case int32(proto.DruidRune_RuneShouldersRipper):
		druid.applyT1Feral6PBonus()
	case int32(proto.DruidRune_RuneShouldersClaw):
		druid.applyT2Feral2PBonus()
	case int32(proto.DruidRune_RuneShouldersPrideful):
		druid.applyT2Feral4PBonus()
	case int32(proto.DruidRune_RuneShouldersBarbaric):
		druid.applyT2Feral6PBonus()
	case int32(proto.DruidRune_RuneShouldersFrenetic):
		druid.applyTAQFeral2PBonus()
	case int32(proto.DruidRune_RuneShouldersExsanguinator):
		druid.applyTAQFeral4PBonus()
	case int32(proto.DruidRune_RuneShouldersAnimalisticExpertise):
		druid.applyRAQFeral3PBonus()

	// Guardian
	case int32(proto.DruidRune_RuneShouldersShifter):
		druid.applyT1Guardian6PBonus()
	case int32(proto.DruidRune_RuneShouldersTerritorial):
		druid.applyT2Guardian2PBonus()
	case int32(proto.DruidRune_RuneShouldersBeast):
		druid.applyT2Guardian4PBonus()
	case int32(proto.DruidRune_RuneShouldersLacerator):
		druid.applyT2Guardian6PBonus()
	case int32(proto.DruidRune_RuneShouldersFurious):
		druid.applyTAQGuardian2PBonus()
	case int32(proto.DruidRune_RuneShouldersMangler):
		druid.applyTAQGuardian4PBonus()
	}
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
	Gore_CatResetProcChance  = .15
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
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			switch spell.SpellCode {
			case SpellCode_DruidWrath:
				druid.tryElunesFiresSunfireExtension(sim, result.Target)
			case SpellCode_DruidStarfire:
				druid.tryElunesFiresMoonfireExtension(sim, result.Target)
			case SpellCode_DruidStarsurge: // Starsurge now benefits from the effects of Wrath and Starfire
				druid.tryElunesFiresSunfireExtension(sim, result.Target)
				druid.tryElunesFiresMoonfireExtension(sim, result.Target)
			case SpellCode_DruidShred:
				druid.tryElunesFiresRipExtension(sim, result.Target)
			}
		},
	})
}

const (
	ElunesFires_BonusMoonfireTime = time.Second * 6
	ElunesFires_BonusSunfireTime  = time.Second * 3
	ElunesFires_BonusRipTime      = time.Second * 2
)

func (druid *Druid) tryElunesFiresMoonfireExtension(sim *core.Simulation, unit *core.Unit) {
	for _, spell := range druid.Moonfire {
		if dot := spell.Dot(unit); dot.IsActive() {
			maxExpiresAt := sim.CurrentTime + dot.Duration
			dot.UpdateExpires(sim, min(maxExpiresAt, dot.ExpiresAt()+ElunesFires_BonusMoonfireTime))
		}
	}
}

func (druid *Druid) tryElunesFiresSunfireExtension(sim *core.Simulation, unit *core.Unit) {
	if druid.Sunfire == nil {
		return
	}
	if dot := druid.Sunfire.Dot(unit); dot.IsActive() {
		maxExpiresAt := sim.CurrentTime + dot.Duration
		dot.UpdateExpires(sim, min(maxExpiresAt, dot.ExpiresAt()+ElunesFires_BonusSunfireTime))
	}
}

func (druid *Druid) tryElunesFiresRipExtension(sim *core.Simulation, unit *core.Unit) {
	if dot := druid.Rip.Dot(unit); dot.IsActive() {
		maxExpiresAt := sim.CurrentTime + dot.Duration
		dot.UpdateExpires(sim, min(maxExpiresAt, dot.ExpiresAt()+ElunesFires_BonusRipTime))
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

	dreamstateAuras := druid.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
		return core.DreamstateAura(target)
	})

	druid.DreamstateManaRegenAura = druid.RegisterAura(core.Aura{
		Label:    "Dreamstate Mana Regen",
		ActionID: core.ActionID{SpellID: int32(proto.DruidRune_RuneFeetDreamstate)},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.SpiritRegenRateCasting += .5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.SpiritRegenRateCasting -= .5
		},
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: "Dreamstate Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskSpellDamage) && result.DidCrit() || spell.SpellCode == SpellCode_DruidStarsurge {
				druid.DreamstateManaRegenAura.Activate(sim)
				dreamstateAuras.Get(result.Target).Activate(sim)
			}
		},
	}))
}
