package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// This could be value or bitflag if we ended up needing multiple flags at the same time.
// 1 to 10 are used by MaelstromWeapon Stacks
const CastTagOverload = 11

func (shaman *Shaman) applyOverload() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestOverload) {
		return
	}

	shaman.overloadProcChance += 0.60
}

func (shaman *Shaman) procOverload(sim *core.Simulation, label string, multiplier float64) bool {
	if shaman.overloadProcChance == 0 {
		return false
	}

	return sim.Proc(shaman.overloadProcChance*multiplier, label)
}

func (shaman *Shaman) applyOverloadModifiers(spell *core.SpellConfig) {
	spell.ActionID.Tag = int32(CastTagOverload)
	spell.ProcMask |= core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc
	spell.Flags |= core.SpellFlagPassiveSpell
	spell.Cast.CD = core.Cooldown{}
	spell.Cast.DefaultCast.CastTime = 0
	spell.Cast.DefaultCast.GCD = 0
	spell.Cast.DefaultCast.Cost = 0
	spell.Cast.ModifyCast = nil
	spell.ManaCost.BaseCost = 0
	spell.ManaCost.FlatCost = 0
	spell.MetricSplits = 0
	spell.DamageMultiplier *= 0.5
	spell.ThreatMultiplier = 0
}
