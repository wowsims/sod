package item_sets

import (
	"github.com/wowsims/sod/sim/core"
)

///////////////////////////////////////////////////////////////////////////
//                                 Cloth
///////////////////////////////////////////////////////////////////////////

var ItemSetRegaliaOfUndeadCleansing = core.NewItemSet(core.ItemSet{
	Name: "Regalia of Undead Cleansing",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

var ItemSetRegaliaOfUndeadPurification = core.NewItemSet(core.ItemSet{
	Name: "Regalia of Undead Purification",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 3 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

var ItemSetRegaliaOfUndeadWarding = core.NewItemSet(core.ItemSet{
	Name: "Regalia of Undead Warding",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////

var ItemSetUndeadCleansersArmor = core.NewItemSet(core.ItemSet{
	Name: "Undead Cleanser's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

var ItemSetUndeadPurifiersArmor = core.NewItemSet(core.ItemSet{
	Name: "Undead Purifier's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

var ItemSetUndeadSlayersArmor = core.NewItemSet(core.ItemSet{
	Name: "Undead Slayer's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

var ItemSetUndeadWardersArmor = core.NewItemSet(core.ItemSet{
	Name: "Undead Warder's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

var ItemSetGarbOfTheUndeadCleansing = core.NewItemSet(core.ItemSet{
	Name: "Garb of the Undead Cleansing",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

var ItemSetGarbOfTheUndeadPurifier = core.NewItemSet(core.ItemSet{
	Name: "Garb of the Undead Purifier",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

var ItemSetGarbOfTheUndeadSlayer = core.NewItemSet(core.ItemSet{
	Name: "Garb of the Undead Slayer",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

var ItemSetGarbOfTheUndeadWarder = core.NewItemSet(core.ItemSet{
	Name: "Garb of the Undead Warder",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////

var ItemSetBattlegearOfUndeadPurification = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Undead Purification",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

var ItemSetBattlegearOfUndeadSlaying = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Undead Slaying",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})

var ItemSetBattlegearOfUndeadWarding = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Undead Warding",
	Bonuses: map[int32]core.ApplyEffect{
		// Treats your Seal of the Dawn bonus as if you were wearing 2 additional pieces of Sanctified armor. (Your total number of Sanctified armor pieces cannot exceed 8)
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SanctifiedBonus += 2
		},
	},
})
