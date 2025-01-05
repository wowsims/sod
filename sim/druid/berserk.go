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
	hasMangle := druid.HasRune(DruidRune_RuneHandsMangle)

	actionId := core.ActionID{SpellID: 417141}
	var affectedSpells []*DruidSpell

	druid.BerserkAura = druid.RegisterAura(core.Aura{
		Label:    "Berserk",
		ActionID: actionId,
		Duration: time.Second * 15,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice([]*DruidSpell{
				druid.Rip,
				// druid.Claw, // If it would exist
				druid.Rake,
				druid.TigersFury,
				druid.Shred,
				// druid.Ravage, // If it would exist
				// druid.Pounce, // If it would exist
				druid.FerociousBite,
				druid.MangleCat,
				// druid.Sunfire, // If it would exist
				// druid.Skullbash, // If it would exist
				druid.SavageRoar,
				druid.SwipeCat,
			}, func(spell *DruidSpell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				if spell.Cost != nil {
					spell.Cost.Multiplier -= 50
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				if spell.Cost != nil {
					spell.Cost.Multiplier += 50
				}
			}
		},
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
