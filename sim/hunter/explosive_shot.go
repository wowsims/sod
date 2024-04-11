package hunter

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) registerExplosiveShotSpell(timer *core.Timer) {
	if !hunter.HasRune(proto.HunterRune_RuneHandsExplosiveShot) {
		return
	}

	actionID := core.ActionID{SpellID: 409552}
	numHits := hunter.Env.GetNumTargets()

	baseLowDamage := hunter.baseRuneAbilityDamage() * 0.36 * 1.15 // Buff from 1/3/2024 - verify with new build and update numbers
	baseHighDamage := hunter.baseRuneAbilityDamage() * 0.54 * 1.15

	manaCostMultiplier := 1 - 0.02*float64(hunter.Talents.Efficiency)
	if hunter.HasRune(proto.HunterRune_RuneChestMasterMarksman) {
		manaCostMultiplier -= 0.25
	}
	hunter.ExplosiveShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		DefenseType:  core.DefenseTypeRanged,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIgnoreResists | core.SpellFlagBinary,
		CastType:     proto.CastType_CastTypeRanged,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.035,
			Multiplier: manaCostMultiplier,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget >= 8
		},

		CritDamageBonus: hunter.mortalShots(),

		DamageMultiplier: hunter.applyTntDamageMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("ExplosiveShot-%d", actionID.SpellID),
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + 0.039*dot.Spell.RangedAttackPower(target)
				dot.Snapshot(target, baseDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + 0.039*spell.RangedAttackPower(target)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)

				if result.Landed() {
					curTarget := target
					for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
						if curTarget != target {
							baseDamage = sim.Roll(baseLowDamage, baseHighDamage) + 0.039*spell.RangedAttackPower(curTarget)
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

func (hunter *Hunter) calcExplosiveShotDamageMultipliers() float64 {
	damageMultiplier := 1.0
	if hunter.HasRune(proto.HunterRune_RuneBracersTNT) {
		damageMultiplier *= 1.1
	}
	return damageMultiplier
}