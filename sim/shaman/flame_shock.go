package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const FlameShockRanks = 6

var FlameShockSpellId = [FlameShockRanks + 1]int32{0, 8050, 8052, 8053, 10447, 10448, 29228}
var FlameShockBaseDamage = [FlameShockRanks + 1]float64{0, 25, 51, 95, 164, 245, 292}
var FlameShockBaseDotDamage = [FlameShockRanks + 1]float64{0, 28, 48, 96, 168, 256, 320}
var FlameShockBaseSpellCoef = [FlameShockRanks + 1]float64{0, .134, .198, .214, .214, .214, .214}
var FlameShockDotSpellCoef = [FlameShockRanks + 1]float64{0, .063, .093, .1, .1, .1, .1}
var FlameShockManaCost = [FlameShockRanks + 1]float64{0, 55, 95, 160, 250, 345, 410}
var FlameShockLevel = [FlameShockRanks + 1]int{0, 10, 18, 28, 40, 52, 60}

func (shaman *Shaman) registerFlameShockSpell(shockTimer *core.Timer) {
	shaman.FlameShock = make([]*core.Spell, FlameShockRanks+1)

	for rank := 1; rank <= FlameShockRanks; rank++ {
		if FlameShockLevel[rank] <= int(shaman.Level) {
			shaman.FlameShock[rank] = shaman.RegisterSpell(shaman.newFlameShockSpell(rank, shockTimer))
		}
	}
}

func (shaman *Shaman) newFlameShockSpell(rank int, shockTimer *core.Timer) core.SpellConfig {
	hasBurnRune := shaman.HasRune(proto.ShamanRune_RuneHelmBurn)

	numTicks := 4 + core.TernaryInt32(hasBurnRune, BurnFlameShockBonusTicks, 0)
	tickDuration := time.Second * 3

	spellId := FlameShockSpellId[rank]
	baseDamage := FlameShockBaseDamage[rank]
	baseDotDamage := FlameShockBaseDotDamage[rank] / float64(4)
	baseSpellCoeff := FlameShockBaseSpellCoef[rank]
	dotSpellCoeff := FlameShockDotSpellCoef[rank]
	manaCost := FlameShockManaCost[rank]
	level := FlameShockLevel[rank]

	spell := shaman.newShockSpellConfig(
		core.ActionID{SpellID: spellId},
		core.SpellSchoolFire,
		manaCost,
		shockTimer,
	)

	spell.ClassSpellMask = ClassSpellMask_ShamanFlameShock
	spell.RequiredLevel = level
	spell.Rank = rank

	spell.BonusCoefficient = baseSpellCoeff

	spell.Dot = core.DotConfig{
		Aura: core.Aura{
			Label: fmt.Sprintf("Flame Shock (Rank %d)", rank),
		},

		NumberOfTicks:    int32(numTicks),
		TickLength:       tickDuration,
		BonusCoefficient: dotSpellCoeff,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.Snapshot(target, baseDotDamage, isRollover)
		},

		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
		},
	}

	results := make([]*core.SpellResult, min(core.TernaryInt32(hasBurnRune, BurnFlameShockTargetCount, 1), shaman.Env.GetNumTargets()))
	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		for idx := range results {
			results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			target = sim.Environment.NextTargetUnit(target)
		}

		for _, result := range results {
			spell.DealDamage(sim, result)
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)

				if shaman.HasRune(proto.ShamanRune_RuneLegsAncestralGuidance) {
					shaman.lastFlameShockTarget = target
				}
			}
		}
	}

	return spell
}
