package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (druid *Druid) registerLacerateSpell() {
	tickDamage := 20.0
	initialDamage := 88.0
	initialDamageMul := 1.0
	additionalTickDamage := 0

	switch druid.Ranged().ID {
	case IdolOfBrutality:
		initialDamageMul *= .07
		additionalTickDamage += 7
	}

	druid.Lacerate = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 414644},
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

		DamageMultiplier: initialDamageMul,
		ThreatMultiplier: 3.5,

		// everything up to here should be accurate

		Dot: core.DotConfig{
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = tickDamage + float64(dot.Aura.GetStacks() * int32(additionalTickDamage))
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
			baseDamage := initialDamage + 0.01*spell.MeleeAttackPower()
			if druid.BleedCategories.Get(target).AnyActive() {
				baseDamage *= 1.3
			}

			spell.DamageMultiplier = initialDamageMul
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.FlatThreatBonus = 0

			if result.Landed() {
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

			spell.DealDamage(sim, result)
		},
	})
}
