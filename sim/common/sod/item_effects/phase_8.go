package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
	"github.com/wowsims/sod/sim/druid"
	"github.com/wowsims/sod/sim/hunter"
	"github.com/wowsims/sod/sim/mage"
	"github.com/wowsims/sod/sim/paladin"
	"github.com/wowsims/sod/sim/priest"
	"github.com/wowsims/sod/sim/rogue"
	"github.com/wowsims/sod/sim/shaman"
	"github.com/wowsims/sod/sim/warlock"
	"github.com/wowsims/sod/sim/warrior"
)

const (
	/* ! Please keep constants ordered by ID ! */

	CaladbolgSword           = 238961
	WillOfTheMountain        = 239060
	HighCommandersGuard      = 240841
	StartersPistol           = 240843
	LightfistHammer          = 240850
	Regicide                 = 240851
	CrimsonCleaver           = 240852
	Queensfall               = 240853
	Mercy                    = 240854
	Ravagane                 = 240919
	Deception                = 240922
	Duplicity                = 240923
	Experiment800M           = 240925
	SoporificBlade           = 240998
	TyrsFall                 = 241001
	RemnantsOfTheRed         = 241002
	MirageRodOfIllusion      = 241003
	Condemnation             = 241008
	GreatstaffOfFealty       = 241011
	AegisOfTheScarletBastion = 241015
	HeartOfLight             = 241034
	AbandonedExperiment      = 241037
	SirDornelsDidgeridoo     = 241038
	InfusionOfSouls          = 241039
	StiltzsStandard          = 241068
	ChokeChain               = 241069
	LuckyDoubloon            = 241241
	HandOfRebornJustice      = 242310
	CaladbolgMace            = 244460
)

