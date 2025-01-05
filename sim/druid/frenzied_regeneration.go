package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

type FrenziedRegenerationRankInfo struct {
	id             int32
	level          int32
	rageConversion float64
}

var frenziedRegenerationSpells = []FrenziedRegenerationRankInfo{
	{
		id:             22842,
		level:          36,
		rageConversion: 10.0,
	},
	{
		id:             22895,
		level:          46,
		rageConversion: 15.0,
	},
	{

		id:             22896,
		level:          56,
		rageConversion: 20.0,
	},
}

func (druid *Druid) registerFrenziedRegenerationCD() {
	// Add highest available rank for level.
	for rank := len(frenziedRegenerationSpells) - 1; rank >= 0; rank-- {
		if druid.Level >= frenziedRegenerationSpells[rank].level {
			config := druid.newFrenziedRegenSpellConfig(frenziedRegenerationSpells[rank])
			druid.FrenziedRegeneration = druid.RegisterSpell(Bear, config)
			healingMulti := 1.0

			druid.FrenziedRegenerationAura = druid.RegisterAura(core.Aura{
				Label:    "Frenzied Regeneration",
				ActionID: druid.FrenziedRegeneration.ActionID,
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					druid.PseudoStats.HealingTakenMultiplier *= healingMulti
				},

				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					druid.PseudoStats.HealingTakenMultiplier /= healingMulti
				},
			})
			break
		}
	}

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.FrenziedRegeneration.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}

func (druid *Druid) newFrenziedRegenSpellConfig(frenziedRegenRank FrenziedRegenerationRankInfo) core.SpellConfig {
	actionID := core.ActionID{SpellID: frenziedRegenRank.id}
	healthMetrics := druid.NewHealthMetrics(actionID)
	rageMetrics := druid.NewRageMetrics(actionID)
	hasImprovedFrenziedRegen := druid.HasRune(DruidRune_RuneBracersImprovedFrenziedRegeneration)

	cdTimer := druid.NewTimer()
	cd := time.Minute * 3

	return core.SpellConfig{
		ActionID: actionID,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 10,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					rageDumped := min(druid.CurrentRage(), 10.0)
					healthGained := core.TernaryFloat64(hasImprovedFrenziedRegen, druid.GetStat(stat.Health)*.1, rageDumped*frenziedRegenRank.rageConversion*druid.PseudoStats.HealingTakenMultiplier)

					if druid.FrenziedRegenerationAura.IsActive() {
						druid.SpendRage(sim, rageDumped, rageMetrics)
						druid.GainHealth(sim, healthGained, healthMetrics)
					}
				},
			})

			druid.FrenziedRegenerationAura.Activate(sim)
		},
	}
}
