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

	// Atiesh
	AtieshSpellPower = 236398
	AtieshHealing    = 236399
	AtieshCastSpeed  = 236400
	AtieshSpellCrit  = 236401

	// Seals of the Dawn
	AspirantsSealOfTheDawnDamage  = 236354
	InitiatesSealOfTheDawnDamage  = 236355
	SquiresSealOfTheDawnDamage    = 236356
	KnightsSealOfTheDawnDamage    = 236357
	TemplarsSealOfTheDawnDamage   = 236358
	ChampionsSealOfTheDawnDamage  = 236360
	VanguardsSealOfTheDawnDamage  = 236361
	CrusadersSealOfTheDawnDamage  = 236362
	CommandersSealOfTheDawnDamage = 236363
	HighlordsSSealOfTheDawnDamage = 236364

	AspirantsSealOfTheDawnHealing  = 236385
	InitiatesSealOfTheDawnHealing  = 236384
	SquiresSealOfTheDawnHealing    = 236383
	KnightsSealOfTheDawnHealing    = 236382
	TemplarsSealOfTheDawnHealing   = 236380
	ChampionsSealOfTheDawnHealing  = 236379
	VanguardsSealOfTheDawnHealing  = 236378
	CrusadersSealOfTheDawnHealing  = 236376
	CommandersSealOfTheDawnHealing = 236375
	HighlordsSSealOfTheDawnHealing = 236374

	AspirantsSealOfTheDawnTanking  = 236396
	InitiatesSealOfTheDawnTanking  = 236395
	SquiresSealOfTheDawnTanking    = 236394
	KnightsSealOfTheDawnTanking    = 236393
	TemplarsSealOfTheDawnTanking   = 236392
	ChampionsSealOfTheDawnTanking  = 236391
	VanguardsSealOfTheDawnTanking  = 236390
	CrusadersSealOfTheDawnTanking  = 236389
	CommandersSealOfTheDawnTanking = 236388
	HighlordsSSealOfTheDawnTanking = 236386
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
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=236356/squires-seal-of-the-dawn
	core.NewItemEffect(AspirantsSealOfTheDawnDamage, sanctifiedDamageEffect(1219539, 1.25))
	core.NewItemEffect(InitiatesSealOfTheDawnDamage, sanctifiedDamageEffect(1223348, 4.38))
	core.NewItemEffect(SquiresSealOfTheDawnDamage, sanctifiedDamageEffect(1223349, 6.25))
	core.NewItemEffect(KnightsSealOfTheDawnDamage, sanctifiedDamageEffect(1223350, 10.0))
	core.NewItemEffect(TemplarsSealOfTheDawnDamage, sanctifiedDamageEffect(1223351, 12.5))
	core.NewItemEffect(ChampionsSealOfTheDawnDamage, sanctifiedDamageEffect(1223352, 18.13))
	core.NewItemEffect(VanguardsSealOfTheDawnDamage, sanctifiedDamageEffect(1223353, 21.25))
	core.NewItemEffect(CrusadersSealOfTheDawnDamage, sanctifiedDamageEffect(1223354, 28.13))
	core.NewItemEffect(CommandersSealOfTheDawnDamage, sanctifiedDamageEffect(1223355, 32.5))
	core.NewItemEffect(HighlordsSSealOfTheDawnDamage, sanctifiedDamageEffect(1223357, 37.5))

	// https://www.wowhead.com/classic/item=236383/squires-seal-of-the-dawn
	core.NewItemEffect(AspirantsSealOfTheDawnHealing, sanctifiedHealingEffect(1219548, 1.25))
	core.NewItemEffect(InitiatesSealOfTheDawnHealing, sanctifiedHealingEffect(1223379, 4.38))
	core.NewItemEffect(SquiresSealOfTheDawnHealing, sanctifiedHealingEffect(1223380, 6.25))
	core.NewItemEffect(KnightsSealOfTheDawnHealing, sanctifiedHealingEffect(1223381, 10.0))
	core.NewItemEffect(TemplarsSealOfTheDawnHealing, sanctifiedHealingEffect(1223382, 12.5))
	core.NewItemEffect(ChampionsSealOfTheDawnHealing, sanctifiedHealingEffect(1223383, 18.13))
	core.NewItemEffect(VanguardsSealOfTheDawnHealing, sanctifiedHealingEffect(1223384, 21.25))
	core.NewItemEffect(CrusadersSealOfTheDawnHealing, sanctifiedHealingEffect(1223385, 28.13))
	core.NewItemEffect(CommandersSealOfTheDawnHealing, sanctifiedHealingEffect(1223386, 32.5))
	core.NewItemEffect(HighlordsSSealOfTheDawnHealing, sanctifiedHealingEffect(1223387, 37.5))

	// https://www.wowhead.com/classic/item=236394/squires-seal-of-the-dawn
	core.NewItemEffect(AspirantsSealOfTheDawnTanking, sanctifiedTankingEffect(1220514, 3.13, 1.25))
	core.NewItemEffect(InitiatesSealOfTheDawnTanking, sanctifiedTankingEffect(1223367, 4.38, 4.38))
	core.NewItemEffect(SquiresSealOfTheDawnTanking, sanctifiedTankingEffect(1223368, 5.0, 6.25))
	core.NewItemEffect(KnightsSealOfTheDawnTanking, sanctifiedTankingEffect(1223370, 5.63, 10))
	core.NewItemEffect(TemplarsSealOfTheDawnTanking, sanctifiedTankingEffect(1223371, 6.25, 12.5))
	core.NewItemEffect(ChampionsSealOfTheDawnTanking, sanctifiedTankingEffect(1223372, 7.5, 18.13))
	core.NewItemEffect(VanguardsSealOfTheDawnTanking, sanctifiedTankingEffect(1223373, 8.13, 21.25))
	// Updated tooltip for Crusader's actually says 21.83% for health but that's not consistent with others
	// and most likely a typo where the 8 and the 1 has swapped places.
	core.NewItemEffect(CrusadersSealOfTheDawnTanking, sanctifiedTankingEffect(1223374, 9.38, 28.13))
	core.NewItemEffect(CommandersSealOfTheDawnTanking, sanctifiedTankingEffect(1223375, 10.0, 32.5))
	core.NewItemEffect(HighlordsSSealOfTheDawnTanking, sanctifiedTankingEffect(1223376, 10.63, 37.5))

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

