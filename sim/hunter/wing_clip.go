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
		ClassSpellMask: ClassSpellMask_HunterWingClip,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagBinary,

		Rank:          rank,
		RequiredLevel: level,
		MaxRange:      core.MaxMeleeAttackRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	}
}

func (hunter *Hunter) registerWingClipSpell() {
	rank := map[int32]int{
		25: 1,
		40: 2,
		50: 3,
		60: 3,
	}[hunter.Level]

	config := hunter.getWingClipConfig(rank)
	hunter.WingClip = hunter.GetOrRegisterSpell(config)
}
