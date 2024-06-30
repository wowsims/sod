package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const DrainLifeRanks = 6

func (warlock *Warlock) getDrainLifeBaseConfig(rank int) core.SpellConfig {
	hasMasterChannelerRune := warlock.HasRune(proto.WarlockRune_RuneChestMasterChanneler)
	hasSoulSiphonRune := warlock.HasRune(proto.WarlockRune_RuneChestSoulSiphon)

	numTicks := core.TernaryInt32(hasMasterChannelerRune, 15, 5)

	spellId := [DrainLifeRanks + 1]int32{0, 689, 699, 709, 7651, 11699, 11700}[rank]
	spellCoeff := [DrainLifeRanks + 1]float64{0, .078, .1, .1, .1, .1, .1}[rank]
	baseDamage := [DrainLifeRanks + 1]float64{0, 10, 17, 29, 41, 55, 71}[rank]
	manaCost := [DrainLifeRanks + 1]float64{0, 55, 85, 135, 185, 240, 300}[rank]
	level := [DrainLifeRanks + 1]int{0, 14, 22, 30, 38, 46, 54}[rank]

	if hasMasterChannelerRune {
		manaCost *= 2
	}

	baseDamage *= 1 + warlock.shadowMasteryBonus() + 0.02*float64(warlock.Talents.ImprovedDrainLife)

	actionID := core.ActionID{SpellID: spellId}
	healthMetrics := warlock.NewHealthMetrics(actionID)

	spellConfig := core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		SpellCode:   SpellCode_WarlockDrainLife,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       WarlockFlagHaunt | core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagAffliction,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DrainLife-" + warlock.Label + strconv.Itoa(rank),
			},
			NumberOfTicks:    numTicks,
			TickLength:       1 * time.Second,
			BonusCoefficient: spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)

				if hasSoulSiphonRune {
					dot.SnapshotAttackerMultiplier *= warlock.calcSoulSiphonMultiplier(target, false)
				}

				// Drain Life heals so it snapshots target modifiers
				// Update 2024-06-29: It no longer snapshots on PTR
				// dot.SnapshotAttackerMultiplier *= dot.Spell.TargetDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType], true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)

				health := result.Damage
				if hasMasterChannelerRune {
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

	if hasMasterChannelerRune {
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
	warlock.DrainLife = make([]*core.Spell, 0)
	for rank := 1; rank <= DrainLifeRanks; rank++ {
		config := warlock.getDrainLifeBaseConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.DrainLife = append(warlock.DrainLife, warlock.GetOrRegisterSpell(config))
		}
	}
}
