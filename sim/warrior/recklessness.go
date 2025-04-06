package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	DefaultRecklessnessDamageTakenMultiplier = 1.20
	DefaultRecklessnessStance                = BerserkerStance
	DefaultRecklessnessDuration              = time.Second * 12
)

// Recklessness now increases critical strike chance by 50% (was 100%) and the duration is reduced to 12 seconds, but the cooldown is reduced to 5 minutes.
func (warrior *Warrior) RegisterRecklessnessCD(sharedTimer *core.Timer) {
	if warrior.Level < 50 {
		return
	}

	actionID := core.ActionID{SpellID: 1719}
	warrior.recklessnessDamageTakenMultiplier = DefaultRecklessnessDamageTakenMultiplier

	reckAura := warrior.RegisterAura(core.Aura{
		Label:    "Recklessness",
		ActionID: actionID,
		Duration: DefaultRecklessnessDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier *= warrior.recklessnessDamageTakenMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier /= warrior.recklessnessDamageTakenMultiplier
		},
	}).AttachStatBuff(stats.MeleeCrit, 50*core.CritRatingPerCritChance)

	warrior.Recklessness = warrior.RegisterSpell(DefaultRecklessnessStance, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_WarriorRecklesness,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 5,
			},
			SharedCD: core.Cooldown{
				Timer:    sharedTimer,
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			reckAura.Activate(sim)
			warrior.Retaliation.RelatedSelfBuff.Deactivate(sim)
			warrior.ShieldWall.RelatedSelfBuff.Deactivate(sim)
		},

		RelatedSelfBuff: reckAura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.Recklessness.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
