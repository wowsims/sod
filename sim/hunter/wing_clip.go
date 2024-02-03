package hunter

import (
	"github.com/wowsims/sod/sim/core"
)

func (hunter *Hunter) getWingClipConfig(rank int) core.SpellConfig {
	spellId := [4]int32{0, 2974, 14267, 14268}[rank]
	baseDamage := [4]float64{0, 5, 25, 50}[rank]
	manaCost := [4]float64{0, 40, 60, 80}[rank]
	level := [4]int{0, 12, 38, 60}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolPhysical,
		ProcMask:      core.ProcMaskMeleeMHSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= 5
		},

		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   hunter.critMultiplier(false, hunter.CurrentTarget),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	}
}

func (hunter *Hunter) registerWingClipSpell() {
	maxRank := 3

	for i := 1; i <= maxRank; i++ {
		config := hunter.getWingClipConfig(i)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.WingClip = hunter.GetOrRegisterSpell(config)
		}
	}
}