const MaxSanctifiedBonus = 8

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your damage by X% and your health by X% for each piece of Sanctified armor equipped.
func sanctifiedDamageEffect(spellID int32, percentIncrease float64) core.ApplyEffect {
	return func(agent core.Agent) {
		character := agent.GetCharacter()

		for _, unit := range getSanctifiedUnits(character) {
			sanctifiedBonus := int32(0)
			multiplier := 1.0
			healthDeps := buildSanctifiedHealthDeps(unit, percentIncrease)

			core.MakePermanent(unit.GetOrRegisterAura(core.Aura{
				Label:      "Seal of the Dawn (Damage)",
				ActionID:   core.ActionID{SpellID: spellID},
				BuildPhase: core.CharacterBuildPhaseGear,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					sanctifiedBonus = max(min(MaxSanctifiedBonus, character.PseudoStats.SanctifiedBonus), 0)
					multiplier = 1.0 + percentIncrease/100.0*float64(sanctifiedBonus)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
						aura.Unit.StatDependencyManager.EnableDynamicStatDep(healthDeps[sanctifiedBonus])
					} else {
						aura.Unit.EnableDynamicStatDep(sim, healthDeps[sanctifiedBonus])
					}

					aura.Unit.PseudoStats.DamageDealtMultiplier *= multiplier
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
						aura.Unit.StatDependencyManager.DisableDynamicStatDep(healthDeps[sanctifiedBonus])
					} else {
						aura.Unit.DisableDynamicStatDep(sim, healthDeps[sanctifiedBonus])
					}

					aura.Unit.PseudoStats.DamageDealtMultiplier /= multiplier
				},
			}))
		}
	}
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your healing and shielding by X% and your health by X% for each piece of Sanctified armor equipped.
func sanctifiedHealingEffect(spellID int32, percentIncrease float64) core.ApplyEffect {
	return func(agent core.Agent) {
		character := agent.GetCharacter()

		for _, unit := range getSanctifiedUnits(character) {
			sanctifiedBonus := int32(0)
			multiplier := 1.0
			healthDeps := buildSanctifiedHealthDeps(unit, percentIncrease)

			core.MakePermanent(unit.GetOrRegisterAura(core.Aura{
				Label:      "Seal of the Dawn (Healing)",
				ActionID:   core.ActionID{SpellID: spellID},
				BuildPhase: core.CharacterBuildPhaseGear,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					sanctifiedBonus = max(min(MaxSanctifiedBonus, character.PseudoStats.SanctifiedBonus), 0)
					multiplier = 1.0 + percentIncrease/100.0*float64(sanctifiedBonus)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
						aura.Unit.StatDependencyManager.EnableDynamicStatDep(healthDeps[sanctifiedBonus])
					} else {
						aura.Unit.EnableDynamicStatDep(sim, healthDeps[sanctifiedBonus])
					}

					aura.Unit.PseudoStats.HealingDealtMultiplier *= multiplier
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
						aura.Unit.StatDependencyManager.DisableDynamicStatDep(healthDeps[sanctifiedBonus])
					} else {
						aura.Unit.DisableDynamicStatDep(sim, healthDeps[sanctifiedBonus])
					}

					aura.Unit.PseudoStats.HealingDealtMultiplier /= multiplier
				},
			}))
		}
	}
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your threat caused by X%, your damage by Y%, and your health by Y% for each piece of Sanctified armor equipped.
func sanctifiedTankingEffect(spellID int32, threatPercentIncrease float64, damageHealthPercentIncrease float64) core.ApplyEffect {
	return func(agent core.Agent) {
		character := agent.GetCharacter()

		units := []*core.Unit{&character.Unit}
		if character.Class == proto.Class_ClassHunter || character.Class == proto.Class_ClassWarlock {
			for _, pet := range character.Pets {
				if pet.IsGuardian() {
					return
				}

				units = append(units, &pet.Unit)
			}
		}

		for _, unit := range units {
			sanctifiedBonus := int32(0)
			damageHealthMultiplier := 1.0
			threatMultiplier := 1.0
			healthDeps := buildSanctifiedHealthDeps(unit, damageHealthPercentIncrease)

			core.MakePermanent(unit.GetOrRegisterAura(core.Aura{
				Label:      "Seal of the Dawn (Tanking)",
				ActionID:   core.ActionID{SpellID: spellID},
				BuildPhase: core.CharacterBuildPhaseGear,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					sanctifiedBonus = max(min(MaxSanctifiedBonus, character.PseudoStats.SanctifiedBonus), 0)
					damageHealthMultiplier = 1.0 + damageHealthPercentIncrease/100.0*float64(sanctifiedBonus)
					threatMultiplier = 1.0 + threatPercentIncrease/100.0*float64(sanctifiedBonus)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
						aura.Unit.StatDependencyManager.EnableDynamicStatDep(healthDeps[sanctifiedBonus])
					} else {
						aura.Unit.EnableDynamicStatDep(sim, healthDeps[sanctifiedBonus])
					}

					aura.Unit.PseudoStats.ThreatMultiplier *= threatMultiplier
					aura.Unit.PseudoStats.DamageDealtMultiplier *= damageHealthMultiplier
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
						aura.Unit.StatDependencyManager.DisableDynamicStatDep(healthDeps[sanctifiedBonus])
					} else {
						aura.Unit.DisableDynamicStatDep(sim, healthDeps[sanctifiedBonus])
					}

					aura.Unit.PseudoStats.ThreatMultiplier /= threatMultiplier
					aura.Unit.PseudoStats.DamageDealtMultiplier /= damageHealthMultiplier
				},
			}))
		}
	}
}

// Gets all units that the Sanctified buff should apply to. This includes the player and Hunter/Warlock pets
func getSanctifiedUnits(character *core.Character) []*core.Unit {
	units := []*core.Unit{&character.Unit}
	if character.Class == proto.Class_ClassHunter || character.Class == proto.Class_ClassWarlock {
		for _, pet := range character.Pets {
			if pet.IsGuardian() {
				continue
			}

			units = append(units, &pet.Unit)
		}
	}

	return units
}

func buildSanctifiedHealthDeps(unit *core.Unit, percentIncrease float64) []*stats.StatDependency {
	healthDeps := []*stats.StatDependency{}
	for i := 0; i < MaxSanctifiedBonus+1; i++ {
		healthDeps = append(healthDeps, unit.NewDynamicMultiplyStat(stats.Health, 1.0+percentIncrease/100.0*float64(i)))
	}

	return healthDeps
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
