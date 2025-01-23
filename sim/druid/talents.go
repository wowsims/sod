package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (druid *Druid) ApplyTalents() {
	// Balance
	druid.registerMoonkinFormSpell()
	druid.applyOmenOfClarity()
	druid.applyImprovedMoonfire()
	druid.applyVengeance()
	druid.applyNaturesGrace()
	druid.applyMoonglow()
	druid.applyMoonfury()

	// SoD tuning made this all damage, not just physical damage
	druid.PseudoStats.DamageDealtMultiplier *= 1 + 0.02*float64(druid.Talents.NaturalWeapons)

	// Feral
	druid.applyBloodFrenzy()

	druid.ApplyEquipScaling(stats.Armor, druid.ThickHideMultiplier())

	if druid.Talents.HeartOfTheWild > 0 {
		bonus := 0.04 * float64(druid.Talents.HeartOfTheWild)
		druid.MultiplyStat(stats.Intellect, 1.0+bonus)
	}

	// Restoration
	druid.applyFuror()

	druid.PseudoStats.SpiritRegenRateCasting += .05 * float64(druid.Talents.Reflection)
}

func (druid *Druid) ThickHideMultiplier() float64 {
	thickHideMulti := 1.0

	if druid.Talents.ThickHide > 0 {
		thickHideMulti += 0.04 + 0.03*float64(druid.Talents.ThickHide-1)
	}

	return thickHideMulti
}

func (druid *Druid) BearArmorMultiplier() float64 {
	sotfMulti := 1.0 + 0.33/3.0
	return 4.7 * sotfMulti
}

func (druid *Druid) applyNaturesGrace() {
	if !druid.Talents.NaturesGrace {
		return
	}

	affectedSpells := []*DruidSpell{}
	druid.NaturesGraceProcAura = druid.RegisterAura(core.Aura{
		Label:     "Natures Grace Proc",
		ActionID:  core.ActionID{SpellID: 16886},
		Duration:  time.Second * 15,
		MaxStacks: 1,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(druid.DruidSpells, func(ds *DruidSpell) bool {
				return ds.DefaultCast.CastTime > 0
			})
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DefaultCast.CastTime -= time.Millisecond * 500

				if spell.Matches(ClassSpellMask_DruidWrath) {
					spell.DefaultCast.GCD -= time.Millisecond * 500
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DefaultCast.CastTime += time.Millisecond * 500

				if spell.Matches(ClassSpellMask_DruidWrath) {
					spell.DefaultCast.GCD += time.Millisecond * 500
				}
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}

			// Make sure the aura actually applied to the spell being cast before deactivating
			if spell.CurCast.CastTime > 0 && (sim.CurrentTime-spell.CurCast.CastTime >= aura.StartedAt()) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: "Natures Grace",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Spells with travel times have their own implementation because the proc occurs as the cast finishes
			if spell.MissileSpeed == 0 && spell.ProcMask.Matches(core.ProcMaskSpellDamage) && result.DidCrit() {
				druid.NaturesGraceProcAura.Activate(sim)
				druid.NaturesGraceProcAura.SetStacks(sim, druid.NaturesGraceProcAura.MaxStacks)
			}
		},
	}))
}

// func (druid *Druid) registerNaturesSwiftnessCD() {
// 	if !druid.Talents.NaturesSwiftness {
// 		return
// 	}
// 	actionID := core.ActionID{SpellID: 17116}

// 	var nsAura *core.Aura
// 	nsSpell := druid.RegisterSpell(Humanoid|Moonkin|Tree, core.SpellConfig{
// 		ActionID: actionID,
// 		Flags:    core.SpellFlagNoOnCastComplete,
// 		Cast: core.CastConfig{
// 			CD: core.Cooldown{
// 				Timer:    druid.NewTimer(),
// 				Duration: time.Minute * 3,
// 			},
// 		},
// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
// 			nsAura.Activate(sim)
// 		},
// 	})

