package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) applyRiposte() {
	if !rogue.Talents.Riposte {
		return
	}

	var riposteReady *core.Aura

	riposte := rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 14251},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       SpellFlagCarnage | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return riposteReady.IsActive()
		},

		DamageMultiplier: 1.5,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			riposteReady.Deactivate(sim)

			damage := rogue.MHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})

	riposteReady = rogue.RegisterAura(core.Aura{
		Label:    "Riposte Ready Aura",
		ActionID: riposte.ActionID,
		Duration: time.Second * 5,
	})

	rogue.RegisterAura(core.Aura{
		Label:    "Riposte Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome == core.OutcomeParry {
				riposteReady.Activate(sim)
			}
		},
	})
}
