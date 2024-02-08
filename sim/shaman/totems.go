package shaman

import (
	"github.com/wowsims/sod/sim/core"
)

func (shaman *Shaman) newTotemSpellConfig(flatCost float64, spellID int32) core.SpellConfig {
	return core.SpellConfig{
		ActionID: core.ActionID{SpellID: spellID},
		Flags:    SpellFlagTotem | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   flatCost,
			Multiplier: shaman.TotemManaMultiplier(),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
	}
}
