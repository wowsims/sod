package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getSiphonLifeBaseConfig(rank int) core.SpellConfig {
	spellId := [5]int32{0, 18265, 18879, 18880, 18881}[rank]
	baseDamage := [5]float64{0, 15, 22, 33, 45}[rank]
	manaCost := [5]float64{0, 150, 205, 285, 365}[rank]
	level := [5]int{0, 0, 38, 48, 58}[rank]

	spellCoeff := 0.05
	actionID := core.ActionID{SpellID: spellId}
	healthMetrics := warlock.NewHealthMetrics(actionID)

	baseDamage *= 1 + 0.02*float64(warlock.Talents.ShadowMastery)

	hasInvocationRune := warlock.HasRune(proto.WarlockRune_RuneBeltInvocation)
	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)

	return core.SpellConfig{
		ActionID:      actionID,
		SpellSchool:   core.SpellSchoolShadow,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         SpellFlagHaunt | core.SpellFlagAPL | core.SpellFlagResetAttackSwing | core.SpellFlagBinary,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusHitRating: float64(warlock.Talents.Suppression) * 2 * core.SpellHitRatingPerHitChance,

		CritDamageBonus: core.TernaryFloat64(hasPandemicRune, 1, 0),

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SiphonLife-" + warlock.Label + strconv.Itoa(rank),
			},
			NumberOfTicks:       10,
			TickLength:          3 * time.Second,
			AffectedByCastSpeed: false,
			BonusCoefficient:    spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)

				if !isRollover {
					// Siphon Life heals so it snapshots target modifiers
					dot.SnapshotAttackerMultiplier *= dot.Spell.TargetDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType], true)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// TODO: interaction with bonus damage taken?
				// Remove target modifiers for the tick only
				dot.Spell.Flags |= core.SpellFlagIgnoreTargetModifiers

				var result *core.SpellResult
				if hasPandemicRune {
					result = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickSnapshotCritCounted)
				} else {
					result = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
				}

				// revert flag changes
				dot.Spell.Flags ^= core.SpellFlagIgnoreTargetModifiers

				health := result.Damage
				warlock.GainHealth(sim, health, healthMetrics)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--

				if hasInvocationRune && spell.Dot(target).IsActive() {
					warlock.InvocationRefresh(sim, spell.Dot(target))
				}

				spell.Dot(target).Apply(sim)
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}
}

func (warlock *Warlock) registerSiphonLifeSpell() {
	if !warlock.Talents.SiphonLife {
		return
	}

	maxRank := 4

	for i := 1; i <= maxRank; i++ {
		config := warlock.getSiphonLifeBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.SiphonLife = warlock.GetOrRegisterSpell(config)
		}
	}
}
