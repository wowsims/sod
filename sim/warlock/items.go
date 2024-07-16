package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	InfernalPactEssence = 216509
	ZIlaGular           = 223214
)

func init() {
	// Infernal Pact Essence
	core.NewItemEffect(InfernalPactEssence, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()

		stats := stats.Stats{
			stats.Stamina:   20,
			stats.Intellect: 80,
		}

		// TODO: Does this affect Infernal or Doomguard?
		warlock.Felhunter.AddStats(stats)
		warlock.Imp.AddStats(stats)
		warlock.Succubus.AddStats(stats)
		warlock.Voidwalker.AddStats(stats)
		if warlock.Felguard != nil {
			warlock.Felguard.AddStats(stats)
		}

		spell := warlock.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 436479},
			SpellSchool: core.SpellSchoolShadow | core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warlock.NewTimer(),
					Duration: time.Minute * 5,
				},
			},

			DamageMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, 150, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		warlock.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Zila Gular
	core.NewItemEffect(ZIlaGular, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()

		warlock.zilaGularAura = warlock.GetOrRegisterAura(core.Aura{
			Label:    "Zila Gular",
			ActionID: core.ActionID{SpellID: 448686},
			Duration: time.Second * 20,
		})

		spell := warlock.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 448686},
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warlock.NewTimer(),
					Duration: time.Minute * 5,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				warlock.zilaGularAura.Activate(sim)
			},
		})

		warlock.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})
}
