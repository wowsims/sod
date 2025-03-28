package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (shaman *Shaman) ApplyTalents() {
	// Elemental Talents
	shaman.applyConcussion()
	shaman.applyCallOfFlame()
	shaman.applyElementalFocus()
	shaman.applyElementalDevastation()
	shaman.applyElementalFury()
	shaman.registerElementalMasteryCD()

	// Enhancement Talents
	shaman.applyFlurry()

	if shaman.Talents.AncestralKnowledge > 0 {
		shaman.MultiplyStat(stats.Mana, 1.0+0.01*float64(shaman.Talents.AncestralKnowledge))
	}

	shaman.AddStat(stats.Block, 1*float64(shaman.Talents.ShieldSpecialization))

	shaman.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(shaman.Talents.ThunderingStrikes))

	shaman.AddStat(stats.Dodge, 1*float64(shaman.Talents.Anticipation))

	shaman.ApplyEquipScaling(stats.Armor, 1+.02*float64(shaman.Talents.Toughness))

	if shaman.Talents.Parry {
		shaman.PseudoStats.CanParry = true
	}

	// TODO: Check whether this does what it should.
	// From all I've seen this appears to not actually be a school modifier at all, but instead simply applies
	// to all attacks done with a weapon. The weaponmask seems to take precedence and the school mask is actually ignored.
	// Will also be the case for similar talents like the one for retribution.
	shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + (.02 * float64(shaman.Talents.WeaponMastery))

	// Restoration Talents
	// TODO: Healing Way
	// TODO: Ancestral Healing
	shaman.registerNaturesSwiftnessCD()
	// shaman.registerManaTideTotemCD()

	if shaman.Talents.TidalFocus > 0 {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: ClassSpellMask_ShamanAll,
			ProcMask:  core.ProcMaskSpellHealing,
			Kind:      core.SpellMod_PowerCost_Pct,
			IntValue:  -int64(shaman.Talents.TidalFocus),
		})
	}

	shaman.AddStat(stats.MeleeHit, float64(shaman.Talents.NaturesGuidance))
	shaman.AddStat(stats.SpellHit, float64(shaman.Talents.NaturesGuidance))

	if shaman.Talents.HealingGrace > 0 {
		shaman.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_Threat_Flat,
			ClassMask:  ClassSpellMask_ShamanAll,
			ProcMask:   core.ProcMaskSpellHealing,
			FloatValue: -0.05 * float64(shaman.Talents.HealingGrace),
		})
	}

	if shaman.Talents.TidalMastery > 0 {
		critBonus := float64(shaman.Talents.TidalMastery) * core.CritRatingPerCritChance
		shaman.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Matches(ClassSpellMask_ShamanHealingSpell | ClassSpellMask_ShamanLightningSpell) {
				spell.BonusCritRating += critBonus
			}
		})
	}
}

func (shaman *Shaman) applyConcussion() {
	if shaman.Talents.Concussion == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanLightningBolt | ClassSpellMask_ShamanChainLightning |
			ClassSpellMask_ShamanEarthShock | ClassSpellMask_ShamanFlameShock | ClassSpellMask_ShamanFrostShock,
		Kind:     core.SpellMod_DamageDone_Flat,
		IntValue: int64(1 * shaman.Talents.Concussion),
	})
}

func (shaman *Shaman) applyCallOfFlame() {
	if shaman.Talents.CallOfFlame == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanSearingTotemAttack | ClassSpellMask_ShamanFireNovaTotemAttack | ClassSpellMask_ShamanFireNova,
		Kind:      core.SpellMod_DamageDone_Flat,
		IntValue:  int64(5 * shaman.Talents.CallOfFlame),
	})
}

func (shaman *Shaman) applyElementalFocus() {
	if !shaman.Talents.ElementalFocus {
		return
	}

	shaman.elementalFocusProcChance = 0.1

	costMod := shaman.AddDynamicMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanAll,
		ProcMask:  core.ProcMaskSpellDamage,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
	})

	shaman.ClearcastingAura = shaman.RegisterAura(core.Aura{
		Label:     "Clearcasting",
		ActionID:  core.ActionID{SpellID: 16246},
		Duration:  time.Second * 15,
		MaxStacks: 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if newStacks == 0 {
				aura.Deactivate(sim)
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}

			if aura.GetStacks() > 0 && shaman.isShamanDamagingSpell(spell) {
				aura.RemoveStack(sim)
			}
		},
	})

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: "Elemental Focus Trigger",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if shaman.isShamanDamagingSpell(spell) && sim.Proc(shaman.elementalFocusProcChance, "Elemental Focus") {
				shaman.ClearcastingAura.Activate(sim)
				shaman.ClearcastingAura.SetStacks(sim, shaman.ClearcastingAura.MaxStacks)
			}
		},
	}))
}

