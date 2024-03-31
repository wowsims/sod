package item_sets

import (
	"time"

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

var ItemSetBenevolentProphetsVestments = core.NewItemSet(core.ItemSet{
	Name: "Benevolent Prophet's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 4)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 33)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////

var ItemSetExiledProphetsRaiment = core.NewItemSet(core.ItemSet{
	Name: "Exiled Prophet's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 4)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 33)
		},
	},
})

var ItemSetLostWorshippersArmor = core.NewItemSet(core.ItemSet{
	Name: "Lost Worshipper's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 4)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 18)
		},
	},
})

var ItemSetBloodCorruptedLeathers = core.NewItemSet(core.ItemSet{
	Name: "Blood Corrupted Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance)
			c.AddStat(stats.SpellCrit, 1*core.CritRatingPerCritChance)
		},
	},
})

// Double check the name here. Feels like a typo that could be changed
// https://www.wowhead.com/classic/item-set=1640/coagulate-bloodguards-leathers
var ItemSetCoagulateBloodguardsLeathers = core.NewItemSet(core.ItemSet{
	Name: "Coagulate Bloodguard's Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Agility, 15)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

var ItemSetDreadHuntersChain = core.NewItemSet(core.ItemSet{
	Name: "Dread Hunter's Chain",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.RangedAttackPower, 24)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Agility, 15)
		},
	},
})

var ItemSetCorruptedSpiritweaversMail = core.NewItemSet(core.ItemSet{
	Name: "Corrupted Spiritweaver's Mail",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 4)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 33)
		},
	},
})

var ItemSetShunnedDevoteesChainmail = core.NewItemSet(core.ItemSet{
	Name: "Shunned Devotee's Chainmail",
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

var ItemSetOstracizedBerserksBattlemail = core.NewItemSet(core.ItemSet{
	Name: "Ostracized Berserk's Battlemail",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Strength, 15)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////

var ItemSetWailingBerserkersPlateArmor = core.NewItemSet(core.ItemSet{
	Name: "Wailing Berserker's Plate Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Strength, 15)
		},
	},
})

var ItemSetObsessedProphetsPlate = core.NewItemSet(core.ItemSet{
	Name: "Obsessed Prophet's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 4)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 33)
		},
	},
})

var ItemSetBanishedMartyrsFullPlate = core.NewItemSet(core.ItemSet{
	Name: "Banished Martyr's Full Plate",
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
				ProcChance: 0.3,
				ICD:        time.Second * 120,
				Handler:    handler,
			})
		},
	},
})
