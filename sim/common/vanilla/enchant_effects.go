package vanilla

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                        All effects ordered by ID
	///////////////////////////////////////////////////////////////////////////

	// Ranged Scopes
	core.AddWeaponEffect(32, func(agent core.Agent, _ proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 2
		w.BaseDamageMax += 2
	})

	// Accurate Scope
	core.AddWeaponEffect(33, func(agent core.Agent, _ proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 3
		w.BaseDamageMax += 3
	})

	// Weapon - Fiery Blaze
	core.NewEnchantEffect(36, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(36)
		procChance := 0.15

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 6296},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					damage := sim.Roll(9, 13)
					spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
				}

			},
		})

		procMaskOnAuto := core.ProcMaskDamageProc     // Both spell and melee proc combo
		procMaskOnSpecial := core.ProcMaskSpellDamage // TODO: check if core.ProcMaskSpellDamage remains on special

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Fiery Blaze",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(procMask) || spell.Flags.Matches(core.SpellFlagSuppressWeaponProcs) {
					return
				}

				if sim.RandomFloat("Fiery Blaze") < procChance {
					if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
						procSpell.ProcMask = procMaskOnSpecial
					} else {
						procSpell.ProcMask = procMaskOnAuto
					}
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffect(36, aura)
	})

	// Weapon - Lesser Striking
	core.AddWeaponEffect(241, func(agent core.Agent, slot proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.MH()
		if slot == proto.ItemSlot_ItemSlotOffHand {
			w = agent.GetCharacter().AutoAttacks.OH()
		}
		w.BaseDamageMin += 2
		w.BaseDamageMax += 2
	})

	// Weapon - Beast Slaying
	core.AddWeaponEffect(249, func(agent core.Agent, slot proto.ItemSlot) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			w := character.AutoAttacks.MH()
			if slot == proto.ItemSlot_ItemSlotOffHand {
				w = character.AutoAttacks.OH()
			}

			w.BaseDamageMin += 2
			w.BaseDamageMax += 2

			w = character.AutoAttacks.Ranged()
			w.BaseDamageMin += 2
			w.BaseDamageMax += 2
		}
	})

	// Weapon - Minor Striking
	core.AddWeaponEffect(250, func(agent core.Agent, slot proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.MH()
		if slot == proto.ItemSlot_ItemSlotOffHand {
			w = agent.GetCharacter().AutoAttacks.OH()
		}
		w.BaseDamageMin += 1
		w.BaseDamageMax += 1
	})

	// Deadly Scope
	core.AddWeaponEffect(663, func(agent core.Agent, _ proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 5
		w.BaseDamageMax += 5
	})

	// Sniper Scope
	core.AddWeaponEffect(664, func(agent core.Agent, _ proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 7
		w.BaseDamageMax += 7
	})

	// Weapon - Fiery Weapon
	core.AddWeaponEffect(803, func(agent core.Agent, _ proto.ItemSlot) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(803)
		ppmm := character.AutoAttacks.NewPPMManager(6.0, procMask)

		procMaskOnAuto := core.ProcMaskDamageProc     // Both spell and melee proc combo
		procMaskOnSpecial := core.ProcMaskSpellDamage // TODO: check if core.ProcMaskSpellDamage remains on special

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 13897},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    procMaskOnAuto,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 40, spell.OutcomeMagicHitAndCrit)
			},
		})

		aura := core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Fiery Weapon",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || spell.Flags.Matches(core.SpellFlagSuppressWeaponProcs) {
					return
				}
				if ppmm.Proc(sim, spell.ProcMask, "Fiery Weapon") {
					if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
						procSpell.ProcMask = procMaskOnSpecial
					} else {
						procSpell.ProcMask = procMaskOnAuto
					}
					procSpell.Cast(sim, result.Target)
				}
			},
		}))

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(803, 6.0, &ppmm, aura)
	})

	// Weapon - Greater Striking
	core.AddWeaponEffect(805, func(agent core.Agent, slot proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.MH()
		if slot == proto.ItemSlot_ItemSlotOffHand {
			w = agent.GetCharacter().AutoAttacks.OH()
		}
		w.BaseDamageMin += 4
		w.BaseDamageMax += 4
	})

	// Weapon - Lesser Beastslayer
	core.AddWeaponEffect(853, func(agent core.Agent, slot proto.ItemSlot) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			w := character.AutoAttacks.MH()
			if slot == proto.ItemSlot_ItemSlotOffHand {
				w = character.AutoAttacks.OH()
			}

			w.BaseDamageMin += 6
			w.BaseDamageMax += 6

			w = character.AutoAttacks.Ranged()
			w.BaseDamageMin += 6
			w.BaseDamageMax += 6
		}
	})

	// Weapon - Lesser Elemental Slayer
	core.AddWeaponEffect(854, func(agent core.Agent, slot proto.ItemSlot) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeElemental {
			w := character.AutoAttacks.MH()
			if slot == proto.ItemSlot_ItemSlotOffHand {
				w = character.AutoAttacks.OH()
			}

			w.BaseDamageMin += 6
			w.BaseDamageMax += 6

			w = character.AutoAttacks.Ranged()
			w.BaseDamageMin += 6
			w.BaseDamageMax += 6
		}
	})

	// Boots - Minor Speed
	core.NewEnchantEffect(911, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.RegisterAura(core.Aura{
			Label: "Minor Speed",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				character.AddMoveSpeedModifier(&core.ActionID{SpellID: 13889}, 1.08)
			},
		})
	})

	// Gloves - Minor Haste
	// Effect #931 explicitly does NOT affect ranged attack speed
	core.NewEnchantEffect(931, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.PseudoStats.MeleeSpeedMultiplier *= 1.01
		character.PseudoStats.RangedSpeedMultiplier *= 1.01
	})

	// Weapon - Striking
	core.AddWeaponEffect(943, func(agent core.Agent, slot proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.MH()
		if slot == proto.ItemSlot_ItemSlotOffHand {
			w = agent.GetCharacter().AutoAttacks.OH()
		}
		w.BaseDamageMin += 3
		w.BaseDamageMax += 3
	})

	// Weapon - Superior Striking
	core.AddWeaponEffect(1897, func(agent core.Agent, slot proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.MH()
		if slot == proto.ItemSlot_ItemSlotOffHand {
			w = agent.GetCharacter().AutoAttacks.OH()
		}
		w.BaseDamageMin += 5
		w.BaseDamageMax += 5
	})

	// Weapon - Lifestealing
	core.AddWeaponEffect(1898, func(agent core.Agent, slot proto.ItemSlot) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(1898)
		ppmm := character.AutoAttacks.NewPPMManager(6.66, procMask)

		procMaskOnAuto := core.ProcMaskDamageProc     // Both spell and melee proc combo
		procMaskOnSpecial := core.ProcMaskSpellDamage // TODO: check if core.ProcMaskSpellDamage remains on special

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 20004},
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    procMaskOnAuto,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 30, spell.OutcomeMagicHitAndCrit)
			},
		})

		aura := core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Lifestealing Weapon",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && !spell.Flags.Matches(core.SpellFlagSuppressWeaponProcs) && ppmm.Proc(sim, spell.ProcMask, "Lifestealing Weapon") {
					if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
						procSpell.ProcMask = procMaskOnSpecial
					} else {
						procSpell.ProcMask = procMaskOnAuto
					}
					procSpell.Cast(sim, result.Target)
				}
			},
		}))

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(1898, 6.66, &ppmm, aura)
	})

	// TODO: Crusader, Mongoose, and Executioner could also be modelled as AddWeaponEffect instead
	// ApplyCrusaderEffect will be applied twice if there is two weapons with this enchant.
	//   However, it will automatically overwrite one of them, so it should be ok.
	//   A single application of the aura will handle both mh and oh procs.
	core.NewEnchantEffect(1900, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(1900)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		// -4 str per level over 60
		strBonus := 100.0 - 4.0*float64(character.Level-60)
		mhAura := character.NewTemporaryStatsAura("Crusader Enchant MH", core.ActionID{SpellID: 20007, Tag: 1}, stats.Stats{stats.Strength: strBonus}, time.Second*15)
		ohAura := character.NewTemporaryStatsAura("Crusader Enchant OH", core.ActionID{SpellID: 20007, Tag: 2}, stats.Stats{stats.Strength: strBonus}, time.Second*15)

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Crusader Enchant",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || spell.Flags.Matches(core.SpellFlagSuppressWeaponProcs) {
					return
				}
				if ppmm.Proc(sim, spell.ProcMask, "Crusader") {
					if spell.IsMH() {
						mhAura.Activate(sim)
					} else {
						ohAura.Activate(sim)
					}
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(1900, 1.0, &ppmm, aura)
	})

	// Biznicks 247x128 Accurascope
	core.AddWeaponEffect(2523, func(agent core.Agent, _ proto.ItemSlot) {
		character := agent.GetCharacter()
		character.AddBonusRangedHitRating(3)
	})

	// Gloves - Libram of Rapidity
	// Confirmed to mod both melee and ranged speed
	core.NewEnchantEffect(2543, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.PseudoStats.MeleeSpeedMultiplier *= 1.01
		character.PseudoStats.RangedSpeedMultiplier *= 1.01
	})

	// Gloves - Threat
	core.NewEnchantEffect(2613, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.RegisterAura(core.Aura{
			Label: "Threat +2%",
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.ThreatMultiplier *= 1.02
			},
		})
	})

	// Cloak - Subtlety
	core.NewEnchantEffect(2621, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.RegisterAura(core.Aura{
			Label: "Subtlety",
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.ThreatMultiplier /= 1.02
			},
		})
	})

	core.AddEffectsToTest = true
}
