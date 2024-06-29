package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const SearingPainRanks = 6

func (warlock *Warlock) getSearingPainBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [SearingPainRanks + 1]float64{0, .396, .429, .429, .429, .429, .429}[rank]
	baseDamage := [SearingPainRanks + 1][]float64{{0}, {38, 47}, {65, 77}, {93, 112}, {131, 155}, {168, 199}, {208, 244}}[rank]
	spellId := [SearingPainRanks + 1]int32{0, 5676, 17919, 17920, 17921, 17922, 17923}[rank]
	manaCost := [SearingPainRanks + 1]float64{0, 45, 68, 91, 118, 141, 168}[rank]
	level := [SearingPainRanks + 1]int{0, 18, 26, 36, 42, 50, 58}[rank]
	castTime := time.Millisecond * 1500

	return core.SpellConfig{
		SpellCode:     SpellCode_WarlockSearingPain,
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagDestruction,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if warlock.MetamorphosisAura != nil && warlock.MetamorphosisAura.IsActive() {
					spell.DefaultCast.CastTime = 0
				} else {
					spell.DefaultCast.CastTime = castTime
				}
			},
		},
		BonusCritRating: 2.0 * float64(warlock.Talents.ImprovedSearingPain) * core.CritRatingPerCritChance,

		DamageMultiplier: 1,
		ThreatMultiplier: 2,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamage[0], baseDamage[1])
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		},
	}
}

func (warlock *Warlock) registerSearingPainSpell() {
	warlock.SearingPain = make([]*core.Spell, 0)
	for rank := 1; rank <= SearingPainRanks; rank++ {
		config := warlock.getSearingPainBaseConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.SearingPain = append(warlock.SearingPain, warlock.GetOrRegisterSpell(config))
		}
	}
}
