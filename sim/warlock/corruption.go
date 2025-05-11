package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const CorruptionRanks = 7

func (warlock *Warlock) getCorruptionConfig(rank int) core.SpellConfig {
	dotTickCoeff := [CorruptionRanks + 1]float64{0, .08, .155, .167, .167, .167, .167, .167}[rank] // per tick
	ticks := [CorruptionRanks + 1]int32{0, 4, 5, 6, 6, 6, 6, 6}[rank]
	baseDamage := [CorruptionRanks + 1]float64{0, 40, 90, 222, 324, 486, 666, 822}[rank] / float64(ticks)
	spellId := [CorruptionRanks + 1]int32{0, 172, 6222, 6223, 7648, 11671, 11672, 25311}[rank]
	manaCost := [CorruptionRanks + 1]float64{0, 35, 55, 100, 160, 225, 290, 340}[rank]
	level := [CorruptionRanks + 1]int{0, 4, 14, 24, 34, 44, 54, 60}[rank]

	castTime := time.Millisecond * (2000 - (400 * time.Duration(warlock.Talents.ImprovedCorruption)))
	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)

	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolShadow,
		ClassSpellMask: ClassSpellMask_WarlockCorruption,
		ProcMask:       core.ProcMaskSpellDamage,
		DefenseType:    core.DefenseTypeMagic,
		Flags:          core.SpellFlagAPL | core.SpellFlagResetAttackSwing | core.SpellFlagPureDot | WarlockFlagAffliction | WarlockFlagHaunt,
		Rank:           rank,
		RequiredLevel:  level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: castTime,
				GCD:      core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Corruption-" + warlock.Label + strconv.Itoa(rank),
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.activeEffects[aura.Unit.UnitIndex]++
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.activeEffects[aura.Unit.UnitIndex]--
				},
			},

			NumberOfTicks:    ticks,
			TickLength:       time.Second * 3,
			BonusCoefficient: dotTickCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)
				if !isRollover {
					if warlock.zilaGularAura.IsActive() {
						dot.SnapshotAttackerMultiplier *= 1.10
						warlock.zilaGularAura.Deactivate(sim)
					}
				}
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
				spell.Dot(target).ApplyOrReset(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				if hasPandemicRune {
					return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicCrit)
				}
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDamage := baseDamage / float64(ticks)
				if hasPandemicRune {
					return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
				}
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}
}

func (warlock *Warlock) registerCorruptionSpell() {
	warlock.Corruption = make([]*core.Spell, 0)

	maxRank := core.TernaryInt(core.IncludeAQ, CorruptionRanks, CorruptionRanks-1)
	for rank := 1; rank <= maxRank; rank++ {
		config := warlock.getCorruptionConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.Corruption = append(warlock.Corruption, warlock.GetOrRegisterSpell(config))
		}
	}
}
