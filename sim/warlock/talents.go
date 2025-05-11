package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) ApplyTalents() {
	warlock.applyWeaponImbue()

	// Affliction
	warlock.applySuppression()
	warlock.applyNightfall()
	warlock.applyShadowMastery()

	// Demonology
	warlock.applyDemonicEmbrace()
	warlock.applyFelIntellect()
	warlock.registerFelDominationCD()
	warlock.applyFelStamina()
	warlock.applyMasterSummoner()
	warlock.applyMasterDemonologist()
	warlock.applyDemonicSacrifice()
	warlock.applySoulLink()

	// Destruction
	warlock.applyImprovedShadowBolt()
	warlock.applyCataclysm()
	warlock.applyBane()
	warlock.applyDevastation()
	warlock.applyImprovedImmolate()
	warlock.applyRuin()
	warlock.applyEmberstorm()
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

			DamageMultiplier: firestoneMulti,
			ThreatMultiplier: 1,
			BonusCoefficient: spellCoeff,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(damageMin, damageMax)

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

///////////////////////////////////////////////////////////////////////////
//                            Affliction
///////////////////////////////////////////////////////////////////////////

func (warlock *Warlock) applySuppression() {
	if warlock.Talents.Suppression == 0 {
		return
	}

	points := float64(warlock.Talents.Suppression)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(WarlockFlagAffliction) {
			spell.BonusHitRating += 2 * points * core.CritRatingPerCritChance
		}
	})
}

func (warlock *Warlock) applyImprovedDrainLife() {
	if warlock.Talents.ImprovedDrainLife == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_BaseDamageDone_Flat,
		School:    core.SpellSchoolShadow,
		ClassMask: ClassSpellMask_WarlockDrainLife,
		IntValue:  int64(2 * warlock.Talents.ImprovedDrainLife),
	})
}

func (warlock *Warlock) applyNightfall() {
	// This aura can be procced by some item sets without having it talented
	warlock.ShadowTranceAura = warlock.RegisterAura(core.Aura{
		Label:    "Nightfall Shadow Trance",
		ActionID: core.ActionID{SpellID: 17941},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.ShadowCleave {
				spell.CD.Reset()
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Check if the shadowbolt was instant cast and not a normal one
			if (spell.Matches(ClassSpellMask_WarlockShadowBolt) && spell.CurCast.CastTime == 0) || spell.Matches(ClassSpellMask_WarlockShadowCleave) {
				aura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  ClassSpellMask_WarlockShadowBolt,
		FloatValue: -1,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_WarlockShadowCleave,
		IntValue:  100,
	})

	if warlock.Talents.Nightfall <= 0 {
		return
	}

	warlock.nightfallProcChance = 0.02 * float64(warlock.Talents.Nightfall)

	hasSoulSiphonRune := warlock.HasRune(proto.WarlockRune_RuneCloakSoulSiphon)

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "Nightfall Hidden Aura",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.Matches(ClassSpellMask_WarlockCorruption|ClassSpellMask_WarlockDrainLife) || (hasSoulSiphonRune && spell.Matches(ClassSpellMask_WarlockDrainSoul))) && sim.Proc(warlock.nightfallProcChance, "Nightfall") {
				warlock.ShadowTranceAura.Activate(sim)
			}
		},
	}))
}

func (warlock *Warlock) applyShadowMastery() {
	if warlock.Talents.ShadowMastery == 0 {
		return
	}

	// These spells have their base damage modded instead
	baseModClassMasks := ClassSpellMask_WarlockCurseOfAgony | ClassSpellMask_WarlockDeathCoil | ClassSpellMask_WarlockDrainLife | ClassSpellMask_WarlockDrainSoul

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		School:    core.SpellSchoolShadow,
		ClassMask: ClassSpellMask_WarlockAll ^ baseModClassMasks,
		IntValue:  int64(2 * warlock.Talents.ShadowMastery),
	})

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_BaseDamageDone_Flat,
		School:    core.SpellSchoolShadow,
		ClassMask: baseModClassMasks,
		IntValue:  int64(2 * warlock.Talents.ShadowMastery),
	})
}

///////////////////////////////////////////////////////////////////////////
//                            Demonology Talents
///////////////////////////////////////////////////////////////////////////

func (warlock *Warlock) applyDemonicEmbrace() {
	if warlock.Talents.DemonicEmbrace == 0 {
		return
	}

	points := float64(warlock.Talents.DemonicEmbrace)
	warlock.MultiplyStat(stats.Stamina, 1+.03*(points))
	warlock.MultiplyStat(stats.Spirit, 1-.01*(points))
}

