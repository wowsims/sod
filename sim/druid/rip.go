package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

type RipRankInfo struct {
	id              int32
	level           int32
	dmgTickBase     float64
	dmgTickPerCombo float64
}

var ripRanks = []RipRankInfo{
	{
		id:              1079,
		level:           20,
		dmgTickBase:     3.0,
		dmgTickPerCombo: 4.0,
	},
	{
		id:              9492,
		level:           28,
		dmgTickBase:     4.0,
		dmgTickPerCombo: 7.0,
	},
	{
		id:              9493,
		level:           36,
		dmgTickBase:     6.0,
		dmgTickPerCombo: 9.0,
	},
	{
		id:              9752,
		level:           44,
		dmgTickBase:     9.0,
		dmgTickPerCombo: 14.0,
	},
	{
		id:              9894,
		level:           52,
		dmgTickBase:     12.0,
		dmgTickPerCombo: 20.0,
	},
	{
		id:              9896,
		level:           60,
		dmgTickBase:     17.0,
		dmgTickPerCombo: 28.0,
	},
}

var RipTicks int32 = 8

func (druid *Druid) registerRipSpell() {
	// Add highest available Rip rank for level.
	for rank := len(ripRanks) - 1; rank >= 0; rank-- {
		if druid.Level >= ripRanks[rank].level {
			config := druid.newRipSpellConfig(ripRanks[rank])
			druid.Rip = druid.RegisterSpell(Cat, config)
			return
		}
	}
}

func (druid *Druid) newRipSpellConfig(ripRank RipRankInfo) core.SpellConfig {
	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: ripRank.id},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial | core.ProcMaskSuppressedExtraAttackAura,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   30,
			Refund: 0,
			//RefundMetrics: druid.PrimalPrecisionRecoveryMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return druid.ComboPoints() > 0
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			NumberOfTicks: RipTicks,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				cp := float64(druid.ComboPoints())
				ap := dot.Spell.MeleeAttackPower()

				cpScaling := core.TernaryFloat64(cp == 5, 4, cp)

				dot.SnapshotBaseDamage = (ripRank.dmgTickBase + ripRank.dmgTickPerCombo*cp + 0.01*ap*cpScaling)

				if !isRollover {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType]
					dot.SnapshotCritChance = 0
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				dot := spell.Dot(target)
				dot.NumberOfTicks = RipTicks
				dot.Apply(sim)
				druid.SpendComboPoints(sim, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	}
}

func (druid *Druid) CurrentRipCost() float64 {
	return druid.Rip.ApplyCostModifiers(druid.Rip.DefaultCast.Cost)
}
