package priest

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const MindSearTicks = 5

func (priest *Priest) registerMindSearSpell() {
	if !priest.HasRune(proto.PriestRune_RuneHandsMindSear) {
		return
	}

	priest.MindSear = make([]*core.Spell, MindSearTicks)

	var tick int32
	for tick = 0; tick < MindSearTicks; tick++ {
		priest.MindSear[tick] = priest.RegisterSpell(priest.newMindSearSpellConfig(tick))
	}
}

func (priest *Priest) newMindSearSpellConfig(tickIdx int32) core.SpellConfig {
	spellId := int32(proto.PriestRune_RuneHandsMindSear)
	manaCost := .28

	numTicks := tickIdx
	flags := SpellFlagPriest | core.SpellFlagChanneled | core.SpellFlagNoMetrics
	if tickIdx == 0 {
		numTicks = 5
		flags |= core.SpellFlagAPL
	}
	tickLength := time.Second

	mindSearTickSpell := priest.newMindSearTickSpell(tickIdx)

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId}.WithTag(tickIdx),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       flags,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: fmt.Sprintf("MindSear-%d", tickIdx),
			},
			NumberOfTicks: numTicks,
			TickLength:    tickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					mindSearTickSpell.Cast(sim, aoeTarget)
					mindSearTickSpell.SpellMetrics[target.UnitIndex].Casts -= 1
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			mindSearTickSpell.SpellMetrics[target.UnitIndex].Casts += 1

			if result.Landed() {
				spell.AOEDot().Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	}
}

func (priest *Priest) newMindSearTickSpell(numTicks int32) *core.Spell {
	baseDamageLow := priest.baseRuneAbilityDamage() * .7 * priest.darknessDamageModifier()
	baseDamageHigh := priest.baseRuneAbilityDamage() * .78 * priest.darknessDamageModifier()
	spellCoeff := 0.15 // classic penalty for mf having a slow effect

	return priest.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 413260}.WithTag(numTicks),
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskProc,

		CritDamageBonus: priest.periodicCritBonus(),
		BonusHitRating:  1, // Not an independent hit once initial lands

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamageLow, baseDamageHigh)
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim, target)
			}
		},
	})
}
