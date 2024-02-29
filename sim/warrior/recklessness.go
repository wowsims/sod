package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warrior *Warrior) RegisterRecklessnessCD() {
	if warrior.Level < 50 {
		return
	}

	actionID := core.ActionID{SpellID: 1719}

	reckAura := warrior.RegisterAura(core.Aura{
		Label:    "Recklessness",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier *= 1.2
			warrior.AddStatDynamic(sim, stats.MeleeCrit, 100)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier /= 1.2
			warrior.AddStatDynamic(sim, stats.MeleeCrit, -100)

		},
	})

	Recklessness := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Cast: core.CastConfig{
			IgnoreHaste: true,
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 30,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if !warrior.StanceMatches(BerserkerStance) {
				warrior.BerserkerStance.Cast(sim, nil)
			}

			reckAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: Recklessness,
		Type:  core.CooldownTypeDPS,
	})
}
