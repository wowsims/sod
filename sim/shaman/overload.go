package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// This could be value or bitflag if we ended up needing multiple flags at the same time.
// 1 to 5 are used by MaelstromWeapon Stacks
const CastTagOverload = 6

const ShamanOverloadChance = .50

func (shaman *Shaman) applyOverload() {
	overloadRuneEquipped := shaman.HasRune(proto.ShamanRune_RuneChestOverload)

	if !overloadRuneEquipped {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Overload",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneChestOverload)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyOverloadModifiers(spell *core.SpellConfig) {
	spell.ActionID.Tag = int32(CastTagOverload)
	spell.ProcMask = core.ProcMaskProc
	spell.Cast.DefaultCast.CastTime = 0
	spell.Cast.DefaultCast.GCD = 0
	spell.Cast.DefaultCast.Cost = 0
	spell.Cast.ModifyCast = nil
	spell.ManaCost.BaseCost = 0
	spell.MetricSplits = 0
	spell.DamageMultiplier *= 0.5
	spell.ThreatMultiplier = 0
}
