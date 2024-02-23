package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) applyShamanisticRage() {
	if !shaman.HasRune(proto.ShamanRune_RuneLegsShamanisticRage) {
		return
	}

	damageTakenMultiplier := .8
	duration := time.Second * 15
	cooldown := time.Minute * 1

	actionID := core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsShamanisticRage)}
	manaMetrics := shaman.NewManaMetrics(actionID)
	srAura := shaman.GetOrRegisterAura(core.Aura{
		Label:    "Shamanistic Rage",
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMultiplier

			// Sham rage mana gain is snapshotted on cast
			// TODO: Raid mana regain
			var manaPerTick = shaman.GetCharacter().MaxMana() * .05

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 15,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					shaman.AddMana(sim, manaPerTick, manaMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMultiplier
		},
	})

	srSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cooldown,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			srAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: srSpell,
		Type:  core.CooldownTypeMana,
	})
}
