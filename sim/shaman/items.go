package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// Totem Item IDs
const (
	TotemCarvedDriftwoodIcon = 209575
	TotemInvigoratingFlame   = 215436
)

func init() {
	core.NewItemEffect(TotemCarvedDriftwoodIcon, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.AddStat(stats.MP5, 2)
	})

	// Ancestral Bloodstorm Beacon
	core.NewItemEffect(216615, func(agent core.Agent) {
		character := agent.GetCharacter()

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 436413},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagAPL,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 5,
				},
			},

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				// The spell is also affected by shadow school mods because it's shadow + nature school.
				damage := 150 * character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *
					target.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow]
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})
}
