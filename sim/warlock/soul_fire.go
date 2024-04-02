package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warlock *Warlock) getSoulFireBaseConfig(rank int) core.SpellConfig {
	spellId := [3]int32{0, 6353, 17924}[rank]
	baseDamage := [3][]float64{{0, 0}, {628, 789}, {715, 894}}[rank]
	manaCost := [3]float64{0, 305, 335}[rank]
	level := [3]int{0, 48, 56}[rank]
	spellCoeff := 1.0

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		RequiredLevel: level,
		Rank:          rank,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (6000 - 400*time.Duration(warlock.Talents.Bane)),
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute,
			},
		},

		BonusCritRating: float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,

		CritDamageBonus: warlock.ruin(),

		DamageMultiplierAdditive: 1 + 0.02*float64(warlock.Talents.Emberstorm),
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellDamage()

			if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
				damage *= warlock.getLakeOfFireMultiplier()
			}

			results := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, results)
			})
		},
	}
}

func (warlock *Warlock) registerSoulFireSpell() {
	maxRank := 2

	for i := 1; i <= maxRank; i++ {
		config := warlock.getSoulFireBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.SoulFire = warlock.GetOrRegisterSpell(config)
		}
	}
}
