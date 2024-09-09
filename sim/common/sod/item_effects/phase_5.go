package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	Heartstriker               = 230253
	DrakeTalonCleaver          = 230271 // 19353
	JekliksCrusher             = 230911
	HaldberdOfSmiting          = 230991
	TigulesHarpoon             = 231272
	GrileksCarver              = 231273
	PitchforkOfMadness         = 231277
	Stormwrath                 = 231387
	WrathOfWray                = 231779
	LightningsCell             = 231784
	GrileksCarverBloodied      = 231846
	TigulesHarpoonBloodied     = 231849
	JekliksCrusherBloodied     = 231861
	PitchforkOfMadnessBloodied = 231864
	HaldberdOfSmitingBloodied  = 231870
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=230271/drake-talon-cleaver
	// Chance on hit: Delivers a fatal wound for 300 damage.
	// Original proc rate 1.0 increased to approximately 1.60 in SoD phase 5
	itemhelpers.CreateWeaponCoHProcDamage(DrakeTalonCleaver, "Drake Talon Cleaver", 1.6, 467167, core.SpellSchoolPhysical, 300, 0, 0.0, core.DefenseTypeMelee) // TBD confirm 1 ppm in SoD

	// https://www.wowhead.com/classic/item=231273/grileks-carver
	// +141 Attack Power when fighting Dragonkin.
	core.NewItemEffect(GrileksCarver, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			character.PseudoStats.MobTypeAttackPower += 141
		}
	})
	core.NewItemEffect(GrileksCarverBloodied, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			character.PseudoStats.MobTypeAttackPower += 141
		}
	})

	// https://www.wowhead.com/classic/item=230991/halberd-of-smiting
	// Equip: Chance to decapitate the target on a melee swing, causing 452 to 676 damage.
	itemhelpers.CreateWeaponEquipProcDamage(HaldberdOfSmiting, "Halberd of Smiting", 2.1, 467819, core.SpellSchoolPhysical, 452, 224, 0.0, core.DefenseTypeMelee)         // Works as phantom strike
	itemhelpers.CreateWeaponEquipProcDamage(HaldberdOfSmitingBloodied, "Halberd of Smiting", 2.1, 467819, core.SpellSchoolPhysical, 452, 224, 0.0, core.DefenseTypeMelee) // Works as phantom strike

	// https://www.wowhead.com/classic/item=230911/jekliks-crusher
	// Chance on hit: Wounds the target for 200 to 220 damage.
	// Original proc rate 4.0 lowered to 1.5 in SoD phase 5
	itemhelpers.CreateWeaponCoHProcDamage(JekliksCrusher, "Jeklik's Crusher", 1.5, 467642, core.SpellSchoolPhysical, 200, 20, 0.0, core.DefenseTypeMelee)
	itemhelpers.CreateWeaponCoHProcDamage(JekliksCrusherBloodied, "Jeklik's Crusher", 1.5, 467642, core.SpellSchoolPhysical, 200, 20, 0.0, core.DefenseTypeMelee)

	// https://www.wowhead.com/classic/item=231277/pitchfork-of-madness
	// +141 Attack Power when fighting Demons.
	core.NewItemEffect(PitchforkOfMadness, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			character.PseudoStats.MobTypeAttackPower += 141
		}
	})
	core.NewItemEffect(PitchforkOfMadnessBloodied, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			character.PseudoStats.MobTypeAttackPower += 141
		}
	})

	// https://www.wowhead.com/classic/item=231387/stormwrath-sanctified-shortblade-of-the-galefinder
	// Equip: Damaging non-periodic spells have a chance to blast up to 3 targets for 181 to 229.
	// (Proc chance: 10%, 100ms cooldown)
	core.NewItemEffect(Stormwrath, func(agent core.Agent) {
		character := agent.GetCharacter()

		maxHits := int(min(3, character.Env.GetNumTargets()))
		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 468670},
			SpellSchool:      core.SpellSchoolNature,
			DefenseType:      core.DefenseTypeMagic,
			ProcMask:         core.ProcMaskEmpty,
			BonusCoefficient: 0.15,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for numHits := 0; numHits < maxHits; numHits++ {
					spell.CalcAndDealDamage(sim, target, sim.Roll(180, 230), spell.OutcomeMagicHitAndCrit)
					target = character.Env.NextTargetUnit(target)
				}
			},
		})

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Millisecond * 100,
		}
		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Chain Lightning (Stormwrath)",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: .10,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				if !icd.IsReady(sim) {
					return
				}
				procSpell.Cast(sim, result.Target)
				icd.Use(sim)
			},
		})
	})

	// https://www.wowhead.com/classic/item=231272/tigules-harpoon
	// +99 Attack Power when fighting Beasts.
	core.NewItemEffect(TigulesHarpoon, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.MobTypeAttackPower += 99
		}
	})
	core.NewItemEffect(TigulesHarpoonBloodied, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.MobTypeAttackPower += 99
		}
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=231784/lightnings-cell
	// You gain a charge of Gathering Storm each time you cause a damaging spell critical strike.
	// When you reach 3 charges of Gathering Storm, they will release, firing an Unleashed Storm for 277 to 323 damage.
	// Gathering Storm cannot be gained more often than once every 2 sec. (2s cooldown)
	core.NewItemEffect(LightningsCell, func(agent core.Agent) {
		character := agent.GetCharacter()

		unleashedStormSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 468782},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(277, 323), spell.OutcomeMagicHitAndCrit)
			},
		})

		chargeAura := character.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 468780},
			Label:     "Lightning's Cell",
			Duration:  core.NeverExpires,
			MaxStacks: 3,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				if aura.GetStacks() == aura.MaxStacks {
					unleashedStormSpell.Cast(sim, aura.Unit.CurrentTarget)
					aura.Deactivate(sim)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Lightning's Cell Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeCrit,
			ProcMask: core.ProcMaskSpellDamage,
			ICD:      time.Millisecond * 2000,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				chargeAura.Activate(sim)
				chargeAura.AddStack(sim)
			},
		})
	})

	// https://www.wowhead.com/classic/item=230253/hearstriker
	// Equip: 2% chance on ranged hit to gain 1 extra attack. (Proc chance: 1%, 1s cooldown) // obviously something wrong here lol
	core.NewItemEffect(Heartstriker, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.AutoSwingRanged {
			return
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Heartstrike",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskRanged,
			ProcChance:        0.02,
			ICD:               time.Second * 1,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,

			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				spell.Unit.AutoAttacks.ExtraRangedAttack(sim, 1, core.ActionID{SpellID: 461164})
			},
		})
	})

	core.NewSimpleStatOffensiveTrinketEffect(WrathOfWray, stats.Stats{stats.Strength: 92}, time.Second*20, time.Minute*2)

	core.AddEffectsToTest = true
}
