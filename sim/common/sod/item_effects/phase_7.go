package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	BulwarkOfIre                 = 235868
	OlReliable                   = 235891
	DoomsayersDemise             = 235894
	TunicOfUndeadSlaying         = 236707
	BreastplateOfUndeadSlaying   = 236708
	ChestguardOfUndeadSlaying    = 236709
	WristguardsOfUndeadSlaying   = 236710
	WristwrapsOfUndeadSlaying    = 236711
	BracersOfUndeadSlaying       = 236712
	HandwrapsOfUndeadSlaying     = 236713
	GauntletsOfUndeadSlaying     = 236714
	HandguardsOfUndeadSlaying    = 236715
	BracersOfUndeadCleansing     = 236716
	GlovesOfUndeadCleansing      = 236717
	RobeofUndeadCleansing        = 236718
	BracersOfUndeadWardingCloth  = 236722
	GlovesOfUndeadWarding        = 236723
	RobeOfUndeadWarding          = 236724
	WristwrapsOfUndeadCleansing  = 236725
	HandwrapsOfUndeadCleansing   = 236726
	TunicOfUndeadCleansing       = 236727
	WristwrapsOfUndeadWarding    = 236731
	HandwrapsOfUndeadWarding     = 236732
	TunicOfUndeadWarding         = 236733
	WristguardsOfUndeadCleansing = 236734
	HandguardsOfUndeadCleansing  = 236735
	ChestguardOfUndeadCleansing  = 236736
	WristguardsOfUndeadWarding   = 236737
	HandguardsOfUndeadWarding    = 236738
	ChestguardOfUndeadWarding    = 236739
	BracersOfUndeadWardingPlate  = 236746
	GauntletsOfUndeadWarding     = 236747
	BreastplateOfUndeadWarding   = 236748
	BladeOfInquisition           = 237512
	TheHungeringCold             = 236341
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=236341/the-hungering-cold
	// Equip: Gives you a 2% chance to get an extra attack on the same target after dealing damage with your weapon.
	core.NewItemEffect(TheHungeringCold, func(agent core.Agent) {
		character := agent.GetCharacter()
		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "The Hungering Cold Trigger",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskMelee,
			ProcChance: 0.02,
			ICD:        time.Millisecond * 200,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				character.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 1223010}, spell)
			},
		})
		character.ItemSwap.RegisterProc(TheHungeringCold, aura)
	})

	// https://www.wowhead.com/classic/item=237512/blade-of-inquisition
	// Equip: Chance on hit to Increase your Strength by 250 and movement speed by 15% for 15 sec. (15s cooldown)
	// TODO: Verify proc chance, 1ppm for now
	core.NewItemEffect(BladeOfInquisition, func(agent core.Agent) {
		character := agent.GetCharacter()

		dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(BladeOfInquisition, 1.0, 0)

		buffAura := character.NewTemporaryStatsAura("Scarlet Inquisition", core.ActionID{SpellID: 1223342}, stats.Stats{stats.Strength: 250}, time.Second*15)

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Blade of Inquisition Trigger",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			ICD:               time.Second * 15,
			DPM:               dpm,
			DPMProcCheck:      core.DPMProc,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterProc(BladeOfInquisition, triggerAura)
	})

	// https://www.wowhead.com/classic/item=235894/doomsayers-demise
	// Equip: Periodic shadow effects have a chance to apply Creeping Darkness up to 5 times.
	// Spells which deal direct Shadow damage detonate this effect, dealing 100 damage per stack. (1.5s cooldown)
	core.NewItemEffect(DoomsayersDemise, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 1219020}

		debuffs := character.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
			return unit.RegisterAura(core.Aura{
				ActionID:  actionID,
				Label:     "Creeping Darkness",
				MaxStacks: 5,
				Duration:  time.Second * 30,
			})
		})

		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1219024},
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty, // Seems to be considered non-harmful
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				debuff := debuffs.Get(target)
				spell.CalcAndDealDamage(sim, target, float64(100*debuff.GetStacks()), spell.OutcomeMagicCrit)
				debuff.Deactivate(sim)
			},
		})

		icds := make(map[int32]core.Cooldown, len(character.Env.Encounter.TargetUnits))
		for _, target := range character.Env.Encounter.TargetUnits {
			icds[target.UnitIndex] = core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Millisecond * 1500,
			}
		}

		core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Creeping Darkness Trigger",
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if icd := icds[result.Target.UnitIndex]; icd.IsReady(sim) && spell.ProcMask.Matches(core.ProcMaskSpellDamage) && spell.SpellSchool.Matches(core.SpellSchoolShadow) {
					debuff := debuffs.Get(result.Target)
					debuff.Activate(sim)
					debuff.AddStack(sim)
					icd.Use(sim)
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Damage > 0 && spell.ProcMask.Matches(core.ProcMaskSpellDamage) && spell.SpellSchool.Matches(core.SpellSchoolShadow) && debuffs.Get(result.Target).IsActive() {
					damageSpell.Cast(sim, result.Target)
				}
			},
		}))
	})

	// https://www.wowhead.com/classic/item=235891/ol-reliable
	// Use: Smash the corpse of a fallen friend or foe, dealing 588 damage to nearby enemies. (2 Min Cooldown)
	core.NewItemEffect(OlReliable, func(agent core.Agent) {
		character := agent.GetCharacter()

		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1219043},
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, 588, spell.OutcomeMeleeSpecialHitAndCrit)
				}
			},
		})

		character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: OlReliable},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagAPL,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				damageSpell.Cast(sim, target)
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Other
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=235868/bulwark-of-ire
	// Deal 100 Shadow damage to melee attackers.
	// Causes twice as much threat as damage dealt.
	core.NewItemEffect(BulwarkOfIre, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThornsDamage += 100

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: BulwarkOfIre},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 2,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 100, spell.OutcomeMagicHit)
			},
		})

		aura := core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Splintered Shield",
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
					procSpell.Cast(sim, spell.Unit)
				}
			},
		}))

		character.ItemSwap.RegisterProc(BulwarkOfIre, aura)
	})

	// Cloth Sets

	// https://www.wowhead.com/classic/item=236718/robe-of-undead-cleansing/
	// Equip: Increases damage done to Undead by magical spells and effects by up to 65.
	core.NewMobTypeSpellPowerEffect(RobeofUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 65)
	// https://www.wowhead.com/classic/item=236716/bracers-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 35.
	core.NewMobTypeSpellPowerEffect(BracersOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 35)
	// https://www.wowhead.com/classic/item=236717/gloves-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 48.
	core.NewMobTypeSpellPowerEffect(GlovesOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 48)

	// https://www.wowhead.com/classic/item=236724/robe-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(RobeOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 26)
	// https://www.wowhead.com/classic/item=236722/bracers-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(BracersOfUndeadWardingCloth, []proto.MobType{proto.MobType_MobTypeUndead}, 26)
	// https://www.wowhead.com/classic/item=236723/gloves-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(GlovesOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 26)

	// Leather sets

	// https://www.wowhead.com/classic/item=236727/tunic-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 65.
	core.NewMobTypeSpellPowerEffect(TunicOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 65)
	// https://www.wowhead.com/classic/item=236725/wristwraps-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 35.
	core.NewMobTypeSpellPowerEffect(WristwrapsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 35)
	// https://www.wowhead.com/classic/item=236726/handwraps-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 48.
	core.NewMobTypeSpellPowerEffect(HandwrapsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 48)

	// https://www.wowhead.com/classic/item=236707/tunic-of-undead-slaying
	// Equip: +108 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(TunicOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 108)
	// https://www.wowhead.com/classic/item=236711/wristwraps-of-undead-slaying
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristwrapsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 60)
	// https://www.wowhead.com/classic/item=236713/handwraps-of-undead-slaying
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandwrapsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236733/tunic-of-undead-warding
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(TunicOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)
	// https://www.wowhead.com/classic/item=236731/wristwraps-of-undead-warding
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristwrapsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)
	// https://www.wowhead.com/classic/item=236732/handwraps-of-undead-warding
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandwrapsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// Mail

	// https://www.wowhead.com/classic/item=236736/chestguard-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 65.
	core.NewMobTypeSpellPowerEffect(ChestguardOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 65)
	// https://www.wowhead.com/classic/item=236734/wristguards-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 35.
	core.NewMobTypeSpellPowerEffect(WristguardsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 35)
	// https://www.wowhead.com/classic/item=236735/handguards-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 48.
	core.NewMobTypeSpellPowerEffect(HandguardsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 48)

	// https://www.wowhead.com/classic/item=236709/chestguard-of-undead-slaying
	// Equip: +108 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(ChestguardOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 108)
	// https://www.wowhead.com/classic/item=236710/wristguards-of-undead-slaying
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristguardsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 60)
	// https://www.wowhead.com/classic/item=236715/handguards-of-undead-slaying
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandguardsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236739/chestguard-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(ChestguardOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 26)
	// https://www.wowhead.com/classic/item=236737/wristguards-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(WristguardsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 26)
	// https://www.wowhead.com/classic/item=236738/handguards-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(HandguardsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 26)

	// Plate

	// https://www.wowhead.com/classic/item=236708/breastplate-of-undead-slaying
	// Equip: +108 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BreastplateOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 108)
	// https://www.wowhead.com/classic/item=236712/bracers-of-undead-slaying
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BracersOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 60)
	// https://www.wowhead.com/classic/item=236714/gauntlets-of-undead-slaying
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(GauntletsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236748/breastplate-of-undead-warding
	// +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BreastplateOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)
	// https://www.wowhead.com/classic/item=236746/bracers-of-undead-warding
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BracersOfUndeadWardingPlate, []proto.MobType{proto.MobType_MobTypeUndead}, 45)
	// https://www.wowhead.com/classic/item=236747/gauntlets-of-undead-warding
	//  Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(GauntletsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	core.AddEffectsToTest = true
}
