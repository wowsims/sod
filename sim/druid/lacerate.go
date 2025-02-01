package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var LacerateDotMaxStacks int32 = 5

func (druid *Druid) registerLacerateSpell() {
	if !druid.HasRune(proto.DruidRune_RuneLegsLacerate) {
		return
	}

	initialDamageMul := 0.0

	switch druid.Ranged().ID {
	case IdolOfCruelty:
		initialDamageMul += .07
	}

	druid.Lacerate = druid.RegisterSpell(Bear, core.SpellConfig{
		ClassSpellMask: ClassSpellMask_DruidLacerate,
		ActionID:       core.ActionID{SpellID: 414644},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   10,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 3.33,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			berserking := druid.BerserkAura.IsActive()
			targetCount := core.TernaryInt32(berserking, 3, 1)
			numHits := min(targetCount, druid.Env.GetNumTargets())
			results := make([]*core.SpellResult, numHits)

			for idx := range results {
				activeStacks := druid.LacerateBleed.Dot(target).GetStacks() + 1
				activeStacks = core.TernaryInt32(activeStacks > LacerateDotMaxStacks, LacerateDotMaxStacks, activeStacks)
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) * (.2*float64(activeStacks) + initialDamageMul)

				results[idx] = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

				if results[idx].Landed() {
					druid.LacerateBleed.Cast(sim, target)
				}
				target = sim.Environment.NextTargetUnit(target)
			}
			if !berserking && !results[0].Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (druid *Druid) registerLacerateBleedSpell() {
	if !druid.HasRune(proto.DruidRune_RuneLegsLacerate) {
		return
	}
	tickDamage := 29.8312

	switch druid.Ranged().ID {
	case IdolOfCruelty:
		tickDamage += 7.0
	}

	druid.LacerateBleed = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 414647},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 3.33,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Lacerate",
				MaxStacks: LacerateDotMaxStacks,
				Duration:  time.Second * 15,
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = tickDamage
				dot.SnapshotBaseDamage *= float64(dot.Aura.GetStacks())

				if !isRollover {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType]
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.Dot(target)
			if dot.IsActive() {
				dot.Refresh(sim)
				dot.AddStack(sim)
				dot.TakeSnapshot(sim, true)
			} else {
				dot.Apply(sim)
				dot.SetStacks(sim, 1)
				dot.TakeSnapshot(sim, true)
			}
		},
	})
}