func (warlock *Warlock) applyFelIntellect() {
	if warlock.Talents.FelIntellect == 0 {
		return
	}

	multiplier := 1 + 0.03*float64(warlock.Talents.FelIntellect)
	for _, pet := range warlock.BasePets {
		pet.MultiplyStat(stats.Mana, multiplier)
	}
}

func (warlock *Warlock) applyFelStamina() {
	if warlock.Talents.FelStamina == 0 {
		return
	}

	multiplier := 1 + 0.03*float64(warlock.Talents.FelStamina)
	for _, pet := range warlock.BasePets {
		pet.MultiplyStat(stats.Health, multiplier)
	}
}

func (warlock *Warlock) applyMasterSummoner() {
	if warlock.Talents.MasterSummoner == 0 {
		return
	}

	// Use an aura because the summon spells aren't registered by this point
	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "Master Summoner Hidden Aura",
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: ClassSpellMask_WarlockSummons,
		IntValue:  -20 * int64(warlock.Talents.MasterSummoner),
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		ClassMask: ClassSpellMask_WarlockSummons,
		TimeValue: -time.Second * 2 * time.Duration(warlock.Talents.MasterSummoner),
	}))
}

func (warlock *Warlock) applyMasterDemonologist() {
	if warlock.Talents.MasterDemonologist == 0 {
		return
	}

	warlock.masterDemonologistMultiplier = 1

	hasMetaRune := warlock.HasRune(proto.WarlockRune_RuneHandsMetamorphosis)

	points := float64(warlock.Talents.MasterDemonologist)
	damageDealtPointsMultiplier := 0.02 * points
	damageTakenPointsMultiplier := 0.02 * points
	threatPointsMultiplier := core.TernaryFloat64(hasMetaRune, 0.04, -0.04) * points
	resistancePointsMultiplier := 2 * points

	impConfig := core.Aura{
		Label:    "Master Demonologist (Imp)",
		ActionID: core.ActionID{SpellID: 23825, Tag: 1},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= 1 + (threatPointsMultiplier * warlock.masterDemonologistMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= 1 + (threatPointsMultiplier * warlock.masterDemonologistMultiplier)
		},
	}

	voidwalkerConfig := core.Aura{
		Label:    "Master Demonologist (Voidwalker)",
		ActionID: core.ActionID{SpellID: 23825, Tag: 2},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 1 - (damageTakenPointsMultiplier * warlock.masterDemonologistMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 1 - (damageTakenPointsMultiplier * warlock.masterDemonologistMultiplier)
		},
	}

	succubusConfig := core.Aura{
		Label:    "Master Demonologist (Succubus)",
		ActionID: core.ActionID{SpellID: 23825, Tag: 3},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1 + (damageDealtPointsMultiplier * warlock.masterDemonologistMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1 + (damageDealtPointsMultiplier * warlock.masterDemonologistMultiplier)
		},
	}

	felhunterConfig := core.Aura{
		Label:    "Master Demonologist (Felhunter)",
		ActionID: core.ActionID{SpellID: 23825, Tag: 4},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddResistancesDynamic(sim, resistancePointsMultiplier*warlock.masterDemonologistMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddResistancesDynamic(sim, -resistancePointsMultiplier*warlock.masterDemonologistMultiplier)
		},
	}

	for _, pet := range warlock.BasePets {
		pet.ApplyOnPetEnable(func(sim *core.Simulation) {
			if warlock.MasterDemonologistAura != nil {
				warlock.MasterDemonologistAura.Deactivate(sim)
			}
		})
	}

	warlockImpAura := warlock.RegisterAura(impConfig)
	impAura := warlock.Imp.RegisterAura(impConfig)
	warlock.Imp.ApplyOnPetEnable(func(sim *core.Simulation) {
		impAura.Activate(sim)
		warlock.MasterDemonologistAura = warlockImpAura
	})
	warlock.Imp.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
		impAura.Deactivate(sim)
	})

	warlockVoidwalkerAura := warlock.RegisterAura(voidwalkerConfig)
	voidwalkerAura := warlock.Voidwalker.RegisterAura(voidwalkerConfig)
	warlock.Voidwalker.ApplyOnPetEnable(func(sim *core.Simulation) {
		voidwalkerAura.Activate(sim)
		warlock.MasterDemonologistAura = warlockVoidwalkerAura
	})
	warlock.Voidwalker.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
		voidwalkerAura.Deactivate(sim)
	})

	warlockSuccubusAura := warlock.RegisterAura(succubusConfig)
	succubusAura := warlock.Succubus.RegisterAura(succubusConfig)
	warlock.Succubus.ApplyOnPetEnable(func(sim *core.Simulation) {
		succubusAura.Activate(sim)
		warlock.MasterDemonologistAura = warlockSuccubusAura
	})
	warlock.Succubus.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
		succubusAura.Deactivate(sim)
	})

	warlockFelhunterAura := warlock.RegisterAura(felhunterConfig)
	felhunterAura := warlock.Felhunter.RegisterAura(felhunterConfig)
	warlock.Felhunter.ApplyOnPetEnable(func(sim *core.Simulation) {
		felhunterAura.Activate(sim)
		warlock.MasterDemonologistAura = warlockFelhunterAura
	})
	warlock.Felhunter.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
		felhunterAura.Deactivate(sim)
	})

	if warlock.HasRune(proto.WarlockRune_RuneBracerSummonFelguard) {
		felguardConfig := core.Aura{
			Label:    "Master Demonologist (Felguard)",
			ActionID: core.ActionID{SpellID: 23825, Tag: 5},
			Duration: core.NeverExpires,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.ThreatMultiplier *= 1 + (threatPointsMultiplier * warlock.masterDemonologistMultiplier)
				aura.Unit.PseudoStats.DamageTakenMultiplier *= 1 - (damageTakenPointsMultiplier * warlock.masterDemonologistMultiplier)
				aura.Unit.PseudoStats.DamageDealtMultiplier *= 1 + (damageDealtPointsMultiplier * warlock.masterDemonologistMultiplier)
				aura.Unit.AddResistancesDynamic(sim, resistancePointsMultiplier*warlock.masterDemonologistMultiplier)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.ThreatMultiplier /= 1 + (threatPointsMultiplier * warlock.masterDemonologistMultiplier)
				aura.Unit.PseudoStats.DamageTakenMultiplier /= 1 - (damageTakenPointsMultiplier * warlock.masterDemonologistMultiplier)
				aura.Unit.PseudoStats.DamageDealtMultiplier /= 1 + (damageDealtPointsMultiplier * warlock.masterDemonologistMultiplier)
				aura.Unit.AddResistancesDynamic(sim, -resistancePointsMultiplier*warlock.masterDemonologistMultiplier)
			},
		}

		warlockFelguardAura := warlock.RegisterAura(felguardConfig)
		felguardAura := warlock.Felguard.RegisterAura(felguardConfig)
		warlock.Felguard.ApplyOnPetEnable(func(sim *core.Simulation) {
			felguardAura.Activate(sim)
			warlock.MasterDemonologistAura = warlockFelguardAura
		})
		warlock.Felguard.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
			felguardAura.Deactivate(sim)
		})
	}

	for _, pet := range warlock.BasePets {
		pet.ApplyOnPetEnable(func(sim *core.Simulation) {
			warlock.MasterDemonologistAura.Activate(sim)
		})

		pet.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
			if !isSacrifice || !warlock.maintainBuffsOnSacrifice {
				warlock.MasterDemonologistAura.Deactivate(sim)
				warlock.MasterDemonologistAura = nil
			}
		})
	}
}

