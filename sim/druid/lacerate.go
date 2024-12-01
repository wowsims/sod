package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (druid *Druid) registerLacerateDirectSpell() {
	initialDamageMul := 1.0

	switch druid.Ranged().ID {
	case IdolOfCruelty:
		initialDamageMul += .07
	}

	druid.LacerateDirect = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 414644},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: initialDamageMul,
		ThreatMultiplier: 3.25,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) * .2) * float64(druid.Lacerate.Dot(target).GetStacks())

			spell.DamageMultiplier = initialDamageMul
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

func (druid *Druid) registerLacerateSpell() {
	tickDamage := 20.0

	switch druid.Ranged().ID {
	case IdolOfCruelty:
		tickDamage += 7.0
	}

	druid.Lacerate = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 414647},
		SpellSchool: core.SpellSchoolPhysical,
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
		ThreatMultiplier: 3.4,

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
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, 0, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.LacerateDirect.Cast(sim, target)
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
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