func (shaman *Shaman) isShamanDamagingSpell(spell *core.Spell) bool {
	return spell.Matches(ClassSpellMask_ShamanAll) && spell.ProcMask.Matches(core.ProcMaskSpellDamage)
}

func (shaman *Shaman) applyElementalDevastation() {
	if shaman.Talents.ElementalDevastation == 0 {
		return
	}

	spellID := []int32{0, 30165, 29177, 29178}[shaman.Talents.ElementalDevastation]
	critBonus := 3.0 * float64(shaman.Talents.ElementalDevastation) * core.CritRatingPerCritChance
	procAura := shaman.NewTemporaryStatsAura("Elemental Devastation Proc", core.ActionID{SpellID: spellID}, stats.Stats{stats.MeleeCrit: critBonus}, time.Second*10)

	shaman.RegisterAura(core.Aura{
		Label:    "Elemental Devastation",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskSpellDamage) && result.Outcome.Matches(core.OutcomeCrit) {
				procAura.Activate(sim)
			}
		},
	})
}

func (shaman *Shaman) applyElementalFury() {
	if !shaman.Talents.ElementalFury {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:        core.SpellMod_CritDamageBonus_Flat,
		ClassMask:   ClassSpellMask_ShamanAll | ClassSpellMask_ShamanTotems,
		DefenseType: core.DefenseTypeMagic,
		FloatValue:  1,
	})
}

func (shaman *Shaman) registerElementalMasteryCD() {
	if !shaman.Talents.ElementalMastery {
		return
	}

	actionID := core.ActionID{SpellID: 16166}

	cdTimer := shaman.NewTimer()
	cd := time.Minute * 3

	emAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ElementalMastery.CD.Use(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_ShamanHarmfulGCDSpells) && !spell.ProcMask.Matches(core.ProcMaskSpellProc) {
				// Elemental mastery can be batched
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + core.SpellBatchWindow,
					OnAction: func(sim *core.Simulation) {
						if aura.IsActive() {
							// Remove the buff and put skill on CD
							aura.Deactivate(sim)
							cdTimer.Set(sim.CurrentTime + cd)
							shaman.UpdateMajorCooldowns()
						}
					},
				})
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Flat,
		ClassMask:  ClassSpellMask_ShamanHarmfulGCDSpells,
		ProcMask:   core.ProcMaskSpellDamage,
		FloatValue: core.CritRatingPerCritChance * 100,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  ClassSpellMask_ShamanHarmfulGCDSpells,
		ProcMask:   core.ProcMaskSpellDamage,
		FloatValue: -100,
	})

	shaman.ElementalMastery = shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			emAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: shaman.ElementalMastery,
		Type:  core.CooldownTypeDPS,
	})
}

func (shaman *Shaman) registerNaturesSwiftnessCD() {
	if !shaman.Talents.NaturesSwiftness {
		return
	}
	actionID := core.ActionID{SpellID: 16188}
	cdTimer := shaman.NewTimer()
	cd := time.Minute * 3

	nsAura := shaman.RegisterAura(core.Aura{
		Label:    "Natures Swiftness",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolNature) && spell.DefaultCast.CastTime > 0 {
				// Remove the buff and put skill on CD
				aura.Deactivate(sim)
				cdTimer.Set(sim.CurrentTime + cd)
				shaman.UpdateMajorCooldowns()
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		School:     core.SpellSchoolNature,
		FloatValue: -1,
	})

	nsSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			nsAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: nsSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (shaman *Shaman) applyFlurry() {
	if shaman.Talents.Flurry == 0 {
		return
	}

	talentAura := shaman.makeFlurryAura(shaman.Talents.Flurry)

	// This must be registered before the below trigger because in-game a crit weapon swing consumes a stack before the refresh, so you end up with:
	// 3 => 2
	// refresh
	// 2 => 3
	shaman.makeFlurryConsumptionTrigger(talentAura)

	shaman.RegisterAura(core.Aura{
		Label:    "Flurry Proc Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee) && result.Outcome.Matches(core.OutcomeCrit) {
				talentAura.Activate(sim)
				if talentAura.IsActive() {
					talentAura.SetStacks(sim, 3)
				}
				return
			}
		},
	})
}

