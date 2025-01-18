package item_effects

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	BulwarkOfIre                     = 235868
	OlReliable                       = 235891
	TunicOfUndeadSlaying             = 236707
	BreastplateOfUndeadSlaying       = 236708
	ChestguardOfUndeadSlaying        = 236709
	WristguardsOfUndeadSlaying       = 236710
	WristwrapsOfUndeadSlaying        = 236711
	BracersOfUndeadSlaying           = 236712
	HandwrapsOfUndeadSlaying         = 236713
	GauntletsOfUndeadSlaying         = 236714
	HandguardsOfUndeadSlaying        = 236715
	BracersOfUndeadCleansing         = 236716
	GlovesOfUndeadCleansing          = 236717
	RobeofUndeadCleansing            = 236718
	BracersOfUndeadPurificationCloth = 236719
	GlovesOfUndeadPurification       = 236720
	RobeOfUndeadPurification         = 236721
	BracersOfUndeadWardingCloth      = 236722
	GlovesOfUndeadWarding            = 236723
	RobeOfUndeadWarding              = 236724
	WristwrapsOfUndeadCleansing      = 236725
	HandwrapsOfUndeadCleansing       = 236726
	TunicOfUndeadCleansing           = 236727
	WristwrapsOfUndeadWarding        = 236731
	HandwrapsOfUndeadWarding         = 236732
	TunicOfUndeadWarding             = 236733
	WristguardsOfUndeadCleansing     = 236734
	HandguardsOfUndeadCleansing      = 236735
	ChestguardOfUndeadCleansing      = 236736
	WristguardsOfUndeadWarding       = 236737
	HandguardsOfUndeadWarding        = 236738
	ChestguardOfUndeadWarding        = 236739
	WristguardsOfUndeadPurification  = 236740
	HandguardsOfUndeadPurification   = 236741
	ChestguardOfUndeadPurification   = 236742
	BracersOfUndeadPurificationPlate = 236743
	GauntletsOfUndeadPurification    = 236744
	BreastplateOfUndeadPurification  = 236745
	BracersOfUndeadWardingPlate      = 236746
	GauntletsOfUndeadWarding         = 236747
	BreastplateOfUndeadWarding       = 236748
	BladeOfInquisition               = 237512

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
	// TODO: PPM assumed and needs testing
	core.NewItemEffect(BladeOfInquisition, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(BladeOfInquisition)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		buffAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1223342},
			Label:    "Scarlet Inquisition",
			Duration: time.Second * 15,
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Blade of Inquisition",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			ICD:               time.Second * 15,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if ppmm.Proc(sim, spell.ProcMask, "Scarlet Inquisition") {
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
	core.NewItemEffect(AspirantsSealOfTheDawnDamage, sanctifiedDamageEffect(1219539, 0.83))
	core.NewItemEffect(InitiatesSealOfTheDawnDamage, sanctifiedDamageEffect(1223348, 2.92))
	core.NewItemEffect(SquiresSealOfTheDawnDamage, sanctifiedDamageEffect(1223349, 4.17))
	core.NewItemEffect(KnightsSealOfTheDawnDamage, sanctifiedDamageEffect(1223350, 6.67))
	core.NewItemEffect(TemplarsSealOfTheDawnDamage, sanctifiedDamageEffect(1223351, 8.33))
	core.NewItemEffect(ChampionsSealOfTheDawnDamage, sanctifiedDamageEffect(1223352, 12.08))
	core.NewItemEffect(VanguardsSealOfTheDawnDamage, sanctifiedDamageEffect(1223353, 14.17))
	core.NewItemEffect(CrusadersSealOfTheDawnDamage, sanctifiedDamageEffect(1223354, 18.75))
	core.NewItemEffect(CommandersSealOfTheDawnDamage, sanctifiedDamageEffect(1223355, 21.67))
	core.NewItemEffect(HighlordsSSealOfTheDawnDamage, sanctifiedDamageEffect(1223357, 25.0))

	// https://www.wowhead.com/classic/item=236383/squires-seal-of-the-dawn
	core.NewItemEffect(AspirantsSealOfTheDawnHealing, sanctifiedHealingEffect(1219548, 0.83))
	core.NewItemEffect(InitiatesSealOfTheDawnHealing, sanctifiedHealingEffect(1223379, 2.92))
	core.NewItemEffect(SquiresSealOfTheDawnHealing, sanctifiedHealingEffect(1223380, 4.17))
	core.NewItemEffect(KnightsSealOfTheDawnHealing, sanctifiedHealingEffect(1223381, 6.67))
	core.NewItemEffect(TemplarsSealOfTheDawnHealing, sanctifiedHealingEffect(1223382, 8.33))
	core.NewItemEffect(ChampionsSealOfTheDawnHealing, sanctifiedHealingEffect(1223383, 12.08))
	core.NewItemEffect(VanguardsSealOfTheDawnHealing, sanctifiedHealingEffect(1223384, 14.17))
	core.NewItemEffect(CrusadersSealOfTheDawnHealing, sanctifiedHealingEffect(1223385, 18.75))
	core.NewItemEffect(CommandersSealOfTheDawnHealing, sanctifiedHealingEffect(1223386, 21.67))
	core.NewItemEffect(HighlordsSSealOfTheDawnHealing, sanctifiedHealingEffect(1223387, 25.0))

	// https://www.wowhead.com/classic/item=236394/squires-seal-of-the-dawn
	core.NewItemEffect(AspirantsSealOfTheDawnTanking, sanctifiedTankingEffect(1220514, 2.08, 0.83))
	core.NewItemEffect(InitiatesSealOfTheDawnTanking, sanctifiedTankingEffect(1223367, 2.92, 2.92))
	core.NewItemEffect(SquiresSealOfTheDawnTanking, sanctifiedTankingEffect(1223368, 3.33, 4.17))
	core.NewItemEffect(KnightsSealOfTheDawnTanking, sanctifiedTankingEffect(1223370, 3.75, 6.67))
	core.NewItemEffect(TemplarsSealOfTheDawnTanking, sanctifiedTankingEffect(1223371, 4.17, 8.33))
	core.NewItemEffect(ChampionsSealOfTheDawnTanking, sanctifiedTankingEffect(1223372, 5.0, 12.08))
	core.NewItemEffect(VanguardsSealOfTheDawnTanking, sanctifiedTankingEffect(1223373, 5.42, 14.17))
	core.NewItemEffect(CrusadersSealOfTheDawnTanking, sanctifiedTankingEffect(1223374, 6.25, 18.75))
	core.NewItemEffect(CommandersSealOfTheDawnTanking, sanctifiedTankingEffect(1223375, 6.67, 21.67))
	core.NewItemEffect(HighlordsSSealOfTheDawnTanking, sanctifiedTankingEffect(1223376, 7.08, 25.0))

	///////////////////////////////////////////////////////////////////////////
	//                                 Other
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=236716/bracers-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(BracersOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 26)

	// https://www.wowhead.com/classic/item=236719/bracers-of-undead-purification
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(BracersOfUndeadPurificationCloth, []proto.MobType{proto.MobType_MobTypeUndead}, 26)

	// https://www.wowhead.com/classic/item=236743/bracers-of-undead-purification
	// +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BracersOfUndeadPurificationPlate, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236712/bracers-of-undead-slaying
	// +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BracersOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236722/bracers-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 26.
	core.NewMobTypeSpellPowerEffect(BracersOfUndeadWardingCloth, []proto.MobType{proto.MobType_MobTypeUndead}, 26)

	// https://www.wowhead.com/classic/item=236746/bracers-of-undead-warding
	// +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BracersOfUndeadWardingPlate, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236745/breastplate-of-undead-purification
	// +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BreastplateOfUndeadPurification, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236708/breastplate-of-undead-slaying
	// +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BreastplateOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236748/breastplate-of-undead-warding
	// +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(BreastplateOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

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

	// https://www.wowhead.com/classic/item=236736/chestguard-of-undead-cleansing
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(ChestguardOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236742/chestguard-of-undead-purification
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(ChestguardOfUndeadPurification, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236709/chestguard-of-undead-slaying
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(ChestguardOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236739/chestguard-of-undead-warding
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(ChestguardOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236744/gauntlets-of-undead-purification
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(GauntletsOfUndeadPurification, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236714/gauntlets-of-undead-slaying
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(GauntletsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236738/handguards-of-undead-warding
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(GauntletsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236717/gloves-of-undead-cleansing
	// Equip: Increases damage done to Undead by magical spells and effects by up to 35.
	core.NewMobTypeSpellPowerEffect(GlovesOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 35)

	// https://www.wowhead.com/classic/item=236720/gloves-of-undead-purification
	// Equip: Increases damage done to Undead by magical spells and effects by up to 35.
	core.NewMobTypeSpellPowerEffect(GlovesOfUndeadPurification, []proto.MobType{proto.MobType_MobTypeUndead}, 35)

	// https://www.wowhead.com/classic/item=236723/gloves-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 35.
	core.NewMobTypeSpellPowerEffect(GlovesOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 35)

	// https://www.wowhead.com/classic/item=236735/handguards-of-undead-cleansing
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandguardsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236741/handguards-of-undead-purification
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandguardsOfUndeadPurification, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236715/handguards-of-undead-slaying
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandguardsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236738/handguards-of-undead-warding
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandguardsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236726/handwraps-of-undead-cleansing
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandwrapsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236713/handwraps-of-undead-slaying
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandwrapsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236732/handwraps-of-undead-warding
	// Equip: +60 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(HandwrapsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 60)

	// https://www.wowhead.com/classic/item=236718/robe-of-undead-cleansing/
	// Equip: Increases damage done to Undead by magical spells and effects by up to 48.
	core.NewMobTypeSpellPowerEffect(RobeofUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 48)

	// https://www.wowhead.com/classic/item=236721/robe-of-undead-purification
	// Equip: Increases damage done to Undead by magical spells and effects by up to 48.
	core.NewMobTypeSpellPowerEffect(RobeOfUndeadPurification, []proto.MobType{proto.MobType_MobTypeUndead}, 48)

	// https://www.wowhead.com/classic/item=236724/robe-of-undead-warding
	// Equip: Increases damage done to Undead by magical spells and effects by up to 48.
	core.NewMobTypeSpellPowerEffect(RobeOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 48)

	// https://www.wowhead.com/classic/item=236727/tunic-of-undead-cleansing
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(TunicOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236707/tunic-of-undead-slaying
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(TunicOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236733/tunic-of-undead-warding
	// Equip: +81 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(TunicOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 81)

	// https://www.wowhead.com/classic/item=236734/wristguards-of-undead-cleansing
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristguardsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236740/wristguards-of-undead-purification
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristguardsOfUndeadPurification, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236710/wristguards-of-undead-slaying
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristguardsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236737/wristguards-of-undead-warding
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristguardsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236725/wristwraps-of-undead-cleansing
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristwrapsOfUndeadCleansing, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236711/wristwraps-of-undead-slaying
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristwrapsOfUndeadSlaying, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	// https://www.wowhead.com/classic/item=236732/handwraps-of-undead-warding
	// Equip: +45 Attack Power when fighting Undead.
	core.NewMobTypeAttackPowerEffect(WristwrapsOfUndeadWarding, []proto.MobType{proto.MobType_MobTypeUndead}, 45)

	core.AddEffectsToTest = true
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your damage by X% and your health by X% for each piece of Sanctified armor equipped.
func sanctifiedDamageEffect(spellID int32, percentIncrease float64) core.ApplyEffect {
	return func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.PseudoStats.SanctifiedBonus == 0 {
			return
		}

		sanctifiedBonus := math.Min(12, float64(character.PseudoStats.SanctifiedBonus))
		multiplier := 1.0 + percentIncrease/100.0*sanctifiedBonus
		healthDep := character.NewDynamicMultiplyStat(stats.Health, multiplier)

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label:      "Seal of the Dawn (Damage)",
			ActionID:   core.ActionID{SpellID: spellID},
			BuildPhase: core.CharacterBuildPhaseGear,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.EnableDynamicStatDep(healthDep)
				} else {
					character.EnableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.DamageDealtMultiplier *= multiplier
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.DisableDynamicStatDep(healthDep)
				} else {
					character.DisableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.DamageDealtMultiplier /= multiplier
			},
		}))
	}
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your healing and shielding by X% and your health by X% for each piece of Sanctified armor equipped.
func sanctifiedHealingEffect(spellID int32, percentIncrease float64) core.ApplyEffect {
	return func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.PseudoStats.SanctifiedBonus == 0 {
			return
		}

		sanctifiedBonus := math.Min(12, float64(character.PseudoStats.SanctifiedBonus))
		multiplier := 1.0 + percentIncrease/100.0*sanctifiedBonus
		healthDep := character.NewDynamicMultiplyStat(stats.Health, multiplier)

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label:      "Seal of the Dawn (Healing)",
			ActionID:   core.ActionID{SpellID: spellID},
			BuildPhase: core.CharacterBuildPhaseGear,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.EnableDynamicStatDep(healthDep)
				} else {
					character.EnableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.HealingDealtMultiplier *= multiplier
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.DisableDynamicStatDep(healthDep)
				} else {
					character.DisableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.HealingDealtMultiplier /= multiplier
			},
		}))
	}
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your threat caused by X%, your damage by Y%, and your health by Y% for each piece of Sanctified armor equipped.
func sanctifiedTankingEffect(spellID int32, threatPercentIncrease float64, damageHealthPercentIncrease float64) core.ApplyEffect {
	return func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.PseudoStats.SanctifiedBonus == 0 {
			return
		}

		sanctifiedBonus := math.Min(12, float64(character.PseudoStats.SanctifiedBonus))
		damageHealthMultiplier := 1.0 + damageHealthPercentIncrease/100.0*sanctifiedBonus
		threatMultiplier := 1.0 + threatPercentIncrease/100.0*sanctifiedBonus
		healthDep := character.NewDynamicMultiplyStat(stats.Health, damageHealthMultiplier)

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label:      "Seal of the Dawn (Tanking)",
			ActionID:   core.ActionID{SpellID: spellID},
			BuildPhase: core.CharacterBuildPhaseGear,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.EnableDynamicStatDep(healthDep)
				} else {
					character.EnableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.ThreatMultiplier *= threatMultiplier
				character.PseudoStats.DamageDealtMultiplier *= damageHealthMultiplier
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.DisableDynamicStatDep(healthDep)
				} else {
					character.DisableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.ThreatMultiplier /= threatMultiplier
				character.PseudoStats.DamageDealtMultiplier /= damageHealthMultiplier
			},
		}))
	}
}

func UnholyMightAura(character *core.Character) *core.Aura {
	return character.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1220668},
		Label:    "Unholy Might",
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Strength, 400)
			character.PseudoStats.DamageTakenMultiplier *= 1.20
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Strength, -400)
			character.PseudoStats.DamageTakenMultiplier /= 1.20
		},
	})
}
