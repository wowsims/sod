package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warlock *Warlock) getConflagrateConfig(rank int) core.SpellConfig {
	spellId := [5]int32{0, 17962, 18930, 18931, 18932}[rank]
	baseDamageMin := [5]float64{0, 249, 319, 395, 447}[rank]
	baseDamageMax := [5]float64{0, 316, 400, 491, 557}[rank]
	manaCost := [5]float64{0, 165, 200, 230, 255}[rank]
	level := [5]int{0, 0, 48, 54, 60}[rank]

	spCoeff := 0.429

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.Immolate.Dot(target).IsActive()
		},

		BonusCritRating:          0 + float64(warlock.Talents.Devastation)*core.CritRatingPerCritChance,
		DamageMultiplierAdditive: 1 + 0.02*float64(warlock.Talents.Emberstorm),
		CritMultiplier:           warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Ruin, 1, 0)),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// takes the SP of the immolate (or shadowflame) dot on the target
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax) + spCoeff*spell.SpellPower()

			if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
				baseDamage *= 1.4
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if !result.Landed() {
				return
			}
			warlock.Immolate.Dot(target).Deactivate(sim)
			//warlock.ShadowflameDot.Deactivate(sim)
		},
	}
}

func (warlock *Warlock) registerConflagrateSpell() {
	if !warlock.Talents.Conflagrate {
		return
	}

	maxRank := 4

	for i := 1; i <= maxRank; i++ {
		config := warlock.getConflagrateConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.Conflagrate = warlock.GetOrRegisterSpell(config)
		}
	}
}
