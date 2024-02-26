package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	// Exsanguinar
	core.NewItemEffect(216497, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionId := core.ActionID{SpellID: 436332}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionId,
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagAPL,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			DamageMultiplier: 1,
			CritMultiplier:   1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Exsanguination",
				},
				TickLength:    2 * time.Second,
				NumberOfTicks: 15,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = 5 * character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType]
					dot.SnapshotCritChance = 0
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				},

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				// The spell is also affected by phys school mods because it's shadow + physical school.
				damage := 65 * character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *
					target.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical]
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					// Has no DefenseType, also haven't seen a miss in logs.
					result := spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeAlwaysHit)
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
