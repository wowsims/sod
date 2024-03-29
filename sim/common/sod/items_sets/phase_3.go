package item_sets

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                                 Cloth
///////////////////////////////////////////////////////////////////////////

var ItemSetNightmareProphetsGarb = core.NewItemSet(core.ItemSet{
	Name: "Nightmare Prophet's Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
	},
})

var ItemSetMalevolentProphetsVestments = core.NewItemSet(core.ItemSet{
	Name: "Malevolent Prophet's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 18)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////
