package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

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

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetChampionsWartide = core.NewItemSet(core.ItemSet{
	Name: "Champion's Wartide",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases healing done by spells and effects by up to 44.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 44)
		},
		// Improves your chance to get a critical strike with all Shock spells by 2%.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.GetOrRegisterAura(core.Aura{
				Label:    "Shaman Shock Crit Bonus",
				ActionID: core.ActionID{SpellID: 22804},
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range core.Flatten([][]*core.Spell{shaman.EarthShock, shaman.FlameShock, shaman.FrostShock}) {
						if spell != nil {
							spell.BonusCritRating += 2 * core.CritRatingPerCritChance
						}
					}
				},
			})
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

var ItemSetChampionsThunderfist = core.NewItemSet(core.ItemSet{
	Name: "Champion's Thunderfist",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Improves your chance to get a critical strike with all Shock spells by 2%.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.GetOrRegisterAura(core.Aura{
				Label:    "Shaman Shock Crit Bonus",
				ActionID: core.ActionID{SpellID: 22804},
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range core.Flatten([][]*core.Spell{shaman.EarthShock, shaman.FlameShock, shaman.FrostShock}) {
						if spell != nil {
							spell.BonusCritRating += 2 * core.CritRatingPerCritChance
						}
					}
				},
			})
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

var ItemSetChampionsEarthshaker = core.NewItemSet(core.ItemSet{
	Name: "Champion's Earthshaker",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Improves your chance to get a critical strike with all Shock spells by 2%.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.GetOrRegisterAura(core.Aura{
				Label:    "Shaman Shock Crit Bonus",
				ActionID: core.ActionID{SpellID: 22804},
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range core.Flatten([][]*core.Spell{shaman.EarthShock, shaman.FlameShock, shaman.FrostShock}) {
						if spell != nil {
							spell.BonusCritRating += 2 * core.CritRatingPerCritChance
						}
					}
				},
			})
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetWarlordsWartide = core.NewItemSet(core.ItemSet{
	Name: "Warlord's Wartide",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Improves your chance to get a critical strike with all Shock spells by 2%.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.GetOrRegisterAura(core.Aura{
				Label:    "Shaman Shock Crit Bonus",
				ActionID: core.ActionID{SpellID: 22804},
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range core.Flatten([][]*core.Spell{shaman.EarthShock, shaman.FlameShock, shaman.FrostShock}) {
						if spell != nil {
							spell.BonusCritRating += 2 * core.CritRatingPerCritChance
						}
					}
				},
			})
		},
		// Increases healing done by spells and effects by up to 44.
		// Increases healing done by up to 44 and damage done by up to 15 for all magical spells and effects.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.HealingPower: 88,
				stats.SpellDamage:  15,
			})
		},
	},
})

var ItemSetWarlordsThunderfist = core.NewItemSet(core.ItemSet{
	Name: "Warlord's Thunderfist",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Improves your chance to get a critical strike with all Shock spells by 2%.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.GetOrRegisterAura(core.Aura{
				Label:    "Shaman Shock Crit Bonus",
				ActionID: core.ActionID{SpellID: 22804},
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range core.Flatten([][]*core.Spell{shaman.EarthShock, shaman.FlameShock, shaman.FrostShock}) {
						if spell != nil {
							spell.BonusCritRating += 2 * core.CritRatingPerCritChance
						}
					}
				},
			})
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
	},
})

var ItemSetWarlordsEarthshaker = core.NewItemSet(core.ItemSet{
	Name: "Warlord's Earthshaker",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Improves your chance to get a critical strike with all Shock spells by 2%.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.GetOrRegisterAura(core.Aura{
				Label:    "Shaman Shock Crit Bonus",
				ActionID: core.ActionID{SpellID: 22804},
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range core.Flatten([][]*core.Spell{shaman.EarthShock, shaman.FlameShock, shaman.FrostShock}) {
						if spell != nil {
							spell.BonusCritRating += 2 * core.CritRatingPerCritChance
						}
					}
				},
			})
		},
		// +40 Attack Power.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
	},
})
