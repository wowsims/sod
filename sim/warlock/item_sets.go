package warlock

import (
	"github.com/wowsims/sod/sim/core"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

// TODO: New Set Bonuses
var ItemSetNightmareProphetsGarb = core.NewItemSet(core.ItemSet{
	Name: "Nightmare Prophet's Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
		3: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
	},
})
