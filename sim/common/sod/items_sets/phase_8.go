package item_sets

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/hunter"
	"github.com/wowsims/sod/sim/rogue"
	"github.com/wowsims/sod/sim/shaman"
	"github.com/wowsims/sod/sim/warrior"
)

var ItemSetFallenRegality = core.NewItemSet(core.ItemSet{
	Name: "Fallen Regality",
	Bonuses: map[int32]core.ApplyEffect{
		// Damaging finishing moves have a 20% chance per combo point to restore 20 energy.
		// Flanking Strike's damage buff is increased by an additional 2% per stack. When striking from behind, your target takes 150% increased damage from Flanking Strike.
		// If Cleave hits fewer than its maximum number of targets, it deals 35% more damage for each unused bounce.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			aura := core.MakePermanent(character.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1232184},
				Label:    "Fallen Regality",
			}))

			switch character.Class {
			case proto.Class_ClassRogue:
				agent.(rogue.RogueAgent).GetRogue().ApplyFallenRegalityRogueBonus(aura)
			case proto.Class_ClassHunter:
				agent.(hunter.HunterAgent).GetHunter().ApplyFallenRegalityHunterBonus(aura)
			case proto.Class_ClassWarrior:
				agent.(warrior.WarriorAgent).GetWarrior().ApplyFallenRegalityWarriorBonus(aura)
			}
		},
	},
})

var ItemSetHackAndSmash = core.NewItemSet(core.ItemSet{
	Name: "Hack and Smash",
	Bonuses: map[int32]core.ApplyEffect{
		// Hunter - The damage increaes from Mercy's and Crimson Cleaver's effects are increased by 10%.
		// Shaman - The Fire and Nature damage increases from Mercy and Crimson Cleaver are increased by 10%.
		// Warrior - The damage increaes from Mercy's and Crimson Cleaver's effects are increased by 10%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			core.MakePermanent(character.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1234318},
				Label:    "Hack and Smash",
			}))

			switch character.Class {
			case proto.Class_ClassHunter:
				agent.(hunter.HunterAgent).GetHunter().ApplyHackAndSmashHunterBonus()
			case proto.Class_ClassShaman:
				agent.(shaman.ShamanAgent).GetShaman().ApplyHackAndSmashShamanBonus()
			case proto.Class_ClassWarrior:
				agent.(warrior.WarriorAgent).GetWarrior().ApplyHackAndSmashWarriorBonus()
			}
		},
	},
})

const (
	Deception                = 240922
	Duplicity                = 240923
)

// https://www.wowhead.com/classic/item-set=1956/tools-of-the-nathrezim
var ItemSetToolsOfTheNathrezim = core.NewItemSet(core.ItemSet{
	Name: "Tools of the Nathrezim",
	Bonuses: map[int32]core.ApplyEffect{
		// Duplicity and Deception's extra attacks now trigger 2 extra attacks.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			core.MakePermanent(character.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1231556},
				Label:    "Tools of the Nathrezim",
			}))

			// Duplicity
			isMhDuplicity  := character.GetMHWeapon().ID == Duplicity
			castTypeDuplicity := core.Ternary(isMhDuplicity, proto.CastType_CastTypeMainHand, proto.CastType_CastTypeOffHand)
			procMaskDuplicity := core.Ternary(isMhDuplicity, core.ProcMaskMeleeMHAuto, core.ProcMaskMeleeOHAuto)

			spellProcSetDuplicity := character.RegisterSpell(core.SpellConfig{
				ActionID:       core.ActionID{SpellID: 1231557},
				SpellSchool:    core.SpellSchoolPhysical,
				DefenseType:    core.DefenseTypeMelee,
				ProcMask:       procMaskDuplicity, // Normal Melee Attack Flag
				Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagSuppressWeaponProcs, // Cannot proc Oil, Poisons, and presumably Weapon Enchants or Procs
				CastType:       castTypeDuplicity,
	
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
	
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					if isMhDuplicity {
						baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
					} else {
						baseDamage := spell.Unit.AutoAttacks.OH().CalculateWeaponDamage(sim, spell.MeleeAttackPower()) // Avoid using OHWeaponDamage(sim, spell.MeleeAttackPower() to avoid 50% DW Penalty
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
					}
				},
			})

			character.NewRageMetrics(spellProcSetDuplicity.ActionID)
			spellProcSetDuplicity.ResourceMetrics = character.NewRageMetrics(spellProcSetDuplicity.ActionID)

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:          core.ActionID{SpellID: 1231556},
				Name:              "Tools of the Nathrezim (Duplicity)", 
				Callback:          core.CallbackOnSpellHitDealt,
				Outcome:           core.OutcomeLanded,
				ProcMask:          core.ProcMaskMelee,
				SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
				ProcChance:        0.02,
				ICD:               time.Millisecond * 100,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					spellProcSetDuplicity.Cast(sim, result.Target)
					spellProcSetDuplicity.Cast(sim, result.Target)
				},
			})

			// Deception
			isMhDeception  := character.GetMHWeapon().ID == Deception
			castTypeDeception := core.Ternary(isMhDeception, proto.CastType_CastTypeMainHand, proto.CastType_CastTypeOffHand)
			procMaskDeception := core.Ternary(isMhDeception, core.ProcMaskMeleeMHAuto, core.ProcMaskMeleeOHAuto)

			spellProcSetDeception := character.RegisterSpell(core.SpellConfig{
				ActionID:       core.ActionID{SpellID: 1231558},
				SpellSchool:    core.SpellSchoolPhysical,
				DefenseType:    core.DefenseTypeMelee,
				ProcMask:       procMaskDeception, // Normal Melee Attack Flag
				Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagSuppressWeaponProcs, // Cannot proc Oil, Poisons, and presumably Weapon Enchants or Procs
				CastType:       castTypeDeception,
	
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
	
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					if isMhDeception {
						baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
					} else {
						baseDamage := spell.Unit.AutoAttacks.OH().CalculateWeaponDamage(sim, spell.MeleeAttackPower()) // Avoid using OHWeaponDamage(sim, spell.MeleeAttackPower() to avoid 50% DW Penalty
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
					}
				},
			})

			character.NewRageMetrics(spellProcSetDeception.ActionID)
			spellProcSetDeception.ResourceMetrics = character.NewRageMetrics(spellProcSetDeception.ActionID)

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:          core.ActionID{SpellID: 1231556},
				Name:              "Tools of the Nathrezim (Deception)", 
				Callback:          core.CallbackOnSpellHitDealt,
				Outcome:           core.OutcomeLanded,
				ProcMask:          core.ProcMaskMelee,
				SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
				ProcChance:        0.02,
				ICD:               time.Millisecond * 100,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					spellProcSetDeception.Cast(sim, result.Target)
					spellProcSetDeception.Cast(sim, result.Target)
				},
			})
		},
	},
})
