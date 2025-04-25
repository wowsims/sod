package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (druid *Druid) registerTigersFurySpell() {
	// Spell may have been added by King of the Jungle rune already.
	if druid.TigersFury != nil {
		return
	}

	actionID := core.ActionID{SpellID: map[int32]int32{
		25: 5217,
		40: 6793,
		50: 9845,
		60: 9846,
	}[druid.Level]}

	dmgBonus := map[int32]float64{
		25: 10.0,
		40: 20.0,
		50: 30.0,
		60: 40.0,
	}[druid.Level]

	druid.TigersFuryAura = druid.RegisterAura(core.Aura{
		Label:    "Tiger's Fury Aura",
		ActionID: actionID,
		Duration: 6 * time.Second,
	}).AttachAdditivePseudoStatBuff(&druid.PseudoStats.BonusPhysicalDamage, dmgBonus)

	spell := druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_DruidTigersFury,
		Flags:          core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			// This condition is here because of T3.5 6pc.
			// Using TF doesn't reset the duration back to 6s if it currently lasts longer.
			if !druid.TigersFuryAura.IsActive() || druid.TigersFuryAura.RemainingDuration(sim) < druid.TigersFuryAura.Duration {
				druid.TigersFuryAura.Activate(sim)
			}
		},
	})

	druid.TigersFury = spell
}

// For King of the Jungle rune.
func (druid *Druid) registerTigersFurySpellKotJ() {
	actionID := core.ActionID{SpellID: 417045}
	energyMetrics := druid.NewEnergyMetrics(actionID)

	druid.TigersFuryAura = druid.RegisterAura(core.Aura{
		Label:    "Tiger's Fury Aura",
		ActionID: actionID,
		Duration: 6 * time.Second,
	}).AttachMultiplicativePseudoStatBuff(&druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical], 1.15)

	spell := druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_DruidTigersFury,
		Flags:          core.SpellFlagAPL,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.AddEnergy(sim, 60.0, energyMetrics)
			// This condition is here because of T3.5 6pc.
			// Using TF doesn't reset the duration back to 6s if it currently lasts longer.
			if !druid.TigersFuryAura.IsActive() || druid.TigersFuryAura.RemainingDuration(sim) < druid.TigersFuryAura.Duration {
				druid.TigersFuryAura.Activate(sim)
			}
		},
	})

	druid.TigersFury = spell
}
