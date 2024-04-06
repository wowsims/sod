package priest

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ShadowWordPainRanks = 8

var ShadowWordPainSpellId = [ShadowWordPainRanks + 1]int32{0, 589, 594, 970, 992, 2767, 10892, 10893, 10894}
var ShadowWordPainBaseDamage = [ShadowWordPainRanks + 1]float64{0, 30, 66, 132, 234, 366, 510, 672, 852}
var ShadowWordPainSpellCoef = [ShadowWordPainRanks + 1]float64{0, 0.067, 0.104, 0.154, 0.167, 0.167, 0.167, 0.167, 0.167} // per tick
var ShadowWordPainManaCost = [ShadowWordPainRanks + 1]float64{0, 25, 50, 95, 155, 230, 305, 385, 470}
var ShadowWordPainLevel = [ShadowWordPainRanks + 1]int{0, 4, 10, 18, 26, 34, 42, 50, 58}

func (priest *Priest) registerShadowWordPainSpell() {
	priest.ShadowWordPain = make([]*core.Spell, ShadowWordPainRanks+1)

	for rank := 1; rank <= ShadowWordPainRanks; rank++ {
		config := priest.getShadowWordPainConfig(rank)

		if config.RequiredLevel <= int(priest.Level) {
			priest.ShadowWordPain[rank] = priest.GetOrRegisterSpell(config)
		}
	}
}

func (priest *Priest) getShadowWordPainConfig(rank int) core.SpellConfig {
	ticks := int32(6)

	spellId := ShadowWordPainSpellId[rank]
	baseDotDamage := ShadowWordPainBaseDamage[rank] / float64(ticks)
	spellCoeff := ShadowWordPainSpellCoef[rank]
	manaCost := ShadowWordPainManaCost[rank]
	level := ShadowWordPainLevel[rank]

	numHits := core.TernaryInt32(priest.HasRune(proto.PriestRune_RuneLegsSharedPain), 3, 1)

	hasDespairRune := priest.HasRune(proto.PriestRune_RuneBracersDespair)

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagPureDot,

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

		BonusCritRating: priest.forceOfWillCritRating(),
		BonusHitRating:  priest.shadowHitModifier(),

		CritDamageBonus: core.TernaryFloat64(hasDespairRune, 1, 0),

		DamageMultiplier: priest.forceOfWillDamageModifier() * priest.darknessDamageModifier(),
		ThreatMultiplier: priest.shadowThreatModifier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("Shadow Word: Pain (Rank %d)", rank),
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if priest.HasRune(proto.PriestRune_RuneChestTwistedFaith) {
						priest.MindBlastModifier = 1.5
						priest.MindFlayModifier = 1.5
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					priest.MindBlastModifier = 1
					priest.MindFlayModifier = 1
				},
			},

			NumberOfTicks:    ticks + (priest.Talents.ImprovedShadowWordPain),
			TickLength:       time.Second * 3,
			BonusCoefficient: spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if hasDespairRune {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickSnapshotCritCounted)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				result := spell.CalcAndDealOutcome(sim, curTarget, spell.OutcomeMagicHit)
				if result.Landed() {
					spell.SpellMetrics[target.UnitIndex].Hits--
					priest.AddShadowWeavingStack(sim, curTarget)
					spell.Dot(curTarget).Apply(sim)
					spell.DealOutcome(sim, result)
				}
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				return spell.CalcPeriodicDamage(sim, target, baseDotDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}
}
