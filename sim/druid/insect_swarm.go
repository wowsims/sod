package druid

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
)

const InsectSwarmRanks = 5

var InsectSwarmSpellId = [InsectSwarmRanks + 1]int32{0, 5570, 24974, 24975, 24976, 24977}
var InsectSwarmBaseDamage = [InsectSwarmRanks + 1]float64{0, 66, 138, 174, 264, 324}
var InsectSwarmManaCost = [InsectSwarmRanks + 1]float64{0, 45, 85, 100, 140, 160}
var InsectSwarmLevel = [InsectSwarmRanks + 1]int{0, 20, 30, 40, 50, 60}

func (druid *Druid) registerInsectSwarmSpell() {
	druid.InsectSwarm = make([]*DruidSpell, InsectSwarmRanks+1)

	for rank := 1; rank <= InsectSwarmRanks; rank++ {
		level := InsectSwarmLevel[rank]
		if int32(level) <= druid.Level {
			numTicks := int32(6)
			tickLength := time.Second * 2

			spellID := InsectSwarmSpellId[rank]
			baseDamage := InsectSwarmBaseDamage[rank] / float64(numTicks)
			manaCost := InsectSwarmManaCost[rank]
			spellCoef := .158

			druid.InsectSwarmAuras = druid.NewEnemyAuraArray(core.InsectSwarmAura)

			druid.InsectSwarm[rank] = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
				SpellCode:   SpellCode_DruidInsectSwarm,
				ActionID:    core.ActionID{SpellID: spellID},
				SpellSchool: core.SpellSchoolNature,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskSpellDamage,
				Flags:       SpellFlagOmen | core.SpellFlagAPL,

				ManaCost: core.ManaCostOptions{
					FlatCost: manaCost,
				},
				Cast: core.CastConfig{
					DefaultCast: core.Cast{
						GCD: core.GCDDefault,
					},
				},

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label: fmt.Sprintf("Insect Swarm (Rank %d)", rank),
						OnGain: func(aura *core.Aura, sim *core.Simulation) {
							druid.InsectSwarmAuras.Get(aura.Unit).Activate(sim)
						},
						OnExpire: func(aura *core.Aura, sim *core.Simulation) {
							druid.InsectSwarmAuras.Get(aura.Unit).Deactivate(sim)
						},
					},

					NumberOfTicks:    numTicks,
					TickLength:       tickLength,
					BonusCoefficient: spellCoef,

					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
						dot.Snapshot(target, baseDamage, isRollover)
						if !druid.form.Matches(Moonkin) {
							dot.SnapshotCritChance = 0
						}
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickSnapshotCritCounted)
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
					if result.Landed() {
						spell.Dot(target).Apply(sim)
						spell.SpellMetrics[result.Target.UnitIndex].Hits--
					}
					spell.DealOutcome(sim, result)
				},

				RelatedAuras: []core.AuraArray{druid.InsectSwarmAuras},
			})
		}
	}
}
