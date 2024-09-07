package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	DevilsaurEye             = 19991
	DevilsaurTooth           = 19992
	SignetOfBeasts           = 209823
	BloodlashBow             = 216516
	GurubashiPitFightersBow  = 221450
	BloodChainVices          = 227075
	KnightChainVices         = 227077
	BloodChainGrips          = 227081
	KnightChainGrips         = 227087
	WhistleOfTheBeast        = 228432
	ArcaneInfusedGem		 = 230237
	MaelstromsWrath          = 231320
	ZandalarPredatorsMantle  = 231321
	ZandalarPredatorsBelt    = 231322
	ZandalarPredatorsBracers = 231323
	MarshalChainGrips        = 231560
	GeneralChainGrips        = 231569
	GeneralChainVices        = 231575
	MarshalChainVices        = 231578
	Kestrel                  = 231754
	Peregrine                = 231755
)

func applyRaptorStrikeDamageEffect(agent core.Agent, multiplier float64) {
	hunter := agent.(HunterAgent).GetHunter()
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellCode == SpellCode_HunterRaptorStrikeHit {
			spell.DamageMultiplier *= multiplier
		}
	})
}

func applyMultiShotDamageEffect(agent core.Agent, multiplier float64) {
	hunter := agent.(HunterAgent).GetHunter()
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellCode == SpellCode_HunterMultiShot {
			spell.DamageMultiplier *= multiplier
		}
	})
}

