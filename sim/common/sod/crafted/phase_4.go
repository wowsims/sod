package crafted

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// Keep these in alphabetical order.

// https://www.wowhead.com/classic/item-set=1792/black-dragon-mail
var ItemSetBlackDragonMail = core.NewItemSet(core.ItemSet{
	Name: "Black Dragon Mail",
	ID:   1792,
	Bonuses: map[int32]core.ApplyEffect{
		// Improves your chance to hit with all spells and attacks by 1%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 1)
			character.AddStat(stats.SpellHit, 1)
		},
		// Improves your chance to get a critical strike with all spells and attacks by 2%.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 2*core.CritRatingPerCritChance)
			character.AddStat(stats.SpellCrit, 2*core.SpellCritRatingPerCritChance)
		},
		// +10 Fire Resistance.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.FireResistance, 10)
		},
	},
})

// https://www.wowhead.com/classic/item-set=1790/blue-dragon-mail
var ItemSetBlueDragonMail = core.NewItemSet(core.ItemSet{
	Name: "Blue Dragon Mail",
	ID:   1790,
	Bonuses: map[int32]core.ApplyEffect{
		// +4 All Resistances.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(4)
		},
		// Increases damage and healing done by magical spells and effects by up to 28.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 28)
		},
	},
})

// https://www.wowhead.com/classic/item-set=1793/devilsaur-armor
var ItemSetDevilsaurArmor = core.NewItemSet(core.ItemSet{
	Name: "Devilsaur Armor",
	ID:   1793,
	Bonuses: map[int32]core.ApplyEffect{
		// +10 Fire Resistance.
		// Improves your chance to hit with all spells and attacks by 2%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.FireResistance, 10)
			character.AddStat(stats.MeleeHit, 2*core.MeleeHitRatingPerHitChance)
			character.AddStat(stats.SpellHit, 2*core.SpellHitRatingPerHitChance)
		},
	},
})

// https://www.wowhead.com/classic/item-set=1791/green-dragon-mail
var ItemSetGreenDragonMail = core.NewItemSet(core.ItemSet{
	Name: "Green Dragon Mail",
	ID:   1791,
	Bonuses: map[int32]core.ApplyEffect{
		// Restores 3 mana per 5 sec.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MP5, 3)
		},
		// Allows 15% of your Mana regeneration to continue while casting.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SpiritRegenRateCasting += .15
		},
	},
})

// https://www.wowhead.com/classic/item-set=1789/volcanic-armor
var ItemSetVolcanicArmor = core.NewItemSet(core.ItemSet{
	Name: "Volcanic Armor",
	ID:   1789,
	Bonuses: map[int32]core.ApplyEffect{
		// +10 Fire Resistance.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.FireResistance, 10)
		},
		// 5% chance of dealing 15 to 25 Fire damage on a successful melee attack.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			procSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 9057},
				SpellSchool: core.SpellSchoolFire,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, sim.Roll(15, 25), spell.OutcomeMagicHitAndCrit)
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Firebolt Trigger (Volcanic Armor)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: .05,
				Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					procSpell.Cast(sim, result.Target)
				},
			})
		},
	},
})
