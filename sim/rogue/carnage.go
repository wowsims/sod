package rogue

import (
	"math"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func carnageMultiplier(spell *core.Spell, _ *core.AttackTable) float64 {
	return core.TernaryFloat64(spell.Flags.Matches(SpellFlagCarnage), 1.15, 1)
}

func (rogue *Rogue) applyCarnage() {
	if !rogue.HasRune(proto.RogueRune_RuneCarnage) {
		return
	}

	var carnageAuras core.AuraArray

	carnageAura := core.Aura{
		Label:     "Carnage",
		ActionID:  core.ActionID{SpellID: int32(proto.RogueRune_RuneCarnage)},
		Duration:  core.NeverExpires,
		MaxStacks: math.MaxInt32,

		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, bleedAura := range aura.Unit.GetAurasWithTag(RogueBleedTag) {
				bleedAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
					ca := carnageAuras[bleedAura.Unit.UnitIndex]
					ca.Activate(sim)
					ca.AddStack(sim)
				}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
					ca := carnageAuras[bleedAura.Unit.UnitIndex]
					if ca.IsActive() { // carnage aura might already be expired by doneIteration
						ca.RemoveStack(sim)
					}
				})
			}
		},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, at := range rogue.AttackTables[aura.Unit.UnitIndex] {
				at.DamageDoneByCasterMultiplier = carnageMultiplier
			}
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, at := range rogue.AttackTables[aura.Unit.UnitIndex] {
				at.DamageDoneByCasterMultiplier = nil
			}
		},
	}

	carnageAuras = rogue.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return target.RegisterAura(carnageAura)
	})
}
