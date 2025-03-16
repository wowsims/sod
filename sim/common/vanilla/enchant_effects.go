package vanilla

import (
	"math"
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
	core.AddWeaponEffect(32, func(agent core.Agent) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 2
		w.BaseDamageMax += 2
	})

	// Accurate Scope
	core.AddWeaponEffect(33, func(agent core.Agent) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 3
		w.BaseDamageMax += 3
	})

	// Weapon - Fiery Blaze
	core.NewEnchantEffect(36, func(agent core.Agent) {
		character := agent.GetCharacter()

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

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Fiery Blaze",
			Callback:          core.CallbackOnSpellHitDealt,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			Outcome:           core.OutcomeLanded,
			DPM:               character.AutoAttacks.NewDynamicProcManagerForEnchant(36, 0, 0.15),
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.ProcMask = core.Ternary(spell.ProcMask.Matches(core.ProcMaskMeleeSpecial), procMaskOnSpecial, procMaskOnAuto)
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterEnchantProc(36, triggerAura)
	})

	// Weapon - Lesser Striking
	core.AddWeaponEffect(241, func(agent core.Agent) {
		registerStrikingEnchantEffect(agent.GetCharacter(), "Lesser Striking", 241, 2)
	})

	// Weapon - Beast Slaying
	core.AddWeaponEffect(249, func(agent core.Agent) {
		registerSlayerEnchantEffect(agent.GetCharacter(), "Beast Slaying", 249, proto.MobType_MobTypeBeast, 2)
	})

	// Weapon - Minor Striking
	core.AddWeaponEffect(250, func(agent core.Agent) {
		registerStrikingEnchantEffect(agent.GetCharacter(), "Minor Striking", 250, 1)
	})

	// Deadly Scope
	core.AddWeaponEffect(663, func(agent core.Agent) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 5
		w.BaseDamageMax += 5
	})

	// Sniper Scope
	core.AddWeaponEffect(664, func(agent core.Agent) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 7
		w.BaseDamageMax += 7
	})

	// Weapon - Fiery Weapon
	core.AddWeaponEffect(803, func(agent core.Agent) {
		character := agent.GetCharacter()

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

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Fiery Weapon",
			Callback:          core.CallbackOnSpellHitDealt,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			Outcome:           core.OutcomeLanded,
			DPM:               character.AutoAttacks.NewDynamicProcManagerForEnchant(803, 6.0, 0),
			DPMProcCheck:      core.DPMProc,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.ProcMask = core.Ternary(spell.ProcMask.Matches(core.ProcMaskMeleeSpecial), procMaskOnSpecial, procMaskOnAuto)
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterEnchantProc(803, triggerAura)
	})

	// Weapon - Greater Striking
	core.AddWeaponEffect(805, func(agent core.Agent) {
		registerStrikingEnchantEffect(agent.GetCharacter(), "Greater Striking", 805, 4)
	})

	// Weapon - Lesser Beastslayer
	core.AddWeaponEffect(853, func(agent core.Agent) {
		registerSlayerEnchantEffect(agent.GetCharacter(), "Lesser Beast Slaying", 853, proto.MobType_MobTypeBeast, 6)
	})

	// Weapon - Lesser Elemental Slayer
	core.AddWeaponEffect(854, func(agent core.Agent) {
		registerSlayerEnchantEffect(agent.GetCharacter(), "Lesser Elemental Slaying", 854, proto.MobType_MobTypeElemental, 6)
	})

	// Boots - Minor Speed
	core.NewEnchantEffect(911, func(agent core.Agent) {
		character := agent.GetCharacter()

		aura := character.RegisterAura(core.Aura{
			Label: "Minor Speed",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				character.AddMoveSpeedModifier(&core.ActionID{SpellID: 13889}, 1.08)
			},
		})

		character.ItemSwap.RegisterEnchantProc(911, aura)
	})

	// Weapon - Iron Counterweight
	// Effect #34 explicitly does NOT affect ranged attack speed
	core.NewEnchantEffect(34, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.NewDynamicEquipEffectAura(core.DynamicEquipEffectConfig{
			EffectID: 34,
			Label:    "Weapon Counterweight",
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				character.PseudoStats.MeleeSpeedMultiplier *= 1 / math.Pow(1.03, float64(oldStacks))
				character.PseudoStats.RangedSpeedMultiplier *= 1 / math.Pow(1.03, float64(oldStacks))

				character.PseudoStats.MeleeSpeedMultiplier *= math.Pow(1.03, float64(newStacks))
				character.PseudoStats.RangedSpeedMultiplier *= math.Pow(1.03, float64(newStacks))
			},
		})
	})

	// Weapon - Striking
	core.AddWeaponEffect(943, func(agent core.Agent) {
		registerStrikingEnchantEffect(agent.GetCharacter(), "Striking", 943, 3)
	})

	// Weapon - Superior Striking
	core.AddWeaponEffect(1897, func(agent core.Agent) {
		registerStrikingEnchantEffect(agent.GetCharacter(), "Superior Striking", 1897, 5)
	})

	// Weapon - Lifestealing
	core.NewEnchantEffect(1898, func(agent core.Agent) {
		character := agent.GetCharacter()

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

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Lifestealing Weapon",
			Callback:          core.CallbackOnSpellHitDealt,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			Outcome:           core.OutcomeLanded,
			DPM:               character.AutoAttacks.NewDynamicProcManagerForEnchant(1898, 6.66, 0),
			DPMProcCheck:      core.DPMProc,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.ProcMask = core.Ternary(spell.ProcMask.Matches(core.ProcMaskMeleeSpecial), procMaskOnSpecial, procMaskOnAuto)
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterEnchantProc(1898, triggerAura)
	})

	// TODO: Crusader, Mongoose, and Executioner could also be modelled as AddWeaponEffect instead
	// ApplyCrusaderEffect will be applied twice if there is two weapons with this enchant.
	//   However, it will automatically overwrite one of them, so it should be ok.
	//   A single application of the aura will handle both mh and oh procs.
	core.NewEnchantEffect(1900, func(agent core.Agent) {
		character := agent.GetCharacter()

		dpm := character.AutoAttacks.NewDynamicProcManagerForEnchant(1900, 1.0, 0)

		// -4 str per level over 60
		strBonus := 100.0 - 4.0*float64(character.Level-60)
		mhAura := character.NewTemporaryStatsAura("Crusader Enchant MH", core.ActionID{SpellID: 20007, Tag: 1}, stats.Stats{stats.Strength: strBonus}, time.Second*15)
		ohAura := character.NewTemporaryStatsAura("Crusader Enchant OH", core.ActionID{SpellID: 20007, Tag: 2}, stats.Stats{stats.Strength: strBonus}, time.Second*15)
		healthMetrics := character.NewHealthMetrics(core.ActionID{SpellID: 20007})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Crusader",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			DPM:               dpm,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.IsMH() {
					mhAura.Activate(sim)
				} else {
					ohAura.Activate(sim)
				}
				character.GainHealth(sim, sim.RollWithLabel(75, 125, "Crusader Heal"), healthMetrics)
			},
		})

		character.ItemSwap.RegisterEnchantProc(1900, aura)
	})

	// Biznicks 247x128 Accurascope
	core.AddWeaponEffect(2523, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.AddBonusRangedHitRating(3)
	})

	// Gloves - Threat
	core.NewEnchantEffect(2613, func(agent core.Agent) {
		character := agent.GetCharacter()

		core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Threat +2%",
		}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.ThreatMultiplier, 1.02))
	})

	// Cloak - Subtlety
	core.NewEnchantEffect(2621, func(agent core.Agent) {
		character := agent.GetCharacter()

		core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Subtlety",
		}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.ThreatMultiplier, 1/1.02))
	})

	core.AddEffectsToTest = true
}

