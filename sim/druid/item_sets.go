package druid

import (
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
			c.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance)
			c.AddStat(stats.SpellCrit, 1*core.CritRatingPerCritChance)
		},
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			for _, spell := range druid.Wrath {
				if spell != nil {
					spell.BonusCritRating += 3 * core.CritRatingPerCritChance
				}
			}
			for _, spell := range druid.Starfire {
				if spell != nil {
					spell.BonusCritRating += 3 * core.CritRatingPerCritChance
				}
			}
		},
	},
})

// TODO: New Set Bonuses
var ItemSetCoagulateBloodguardsLeathers = core.NewItemSet(core.ItemSet{
	Name: "Coagulate Bloodguard's Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
		3: func(agent core.Agent) {
			// c := agent.GetCharacter()
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
					Name:     "Exiled Dreamer",
					ActionID: core.ActionID{SpellID: 449929},
					Callback: core.CallbackOnHealDealt,
					Duration: core.NeverExpires,
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if spell.ProcMask.Matches(core.ProcMaskSpellHealing) && sim.RandomFloat("Trigger Dreamstate") < .5 {
							druid.DreamstateManaRegenAura.Activate(sim)
							core.DreamstateAura(result.Target).Activate(sim)
						}
					},
				})
			}
		},
	},
})

var ItemSetKnightLieutenantsCracklingLeather = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Crackling Leather",
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

var ItemSetBloodGuardsCracklingLeather = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Crackling Leather",
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

var ItemSetKnightLieutenantsRestoredLeather = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Restored Leather",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 33)
		},
	},
})

var ItemSetBloodGuardsRestoredLeather = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Restored Leather",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 33)
		},
	},
})

var ItemSetEmeraldWatcherVestments = core.NewItemSet(core.ItemSet{
	Name: "Emerald Watcher Vestments",
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

var ItemSetEmeraldDreamkeeperGarb = core.NewItemSet(core.ItemSet{
	Name: "Emerald Dreamkeeper Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 22)
		},
	},
})
