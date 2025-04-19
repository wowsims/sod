package hunter

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) registerExplosiveShotSpell() {
	if !hunter.HasRune(proto.HunterRune_RuneHandsExplosiveShot) {
		return
	}

	actionID := core.ActionID{SpellID: 409552}
	numHits := hunter.Env.GetNumTargets()

	baseLowDamage := hunter.baseRuneAbilityDamage() * 0.36 * 1.15 * 1.5  // 15% Buff from 1/3/2024 - verify with new build and update numbers
	baseHighDamage := hunter.baseRuneAbilityDamage() * 0.54 * 1.15 * 1.5 // Second 50% buff from 23/4/2024

	hunter.ExplosiveShot = hunter.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterExplosiveShot,
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		CastType:       proto.CastType_CastTypeRanged,
		DefenseType:    core.DefenseTypeRanged,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		MinRange:     core.MinRangedAttackRange,
		MaxRange:     core.MaxRangedAttackRange,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.035,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		CritDamageBonus: hunter.mortalShots(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("ExplosiveShot-%d", actionID.SpellID),
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + 0.039*dot.Spell.RangedAttackPower(target, false)
				dot.Snapshot(target, baseDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + 0.039*spell.RangedAttackPower(target, false)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)

				if result.Landed() {
					curTarget := target
					for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
						if curTarget != target {
							baseDamage = sim.Roll(baseLowDamage, baseHighDamage) + 0.039*spell.RangedAttackPower(curTarget, false)
							spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeRangedCritOnly)
						}

						dot := spell.Dot(curTarget)
						dot.Apply(sim)

						curTarget = sim.Environment.NextTargetUnit(curTarget)
					}
				}
			})
		},
	})
}
