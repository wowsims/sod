package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) registerShadowflameSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneBootsShadowflame) {
		return
	}

	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)

	baseSpellCoeff := 0.107
	dotSpellCoeff := 0.107
	baseDamage := warlock.baseRuneAbilityDamage() * 2.26
	dotDamage := warlock.baseRuneAbilityDamage() * 0.61

	// TODO: Probably merge this into the base spell now
	warlock.ShadowflameDot = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 426325},
		SpellSchool: core.SpellSchoolFire | core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplierAdditive: 1 + 0.02*float64(warlock.Talents.Emberstorm),
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Shadowflame" + warlock.Label,
			},

			NumberOfTicks:    4,
			TickLength:       time.Second * 2,
			BonusCoefficient: dotSpellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				var result *core.SpellResult
				if hasPandemicRune {
					// We add the crit damage bonus and remove it after the call to not affect the initial damage portion of the spell
					dot.Spell.CritDamageBonus += 1
					result = dot.CalcSnapshotDamage(sim, target, dot.OutcomeTickSnapshotCrit)
					dot.Spell.CritDamageBonus -= 1
				} else {
					result = dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
				}
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			spell.Dot(target).Apply(sim)
		},
	})

	warlock.Shadowflame = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 426320},
		SpellSchool: core.SpellSchoolFire | core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.27,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: 15 * time.Second,
			},
		},

		BonusCritRating: float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,

		CritDamageBonus: warlock.ruin(),

		DamageMultiplierAdditive: 1 + 0.02*float64(warlock.Talents.ShadowMastery),
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         baseSpellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				warlock.ShadowflameDot.Cast(sim, target)
			}
		},
	})
}
