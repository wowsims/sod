package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerMasterOfSubtlety() {
	if !rogue.HasRune(proto.RogueRune_RuneMasterOfSubtlety) {
		return
	}

	percent := 1.1

	effectDuration := time.Second * 6
	if rogue.StealthAura.IsActive() {
		effectDuration = core.NeverExpires
	}

	rogue.MasterOfSubtletyAura = rogue.RegisterAura(core.Aura{
		Label:    "Master of Subtlety",
		ActionID: core.ActionID{SpellID: int32(proto.RogueRune_RuneMasterOfSubtlety)},
		Duration: effectDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= percent
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= 1 / percent
		},
	})
}
