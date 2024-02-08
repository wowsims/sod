package druid

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
)

const InsectSwarmRanks = 5

var InsectSwarmSpellId = [InsectSwarmRanks + 1]int32{0, 5570, 24974, 24975, 24976, 24977}
var InsectSwarmsSellDotCoeff = [InsectSwarmRanks + 1]float64{0, .158, .158, .158, .158, .158}
var InsectSwarmBaseDotDamage = [InsectSwarmRanks + 1]float64{0, 66, 138, 174, 264, 324}
var InsectSwarmManaCost = [InsectSwarmRanks + 1]float64{0, 45, 85, 100, 140, 160}
var InsectSwarmLevel = [InsectSwarmRanks + 1]int{0, 1, 30, 40, 50, 60}

func (druid *Druid) registerInsectSwarmSpell() {
	druid.InsectSwarm = make([]*DruidSpell, InsectSwarmRanks+1)

	for rank := 1; rank <= InsectSwarmRanks; rank++ {
		config := druid.getInsectSwarmBaseConfig(rank)

		if config.RequiredLevel <= int(druid.Level) {
			druid.InsectSwarm[rank] = druid.RegisterSpell(Humanoid|Moonkin, config)
		}
	}
}

func (druid *Druid) getInsectSwarmBaseConfig(rank int) core.SpellConfig {
	spellId := InsectSwarmSpellId[rank]
	spellDotCoeff := InsectSwarmsSellDotCoeff[rank]
	baseDotDamage := InsectSwarmBaseDotDamage[rank]
	manaCost := InsectSwarmManaCost[rank]
	level := InsectSwarmLevel[rank]

	ticks := int32(6)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolNature,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 0,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    fmt.Sprintf("InsectSwarm (Rank %d)", rank),
				ActionID: core.ActionID{SpellID: spellId},
			},
			NumberOfTicks: ticks,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = (baseDotDamage / float64(ticks)) + spellDotCoeff*dot.Spell.SpellDamage()
				dot.SnapshotAttackerMultiplier = 1 // dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		BonusCritRating:  core.SpellCritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)

			if result.Landed() {
				spell.Dot(target).Apply(sim)

				missAuras := core.InsectSwarmAura(target, spell.SpellID)
				missAuras.Activate(sim)
				spell.RelatedAuras = []core.AuraArray{{missAuras}}
			}

			spell.DealOutcome(sim, result)
		},
	}
}
