package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerHolyWrath() {
	ranks := []struct {
		level      int32
		spellID    int32
		manaCost   float64
		scaleLevel int32
		minDamage  float64
		maxDamage  float64
		scale      float64
	}{
		{level: 50, spellID: 2812, manaCost: 645, scaleLevel: 54, minDamage: 362, maxDamage: 428, scale: 1.6},
		{level: 60, spellID: 10318, manaCost: 805, scaleLevel: 60, minDamage: 490, maxDamage: 576, scale: 1.9},
	}

	hasPurifyingPower := paladin.hasRune(proto.PaladinRune_RuneWristPurifyingPower)
	hasWrath := paladin.hasRune(proto.PaladinRune_RuneHeadWrath)

	var results []*core.SpellResult

	for i, rank := range ranks {
		rank := rank
		if paladin.Level < rank.level {
			break
		}

		minDamage := rank.minDamage + float64(min(paladin.Level, rank.scaleLevel)-rank.level)*rank.scale
		maxDamage := rank.maxDamage + float64(min(paladin.Level, rank.scaleLevel)-rank.level)*rank.scale

		holyWrathSpell := paladin.GetOrRegisterSpell(core.SpellConfig{
			ClassSpellMask: ClassSpellMask_PaladinHolyWrath,
			ActionID:       core.ActionID{SpellID: rank.spellID},
			SpellSchool:    core.SpellSchoolHoly,
			DefenseType:    core.DefenseTypeMagic,
			ProcMask:       core.ProcMaskSpellDamage, // TODO to be tested
			Flags:          core.SpellFlagAPL | core.SpellFlagBatchStartAttackMacro,
			MissileSpeed:   20,

			RequiredLevel: int(rank.level),
			Rank:          i + 1,

			ManaCost: core.ManaCostOptions{
				FlatCost: rank.manaCost,
			},
			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD:      core.GCDDefault,
					CastTime: time.Second * 2,
				},

				CD: core.Cooldown{
					Timer:    paladin.NewTimer(),
					Duration: time.Second * 60,
				},
			},

			DamageMultiplier: 1.0,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.19,

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				bonusCrit := core.TernaryFloat64(hasWrath, paladin.GetStat(stats.MeleeCrit), 0)
				spell.BonusCritRating += bonusCrit

				results = results[:0]
				for _, target := range paladin.Env.Encounter.TargetUnits {
					if hasPurifyingPower || (target.MobType == proto.MobType_MobTypeDemon || target.MobType == proto.MobType_MobTypeUndead) {
						damage := sim.Roll(minDamage, maxDamage)
						result := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
						results = append(results, result)
					}
				}

				spell.BonusCritRating -= bonusCrit

				if len(results) == 0 {
					return
				}

				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					for _, result := range results {
						spell.DealDamage(sim, result)
					}
				})
			},
		})

		paladin.holyWrath = append(paladin.holyWrath, holyWrathSpell)
	}
}
