package druid

import (
	"time"

	"github.com/wowsims/sod/sim/common/sod/item_sets"
	"github.com/wowsims/sod/sim/core"
)

const WrathRanks = 8

var WrathSpellId = [WrathRanks + 1]int32{0, 5176, 5177, 5178, 5179, 5180, 6780, 8905, 9912}
var WrathBaseDamage = [WrathRanks + 1][]float64{{0}, {13, 16}, {28, 33}, {48, 57}, {69, 79}, {108, 123}, {148, 167}, {198, 221}, {248, 277}}
var WrathSpellCoeff = [WrathRanks + 1]float64{0, 0.123, 0.231, 0.443, 0.571, 0.571, 0.571, 0.571, 0.571}
var WrathManaCost = [WrathRanks + 1]float64{0, 20, 35, 55, 70, 100, 125, 155, 180}
var WrathCastTime = [WrathRanks + 1]int{0, 1500, 1700, 2000, 2000, 2000, 2000, 2000, 2000}
var WrathLevel = [WrathRanks + 1]int{0, 1, 6, 14, 22, 30, 38, 46, 54}

func (druid *Druid) registerWrathSpell() {
	druid.Wrath = make([]*DruidSpell, WrathRanks+1)

	for rank := 1; rank <= WrathRanks; rank++ {
		config := druid.newWrathSpellConfig(rank)

		if config.RequiredLevel <= int(druid.Level) {
			druid.Wrath[rank] = druid.RegisterSpell(Humanoid|Moonkin, config)
		}
	}
}

func (druid *Druid) newWrathSpellConfig(rank int) core.SpellConfig {
	spellId := WrathSpellId[rank]
	baseDamageLow := WrathBaseDamage[rank][0]
	baseDamageHigh := WrathBaseDamage[rank][1]
	spellCoeff := WrathSpellCoeff[rank]
	manaCost := WrathManaCost[rank]
	castTime := WrathCastTime[rank]
	level := WrathLevel[rank]

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellCode:   SpellCode_DruidWrath,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing,

		RequiredLevel: level,
		Rank:          rank,
		MissileSpeed:  20,

		ManaCost: core.ManaCostOptions{
			FlatCost: core.TernaryFloat64(druid.FuryOfStormrageAura != nil, 0, manaCost),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*time.Duration(castTime) - time.Millisecond*100*time.Duration(druid.Talents.ImprovedWrath),
			},
			CastTime: druid.NaturesGraceCastTime(),
		},

		DamageMultiplier: 1,
		CritMultiplier:   druid.VengeanceCritMultiplier(),
		BonusCritRating:  core.TernaryFloat64(druid.HasSetBonus(item_sets.ItemSetInsulatedSorcerorLeather, 3), 2, 0) * core.CritRatingPerCritChance,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)*druid.MoonfuryDamageMultiplier() + spellCoeff*spell.SpellDamage()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.DidCrit() && druid.NaturesGraceProcAura != nil {
				druid.NaturesGraceProcAura.Activate(sim)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	}
}
