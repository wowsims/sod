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
	for tick = 1; tick < MindSearTicks; tick++ {
		priest.MindSear[tick] = priest.GetOrRegisterSpell(priest.newMindSearSpellConfig(tick))
	}
}

func (priest *Priest) newMindSearSpellConfig(tickIdx int32) core.SpellConfig {
	spellId := int32(proto.PriestRune_RuneHandsMindSear)
	manaCost := .28

	numTicks := tickIdx
	flags := core.SpellFlagChanneled | core.SpellFlagNoMetrics
	if tickIdx == 0 {
		numTicks = 5
		flags |= core.SpellFlagAPL
	}
	tickLength := time.Second

	mindSearTickSpell := priest.newMindSearTickSpell(tickIdx)

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
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

		BonusHitRating:   priest.shadowHitModifier(),
		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("MindSear-%d", tickIdx),
			},
			NumberOfTicks:       numTicks,
			TickLength:          tickLength,
			AffectedByCastSpeed: false,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					if aoeTarget != target {
						mindSearTickSpell.Cast(sim, aoeTarget)
						mindSearTickSpell.SpellMetrics[target.UnitIndex].Casts -= 1
					}
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			mindSearTickSpell.SpellMetrics[target.UnitIndex].Casts += 1

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	}
}

func (priest *Priest) newMindSearTickSpell(numTicks int32) *core.Spell {
	level := float64(priest.Level)
	spellId := int32(proto.PriestRune_RuneHandsMindSear)
	baseDamage := (9.456667 + 0.635108*level + 0.039063*level*level)
	baseDamageLow := baseDamage * .7
	baseDamageHigh := baseDamage * .78
	spellCoeff := 0.15 // classic penalty for mf having a slow effect

	return priest.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId}.WithTag(numTicks),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskProc | core.ProcMaskNotInSpellbook,

		BonusHitRating:   1, // Not an independent hit once initial lands
		BonusCritRating:  priest.forceOfWillCritRating(),
		DamageMultiplier: priest.forceOfWillDamageModifier(),
		CritMultiplier:   1,
		ThreatMultiplier: priest.shadowThreatModifier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamageLow, baseDamageHigh) + (spellCoeff * spell.SpellDamage())
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeExpectedMagicAlwaysHit)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim, target)
			}
		},
	})
}
