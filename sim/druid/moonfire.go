package druid

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
)

const MoonfireRanks = 10

var MoonfireSpellId = [MoonfireRanks + 1]int32{0, 8921, 8924, 8925, 8926, 8927, 8928, 8929, 9833, 9834, 9835}
var MoonfiresSpellCoeff = [MoonfireRanks + 1]float64{0, .06, .094, .128, .15, .15, .15, .15, .15, .15, .15}
var MoonfiresSellDotCoeff = [MoonfireRanks + 1]float64{0, .052, .081, .111, .13, .13, .13, .13, .13, .13, .13}
var MoonfireBaseDamage = [MoonfireRanks + 1][]float64{{0}, {9, 12}, {17, 21}, {30, 37}, {44, 53}, {70, 82}, {91, 108}, {117, 137}, {143, 168}, {172, 200}, {195, 228}}
var MoonfireBaseDotDamage = [MoonfireRanks + 1]float64{0, 12, 32, 52, 80, 124, 164, 212, 264, 320, 384}
var MoonfireManaCost = [MoonfireRanks + 1]float64{0, 25, 50, 75, 105, 150, 190, 235, 280, 325, 375}
var MoonfireLevel = [MoonfireRanks + 1]int{0, 4, 10, 16, 22, 28, 34, 40, 46, 52, 58}

func (druid *Druid) registerMoonfireSpell() {
	druid.Moonfire = make([]*DruidSpell, MoonfireRanks+1)

	for rank := 1; rank <= MoonfireRanks; rank++ {
		config := druid.getMoonfireBaseConfig(rank)

		if config.RequiredLevel <= int(druid.Level) {
			druid.Moonfire[rank] = druid.RegisterSpell(Humanoid|Moonkin, config)
		}
	}
}

func (druid *Druid) getMoonfireBaseConfig(rank int) core.SpellConfig {
	spellId := MoonfireSpellId[rank]
	spellCoeff := MoonfiresSpellCoeff[rank]
	spellDotCoeff := MoonfiresSellDotCoeff[rank]
	baseDamageLow := MoonfireBaseDamage[rank][0]
	baseDamageHigh := MoonfireBaseDamage[rank][1]
	baseDotDamage := MoonfireBaseDotDamage[rank]
	manaCost := MoonfireManaCost[rank]
	level := MoonfireLevel[rank]

	ticks := core.TernaryInt32(rank < 2, 3, 4)

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost * (1 - 0.03*float64(druid.Talents.Moonglow)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 0,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    fmt.Sprintf("Moonfire (Rank %d)", rank),
				ActionID: core.ActionID{SpellID: spellId},
			},
			NumberOfTicks: ticks,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = (baseDotDamage/float64(ticks))*druid.MoonfuryDamageMultiplier() + spellDotCoeff*dot.Spell.SpellDamage()
				dot.SnapshotAttackerMultiplier = 1 // dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		BonusCritRating:  druid.ImprovedMoonfireCritBonus() * core.SpellCritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   druid.VengeanceCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)*druid.MoonfuryDamageMultiplier()*druid.ImprovedMoonfireDamageMultiplier() + spellCoeff*spell.SpellDamage()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	}
}
