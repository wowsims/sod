package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (mage *Mage) registerIcyVeinsSpell() {
	if !mage.HasRune(proto.MageRune_RuneLegsIcyVeins) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.MageRune_RuneLegsIcyVeins)}
	castSpeedMultiplier := 1.2
	manaCost := .03
	duration := time.Second * 20
	cooldown := time.Minute * 3

	mage.IcyVeinsAura = mage.RegisterAura(core.Aura{
		Label:    "Icy Veins",
		ActionID: actionID,
		Duration: duration,
	}).AttachMultiplyCastSpeed(&mage.Unit, castSpeedMultiplier)

	mage.IcyVeins = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.IcyVeinsAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.IcyVeins,
		Type:  core.CooldownTypeDPS,
	})
}
