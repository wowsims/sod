package sod

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const (
	WolfsheadTrophy = 7124
)

func init() {
	core.AddEffectsToTest = false

	// Weapon - Dismantle
	core.NewEnchantEffect(7210, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 439164},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(60, 90), spell.OutcomeMagicHitAndCrit)
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

				// Dismantle only procs on direct attacks, not proc effects or DoT ticks
				if !spell.Flags.Matches(core.SpellFlagNotAProc) && spell.ProcMask.Matches(core.ProcMaskProc|core.ProcMaskSpellDamageProc) {
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
					if sim.RandomFloat("Dismantle") < 0.10 {
						// Spells proc both Main-Hand and Off-Hand if both are enchanted
						if character.GetMHWeapon() != nil && character.GetMHWeapon().Enchant.EffectID == 7210 {
							procSpell.Cast(sim, result.Target)
						}
						if character.GetOHWeapon() != nil && character.GetOHWeapon().Enchant.EffectID == 7210 {
							procSpell.Cast(sim, result.Target)
						}
					}
				} else if sim.RandomFloat("Dismantle") < 0.10 {
					// Physical hits only proc on the hand that was hit with
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(7210, aura)
	})

	// Sharpened Chitin Armor Kit
	// Permanently cause an item worn on the chest, legs, hands or feet to cause 20 Nature damage to the attacker when struck in combat.
	// Only usable on items level 45 and above.
	core.NewEnchantEffect(7649, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 233803}

		damage := 20.0
		numEnchants := 0
		for _, item := range character.Equipment {
			if item.Enchant.EffectID == 7649 {
				numEnchants++
			}
		}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHit)
			},
		})

		procAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Thorns +20",
			Callback: core.CallbackOnSpellHitTaken,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				for i := 0; i < numEnchants; i++ {
					procSpell.Cast(sim, spell.Unit)
				}
			},
		}).AttachAdditivePseudoStatBuff(&character.PseudoStats.ThornsDamage, damage*float64(numEnchants))

		character.ItemSwap.RegisterEnchantProc(7649, procAura)
	})

	// Obsidian Scope
	core.AddWeaponEffect(7657, func(agent core.Agent) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 10
		w.BaseDamageMax += 10
	})

	core.AddEffectsToTest = true
}
