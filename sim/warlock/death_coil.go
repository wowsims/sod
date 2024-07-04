package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const DeathCoilRanks = 3

func (warlock *Warlock) getDeathCoilBaseConfig(rank int) core.SpellConfig {
	spellId := [DeathCoilRanks + 1]int32{0, 6789, 17925, 17926}[rank]
	baseDamage := [DeathCoilRanks + 1]float64{0, 301, 375, 476}[rank]
	manaCost := [DeathCoilRanks + 1]float64{0, 430, 495, 565}[rank]
	level := [DeathCoilRanks + 1]int{0, 42, 50, 58}[rank]
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
	warlock.DeathCoil = make([]*core.Spell, 0)
	for rank := 1; rank <= DeathCoilRanks; rank++ {
		config := warlock.getDeathCoilBaseConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.DeathCoil = append(warlock.DeathCoil, warlock.GetOrRegisterSpell(config))
		}
	}
}
