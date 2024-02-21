package priest

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
)

const MindFlayRanks = 6
const MindFlayTicks = 3

var MindFlaySpellId = [MindFlayRanks + 1]int32{0, 15407, 17311, 17312, 17313, 17314, 18807}
var MindFlayTickSpellId = [MindFlayRanks + 1]int32{0, 16568, 7378, 17316, 17317, 17318, 18808}
var MindFlayBaseDamage = [MindFlayRanks + 1]float64{0, 75, 126, 186, 261, 330, 426}
var MindFlayManaCost = [MindFlayRanks + 1]float64{0, 45, 70, 100, 135, 165, 205}
var MindFlayLevel = [MindFlayRanks + 1]int{0, 20, 28, 36, 44, 52, 60}

func (priest *Priest) registerMindFlay() {
	if !priest.Talents.MindFlay {
		return
	}

	priest.MindFlay = make([][]*core.Spell, MindFlayRanks+1)

	for rank := 1; rank < MindFlayRanks; rank++ {
		priest.MindFlay[rank] = make([]*core.Spell, MindFlayTicks+1)

		var tick int32
		for tick = 0; tick < MindFlayTicks; tick++ {
			config := priest.newMindFlaySpellConfig(rank, tick)

			if config.RequiredLevel <= int(priest.Level) {
				priest.MindFlay[rank][tick] = priest.GetOrRegisterSpell(config)
			}
		}
	}
}

func (priest *Priest) newMindFlaySpellConfig(rank int, tickIdx int32) core.SpellConfig {
	spellId := MindFlaySpellId[rank]
	baseDamage := MindFlayBaseDamage[rank]
	manaCost := MindFlayManaCost[rank]
	level := MindFlayLevel[rank]

	spellCoeff := 0.15 // classic penalty for mf having a slow effect

	numTicks := tickIdx
	flags := core.SpellFlagNoMetrics | core.SpellFlagChanneled
	if tickIdx == 0 {
		numTicks = 3
		flags |= core.SpellFlagAPL
	}
	tickLength := time.Second
	mindFlayTickSpell := priest.newMindFlagTickSpell(rank, tickIdx)

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       flags,

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

		BonusHitRating:   priest.shadowHitModifier(),
		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("MindFlay-%d-%d", rank, tickIdx),
			},
			NumberOfTicks:       numTicks,
			TickLength:          tickLength,
			AffectedByCastSpeed: false,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mindFlayTickSpell.Cast(sim, target)
				mindFlayTickSpell.SpellMetrics[target.UnitIndex].Casts -= 1
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			mindFlayTickSpell.SpellMetrics[target.UnitIndex].Casts += 1

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := baseDamage/MindFlayTicks + (spellCoeff * spell.SpellDamage())
			return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
		},
	}
}

func (priest *Priest) newMindFlagTickSpell(rank int, numTicks int32) *core.Spell {
	spellId := MindFlayTickSpellId[rank]
	baseDamage := MindFlayBaseDamage[rank]
	spellCoeff := 0.15 // classic penalty for mf having a slow effect

	return priest.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId}.WithTag(numTicks),
		Rank:        rank,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskProc | core.ProcMaskNotInSpellbook,

		BonusHitRating:   1, // Not an independent hit once initial lands
		BonusCritRating:  0,
		DamageMultiplier: priest.forceOfWillDamageModifier(),
		CritMultiplier:   1.0,
		ThreatMultiplier: priest.shadowThreatModifier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := baseDamage/MindFlayTicks + (spellCoeff*spell.SpellDamage())*priest.MindFlayModifier
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeExpectedMagicAlwaysHit)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim, target)
			}
		},
	})
}
