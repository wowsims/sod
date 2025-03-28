package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core/stats"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (paladin *Paladin) registerExorcism() {
	ranks := []struct {
		level      int32
		spellID    int32
		manaCost   float64
		scaleLevel int32
		minDamage  float64
		maxDamage  float64
		scale      float64
	}{
		{level: 20, spellID: 415068, manaCost: 85, scaleLevel: 25, minDamage: 84, maxDamage: 96, scale: 1.2},
		{level: 28, spellID: 415069, manaCost: 135, scaleLevel: 33, minDamage: 152, maxDamage: 172, scale: 1.6},
		{level: 36, spellID: 415070, manaCost: 180, scaleLevel: 41, minDamage: 217, maxDamage: 245, scale: 2.0},
		{level: 44, spellID: 415071, manaCost: 235, scaleLevel: 49, minDamage: 304, maxDamage: 342, scale: 2.4},
		{level: 52, spellID: 415072, manaCost: 285, scaleLevel: 57, minDamage: 393, maxDamage: 439, scale: 2.8},
		{level: 60, spellID: 415073, manaCost: 345, scaleLevel: 60, minDamage: 505, maxDamage: 563, scale: 3.2},
	}

	hasWrath := paladin.hasRune(proto.PaladinRune_RuneHeadWrath)

	paladin.exorcismCooldown = &core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 15,
	}

	for i, rank := range ranks {
		rank := rank
		if paladin.Level < rank.level {
			break
		}

		minDamage := rank.minDamage + float64(min(paladin.Level, rank.scaleLevel)-rank.level)*rank.scale
		maxDamage := rank.maxDamage + float64(min(paladin.Level, rank.scaleLevel)-rank.level)*rank.scale

		spell := paladin.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: rank.spellID},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagBinary | core.SpellFlagBatchStartAttackMacro, //Logs show it never has partial resists, No clue why, still misses

			RequiredLevel: int(rank.level),
			Rank:          i + 1,

			ClassSpellMask: ClassSpellMask_PaladinExorcism,

			ManaCost: core.ManaCostOptions{
				FlatCost:   rank.manaCost,
				Multiplier: core.TernaryInt32(paladin.hasRune(proto.PaladinRune_RuneFeetTheArtOfWar), 20, 100),
			},

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				CD: *paladin.exorcismCooldown,
			},

			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				if paladin.artOfWarDelayAura == nil {
					return true
				}

				if paladin.artOfWarDelayAura.IsActive() {
					if sim.Log != nil {
						paladin.Log(sim, "The Art of War Delay prevents Exorcism from being cast for another %s", paladin.artOfWarDelayAura.RemainingDuration(sim))
					}
					return false
				}

				return true
			},

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BonusCoefficient: 0.429,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				bonusCrit := 0.0
				if target.MobType == proto.MobType_MobTypeDemon || target.MobType == proto.MobType_MobTypeUndead {
					bonusCrit += 100 * core.CritRatingPerCritChance
				}
				if hasWrath {
					bonusCrit += paladin.GetStat(stats.MeleeCrit)
				}

				spell.BonusCritRating += bonusCrit
				spell.CalcAndDealDamage(sim, target, sim.Roll(minDamage, maxDamage), spell.OutcomeMagicHitAndCrit)
				spell.BonusCritRating -= bonusCrit
			},
		})

		paladin.exorcism = append(paladin.exorcism, spell)
	}
}
