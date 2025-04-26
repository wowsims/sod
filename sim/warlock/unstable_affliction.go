package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) registerUnstableAfflictionSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneBracerUnstableAffliction) {
		return
	}

	warlock.UnstableAffliction = warlock.GetOrRegisterSpell(warlock.getUnstableAfflictionConfig())
}

func (warlock *Warlock) getUnstableAfflictionConfig() core.SpellConfig {
	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)

	// TODO: Verify numbers after tooltips update
	// 2024-11-22 +120% damage
	baseDamage := warlock.baseRuneAbilityDamage() * 1.1 * 2.20

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_WarlockUnstableAffliction,
		ActionID:       core.ActionID{SpellID: int32(proto.WarlockRune_RuneBracerUnstableAffliction)},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		DefenseType:    core.DefenseTypeMagic,
		Flags:          core.SpellFlagAPL | WarlockFlagHaunt | core.SpellFlagBinary | core.SpellFlagResetAttackSwing | core.SpellFlagPureDot | WarlockFlagAffliction,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.2,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "UnstableAffliction-" + warlock.Label,
			},

			NumberOfTicks:    6,
			TickLength:       time.Second * 3,
			BonusCoefficient: 0.2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if hasPandemicRune {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				// UA is mutually exclusive with Immolate
				immoDot := warlock.getActiveImmolateSpell(target)
				if immoDot != nil {
					immoDot.Dot(target).Deactivate(sim)
				}

				spell.Dot(target).ApplyOrReset(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}
}
