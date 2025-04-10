package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
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
	ArcaneInfusedGem         = 230237
	RenatakisCharmOfRavaging = 231288
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
	CloakOfTheUnseenPath     = 233420
	ScytheOfTheUnseenPath    = 233421
	SignetOfTheUnseenPath    = 233422
	StringsOfFate            = 240837
	PoleaxeOfTheBeast        = 240924
)

func applyRaptorStrikeDamageEffect(agent core.Agent, modifier int64) {
	hunter := agent.(HunterAgent).GetHunter()

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterRaptorStrikeHit,
		IntValue:  modifier,
	})
}

func applyMultiShotDamageEffect(agent core.Agent, modifier int64) {
	hunter := agent.(HunterAgent).GetHunter()

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterMultiShot,
		IntValue:  modifier,
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

		hunter.AddMajorEquipmentCooldown(core.MajorCooldown{
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

		hunter.AddMajorEquipmentCooldown(core.MajorCooldown{
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
			aura := core.MakePermanent(hunter.GetOrRegisterAura(core.Aura{
				Label: "Increased Hunter Pet Damage (Signet of Beasts)",
			}).AttachAdditivePseudoStatBuff(&hunter.pet.PseudoStats.DamageDealtMultiplierAdditive, 0.01))

			hunter.ItemSwap.RegisterProc(SignetOfBeasts, aura)
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

		hunter.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				return hunter.pet != nil && hunter.pet.IsEnabled()
			},
		})
	})

	core.NewItemEffect(BloodChainGrips, func(agent core.Agent) {
		applyRaptorStrikeDamageEffect(agent, 4)
	})

	core.NewItemEffect(KnightChainGrips, func(agent core.Agent) {
		applyRaptorStrikeDamageEffect(agent, 4)
	})

	core.NewItemEffect(GeneralChainGrips, func(agent core.Agent) {
		applyRaptorStrikeDamageEffect(agent, 4)
	})

	core.NewItemEffect(MarshalChainGrips, func(agent core.Agent) {
		applyRaptorStrikeDamageEffect(agent, 4)
	})

	core.NewItemEffect(BloodChainVices, func(agent core.Agent) {
		applyMultiShotDamageEffect(agent, 4)
	})

	core.NewItemEffect(KnightChainVices, func(agent core.Agent) {
		applyMultiShotDamageEffect(agent, 4)
	})

	core.NewItemEffect(GeneralChainVices, func(agent core.Agent) {
		applyMultiShotDamageEffect(agent, 4)
	})

	core.NewItemEffect(MarshalChainVices, func(agent core.Agent) {
		applyMultiShotDamageEffect(agent, 4)
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
	// Main-hand attack is treated like a normal extra-attack, Off-hand attack is a spell that uses your off-hand damage but won't glance
	core.NewItemEffect(Peregrine, func(agent core.Agent) {
		character := agent.GetCharacter()
		peregrineOHAttack := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 469140},
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeOHSpecial,
			Flags:       core.SpellFlagMeleeMetrics,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				damage := character.OHWeaponDamage(sim, spell.MeleeAttackPower()) * character.AutoAttacks.OHConfig().DamageMultiplier
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			},
		})

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Peregrine Trigger",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskMeleeOH,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			PPM:               1.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				character.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 469140}, spell)
				peregrineOHAttack.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterProc(Peregrine, triggerAura)
	})

	itemhelpers.CreateWeaponProcAura(Kestrel, "Kestrel", 1, func(character *core.Character) *core.Aura {
		return character.GetOrRegisterAura(core.Aura{
			Label:    "Kestrel Move Speed Aura",
			ActionID: core.ActionID{SpellID: 469148},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddMoveSpeedModifier(&aura.ActionID, 1.40)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.RemoveMoveSpeedModifier(&aura.ActionID)
			},
		})
	})

	// https://www.wowhead.com/classic/item=231288/renatakis-charm-of-ravaging
	core.NewItemEffect(RenatakisCharmOfRavaging, func(agent core.Agent) {
		character := agent.GetCharacter()

		lockedIn := character.RegisterAura(core.Aura{
			Label:     "Locked In",
			ActionID:  core.ActionID{SpellID: 468388},
			Duration:  time.Second * 20,
			MaxStacks: 2,
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.Matches(ClassSpellMask_HunterShots) || spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) && spell.CD.Timer != nil {
					if spell.Matches(ClassSpellMask_HunterRaptorStrike) {
						spell.CD.QueueReset(sim.CurrentTime)
					} else {
						spell.CD.Reset()
					}
					aura.RemoveStack(sim)
				}
			},
		})

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 468388},

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				lockedIn.Activate(sim)
				lockedIn.SetStacks(sim, lockedIn.MaxStacks)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	// https://www.wowhead.com/classic/item=230237/arcane-infused-gem
	core.NewItemEffect(ArcaneInfusedGem, func(agent core.Agent) {
		character := agent.GetCharacter()

		arcaneDetonation := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 467447},
			SpellSchool: core.SpellSchoolArcane,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					damage := sim.Roll(185, 210)
					spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		maxCarveTargetsPerCast := int32(5)
		maxMultishotTargetsPerCast := int32(3)

		arcaneInfused := character.RegisterAura(core.Aura{
			Label:    "Arcane Infused",
			ActionID: core.ActionID{SpellID: 467446},
			Duration: time.Second * 15,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				maxCarveTargetsPerCast = min(sim.Environment.GetNumTargets(), 5)
				maxMultishotTargetsPerCast = min(sim.Environment.GetNumTargets(), 3)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				// Uses same targeting code as multi-shot however the detonations occur at cast time rather than when the shots land
				if spell.Matches(ClassSpellMask_HunterMultiShot) {
					curTarget := sim.Environment.Encounter.TargetUnits[0]
					for hitIndex := int32(0); hitIndex < maxMultishotTargetsPerCast; hitIndex++ {
						arcaneDetonation.Cast(sim, curTarget)
						curTarget = sim.Environment.NextTargetUnit(curTarget)
					}
				}
				// 1 explosion per target up to 5 targets per carve cast
				if spell.Matches(ClassSpellMask_HunterCarve) {
					curTarget := sim.Environment.Encounter.TargetUnits[0]
					for hitIndex := int32(0); hitIndex < maxCarveTargetsPerCast; hitIndex++ {
						arcaneDetonation.Cast(sim, curTarget)
						curTarget = sim.Environment.NextTargetUnit(curTarget)
					}
				}
			},
		})

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: arcaneInfused.ActionID.SpellID},

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 90,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: arcaneInfused.Duration,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				arcaneInfused.Activate(sim)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(CloakOfTheUnseenPath, func(a core.Agent) {
		hunter := a.(HunterAgent).GetHunter()
		if hunter.pet == nil {
			return
		}

		procAura := core.MakePermanent(hunter.GetOrRegisterAura(core.Aura{
			Label:    "Increased Pet Damage +2% (Cloak of the Unseen Path)",
			ActionID: core.ActionID{SpellID: 468270},
		}).AttachMultiplicativePseudoStatBuff(&hunter.pet.PseudoStats.DamageDealtMultiplier, 1.02))

		hunter.ItemSwap.RegisterProc(CloakOfTheUnseenPath, procAura)
	})

	core.NewItemEffect(ScytheOfTheUnseenPath, func(a core.Agent) {
		hunter := a.(HunterAgent).GetHunter()
		if hunter.pet == nil {
			return
		}
		procAura := core.MakePermanent(hunter.GetOrRegisterAura(core.Aura{
			Label:    "Increased Pet Damage +3% (Scythe of the Unseen Path)",
			ActionID: core.ActionID{SpellID: 468268},
		}).AttachMultiplicativePseudoStatBuff(&hunter.pet.PseudoStats.DamageDealtMultiplier, 1.03))

		hunter.ItemSwap.RegisterProc(ScytheOfTheUnseenPath, procAura)
	})

	core.NewItemEffect(SignetOfTheUnseenPath, func(a core.Agent) {
		hunter := a.(HunterAgent).GetHunter()
		if hunter.pet == nil {
			return
		}

		procAura := core.MakePermanent(hunter.GetOrRegisterAura(core.Aura{
			Label:    "Increased Pet Damage +2% (Signet of the Unseen Path)",
			ActionID: core.ActionID{SpellID: 468270},
		}).AttachMultiplicativePseudoStatBuff(&hunter.pet.PseudoStats.DamageDealtMultiplier, 1.02))

		hunter.ItemSwap.RegisterProc(SignetOfTheUnseenPath, procAura)
	})

	core.NewItemEffect(StringsOfFate, func(a core.Agent) {
		hunter := a.(HunterAgent).GetHunter()
		hasMeleeSpecialist := hunter.HasRune(proto.HunterRune_RuneBeltMeleeSpecialist)

		// Tracks the number of strands of fate, up to 5 stacks if you're ranged or 4 if you're using the melee version.
		// TODO: Ranged Hunter's version lets them consume these stacks to move while shooting
		stacksAura := hunter.RegisterAura(core.Aura{
			Label:     "Strand Of Fate - Stacks",
			ActionID:  core.ActionID{SpellID: 1232946},
			MaxStacks: core.TernaryInt32(hasMeleeSpecialist, 4, 5),
			Duration:  time.Second * 20,
		})

		icd := core.Cooldown{
			Timer:    hunter.NewTimer(),
			Duration: time.Second * 30,
		}

		if !hasMeleeSpecialist {
			rangedTrigger := core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
				Name:       "Strand Of Fate - Ranged Trigger",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskRanged,
				ProcChance: 0.1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !icd.IsReady(sim) {
						return
					}
					stacksAura.Activate(sim)
					stacksAura.AddStack(sim)
				},
			})

			hunter.ItemSwap.RegisterProc(StringsOfFate, rangedTrigger)
		} else {
			stacksAura.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				strengthChange := float64(20 * (newStacks - oldStacks))
				hunter.Unit.AddStatDynamic(sim, stats.Strength, strengthChange)
			}

			meleeTrigger := core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
				Name:       "Strand Of Fate - Melee Trigger",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 0.1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !icd.IsReady(sim) {
						return
					}
					stacksAura.Activate(sim)
					stacksAura.AddStack(sim)
				},
			})

			hunter.ItemSwap.RegisterProc(StringsOfFate, meleeTrigger)
		}

		// When using the melee version and activating the bow your next strike applies a special Serpent Sting effect.
		serpentStrikeAura := hunter.RegisterAura(core.Aura{
			Label:     "Strand Of Fate - Serpent Sting on Next Strike",
			ActionID:  core.ActionID{SpellID: 1232976},
			Duration:  time.Second * 20,
			MaxStacks: 4,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.SetStacks(sim, stacksAura.GetStacks())
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Matches(ClassSpellMask_HunterStrikes) && result.Landed() {
					hunter.SoFSerpentSting[aura.GetStacks()].Cast(sim, result.Target)
					aura.Deactivate(sim)
				}
			},
		})

		// Gain 40 strength per stack of Strand of Fate consumed
		strengthAura := hunter.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 1232969},
			Label:     "Strand Of Fate - Strength",
			Duration:  time.Second * 20,
			MaxStacks: 4,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.SetStacks(sim, stacksAura.GetStacks())
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				hunter.AddStatDynamic(sim, stats.Strength, float64(40*(newStacks-oldStacks)))
			},
		})

		// The bow active
		spell := hunter.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 1231604},
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    hunter.NewTimer(),
					Duration: time.Second * 80,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				if hasMeleeSpecialist {
					if stacksAura.IsActive() {
						strengthAura.Activate(sim)
						serpentStrikeAura.Activate(sim)
						stacksAura.Deactivate(sim)
					}
				} else {
					stacksAura.Activate(sim)
					stacksAura.AddStacks(sim, 3)
				}
			},
		})

		hunter.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})

		hunter.ItemSwap.RegisterActive(StringsOfFate)
		hunter.ItemSwap.RegisterProc(StringsOfFate, stacksAura)
	})

	// https://www.wowhead.com/classic/item=240924/poleaxe-of-the-beast
	// Equip: Focus Fire now grants you and your pet 5% increased damage per stack consumed for 20 sec.
	core.NewItemEffect(PoleaxeOfTheBeast, func(agent core.Agent) {
		if agent.GetCharacter().Class != proto.Class_ClassHunter {
			return
		}

		hunter := agent.(HunterAgent).GetHunter()

		if hunter.pet == nil || !hunter.HasRune(proto.HunterRune_RuneBracersFocusFire) {
			return
		}

		bestialFocusAura := newBestialFocusAura(&hunter.Unit, 1231591)
		bestialFocusPetAura := newBestialFocusAura(&hunter.pet.Unit, 1231590)

		triggerAura := core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
			Name:           "Poleaxe of the Beast Trigger",
			Callback:       core.CallbackOnCastComplete,
			ClassSpellMask: ClassSpellMask_HunterFocusFire,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				stacks := hunter.FocusFireAura.GetStacks()
				bestialFocusAura.Activate(sim)
				bestialFocusAura.SetStacks(sim, stacks)
				bestialFocusPetAura.Activate(sim)
				bestialFocusPetAura.SetStacks(sim, stacks)
			},
		})

		hunter.ItemSwap.RegisterProc(PoleaxeOfTheBeast, triggerAura)
	})
}

