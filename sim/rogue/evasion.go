package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (rogue *Rogue) RegisterEvasionSpell() {
	//Used to double evasion due it ignoring the dynamic -50% dodge suppresion aura from JAFW
	hasJAFW := rogue.HasRune(proto.RogueRune_RuneJustAFleshWound)

	rogue.EvasionAura = rogue.NewTemporaryStatsAura(
		"Evasion",
		core.ActionID{SpellID: 5277},
		stats.Stats{stats.Dodge: core.TernaryFloat64(hasJAFW, 100*core.DodgeRatingPerDodgeChance, 50*core.DodgeRatingPerDodgeChance)},
		time.Second*15,
	)

	rogue.Evasion = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5277},
		ClassSpellMask: SpellClassMask_RogueEvasion,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: []time.Duration{time.Minute * 5, time.Minute*5 - time.Second*45, time.Minute*5 - time.Second*90}[rogue.Talents.Endurance],
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
