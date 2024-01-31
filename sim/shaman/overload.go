package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

type CastTagOverload int32

const (
	// This could be value or bitflag if we ended up needing multiple flags at the same time.
	//1 to 5 are used by MaelstromWeapon Stacks
	CastTagLightningBoltOverload CastTagOverload = iota + 6
	CastTagChainLightningOverload
	CastTagLavaBurstOverload

	CastTagHealingWaveOverload
	CastTagChainHealOverload
)

func (shaman *Shaman) applyOverload() {
	overloadRuneEquipped := shaman.HasRune(proto.ShamanRune_RuneChestOverload)
	shaman.OverloadChance = core.TernaryFloat64(overloadRuneEquipped, .50, 0)

	if !overloadRuneEquipped {
		return
	}

	shaman.OverloadAura = shaman.RegisterAura(core.Aura{
		Label:    "Overload",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneChestOverload)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyOverloadModifiers(spell *core.SpellConfig, tag CastTagOverload) {
	spell.ActionID.Tag = int32(tag)
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
