package warlock

import (
	"slices"
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

			DamageMultiplier:         firestoneMulti,
			ThreatMultiplier:         1,
			DamageMultiplierAdditive: 1,
			BonusCoefficient:         spellCoeff,

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

func (warlock *Warlock) applyNightfall() {
	if warlock.Talents.Nightfall <= 0 {
		return
	}

	hasSoulSiphonRune := warlock.HasRune(proto.WarlockRune_RuneCloakSoulSiphon)

	warlock.nightfallProcChance = 0.02 * float64(warlock.Talents.Nightfall)

	warlock.NightfallProcAura = warlock.RegisterAura(core.Aura{
		Label:    "Nightfall Shadow Trance",
		ActionID: core.ActionID{SpellID: 17941},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.ShadowBolt {
				spell.CastTimeMultiplier -= 1
			}
			for _, spell := range warlock.ShadowCleave {
				spell.CD.Reset()
				spell.DamageMultiplier *= 2
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.ShadowBolt {
				spell.CastTimeMultiplier += 1
			}
			for _, spell := range warlock.ShadowCleave {
				spell.DamageMultiplier /= 2
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Check if the shadowbolt was instant cast and not a normal one
			if (spell.SpellCode == SpellCode_WarlockShadowBolt && spell.CurCast.CastTime == 0) || spell.SpellCode == SpellCode_WarlockShadowCleave {
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
			if (spell.SpellCode == SpellCode_WarlockCorruption || spell.SpellCode == SpellCode_WarlockDrainLife || (hasSoulSiphonRune && spell.SpellCode == SpellCode_WarlockDrainSoul)) && sim.Proc(warlock.nightfallProcChance, "Nightfall") {
				warlock.NightfallProcAura.Activate(sim)
			}
		},
	})
}

func (warlock *Warlock) applyShadowMastery() {
	if warlock.Talents.ShadowMastery == 0 {
		return
	}

	warlock.OnSpellRegistered(func(spell *core.Spell) {
		// Shadow Mastery applies a base damage modifier to all dots / channeled spells instead
		if spell.SpellSchool.Matches(core.SpellSchoolShadow) && isWarlockSpell(spell) && !spell.Flags.Matches(core.SpellFlagPureDot) && !spell.Flags.Matches(WarlockFlagHaunt) {
			spell.DamageMultiplierAdditive += warlock.shadowMasteryBonus()
		}
	})
}

func (warlock *Warlock) shadowMasteryBonus() float64 {
	return .02 * float64(warlock.Talents.ShadowMastery)
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

	castTimeReduction := time.Second * 2 * time.Duration(warlock.Talents.MasterSummoner)
	costReduction := 20 * warlock.Talents.MasterSummoner

	// Use an aura because the summon spells aren't registered by this point
	warlock.RegisterAura(core.Aura{
		Label:    "Master Summoner Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.SummonDemonSpells {
				spell.DefaultCast.CastTime -= castTimeReduction
				spell.CostValues.Multiplier -= costReduction
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.SummonDemonSpells {
				spell.DefaultCast.CastTime += castTimeReduction
				spell.CostValues.Multiplier += costReduction
			}
		},
	})
}

func (warlock *Warlock) applyMasterDemonologist() {
	if warlock.Talents.MasterDemonologist == 0 {
		return
	}

	hasMetaRune := warlock.HasRune(proto.WarlockRune_RuneHandsMetamorphosis)

	points := float64(warlock.Talents.MasterDemonologist)
	damageDealtMultiplier := 1 + 0.02*points
	damageTakenMultiplier := 1 - 0.02*points
	threatMultiplier := core.TernaryFloat64(hasMetaRune, 1+0.04*points, 1-0.04*points)
	bonusResistance := 2 * points

	masterDemonologistConfig := core.Aura{
		Label:    "Master Demonologist",
		ActionID: core.ActionID{SpellID: 23825},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			switch warlock.ActivePet {
			case warlock.Felguard:
				aura.Unit.PseudoStats.DamageDealtMultiplier *= damageDealtMultiplier
				aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMultiplier
				aura.Unit.PseudoStats.ThreatMultiplier *= threatMultiplier
				aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.ArcaneResistance: bonusResistance,
					stats.FireResistance:   bonusResistance,
					stats.FrostPower:       bonusResistance,
					stats.NatureResistance: bonusResistance,
					stats.ShadowResistance: bonusResistance,
				})
			case warlock.Felhunter:
				aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.ArcaneResistance: bonusResistance,
					stats.FireResistance:   bonusResistance,
					stats.FrostPower:       bonusResistance,
					stats.NatureResistance: bonusResistance,
					stats.ShadowResistance: bonusResistance,
				})
			case warlock.Imp:
				aura.Unit.PseudoStats.ThreatMultiplier *= threatMultiplier
			case warlock.Succubus:
				aura.Unit.PseudoStats.DamageDealtMultiplier *= damageDealtMultiplier
			case warlock.Voidwalker:
				aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMultiplier
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			switch warlock.ActivePet {
			case warlock.Felguard:
				aura.Unit.PseudoStats.DamageDealtMultiplier /= damageDealtMultiplier
				aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMultiplier
				aura.Unit.PseudoStats.ThreatMultiplier /= threatMultiplier
				aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.ArcaneResistance: -bonusResistance,
					stats.FireResistance:   -bonusResistance,
					stats.FrostPower:       -bonusResistance,
					stats.NatureResistance: -bonusResistance,
					stats.ShadowResistance: -bonusResistance,
				})
			case warlock.Felhunter:
				aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.ArcaneResistance: -bonusResistance,
					stats.FireResistance:   -bonusResistance,
					stats.FrostPower:       -bonusResistance,
					stats.NatureResistance: -bonusResistance,
					stats.ShadowResistance: -bonusResistance,
				})
			case warlock.Imp:
				aura.Unit.PseudoStats.ThreatMultiplier /= threatMultiplier
			case warlock.Succubus:
				aura.Unit.PseudoStats.DamageDealtMultiplier /= damageDealtMultiplier
			case warlock.Voidwalker:
				aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMultiplier
			}
		},
	}

	lockAura := warlock.RegisterAura(masterDemonologistConfig)
	for _, pet := range warlock.BasePets {
		petAura := pet.RegisterAura(masterDemonologistConfig)

		oldOnPetEnable := pet.OnPetEnable
		pet.OnPetEnable = func(sim *core.Simulation) {
			oldOnPetEnable(sim)
			lockAura.Activate(sim)
			petAura.Activate(sim)
		}

		oldOnPetDisable := pet.OnPetDisable
		pet.OnPetDisable = func(sim *core.Simulation) {
			oldOnPetDisable(sim)
			lockAura.Deactivate(sim)
			petAura.Deactivate(sim)
		}
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
		pet.OnPetDisable = func(sim *core.Simulation) {
			oldOnPetDisable(sim)
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
	for _, pet := range warlock.BasePets {
		oldOnPetEnable := pet.OnPetEnable
		pet.OnPetEnable = func(sim *core.Simulation) {
			oldOnPetEnable(sim)
			for _, dsAura := range dsAuras {
				dsAura.Deactivate(sim)
			}
		}
	}

	warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 18788},
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.ActivePet != nil
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
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

			warlock.changeActivePet(sim, nil)
		},
	})
}

