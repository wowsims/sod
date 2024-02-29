package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const HurricaneRanks = 3

var HurricaneSpellId = [HurricaneRanks + 1]int32{0, 16914, 17401, 17402}
var HurricaneBaseDamage = [HurricaneRanks + 1]float64{0, 72, 102, 134}
var HurricaneSpellCoef = [HurricaneRanks + 1]float64{0, .03, .03, .03}
var HurricaneManaCost = [HurricaneRanks + 1]float64{0, 880, 1180, 1495}
var HurricaneLevel = [HurricaneRanks + 1]int{0, 40, 50, 60}

func (druid *Druid) registerHurricaneSpell() {
	druid.Hurricane = make([]*DruidSpell, HurricaneRanks+1)
	druid.HurricaneTickSpell = make([]*DruidSpell, HurricaneRanks+1)

	for rank := 1; rank <= HurricaneRanks; rank++ {
		config := druid.newHurricaneSpellConfig(rank)

		if config.RequiredLevel <= int(druid.Level) {
			druid.Hurricane[rank] = druid.RegisterSpell(Humanoid|Moonkin, config)
			druid.HurricaneTickSpell[rank] = druid.RegisterSpell(Humanoid|Moonkin, druid.newHurricaneTickSpellConfig(rank))
		}
	}
}

func (druid *Druid) newHurricaneSpellConfig(rank int) core.SpellConfig {
	spellId := HurricaneSpellId[rank]
	manaCost := HurricaneManaCost[rank]
	level := HurricaneLevel[rank]

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagChanneled | core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Hurricane",
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				druid.HurricaneTickSpell[rank].Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			druid.AutoAttacks.CancelAutoSwing(sim)
			spell.AOEDot().Apply(sim)
		},
	}
}

func (druid *Druid) newHurricaneTickSpellConfig(rank int) core.SpellConfig {
	spellId := HurricaneSpellId[rank]
	baseDamage := HurricaneBaseDamage[rank]
	spellCoef := HurricaneSpellCoef[rank]

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskProc,

		CritMultiplier:   1,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := baseDamage + spellCoef*spell.SpellDamage()
			// damage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
			}
		},
	}
}
