package warrior

import (
	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) registerHamstringSpell() {
	damage := map[int32]float64{
		25: 5,
		40: 18,
		50: 18,
		60: 45,
	}[warrior.Level]

	spellID := map[int32]int32{
		25: 1715,
		40: 7372,
		50: 7372,
		60: 27584,
	}[warrior.Level]

	warrior.Hamstring = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   10 - warrior.FocusedRageDiscount,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
