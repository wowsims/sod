package item_sets

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/hunter"
	"github.com/wowsims/sod/sim/rogue"
	"github.com/wowsims/sod/sim/warrior"
)

var ItemSetFallenRegality = core.NewItemSet(core.ItemSet{
	Name: "Fallen Regality",
	Bonuses: map[int32]core.ApplyEffect{
		// Damaging finishing moves have a 20% chance per combo point to restore 20 energy.
		// Flanking Strike's damage buff is increased by an additional 2% per stack. When striking from behind, your target takes 150% increased damage from Flanking Strike.
		// If Cleave hits fewer than its maximum number of targets, it deals 35% more damage for each unused bounce.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			aura := core.MakePermanent(character.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1232184},
				Label:    "Fallen Regality",
			}))

			switch character.Class {
			case proto.Class_ClassRogue:
				agent.(rogue.RogueAgent).GetRogue().ApplyFallenRegalityRogueBonus(aura)
			case proto.Class_ClassHunter:
				agent.(hunter.HunterAgent).GetHunter().ApplyFallenRegalityHunterBonus(aura)
			case proto.Class_ClassWarrior:
				agent.(warrior.WarriorAgent).GetWarrior().ApplyFallenRegalityWarriorBonus(aura)
			}
		},
	},
})

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