func (warlock *Warlock) applySoulLink() {
	if !warlock.Talents.SoulLink {
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

	warlock.SoulLinkAura = warlock.RegisterAura(soulLinkConfig)
	for _, pet := range warlock.BasePets {
		pet.SoulLinkAura = pet.RegisterAura(soulLinkConfig)

		oldOnPetDisable := pet.OnPetDisable
		pet.OnPetDisable = func(sim *core.Simulation, isSacrifice bool) {
			oldOnPetDisable(sim, isSacrifice)
			warlock.SoulLinkAura.Deactivate(sim)
			pet.SoulLinkAura.Deactivate(sim)
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
			return warlock.ActivePet != nil
		},

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.2,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warlock.SoulLinkAura.Activate(sim)
			warlock.ActivePet.SoulLinkAura.Activate(sim)
		},
	})
}

func (warlock *Warlock) applyDemonicSacrifice() {
	if !warlock.Talents.DemonicSacrifice {
		return
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

	dsAuras := []*core.Aura{felhunterAura, impAura, succubusAura, voidwalkerAura}
	for idx := range warlock.BasePets {
		pet := warlock.BasePets[idx]

		pet.ApplyOnPetEnable(func(sim *core.Simulation) {
			if !warlock.maintainBuffsOnSacrifice || pet == warlock.SacrificedPet {
				for _, dsAura := range dsAuras {
					dsAura.Deactivate(sim)
				}
			}
		})
	}

	warlock.GetOrRegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_WarlockDemonicSacrifice,
		ActionID:       core.ActionID{SpellID: 18788},
		SpellSchool:    core.SpellSchoolShadow,
		Flags:          core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.ActivePet != nil
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			felhunterAura.Deactivate(sim)
			impAura.Deactivate(sim)
			succubusAura.Deactivate(sim)
			voidwalkerAura.Deactivate(sim)

			switch warlock.ActivePet {
			case warlock.Felguard:
				felhunterAura.Activate(sim)
				impAura.Activate(sim)
				succubusAura.Activate(sim)
				voidwalkerAura.Activate(sim)
			case warlock.Felhunter:
				felhunterAura.Activate(sim)
			case warlock.Imp:
				impAura.Activate(sim)
			case warlock.Succubus:
				succubusAura.Activate(sim)
			case warlock.Voidwalker:
				voidwalkerAura.Activate(sim)
			}

			warlock.changeActivePet(sim, nil, true)
		},
	})
}

