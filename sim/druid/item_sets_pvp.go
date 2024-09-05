package druid

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetKnightLieutenantsCracklingLeather = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Crackling Leather",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 18)
		},
	},
})

var ItemSetBloodGuardsCracklingLeather = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Crackling Leather",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 18)
		},
	},
})

var ItemSetKnightLieutenantsRestoredLeather = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Restored Leather",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 33)
		},
	},
})

var ItemSetBloodGuardsRestoredLeather = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Restored Leather",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 33)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetChampionsWildhide = core.NewItemSet(core.ItemSet{
	Name: "Champion's Wildhide",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

var ItemSetChampionsSanctuary = core.NewItemSet(core.ItemSet{
	Name: "Champion's Sanctuary",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 40)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

var ItemSetChampionsRefuge = core.NewItemSet(core.ItemSet{
	Name: "Champion's Refuge",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases healing done by spells and effects by up to 44.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 44)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

var ItemSetLieutenantCommandersWildhide = core.NewItemSet(core.ItemSet{
	Name: "Lieutenant Commander's Wildhide",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

var ItemSetLieutenantCommandersSanctuary = core.NewItemSet(core.ItemSet{
	Name: "Lieutenant Commander's Sanctuary",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 40)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

var ItemSetLieutenantCommandersRefuge = core.NewItemSet(core.ItemSet{
	Name: "Lieutenant Commander's Refuge",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases healing done by spells and effects by up to 44.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 44)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 5 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetFieldMarshalsWildhide = core.NewItemSet(core.ItemSet{
	Name: "Field Marshal's Wildhide",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		3: func(agent core.Agent) {
			// Nothing to do
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
	},
})

var ItemSetWarlordsWildhide = core.NewItemSet(core.ItemSet{
	Name: "Warlord's Wildhide",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		3: func(agent core.Agent) {
			// Nothing to do
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
	},
})

var ItemSetFieldMarshalsRefuge = core.NewItemSet(core.ItemSet{
	Name: "Field Marshal's Refuge",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		3: func(agent core.Agent) {
			// Nothing to do
		},
		// Increases healing done by spells and effects by up to 44.
		// Increases healing done by up to 44 and damage done by up to 15 for all magical spells and effects.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 88)
			c.AddStat(stats.SpellDamage, 15)
		},
	},
})

var ItemSetWarlordsRefuge = core.NewItemSet(core.ItemSet{
	Name: "Warlord's Refuge",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		3: func(agent core.Agent) {
			// Nothing to do
		},
		// Increases healing done by spells and effects by up to 44.
		// Increases healing done by up to 44 and damage done by up to 15 for all magical spells and effects.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 88)
			c.AddStat(stats.SpellDamage, 15)
		},
	},
})

var ItemSetFieldMarshalsSanctuary = core.NewItemSet(core.ItemSet{
	Name: "Field Marshal's Sanctuary",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		3: func(agent core.Agent) {
			// Nothing to do
		},
		// +40 Attack Power.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 40)
			c.AddStat(stats.RangedAttackPower, 40)
		},
	},
})

var ItemSetWarlordsSanctuary = core.NewItemSet(core.ItemSet{
	Name: "Warlord's Sanctuary",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Increases your movement speed by 15% while in Bear, Cat, or Travel Form. Only active outdoors.
		3: func(agent core.Agent) {
			// Nothing to do
		},
		// +40 Attack Power.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 40)
			c.AddStat(stats.RangedAttackPower, 40)
		},
	},
})
