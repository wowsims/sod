package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

// Shared logic for all shocks.
func (shaman *Shaman) newShockSpellConfig(actionId core.ActionID, spellSchool core.SpellSchool, baseCost float64, shockTimer *core.Timer) core.SpellConfig {
	cdDuration := time.Second*6 - time.Millisecond*200*time.Duration(shaman.Talents.Reverberation)

	return core.SpellConfig{
		ActionID:    actionId,
		SpellSchool: spellSchool,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagShaman | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   baseCost,
			Multiplier: 100 - 2*shaman.Talents.Convection,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cdDuration,
			},
			SharedCD: core.Cooldown{
				Timer:    shockTimer,
				Duration: cdDuration,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
	}
}

func (shaman *Shaman) registerShocks() {
	shockTimer := shaman.NewTimer()
	shaman.registerEarthShockSpell(shockTimer)
	shaman.registerFlameShockSpell(shockTimer)
	shaman.registerFrostShockSpell(shockTimer)
}
