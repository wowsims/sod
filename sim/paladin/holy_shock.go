package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core/stats"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (paladin *Paladin) registerHolyShock() {

	hasInfusionOfLight := paladin.hasRune(proto.PaladinRune_RuneWaistInfusionOfLight)

	cdTime := time.Second * 30
	if hasInfusionOfLight {
		cdTime = time.Second * 6
	}

	paladin.holyShockCooldown = &core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: cdTime,
	}

	if !paladin.Talents.HolyShock {
		return
	}

	ranks := []struct {
		level     int32
		spellID   int32
		manaCost  float64
		minDamage float64
		maxDamage float64
	}{
		{level: 40, spellID: 20473, manaCost: 225, minDamage: 204, maxDamage: 220},
		{level: 48, spellID: 20929, manaCost: 275, minDamage: 279, maxDamage: 301},
		{level: 56, spellID: 20930, manaCost: 325, minDamage: 365, maxDamage: 395},
	}

	damageMultiplier := core.TernaryFloat64(hasInfusionOfLight, 1.5, 1.0)

	//hasArtOfWar := paladin.hasRune(proto.PaladinRune_RuneFeetTheArtOfWar)
	manaCostMultiplier := int32(100) //core.TernaryFloat64(hasArtOfWar, 0.2, 1.0)

	hasWrath := paladin.hasRune(proto.PaladinRune_RuneHeadWrath)

	manaMetrics := paladin.NewManaMetrics(core.ActionID{SpellID: 437063}) // Infusion of Light mana restore

	for i, rank := range ranks {
		rank := rank
		if paladin.Level < rank.level {
			break
		}

		paladin.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: rank.spellID},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

			RequiredLevel: int(rank.level),
			Rank:          i + 1,

			SpellCode: SpellCode_PaladinHolyShock,

			ManaCost: core.ManaCostOptions{
				FlatCost:   rank.manaCost,
				Multiplier: manaCostMultiplier,
			},

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				CD: *paladin.holyShockCooldown,
			},

			DamageMultiplier: damageMultiplier,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.429,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(rank.minDamage, rank.maxDamage)

				bonusCrit := core.TernaryFloat64(hasWrath, paladin.GetStat(stats.MeleeCrit), 0)
				spell.BonusCritRating += bonusCrit
				result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				spell.BonusCritRating -= bonusCrit

				// If we crit, Infusion of Light refunds base mana cost and reduces next Holy Shock Cooldown by 3 seconds
				if hasInfusionOfLight && result.Outcome.Matches(core.OutcomeCrit) {
					paladin.AddMana(sim, rank.manaCost, manaMetrics)
					paladin.holyShockCooldown.Set(sim.CurrentTime + max(0, paladin.holyShockCooldown.TimeToReady(sim)-(time.Second*3)))

				}
			},
		})
	}
}
