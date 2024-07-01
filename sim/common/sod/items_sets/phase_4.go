package item_sets

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                                 Cloth
///////////////////////////////////////////////////////////////////////////

var ItemSetThePostmaster = core.NewItemSet(core.ItemSet{
	Name: "The Postmaster",
	Bonuses: map[int32]core.ApplyEffect{
		// +50 Armor.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Armor, 50)
		},
		// +10 Fire Resistance.
		// +10 Arcane Resistance.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.ArcaneResistance, 10)
			character.AddStat(stats.FireResistance, 10)
		},
		// Increases damage and healing done by magical spells and effects by up to 12.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 12)
		},
		// Increases run speed by 5%.
		// +10 Intellect.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Intellect, 10)
		},
	},
})

var ItemSetNecropileRaiment = core.NewItemSet(core.ItemSet{
	Name: "Necropile Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// +5 Stamina.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +5 Intellect.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Intellect, 5)
		},
		// +15 All Resistances.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.ArcaneResistance, 15)
			character.AddStat(stats.FireResistance, 15)
			character.AddStat(stats.FrostResistance, 15)
			character.AddStat(stats.NatureResistance, 15)
			character.AddStat(stats.ShadowResistance, 15)
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 23)
		},
	},
})

var ItemSetIronweaveBattlesuit = core.NewItemSet(core.ItemSet{
	Name: "Ironweave Battlesuit",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases your chance to resist Silence and Interrupt effects by 10%.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// +200 Armor.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Armor, 200)
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		6: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 23)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////

var ItemSetCadaverousGarb = core.NewItemSet(core.ItemSet{
	Name: "Cadaverous Garb",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +3.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 3)
		},
		// +10 Attack Power.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.AttackPower, 10)
			character.AddStat(stats.RangedAttackPower, 10)
		},
		// +15 All Resistances.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.ArcaneResistance, 15)
			character.AddStat(stats.FireResistance, 15)
			character.AddStat(stats.FrostResistance, 15)
			character.AddStat(stats.NatureResistance, 15)
			character.AddStat(stats.ShadowResistance, 15)
		},
		// Improves your chance to hit by 2%.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 2)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

var ItemSetBloodmailRegalia = core.NewItemSet(core.ItemSet{
	Name: "Bloodmail Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +3.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 3)
		},
		// +10 Attack Power.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.AttackPower, 10)
			character.AddStat(stats.RangedAttackPower, 10)
		},
		// +15 All Resistances.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.ArcaneResistance, 15)
			character.AddStat(stats.FireResistance, 15)
			character.AddStat(stats.FrostResistance, 15)
			character.AddStat(stats.NatureResistance, 15)
			character.AddStat(stats.ShadowResistance, 15)
		},
		// Increases your chance to parry an attack by 1%.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Parry, 1)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////

var ItemSetDeathboneGuardian = core.NewItemSet(core.ItemSet{
	Name: "Deathbone Guardian",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +3.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 3)
		},
		// +50 Armor.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Armor, 50)
		},
		// +15 All Resistances.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.ArcaneResistance, 15)
			character.AddStat(stats.FireResistance, 15)
			character.AddStat(stats.FrostResistance, 15)
			character.AddStat(stats.NatureResistance, 15)
			character.AddStat(stats.ShadowResistance, 15)
		},
		// Increases run speed by 5%.
		// +10 Intellect.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Parry, 1)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Other
///////////////////////////////////////////////////////////////////////////

var ItemSetSpidersKiss = core.NewItemSet(core.ItemSet{
	Name: "Spider's Kiss",
	Bonuses: map[int32]core.ApplyEffect{
		// Chance on Hit: Immobilizes the target and lowers their armor by 100 for 10 sec.
		// Increased Defense +7.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura("Spider's Kiss", core.ActionID{SpellID: 17333}, stats.Stats{stats.Armor: -100}, time.Second*10)
			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{SpellID: 446570},
				Name:       "Echoes of the Depraved",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 0.05,
				Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
			character.AddStat(stats.Defense, 7)
		},
	},
})

var ItemSetDalRendsArms = core.NewItemSet(core.ItemSet{
	Name: "Dal'Rend's Arms",
	Bonuses: map[int32]core.ApplyEffect{
		// +50 Attack Power.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.AttackPower, 50)
			character.AddStat(stats.RangedAttackPower, 50)
		},
	},
})
