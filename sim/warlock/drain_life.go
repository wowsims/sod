package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getDrainLifeBaseConfig(rank int) core.SpellConfig {
	masterChanneler := warlock.HasRune(proto.WarlockRune_RuneChestMasterChanneler)

	spellId := [7]int32{0, 689, 699, 709, 7651, 11699, 11700}[rank]
	spellCoeff := [7]float64{0, .078, .1, .1, .1, .1, .1}[rank]
	baseDamage := [7]float64{0, 10, 17, 29, 41, 55, 71}[rank]
	manaCost := [7]float64{0, 55, 85, 135, 185, 240, 300}[rank]
	level := [7]int{0, 14, 22, 30, 38, 46, 54}[rank]

	ticks := core.TernaryInt32(masterChanneler, 15, 5)

	if masterChanneler {
		manaCost *= 2
	}

	baseDamage *= 1 + 0.02*float64(warlock.Talents.ShadowMastery) + 0.02*float64(warlock.Talents.ImprovedDrainLife)

	actionID := core.ActionID{SpellID: spellId}
	healthMetrics := warlock.NewHealthMetrics(actionID)

	spellConfig := core.SpellConfig{
		ActionID:      actionID,
		SpellSchool:   core.SpellSchoolShadow,
		SpellCode:     SpellCode_DrainLife,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         SpellFlagHaunt | core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				// ChannelTime: channelTime,
			},
		},

		BonusHitRating: float64(warlock.Talents.Suppression) * 2 * core.SpellHitRatingPerHitChance,

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DrainLife-" + warlock.Label + strconv.Itoa(rank),
			},
			NumberOfTicks:       ticks,
			TickLength:          1 * time.Second,
			AffectedByCastSpeed: false,
			BonusCoefficient:    spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)

				// Drain Life heals so it snapshots target modifiers
				dot.SnapshotAttackerMultiplier *= dot.Spell.TargetDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType], true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// TODO: How does it interact with bonus damage taken?
				// Remove target modifiers for the tick only
				dot.Spell.Flags |= core.SpellFlagIgnoreTargetModifiers
				//dot.Spell.Flags ^= core.SpellFlagBinary

				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)

				// revert flag changes
				dot.Spell.Flags ^= core.SpellFlagIgnoreTargetModifiers
				//dot.Spell.Flags |= core.SpellFlagBinary

				health := result.Damage
				if masterChanneler {
					health *= 1.5
				}
				warlock.GainHealth(sim, health, healthMetrics)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--

				dot := spell.Dot(target)
				dot.Apply(sim)

				warlock.EverlastingAfflictionRefresh(sim, target)
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

	if masterChanneler {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    warlock.NewTimer(),
			Duration: 15 * time.Second,
		}
	} else {
		spellConfig.Flags |= core.SpellFlagChanneled
	}

	return spellConfig
}

func (warlock *Warlock) registerDrainLifeSpell() {
	maxRank := 6

	for i := 1; i <= maxRank; i++ {
		config := warlock.getDrainLifeBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.DrainLife = warlock.GetOrRegisterSpell(config)
		}
	}
}
