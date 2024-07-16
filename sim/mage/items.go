package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	FireRuby = 20036
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(FireRuby, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{ItemID: FireRuby}
		manaMetrics := character.NewManaMetrics(actionID)

		damageAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Chaos Fire",
			ActionID: core.ActionID{SpellID: 24389},
			Duration: time.Minute * 1,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.FirePower, 100)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.FirePower, -100)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.SpellSchool.Matches(core.SpellSchoolFire) {
					aura.Deactivate(sim)
				}
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				character.AddMana(sim, sim.Roll(1, 500), manaMetrics)
				damageAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	core.AddEffectsToTest = true
}
