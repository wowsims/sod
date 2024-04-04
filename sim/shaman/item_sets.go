package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var OstracizedBerserksBattlemail = core.NewItemSet(core.ItemSet{
	Name: "Ostracized Berserk's Battlemail",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()

			// buffAura := c.RegisterAura(core.Aura{
			// 	Label:     "Fiery Strength",
			// 	ActionID:  core.ActionID{SpellID: 449931},
			// 	Duration:  time.Second * 12,
			// 	MaxStacks: 10,
			// 	OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			// 		c.AddStatDynamic(sim, stats.AttackPower, float64(-5*oldStacks))
			// 		c.AddStatDynamic(sim, stats.AttackPower, float64(5*newStacks))
			// 	},
			// })

			// procAura := c.RegisterAura(core.Aura{
			// 	Label:    "Fiery Strength Trigger",
			// 	Duration: core.NeverExpires,
			// })

			c.AddStat(stats.SpellPower, 18)
		},
	},
})

var ItemSetBloodGuardsPulsingMail = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Pulsing Mail",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 18)
		},
	},
})

var ItemSetEmeraldChainmail = core.NewItemSet(core.ItemSet{
	Name: "Emerald Chainmail",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 12)
		},
	},
})
