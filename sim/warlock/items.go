package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	// core.NewItemEffect(32493, func(agent core.Agent) {
	// 	warlock := agent.(WarlockAgent).GetWarlock()
	// 	procAura := warlock.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{SpellID: 40478}, stats.Stats{stats.SpellPower: 220}, time.Second*5)

	// 	warlock.RegisterAura(core.Aura{
	// 		Label:    "Ashtongue Talisman",
	// 		Duration: core.NeverExpires,
	// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
	// 			aura.Activate(sim)
	// 		},
	// 		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 			if spell == warlock.Corruption && sim.Proc(0.2, "Ashtongue Talisman of Insight") {
	// 				procAura.Activate(sim)
	// 			}
	// 		},
	// 	})
	// })

	// Infernal Pact Essence
	core.NewItemEffect(216509, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()

		if warlock.Pet != nil {
			warlock.Pet.AddStat(stats.Stamina, 20)
			warlock.Pet.AddStat(stats.Intellect, 80)
		}

		spell := warlock.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 436479},
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagAPL,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warlock.NewTimer(),
					Duration: time.Minute * 5,
				},
			},

			DamageMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				// The spell is also affected by fire school mods because it's shadow + fire school.
				damage := 150 * warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *
					target.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire]
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		warlock.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})
}
