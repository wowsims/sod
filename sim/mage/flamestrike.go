package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
)

const FlamestrikeRanks = 6

var FlamestrikeSpellId = [FlamestrikeRanks + 1]int32{0, 2120, 2121, 8422, 8423, 10215, 10216}
var FlamestrikeBaseDamage = [FlamestrikeRanks + 1][]float64{{0}, {55, 71}, {96, 123}, {159, 197}, {220, 272}, {294, 362}, {381, 466}}
var FlamestrikeDotDamage = [FlamestrikeRanks + 1]float64{0, 48, 88, 140, 196, 264, 340}
var FlamestrikeSpellCoeff = [FlamestrikeRanks + 1]float64{0, .134, .157, .157, .157, .157, .157}
var FlamestrikeDotCoeff = [FlamestrikeRanks + 1]float64{0, .017, .02, .02, .02, .02, .02}
var FlamestrikeManaCost = [FlamestrikeRanks + 1]float64{0, 195, 330, 490, 650, 815, 990}
var FlamestrikeLevel = [FlamestrikeRanks + 1]int{0, 16, 24, 32, 40, 48, 56}

func (mage *Mage) registerFlamestrikeSpell() {
	mage.Flamestrike = make([]*core.Spell, FlamestrikeRanks+1)

	for rank := 1; rank <= FlamestrikeRanks; rank++ {
		config := mage.newFlamestrikeSpellConfig(rank)

		if config.RequiredLevel <= int(mage.Level) {
			mage.Flamestrike[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) newFlamestrikeSpellConfig(rank int) core.SpellConfig {
	numTicks := int32(4)
	tickLength := time.Second * 2

	spellId := FlamestrikeSpellId[rank]
	baseDamageLow := FlamestrikeBaseDamage[rank][0]
	baseDamageHigh := FlamestrikeBaseDamage[rank][1]
	baseDotDamage := FlamestrikeDotDamage[rank] / float64(numTicks)
	spellCoeff := FlamestrikeSpellCoeff[rank]
	dotCoeff := FlamestrikeDotCoeff[rank]
	manaCost := FlamestrikeManaCost[rank]
	level := FlamestrikeLevel[rank]

	castTime := time.Second * 3

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
		},

		BonusCritRating: float64(5 * mage.Talents.ImprovedFlamestrike * core.CritRatingPerCritChance),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: fmt.Sprintf("Flamestrike (Rank %d)", rank),
			},
			NumberOfTicks:    numTicks,
			TickLength:       tickLength,
			BonusCoefficient: dotCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicCrit)
			}
			spell.AOEDot().Apply(sim)
		},
	}
}
