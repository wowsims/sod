package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerAvengingWrath() {

	cd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Minute * 3,
	}

	aura := paladin.RegisterAura(core.Aura{
		Label:    "Avenging Wrath",
		ActionID: core.ActionID{SpellID: 407788},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
                        paladin.PseudoStats.DamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
                        paladin.PseudoStats.DamageDealtMultiplier /= 1.2
		},
	})

	avengingWrath := paladin.RegisterSpell(core.SpellConfig{
		ActionID: aura.ActionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: cd,
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aura.Activate(sim)
		},
	})
	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: avengingWrath,
		Type:  core.CooldownTypeDPS,
	})
}
