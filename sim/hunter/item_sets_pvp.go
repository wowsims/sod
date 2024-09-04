package hunter

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetBloodGuardsChain = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Chain",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
			c.AddStat(stats.RangedAttackPower, 20)
		},
	},
})

var ItemSetKnightLieutenantsChain = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Chain",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
			c.AddStat(stats.RangedAttackPower, 20)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetChampionsPursuit = core.NewItemSet(core.ItemSet{
	Name: "Champion's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Agility.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Agility, 20)
		},
		// Reduces the cooldown of your Concussive Shot by 1 sec.
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

var ItemSetChampionsProwess = core.NewItemSet(core.ItemSet{
	Name: "Champion's Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Increases the duration of your Wing Clip by 2 sec.
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

var ItemSetLieutenantCommandersPursuit = core.NewItemSet(core.ItemSet{
	Name: "Lieutenant Commander's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Agility.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Agility, 20)
		},
		// Reduces the cooldown of your Concussive Shot by 1 sec.
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

var ItemSetLieutenantCommandersProwess = core.NewItemSet(core.ItemSet{
	Name: "Lieutenant Commander's Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Increases the duration of your Wing Clip by 2 sec.
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

var ItemSetWarlordsPursuit = core.NewItemSet(core.ItemSet{
	Name: "Warlord's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Agility, 20 Stamina
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.Agility:       20,
				stats.Stamina:       20,
			})
		},
		// Reduces the cooldown of your Concussive Shot by 1 sec.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Agi, 23 Spellpower.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.Agility:           20,
				stats.Spellpower:        23,
			})
		},
	},
})

var ItemSetWarlordsProwess = core.NewItemSet(core.ItemSet{
	Name: "Warlord's Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Agility.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Agility, 20)
		},
		// Increases the duration of your Wing Clip by 2 sec.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +40 Attack Power.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 40)
		},
	},
})

var ItemSetFieldMarshalsPursuit = core.NewItemSet(core.ItemSet{
	Name: "Field Marshal's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Agility, 20 stam
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.Agility:       20,
				stats.Stamina:       20,
			})
		},
		// Reduces the cooldown of your Concussive Shot by 1 sec.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Agi, 23 Spellpower.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.Agility:           20,
				stats.Spellpower:        23,
			})
		},
	},
})

var ItemSetFieldMarshalsProwess = core.NewItemSet(core.ItemSet{
	Name: "Field Marshal's Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Agility.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Agility, 20)
		},
		// Increases the duration of your Wing Clip by 2 sec.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +40 Attack Power.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 40)
		},
	},
})