// 	nsAura = druid.RegisterAura(core.Aura{
// 		Label:    "Natures Swiftness",
// 		ActionID: actionID,
// 		Duration: core.NeverExpires,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			if druid.Starfire != nil {
// 				druid.Starfire.CastTimeMultiplier -= 1
// 			}
// 			if druid.Wrath != nil {
// 				druid.Wrath.CastTimeMultiplier -= 1
// 			}
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			if druid.Starfire != nil {
// 				druid.Starfire.CastTimeMultiplier += 1
// 			}
// 			if druid.Wrath != nil {
// 				druid.Wrath.CastTimeMultiplier += 1
// 			}
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if !druid.Wrath.IsEqual(spell) && !druid.Starfire.IsEqual(spell) {
// 				return
// 			}

// 			// Remove the buff and put skill on CD
// 			aura.Deactivate(sim)
// 			nsSpell.CD.Use(sim)
// 			druid.UpdateMajorCooldowns()
// 		},
// 	})

// 	druid.AddMajorCooldown(core.MajorCooldown{
// 		Spell: nsSpell.Spell,
// 		Type:  core.CooldownTypeDPS,
// 		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
// 			// Don't use NS unless we're casting a full-length starfire or wrath.
// 			return !character.HasTemporarySpellCastSpeedIncrease()
// 		},
// 	})
// }

// TODO: Classic bear
// func (druid *Druid) applyPrimalFury() {
// 	if druid.Talents.PrimalFury == 0 {
// 		return
// 	}

// 	procChance := []float64{0, 0.5, 1}[druid.Talents.PrimalFury]
// 	actionID := core.ActionID{SpellID: 37117}
// 	rageMetrics := druid.NewRageMetrics(actionID)
// 	cpMetrics := druid.NewComboPointMetrics(actionID)

// 	druid.RegisterAura(core.Aura{
// 		Label:    "Primal Fury",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if druid.InForm(Bear) {
// 				if result.Outcome.Matches(core.OutcomeCrit) {
// 					if sim.Proc(procChance, "Primal Fury") {
// 						druid.AddRage(sim, 5, rageMetrics)
// 					}
// 				}
// 			} else if druid.InForm(Cat) {
// 				if druid.IsMangle(spell) || druid.Shred.IsEqual(spell) || druid.Rake.IsEqual(spell) {
// 					if result.Outcome.Matches(core.OutcomeCrit) {
// 						if sim.Proc(procChance, "Primal Fury") {
// 							druid.AddComboPoints(sim, 1, cpMetrics)
// 						}
// 					}
// 				}
// 			}
// 		},
// 	})
// }

func (druid *Druid) applyBloodFrenzy() {
	if druid.Talents.BloodFrenzy == 0 {
		return
	}

	procChance := []float64{0, 0.5, 1}[druid.Talents.BloodFrenzy]
	actionID := core.ActionID{SpellID: 16953}
	cpMetrics := druid.NewComboPointMetrics(actionID)

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: "Blood Frenzy",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if druid.InForm(Cat) &&
				result.Target == aura.Unit.CurrentTarget &&
				spell.Flags.Matches(SpellFlagBuilder) &&
				result.Outcome.Matches(core.OutcomeCrit) &&
				sim.Proc(procChance, "Blood Frenzy") {
				druid.AddComboPoints(sim, 1, result.Target, cpMetrics)
			}
		},
	}))
}