// Your Raptor Strike and Mongoose Bite critical strikes set the duration of your Serpent Sting on the target to 15 sec
func (hunter *Hunter) ApplyQueensfallHunterEffect(aura *core.Aura) {
	aura.AttachProcTrigger(core.ProcTrigger{
		Name:           "Queensfall Trigger - Hunter",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeCrit,
		ClassSpellMask: ClassSpellMask_HunterRaptorStrikeHit | ClassSpellMask_HunterMongooseBite,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if dot := hunter.SerpentSting.Dot(result.Target); dot.IsActive() {
				dot.NumberOfTicks = int32(16 / dot.TickLength.Seconds())
				dot.RecomputeAuraDuration()
				dot.Rollover(sim)
			}
		},
	})
}

func newBestialFocusAura(unit *core.Unit, spellID int32) *core.Aura {
	return unit.RegisterAura(core.Aura{
		Label:     "Bestial Focus",
		ActionID:  core.ActionID{SpellID: spellID},
		Duration:  time.Second * 20,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= ((1.0 + (0.04 * float64(newStacks))) / (1.0 + (0.04 * float64(oldStacks))))
		},
	})
}

func (hunter *Hunter) newBloodlashProcItem(bonusStrength float64, spellID int32) {
	procAura := hunter.NewTemporaryStatsAura("Bloodlash", core.ActionID{SpellID: spellID}, stats.Stats{stats.Strength: bonusStrength}, time.Second*15)
	dpm := hunter.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMeleeOrRanged)

	aura := core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
		Name:     "Bloodlash Proc",
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		DPM:      dpm,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAura.Activate(sim)
		},
	})

	hunter.ItemSwap.RegisterProc(BloodlashBow, aura)
}

