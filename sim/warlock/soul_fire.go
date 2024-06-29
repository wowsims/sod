package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getSoulFireBaseConfig(rank int) core.SpellConfig {
	hasDecimationRune := warlock.HasRune(proto.WarlockRune_RuneCloakDecimation)

	spellId := [3]int32{0, 6353, 17924}[rank]
	baseDamage := [3][]float64{{0, 0}, {628, 789}, {715, 894}}[rank]
	manaCost := [3]float64{0, 305, 335}[rank]
	level := [3]int{0, 48, 56}[rank]
	spellCoeff := 1.0

	config := core.SpellConfig{
		SpellCode:     SpellCode_WarlockSoulFire,
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagDestruction,
		RequiredLevel: level,
		Rank:          rank,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 6000,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamage[0], baseDamage[1])
			results := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, results)
			})
		},
	}

	if !hasDecimationRune {
		config.Cast.CD = core.Cooldown{
			Timer:    warlock.NewTimer(),
			Duration: time.Minute,
		}
	}

	return config
}

func (warlock *Warlock) registerSoulFireSpell() {
	maxRank := 2
	warlock.SoulFire = make([]*core.Spell, maxRank+1)

	for i := 1; i <= maxRank; i++ {
		config := warlock.getSoulFireBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.SoulFire[i] = warlock.GetOrRegisterSpell(config)
		}
	}
}
