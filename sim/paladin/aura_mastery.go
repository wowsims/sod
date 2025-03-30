package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (paladin *Paladin) registerAuraMastery() {

	if !paladin.hasRune(proto.PaladinRune_RuneLegsAuraMastery) {
		return
	}

	cd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Minute * 2,
	}

	aura := paladin.RegisterAura(core.Aura{
		Label:    "Aura Mastery",
		ActionID: core.ActionID{SpellID: int32(proto.PaladinRune_RuneLegsAuraMastery)},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.currentPaladinAura != nil {
				if paladin.currentPaladinAura.Label == "Sanctity Aura" {
					paladin.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= (1.2 / 1.1)
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.currentPaladinAura != nil {
				if paladin.currentPaladinAura.Label == "Sanctity Aura" {
					paladin.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= (1.2 / 1.1)
				}
			}
		},
	})

	auraMastery := paladin.RegisterSpell(core.SpellConfig{
		ActionID: aura.ActionID,
		Cast: core.CastConfig{
			CD: cd,
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aura.Activate(sim)
		},
	})
	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: auraMastery,
		Type:  core.CooldownTypeDPS,
	})
}
