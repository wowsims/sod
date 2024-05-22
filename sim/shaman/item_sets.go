package shaman

import (
	"time"

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

			procAura := c.GetOrRegisterAura(core.Aura{
				Label:     "Fiery Strength Proc",
				ActionID:  core.ActionID{SpellID: 449932},
				Duration:  time.Second * 12,
				MaxStacks: 10,

				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					statsDelta := float64(newStacks-oldStacks) * 5.0
					aura.Unit.AddStatDynamic(sim, stats.AttackPower, statsDelta)
				},
			})

			handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.SpellSchool.Matches(core.SpellSchoolFire) {
					procAura.Activate(sim)
					procAura.AddStack(sim)
				}
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID: core.ActionID{SpellID: 449931},
				Name:     "Fiery Strength",
				Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
				ProcMask: core.ProcMaskDirect,
				Outcome:  core.OutcomeLanded,
				Handler:  handler,
			})
		},
	},
})

var ItemSetBloodGuardsMail = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Mail",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
		},
	},
})

var ItemSetBloodGuardsInscribedMail = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Inscribed Mail",
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

var ItemSetBloodGuardsPulsingMail = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Pulsing Mail",
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

var ItemSetEmeraldChainmail = core.NewItemSet(core.ItemSet{
	Name: "Emerald Chainmail",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 12)
		},
	},
})

var ItemSetEmeraldScalemail = core.NewItemSet(core.ItemSet{
	Name: "Emerald Scalemail",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
		},
	},
})

var ItemSetEmeraldLadenChain = core.NewItemSet(core.ItemSet{
	Name: "Emerald Laden Chain",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 22)
		},
	},
})
