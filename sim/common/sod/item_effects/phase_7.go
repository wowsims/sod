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
	AtieshSpellPower             = 236398
	AtieshHealing                = 236399
	AtieshCastSpeed              = 236400
	AtieshSpellCrit              = 236401
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
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=236400/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshCastSpeed, func(agent core.Agent) {
		core.AtieshCastSpeedEffect(&agent.GetCharacter().Unit)
	})
	// https://www.wowhead.com/classic/item=236399/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshHealing, func(agent core.Agent) {
		core.AtieshHealingEffect(&agent.GetCharacter().Unit)
	})
	// https://www.wowhead.com/classic/item=236401/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshSpellCrit, func(agent core.Agent) {
		core.AtieshSpellCritEffect(&agent.GetCharacter().Unit)
	})
	// https://www.wowhead.com/classic/item=236398/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshSpellPower, func(agent core.Agent) {
		core.AtieshSpellPowerEffect(&agent.GetCharacter().Unit)
	})

	// https://www.wowhead.com/classic/item=237512/blade-of-inquisition
	// Equip: Chance on hit to Increase your Strength by 250 and movement speed by 15% for 15 sec. (15s cooldown)
	// TODO: Verify proc chance, 1ppm for now
	core.NewItemEffect(BladeOfInquisition, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(BladeOfInquisition)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 15,
		}

		buffAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1223342},
			Label:    "Scarlet Inquisition",
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.Strength, 250)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.Strength, -250)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Blade of Inquisition Trigger",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			ProcMask:          procMask,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if icd.IsReady(sim) && ppmm.Proc(sim, procMask, "Scarlet Inquisition") {
					icd.Use(sim)
					buffAura.Activate(sim)
				}
			},
		})
	})

	// https://www.wowhead.com/classic/item=235894/doomsayers-demise
	// Equip: Periodic shadow effects have a chance to apply Creeping Darkness up to 5 times.
	// Spells which deal direct Shadow damage detonate this effect, dealing 45 damage per stack. (1.5s cooldown)
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
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				debuff := debuffs.Get(target)
				spell.CalcAndDealDamage(sim, target, float64(45*debuff.GetStacks()), spell.OutcomeMagicHitAndCrit)
				debuff.Deactivate(sim)
			},
		})

		// TODO: Made up proc rate TBD
		procChance := 0.20
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Millisecond * 1500,
		}

		core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Creeping Darkness Trigger",
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskSpellDamage) && spell.SpellSchool.Matches(core.SpellSchoolShadow) && sim.Proc(procChance, "Creeping Darkness") {
					debuff := debuffs.Get(result.Target)
					debuff.Activate(sim)
					debuff.AddStack(sim)
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskSpellDamage) && spell.SpellSchool.Matches(core.SpellSchoolShadow) && icd.IsReady(sim) && debuffs.Get(result.Target).IsActive() {
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

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: OlReliable},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

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

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
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

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Splintered Shield",
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
					procSpell.Cast(sim, spell.Unit)
				}
			},
		}))
	})

	// https://www.wowhead.com/classic/item=236707/tunic-of-undead-slaying
	// Equip: +108 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(TunicOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 108)

	// https://www.wowhead.com/classic/item=236708/breastplate-of-undead-slaying
	// Equip: +108 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BreastplateOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 108)

	// https://www.wowhead.com/classic/item=236709/chestguard-of-undead-slaying
	// Equip: +108 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(ChestguardOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 108)

	// https://www.wowhead.com/classic/item=236710/wristguards-of-undead-slaying
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristguardsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236711/wristwraps-of-undead-slaying
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristwrapsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236712/bracers-of-undead-slaying
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BracersOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236713/handwraps-of-undead-slaying
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandwrapsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236714/gauntlets-of-undead-slaying
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(GauntletsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236715/handguards-of-undead-slaying
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandguardsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236716/bracers-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 35.
	core.NewMobTypeSpellPowerEffect(BracersOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 35)

	// https://www.wowhead.com/classic/item=236717/gloves-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 48.
	core.NewMobTypeSpellPowerEffect(GlovesOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 48)

	// https://www.wowhead.com/classic/item=236718/robe-of-undead-cleansing/
	// Equip: Increases damage done to Undead by magical spells and effects by up to 65.
	core.NewMobTypeSpellPowerEffect(RobeofUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 65)

	// https://www.wowhead.com/classic/item=236722/bracers-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(BracersOfUndeadWardingCloth, []proto.MobType{proto.MobType_MobTypeUndead}, 26)

	// https://www.wowhead.com/classic/item=236723/gloves-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(GlovesOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 26)

	// https://www.wowhead.com/classic/item=236724/robe-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(RobeOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 26)

	// https://www.wowhead.com/classic/item=236725/wristwraps-of-undead-cleansing
	// Equip: +35 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristwrapsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 35)

	// https://www.wowhead.com/classic/item=236726/handwraps-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 48.
	core.NewMobTypeSpellPowerEffect(HandwrapsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 48)

	// https://www.wowhead.com/classic/item=236727/tunic-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 65.
	core.NewMobTypeSpellPowerEffect(TunicOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 65)

	// https://www.wowhead.com/classic/item=236731/wristwraps-of-undead-warding
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristwrapsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236732/handwraps-of-undead-warding
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandwrapsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236733/tunic-of-undead-warding
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(TunicOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236734/wristguards-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 35.
	core.NewMobTypeSpellPowerEffect(WristguardsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 35)

	// https://www.wowhead.com/classic/item=236735/handguards-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 48.
	core.NewMobTypeSpellPowerEffect(HandguardsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 48)

	// https://www.wowhead.com/classic/item=236736/chestguard-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 65.
	core.NewMobTypeSpellPowerEffect(ChestguardOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 65)

	// https://www.wowhead.com/classic/item=236737/wristguards-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(WristguardsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 36)

	// https://www.wowhead.com/classic/item=236738/handguards-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(HandguardsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 26)

	// https://www.wowhead.com/classic/item=236739/chestguard-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(ChestguardOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 26)

	// https://www.wowhead.com/classic/item=236746/bracers-of-undead-warding
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BracersOfUndeadWardingPlate, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236747/gauntlets-of-undead-warding
	//  Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(GauntletsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236748/breastplate-of-undead-warding
	// +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BreastplateOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	core.AddEffectsToTest = true
}

func UnholyMightAura(character *core.Character) *core.Aura {
	return character.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1220668},
		Label:    "Unholy Might",
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Strength, 350)
			character.PseudoStats.DamageTakenMultiplier *= 1.05
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Strength, -350)
			character.PseudoStats.DamageTakenMultiplier /= 1.05
		},
	})
}