func registerStrikingEnchantEffect(character *core.Character, label string, effectID int32, bonus float64) {
	applyWeaponBonus := func(weapon *core.Weapon, bonus float64) {
		weapon.BaseDamageMin += bonus
		weapon.BaseDamageMax += bonus
	}

	makeItemSlotAura := func(itemSlotProcMask core.ProcMask) *core.Aura {
		var weapon *core.Weapon
		switch {
		case itemSlotProcMask.Matches(core.ProcMaskMeleeMH):
			weapon = character.AutoAttacks.MH()
			label += " MH"
		case itemSlotProcMask.Matches(core.ProcMaskMeleeOH):
			weapon = character.AutoAttacks.OH()
			label += " OH"
		}

		return core.MakePermanent(character.RegisterAura(core.Aura{
			Label: label,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				applyWeaponBonus(weapon, bonus)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				applyWeaponBonus(weapon, -bonus)
			},
		}))
	}

	mhAura := makeItemSlotAura(core.ProcMaskMeleeMH)
	ohAura := makeItemSlotAura(core.ProcMaskMeleeOH)

	character.ItemSwap.RegisterEnchantProcWithSlots(effectID, mhAura, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand})
	character.ItemSwap.RegisterEnchantProcWithSlots(effectID, ohAura, []proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand})
}

func registerSlayerEnchantEffect(character *core.Character, label string, effectID int32, mobType proto.MobType, bonus float64) {
	applyWeaponBonus := func(weapon *core.Weapon, bonus float64) {
		weapon.BaseDamageMin += bonus
		weapon.BaseDamageMax += bonus
	}

	makeItemSlotAura := func(itemSlotProcMask core.ProcMask) *core.Aura {
		var weapon *core.Weapon
		rangedWeapon := character.AutoAttacks.Ranged()
		switch {
		case itemSlotProcMask.Matches(core.ProcMaskMeleeMH):
			weapon = character.AutoAttacks.MH()
			label += " MH"
		case itemSlotProcMask.Matches(core.ProcMaskMeleeOH):
			weapon = character.AutoAttacks.OH()
			label += " OH"
		}

		return character.RegisterAura(core.Aura{
			Label: label,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				if character.CurrentTarget.MobType == mobType {
					aura.Activate(sim)
				}
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if character.CurrentTarget.MobType != mobType {
					aura.Deactivate(sim)
				}
				applyWeaponBonus(weapon, bonus)
				applyWeaponBonus(rangedWeapon, bonus)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				applyWeaponBonus(weapon, -bonus)
				applyWeaponBonus(rangedWeapon, -bonus)
			},
		})
	}

	mhAura := makeItemSlotAura(core.ProcMaskMeleeMH)
	ohAura := makeItemSlotAura(core.ProcMaskMeleeOH)

	character.ItemSwap.RegisterEnchantProcWithSlots(effectID, mhAura, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand})
	character.ItemSwap.RegisterEnchantProcWithSlots(effectID, ohAura, []proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand})
}
