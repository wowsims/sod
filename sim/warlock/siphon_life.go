package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const SiphonLifeRanks = 4

func (warlock *Warlock) getSiphonLifeBaseConfig(rank int) core.SpellConfig {
	spellId := [SiphonLifeRanks + 1]int32{0, 18265, 18879, 18880, 18881}[rank]
	baseDamage := [SiphonLifeRanks + 1]float64{0, 15, 22, 33, 45}[rank]
	manaCost := [SiphonLifeRanks + 1]float64{0, 150, 205, 285, 365}[rank]
	level := [SiphonLifeRanks + 1]int{0, 0, 38, 48, 58}[rank]

	spellCoeff := 0.05
	actionID := core.ActionID{SpellID: spellId}

	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)

	healingSpell := warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId}.WithTag(1),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

		DamageMultiplier: 1,
		ThreatMultiplier: 0,
	})
	return core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_WarlockSiphonLife,
		SpellSchool:    core.SpellSchoolShadow,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagResetAttackSwing | core.SpellFlagBinary | WarlockFlagAffliction | WarlockFlagHaunt,
		RequiredLevel:  level,
		Rank:           rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

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
					result = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					result = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				}

				// revert flag changes
				dot.Spell.Flags ^= core.SpellFlagIgnoreTargetModifiers

				healingSpell.CalcAndDealHealing(sim, healingSpell.Unit, result.Damage, healingSpell.OutcomeHealing)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).ApplyOrReset(sim)
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

	warlock.SiphonLife = make([]*core.Spell, 0)
	for rank := 1; rank <= SiphonLifeRanks; rank++ {
		config := warlock.getSiphonLifeBaseConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.SiphonLife = append(warlock.SiphonLife, warlock.GetOrRegisterSpell(config))
		}
	}
}