func init() {
	core.NewItemEffect(DevilsaurEye, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()

		procBonus := stats.Stats{
			stats.AttackPower:       150,
			stats.RangedAttackPower: 150,
			stats.MeleeHit:          2,
		}
		aura := hunter.GetOrRegisterAura(core.Aura{
			Label:    "Devilsaur Fury",
			ActionID: core.ActionID{SpellID: 24352},
			Duration: time.Second * 20,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatsDynamic(sim, procBonus)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatsDynamic(sim, procBonus.Invert())
			},
		})

		spell := hunter.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 24352},
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    hunter.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		hunter.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(DevilsaurTooth, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		if hunter.pet == nil {
			return
		}

		// Hunter aura so its visible in the timeline
		// TODO: Probably should add pet auras in the timeline at some point
		trackingAura := hunter.GetOrRegisterAura(core.Aura{
			Label:    "Primal Instinct Hunter",
			ActionID: core.ActionID{SpellID: 24353},
			Duration: core.NeverExpires,
		})

		aura := hunter.pet.GetOrRegisterAura(core.Aura{
			Label:    "Primal Instinct",
			ActionID: core.ActionID{SpellID: 24353},
			Duration: core.NeverExpires,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if hunter.pet.focusDump != nil {
					hunter.pet.focusDump.BonusCritRating += 100
				}
				if hunter.pet.specialAbility != nil {
					hunter.pet.specialAbility.BonusCritRating += 100
				}
				trackingAura.Activate(sim)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if hunter.pet.focusDump != nil {
					hunter.pet.focusDump.BonusCritRating -= 100
				}
				if hunter.pet.specialAbility != nil {
					hunter.pet.specialAbility.BonusCritRating -= 100
				}
				trackingAura.Deactivate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == hunter.pet.focusDump || spell == hunter.pet.specialAbility {
					aura.Deactivate(sim)
				}
			},
		})

		spell := hunter.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 24353},
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				CD: core.Cooldown{
					Timer:    hunter.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.pet.IsEnabled()
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		hunter.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				return hunter.pet != nil && hunter.pet.IsEnabled()
			},
		})
	})

	core.NewItemEffect(SignetOfBeasts, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		if hunter.pet != nil {
			hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.01
		}
	})

	core.NewItemEffect(BloodlashBow, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		hunter.newBloodlashProcItem(50, 436471)
	})

	core.NewItemEffect(GurubashiPitFightersBow, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		hunter.newBloodlashProcItem(75, 446723)
	})

	// https://www.wowhead.com/classic/item=228432/whistle-of-the-beast
	// Use: Your pet's next attack is guaranteed to critically strike if that attack is capable of striking critically. (1 Min Cooldown)
	core.NewItemEffect(WhistleOfTheBeast, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()

		if hunter.pet == nil {
			return
		}

		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.03
		hunter.pet.MultiplyStat(stats.Health, 1.03)
		hunter.pet.MultiplyStat(stats.Armor, 1.10)
		hunter.pet.AddStat(stats.MeleeCrit, 2*core.CritRatingPerCritChance)
		hunter.pet.AddStat(stats.SpellCrit, 2*core.SpellCritRatingPerCritChance)

		actionID := core.ActionID{ItemID: WhistleOfTheBeast}

		trackingAura := hunter.GetOrRegisterAura(core.Aura{
			Label:    "Whistle of the Beast Hunter",
			ActionID: actionID,
			Duration: core.NeverExpires,
		})

		aura := hunter.pet.GetOrRegisterAura(core.Aura{
			Label:    "Whistle of the Beast",
			ActionID: actionID,
			Duration: core.NeverExpires,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if hunter.pet.focusDump != nil {
					hunter.pet.focusDump.BonusCritRating += 100
				}
				if hunter.pet.specialAbility != nil {
					hunter.pet.specialAbility.BonusCritRating += 100
				}
				trackingAura.Activate(sim)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if hunter.pet.focusDump != nil {
					hunter.pet.focusDump.BonusCritRating -= 100
				}
				if hunter.pet.specialAbility != nil {
					hunter.pet.specialAbility.BonusCritRating -= 100
				}
				trackingAura.Deactivate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == hunter.pet.focusDump || spell == hunter.pet.specialAbility {
					aura.Deactivate(sim)
				}
			},
		})

		spell := hunter.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    hunter.NewTimer(),
					Duration: time.Minute * 1,
				},
			},
			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return hunter.pet.IsEnabled()
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		hunter.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				return hunter.pet != nil && hunter.pet.IsEnabled()
			},
		})
	})

	core.NewItemEffect(BloodChainGrips, func(agent core.Agent) {
		applyRaptorStrikeDamageEffect(agent, 1.04)
	})

	core.NewItemEffect(KnightChainGrips, func(agent core.Agent) {
		applyRaptorStrikeDamageEffect(agent, 1.04)
	})

	core.NewItemEffect(GeneralChainGrips, func(agent core.Agent) {
		applyRaptorStrikeDamageEffect(agent, 1.04)
	})

	core.NewItemEffect(MarshalChainGrips, func(agent core.Agent) {
		applyRaptorStrikeDamageEffect(agent, 1.04)
	})

	core.NewItemEffect(BloodChainVices, func(agent core.Agent) {
		applyMultiShotDamageEffect(agent, 1.04)
	})

	core.NewItemEffect(KnightChainVices, func(agent core.Agent) {
		applyMultiShotDamageEffect(agent, 1.04)
	})

	core.NewItemEffect(GeneralChainVices, func(agent core.Agent) {
		applyMultiShotDamageEffect(agent, 1.04)
	})

	core.NewItemEffect(MarshalChainVices, func(agent core.Agent) {
		applyMultiShotDamageEffect(agent, 1.04)
	})

	core.NewItemEffect(MaelstromsWrath, func(a core.Agent) {
		hunter := a.(HunterAgent).GetHunter()
		if hunter.pet == nil {
			return
		}

		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.02

		if !hunter.Talents.BestialWrath {
			return
		}

		hunter.RegisterAura(core.Aura{
			Label: "Maelstroms's Wrath Bestial Wrath",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				hunter.BestialWrathPetAura.Duration += (time.Second * 3)
			},
		})
	})

	core.NewItemEffect(ZandalarPredatorsMantle, func(a core.Agent) {
		hunter := a.(HunterAgent).GetHunter()
		if hunter.pet == nil {
			return
		}

		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.03
	})

	core.NewItemEffect(ZandalarPredatorsBelt, func(a core.Agent) {
		hunter := a.(HunterAgent).GetHunter()
		if hunter.pet == nil {
			return
		}

		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.02
	})

	core.NewItemEffect(ZandalarPredatorsBracers, func(a core.Agent) {
		hunter := a.(HunterAgent).GetHunter()

		if hunter.pet == nil {
			return
		}

		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.01
	})

	// https://www.wowhead.com/classic/item=231755/peregrine
	// Chance on hit: Instantly gain 1 extra attack with both weapons.
	// TODO: Proc rate assumed and needs testing
	core.NewItemEffect(Peregrine, func(agent core.Agent) {
		character := agent.GetCharacter()
		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Peregrine Trigger",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskMeleeOH,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			PPM:               1.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				character.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 469140}, spell)
				character.AutoAttacks.ExtraOHAttack(sim, 1, core.ActionID{SpellID: 469140}, spell.ActionID)
			},
		})
	})

	core.NewItemEffect(Kestrel, func(agent core.Agent) {
		character := agent.GetCharacter()

		kestrelAura := character.RegisterAura(core.Aura{
			Label: "Kestrel",
			ActionID: core.ActionID{SpellID: 469148},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MoveSpeed *= 1.40
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MoveSpeed /= 1.40
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Kestrel Trigger",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskMeleeMH,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			PPM:               1.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				kestrelAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(ArcaneInfusedGem, func(agent core.Agent) {
		character := agent.GetCharacter()

		arcaneDetonation := character.RegisterSpell(core.SpellConfig{
			ActionID: 		core.ActionID{SpellID: 467447},
			SpellSchool: 	core.SpellSchoolArcane,
			DefenseType: 	core.DefenseTypeMagic,
			ProcMask: 		core.ProcMaskSpellDamage,
			Flags: 			core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 0,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					damage := sim.Roll(185, 210)
					spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicCrit)
				}
			},
		})

		arcaneInfused := character.RegisterAura(core.Aura{
			Label: "Arcane Infused",
			ActionID: core.ActionID{SpellID: 467446},
			Duration: time.Second * 15,
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				maxTargetsPerCast := int32(5)
				// Uses same targeting code as multi-shot however the detonations occur at cast time rather than when the shots land
				if spell.SpellCode == SpellCode_HunterMultiShot {
					curTarget := sim.Environment.Encounter.TargetUnits[0]
					for hitIndex := int32(0); hitIndex < min(3, character.Env.GetNumTargets()); hitIndex++ {
						arcaneDetonation.Cast(sim, curTarget)
						curTarget = sim.Environment.NextTargetUnit(curTarget)
					}
				}
				// 1 explosion per target up to 5 targets per carve cast
				if spell.SpellCode == SpellCode_HunterCarve {
					curTarget := sim.Environment.Encounter.TargetUnits[0]
					for hitIndex := int32(0); hitIndex < min(maxTargetsPerCast, character.Env.GetNumTargets()); hitIndex++ {
						arcaneDetonation.Cast(sim, curTarget)
						curTarget = sim.Environment.NextTargetUnit(curTarget)
					}
				}
			},
		})

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: arcaneInfused.ActionID.SpellID},
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 90,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				arcaneInfused.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})
}

func (hunter *Hunter) newBloodlashProcItem(bonusStrength float64, spellId int32) {
	procAura := hunter.NewTemporaryStatsAura("Bloodlash", core.ActionID{SpellID: spellId}, stats.Stats{stats.Strength: bonusStrength}, time.Second*15)
	ppm := hunter.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMeleeOrRanged)
	core.MakePermanent(hunter.GetOrRegisterAura(core.Aura{
		Label: "Bloodlash Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && ppm.Proc(sim, spell.ProcMask, "Bloodlash Proc") {
				procAura.Activate(sim)
			}
		},
	}))
}