// Striking a higher level enemy applies a stack of Coup, increasing their damage taken from your next Kill Shot by 5% per stack, stacking up to 20 times.
func (hunter *Hunter) ApplyRegicideHunterEffect(itemID int32, aura *core.Aura) {
	// Coup debuff array
	debuffAuras := hunter.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
		return unit.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 1231765},
			Label:     "Coup",
			MaxStacks: core.TernaryInt32(unit.Level > hunter.Level, 20, 0),
			Duration:  time.Second * 15,
		})
	})

	killshotDamageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: ClassSpellMask_HunterKillShot,
	})

	damageModTrigger := core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
		Name:           "Coup - Kill Shot Damage Mod Trigger",
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: ClassSpellMask_HunterKillShot,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			killshotDamageMod.UpdateFloatValue(1 + float64(debuffAuras.Get(result.Target).GetStacks())*0.05)
			killshotDamageMod.Activate()
		},
	})
	hunter.ItemSwap.RegisterProc(itemID, damageModTrigger)

	consumptionTrigger := core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
		Name:           "Coup - Consume Stacks Trigger",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: ClassSpellMask_HunterKillShot,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			debuffAuras.Get(result.Target).Deactivate(sim)
		},
	})
	hunter.ItemSwap.RegisterProc(itemID, consumptionTrigger)

	// Apply the Coup debuff to the target hit by melee abilities
	aura.AttachProcTrigger(core.ProcTrigger{
		Name:              "Regicide Trigger - Hunter",
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		ProcMask:          core.ProcMaskMelee,
		SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			debuff := debuffAuras.Get(result.Target)
			debuff.Activate(sim)
			if debuff.MaxStacks > 0 {
				debuff.AddStack(sim)
			}
		},
	})
}

