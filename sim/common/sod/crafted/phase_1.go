package crafted

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Cloth
	///////////////////////////////////////////////////////////////////////////

	// Extraplanar Spidersilk Boots
	core.NewItemEffect(210795, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 428489}

		buffAura := character.RegisterAura(core.Aura{
			Label:    "Planar Shift",
			ActionID: actionID,
			Duration: time.Second * 6,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.DamageDealtMultiplier *= .7
				character.PseudoStats.DamageTakenMultiplier *= .7
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.DamageDealtMultiplier /= .7
				character.PseudoStats.DamageTakenMultiplier /= .7
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagDefensiveEquipment,

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
			Type:     core.CooldownTypeSurvival,
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return character.IsTanking()
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Leather
	///////////////////////////////////////////////////////////////////////////

	// Void-Touched Leather Gloves
	core.NewItemEffect(211423, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 429867}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Void Madness",
			ActionID: actionID,
			Duration: time.Second * 10,
		}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.ThreatMultiplier, 1.2)

		ee := NewSodCraftedAttackSpeedEffect(buffAura, 1.1)

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				buffAura.Activate(sim)
			},

			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return !ee.Category.AnyActive()
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Void-Touched Leather Gauntlets
	core.NewItemEffect(211502, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 429868}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Void Madness",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1.1)
				character.PseudoStats.ThreatMultiplier *= 1.2
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1 / 1.1)
				character.PseudoStats.ThreatMultiplier /= 1.2
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,

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

	// Shifting Silver Breastplate
	// core.NewItemEffect(210794, func(agent core.Agent) {
	// 	// Nothing to do
	// })

	core.AddEffectsToTest = true
}

func NewSodCraftedAttackSpeedEffect(aura *core.Aura, attackSpeedBonus float64) *core.ExclusiveEffect {
	return aura.NewExclusiveEffect("SoD Crafted Attack Speed", false, core.ExclusiveEffect{
		Priority: attackSpeedBonus,
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(sim, ee.Priority)
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(sim, 1/ee.Priority)
		},
	})
}
