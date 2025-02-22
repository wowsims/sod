package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

type FerociousBiteRankInfo struct {
	id           int32
	level        int32
	dmgBase      float64
	dmgRange     float64
	dmgPerCombo  float64
	dmgPerEnergy float64
}

var ferociousBiteRanks = []FerociousBiteRankInfo{
	{
		id:           22568,
		level:        32,
		dmgBase:      14.0,
		dmgRange:     16.0,
		dmgPerCombo:  36.0,
		dmgPerEnergy: 1.0,
	},
	{
		id:           22827,
		level:        40,
		dmgBase:      20.0,
		dmgRange:     24.0,
		dmgPerCombo:  59.0,
		dmgPerEnergy: 1.5,
	},
	{
		id:           22828,
		level:        48,
		dmgBase:      30.0,
		dmgRange:     40.0,
		dmgPerCombo:  92.0,
		dmgPerEnergy: 2.0,
	},
	{
		id:           22829,
		level:        56,
		dmgBase:      45.0,
		dmgRange:     50.0,
		dmgPerCombo:  128.0,
		dmgPerEnergy: 2.5,
	},
	{
		id:           31018,
		level:        60,
		dmgBase:      52.0,
		dmgRange:     60.0,
		dmgPerCombo:  147.0,
		dmgPerEnergy: 2.7,
	},
}

func (druid *Druid) registerFerociousBiteSpell() {
	// Add highest available rank for level.
	for rank := len(ferociousBiteRanks) - 1; rank >= 0; rank-- {
		if druid.Level >= ferociousBiteRanks[rank].level {
			config := druid.newFerociousBiteSpellConfig(ferociousBiteRanks[rank])
			druid.FerociousBite = druid.RegisterSpell(Cat, config)
			return
		}
	}
}

func (druid *Druid) newFerociousBiteSpellConfig(rank FerociousBiteRankInfo) core.SpellConfig {
	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_DruidFerociousBite,
		ActionID:       core.ActionID{SpellID: rank.id},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   35,
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

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			comboPoints := float64(druid.ComboPoints())
			attackPower := spell.MeleeAttackPower()
			excessEnergy := druid.CurrentEnergy()

			baseDamage := rank.dmgBase + rank.dmgRange*sim.RandomFloat("Ferocious Bite") +
				rank.dmgPerCombo*comboPoints +
				attackPower*0.05*comboPoints

			if !druid.FerociousBiteExcessEnergyOverride {
				baseDamage += rank.dmgPerEnergy * excessEnergy
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.SpendEnergy(sim, excessEnergy, spell.EnergyMetrics())
				druid.SpendComboPoints(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}
		},
	}
}

func (druid *Druid) CurrentFerociousBiteCost() float64 {
	return druid.FerociousBite.Cost.GetCurrentCost()
}