func init() {
	core.AddEffectsToTest = false

	/* ! Please keep items ordered alphabetically ! */

	// https://www.wowhead.com/classic/item=241037/abandoned-experiment
	// Use: After drinking the experiment, ranged or melee attacks increase your attack speed by 2% for 30 sec.
	// This effect stacks up to 15 times. (2 Min Cooldown)
	core.NewItemEffect(AbandonedExperiment, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{ItemID: AbandonedExperiment}
		duration := time.Second * 30

		buffAura := character.RegisterAura(core.Aura{
			ActionID:  actionID,
			Label:     "Abandoned Experiment",
			MaxStacks: 15,
			Duration:  duration,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				character.MultiplyAttackSpeed(sim, 1/(1+0.02*float64(oldStacks)))
				character.MultiplyAttackSpeed(sim, 1+0.02*float64(newStacks))
			},
		})

		buffAura.MakeDependentProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Abandoned Experiment Trigger",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskMelee | core.ProcMaskRanged,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !buffAura.IsActive() {
					// Should not refresh
					buffAura.Activate(sim)
				}
				buffAura.AddStack(sim)
			},
		})

		cdSpell := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: cdSpell,
		})

		character.ItemSwap.RegisterActive(AbandonedExperiment)
		character.ItemSwap.RegisterProc(AbandonedExperiment, buffAura)
	})

	// https://www.wowhead.com/classic/item=241015/aegis-of-the-scarlet-bastion
	// Use: Increases the amount of damage absorbed by your shield by 20% for 15 sec. (1 Min, 30 Sec Cooldown)
	core.NewItemEffect(AegisOfTheScarletBastion, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: AegisOfTheScarletBastion}
		statDep := character.NewDynamicMultiplyStat(stats.BlockValue, 1.20)

		buffAura := character.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Scarlet Bastion",
			Duration: time.Second * 15,
		}).AttachStatDependency(statDep)

		cdSpell := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			ProcMask: core.ProcMaskEmpty,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 90,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell: cdSpell,
			Type:  core.CooldownTypeSurvival,
		})

		character.ItemSwap.RegisterActive(AegisOfTheScarletBastion)
	})

	// https://www.wowhead.com/classic/item=244460/caladbolg
	core.NewItemEffect(CaladbolgMace, func(agent core.Agent) {
		makeCaladbolgEffect(CaladbolgMace, agent)
	})
	// https://www.wowhead.com/classic/item=238961/caladbolg
	core.NewItemEffect(CaladbolgSword, func(agent core.Agent) {
		makeCaladbolgEffect(CaladbolgSword, agent)
	})

	// https://www.wowhead.com/classic/item=241069/choke-chain
	core.NewItemEffect(ChokeChain, func(agent core.Agent) {
		character := agent.GetCharacter()

		// There's a hidden effect that causes Shamans and Warlocks to receive 2 hit instead of 2 expertise
		if character.Class == proto.Class_ClassShaman || character.Class == proto.Class_ClassWarlock {
			character.AddStats(stats.Stats{
				stats.Expertise: -2 * core.ExpertiseRatingPerExpertiseChance,
				stats.MeleeHit:  2 * core.MeleeHitRatingPerHitChance,
				stats.SpellHit:  2 * core.SpellHitRatingPerHitChance,
			})
		}
	})

	// https://www.wowhead.com/classic/item=241008/condemnation
	// Equip: Damaging spell casts on enemies Condemns them for 20 sec.
	// Whenever they deal damage, their target is healed for 3% of their maximum health.
	// Lasts 20 sec or up to 10 hits. (60 Secs Cooldown) (1m cooldown)
	core.NewItemEffect(Condemnation, func(agent core.Agent) {
		character := agent.GetCharacter()

		healthMetrics := character.NewHealthMetrics(core.ActionID{ItemID: Condemnation})

		debuffs := character.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
			return core.MakeProcTriggerAura(unit, core.ProcTrigger{
				ActionID: core.ActionID{SpellID: 1231695},
				Name:     "Condemned",
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
				ProcMask: core.ProcMaskWhiteHit, // Confirmed via Wago https://wago.tools/db2/SpellAuraOptions?build=1.15.7.60000&filter%5BSpellID%5D=1231695&page=1
				ICD:      time.Second,
				Duration: time.Second * 20,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					result.Target.GainHealth(sim, 0.03*result.Target.MaxHealth(), healthMetrics)
				},
			})
		})

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Condemnation",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskSpellDamage,
			ICD:      time.Minute,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				debuffs.Get(result.Target).Activate(sim)
			},
		})

		character.ItemSwap.RegisterProc(Condemnation, triggerAura)
	})

	// https://www.wowhead.com/classic/item=240852/crimson-cleaver
	// Hunter - Equip: Chance on hit to cause your next 2 instances of Raptor Strike damage to be increased by 20%. Lasts 12 sec.
	// Paladin - Equip: Chance on hit to cause your next 2 instances of Holy damage to be increased by 20%. Lasts 12 sec.
	// Shaman - Equip: Chance on hit to cause your next 2 instances of Nature damage are increased by 20%. Lasts 12 sec.
	// Warrior - Equip: Chance on hit to cause your next 2 instances of Cleave damage to be increased by 20%. Lasts 12 sec.
	core.NewItemEffect(CrimsonCleaver, func(agent core.Agent) {
		character := agent.GetCharacter()

		aura := core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Crimson Cleaver Trigger",
		}))

		switch character.Class {
		case proto.Class_ClassHunter:
			agent.(hunter.HunterAgent).GetHunter().ApplyCrimsonCleaverHunterEffect(aura)
		case proto.Class_ClassPaladin:
			agent.(paladin.PaladinAgent).GetPaladin().ApplyCrimsonCleaverPaladinEffect(aura)
		case proto.Class_ClassShaman:
			agent.(shaman.ShamanAgent).GetShaman().ApplyCrimsonCleaverShamanEffect(aura)
		case proto.Class_ClassWarrior:
			agent.(warrior.WarriorAgent).GetWarrior().ApplyCrimsonCleaverWarriorEffect(aura)
		}

		character.ItemSwap.RegisterProc(CrimsonCleaver, aura)
	})

	// https://www.wowhead.com/classic/item=240922/deception
	// Equip: 2% chance on melee hit to gain 1 extra attack. (Proc chance: 2%, 100ms cooldown)
	core.NewItemEffect(Deception, func(agent core.Agent) {
		character := agent.GetCharacter()

		isMH  := character.GetMHWeapon().ID == Deception
		procMaskAura := core.Ternary(isMH, core.ProcMaskMeleeMH, core.ProcMaskMeleeOH)

		spellProc := character.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: 1231552},   
			SpellSchool:    core.SpellSchoolPhysical,
			DefenseType:    core.DefenseTypeMelee,
			ProcMask:       core.ProcMaskMeleeMHAuto, // Normal Melee Attack Flag
			Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagSuppressWeaponProcs, // Cannot proc Oil, Poisons, and presumably Weapon Enchants or Procs(Chance on Hit)
			CastType:       proto.CastType_CastTypeMainHand,
	
			DamageMultiplier: 1.0,
			ThreatMultiplier: 1,
	
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			},
		})

		character.NewRageMetrics(spellProc.ActionID)
		spellProc.ResourceMetrics = character.NewRageMetrics(spellProc.ActionID)

		// Use a dummy to set a flag for the set bonus that doubles the extra attacks
		var setAura *core.Aura
		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:          core.ActionID{SpellID: 1231553},
			Name:              "Deception Proc",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          procMaskAura,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			ProcChance:        0.02,
			ICD:               time.Millisecond * 100,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if (setAura != nil && setAura.IsActive()) {
					return
				} else {
					spellProc.Cast(sim, result.Target)
				}
			},
		}).ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			setAura = character.GetAura("Tools of the Nathrezim")

			if setAura != nil {
				procMaskAura = core.ProcMaskMelee
			}
		})

		character.ItemSwap.RegisterProc(Deception, triggerAura)
	})

	// https://www.wowhead.com/classic/item=240923/duplicity
	// Equip: 2% chance on melee hit to gain 1 extra attack. (Proc chance: 2%, 100ms cooldown)
	core.NewItemEffect(Duplicity, func(agent core.Agent) {
		character := agent.GetCharacter()

		isMH  := character.GetMHWeapon().ID == Duplicity
		procMaskAura := core.Ternary(isMH, core.ProcMaskMeleeMH, core.ProcMaskMeleeOH)

		spellProc := character.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: 1231555},
			SpellSchool:    core.SpellSchoolPhysical,
			DefenseType:    core.DefenseTypeMelee,
			ProcMask:       core.ProcMaskMeleeMHAuto, // Normal Melee Attack Flag
			Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagSuppressWeaponProcs, // Cannot proc Oil, Poisons, and presumably Weapon Enchants or Procs(Chance on Hit)
			CastType:       proto.CastType_CastTypeMainHand,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			},
		})

		character.NewRageMetrics(spellProc.ActionID)
		spellProc.ResourceMetrics = character.NewRageMetrics(spellProc.ActionID)

		// Use a dummy to set a flag for the set bonus that doubles the extra attacks
		var setAura *core.Aura
		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:          core.ActionID{SpellID: 1231554},
			Name:              "Duplicity Proc",  
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          procMaskAura,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			ProcChance:        0.02,
			ICD:               time.Millisecond * 100,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if (setAura != nil && setAura.IsActive()) {
					return
				} else {
					spellProc.Cast(sim, result.Target)
				}
			},
		}).ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			setAura = character.GetAura("Tools of the Nathrezim")
		})

		character.ItemSwap.RegisterProc(Duplicity, triggerAura)
	})

	// https://www.wowhead.com/classic/item=240925/experiment-800m
	// Use: Transform into a Scarlet Cannon for 20 sec, causing your ranged attacks to explode for 386 to 472 Fire damage to all enemies within 8 yards of your target. (2 Min Cooldown)
	core.NewItemEffect(Experiment800M, func(agent core.Agent) {
		character := agent.GetCharacter()

		numHits := character.Env.GetNumTargets()
		explosionSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1231607},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagIgnoreAttackerModifiers | core.SpellFlagIgnoreTargetModifiers,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					damage := sim.Roll(386, 472)
					spell.CalcAndDealDamage(sim, curTarget, damage, spell.OutcomeMagicCrit)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		})

		buffAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231605},
			Label:    "EXPERIMENT-8OOM!!!",
			Duration: time.Second * 20,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellResult *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskRanged) && spellResult.Landed() {
					explosionSpell.Cast(sim, spellResult.Target)
				}
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 1231605},
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})

		character.ItemSwap.RegisterActive(Experiment800M)
	})

	// https://www.wowhead.com/classic/item=241011/greatstaff-of-fealty
	// Equip: If there are no enemies within 20 yards of you, heal all party members within 20 yards for 125 every 3 seconds.
	// A party may only swear Fealty to one player at a time.
	core.NewItemEffect(GreatstaffOfFealty, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{ItemID: GreatstaffOfFealty}
		healthMetrics := character.NewHealthMetrics(actionID)

		var healingPA *core.PendingAction
		healingAura := character.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Fealty",
			Duration: core.NeverExpires,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				healingPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:   time.Second * 3,
					Priority: core.ActionPriorityLow,
					OnAction: func(sim *core.Simulation) {
						character.GainHealth(sim, 125, healthMetrics)
					},
				})
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				healingPA.Cancel(sim)
			},
		})

		var rangeCheckPA *core.PendingAction
		triggerAura := core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Greater of Fealty Periodic Trigger",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				rangeCheckPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:          time.Second * 1,
					Priority:        core.ActionPriorityLow,
					TickImmediately: true,
					OnAction: func(sim *core.Simulation) {
						if character.DistanceFromTarget >= 20 {
							healingAura.Activate(sim)
						} else {
							healingAura.Deactivate(sim)
						}
					},
				})
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				healingAura.Deactivate(sim)
				rangeCheckPA.Cancel(sim)
			},
		}))

		character.ItemSwap.RegisterProc(GreatstaffOfFealty, triggerAura)
	})

	// https://www.wowhead.com/classic/item=242310/hand-of-reborn-justice
	// Equip: 2% chance on melee or ranged hit to gain 1 extra attack. (Proc chance: 2%, 2s cooldown)
	core.NewItemEffect(HandOfRebornJustice, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.AutoSwingMelee {
			return
		}

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Hand of Reborn Justice Trigger (Melee)",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskMelee | core.ProcMaskRanged,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			ProcChance:        0.02,
			ICD:               time.Second * 2,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMelee) {
					spell.Unit.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 1232044}, spell)
				} else {
					character.AutoAttacks.StoreExtraRangedAttack(sim, 1, core.ActionID{SpellID: 1213381}, spell.ActionID)
				}
			},
		})

		character.ItemSwap.RegisterProc(HandOfInjustice, triggerAura)
	})

	// https://www.wowhead.com/classic/item=241034/heart-of-light
	// Use: Increases maximum health by 2500 for 20 sec. (2 Min Cooldown)
	core.NewItemEffect(HeartOfLight, func(agent core.Agent) {
		character := agent.GetCharacter()

		// There's a hidden effect that causes Shamans and Warlocks to receive 2 hit instead of 2 expertise
		if character.Class == proto.Class_ClassShaman || character.Class == proto.Class_ClassWarlock {
			character.AddStats(stats.Stats{
				stats.Expertise: -2 * core.ExpertiseRatingPerExpertiseChance,
				stats.MeleeHit:  2 * core.MeleeHitRatingPerHitChance,
				stats.SpellHit:  2 * core.SpellHitRatingPerHitChance,
			})
		}

		actionID := core.ActionID{ItemID: HeartOfLight}

		buffAura := character.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Heart of Light",
			Duration: time.Second * 20,
		}).AttachStatBuff(stats.Health, 250)

		cdSpell := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeSurvival,
			Spell: cdSpell,
		})
	})

	// https://www.wowhead.com/classic/item=240841/high-commanders-guard
	// Chance on hit: Increase Defense by 20 and Armor by 750 for 10 sec.
	// Confirmed PPM 2.5
	itemhelpers.CreateWeaponProcAura(HighCommandersGuard, "High Commander's Guard", 2.5, func(character *core.Character) *core.Aura {
		return character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231254},
			Label:    "Scarlet Bulwark",
			Duration: time.Second * 10,
		}).AttachStatsBuff(stats.Stats{
			stats.Defense: 20,
			stats.Armor:   750,
		})
	})

	// https://www.wowhead.com/classic/item=241039/infusion-of-souls
	// The Global Cooldown caused by your non-weapon based damaging spells can be reduced by Spell Haste, up to a 0.5 second reduction.
	core.NewItemEffect(InfusionOfSouls, func(agent core.Agent) {
		character := agent.GetCharacter()

		var classMask uint64
		switch character.Class {
		// https://www.wowhead.com/classic/spell=1232094/infusion-of-souls
		case proto.Class_ClassDruid:
			classMask = druid.ClassSpellMask_DruidHarmfulGCDSpells

		// https://www.wowhead.com/classic/spell=1230948/infusion-of-souls
		case proto.Class_ClassMage:
			classMask = mage.ClassSpellMask_MageHarmfulGCDSpells

		// https://www.wowhead.com/classic/spell=1232104/infusion-of-souls
		case proto.Class_ClassPaladin:
			// Explicitly lists that it does not work for Holy Shock in the tooltip https://www.wowhead.com/classic/item=241039/infusion-of-souls?spellModifier=462814
			classMask = paladin.ClassSpellMask_PaladinHarmfulGCDSpells ^ paladin.ClassSpellMask_PaladinHolyShock

		// https://www.wowhead.com/classic/spell=1232095/infusion-of-souls
		case proto.Class_ClassPriest:
			// Explicitly lists that it does not work for Penance in the tooltip https://www.wowhead.com/classic/item=241039/infusion-of-souls?spellModifier=440247
			classMask = priest.ClassSpellMask_PriestHarmfulGCDSpells ^ priest.ClassSpellMask_PriestPenance

		// https://www.wowhead.com/classic/spell=1232096/infusion-of-souls
		case proto.Class_ClassShaman:
			// Explicitly lists that it does not work while Way of Earth is active
			classMask = core.Ternary(agent.(shaman.ShamanAgent).GetShaman().WayOfEarthActive(), 0, shaman.ClassSpellMask_ShamanHarmfulGCDSpells)

		// https://www.wowhead.com/classic/spell=1232093/infusion-of-souls
		case proto.Class_ClassWarlock:
			classMask = warlock.ClassSpellMask_WarlockHarmfulGCDSpells

			// Infusion of souls also affects Warlock pets
			warlockPlayer := agent.(warlock.WarlockAgent).GetWarlock()
			for _, pet := range warlockPlayer.BasePets {
				pet.OnSpellRegistered(func(spell *core.Spell) {
					if spell.Matches(classMask) {
						spell.AllowGCDHasteScaling = true

						if spell.Matches(warlock.ClassSpellMask_WarlockSummonImpFireBolt) {
							spell.DefaultCast.GCDMin = time.Millisecond * 500
						}
					}
				})
			}
		}

		character.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Matches(classMask) {
				spell.AllowGCDHasteScaling = true
			}
		})
	})

	// https://www.wowhead.com/classic/item=240850/lightfist-hammer
	// Chance on hit: Increases your attack speed by 10% for 15 sec.
	// Confirmed PPM 0.7
	itemhelpers.CreateWeaponProcAura(LightfistHammer, "Lightfist Hammer", 0.7, func(character *core.Character) *core.Aura {
		return character.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: LightfistHammer},
			Label:    "Lightfist Hammer",
			Duration: time.Second * 15,
		}).AttachMultiplyMeleeSpeed(&character.Unit, 1.10)
	})

	// https://www.wowhead.com/classic/item=241241/lucky-doubloon
	// Use: Flip the Lucky Doubloon.
	// On heads, the cooldown of Lucky Doubloon is refreshed, and your critical strike chance with all spells and attacks is increased by 5% for 15 sec, stacking up to 5 times.
	// On tails, all stacks of Lucky Doubloon are removed. (Must be in combat.) (30 Sec Cooldown)
	core.NewItemEffect(LuckyDoubloon, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{ItemID: LuckyDoubloon}
		duration := time.Second * 15

		buffAura := character.RegisterAura(core.Aura{
			ActionID:  actionID,
			Label:     "Lucky Doubloon",
			MaxStacks: 5,
			Duration:  duration,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				bonusCrit := 5 * float64(newStacks-oldStacks)
				character.AddStatsDynamic(sim, stats.Stats{stats.MeleeCrit: bonusCrit, stats.SpellCrit: bonusCrit})
			},
		})

		cdSpell := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 30,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				if sim.Proc(0.5, "Lucky Doubloon Heads") {
					buffAura.Activate(sim)
					buffAura.AddStack(sim)

					// Effect has a 1.5s ICD to avoid misclicks
					sim.AddPendingAction(&core.PendingAction{
						NextActionAt: sim.CurrentTime + time.Millisecond*1500,
						OnAction: func(sim *core.Simulation) {
							spell.CD.Reset()
						},
					})
				} else {
					buffAura.Deactivate(sim)
				}
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: cdSpell,
		})

		character.ItemSwap.RegisterActive(LuckyDoubloon)
	})

	// https://www.wowhead.com/classic/item=240854/mercy
	// Hunter - Equip: Chance on hit to cause your next 2 instances of damage from your pet's special abilities to be increased by 20%. Lasts 12 sec.
	// Shaman - Equip: Chance on hit to cause your next 2 instances of Fire damage are increased by 20%.  Lasts 12 sec. (100ms cooldown)
	// Warrior - Equip: Chance on hit to cause your next 2 instances of Whirlwind damage to be increased by 20%. Lasts 12 sec.
	core.NewItemEffect(Mercy, func(agent core.Agent) {
		character := agent.GetCharacter()

		aura := core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Mercy Trigger",
		}))

		switch character.Class {
		case proto.Class_ClassHunter:
			agent.(hunter.HunterAgent).GetHunter().ApplyMercyHunterEffect(aura)
		case proto.Class_ClassShaman:
			agent.(shaman.ShamanAgent).GetShaman().ApplyMercyShamanEffect(aura)
		case proto.Class_ClassWarrior:
			agent.(warrior.WarriorAgent).GetWarrior().ApplyMercyWarriorEffect(aura)
		}

		character.ItemSwap.RegisterProc(Mercy, aura)
	})

	// https://www.wowhead.com/classic/item=241003/mirage-rod-of-illusion
	// Equip: Chance on landing a damaging spell to create a Mirage on top of your target that deals arcane damage to nearby enemies for 30 sec. (Proc chance: 10%, 30s cooldown)
	core.NewItemEffect(MirageRodOfIllusion, func(agent core.Agent) {
		character := agent.GetCharacter()

		explosionSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1231648}, // TODO: Verify real spell ID
			SpellSchool: core.SpellSchoolArcane,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty, // Cast by an NPC so presumably no procs
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			BonusCoefficient: 0.143, // TODO: Taken from the above spell. Verify

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(243, 263), spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		// Use an AOE DoT rather than a real NPC for performance since it hopefully just casts a constant spell
		dot := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1231629}, // https://www.wowhead.com/classic/spell=1231629/mirage
			SpellSchool: core.SpellSchoolArcane,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Mirage",
				},
				NumberOfTicks: 20,
				TickLength:    core.GCDDefault,
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					explosionSpell.Cast(sim, target)
				},
			},

			BonusCoefficient: 0.143, // TODO: Taken from the above spell. Verify

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				if dot := spell.Dot(target); dot != nil {
					dot.Apply(sim)
				}
			},
		})

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Mirage Rod of Illusion Trigger",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.10,
			ICD:        time.Second * 30,
			Harmful:    true,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				dot.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterProc(MirageRodOfIllusion, triggerAura)
	})

	// https://www.wowhead.com/classic/item=240851/regicide
	// Striking a higher level enemy applies a stack of Coup, increasing their damage taken from your next Execute by 10% per stack, stacking up to 20 times. At 20 stacks, Execute may be cast regardless of the target's health.
	// Striking a higher level enemy applies a stack of Coup, increasing their damage taken from your next Envenom by 5% per stack, stacking up to 20 times.
	// Striking a higher level enemy applies a stack of Coup, increasing their damage taken from your next Kill Shot by 5% per stack, stacking up to 20 times.
	core.NewItemEffect(Regicide, func(agent core.Agent) {
		character := agent.GetCharacter()

		aura := core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Regicide Trigger",
		}))

		switch character.Class {
		case proto.Class_ClassWarrior:
			agent.(warrior.WarriorAgent).GetWarrior().ApplyRegicideWarriorEffect(Regicide, aura)
		case proto.Class_ClassRogue:
			agent.(rogue.RogueAgent).GetRogue().ApplyRegicideRogueEffect(Regicide, aura)
		case proto.Class_ClassHunter:
			agent.(hunter.HunterAgent).GetHunter().ApplyRegicideHunterEffect(Regicide, aura)
		}

		character.ItemSwap.RegisterProc(Regicide, aura)
	})

	// https://www.wowhead.com/classic/item=240853/queensfall
	// Your Bloodthirst, Mortal Strike, Shield Slam, Heroic Strike, and Cleave critical strikes set the duration of your Rend on the target to 21 sec.
	// Your Backstab, Mutilate, and Saber Slash critical strikes set the duration of your Rupture on the target to 16 secs
	// Your Raptor Strike and Mongoose Bite critical strikes set the duration of your Serpent Sting on the target to 15 sec
	core.NewItemEffect(Queensfall, func(agent core.Agent) {
		character := agent.GetCharacter()

		aura := core.MakePermanent(character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1232181},
			Label:    "Queensfall Trigger",
		}))

		switch character.Class {
		case proto.Class_ClassWarrior:
			agent.(warrior.WarriorAgent).GetWarrior().ApplyQueensfallWarriorEffect(aura)
		case proto.Class_ClassRogue:
			agent.(rogue.RogueAgent).GetRogue().ApplyQueensfallRogueEffect(aura)
		case proto.Class_ClassHunter:
			agent.(hunter.HunterAgent).GetHunter().ApplyQueensfallHunterEffect(aura)
		}
	})

	// https://www.wowhead.com/classic/item=240919/ravagane
	// Chance on hit: You attack all nearby enemies for 9 sec causing weapon damage plus an additional 200 every 1.5 sec.
	// Confirmed PPM 0.8
	itemhelpers.CreateWeaponProcAura(Ravagane, "Ravagane", 0.8, func(character *core.Character) *core.Aura {
		tickSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1231546},
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				damage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) + 200
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMeleeSpecialHitAndCrit)
				}
			},
		})

		whirlwindSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1231547},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Ravagane Whirlwind",
				},
				NumberOfTicks: 6,
				TickLength:    time.Millisecond * 1500,
				IsAOE:         true,
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					tickSpell.Cast(sim, target)
				},
			},
		})

		return character.RegisterAura(core.Aura{
			Label:    "Ravagane Bladestorm",
			Duration: time.Second * 9,
			Icd: &core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 8,
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if !aura.Icd.IsReady(sim) {
					return
				}
				aura.Icd.Use(sim)

				whirlwindSpell.AOEDot().Apply(sim)
				character.AutoAttacks.CancelAutoSwing(sim)
			},
			OnRefresh: func(aura *core.Aura, sim *core.Simulation) {
				if !aura.Icd.IsReady(sim) {
					return
				}
				aura.Icd.Use(sim)

				whirlwindSpell.AOEDot().ApplyOrReset(sim)
				character.AutoAttacks.CancelAutoSwing(sim)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				whirlwindSpell.AOEDot().Cancel(sim)
				character.AutoAttacks.EnableAutoSwing(sim)
			},
		})
	})

	// https://www.wowhead.com/classic/item=241002/remnants-of-the-red
	// Equip: Dealing non-periodic Fire damage has a 5% chance to increase your Fire damage dealt by 5% for 15 sec. (Proc chance: 5%)
	core.NewItemEffect(RemnantsOfTheRed, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231625},
			Label:    "Flames of the Red",
			Duration: time.Second * 15,
		}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire], 1.05)

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:        "Remnants of the Red Trigger",
			Callback:    core.CallbackOnSpellHitDealt,
			Outcome:     core.OutcomeLanded,
			ProcMask:    core.ProcMaskSpellDamage,
			SpellSchool: core.SpellSchoolFire,
			ProcChance:  0.05,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterProc(RemnantsOfTheRed, triggerAura)
	})

	// https://www.wowhead.com/classic/item=241038/sir-dornels-didgeridoo
	// Use: Playing the didgeridoo summons a random beast to your side, increasing your physical abilities for 30 sec.
	// May only be used while in combat. (2 Min Cooldown)
	core.NewItemEffect(SirDornelsDidgeridoo, func(agent core.Agent) {
		character := agent.GetCharacter()

		duration := time.Second * 30

		bonusCrit := 10.0
		bonusArp := 1000.0
		bonusAgi := 140.0
		bonusStr := 140.0
		bonusHaste := 0.14
		bonusAP := 280.0
		bonusRAP := 308.0

		// https://www.wowhead.com/classic/spell=1231884/accuracy-of-the-owl
		owlBuff := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231884},
			Label:    "Sir Dornel's Didgeridoo - Owl",
			Duration: duration,
		}).AttachStatBuff(stats.MeleeCrit, bonusCrit).AttachStatBuff(stats.SpellCrit, bonusCrit)

		// https://www.wowhead.com/classic/spell=1231894/ferocity-of-the-crocolisk
		crocoliskBuff := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231894},
			Label:    "Sir Dornel's Didgeridoo - Crocolisk",
			Duration: duration,
		}).AttachStatBuff(stats.ArmorPenetration, bonusArp)

		// https://www.wowhead.com/classic/spell=1231888/agility-of-the-raptor
		raptorBuff := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231888},
			Label:    "Sir Dornel's Didgeridoo - Raptor",
			Duration: duration,
		}).AttachStatBuff(stats.Agility, bonusAgi)

		// https://www.wowhead.com/classic/spell=1231883/strength-of-the-bear
		bearBuff := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231883},
			Label:    "Sir Dornel's Didgeridoo - Bear",
			Duration: duration,
		}).AttachStatBuff(stats.Strength, bonusStr)

		// https://www.wowhead.com/classic/spell=1231886/speed-of-the-cat
		catBuff := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231886},
			Label:    "Sir Dornel's Didgeridoo - Cat",
			Duration: duration,
		}).AttachMultiplyAttackSpeed(&character.Unit, 1+bonusHaste)

		// https://www.wowhead.com/classic/spell=1231891/power-of-the-gorilla
		gorillaBuff := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231891},
			Label:    "Sir Dornel's Didgeridoo - Gorilla",
			Duration: duration,
		}).AttachStatBuff(stats.AttackPower, bonusAP).AttachStatBuff(stats.RangedAttackPower, bonusRAP)

		// https://www.wowhead.com/classic/spell=1231896/brilliance-of-mr-bigglesworth
		bigglesworthBuff := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231896},
			Label:    "Sir Dornel's Didgeridoo - Mr. Bigglesworth",
			Duration: duration,
		}).AttachStatBuff(
			stats.MeleeCrit, bonusCrit/2.0,
		).AttachStatBuff(
			stats.SpellCrit, bonusCrit/2.0,
		).AttachStatBuff(
			stats.ArmorPenetration, bonusArp/2.0,
		).AttachStatBuff(
			stats.Agility, bonusAgi/2.0,
		).AttachStatBuff(
			stats.Strength, bonusStr/2.0,
		).AttachMultiplyAttackSpeed(
			&character.Unit, 1+bonusHaste/2.0,
		).AttachStatBuff(
			stats.AttackPower, bonusAP/2.0,
		).AttachStatBuff(
			stats.RangedAttackPower, bonusRAP/2.0,
		)

		var buffAuras []*core.Aura

		switch character.Class {
		case proto.Class_ClassDruid:
			buffAuras = []*core.Aura{bearBuff, catBuff, crocoliskBuff, raptorBuff}
		case proto.Class_ClassHunter:
			buffAuras = []*core.Aura{crocoliskBuff, gorillaBuff, owlBuff, raptorBuff}
		case proto.Class_ClassPaladin:
			buffAuras = []*core.Aura{bearBuff, catBuff, gorillaBuff, owlBuff}
		case proto.Class_ClassRogue:
			buffAuras = []*core.Aura{catBuff, crocoliskBuff, gorillaBuff, raptorBuff}
		case proto.Class_ClassShaman:
			buffAuras = []*core.Aura{bearBuff, catBuff, gorillaBuff, owlBuff}
		case proto.Class_ClassWarrior:
			buffAuras = []*core.Aura{bearBuff, catBuff, crocoliskBuff, owlBuff}
		}

		cdSpell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: SirDornelsDidgeridoo},
			ProcMask: core.ProcMaskEmpty,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				// Just in case people try to sim with an invalid class..
				if len(buffAuras) == 0 {
					return
				}

				if int32(sim.Roll(0, 100)) < 1 {
					bigglesworthBuff.Activate(sim)
				} else {
					buffAuras[int32(sim.Roll(0, 4))].Activate(sim)
				}
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: cdSpell,
		})

		character.ItemSwap.RegisterActive(SirDornelsDidgeridoo)
	})

	// https://www.wowhead.com/classic/item=240998/soporific-blade
	// Use: Your next regular melee attack made within 20 sec deals Arcane damage equal to 2 times your spell power, and puts the target to sleep for 20 sec. (2 Min Cooldown)
	core.NewItemEffect(SoporificBlade, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: SoporificBlade}

		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 1231614},
			SpellSchool:      core.SpellSchoolArcane,
			DefenseType:      core.DefenseTypeMagic,
			ProcMask:         core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc,
			Flags:            core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 2,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 0, spell.OutcomeMagicHitAndCrit)
			},
		})

		procTrigger := character.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Soporific Blade Trigger",
			Duration: time.Second * 20,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && result.Landed() {
					damageSpell.Cast(sim, result.Target)
					aura.Deactivate(sim)
				}
			},
		})

		cdSpell := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			ProcMask: core.ProcMaskEmpty,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				procTrigger.Activate(sim)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: cdSpell,
		})

		character.ItemSwap.RegisterActive(SoporificBlade)
	})

	// https://www.wowhead.com/classic/item=240843/starters-pistol
	// Equip: Firing a regular ranged attack at a target prepares you for battle, increasing your Defense by 20 and melee attack speed by 10% for 15 sec.
	core.NewItemEffect(StartersPistol, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231266},
			Label:    "En Garde!",
			Duration: time.Second * 15,
		}).AttachStatBuff(stats.Defense, 20).AttachMultiplyMeleeSpeed(&character.Unit, 1.10)

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "En Garde!",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskRangedAuto,
			ProcChance: 1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterProc(StartersPistol, triggerAura)
	})

	// https://www.wowhead.com/classic/item=241068/stiltzs-standard
	// Use: Throw down the Standard of Stiltz, increasing the maximum health of all nearby allies by 1000 for 20 sec. (2 Min Cooldown)
	core.NewSimpleStatDefensiveTrinketEffect(StiltzsStandard, stats.Stats{stats.Health: 1000}, time.Second*20, time.Minute*2)

	// https://www.wowhead.com/classic/item=241001/tyrs-fall
	// Equip: Dealing periodic damage has a 5% chance to grant 120 spell damage and healing for 15 sec. (Proc chance: 5%)
	core.NewItemEffect(TyrsFall, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231623},
			Label:    "Tyr's Return",
			Duration: time.Second * 15,
		}).AttachStatBuff(stats.SpellPower, 120)

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Tyr's Fall Trigger",
			Callback:   core.CallbackOnPeriodicDamageDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.05,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterProc(TyrsFall, triggerAura)
	})

	// https://www.wowhead.com/classic/item=239060/will-of-the-mountain
	// Chance on hit: Invoke the Will of the Mountain, increasing physical damage dealt by 5%, armor by 500, and size by 20% for 15 sec.
	// PPM confirmed 0.5
	itemhelpers.CreateWeaponProcAura(WillOfTheMountain, "Will of the Mountain", 0.5, func(character *core.Character) *core.Aura {
		return character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231289},
			Label:    "Avatar of the Mountain",
			Duration: time.Second * 15,
		}).AttachMultiplicativePseudoStatBuff(
			&character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical], 1.05,
		).AttachStatBuff(stats.Armor, 500)
	})

	core.AddEffectsToTest = true
}

