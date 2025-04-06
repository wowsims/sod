package item_sets

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/hunter"
	"github.com/wowsims/sod/sim/rogue"
	"github.com/wowsims/sod/sim/shaman"
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
		// Hunter - The damage increaes from Mercy's and Crimson Cleaver's effects are increased by 10%.
		// Shaman - The Fire and Nature damage increases from Mercy and Crimson Cleaver are increased by 10%.
		// Warrior - The damage increaes from Mercy's and Crimson Cleaver's effects are increased by 10%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			core.MakePermanent(character.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1234318},
				Label:    "Hack and Smash",
			}))

			switch character.Class {
			case proto.Class_ClassHunter:
				agent.(hunter.HunterAgent).GetHunter().ApplyHackAndSmashHunterBonus()
			case proto.Class_ClassShaman:
				agent.(shaman.ShamanAgent).GetShaman().ApplyHackAndSmashShamanBonus()
			case proto.Class_ClassWarrior:
				agent.(warrior.WarriorAgent).GetWarrior().ApplyHackAndSmashWarriorBonus()
			}
		},
	},
})

// https://www.wowhead.com/classic-ptr/item-set=1956/tools-of-the-nathrezim
var ItemSetToolsOfTheNathrezim = core.NewItemSet(core.ItemSet{
	Name: "Tools of the Nathrezim",
	Bonuses: map[int32]core.ApplyEffect{
		// Duplicity and Deception's extra attacks now trigger 2 extra attacks.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			core.MakePermanent(character.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1231556},
				Label:    "Tools of the Nathrezim",
			}))
		},
	},
})
