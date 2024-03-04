package rogue

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerWaylayAura() {
	if !rogue.HasRune(proto.RogueRune_RuneWaylay) {
		return
	}

	rogue.WaylayAuras = rogue.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return core.WaylayAura(target)
	})
}
