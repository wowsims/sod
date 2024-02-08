package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) registerCarveSpell() {
	if !hunter.HasRune(proto.HunterRune_RuneHandsCarve) {
		return
	}

	actionID := core.ActionID{SpellID: 425711}
	numHits := hunter.Env.GetNumTargets()
	results := make([]*core.SpellResult, numHits)

	hasDwRune := hunter.HasRune(proto.HunterRune_RuneBootsDualWieldSpecialization)

	if hunter.AutoAttacks.IsDualWielding {
		hunter.CarveOh = hunter.RegisterSpell(core.SpellConfig{
			ActionID:    actionID.WithTag(1),
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskMeleeOHSpecial,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: core.TernaryFloat64(hasDwRune, 1.5, 1),
			CritMultiplier:   hunter.critMultiplier(true, hunter.CurrentTarget),
		})
	}

	hunter.CarveMh = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 425711},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIncludeTargetBonusDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.04,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= 5
		},

		DamageMultiplier: 1,
		CritMultiplier:   hunter.critMultiplier(true, hunter.CurrentTarget),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := 0.5*spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()

				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			if hunter.CarveOh != nil {
				hunter.CarveOh.Cast(sim, target)

				curTarget = target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					baseDamage := 0.5*spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
						spell.BonusWeaponDamage()*0.5

					results[hitIndex] = hunter.CarveOh.CalcDamage(sim, curTarget, baseDamage, hunter.CarveOh.OutcomeMeleeWeaponSpecialHitAndCrit)

					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}

				curTarget = target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					hunter.CarveOh.DealDamage(sim, results[hitIndex])
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			}
		},
	})
}
