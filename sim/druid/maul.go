package druid

import (
	"github.com/wowsims/sod/sim/core"
)

func (druid *Druid) registerMaulSpell() {
	flatBaseDamage := 128.0
	rageCost := 15 - float64(druid.Talents.Ferocity)

	switch druid.Ranged().ID {
	case IdolOfBrutality:
		rageCost -= 3
	}

	druid.Maul = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 9881},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeMHAuto,
		Flags:       SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		RageCost: core.RageCostOptions{
			Cost:   rageCost,
			Refund: 0.8,
		},

		DamageMultiplier: 1 + .1*float64(druid.Talents.SavageFury),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Need to specially deactivate CC here in case maul is cast simultaneously with another spell.
			if druid.ClearcastingAura != nil {
				druid.ClearcastingAura.Deactivate(sim)
			}

			baseDamage := flatBaseDamage + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			druid.MaulQueueAura.Deactivate(sim)
		},
	})

	druid.MaulQueueAura = druid.RegisterAura(core.Aura{
		Label:    "Maul Queue Aura",
		ActionID: druid.Maul.ActionID,
		Duration: core.NeverExpires,
	})

	druid.MaulQueueSpell = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    druid.Maul.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !druid.MaulQueueAura.IsActive() &&
				druid.CurrentRage() >= druid.Maul.Cost.GetCurrentCost() &&
				!druid.IsCasting(sim)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			druid.MaulQueueAura.Activate(sim)
		},
	})
}

func (druid *Druid) QueueMaul(sim *core.Simulation) {
	if druid.MaulQueueSpell.CanCast(sim, druid.CurrentTarget) {
		druid.MaulQueueSpell.Cast(sim, druid.CurrentTarget)
	}
}

// Returns true if the regular melee swing should be used, false otherwise.
func (druid *Druid) MaulReplaceMH(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if !druid.MaulQueueAura.IsActive() {
		return mhSwingSpell
	}

	if !druid.Maul.Spell.CanCast(sim, druid.CurrentTarget) {
		druid.MaulQueueAura.Deactivate(sim)
		return mhSwingSpell
	}

	return druid.Maul.Spell
}
