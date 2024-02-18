package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

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

func (druid *Druid) ApplyTalents() {
	druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + 0.02*float64(druid.Talents.NaturalWeapons)
	druid.ApplyEquipScaling(stats.Armor, druid.ThickHideMultiplier())

	if druid.Talents.HeartOfTheWild > 0 {
		bonus := 0.04 * float64(druid.Talents.HeartOfTheWild)
		druid.MultiplyStat(stats.Intellect, 1.0+bonus)
	}

	druid.setupNaturesGrace()
	druid.applyMoonkinForm()
	druid.applyOmenOfClarity()
	druid.applyBloodFrenzy()
}

func (druid *Druid) setupNaturesGrace() {
	if !druid.Talents.NaturesGrace {
		return
	}

	ngWrathGCD := time.Millisecond * 1000

	var wrathSpells []*DruidSpell

	druid.NaturesGraceProcAura = druid.RegisterAura(core.Aura{
		Label:    "Natures Grace Proc",
		ActionID: core.ActionID{SpellID: 16886},
		Duration: time.Second * 15,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			wrathSpells = core.FilterSlice(
				core.Flatten([][]*DruidSpell{druid.Wrath}),
				func(spell *DruidSpell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(wrathSpells, func(spell *DruidSpell) {
				spell.Spell.DefaultCast.GCD = ngWrathGCD
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(wrathSpells, func(spell *DruidSpell) {
				spell.Spell.DefaultCast.GCD = core.GCDDefault
			})
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.CastTime() > 0 && aura.TimeActive(sim) >= spell.CastTime() {
				aura.Deactivate(sim)
			}
		},
	})

	druid.RegisterAura(core.Aura{
		Label:    "Natures Grace",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Spells with travel times have their own calculation prior to the travel time to prevent issues with back-to-back casts
			if spell.TravelTime() > 0 && result.DidCrit() {
				druid.NaturesGraceProcAura.Activate(sim)
			}
		},
	})
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

	druid.RegisterAura(core.Aura{
		Label:    "Blood Frenzy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if druid.InForm(Cat) {
				if druid.IsMangle(spell) || druid.Shred.IsEqual(spell) || druid.Rake.IsEqual(spell) {
					if result.Outcome.Matches(core.OutcomeCrit) {
						if sim.Proc(procChance, "Blood Frenzy") {
							druid.AddComboPoints(sim, 1, cpMetrics)
						}
					}
				}
			}
		},
	})
}

// TODO: Class druid omen
func (druid *Druid) applyOmenOfClarity() {

	if !druid.Talents.OmenOfClarity {
		return
	}

	var affectedSpells []*DruidSpell
	druid.ClearcastingAura = druid.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 16870},
		Duration: time.Second * 15,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice([]*DruidSpell{
				druid.DemoralizingRoar,
				druid.FerociousBite,
				druid.Lacerate,
				druid.MangleBear,
				druid.MangleCat,
				druid.Maul,
				druid.Rake,
				druid.Rip,
				druid.Shred,
				druid.SwipeBear,
				druid.SwipeCat,
				druid.SavageRoar,
			}, func(spell *DruidSpell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CostMultiplier -= 1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CostMultiplier += 1
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if aura.RemainingDuration(sim) == aura.Duration {
				// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate
				// if it was just activated.
				return
			}

			for _, as := range affectedSpells {
				if as.IsEqual(spell) {
					aura.Deactivate(sim)
					break
				}
			}
		},
	})

	druid.ProcOoc = func(sim *core.Simulation) {
		druid.ClearcastingAura.Activate(sim)
	}

	ppmm := druid.AutoAttacks.NewPPMManager(2.0, core.ProcMaskMelee)
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
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			if !ppmm.ProcWithWeaponSpecials(sim, spell.ProcMask, "Omen of Clarity") {
				return
			}
			icd.Use(sim)
			druid.ClearcastingAura.Activate(sim)
		},
	})
}

func (druid *Druid) ImprovedMoonfireDamageMultiplier() float64 {
	return 1 + .02*float64(druid.Talents.ImprovedMoonfire)
}

func (druid *Druid) ImprovedMoonfireCritBonus() float64 {
	return 2 * float64(druid.Talents.ImprovedMoonfire)
}

func (druid *Druid) MoonfuryDamageMultiplier() float64 {
	return 1 + 0.02*float64(druid.Talents.Moonfury)
}

func (druid *Druid) VengeanceCritMultiplier() float64 {
	return druid.SpellCritMultiplier(1, 0.2*float64(druid.Talents.Vengeance))
}
