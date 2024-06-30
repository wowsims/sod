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

	results := make([]*core.SpellResult, hunter.Env.GetNumTargets())

	var ohSpell *core.Spell
	if hunter.AutoAttacks.IsDualWielding {
		ohSpell = hunter.RegisterSpell(core.SpellConfig{
			ActionID:    actionID.WithTag(1),
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeOHSpecial,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: hunter.AutoAttacks.OHConfig().DamageMultiplier * 0.65,
		})
	}

	hunter.CarveMh = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 425711},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

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

		DamageMultiplier: 0.65,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for idx := range results {
				baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				if curTarget == target {
					results[idx] = spell.CalcDamage(sim, curTarget, baseDamage * 1.5, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				} else {
					results[idx] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				}
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}

			if ohSpell != nil {
				ohSpell.Cast(sim, target)

				curTarget := target
				for idx := range results {
					baseDamage := ohSpell.Unit.OHNormalizedWeaponDamage(sim, ohSpell.MeleeAttackPower())
					if curTarget == target {
						results[idx] = ohSpell.CalcDamage(sim, curTarget, baseDamage * 1.5, ohSpell.OutcomeMeleeWeaponSpecialHitAndCrit)
					} else {
						results[idx] = ohSpell.CalcDamage(sim, curTarget, baseDamage, ohSpell.OutcomeMeleeWeaponSpecialHitAndCrit)
					}
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}

				for _, result := range results {
					ohSpell.DealDamage(sim, result)
				}
			}
		},
	})
}
