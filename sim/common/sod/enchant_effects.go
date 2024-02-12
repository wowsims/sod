package sod

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	core.AddEffectsToTest = false

	// Fiery Blaze
	core.NewEnchantEffect(36, func(agent core.Agent) {
		character := agent.GetCharacter()
		procChance := 0.15

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 6296},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					damage := sim.Roll(9, 13)
					spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
				}

			},
		})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Fiery Blaze",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if sim.RandomFloat("Fiery Blaze") < procChance {
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffect(36, aura)
	})

	// Fiery Weapon
	// core.NewEnchantEffect(803, func(agent core.Agent) {
	// 	character := agent.GetCharacter()

	// 	procMask := character.GetProcMaskForEnchant(803)
	// 	ppmm := character.AutoAttacks.NewPPMManager(6.0, procMask)

	// 	procSpell := character.RegisterSpell(core.SpellConfig{
	// 		ActionID:    core.ActionID{SpellID: 13898},
	// 		SpellSchool: core.SpellSchoolFire,
	// 		ProcMask:    core.ProcMaskSpellDamage,

	// 		DamageMultiplier: 1,
	// 		CritMultiplier:   character.DefaultSpellCritMultiplier(),
	// 		ThreatMultiplier: 1,

	// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 			spell.CalcAndDealDamage(sim, target, 40, spell.OutcomeMagicHitAndCrit)
	// 		},
	// 	})

	// 	aura := character.GetOrRegisterAura(core.Aura{
	// 		Label:    "Fiery Weapon",
	// 		Duration: core.NeverExpires,
	// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
	// 			aura.Activate(sim)
	// 		},
	// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 			if !result.Landed() {
	// 				return
	// 			}

	// 			if ppmm.Proc(sim, spell.ProcMask, "Fiery Weapon") {
	// 				procSpell.Cast(sim, result.Target)
	// 			}
	// 		},
	// 	})

	// 	character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(803, 6.0, &ppmm, aura)
	// })

	// Minor Striking
	core.AddWeaponEffect(250, func(agent core.Agent, slot proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.MH()
		if slot == proto.ItemSlot_ItemSlotOffHand {
			w = agent.GetCharacter().AutoAttacks.OH()
		}
		w.BaseDamageMin += 1
		w.BaseDamageMax += 1
	})

	// Lesser Striking
	core.AddWeaponEffect(241, func(agent core.Agent, slot proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.MH()
		if slot == proto.ItemSlot_ItemSlotOffHand {
			w = agent.GetCharacter().AutoAttacks.OH()
		}
		w.BaseDamageMin += 2
		w.BaseDamageMax += 2
	})

	// Striking
	core.AddWeaponEffect(943, func(agent core.Agent, slot proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.MH()
		if slot == proto.ItemSlot_ItemSlotOffHand {
			w = agent.GetCharacter().AutoAttacks.OH()
		}
		w.BaseDamageMin += 3
		w.BaseDamageMax += 3
	})

	// Greater Striking
	// core.AddWeaponEffect(805, func(agent core.Agent, slot proto.ItemSlot) {
	// 	w := agent.GetCharacter().AutoAttacks.MH()
	// 	if slot == proto.ItemSlot_ItemSlotOffHand {
	// 		w = agent.GetCharacter().AutoAttacks.OH()
	// 	}
	// 	w.BaseDamageMin += 4
	// 	w.BaseDamageMax += 4
	// })

	// Superior Striking
	// core.AddWeaponEffect(1897, func(agent core.Agent, slot proto.ItemSlot) {
	// 	w := agent.GetCharacter().AutoAttacks.MH()
	// 	if slot == proto.ItemSlot_ItemSlotOffHand {
	// 		w = agent.GetCharacter().AutoAttacks.OH()
	// 	}
	// 	w.BaseDamageMin += 5
	// 	w.BaseDamageMax += 5
	// })

	// TODO: Crusader, Mongoose, and Executioner could also be modelled as AddWeaponEffect instead
	// ApplyCrusaderEffect will be applied twice if there is two weapons with this enchant.
	//   However, it will automatically overwrite one of them, so it should be ok.
	//   A single application of the aura will handle both mh and oh procs.
	// core.NewEnchantEffect(1900, func(agent core.Agent) {
	// 	character := agent.GetCharacter()

	// 	procMask := character.GetProcMaskForEnchant(1900)
	// 	ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

	// 	// -4 str per level over 60
	// 	strBonus := 100.0 - 4.0*float64(agent.GetCharacter().Level /*core.CharacterLevel*/ -60)
	// 	mhAura := character.NewTemporaryStatsAura("Crusader Enchant MH", core.ActionID{SpellID: 20007, Tag: 1}, stats.Stats{stats.Strength: strBonus}, time.Second*15)
	// 	ohAura := character.NewTemporaryStatsAura("Crusader Enchant OH", core.ActionID{SpellID: 20007, Tag: 2}, stats.Stats{stats.Strength: strBonus}, time.Second*15)

	// 	aura := character.GetOrRegisterAura(core.Aura{
	// 		Label:    "Crusader Enchant",
	// 		Duration: core.NeverExpires,
	// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
	// 			aura.Activate(sim)
	// 		},
	// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 			if !result.Landed() {
	// 				return
	// 			}

	// 			if ppmm.Proc(sim, spell.ProcMask, "Crusader") {
	// 				if spell.IsMH() {
	// 					mhAura.Activate(sim)
	// 				} else {
	// 					ohAura.Activate(sim)
	// 				}
	// 			}
	// 		},
	// 	})

	// 	character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(1900, 1.0, &ppmm, aura)
	// })

	// Weapon - Dismantle
	core.NewEnchantEffect(7210, func(agent core.Agent) {
		character := agent.GetCharacter()

		// TODO: Verify PPM
		ppm := 6.0

		procMask := character.GetProcMaskForEnchant(7210)
		ppmm := character.AutoAttacks.NewPPMManager(ppm, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 435481},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagIgnoreTargetModifiers | core.SpellFlagIgnoreAttackerModifiers,

			DamageMultiplier: 1,
			CritMultiplier:   1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				damage := sim.Roll(60, 90)
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeAlwaysHit)
			},
		})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Enchant Weapon - Dismantle",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || result.Target.MobType != proto.MobType_MobTypeMechanical {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Dismantle") {
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(7210, ppm, &ppmm, aura)
	})

	// Cloak - Subtlety
	// core.NewEnchantEffect(2621, func(agent core.Agent) {
	// 	character := agent.GetCharacter()
	// 	character.PseudoStats.ThreatMultiplier *= 0.98
	// })

	// Globes - Threat
	// core.NewEnchantEffect(2613, func(agent core.Agent) {
	// 	character := agent.GetCharacter()
	// 	character.PseudoStats.ThreatMultiplier *= 1.02
	// })

	// Ranged Scopes
	core.AddWeaponEffect(32, func(agent core.Agent, _ proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 2
		w.BaseDamageMax += 2
	})

	core.AddWeaponEffect(33, func(agent core.Agent, _ proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 3
		w.BaseDamageMax += 3
	})

	core.AddWeaponEffect(663, func(agent core.Agent, _ proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 5
		w.BaseDamageMax += 5
	})

	core.AddWeaponEffect(664, func(agent core.Agent, _ proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 7
		w.BaseDamageMax += 7
	})

	//core.AddWeaponEffect(2523, func(agent core.Agent, _ proto.ItemSlot) {
	//character := agent.GetCharacter()
	// TODO: Add ranged hit +3
	//})

	core.AddEffectsToTest = true
}
