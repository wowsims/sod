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

// See https://www.wowhead.com/classic/spell=436895/s03-tuning-and-overrides-passive-druid
// Modifies Buff Duration +4001:
// Modifies Periodic Damage/Healing Done +51%:
// const RipTicks int32 = 6
const RipTicks int32 = 8

// See https://www.wowhead.com/classic/news/development-notes-for-phase-4-ptr-season-of-discovery-new-runes-class-changes-342896
// - Rake and Rip damage contributions from attack power increased by roughly 50%.
// PTR testing comes out to .0165563 AP scaling per CP
// damageCoefPerCP := 0.01
const RipDamageCoefPerAPPerCP = 0.01

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
	energyCost := 30.0

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_DruidRip,
		ActionID:       core.ActionID{SpellID: ripRank.id},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagPureDot,

		EnergyCost: core.EnergyCostOptions{
			Cost:   energyCost,
			Refund: 0,
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

		DamageMultiplier:                    1,
		PeriodicDamageMultiplierAdditivePct: 50, // https://www.wowhead.com/classic/spell=436895/s03-tuning-and-overrides-passive-druid
		ThreatMultiplier:                    1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rip",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					druid.BleedsActive[aura.Unit.UnitIndex]++
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					druid.BleedsActive[aura.Unit.UnitIndex]--
				},
			},
			NumberOfTicks: RipTicks,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				cp := float64(druid.ComboPoints())
				ap := dot.Spell.MeleeAttackPower()

				cpScaling := core.TernaryFloat64(cp == 5, 4, cp)
				baseDamage := (ripRank.dmgTickBase + ripRank.dmgTickPerCombo*cp + RipDamageCoefPerAPPerCP*ap*cpScaling)
				dot.Snapshot(target, baseDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if druid.AllowRakeRipDoTCrits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHitNoHitCounter)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.Apply(sim)
				druid.SpendComboPoints(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	}
}

func (druid *Druid) CurrentRipCost() float64 {
	return druid.Rip.Cost.GetCurrentCost()
}
