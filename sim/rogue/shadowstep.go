package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerShadowstep() {
	if !rogue.HasRune(proto.RogueRune_RuneShadowstep) {
		return
	}

	actionID := core.ActionID{SpellID: 400029}
	baseCost := 0.0

	rogue.Shadowstep = rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		EnergyCost: core.EnergyCostOptions{
			Cost:   baseCost,
			Refund: 0,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second * 1,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {},
	})
}
