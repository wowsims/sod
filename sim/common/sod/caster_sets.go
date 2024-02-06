package sod

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetBlackfathomElementalistHide = core.NewItemSet(core.ItemSet{
	Name: "Blackfathom Elementalist's Hide",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 9)
			c.AddStat(stats.Healing, 9)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
	},
})

var ItemSetBlackfathomInvokerVestaments = core.NewItemSet(core.ItemSet{
	Name: "Twilight Invoker's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 9)
			c.AddStat(stats.Healing, 9)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
	},
})

var ItemSetHyperconductiveMendersMeditation = core.NewItemSet(core.ItemSet{
	Name: "Hyperconductive Mender's Meditation",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Spirit, 14)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 7)
		},
	},
})

var ItemSetHyperconductiveWizardsAttire = core.NewItemSet(core.ItemSet{
	Name: "Hyperconductive Wizard's Attire",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
			c.AddStat(stats.BonusArmor, 100)
		},
		3: func(agent core.Agent) {
			character := agent.GetCharacter()

			procAura := character.NewTemporaryStatsAura("Energized Hyperconductor Proc", core.ActionID{SpellID: 435978}, stats.Stats{stats.SpellPower: 40}, time.Second*10)

			handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{ItemID: 435977},
				Name:       "Energized Hyperconductor",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskDirect,
				ProcChance: 0.10,
				Handler:    handler,
			})
		},
	},
})

var ItemSetIrradiatedGarments = core.NewItemSet(core.ItemSet{
	Name: "Irradiated Garments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeCrit, 1)
			c.AddStat(stats.SpellCrit, 1)
			c.AddStat(stats.Stamina, -5)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 11)
		},
	},
})
