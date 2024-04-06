package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) ApplyTalents() {
	// Demonic Embrace
	if warlock.Talents.DemonicEmbrace > 0 {
		warlock.MultiplyStat(stats.Stamina, 1+.03*(float64(warlock.Talents.DemonicEmbrace)))
		warlock.MultiplyStat(stats.Spirit, 1-.01*(float64(warlock.Talents.DemonicEmbrace)))
	}

	if warlock.Pet != nil && warlock.Talents.FelIntellect > 0 {
		warlock.Pet.MultiplyStat(stats.Mana, 1+0.03*float64(warlock.Talents.FelIntellect))
	}

	if warlock.Talents.ImprovedShadowBolt > 0 {
		warlock.applyImprovedShadowBolt()
	}

	warlock.applyWeaponImbue()
	warlock.applyNightfall()
	warlock.applyMasterDemonologist()
	warlock.applyDemonicSacrifice()
	warlock.applySoulLink()
}

func (warlock *Warlock) applyMasterDemonologist() {
	if warlock.Talents.MasterDemonologist == 0 || warlock.Pet == nil {
		return
	}

	masterDemonologistConfig := core.Aura{
		Label:    "Master Demonologist",
		ActionID: core.ActionID{SpellID: 23825},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, _ *core.Simulation) {
			switch warlock.Options.Summon {
			case proto.WarlockOptions_Imp:
				aura.Unit.PseudoStats.ThreatMultiplier /= 1 + 0.04*float64(warlock.Talents.MasterDemonologist)
			case proto.WarlockOptions_Succubus:
				aura.Unit.PseudoStats.DamageDealtMultiplier *= 1 + 0.02*float64(warlock.Talents.MasterDemonologist)
			case proto.WarlockOptions_Voidwalker:
				aura.Unit.PseudoStats.DamageTakenMultiplier /= 1 + 0.02*float64(warlock.Talents.MasterDemonologist)
			case proto.WarlockOptions_Felguard:
				aura.Unit.PseudoStats.ThreatMultiplier /= 1 + 0.04*float64(warlock.Talents.MasterDemonologist)
				aura.Unit.PseudoStats.DamageDealtMultiplier *= 1 + 0.02*float64(warlock.Talents.MasterDemonologist)
				aura.Unit.PseudoStats.DamageTakenMultiplier /= 1 + 0.02*float64(warlock.Talents.MasterDemonologist)
			}
		},
		OnExpire: func(aura *core.Aura, _ *core.Simulation) {
			switch warlock.Options.Summon {
			case proto.WarlockOptions_Imp:
				aura.Unit.PseudoStats.ThreatMultiplier *= 1 + 0.04*float64(warlock.Talents.MasterDemonologist)
			case proto.WarlockOptions_Succubus:
				aura.Unit.PseudoStats.DamageDealtMultiplier /= 1 + 0.02*float64(warlock.Talents.MasterDemonologist)
			case proto.WarlockOptions_Voidwalker:
				aura.Unit.PseudoStats.DamageTakenMultiplier *= 1 + 0.02*float64(warlock.Talents.MasterDemonologist)
			case proto.WarlockOptions_Felguard:
				aura.Unit.PseudoStats.ThreatMultiplier *= 1 + 0.04*float64(warlock.Talents.MasterDemonologist)
				aura.Unit.PseudoStats.DamageDealtMultiplier /= 1 + 0.02*float64(warlock.Talents.MasterDemonologist)
				aura.Unit.PseudoStats.DamageTakenMultiplier *= 1 + 0.02*float64(warlock.Talents.MasterDemonologist)
			}
		},
	}

	wp := warlock.Pet

	mdLockAura := warlock.RegisterAura(masterDemonologistConfig)
	mdPetAura := wp.RegisterAura(masterDemonologistConfig)

	wp.OnPetEnable = func(sim *core.Simulation) {
		mdLockAura.Activate(sim)
		mdPetAura.Activate(sim)
	}

	wp.OnPetDisable = func(sim *core.Simulation) {
		mdLockAura.Deactivate(sim)
		mdPetAura.Deactivate(sim)
	}
}

