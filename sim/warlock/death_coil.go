package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warlock *Warlock) getDeathCoilBaseConfig(rank int) core.SpellConfig {
	spellId := [4]int32{0, 6789, 17925, 17926}[rank]
	baseDamage := [4]float64{0, 301, 375, 476}[rank]
	manaCost := [4]float64{0, 430, 495, 565}[rank]
	level := [4]int{0, 42, 50, 58}[rank]
	spellCoeff := 0.214

	baseDamage *= 1 + warlock.shadowMasteryBonus()

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing | core.SpellFlagBinary | WarlockFlagAffliction,
		RequiredLevel: level,
		Rank:          rank,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, results)
			})
		},
	}
}

func (warlock *Warlock) registerDeathCoilSpell() {
	maxRank := 3

	for i := 1; i <= maxRank; i++ {
		config := warlock.getDeathCoilBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.DeathCoil = warlock.GetOrRegisterSpell(config)
		}
	}
}
