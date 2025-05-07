package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (paladin *Paladin) activateSancAuraMastery(sim *core.Simulation, aura *core.Aura) {
	paladin.auraMasterySancActive = true
	paladin.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= (1.2 / 1.1)
}

func (paladin *Paladin) deactivateSancAuraMastery(sim *core.Simulation) {
	paladin.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= (1.2 / 1.1)
	paladin.auraMasterySancActive = false
}

func (paladin *Paladin) registerAuraMastery() {
	if !paladin.hasRune(proto.PaladinRune_RuneLegsAuraMastery) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.PaladinRune_RuneLegsAuraMastery)}
	paladin.auraMasteryAura = paladin.RegisterAura(core.Aura{
		Label:    "Aura Mastery",
		ActionID: actionID,
		Duration: time.Second * 6,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.libramOfSanctityAura != nil {
				paladin.libramOfSanctityAura.ApplyOnRefresh(func(aura *core.Aura, sim *core.Simulation) {
					if paladin.auraMasteryAura.IsActive() && !paladin.auraMasterySancActive {
						paladin.activateSancAuraMastery(sim, aura)
					}
				})
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.sanctityAura != nil && paladin.sanctityAura.IsActive() && !paladin.auraMasterySancActive {
				paladin.activateSancAuraMastery(sim, paladin.sanctityAura)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.auraMasterySancActive {
				paladin.deactivateSancAuraMastery(sim)
			}
		},
	})

	auraMastery := paladin.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},
		RelatedSelfBuff: paladin.auraMasteryAura,
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: auraMastery,
		Type:  core.CooldownTypeDPS,
	})
}
