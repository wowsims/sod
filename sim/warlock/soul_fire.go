package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const SoulFireRanks = 2
const SoulFireCastTime = time.Millisecond * 6000

func (warlock *Warlock) getSoulFireBaseConfig(rank int) core.SpellConfig {
	hasDecimationRune := warlock.HasRune(proto.WarlockRune_RuneBootsDecimation)

	spellId := [SoulFireRanks + 1]int32{0, 6353, 17924}[rank]
	baseDamage := [SoulFireRanks + 1][]float64{{0, 0}, {628, 789}, {715, 894}}[rank]
	manaCost := [SoulFireRanks + 1]float64{0, 305, 335}[rank]
	level := [SoulFireRanks + 1]int{0, 48, 56}[rank]
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
				CastTime: SoulFireCastTime,
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
	warlock.SoulFire = make([]*core.Spell, 0)
	for rank := 1; rank <= SoulFireRanks; rank++ {
		config := warlock.getSoulFireBaseConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.SoulFire = append(warlock.SoulFire, warlock.GetOrRegisterSpell(config))
		}
	}
}
