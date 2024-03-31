package sod

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func init() {
	core.AddEffectsToTest = false

	core.AddEffectsToTest = true
}

// For Automatic Crowd Pummeler and Druid's Catnip
func RegisterFiftyPercentHasteBuffCD(character *core.Character, actionID core.ActionID) {
	aura := character.GetOrRegisterAura(core.Aura{
		Label:    "Haste",
		ActionID: core.ActionID{SpellID: 13494},
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			character.MultiplyAttackSpeed(sim, 1.5)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			character.MultiplyAttackSpeed(sim, 1.0/1.5)
		},
	})

	spell := character.GetOrRegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    character.GetFiftyPercentHasteBuffCD(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aura.Activate(sim)
		},
	})

	character.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}