///////////////////////////////////////////////////////////////////////////
//                            Destruction Talents
///////////////////////////////////////////////////////////////////////////

func (warlock *Warlock) applyImprovedShadowBolt() {
	if warlock.Talents.ImprovedShadowBolt == 0 {
		return
	}

	hasShadowflameRune := warlock.HasRune(proto.WarlockRune_RuneBootsShadowflame)

	stackCount := core.TernaryInt32(hasShadowflameRune, core.ISBNumStacksShadowflame, core.ISBNumStacksBase)
	warlock.ImprovedShadowBoltAuras = warlock.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		return core.ImprovedShadowBoltAura(unit, warlock.Talents.ImprovedShadowBolt, stackCount)
	})

	affectedSpellCodes := []int32{SpellCode_WarlockShadowBolt, SpellCode_WarlockShadowCleave, SpellCode_WarlockShadowflame}
	warlock.RegisterAura(core.Aura{
		Label:    "ISB Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && result.DidCrit() && slices.Contains(affectedSpellCodes, spell.SpellCode) {
				impShadowBoltAura := warlock.ImprovedShadowBoltAuras.Get(result.Target)
				impShadowBoltAura.Activate(sim)
				impShadowBoltAura.SetStacks(sim, stackCount)
			}
		},
	})
}

func (warlock *Warlock) applyCataclysm() {
	if warlock.Talents.Cataclysm == 0 {
		return
	}

	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(WarlockFlagDestruction) {
			spell.CostValues.Multiplier -= warlock.Talents.Cataclysm
		}
	})
}

func (warlock *Warlock) applyBane() {
	if warlock.Talents.Bane == 0 {
		return
	}

	points := time.Duration(warlock.Talents.Bane)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellCode == SpellCode_WarlockShadowBolt || spell.SpellCode == SpellCode_WarlockImmolate || spell.SpellCode == SpellCode_WarlockShadowflame {
			spell.DefaultCast.CastTime -= time.Millisecond * 100 * points
		} else if spell.SpellCode == SpellCode_WarlockSoulFire {
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

func (warlock *Warlock) improvedImmolateBonus() float64 {
	return 0.05 * float64(warlock.Talents.ImprovedImmolate)
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

	points := float64(warlock.Talents.Emberstorm)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellSchool.Matches(core.SpellSchoolFire) && isWarlockSpell(spell) {
			spell.DamageMultiplierAdditive += 0.02 * points
		}
	})
}
