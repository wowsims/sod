package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var KnightLieutenantsLamellarPlate = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Lamellar Plate",
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

var KnightLieutenantsImbuedPlate = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Imbued Plate",
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

var ItemSetLieutenantCommandersVindication = core.NewItemSet(core.ItemSet{
	Name: "Lieutenant Commander's Vindication",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Reduces the cooldown of your Hammer of Justice by 10 sec.
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

var ItemSetLieutenantCommandersRedemption = core.NewItemSet(core.ItemSet{
	Name: "Lieutenant Commander's Redemption",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases healing done by spells and effects by up to 44.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 44)
		},
		// Reduces the cooldown of your Hammer of Justice by 10 sec.
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
