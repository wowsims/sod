package sod

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/shaman"
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

		procChance := 0.10
		baseDamageLow := 60.0
		baseDamageHigh := 90.0

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 439164},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			CritMultiplier:   1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(baseDamageLow, baseDamageHigh), spell.OutcomeMagicHitAndCrit)
			},
		})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Enchant Weapon - Dismantle",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// Dismantle only procs on hits that land
				if !result.Landed() {
					return
				}

				// Dismantle only procs on Mechanical units
				if result.Target.MobType != proto.MobType_MobTypeMechanical {
					return
				}

				// Dismantle only procs on attacks from the player character (not pet attacks or totems)
				if spell.Unit != &character.Unit || shaman.SpellFlagTotem.Matches(spell.Flags) {
					return
				}

				// Dismantle only procs on direct attacks, not proc effects or DoT ticks
				if core.ProcMaskProc.Matches(spell.ProcMask) || core.ProcMaskWeaponProc.Matches(spell.ProcMask) {
					return
				}

				// TODO: Confirm: Dismantle can not proc itself
				if spell == procSpell {
					return
				}

				// Main-Hand hits only trigger Dismantle if the MH weapon is enchanted with Dismantle
				if core.ProcMaskMeleeMH.Matches(spell.ProcMask) && (character.GetMHWeapon() == nil || character.GetMHWeapon().Enchant.EffectID != 7210) {
					return
				}

				// Off-Hand hits only trigger Dismantle if the MH weapon is enchanted with Dismantle
				if core.ProcMaskMeleeOH.Matches(spell.ProcMask) && (character.GetOHWeapon() == nil || character.GetOHWeapon().Enchant.EffectID != 7210) {
					return
				}

				if spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
					if sim.RandomFloat("Dismantle") < procChance {
						// Spells proc both Main-Hand and Off-Hand if both are enchanted
						if character.GetMHWeapon() != nil && character.GetMHWeapon().Enchant.EffectID == 7210 {
							procSpell.Cast(sim, result.Target)
						}
						if character.GetOHWeapon() != nil && character.GetOHWeapon().Enchant.EffectID == 7210 {
							procSpell.Cast(sim, result.Target)
						}
					}
				} else if sim.RandomFloat("Dismantle") < procChance {
					// Physical hits only proc on the hand that was hit with
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffect(7210, aura)
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
