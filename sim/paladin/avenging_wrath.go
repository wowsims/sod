package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerAvengingWrath() {
	actionID := core.ActionID{SpellID: 407788}

	AvengingWrathAura := paladin.RegisterAura(core.Aura{
		Label:    "Avenging Wrath",
		ActionID: actionID,
		Duration: time.Second * 20,
	})
	AvengingWrathAura.AttachMultiplicativePseudoStatBuff(&paladin.PseudoStats.DamageDealtMultiplier, 1.2)
	AvengingWrathAura.AttachMultiplicativePseudoStatBuff(&paladin.PseudoStats.HealingDealtMultiplier, 1.2)

	core.RegisterPercentDamageModifierEffect(AvengingWrathAura, 1.2)

	paladin.avengingWrath = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_PaladinavengingWrath,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | SpellFlag_Forbearance,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.08,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			AvengingWrathAura.Activate(sim)
		},
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: paladin.avengingWrath,
		Type:  core.CooldownTypeDPS,
	})
}
