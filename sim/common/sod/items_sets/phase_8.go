package item_sets

import (
	"github.com/wowsims/sod/sim/core"
)

var ItemSetHackAndSmash = core.NewItemSet(core.ItemSet{
	Name: "Hack and Smash",
	Bonuses: map[int32]core.ApplyEffect{
		// The Fire and Nature damage increases from Mercy and Crimson Cleaver are increased by 10%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			fireAura := character.GetAuraByID(core.ActionID{SpellID: 1231498})
			fireAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				School:     core.SpellSchoolFire,
				FloatValue: 1.30 / 1.20, // Revert the 20% and apply 30%
			})

			natureAura := character.GetAuraByID(core.ActionID{SpellID: 1231456})
			natureAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				School:     core.SpellSchoolNature,
				FloatValue: 1.30 / 1.20, // Revert the 20% and apply 30%
			})
		},
	},
})
