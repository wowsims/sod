package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func init() {
	core.AddEffectsToTest = false

	// Exsanguinar
	core.NewItemEffect(216497, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionId := core.ActionID{SpellID: 436332}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionId,
			SpellSchool: core.SpellSchoolShadowstrike,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagAPL,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			DamageMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Exsanguination",
				},
				TickLength:    2 * time.Second,
				NumberOfTicks: 15,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 5, 0, isRollover)
				},

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamageNew(sim, target, 0, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					// Has no DefenseType, also haven't seen a miss in logs.
					result := spell.CalcAndDealDamageNew(sim, aoeTarget, 65, 0, spell.OutcomeAlwaysHit)
					if result.Landed() {
						spell.Dot(aoeTarget).Apply(sim)
					}
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.AddEffectsToTest = true
}
