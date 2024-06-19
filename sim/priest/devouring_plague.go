package priest

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const DevouringPlagueRanks = 6

var DevouringPlagueSpellId = [DevouringPlagueRanks + 1]int32{0, 2944, 19276, 19277, 19278, 19279, 19280}
var DevouringPlagueBaseDamage = [DevouringPlagueRanks + 1]float64{0, 152, 272, 400, 544, 712, 904}
var DevouringPlagueManaCost = [DevouringPlagueRanks + 1]float64{0, 215, 350, 495, 645, 810, 985}
var DevouringPlagueLevel = [DevouringPlagueRanks + 1]int{0, 20, 28, 36, 44, 52, 60}

func (priest *Priest) registerDevouringPlagueSpell() {
	if priest.Race != proto.Race_RaceUndead {
		return
	}

	priest.DevouringPlague = make([]*core.Spell, DevouringPlagueRanks+1)
	cdTimer := priest.NewTimer()

	for rank := 1; rank < DevouringPlagueRanks; rank++ {
		config := priest.getDevouringPlagueConfig(rank, cdTimer)

		if config.RequiredLevel <= int(priest.Level) {
			priest.DevouringPlague[rank] = priest.GetOrRegisterSpell(config)
			priest.AddMajorCooldown(core.MajorCooldown{
				Spell:    priest.DevouringPlague[rank],
				Priority: int32(rank),
				Type:     core.CooldownTypeDPS,
			})
		}
	}
}

func (priest *Priest) getDevouringPlagueConfig(rank int, cdTimer *core.Timer) core.SpellConfig {
	var ticks int32 = 8

	spellId := DevouringPlagueSpellId[rank]
	baseDotDamage := (DevouringPlagueBaseDamage[rank] / float64(ticks)) * priest.darknessDamageModifier()
	manaCost := DevouringPlagueManaCost[rank]
	level := DevouringPlagueLevel[rank]

	spellCoeff := 0.063

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagPriest | core.SpellFlagAPL | core.SpellFlagDisease | core.SpellFlagPureDot,

		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: time.Minute * 3,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("Devouring Plague (Rank %d)", rank),
			},

			NumberOfTicks:    ticks,
			TickLength:       time.Second * 3,
			BonusCoefficient: spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				priest.AddShadowWeavingStack(sim, target)
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
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
