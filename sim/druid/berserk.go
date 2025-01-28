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

	actionId := core.ActionID{SpellID: 417141}
	affectedSpellClassMasks := ClassSpellMask_DruidRip | ClassSpellMask_DruidRake | ClassSpellMask_DruidTigersFury |
		ClassSpellMask_DruidShred | ClassSpellMask_DruidFerociousBite | ClassSpellMask_DruidMangleCat |
		ClassSpellMask_DruidSwipeCat | ClassSpellMask_DruidSavageRoar

	druid.BerserkAura = druid.RegisterAura(core.Aura{
		Label:    "Berserk",
		ActionID: actionId,
		Duration: time.Second * 15,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: affectedSpellClassMasks,
		IntValue:  -50,
	})

	druid.Berserk = druid.RegisterSpell(Cat|Bear, core.SpellConfig{
		ActionID: actionId,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.BerserkAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Berserk.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
