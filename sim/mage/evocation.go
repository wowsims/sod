package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (mage *Mage) registerEvocationCD() {
	actionID := core.ActionID{SpellID: 12051}
	channelTime := time.Second * 8
	cooldown := time.Minute * 8

	tickLength := time.Millisecond * 250
	maxTicks := int32(channelTime / tickLength)

	manaRegenAura := mage.RegisterAura(core.Aura{
		Label:    "Evocation Regen",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.PseudoStats.SpiritRegenMultiplier += 15
			mage.PseudoStats.ForceFullSpiritRegen = true
			mage.UpdateManaRegenRates()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.PseudoStats.SpiritRegenMultiplier -= 15
			mage.PseudoStats.ForceFullSpiritRegen = false
			mage.UpdateManaRegenRates()
		},
	})

	mage.Evocation = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagHelpful | core.SpellFlagChanneled | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: cooldown,
			},
		},

		Dot: core.DotConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label: "Evocation",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					manaRegenAura.Activate(sim)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					manaRegenAura.Deactivate(sim)
				},
			},
			NumberOfTicks: maxTicks,
			TickLength:    tickLength,
			OnTick:        func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.SelfHot().Apply(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.Evocation,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return false // Require APL usage
		},
	})
}
