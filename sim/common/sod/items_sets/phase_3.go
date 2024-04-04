package item_sets

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                                 Cloth
///////////////////////////////////////////////////////////////////////////

// TODO: New Set Bonuses
var ItemSetMalevolentProphetsVestments = core.NewItemSet(core.ItemSet{
	Name: "Malevolent Prophet's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
		3: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
	},
})

var ItemSetKnightLieutenantsDreadweave = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Dreadweave",
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

var ItemSetBloodGuardsDreadweave = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Dreadweave",
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

var ItemSetKnightLieutenantsSatin = core.NewItemSet(core.ItemSet{
	Name: "Knight Lieutenant's Satin",
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

var ItemSetBloodGuardsSatin = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Satin",
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

var ItemSetEmeraldEnchantedVestments = core.NewItemSet(core.ItemSet{
	Name: "Emerald Enchanted Vestments",
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

var ItemSetEmeraldWovenGarb = core.NewItemSet(core.ItemSet{
	Name: "Emerald Woven Garb",
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

///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////

var ItemSetKnightLieutenantsLeather = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Leather",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
		},
	},
})

var ItemSetBloodGuardsLeather = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Leather",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
		},
	},
})

var ItemSetEmeraldLeathers = core.NewItemSet(core.ItemSet{
	Name: "Emerald Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

var ItemSetShunnedDevoteesChainmail = core.NewItemSet(core.ItemSet{
	Name: "Shunned Devotee's Chainmail",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
		3: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////

// TODO: New Set Bonuses
var ItemSetWailingBerserkersPlateArmor = core.NewItemSet(core.ItemSet{
	Name: "Wailing Berserker's Plate Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
		3: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
	},
})

// TODO: New Set Bonuses
var ItemSetBanishedMartyrsFullPlate = core.NewItemSet(core.ItemSet{
	Name: "Banished Martyr's Full Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
		3: func(agent core.Agent) {
			// c := agent.GetCharacter()
		},
	},
})

var ItemSetKnightLieutenantsPlate = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
		},
	},
})

var ItemSetBloodGuardsPlate = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
		},
	},
})

var ItemSetEmeraldDreamPlate = core.NewItemSet(core.ItemSet{
	Name: "Emerald Dream Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Other
///////////////////////////////////////////////////////////////////////////

var ItemSetSerpentsAscension = core.NewItemSet(core.ItemSet{
	Name: "Serpent's Ascension",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura("Serpent's Ascension Proc", core.ActionID{SpellID: 446231}, stats.Stats{stats.AttackPower: 150}, time.Second*12)

			handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{SpellID: 446233},
				Name:       "Serpent's Ascension",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged,
				ProcChance: 0.03,
				ICD:        time.Second * 120,
				Handler:    handler,
			})
		},
	},
})
