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

		character.ItemSwap.RegisterOnSwapItemForEffect(7210, aura)
	})

	core.AddEffectsToTest = true
}
