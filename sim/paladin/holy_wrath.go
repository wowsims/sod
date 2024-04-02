package paladin

import (
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerHolyWrathSpell() {
	ranks := []struct {
		level      int32
		spellID    int32
		damageLow  float64
		damageHigh float64
		manaCost   float64
	}{
		{level: 50, spellID: 2812, damageLow: 362, damageHigh: 428, manaCost: 645}, // 368-435 at level >= 54
		{level: 60, spellID: 10318, damageLow: 490, damageHigh: 576, manaCost: 805},
	}

	cd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 60,
	}

	hasPurifyingPower := paladin.HasRune(proto.PaladinRune_RuneWristPurifyingPower)
	hasWrath := paladin.HasRune(proto.PaladinRune_RuneHeadWrath)

	if hasPurifyingPower {
		cd.Duration /= 2
	}

	var results []*core.SpellResult

	for i, rank := range ranks {
		if paladin.Level < rank.level {
			break
		}

		paladin.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: rank.spellID},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage, // TODO to be tested
			Flags:       core.SpellFlagAPL,

			Rank:          i + 1,
			RequiredLevel: int(rank.level),

			ManaCost: core.ManaCostOptions{
				FlatCost: rank.manaCost,
			},
			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD:      time.Second,
					CastTime: time.Second * 2,
				},
				IgnoreHaste: true,
				CD:          cd,
			},

			BonusCritRating: paladin.holyPower() + paladin.fanaticism(),

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				bonusCrit := core.TernaryFloat64(hasWrath, paladin.GetStat(stats.MeleeCrit), 0)
				spell.BonusCritRating += bonusCrit

				results = results[:0]
				for _, target := range paladin.Env.Encounter.TargetUnits {
					if hasPurifyingPower || (target.MobType == proto.MobType_MobTypeDemon || target.MobType == proto.MobType_MobTypeUndead) {
						damage := sim.Roll(rank.damageLow, rank.damageHigh) + 0.19*spell.SpellDamage()
						result := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
						results = append(results, result)
					}
				}

				for _, result := range results {
					spell.DealDamage(sim, result)
				}

				spell.BonusCritRating -= bonusCrit
			},
		})
	}
}
