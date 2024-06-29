package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getConflagrateConfig(rank int) core.SpellConfig {
	hasBackdraftRune := warlock.HasRune(proto.WarlockRune_RuneHelmBackdraft)

	spellId := [5]int32{0, 17962, 18930, 18931, 18932}[rank]
	baseDamageMin := [5]float64{0, 249, 319, 395, 447}[rank]
	baseDamageMax := [5]float64{0, 316, 400, 491, 557}[rank]
	manaCost := [5]float64{0, 165, 200, 230, 255}[rank]
	level := [5]int{0, 0, 48, 54, 60}[rank]

	spCoeff := 0.429

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | WarlockFlagDestruction,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
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
			return warlock.getActiveImmolateSpell(target) != nil || (warlock.ShadowflameDot != nil && warlock.ShadowflameDot.Dot(target).IsActive())
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax)

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() && hasBackdraftRune {
				warlock.BackdraftAura.Activate(sim)
			}

			// Conflag now doesn't consume Immo or Shadowflame when using Backdraft
			if !hasBackdraftRune {
				immoTime := core.NeverExpires
				shadowflameTime := core.NeverExpires

				immoSpell := warlock.getActiveImmolateSpell(target)
				if immoSpell != nil {
					immoDot := immoSpell.Dot(target)
					immoTime = core.TernaryDuration(immoDot.IsActive(), immoDot.RemainingDuration(sim), core.NeverExpires)
				}

				if warlock.Shadowflame != nil {
					sfDot := warlock.ShadowflameDot.Dot(target)
					shadowflameTime = core.TernaryDuration(sfDot.IsActive(), sfDot.RemainingDuration(sim), core.NeverExpires)
				}

				if immoTime < shadowflameTime {
					immoSpell.Dot(target).Deactivate(sim)
				} else {
					warlock.ShadowflameDot.Dot(target).Deactivate(sim)
				}
			}
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