func (warlock *Warlock) applySoulLink() {
	if !warlock.Talents.SoulLink || warlock.Pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 19028}
	soulLinkConfig := core.Aura{
		Label:    "Soul Link Aura",
		ActionID: actionID,
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 1.3
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.03
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.03
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 1.3
		},
	}

	warlockSlAura := warlock.RegisterAura(soulLinkConfig)
	petSlAura := warlock.Pet.RegisterAura(soulLinkConfig)

	wp := warlock.Pet
	oldPetDisable := wp.OnPetDisable
	wp.OnPetDisable = func(sim *core.Simulation) {
		if oldPetDisable != nil {
			oldPetDisable(sim)
		}
		if warlockSlAura.IsActive() {
			warlockSlAura.Deactivate(sim)
		}
		if petSlAura.IsActive() {
			petSlAura.Deactivate(sim)
		}
	}

	warlock.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.Pet != nil && warlock.Pet.IsActive()
		},

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.2,
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warlockSlAura.Activate(sim)
			petSlAura.Activate(sim)
		},
	})
}

func (warlock *Warlock) applyDemonicSacrifice() {
	if !warlock.Talents.DemonicSacrifice || warlock.Pet == nil {
		return
	}

	wp := warlock.Pet
	oldPetEnable := wp.OnPetEnable
	wp.OnPetEnable = func(sim *core.Simulation) {
		if oldPetEnable != nil {
			oldPetEnable(sim)
		}
		for _, dsAura := range warlock.demonicSacrificeAuras {
			if dsAura.IsActive() {
				dsAura.Deactivate(sim)
			}
		}
	}

	impAura := warlock.GetOrRegisterAura(core.Aura{
		Label:    "Burning Wish",
		ActionID: core.ActionID{SpellID: 18789},
		Duration: 30 * time.Minute,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= 1.15
		},
	})

	var vwPa *core.PendingAction
	healthMetric := warlock.NewHealthMetrics(core.ActionID{SpellID: 18790})
	voidwalkerAura := warlock.GetOrRegisterAura(core.Aura{
		Label:    "Fel Stamina",
		ActionID: core.ActionID{SpellID: 18790},
		Duration: 30 * time.Minute,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			vwPa = core.NewPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 4,
				OnAction: func(s *core.Simulation) {
					warlock.GainHealth(sim, warlock.MaxHealth()*0.03, healthMetric)
				},
			})
			sim.AddPendingAction(vwPa)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			vwPa.Cancel(sim)
		},
	})

	succubusAura := warlock.GetOrRegisterAura(core.Aura{
		Label:    "Touch of Shadow",
		ActionID: core.ActionID{SpellID: 18791},
		Duration: 30 * time.Minute,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.15
		},
	})

	var fhPa *core.PendingAction
	manaMetric := warlock.NewManaMetrics(core.ActionID{SpellID: 18792})
	felhunterAura := warlock.GetOrRegisterAura(core.Aura{
		Label:    "Fel Energy",
		ActionID: core.ActionID{SpellID: 18792},
		Duration: 30 * time.Minute,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			fhPa = core.NewPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 4,
				OnAction: func(s *core.Simulation) {
					warlock.AddMana(sim, warlock.MaxMana()*0.02, manaMetric)
				},
			})
			sim.AddPendingAction(fhPa)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			fhPa.Cancel(sim)
		},
	})

	warlock.demonicSacrificeAuras = make([]*core.Aura, 0)
	switch warlock.Options.Summon {
	case proto.WarlockOptions_Imp:
		warlock.demonicSacrificeAuras = append(warlock.demonicSacrificeAuras, impAura)
	case proto.WarlockOptions_Voidwalker:
		warlock.demonicSacrificeAuras = append(warlock.demonicSacrificeAuras, voidwalkerAura)
	case proto.WarlockOptions_Succubus:
		warlock.demonicSacrificeAuras = append(warlock.demonicSacrificeAuras, succubusAura)
	case proto.WarlockOptions_Felhunter:
		warlock.demonicSacrificeAuras = append(warlock.demonicSacrificeAuras, felhunterAura)
	case proto.WarlockOptions_Felguard:
		warlock.demonicSacrificeAuras = append(warlock.demonicSacrificeAuras, impAura)
		warlock.demonicSacrificeAuras = append(warlock.demonicSacrificeAuras, voidwalkerAura)
		warlock.demonicSacrificeAuras = append(warlock.demonicSacrificeAuras, succubusAura)
		warlock.demonicSacrificeAuras = append(warlock.demonicSacrificeAuras, felhunterAura)
	}

	warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 18788},
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.Pet != nil && warlock.Pet.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, dsAura := range warlock.demonicSacrificeAuras {
				if dsAura != nil {
					dsAura.Activate(sim)
					warlock.Pet.Disable(sim)
				}
			}
		},
	})
}

