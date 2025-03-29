package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
	"github.com/wowsims/sod/sim/druid"
	"github.com/wowsims/sod/sim/mage"
	"github.com/wowsims/sod/sim/paladin"
	"github.com/wowsims/sod/sim/priest"
	"github.com/wowsims/sod/sim/shaman"
	"github.com/wowsims/sod/sim/warlock"
)

const (
	/* ! Please keep constants ordered by ID ! */

	Experiment800M       = 240925
	TyrsFall             = 241001
	RemnantsOfTheRed     = 241002
	HeartOfLight         = 241034
	AbandonedExperiment  = 241037
	SirDornelsDidgeridoo = 241038
	InfusionOfSouls      = 241039
	StiltzsStandard      = 241068
	LuckyDoubloon        = 241241
	HandOfRebornJustice  = 242310
)

func init() {
	core.AddEffectsToTest = false

	/* ! Please keep items ordered alphabetically ! */

	// https://www.wowhead.com/classic-ptr/item=241037/abandoned-experiment
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
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
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
		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: cdSpell,
		})
	})

	// https://www.wowhead.com/classic-ptr/item=242310/hand-of-reborn-justice
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

	// https://www.wowhead.com/classic-ptr/item=241034/heart-of-light
	// Use: Increases maximum health by 2500 for 20 sec. (2 Min Cooldown)
	core.NewSimpleStatDefensiveTrinketEffect(HeartOfLight, stats.Stats{stats.Health: 2500}, time.Second*20, time.Minute*2)

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
			// Explicitly lists that it does not work for Holy Shock in the tooltip https://www.wowhead.com/classic-ptr/item=241039/infusion-of-souls?spellModifier=462814
			classMask = paladin.ClassSpellMask_PaladinHarmfulGCDSpells ^ paladin.ClassSpellMask_PaladinHolyShock

		// https://www.wowhead.com/classic/spell=1232095/infusion-of-souls
		case proto.Class_ClassPriest:
			// Explicitly lists that it does not work for Penance in the tooltip https://www.wowhead.com/classic-ptr/item=241039/infusion-of-souls?spellModifier=440247
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

	// https://www.wowhead.com/classic-ptr/item=241241/lucky-doubloon
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
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
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

					// Delay the reset to simulate real player reaction time
					sim.AddPendingAction(&core.PendingAction{
						NextActionAt: sim.CurrentTime + character.Unit.ReactionTime,
						OnAction: func(sim *core.Simulation) {
							spell.CD.Reset()
						},
					})
				} else {
					buffAura.Deactivate(sim)
				}
			},
		})
		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: cdSpell,
		})
	})

	// https://www.wowhead.com/classic-ptr/item=241002/remnants-of-the-red
	// Equip: Dealing non-periodic Fire damage has a 10% chance to increase your Fire damage dealt by 10% for 20 sec. (Proc chance: 10%)
	core.NewItemEffect(RemnantsOfTheRed, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1231625},
			Label:    "Flames of the Red",
			Duration: time.Second * 20,
		}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire], 1.1)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:        "Remnants of the Red Trigger",
			Callback:    core.CallbackOnSpellHitDealt,
			Outcome:     core.OutcomeLanded,
			ProcMask:    core.ProcMaskSpellDamage,
			SpellSchool: core.SpellSchoolFire,
			ProcChance:  0.10,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})
	})

	// https://www.wowhead.com/classic-ptr/item=241038/sir-dornels-didgeridoo
	// Use: Playing the didgeridoo summons a random beast to your side, increasing your physical abilities for 30 sec.
	// May only be used while in combat. (2 Min Cooldown)
	// core.NewItemEffect(SirDornelsDidgeridoo, func(agent core.Agent) {
	// 	character := agent.GetCharacter()

	// 	actionID := core.ActionID{ItemID: SirDornelsDidgeridoo}
	// 	duration := time.Second * 30

	// 	bonusCrit := 10.0
	// 	bonusArp := 1000.0
	// 	bonusAgi := 140.0
	// 	bonusStr := 140.0
	// 	bonusHaste := 0.14
	// 	bonusAP := 280.0
	// 	bonusRAP := 308.0

	// 	// https://www.wowhead.com/classic-ptr/spell=1231884/accuracy-of-the-owl
	// 	owlAura := character.RegisterAura(core.Aura{
	// 		ActionID: actionID.WithTag(1231884),
	// 		Label:    "Sir Dornel's Didgeridoo - Owl",
	// 		Duration: duration,
	// 	}).AttachStatBuff(stats.MeleeCrit, bonusCrit).AttachStatBuff(stats.SpellCrit, bonusCrit)

	// 	// https://www.wowhead.com/classic-ptr/spell=1231894/ferocity-of-the-crocolisk
	// 	crocoliskBuff := character.RegisterAura(core.Aura{
	// 		ActionID: actionID.WithTag(1231894),
	// 		Label:    "Sir Dornel's Didgeridoo - Crocolisk",
	// 		Duration: duration,
	// 	}).AttachStatBuff(stats.ArmorPenetration, bonusArp)

	// 	// https://www.wowhead.com/classic-ptr/spell=1231888/agility-of-the-raptor
	// 	raptorBuff := character.RegisterAura(core.Aura{
	// 		ActionID: actionID.WithTag(1231888),
	// 		Label:    "Sir Dornel's Didgeridoo - Raptor",
	// 		Duration: duration,
	// 	}).AttachStatBuff(stats.Agility, bonusAgi)

	// 	// https://www.wowhead.com/classic-ptr/spell=1231883/strength-of-the-bear
	// 	bearBuff := character.RegisterAura(core.Aura{
	// 		ActionID: actionID.WithTag(1231883),
	// 		Label:    "Sir Dornel's Didgeridoo - Bear",
	// 		Duration: duration,
	// 	}).AttachStatBuff(stats.Strength, bonusStr)

	// 	// https://www.wowhead.com/classic-ptr/spell=1231886/speed-of-the-cat
	// 	catBuff := character.RegisterAura(core.Aura{
	// 		ActionID: actionID.WithTag(1231886),
	// 		Label:    "Sir Dornel's Didgeridoo - Cat",
	// 		Duration: duration,
	// 	}).AttachMultiplyAttackSpeed(&character.Unit, 1+bonusHaste)

	// 	// https://www.wowhead.com/classic-ptr/spell=1231891/power-of-the-gorilla
	// 	gorillaBuff := character.RegisterAura(core.Aura{
	// 		ActionID: actionID.WithTag(1231891),
	// 		Label:    "Sir Dornel's Didgeridoo - Gorilla",
	// 		Duration: duration,
	// 	}).AttachStatBuff(stats.AttackPower, bonusAP).AttachStatBuff(stats.RangedAttackPower, bonusRAP)

	// 	// https://www.wowhead.com/classic-ptr/spell=1231896/brilliance-of-mr-bigglesworth
	// 	bigglesworthBuff := character.RegisterAura(core.Aura{
	// 		ActionID: actionID.WithTag(1231896),
	// 		Label:    "Sir Dornel's Didgeridoo - Mr. Bigglesworth",
	// 		Duration: duration,
	// 	}).AttachStatBuff(
	// 		stats.MeleeCrit, bonusCrit/2.0,
	// 	).AttachStatBuff(
	// 		stats.SpellCrit, bonusCrit/2.0,
	// 	).AttachStatBuff(
	// 		stats.ArmorPenetration, bonusArp/2.0,
	// 	).AttachStatBuff(
	// 		stats.Agility, bonusAgi/2.0,
	// 	).AttachStatBuff(
	// 		stats.Strength, bonusStr/2.0,
	// 	).AttachMultiplyAttackSpeed(
	// 		&character.Unit, 1+bonusHaste/2.0,
	// 	).AttachStatBuff(
	// 		stats.AttackPower, bonusAP/2.0,
	// 	).AttachStatBuff(
	// 		stats.RangedAttackPower, bonusRAP/2.0,
	// 	)

	// 	cdSpell := character.RegisterSpell(core.SpellConfig{
	// 		ActionID: actionID,
	// 		ProcMask: core.ProcMaskEmpty,
	// 		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
	// 		Cast: core.CastConfig{
	// 			CD: core.Cooldown{
	// 				Timer:    character.NewTimer(),
	// 				Duration: time.Minute * 2,
	// 			},
	// 			SharedCD: core.Cooldown{
	// 				Timer:    character.GetOffensiveTrinketCD(),
	// 				Duration: duration,
	// 			},
	// 		},
	// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 			// TODO: How does this work? Equal chance for all? Or weighted?
	// 		},
	// 	})
	// 	character.AddMajorCooldown(core.MajorCooldown{
	// 		Type:  core.CooldownTypeDPS,
	// 		Spell: cdSpell,
	// 	})
	// })

	// https://www.wowhead.com/classic-ptr/item=241068/stiltzs-standard
	// Use: Throw down the Standard of Stiltz, increasing the maximum health of all nearby allies by 1000 for 20 sec. (2 Min Cooldown)
	core.NewSimpleStatDefensiveTrinketEffect(StiltzsStandard, stats.Stats{stats.Health: 1000}, time.Second*20, time.Minute*2)

	// https://www.wowhead.com/classic-ptr/item=241001/tyrs-fall
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

	core.NewItemEffect(Experiment800M, func(agent core.Agent) {
		character := agent.GetCharacter()

		numHits := character.Env.GetNumTargets()
		explosionSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1231607},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

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
			Flags:    core.SpellFlagOffensiveEquipment,
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

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	core.AddEffectsToTest = true
}