const MercyDamageBonus = 1.20

// Equip: Chance on hit to cause your next 2 instances of damage from your pet's special abilities to be increased by 20%. Lasts 12 sec. (100ms cooldown)
// Confirmed PPM 1.0
func (hunter *Hunter) ApplyMercyHunterEffect(aura *core.Aura) {
	if hunter.pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 1235361}

	// Create a dummy for UI tracking
	dummyAura := hunter.RegisterAura(core.Aura{
		ActionID:  actionID,
		Label:     "Mercy by Fire",
		Duration:  time.Second * 12,
		MaxStacks: 2,
	})

	icd := core.Cooldown{
		Timer:    hunter.NewTimer(),
		Duration: time.Millisecond * 100,
	}

	buffAura := hunter.pet.RegisterAura(core.Aura{
		ActionID:  actionID,
		Label:     "Mercy by Fire",
		Duration:  time.Second * 12,
		MaxStacks: 2,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_HunterPetBasicAttacks) && result.Landed() && icd.IsReady(sim) {
				icd.Use(sim)
				aura.RemoveStack(sim)
			}
		},
	}).AttachDependentAura(dummyAura)
	hunter.applyMercyAuraBonuses(buffAura, CrimsonCleaverDamageBonus)

	aura.AttachProcTrigger(core.ProcTrigger{
		Name:              "Mercy Trigger - Hunter",
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		ProcMask:          core.ProcMaskMelee, // Confirmed procs from either hand
		SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
		PPM:               1.0,
		ICD:               time.Millisecond * 100,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buffAura.Activate(sim)
			buffAura.SetStacks(sim, buffAura.MaxStacks)
		},
	})
}

