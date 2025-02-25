package druid

import (
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
	druid.registerSunfireSpell()
	druid.applyMangle()

	// Belt
	druid.applyBerserk()
	druid.applyEclipse()

	// Legs
	druid.applyStarsurge()
	druid.applySavageRoar()
	druid.registerLacerateBleedSpell()
	druid.registerLacerateSpell()

	// Feet
	druid.applyDreamstate()
	druid.applyKingOfTheJungle()
	druid.applySurvivalInstincts()
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
	case int32(proto.DruidRune_RuneShoulderFerocious):
		druid.applyT1Guardian4PBonus()
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

	affectedBearSpells := ClassSpellMask_DruidLacerate | ClassSpellMask_DruidSwipeBear | ClassSpellMask_DruidMaul
	affectedCatSpells := ClassSpellMask_DruidMangleCat | ClassSpellMask_DruidShred

	actionID := core.ActionID{SpellID: int32(proto.DruidRune_RuneHelmGore)}
	rageMetrics := druid.NewRageMetrics(actionID)

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: "Gore Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell.Matches(affectedBearSpells) {
				if !sim.Proc(0.15, "Gore Bear") {
					return
				}

				if druid.MangleBear != nil {
					druid.MangleBear.CD.Reset()
				}
				druid.AddRage(sim, 10, rageMetrics)
			} else if spell.Matches(affectedCatSpells) && sim.Proc(0.15, "Gore Bear") {
				druid.TigersFury.CD.Reset()
			}
		},
	}))
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

	// Solar
	druid.SolarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:     "Solar Eclipse proc",
		Duration:  time.Second * 15,
		MaxStacks: 4,
		ActionID:  core.ActionID{SpellID: 408250},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(ClassSpellMask_DruidWrath | ClassSpellMask_DruidStarsurge) {
				return
			}

			aura.RemoveStack(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Flat,
		ClassMask:  ClassSpellMask_DruidWrath | ClassSpellMask_DruidStarsurge,
		FloatValue: solarCritBonus,
	})

	// Lunar
	druid.LunarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:     "Lunar Eclipse proc",
		Duration:  time.Second * 15,
		MaxStacks: 4,
		ActionID:  core.ActionID{SpellID: 408255},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(ClassSpellMask_DruidStarfire) {
				return
			}

			aura.RemoveStack(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		ClassMask: ClassSpellMask_DruidStarfire,
		TimeValue: -lunarCastTimeReduction,
	})

	druid.EclipseAura = druid.RegisterAura(core.Aura{
		Label:    "Eclipse",
		Duration: core.NeverExpires,
		ActionID: core.ActionID{SpellID: 408248},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(ClassSpellMask_DruidWrath|ClassSpellMask_DruidStarfire|ClassSpellMask_DruidStarsurge) || !result.Landed() {
				return
			}

			if spell.Matches(ClassSpellMask_DruidWrath | ClassSpellMask_DruidStarsurge) {
				druid.LunarEclipseProcAura.Activate(sim)
				// Solar gives 1 stack of lunar bonus
				druid.LunarEclipseProcAura.AddStack(sim)
			}

			if spell.Matches(ClassSpellMask_DruidStarfire | ClassSpellMask_DruidStarsurge) {
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

			switch spell.ClassSpellMask {
			case ClassSpellMask_DruidWrath:
				druid.tryElunesFiresSunfireExtension(sim, result.Target)
			case ClassSpellMask_DruidStarfire:
				druid.tryElunesFiresMoonfireExtension(sim, result.Target)
			case ClassSpellMask_DruidStarsurge: // Starsurge now benefits from the effects of Wrath and Starfire
				druid.tryElunesFiresSunfireExtension(sim, result.Target)
				druid.tryElunesFiresMoonfireExtension(sim, result.Target)
			case ClassSpellMask_DruidShred:
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
	druid.registerMangleBearSpell()
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
			if spell.ProcMask.Matches(core.ProcMaskSpellDamage) && result.DidCrit() || spell.Matches(ClassSpellMask_DruidStarsurge) {
				druid.DreamstateManaRegenAura.Activate(sim)
				dreamstateAuras.Get(result.Target).Activate(sim)
			}
		},
	}))
}
