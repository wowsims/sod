package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) applyDemonArmor() {
	spellID := map[int32]int32{
		25: 706,
		40: 11733,
		50: 11734,
		60: 11735,
	}[warlock.Level]

	armor := map[int32]float64{
		25: 210.0,
		40: 390.0,
		50: 480.0,
		60: 570.0,
	}[warlock.Level]

	shadowRes := map[int32]float64{
		25: 3.0,
		40: 9.0,
		50: 12.0,
		60: 15.0,
	}[warlock.Level]

	warlock.AddStat(stats.Armor, armor)
	warlock.AddStat(stats.ShadowResistance, shadowRes)

	warlock.GetOrRegisterAura(core.Aura{
		Label:    "Demon Armor",
		ActionID: core.ActionID{SpellID: spellID},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

// Surrounds the caster with fel energy, increasing spell damage and healing by 1 plus additional spell damage and healing equal to 50% of your Spirit.
// In addition, you regain 2% of your maximum health every 5 sec.
func (warlock *Warlock) applyFelArmor() {
	actionID := core.ActionID{SpellID: 403619}

	warlock.AddStat(stats.SpellPower, 60)
	warlock.AddStatDependency(stats.Spirit, stats.SpellPower, .50)

	healthMetrics := warlock.NewHealthMetrics(actionID)
	warlock.GetOrRegisterAura(core.Aura{
		Label:    "Fel Armor",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 5,
				Priority: core.ActionPriorityAuto,
				OnAction: func(sim *core.Simulation) {
					warlock.GainHealth(sim, warlock.MaxHealth()*.02, healthMetrics)
				},
			})
		},
	})
}
