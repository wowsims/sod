package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) newWhirlwindSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: 1680}
	procMask := core.ProcMaskMeleeSpecial
	if !isMH {
		procMask = core.ProcMaskMeleeOHSpecial
	}

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID.WithTag(int32(core.Ternary(isMH, 1, 2))),
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1.25,
		BonusCoefficient: 1,
	})
}

func (warrior *Warrior) registerWhirlwindSpell() {
	if warrior.Level < 36 {
		return
	}

	hasConsumedByRageRune := warrior.HasRune(proto.WarriorRune_RuneConsumedByRage)

	warrior.WhirlwindMH = warrior.newWhirlwindSpell(true)
	if hasConsumedByRageRune {
		warrior.WhirlwindOH = warrior.newWhirlwindSpell(false)
	}

	actionID := core.ActionID{SpellID: 1680}
	results := make([]*core.SpellResult, min(4, warrior.Env.GetNumTargets()))

	warrior.Whirlwind = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagAPL | SpellFlagBloodSurge,

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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			doOHAttack := hasConsumedByRageRune && warrior.AutoAttacks.IsDualWielding && warrior.IsEnraged()

			aoeTarget := target
			for idx := range results {
				warrior.WhirlwindMH.Cast(sim, aoeTarget)
				baseDamage := warrior.WhirlwindMH.Unit.MHNormalizedWeaponDamage(sim, warrior.WhirlwindMH.MeleeAttackPower())
				results[idx] = warrior.WhirlwindMH.CalcDamage(sim, aoeTarget, baseDamage, warrior.WhirlwindMH.OutcomeMeleeWeaponSpecialHitAndCrit)
				aoeTarget = sim.Environment.NextTargetUnit(aoeTarget)
			}

			for _, result := range results {
				warrior.WhirlwindMH.DealDamage(sim, result)
			}

			if doOHAttack {
				aoeTarget := target
				for idx := range results {
					warrior.WhirlwindOH.Cast(sim, aoeTarget)
					baseDamage := warrior.WhirlwindOH.Unit.OHNormalizedWeaponDamage(sim, warrior.WhirlwindOH.MeleeAttackPower())
					results[idx] = warrior.WhirlwindOH.CalcDamage(sim, aoeTarget, baseDamage, warrior.WhirlwindOH.OutcomeMeleeWeaponSpecialHitAndCrit)
					aoeTarget = sim.Environment.NextTargetUnit(aoeTarget)
				}

				for _, result := range results {
					warrior.WhirlwindOH.DealDamage(sim, result)
				}
			}
		},
	})
}
