package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (druid *Druid) registerLacerateSpell() {
	if !druid.HasRune(proto.DruidRune_RuneLegsLacerate) {
		return
	}
	initialDamageMul := 0.0
	hasGore := druid.HasRune(proto.DruidRune_RuneHelmGore)

	switch druid.Ranged().ID {
	case IdolOfCruelty:
		initialDamageMul += .07
	}
	rageMetrics := druid.NewRageMetrics(core.ActionID{SpellID: 431446})

	druid.Lacerate = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 414644},
		SpellSchool: core.SpellSchoolPhysical,
		SpellCode:   SpellCode_DruidLacerate,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

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
		// TODO: Berserk 3 target lacerate cleave - Saeyon

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) * (.2*float64(druid.LacerateBleed.Dot(target).GetStacks()) + initialDamageMul)
			berserking := druid.BerserkAura.IsActive()

			dotBonusCrit := 0.0
			if druid.LacerateBleed.Dot(target).GetStacks() > 0 {
				dotBonusCrit = druid.FuryOfStormrageCritRatingBonus
			}

			spell.BonusCritRating += dotBonusCrit
			spell.Cost.FlatModifier -= core.TernaryInt32(berserking, 10, 0)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.Cost.FlatModifier += core.TernaryInt32(berserking, 10, 0)
			spell.BonusCritRating -= dotBonusCrit

			if result.Landed() {
				druid.LacerateBleed.Cast(sim, target)

				if hasGore && sim.Proc(0.15, "Gore") {
					druid.AddRage(sim, 10.0, rageMetrics)
					druid.MangleBear.CD.Reset()
				}
			} else {
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
				MaxStacks: 5,
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
