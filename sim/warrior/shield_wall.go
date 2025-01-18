package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: Classic Update
func (warrior *Warrior) RegisterShieldWallCD() {
	if warrior.Level < 28 || warrior.OffHand().WeaponType != proto.WeaponType_WeaponTypeShield {
		return
	}

	duration := time.Duration(10+[]float64{0, 3, 5}[warrior.Talents.ImprovedShieldWall]) * time.Second
	//This is the inverse of the tooltip since it is a damage TAKEN coefficient
	damageTaken := 0.25

	actionID := core.ActionID{SpellID: 871}
	swAura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Wall",
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier *= damageTaken
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier /= damageTaken
		},
	})

	cooldownDur := time.Minute * 30

	warrior.ShieldWall = warrior.RegisterSpell(DefensiveStance, core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDur,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			swAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.ShieldWall.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}
