package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (druid *Druid) applyBerserk() {
	if !druid.HasRune(proto.DruidRune_RuneBeltBerserk) {
		return
	}

	hasMangle := druid.HasRune(proto.DruidRune_RuneHandsMangle)

	actionId := core.ActionID{SpellID: 417141}
	druid.BerserkAura = druid.RegisterAura(core.Aura{
		Label:    "Berserk",
		ActionID: actionId,
		Duration: time.Second * 15,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_DruidLacerate,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_DruidMangleBear,
		Kind:      core.SpellMod_Cooldown_Multi_Flat,
		IntValue:  -100,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_DruidCatFormSpells,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -50,
	})

	druid.Berserk = druid.RegisterSpell(Cat|Bear, core.SpellConfig{
		ClassSpellMask: ClassSpellMask_DruidBerserk,
		ActionID:       actionId,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if hasMangle {
				druid.MangleBear.CD.Reset()
			}
			druid.BerserkAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Berserk.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