func (hunter *Hunter) applyMercyAuraBonuses(aura *core.Aura, modifier float64) {
	aura.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  ClassSpellMask_HunterPetBasicAttacks,
		FloatValue: modifier,
	})
}

const CrimsonCleaverDamageBonus = 1.20

// Equip: Chance on hit to cause your next 2 instances of Raptor Strike damage to be increased by 20%. Lasts 12 sec. (100ms cooldown)
// Confirmed PPM 1.0
func (hunter *Hunter) ApplyCrimsonCleaverHunterEffect(aura *core.Aura) {
	icd := core.Cooldown{
		Timer:    hunter.NewTimer(),
		Duration: time.Millisecond * 100,
	}

	buffAura := hunter.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1235341},
		Label:     "Crimson Crusade",
		Duration:  time.Second * 12,
		MaxStacks: 2,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_HunterRaptorStrikeHit) && result.Landed() && icd.IsReady(sim) {
				icd.Use(sim)
				aura.RemoveStack(sim)
			}
		},
	})
	hunter.applyCrimsonCleaverAuraBonuses(buffAura, CrimsonCleaverDamageBonus)

	aura.AttachProcTrigger(core.ProcTrigger{
		Name:              "Crimson Cleaver Trigger - Hunter",
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		ProcMask:          core.ProcMaskMelee, // Confirmed procs from either hand
		SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
		PPM:               1.0,
		ICD:               time.Millisecond * 100,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buffAura.Activate(sim)
			buffAura.SetStacks(sim, buffAura.MaxStacks)
		},
	})
}

func (hunter *Hunter) applyCrimsonCleaverAuraBonuses(aura *core.Aura, modifier float64) {
	aura.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  ClassSpellMask_HunterRaptorStrikeHit,
		FloatValue: modifier,
	})
}
