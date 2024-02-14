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

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Charged Inspiration",
			ActionID: actionId,
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.SpellPower, 50)
				character.PseudoStats.CostMultiplier *= 0.5
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.SpellPower, -50)
				character.PseudoStats.CostMultiplier /= 0.5
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionId,
			Flags:    core.SpellFlagNoOnCastComplete,

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

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Hyperconductive Goldwrap
	core.NewItemEffect(215115, func(agent core.Agent) {
		character := agent.GetCharacter()
		buffAuraCrit := character.GetOrRegisterAura(core.Aura{
			Label:    "Coin Flip: Crit",
			ActionID: core.ActionID{SpellID: 437698},
			Duration: time.Second * 30,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.MeleeCrit, 3)
				character.AddStatDynamic(sim, stats.SpellCrit, 3)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.MeleeCrit, -3)
				character.AddStatDynamic(sim, stats.SpellCrit, -3)
			},
		})

		buffAuraMs := character.GetOrRegisterAura(core.Aura{
			Label:    "Coin Flip: Movement Speed",
			ActionID: core.ActionID{SpellID: 437699},
			Duration: time.Second * 30,
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 437368},
			Flags:    core.SpellFlagNoOnCastComplete,

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

		character.AddMajorCooldown(core.MajorCooldown{
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
		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Gneuro-Logical Shock",
			ActionID: core.ActionID{SpellID: 437349},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.2)
				character.MultiplyRangedSpeed(sim, 1.2)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.0/1.2)
				character.MultiplyRangedSpeed(sim, 1.0/1.2)
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 437349},
			Flags:    core.SpellFlagNoOnCastComplete,

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

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Gneuro-Conductive Channeler's Hood
	core.NewItemEffect(215381, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionId := core.ActionID{SpellID: 437357}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Gneuromantic Meditation",
			ActionID: actionId,
			Duration: time.Second * 20,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.SpellPower, 50)
				character.PseudoStats.ForceFullSpiritRegen = true
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.SpellPower, -50)
				character.PseudoStats.ForceFullSpiritRegen = false
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionId,
			Flags:    core.SpellFlagNoOnCastComplete,

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

		character.AddMajorCooldown(core.MajorCooldown{
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

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Hyperconductive Shock",
			ActionID: actionId,
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyCastSpeed(1.2)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyCastSpeed(1 / 1.2)
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionId,
			Flags:    core.SpellFlagNoOnCastComplete,

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

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Mail
	///////////////////////////////////////////////////////////////////////////

	// Tempered Interference-Negating Helmet
	core.NewItemEffect(215161, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionId := core.ActionID{SpellID: 437377}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Intense Concentration",
			ActionID: actionId,
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.2)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.0/1.2)
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionId,
			Flags:    core.SpellFlagNoOnCastComplete,

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

		character.AddMajorCooldown(core.MajorCooldown{
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
