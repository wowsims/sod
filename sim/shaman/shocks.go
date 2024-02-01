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
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagShock | core.SpellFlagAPL,

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

		DamageMultiplier: 1 + 0.01*float64(shaman.Talents.Concussion),
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
	}
}

// func (shaman *Shaman) registerEarthShockSpell(shockTimer *core.Timer) {
// 	config := shaman.newShockSpellConfig(49231, core.SpellSchoolNature, 0.18, shockTimer)
// 	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 		baseDamage := sim.Roll(854, 900) + 0.386*spell.SpellPower()
// 		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
// 	}

// 	shaman.EarthShock = shaman.RegisterSpell(config)
// }

// func (shaman *Shaman) registerFrostShockSpell(shockTimer *core.Timer) {
// 	config := shaman.newShockSpellConfig(49236, core.SpellSchoolFrost, 0.18, shockTimer)
// 	config.Cast.CD.Duration -= time.Duration(shaman.Talents.BoomingEchoes) * time.Second
// 	config.DamageMultiplier += 0.1 * float64(shaman.Talents.BoomingEchoes)
// 	config.ThreatMultiplier *= 2
// 	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 		baseDamage := sim.Roll(812, 858) + 0.386*spell.SpellPower()
// 		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
// 	}

// 	shaman.FrostShock = shaman.RegisterSpell(config)
// }

func (shaman *Shaman) registerShocks() {
	shockTimer := shaman.NewTimer()
	// shaman.registerEarthShockSpell(shockTimer)
	shaman.registerFlameShockSpell(shockTimer)
	// shaman.registerFrostShockSpell(shockTimer)
}