///////////////////////////////////////////////////////////////////////////
//                            Destruction Talents
///////////////////////////////////////////////////////////////////////////

func (warlock *Warlock) applyImprovedShadowBolt() {
	hasShadowflameRune := warlock.HasRune(proto.WarlockRune_RuneBootsShadowflame)

	// These debuffs get used by the T2.5 DPS 2p bonus and don't require the ISB talent, so always initialize them
	stackCount := core.TernaryInt32(hasShadowflameRune, core.ISBNumStacksShadowflame, core.ISBNumStacksBase)
	warlock.ImprovedShadowBoltAuras = warlock.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		isbAura := core.ImprovedShadowBoltAura(target, warlock.Talents.ImprovedShadowBolt)
		// Use a wrapper to prevent an external ISB from affecting the warlock's effect count
		return target.RegisterAura(core.Aura{
			Label:     "Improved Shadow Bolt Wrapper",
			Duration:  core.ISBDuration,
			MaxStacks: isbAura.MaxStacks,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.activeEffects[aura.Unit.UnitIndex]++
				isbAura.Activate(sim)
				isbAura.SetStacks(sim, stackCount)
			},
			OnRefresh: func(aura *core.Aura, sim *core.Simulation) {
				isbAura.Activate(sim)
				isbAura.SetStacks(sim, stackCount)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.activeEffects[aura.Unit.UnitIndex]--
			},
		})
	})

	if warlock.Talents.ImprovedShadowBolt == 0 {
		return
	}

	improvedShadowBoltSpellClassMasks := ClassSpellMask_WarlockShadowBolt | ClassSpellMask_WarlockShadowCleave | ClassSpellMask_WarlockShadowflame
	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "ISB Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && result.DidCrit() && spell.Matches(improvedShadowBoltSpellClassMasks) {
				isbAura := warlock.ImprovedShadowBoltAuras.Get(result.Target)
				isbAura.Activate(sim)
				isbAura.SetStacks(sim, stackCount)
			}
		},
	}))
}

func (warlock *Warlock) applyCataclysm() {
	if warlock.Talents.Cataclysm == 0 {
		return
	}

	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(WarlockFlagDestruction) && spell.Cost != nil {
			spell.Cost.Multiplier -= warlock.Talents.Cataclysm
		}
	})
}

func (warlock *Warlock) applyBane() {
	if warlock.Talents.Bane == 0 {
		return
	}

	points := time.Duration(warlock.Talents.Bane)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_WarlockShadowBolt | ClassSpellMask_WarlockImmolate | ClassSpellMask_WarlockShadowflame) {
			spell.DefaultCast.CastTime -= time.Millisecond * 100 * points
		} else if spell.Matches(ClassSpellMask_WarlockSoulFire) {
			spell.DefaultCast.CastTime -= time.Millisecond * 400 * points
		}
	})
}

func (warlock *Warlock) applyDevastation() {
	if warlock.Talents.Devastation == 0 {
		return
	}

	points := float64(warlock.Talents.Devastation)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(WarlockFlagDestruction) {
			spell.BonusCritRating += points * core.CritRatingPerCritChance
		}
	})
}

func (warlock *Warlock) applyImprovedImmolate() {
	if warlock.Talents.ImprovedImmolate == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_WarlockImmolate | ClassSpellMask_WarlockShadowflame,
		Kind:      core.SpellMod_ImpactDamageDone_Flat,
		// @Lucenia: Exclude the periodic hits from Invocation and 4pT3 with the "Treat as Periodic" flag.
		// This bug was present in both the sim and in-game but was confirmed to be unintended and fixed in Phase 7
		SpellFlagsExclude: core.SpellFlagTreatAsPeriodic,
		IntValue:          int64(5 * warlock.Talents.ImprovedImmolate),
	})
}

func (warlock *Warlock) applyRuin() {
	if !warlock.Talents.Ruin {
		return
	}
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(WarlockFlagDestruction) {
			spell.CritDamageBonus += 1
		}
	})
}

func (warlock *Warlock) applyEmberstorm() {
	if warlock.Talents.Emberstorm == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		School:    core.SpellSchoolFire,
		ClassMask: ClassSpellMask_WarlockAll,
		IntValue:  int64(2 * warlock.Talents.Emberstorm),
	})
}
