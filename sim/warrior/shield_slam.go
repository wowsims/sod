package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const ShieldSlamRanks = 4

var ShieldSlamSpellId = [ShieldSlamRanks + 1]int32{0, 23922, 23923, 23924, 23925}
var ShieldSlamBaseDamage = [ShieldSlamRanks + 1][]float64{{0, 0}, {225, 352}, {264, 276}, {303, 317}, {342, 358}}
var ShieldSlamLevel = [ShieldSlamRanks + 1]int{0, 40, 48, 54, 60}

// TODO: Classic Update
func (warrior *Warrior) registerShieldSlamSpell() {
	if !warrior.Talents.ShieldSlam || warrior.Level < 40 {
		return
	}

	rank := []int{
		40: 1,
		50: 2,
		60: 4,
	}[warrior.Level]
	actionID := core.ActionID{SpellID: ShieldSlamSpellId[rank]}
	basedamageLow := ShieldSlamBaseDamage[rank][0]
	basedamageHigh := ShieldSlamBaseDamage[rank][1]
	cooldown := time.Second * 6

	warrior.ShieldSlam = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial, // TODO: Is this right?
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   20,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldown,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock
		},

		BonusCritRating:  5 * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1.3,
		FlatThreatBonus:  770,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: Verify that this bypass behavior and DR curve are correct
			sbvMod := warrior.PseudoStats.BlockValueMultiplier
			sbvMod /= 1

			sbv := warrior.BlockValue() / sbvMod
			sbv = sbvMod * (core.TernaryFloat64(sbv <= 1960.0, sbv, 0.0) + core.TernaryFloat64(sbv > 1960.0 && sbv <= 3160.0, 0.09333333333*sbv+1777.06666667, 0.0) + core.TernaryFloat64(sbv > 3160.0, 2072.0, 0.0))

			baseDamage := sim.Roll(basedamageLow, basedamageHigh) + sbv
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
