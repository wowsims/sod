package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (druid *Druid) applySurvivalInstincts() {
	if !druid.HasRune(proto.DruidRune_RuneFeetSurvivalInstincts) {
		return
	}

	actionID := core.ActionID{SpellID: 408024}
	healthMetrics := druid.NewHealthMetrics(actionID)

	var bonusHealth float64
	druid.SurvivalInstinctsAura = druid.RegisterAura(core.Aura{
		Label:    "Survival Instincts",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = druid.MaxHealth() * 0.2
			druid.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})
			druid.GainHealth(sim, bonusHealth, healthMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})
		},
	})

	survivalInsinctsSpell := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			druid.SurvivalInstinctsAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: survivalInsinctsSpell.Spell,
		Type:  core.CooldownTypeSurvival,
	})

	// In addition, you regenerate 5 rage every time you dodge while in Bear Form or Dire Bear Form, 10 energy while in Cat Form, or 1% of your maximum mana while in any other form.
	actionID = core.ActionID{SpellID: 409809}
	manaMetrics := druid.NewManaMetrics(actionID)
	rageMetrics := druid.NewRageMetrics(actionID)
	energyMetrics := druid.NewEnergyMetrics(actionID)
	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: "Survival Instincts - Passive",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidDodge() {
				switch druid.form {
				case Cat:
					druid.AddEnergy(sim, 10, energyMetrics)
				case Bear:
					druid.AddRage(sim, 5, rageMetrics)
				default:
					amount := druid.MaxMana() * 0.01
					druid.AddMana(sim, amount, manaMetrics)
				}
			}
		},
	}))
}
