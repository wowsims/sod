package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerConsecration() {
	if !paladin.Talents.Consecration {
		return
	}

	ranks := []struct {
		level    int32
		spellID  int32
		manaCost float64
		damage   float64
	}{
		{level: 20, spellID: 26573, manaCost: 135, damage: 64 / 8},
		{level: 30, spellID: 20116, manaCost: 235, damage: 120 / 8},
		{level: 40, spellID: 20922, manaCost: 320, damage: 192 / 8},
		{level: 50, spellID: 20923, manaCost: 435, damage: 280 / 8},
		{level: 60, spellID: 20924, manaCost: 565, damage: 384 / 8},
	}

	cd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 8,
	}

	hasWrath := paladin.hasRune(proto.PaladinRune_RuneHeadWrath)
	hasHallowedGround := paladin.hasRune(proto.PaladinRune_RuneChestHallowedGround)
	healthMetrics := paladin.NewHealthMetrics(core.ActionID{SpellID: 458286})

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
			Flags:       core.SpellFlagPureDot | core.SpellFlagAPL | core.SpellFlagBatchStartAttackMacro,

			RequiredLevel: int(rank.level),
			Rank:          i + 1,

			ClassSpellMask: ClassSpellMask_PaladinConsecration,
			ManaCost: core.ManaCostOptions{
				FlatCost: rank.manaCost,
			},
			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				CD: cd,
			},
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.042,
			Dot: core.DotConfig{
				IsAOE: true,
				Aura: core.Aura{
					Label: "Consecration" + paladin.Label + strconv.Itoa(i+1),
				},
				NumberOfTicks: 8,
				TickLength:    time.Second * 1,

				BonusCoefficient: 0.042,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, rank.damage, isRollover)
					if hasWrath {
						dot.Spell.BonusCritRating += paladin.GetStat(stats.MeleeCrit)
						dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
						dot.Spell.BonusCritRating -= paladin.GetStat(stats.MeleeCrit)
					}
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					// Consecration can miss, showing up as either a resist in logs or a
					// silent failure (missing damage tick).
					outcomeApplier := core.Ternary(hasWrath, dot.OutcomeMagicHitAndSnapshotCrit, dot.Spell.OutcomeMagicHit)
					for _, aoeTarget := range sim.Encounter.TargetUnits {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, outcomeApplier)
					}

					if hasHallowedGround {
						rawDamageResult := dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
						healAmount := rawDamageResult.RawDamage() * 2.0
						if sim.RandomFloat("Hallowed Ground") < dot.Spell.HealingCritChance() {
							attackTable := paladin.AttackTables[paladin.UnitIndex][dot.Spell.CastType]
							healAmount *= dot.Spell.CritMultiplier(attackTable)
						}

						paladin.GainHealth(sim, healAmount, healthMetrics)
					}
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.AOEDot().Apply(sim)
			},
		})
	}
}
