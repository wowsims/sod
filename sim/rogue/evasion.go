package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (rogue *Rogue) RegisterEvasionSpell() {
	rogue.EvasionAura = rogue.RegisterAura(core.Aura{
		Label:    "Evasion",
		ActionID: core.ActionID{SpellID: 5277},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.AddStatDynamic(sim, stats.Dodge, 15*core.DodgeRatingPerDodgeChance)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.AddStatDynamic(sim, stats.Dodge, -15*core.DodgeRatingPerDodgeChance)
		},
	})

	rogue.Evasion = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 5277},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: []time.Duration{time.Minute * 5, time.Minute*5 - time.Second*45, time.Second*5 - time.Second*90}[rogue.Talents.Elusiveness],
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Activate aura
			rogue.EvasionAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.Evasion,
		Type:  core.CooldownTypeSurvival,
	})
}