// We're using an aura so that the APL can know if the Druid has furor for powershifting logic
func (druid *Druid) applyFuror() {
	if druid.Talents.Furor == 0 {
		return
	}

	spellID := []int32{0, 17056, 17058, 17059, 17060, 17061}[druid.Talents.Furor]

	druid.FurorAura = druid.RegisterAura(core.Aura{
		Label:    "Furor",
		ActionID: core.ActionID{SpellID: spellID},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (druid *Druid) applyOmenOfClarity() {
	if !druid.Talents.OmenOfClarity {
		return
	}

	var affectedSpells []*core.Spell
	druid.ClearcastingAura = druid.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 16870},
		Duration: time.Second * 15,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(druid.Spellbook, func(spell *core.Spell) bool { return spell.Flags.Matches(SpellFlagOmen) })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.Cost.Multiplier -= 100
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.Cost.Multiplier += 100
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}

			// Hotfix 2024-04-13 Starsurge does not consume clearcasting
			if spell.Flags.Matches(SpellFlagOmen) && !spell.Matches(ClassSpellMask_DruidStarsurge) && spell.DefaultCast.Cost > 0 {
				aura.Deactivate(sim)
			}
		},
	})

	dpm := druid.AutoAttacks.NewPPMManager(2.0, core.ProcMaskMelee)
	icd := core.Cooldown{
		Timer:    druid.NewTimer(),
		Duration: time.Second * 10,
	}

	druid.RegisterAura(core.Aura{
		Label:    "Omen of Clarity",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !icd.IsReady(sim) {
				return
			}
			// TODO: Phase 3 "and non-instant spell casts" but we need to find out how the procs work for those
			if spell.ProcMask.Matches(core.ProcMaskMelee) && dpm.ProcWithWeaponSpecials(sim, spell.ProcMask, "Omen of Clarity") {
				icd.Use(sim)
				druid.ClearcastingAura.Activate(sim)
			}
		},
	})
}

func (druid *Druid) applyMoonfury() {
	if druid.Talents.Moonfury == 0 {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: "Moonfury",
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  ClassSpellMask_DruidWrath | ClassSpellMask_DruidStarfire | ClassSpellMask_DruidMoonfire | ClassSpellMask_DruidStarsurge | ClassSpellMask_DruidSunfire,
		FloatValue: 0.02 * float64(druid.Talents.Moonfury),
	}))
}

func (druid *Druid) applyImprovedMoonfire() {
	if druid.Talents.ImprovedMoonfire == 0 {
		return
	}

	affectedSpells := ClassSpellMask_DruidMoonfire | ClassSpellMask_DruidSunfire | ClassSpellMask_DruidStarfallTick | ClassSpellMask_DruidStarfallSplash

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: "Improved moonfire",
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_BaseDamageDone_Flat,
		ClassMask:  affectedSpells,
		FloatValue: 0.02 * float64(druid.Talents.ImprovedMoonfire),
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Flat,
		ClassMask:  affectedSpells,
		FloatValue: 2 * float64(druid.Talents.ImprovedMoonfire) * core.SpellCritRatingPerCritChance,
	}))
}

func (druid *Druid) applyVengeance() {
	if druid.Talents.Vengeance == 0 {
		return
	}

	critDamageBonus := 0.20 * float64(druid.Talents.Vengeance)

	druid.RegisterAura(core.Aura{
		Label: "Vengeance",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells := core.FilterSlice(
				core.Flatten(
					[][]*DruidSpell{
						druid.Wrath,
						druid.Starfire,
						druid.Moonfire,
						{druid.Starsurge},
						{druid.Sunfire},
						{druid.StarfallTick},
						{druid.StarfallSplash},
					},
				),
				func(spell *DruidSpell) bool { return spell != nil },
			)

			for _, spell := range affectedSpells {
				spell.CritDamageBonus += critDamageBonus
			}
		},
	})
}

func (druid *Druid) applyMoonglow() {
	if druid.Talents.Moonglow == 0 {
		return
	}

	multiplier := 3 * druid.Talents.Moonglow

	druid.RegisterAura(core.Aura{
		Label: "Moonglow",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells := core.FilterSlice(
				core.Flatten(
					[][]*DruidSpell{
						druid.Wrath,
						druid.Starfire,
						druid.Moonfire,
						{druid.Starsurge},
						{druid.Starfall},
					},
				),
				func(spell *DruidSpell) bool { return spell != nil },
			)

			for _, spell := range affectedSpells {
				spell.Cost.Multiplier -= multiplier
			}
		},
	})
}
