package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) registerWhirlwindSpell() {
	if warrior.Level < 36 {
		return
	}

	results := make([]*core.SpellResult, min(4, warrior.Env.GetNumTargets()))

	warrior.Whirlwind = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1680},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL | SpellFlagBloodSurge,

		RageCost: core.RageCostOptions{
			Cost: 25 - warrior.FocusedRageDiscount,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance) || warrior.StanceMatches(GladiatorStance)
		},

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1.25,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}
