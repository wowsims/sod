package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ImmolateRanks = 8
const ImmolateCastTime = time.Millisecond * 2000

func (warlock *Warlock) getImmolateConfig(rank int) core.SpellConfig {
	directCoeff := [ImmolateRanks + 1]float64{0, .058, .125, .2, .2, .2, .2, .2, .2}[rank]
	dotCoeff := [ImmolateRanks + 1]float64{0, .037, .081, .13, .13, .13, .13, .13, .13}[rank]
	baseDamage := [ImmolateRanks + 1]float64{0, 11, 24, 53, 101, 148, 208, 258, 279}[rank]
	dotDamage := [ImmolateRanks + 1]float64{0, 20, 40, 90, 165, 255, 365, 485, 510}[rank] / 5
	spellId := [ImmolateRanks + 1]int32{0, 348, 707, 1094, 2941, 11665, 11667, 11668, 25309}[rank]
	manaCost := [ImmolateRanks + 1]float64{0, 25, 45, 90, 155, 220, 295, 370, 380}[rank]
	level := [ImmolateRanks + 1]int{0, 1, 10, 20, 30, 40, 50, 60, 60}[rank]

	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)
	hasUnstableAffliction := warlock.HasRune(proto.WarlockRune_RuneBracerUnstableAffliction)
	hasShadowflameRune := warlock.HasRune(proto.WarlockRune_RuneBootsShadowflame)

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_WarlockImmolate,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagResetAttackSwing | core.SpellFlagBinary | WarlockFlagDestruction,

		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: ImmolateCastTime,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: directCoeff,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolate-" + warlock.Label + strconv.Itoa(rank),
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.activeEffects[aura.Unit.UnitIndex]++
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.activeEffects[aura.Unit.UnitIndex]--
				},
			},

			NumberOfTicks:    5,
			TickLength:       time.Second * 3,
			BonusCoefficient: dotCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dotDamage, isRollover)
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
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				// UA, Immo, Shadowflame exclusivity
				if hasUnstableAffliction && warlock.UnstableAffliction.Dot(target).IsActive() {
					warlock.UnstableAffliction.Dot(target).Deactivate(sim)
				}
				if hasShadowflameRune && warlock.Shadowflame.Dot(target).IsActive() {
					warlock.Shadowflame.Dot(target).Deactivate(sim)
				}

				spell.Dot(target).ApplyOrReset(sim)
			}
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

func (warlock *Warlock) getActiveImmolateSpell(target *core.Unit) *core.Spell {
	for _, immolateSpell := range warlock.Immolate {
		if immolateSpell.Dot(target).IsActive() {
			return immolateSpell
		}
	}
	return nil
}

func (warlock *Warlock) registerImmolateSpell() {
	warlock.Immolate = make([]*core.Spell, 0)

	maxRank := core.TernaryInt(core.IncludeAQ, ImmolateRanks, ImmolateRanks-1)
	for rank := 1; rank <= maxRank; rank++ {
		config := warlock.getImmolateConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.Immolate = append(warlock.Immolate, warlock.GetOrRegisterSpell(config))
		}
	}
}
