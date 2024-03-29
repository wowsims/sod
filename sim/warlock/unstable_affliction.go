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
	spellCoeff := 0.2

	level := float64(warlock.GetCharacter().Level)
	baseCalc := (6.568597 + 0.672028*level + 0.031721*level*level)
	baseDamage := baseCalc * 0.6

	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)

	warlock.UnstableAffliction = warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.WarlockRune_RuneBracerUnstableAffliction)},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		DefenseType: core.DefenseTypeMagic,
		Flags:       core.SpellFlagAPL | core.SpellFlagHauntSE | core.SpellFlagResetAttackSwing | core.SpellFlagPureDot,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusHitRating: float64(warlock.Talents.Suppression) * 2 * core.SpellHitRatingPerHitChance,

		CritDamageBonus: core.TernaryFloat64(hasPandemicRune, 1, 0),

		DamageMultiplierAdditive: 1 + 0.02*float64(warlock.Talents.ShadowMastery),
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "UnstableAffliction-" + warlock.Label,
			},

			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = baseDamage + (spellCoeff * dot.Spell.SpellDamage())

				if !isRollover {
					dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if hasPandemicRune {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickSnapshotCritCounted)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--

				immoDot := warlock.getActiveImmolateSpell(target)
				if immoDot != nil {
					immoDot.Dot(target).Deactivate(sim)
				}

				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDamage := baseDamage + (spellCoeff * spell.SpellDamage())
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	})
}