// Use: The earth once again swallow up Caladbolg for 20 sec, decreasing movement speed by 50%. If moving, your speed increases every 2 seconds.
// Smashing into an enemy deals 300 to 2700 fire damage, knocks back all nearby enemy players, and ends this effect. (2 Min Cooldown)
func makeCaladbolgEffect(itemID int32, agent core.Agent) {
	character := agent.GetCharacter()

	moveSpeedActionID := &core.ActionID{SpellID: 1232322}

	damageSpell := character.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 1232323},
		SpellSchool:      core.SpellSchoolFire,
		DefenseType:      core.DefenseTypeMagic,
		ProcMask:         core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc,
		Flags:            core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(300, 2700), spell.OutcomeMagicCrit) // Has a flag that it always hits
			}
		},
	})

	// TODO: Increasing move speed every 2 seconds if we really care enough
	var rangeDummyPA *core.PendingAction
	rangeDummyAura := character.RegisterAura(core.Aura{
		Label:    "Caladbolg Range Dummy",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rangeDummyPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Millisecond * 500,
				Priority: core.ActionPriorityLow,
				OnAction: func(sim *core.Simulation) {
					if character.DistanceFromTarget <= core.MaxMeleeAttackRange {
						damageSpell.Cast(sim, character.CurrentTarget)
						character.RemoveMoveSpeedModifier(moveSpeedActionID)
						rangeDummyPA.Cancel(sim)
					}
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rangeDummyPA.Cancel(sim)
		},
	})

	cdSpell := character.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{ItemID: itemID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			character.AddMoveSpeedModifier(moveSpeedActionID, 0.50)
			character.MoveTo(core.MaxMeleeAttackRange, sim)
			rangeDummyAura.Activate(sim)
		},
	})

	character.AddMajorCooldown(core.MajorCooldown{
		Spell: cdSpell,
		Type:  core.CooldownTypeDPS,
	})

	character.ItemSwap.RegisterActive(itemID)
	character.ItemSwap.RegisterProc(itemID, rangeDummyAura)
}
