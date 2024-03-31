package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (shaman *Shaman) ShockCD() time.Duration {
	return time.Second*6 - time.Millisecond*200*time.Duration(shaman.Talents.Reverberation)
}

// Shared logic for all shocks.
func (shaman *Shaman) newShockSpellConfig(actionId core.ActionID, spellSchool core.SpellSchool, baseCost float64, shockTimer *core.Timer) core.SpellConfig {
	return core.SpellConfig{
		ActionID:    actionId,
		SpellSchool: spellSchool,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | SpellFlagFocusable,

		ManaCost: core.ManaCostOptions{
			FlatCost:   baseCost,
			Multiplier: 1 - 0.02*float64(shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shockTimer,
				Duration: shaman.ShockCD(),
			},
		},

		CritDamageBonus: shaman.elementalFury(),

		DamageMultiplier: shaman.ConcussionMultiplier(),
		ThreatMultiplier: 1,
	}
}

func (shaman *Shaman) registerShocks() {
	shockTimer := shaman.NewTimer()
	shaman.registerEarthShockSpell(shockTimer)
	shaman.registerFlameShockSpell(shockTimer)
	shaman.registerFrostShockSpell(shockTimer)
}
