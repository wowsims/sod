package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerSinisterStrikeSpell() {
	flatDamageBonus := map[int32]float64{
		25: 15,
		40: 33,
		50: 52,
		60: 68,
	}[rogue.Level]

	spellID := map[int32]int32{
		25: 1759,
		40: 8621,
		50: 11293,
		60: 11294,
	}[rogue.Level]

	rogue.SinisterStrike = rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueSinisterStrike,
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          rogue.builderFlags(),

		EnergyCost: core.EnergyCostOptions{
			Cost:   45,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			baseDamage := flatDamageBonus + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
