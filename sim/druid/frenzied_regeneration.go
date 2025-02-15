package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
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
			break
		}
	}

	if druid.FrenziedRegeneration == nil {
		return
	}

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.FrenziedRegeneration.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}

func (druid *Druid) newFrenziedRegenSpellConfig(frenziedRegenRank FrenziedRegenerationRankInfo) core.SpellConfig {
	hasImprovedFrenziedRegen := druid.HasRune(proto.DruidRune_RuneBracersImpFrenziedRegen)

	actionID := core.ActionID{SpellID: frenziedRegenRank.id}
	healthMetrics := druid.NewHealthMetrics(actionID)
	rageMetrics := druid.NewRageMetrics(actionID)

	var frenziedRegenPA *core.PendingAction
	druid.FrenziedRegenerationAura = druid.RegisterAura(core.Aura{
		Label:    "Frenzied Regeneration",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			frenziedRegenPA = core.NewPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 10,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					rageDumped := min(druid.CurrentRage(), 10.0)
					healthGained := core.TernaryFloat64(hasImprovedFrenziedRegen, druid.GetStat(stats.Health)*.1, rageDumped*frenziedRegenRank.rageConversion*druid.PseudoStats.HealingTakenMultiplier)

					if druid.CurrentRage() > druid.FrenziedRegenRageThreshold {
						druid.SpendRage(sim, rageDumped, rageMetrics)
					}
					druid.GainHealth(sim, healthGained, healthMetrics)
				},
			})
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			frenziedRegenPA = nil
		},
		OnRefresh: func(aura *core.Aura, sim *core.Simulation) {
			if frenziedRegenPA != nil {
				frenziedRegenPA.Cancel(sim)
			}

			frenziedRegenPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 10,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					rageDumped := min(druid.CurrentRage(), 10.0)
					healthGained := core.TernaryFloat64(hasImprovedFrenziedRegen, druid.GetStat(stats.Health)*.1, rageDumped*frenziedRegenRank.rageConversion*druid.PseudoStats.HealingTakenMultiplier)

					if druid.CurrentRage() > druid.FrenziedRegenRageThreshold {
						druid.SpendRage(sim, rageDumped, rageMetrics)
					}
					druid.GainHealth(sim, healthGained, healthMetrics)
				},
			})
		},
	})

	return core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_DruidFrenziedRegeneration,
		Flags:          core.SpellFlagHelpful | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			druid.FrenziedRegenerationAura.Activate(sim)
		},
	}
}
