package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// Recklessness now increases critical strike chance by 50% (was 100%) and the duration is reduced to 12 seconds, but the cooldown is reduced to 5 minutes.
func (warrior *Warrior) RegisterRecklessnessCD() {
	if warrior.Level < 50 {
		return
	}

	actionID := core.ActionID{SpellID: 1719}

	reckAura := warrior.RegisterAura(core.Aura{
		Label:    "Recklessness",
		ActionID: actionID,
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier *= 1.2
			warrior.AddStatDynamic(sim, stats.MeleeCrit, 50*core.CritRatingPerCritChance)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier /= 1.2
			warrior.AddStatDynamic(sim, stats.MeleeCrit, -50*core.CritRatingPerCritChance)

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
				Duration: time.Minute * 5,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance) || warrior.StanceMatches(GladiatorStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			reckAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: Recklessness,
		Type:  core.CooldownTypeDPS,
	})
}
