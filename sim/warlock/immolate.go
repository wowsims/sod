package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getImmolateConfig(rank int) core.SpellConfig {
	directCoeff := [9]float64{0, .058, .125, .2, .2, .2, .2, .2, .2}[rank]
	dotCoeff := [9]float64{0, .037, .081, .13, .13, .13, .13, .13, .13}[rank]
	baseDamage := [9]float64{0, 11, 24, 53, 101, 148, 208, 258, 279}[rank]
	dotDamage := [9]float64{0, 20, 40, 90, 165, 255, 365, 485, 510}[rank] / 5
	spellId := [9]int32{0, 348, 707, 1094, 2941, 11665, 11667, 11668, 25309}[rank]
	manaCost := [9]float64{0, 25, 45, 90, 155, 220, 295, 370, 380}[rank]
	level := [9]int{0, 1, 10, 20, 30, 40, 50, 60, 60}[rank]

	hasInvocationRune := warlock.HasRune(proto.WarlockRune_RuneBeltInvocation)
	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)
	hasUnstableAffliction := warlock.HasRune(proto.WarlockRune_RuneBracerUnstableAffliction)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (2000 - 100*time.Duration(warlock.Talents.Bane)),
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
			},
			CastTime: func(spell *core.Spell) time.Duration {
				durationDecrease := time.Duration(0)
				if warlock.shadowSparkAura.IsActive() {
					durationDecrease = (spell.DefaultCast.CastTime / 2) * time.Duration(warlock.shadowSparkAura.GetStacks())
				}
				return spell.DefaultCast.CastTime - durationDecrease
			},
		},

		BonusCritRating: float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,

		CritDamageBonus: warlock.ruin(),

		DamageMultiplierAdditive: 1 + 0.02*float64(warlock.Talents.Emberstorm),
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolate-" + warlock.Label + strconv.Itoa(rank),
			},

			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = dotDamage + dotCoeff*dot.Spell.SpellDamage()

				if !isRollover {
					dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
				}
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
				if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
					result.Damage *= warlock.getLakeOfFireMultiplier()
					result.Threat *= warlock.getLakeOfFireMultiplier()
				}
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := (baseDamage + directCoeff*spell.SpellDamage()) * (1 + 0.05*float64(warlock.Talents.ImprovedImmolate))

			if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
				damage *= warlock.getLakeOfFireMultiplier()
			}

			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				if hasUnstableAffliction && warlock.UnstableAffliction.Dot(target).IsActive() {
					warlock.UnstableAffliction.Dot(target).Deactivate(sim)
				}
				if hasInvocationRune && spell.Dot(target).IsActive() {
					warlock.InvocationRefresh(sim, spell.Dot(target))
				}
				spell.Dot(target).Apply(sim)
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDamage := dotDamage + dotCoeff*spell.SpellDamage()
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
	maxRank := 8

	warlock.Immolate = make([]*core.Spell, 0)
	for i := 1; i <= maxRank; i++ {
		config := warlock.getImmolateConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.Immolate = append(warlock.Immolate, warlock.GetOrRegisterSpell(config))
		}
	}
}