// These are separated out because of the T1 Shaman Tank 2P that can proc Flurry separately from the talent.
// It triggers the max-rank Flurry aura but with dodge, parry, or block.
func (shaman *Shaman) makeFlurryAura(points int32) *core.Aura {
	if points == 0 {
		return nil
	}

	spellID := []int32{16257, 16277, 16278, 16279, 16280}[points-1]
	attackSpeed := []float64{1.1, 1.15, 1.2, 1.25, 1.3}[points-1]
	label := fmt.Sprintf("Flurry Proc (%d)", spellID)

	if aura := shaman.GetAura(label); aura != nil {
		return aura
	}

	aura := shaman.RegisterAura(core.Aura{
		Label:     label,
		ActionID:  core.ActionID{SpellID: spellID},
		Duration:  core.NeverExpires,
		MaxStacks: 3,
	})

	aura.NewExclusiveEffect("Flurry", true, core.ExclusiveEffect{
		Priority: attackSpeed,
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, attackSpeed+shaman.bonusFlurrySpeed)
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, 1/(attackSpeed+shaman.bonusFlurrySpeed))
		},
	})

	return aura
}

// With the Warden T1 2pc it's possible to have 2 different Flurry auras if using less than 5/5 points in Flurry.
// The two different buffs don't stack whatsoever. Instead the stronger aura takes precedence and each one is only refreshed by the corresponding triggers.
func (shaman *Shaman) makeFlurryConsumptionTrigger(flurryAura *core.Aura) *core.Aura {
	label := fmt.Sprintf("Flurry Consume Trigger - %d", flurryAura.ActionID.SpellID)
	if aura := shaman.GetAura(label); aura != nil {
		return aura
	}

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond * 500,
	}

	return core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Remove a stack.
			if flurryAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && icd.IsReady(sim) {
				icd.Use(sim)
				flurryAura.RemoveStack(sim)
			}
		},
	}))
}

func (shaman *Shaman) totemManaMultiplier() int32 {
	return 100 - 5*shaman.Talents.TotemicFocus
}

// Restorative Totems uses Mod Spell Effectiveness (Base Value)
func (shaman *Shaman) restorativeTotemsModifier() float64 {
	return 0.05 * float64(shaman.Talents.RestorativeTotems)
}

// Purification uses Mod Spell Effectiveness (Base Healing)
func (shaman *Shaman) purificationHealingModifier() float64 {
	return .02 * float64(shaman.Talents.Purification)
}

// func (shaman *Shaman) registerManaTideTotemCD() {
// 	if !shaman.Talents.ManaTideTotem {
// 		return
// 	}

// 	mttAura := core.ManaTideTotemAura(shaman.GetCharacter(), shaman.Index)
// 	mttSpell := shaman.RegisterSpell(core.SpellConfig{
// 		ActionID: core.ManaTideTotemActionID,
// 		Flags:    core.SpellFlagNoOnCastComplete,
// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: time.Second,
// 			},
// 			IgnoreHaste: true,
// 			CD: core.Cooldown{
// 				Timer:    shaman.NewTimer(),
// 				Duration: time.Minute * 5,
// 			},
// 		},
// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
// 			mttAura.Activate(sim)

// 			// If healing stream is active, cancel it while mana tide is up.
// 			if shaman.HealingStreamTotem.Hot(&shaman.Unit).IsActive() {
// 				for _, agent := range shaman.Party.Players {
// 					shaman.HealingStreamTotem.Hot(&agent.GetCharacter().Unit).Cancel(sim)
// 				}
// 			}

// 			// TODO: Current water totem buff needs to be removed from party/raid.
// 			if shaman.Totems.Water != proto.WaterTotem_NoWaterTotem {
// 				shaman.TotemExpirations[WaterTotem] = sim.CurrentTime + time.Second*12
// 			}
// 		},
// 	})

// 	shaman.AddMajorCooldown(core.MajorCooldown{
// 		Spell: mttSpell,
// 		Type:  core.CooldownTypeDPS,
// 		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
// 			return sim.CurrentTime > time.Second*30
// 		},
// 	})
// }
