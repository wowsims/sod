package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	Truthbearer2H             = 229749
	Truthbearer1H             = 229806
	HammerOfTheLightbringer   = 230003
	TheUntamedBlade           = 230242
	DrakeTalonCleaver         = 230271
	ScrollsOfBlindingLight    = 230272
	JekliksCrusher            = 230911
	HaldberdOfSmiting         = 230991
	Stormwrath                = 231387
	WrathOfWray               = 231779
	LightningsCell            = 231784
	JekliksCrusherBloodied    = 231861
	HaldberdOfSmitingBloodied = 231870
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	itemhelpers.CreateWeaponEquipProcDamage(DrakeTalonCleaver, "Drake Talon Cleaver", 1.0, 467167, core.SpellSchoolPhysical, 300, 0, 0.0, core.DefenseTypeMelee) // TBD confirm 1 ppm in SoD

	itemhelpers.CreateWeaponEquipProcDamage(HaldberdOfSmiting, "Halberd of Smiting", 0.5, 467819, core.SpellSchoolPhysical, 452, 224, 0.0, core.DefenseTypeMelee)         // TBD does this work as phantom strike?, confirm 0.5 ppm in SoD
	itemhelpers.CreateWeaponEquipProcDamage(HaldberdOfSmitingBloodied, "Halberd of Smiting", 0.5, 467819, core.SpellSchoolPhysical, 452, 224, 0.0, core.DefenseTypeMelee) // TBD does this work as phantom strike?, confirm 0.5 ppm in SoD

	itemhelpers.CreateWeaponCoHProcDamage(JekliksCrusher, "Jeklik's Crusher", 4.0, 467642, core.SpellSchoolPhysical, 200, 20, 0.0, core.DefenseTypeMelee)
	itemhelpers.CreateWeaponCoHProcDamage(JekliksCrusherBloodied, "Jeklik's Crusher", 4.0, 467642, core.SpellSchoolPhysical, 200, 20, 0.0, core.DefenseTypeMelee)

	core.NewItemEffect(TheUntamedBlade, func(agent core.Agent) {
		character := agent.GetCharacter()

		strengthAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Untamed Fury",
			ActionID: core.ActionID{SpellID: 23719},
			Duration: time.Second * 8,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.Strength, 300)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.Strength, -300)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "The Untamed Blade (Strength)",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskMelee,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			PPM:               1, // Estimated based on data from WoW Armaments Discord
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				strengthAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(HammerOfTheLightbringer, func(agent core.Agent) {
		character := agent.GetCharacter()
		blazefuryTriggerAura(character, 465412, core.SpellSchoolHoly, 4)
		crusadersZealAura465414(character)
	})

	// https://www.wowhead.com/classic/item=229749/truthbearer
	// Chance on hit: Increases damage done by 15 and attack speed by 30% for 8 sec.
	// TODO: Proc rate assumed and needs testing
	core.NewItemEffect(Truthbearer1H, func(agent core.Agent) {
		character := agent.GetCharacter()
		crusadersZealAura465414(character)
	})

	// https://www.wowhead.com/classic/item=229806/truthbearer
	// Chance on hit: Increases damage done by 15 and attack speed by 30% for 8 sec.
	// TODO: Proc rate assumed and needs testing
	core.NewItemEffect(Truthbearer2H, func(agent core.Agent) {
		character := agent.GetCharacter()
		blazefuryTriggerAura(character, 465412, core.SpellSchoolHoly, 4)
		crusadersZealAura465414(character)
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

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=231784/lightnings-cell
	// You gain a charge of Gathering Storm each time you cause a damaging spell critical strike.
	// When you reach 3 charges of Gathering Storm, they will release, firing an Unleashed Storm for 277 to 323 damage.
	// Gathering Storm cannot be gained more often than once every 2.5 sec. (2.5s cooldown)
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
			ICD:      time.Millisecond * 2500,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				chargeAura.Activate(sim)
				chargeAura.AddStack(sim)
			},
		})
	})

	// https://www.wowhead.com/classic/item=231784/lightnings-cell
	//

	// Use: Energizes a Paladin with light, increasing melee attack speed by 25%
	// and spell casting speed by 33% for 20 sec. (2 Min Cooldown)
	core.NewItemEffect(ScrollsOfBlindingLight, func(agent core.Agent) {
		character := agent.GetCharacter()

		duration := time.Second * 20

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Blinding Light",
			ActionID: core.ActionID{SpellID: 467522},
			Duration: duration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.25)
				character.MultiplyCastSpeed(1.33)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.0/1.25)
				character.MultiplyCastSpeed(1.0 / 1.33)
			},
		})

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: ScrollsOfBlindingLight},
			Flags:    core.SpellFlagNoOnCastComplete,

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

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:     core.CooldownTypeDPS,
			Priority: core.CooldownPriorityDefault,
			Spell:    spell,
		})
	})

	core.NewSimpleStatOffensiveTrinketEffect(WrathOfWray, stats.Stats{stats.Strength: 92}, time.Second*20, time.Minute*2)

	core.AddEffectsToTest = true
}

func blazefuryTriggerAura(character *core.Character, spellID int32, spellSchool core.SpellSchool, damage int32) *core.Aura {

	procSpell := character.GetOrRegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: spellID},
		SpellSchool:      core.SpellSchoolHoly,
		DefenseType:      core.DefenseTypeMagic,
		ProcMask:         core.ProcMaskTriggerInstant,
		Flags:            core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 4, spell.OutcomeMagicCrit)
		},
	})

	return core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:              "Blazefury Trigger",
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		ProcMask:          core.ProcMaskMelee,
		SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
				procSpell.ProcMask = core.ProcMaskEmpty
			} else {
				procSpell.ProcMask = core.ProcMaskTriggerInstant
			}
			procSpell.Cast(sim, result.Target)
		},
	})
}

// https://www.wowhead.com/classic/spell=465414/crusaders-zeal
// Used by:
// - https://www.wowhead.com/classic/item=229806/truthbearer and
// - https://www.wowhead.com/classic/item=229749/truthbearer
func crusadersZealAura465414(character *core.Character) *core.Aura {
	procAura := character.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 465414},
		Label:    "Crusader's Zeal",
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			character.PseudoStats.BonusDamage += 15
			character.MultiplyAttackSpeed(sim, 1.30)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			character.PseudoStats.BonusDamage -= 15
			character.MultiplyAttackSpeed(sim, 1/1.30)
		},
	})

	return core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:              "Truthbearer (Crusader's Zeal)",
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		ProcMask:          core.ProcMaskMelee,
		SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
		PPM:               1, // TBD
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAura.Activate(sim)
		},
	})
}
