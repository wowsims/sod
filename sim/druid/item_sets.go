package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetLostWorshippersArmor = core.NewItemSet(core.ItemSet{
	Name: "Lost Worshipper's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.OnSpellRegistered(func(spell *core.Spell) {
				if spell.SpellCode == SpellCode_DruidWrath || spell.SpellCode == SpellCode_DruidStarfire {
					spell.BonusCritRating += 3 * core.CritRatingPerCritChance
				}
			})
		},
	},
})

var ItemSetCoagulateBloodguardsLeathers = core.NewItemSet(core.ItemSet{
	Name: "Coagulate Bloodguard's Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Strength, 10)
		},
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			// Power Shredder
			procAura := druid.GetOrRegisterAura(core.Aura{
				Label:    "Power Shredder Proc",
				ActionID: core.ActionID{SpellID: 449925},
				Duration: time.Second * 10,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					druid.CatForm.CostMultiplier -= 0.3
					druid.BearForm.CostMultiplier -= 0.3
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					druid.CatForm.CostMultiplier += 0.3
					druid.BearForm.CostMultiplier += 0.3
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == druid.CatForm.Spell || spell == druid.BearForm.Spell {
						aura.Deactivate(sim)
					}
				},
			})

			core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
				Name:     "Power Shredder",
				ActionID: core.ActionID{SpellID: 449924},
				Callback: core.CallbackOnCastComplete,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == druid.Shred.Spell {
						procAura.Activate(sim)
					}
				},
			})

			// Precise Claws should be implemented in the bear form spells when those get added back
			// Adds 2% hit while in bear/dire bear forms
		},
	},
})

var ItemSetExiledProphetsRaiment = core.NewItemSet(core.ItemSet{
	Name: "Exiled Prophet's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 4)
		},
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			// TODO: Not tested because Druid doesn't have healing spells implemented at the moment
			if druid.HasRune(proto.DruidRune_RuneFeetDreamstate) {
				core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
					Name:       "Exiled Dreamer",
					ActionID:   core.ActionID{SpellID: 449929},
					Callback:   core.CallbackOnHealDealt,
					ProcMask:   core.ProcMaskSpellHealing,
					Outcome:    core.OutcomeCrit,
					ProcChance: 0.5,
					Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
						druid.DreamstateManaRegenAura.Activate(sim)
					},
				})
			}
		},
	},
})

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

var ItemSetEmeraldWatcherVestments = core.NewItemSet(core.ItemSet{
	Name: "Emerald Watcher Vestments",
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

var ItemSetEmeraldDreamkeeperGarb = core.NewItemSet(core.ItemSet{
	Name: "Emerald Dreamkeeper Garb",
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