func (warlock *Warlock) applyImprovedShadowBolt() {
	warlock.ImprovedShadowBoltAuras = warlock.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		return core.ImprovedShadowBoltAura(unit, warlock.Talents.ImprovedShadowBolt)
	})
}

func (warlock *Warlock) applyWeaponImbue() {
	if warlock.GetCharacter().Equipment.OffHand().Type != proto.ItemType_ItemTypeUnknown {
		return
	}

	level := warlock.Level
	if warlock.Options.WeaponImbue == proto.WarlockOptions_Firestone {
		warlock.applyFirestone()
	}
	if warlock.Options.WeaponImbue == proto.WarlockOptions_Spellstone {
		if level >= 55 {
			warlock.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
		}
	}
}

func (warlock *Warlock) applyFirestone() {
	level := warlock.Level

	damageMin := 0.0
	damageMax := 0.0

	// TODO: Test for spell scaling
	spellCoeff := 0.0
	spellId := int32(0)

	// TODO: Test PPM
	ppm := warlock.AutoAttacks.NewPPMManager(8, core.ProcMaskMelee)

	firestoneMulti := 1.0 + float64(warlock.Talents.ImprovedFirestone)*0.15

	if level >= 56 {
		warlock.AddStat(stats.FirePower, 21*firestoneMulti)
		damageMin = 80.0
		damageMax = 120.0
		spellId = 17949
	} else if level >= 46 {
		warlock.AddStat(stats.FirePower, 17*firestoneMulti)
		damageMin = 60.0
		damageMax = 90.0
		spellId = 17947
	} else if level >= 36 {
		warlock.AddStat(stats.FirePower, 14*firestoneMulti)
		damageMin = 40.0
		damageMax = 60.0
		spellId = 17945
	} else if level >= 28 {
		warlock.AddStat(stats.FirePower, 10*firestoneMulti)
		damageMin = 25.0
		damageMax = 35.0
		spellId = 758
	}

	if level >= 28 && warlock.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueUnknown {
		fireProcSpell := warlock.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellId},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier:         firestoneMulti,
			ThreatMultiplier:         1,
			DamageMultiplierAdditive: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(damageMin, damageMax) + spellCoeff*spell.SpellDamage()

				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
			},
		})

		core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
			Label: "Firestone Proc",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if !ppm.Proc(sim, core.ProcMaskMelee, "Firestone Proc") {
					return
				}

				fireProcSpell.Cast(sim, result.Target)
			},
		}))
	}
}

func (warlock *Warlock) applyNightfall() {
	if warlock.Talents.Nightfall <= 0 {
		return
	}

	nightfallProcChance := 0.02 * float64(warlock.Talents.Nightfall)

	warlock.NightfallProcAura = warlock.RegisterAura(core.Aura{
		Label:    "Nightfall Shadow Trance",
		ActionID: core.ActionID{SpellID: 17941},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ShadowBolt.CastTimeMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ShadowBolt.CastTimeMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Check if the shadowbolt was instant cast and not a normal one
			if spell == warlock.ShadowBolt && spell.CurCast.CastTime == 0 {
				aura.Deactivate(sim)
			}
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Nightfall Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warlock.Corruption || spell == warlock.DrainLife {
				if sim.Proc(nightfallProcChance, "Nightfall") {
					warlock.NightfallProcAura.Activate(sim)

					for _, spell := range warlock.ShadowCleave {
						spell.CD.Reset()
					}
				}
			}
		},
	})
}

func (warlock *Warlock) ruin() float64 {
	return core.TernaryFloat64(warlock.Talents.Ruin, 1, 0)
}

// func (warlock *Warlock) setupPyroclasm() {
// 	if warlock.Talents.Pyroclasm <= 0 {
// 		return
// 	}

// 	pyroclasmDamageBonus := 1 + 0.02*float64(warlock.Talents.Pyroclasm)

// 	warlock.PyroclasmAura = warlock.RegisterAura(core.Aura{
// 		Label:    "Pyroclasm",
// 		ActionID: core.ActionID{SpellID: 63244},
// 		Duration: time.Second * 10,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= pyroclasmDamageBonus
// 			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= pyroclasmDamageBonus
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= pyroclasmDamageBonus
// 			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= pyroclasmDamageBonus
// 		},
// 	})

// 	warlock.RegisterAura(core.Aura{
// 		Label:    "Pyroclasm Talent Hidden Aura",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if (spell == warlock.Conflagrate || spell == warlock.SearingPain) && result.DidCrit() {
// 				warlock.PyroclasmAura.Activate(sim)
// 			}
// 		},
// 	})
// }
