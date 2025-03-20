package shaman

import (
	"github.com/wowsims/sod/sim/core"
)

func (shaman *Shaman) newTotemSpellConfig(classMask uint64, spellID int32, flatCost float64) core.SpellConfig {
	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		ClassSpellMask: classMask,
		Flags:          core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   flatCost,
			Multiplier: shaman.totemManaMultiplier(),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
	}
}
