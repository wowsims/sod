package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (druid *Druid) registerEnrageSpell() {
	actionID := core.ActionID{SpellID: 5229}
	rageMetrics := druid.NewRageMetrics(actionID)

	instantRage := []float64{0, 5, 10}[druid.Talents.ImprovedEnrage]
	initarmor := druid.BaseEquipStats()[stats.Armor]
	hasCenarionRage4Piece := druid.HasSetBonus(ItemSetCenarionRage, 4)

	druid.EnrageAura = druid.RegisterAura(core.Aura{
		Label:    "Enrage Aura",
		ActionID: actionID,
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !hasCenarionRage4Piece {
				druid.AddStatDynamic(sim, stats.Armor, float64(0.16*initarmor)*-1)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if !hasCenarionRage4Piece {
				druid.AddStatDynamic(sim, stats.Armor, float64(0.16*initarmor))
			}
		},
	})

	druid.Enrage = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.AddRage(sim, instantRage, rageMetrics)

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 10,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					if druid.EnrageAura.IsActive() {
						druid.AddRage(sim, 2, rageMetrics)
					}
				},
			})

			druid.EnrageAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Enrage.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
