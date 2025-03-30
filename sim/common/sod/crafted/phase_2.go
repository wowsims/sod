package crafted

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Cloth
	///////////////////////////////////////////////////////////////////////////

	// Gneuro-Linked Arcano-Filament Monocle
	core.NewItemEffect(215111, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionId := core.ActionID{SpellID: 437327}

		buffAura := character.RegisterAura(core.Aura{
			Label:    "Charged Inspiration",
			ActionID: actionId,
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SchoolCostMultiplier.AddToMagicSchools(-50)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SchoolCostMultiplier.AddToMagicSchools(50)
			},
		}).AttachStatBuff(stats.SpellPower, 50)

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionId,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Hyperconductive Goldwrap
	core.NewItemEffect(215115, func(agent core.Agent) {
		character := agent.GetCharacter()
		buffAuraCrit := character.RegisterAura(core.Aura{
			Label:    "Coin Flip: Crit",
			ActionID: core.ActionID{SpellID: 437698},
			Duration: time.Second * 30,
		}).AttachStatsBuff(stats.Stats{stats.MeleeCrit: 3, stats.SpellCrit: 3})

		buffAuraMs := character.GetOrRegisterAura(core.Aura{
			Label:    "Coin Flip: Movement Speed",
			ActionID: core.ActionID{SpellID: 437699},
			Duration: time.Second * 30,
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 437368},
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 15,
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				if sim.RandomFloat("Coin Flip") > 0.5 {
					buffAuraCrit.Activate(sim)
				} else {
					buffAuraMs.Activate(sim)
				}
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Leather
	///////////////////////////////////////////////////////////////////////////

	// Glowing Gneuro-Linked Cowl
	core.NewItemEffect(215166, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionId := core.ActionID{SpellID: 437349}
		healthMetrics := character.NewHealthMetrics(actionId)

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Gneuro-Logical Shock",
			ActionID: actionId,
			Duration: time.Second * 10,
		})

		ee := NewSodCraftedAttackSpeedEffect(buffAura, 1.2)

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionId,
			SpellSchool: core.SpellSchoolArcane,
			ProcMask:    core.ProcMaskEmpty,
			// TODO: Verify if SP affects the damage
			Flags: core.SpellFlagNoOnCastComplete | core.SpellFlagNoMetrics | core.SpellFlagIgnoreAttackerModifiers,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},

			DamageMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcDamage(sim, &character.Unit, sim.Roll(343, 757), spell.OutcomeAlwaysHit)
				if sim.Log != nil {
					character.Log(sim, "Took %.1f damage from Gneuro-Logical Shock.", result.Damage)
				}
				character.RemoveHealth(sim, result.Damage, healthMetrics)
				buffAura.Activate(sim)
			},

			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return !ee.Category.AnyActive()
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionId,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				spell.Cast(sim, &character.Unit)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Gneuro-Conductive Channeler's Hood
	core.NewItemEffect(215381, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionId := core.ActionID{SpellID: 437357}

		buffAura := character.RegisterAura(core.Aura{
			Label:    "Gneuromantic Meditation",
			ActionID: actionId,
			Duration: time.Second * 20,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.ForceFullSpiritRegen = true
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.ForceFullSpiritRegen = false
			},
		}).AttachStatBuff(stats.SpellPower, 50)

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionId,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Mail
	///////////////////////////////////////////////////////////////////////////

	// Rad-Resistant Scale Hood
	// core.NewItemEffect(215114, func(agent core.Agent) {
	// 	// Nothing to do
	// })

	// Glowing Hyperconductive Scale Coif
	core.NewItemEffect(215114, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionId := core.ActionID{SpellID: 437362}
		healthMetrics := character.NewHealthMetrics(actionId)

		buffAura := character.RegisterAura(core.Aura{
			Label:    "Hyperconductive Shock",
			ActionID: actionId,
			Duration: time.Second * 10,
		}).AttachMultiplyCastSpeed(&character.Unit, 1.2)

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionId,
			SpellSchool: core.SpellSchoolArcane,
			ProcMask:    core.ProcMaskEmpty,
			// TODO: Verify if SP affects the damage
			Flags: core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoMetrics | core.SpellFlagIgnoreAttackerModifiers,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},

			DamageMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcDamage(sim, &character.Unit, sim.Roll(312, 668), spell.OutcomeAlwaysHit)
				character.RemoveHealth(sim, result.Damage, healthMetrics)
				buffAura.Activate(sim)
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionId,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				spell.Cast(sim, &character.Unit)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Plate
	///////////////////////////////////////////////////////////////////////////

	// Tempered Interference-Negating Helmet
	core.NewItemEffect(215161, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.Level > 59 {
			return
		}

		actionId := core.ActionID{SpellID: 437377}

		buffAura := character.RegisterAura(core.Aura{
			Label:    "Intense Concentration",
			ActionID: actionId,
			Duration: time.Second * 10,
		}).AttachMultiplyAttackSpeed(&character.Unit, 1.2)

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionId,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Reflective Truesilver Braincage
	// core.NewItemEffect(215167, func(agent core.Agent) {
	// 	// Nothing to do
	// })

	core.AddEffectsToTest = true
}
