package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetBenevolentProphetsVestments = core.NewItemSet(core.ItemSet{
	Name: "Benevolent Prophet's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		// Restores 4 mana per 5 sec.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 4)
		},
		// Your Holy damage spells cause you to gain 60 increased damage and healing power for 15 sec.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()

			procAura := c.NewTemporaryStatsAura("Faith and Magic Proc", core.ActionID{SpellID: 449923}, stats.Stats{stats.SpellPower: 60}, time.Second*15)

			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if spell.SpellSchool.Matches(core.SpellSchoolHoly) {
					procAura.Activate(sim)
				}
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "Faith and Magic",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage,
				ProcChance: 1,
				Handler:    handler,
			})
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetVestmentsOfTheVirtuous = core.NewItemSet(core.ItemSet{
	Name: "Vestments of the Virtuous",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Your spellcasts have a 6% chance to energize you for 300 mana.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450576}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Mana Proc on Cast - Vestments of the Devout",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.06,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasManaBar() {
						c.AddMana(sim, 300, manaMetrics)
					}
				},
			})
		},
		// +8 All Resistances.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// +200 Armor.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
	},
})

var ItemSetDawnProphecy = core.NewItemSet(core.ItemSet{
	Name: "Dawn Prophecy",
	Bonuses: map[int32]core.ApplyEffect{
		// -0.1 sec to the casting time of Flash Heal and -0.1 sec to the casting time of Greater Heal.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// Increases your critical strike chance with spells and attacks by 2%.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.MeleeCrit: 2 * core.CritRatingPerCritChance,
				stats.SpellCrit: 2 * core.CritRatingPerCritChance,
			})
		},
		// Increases your critical strike chance with Prayer of Healing and Circle of Healing by 25%.
		6: func(agent core.Agent) {
			// Nothing to do
		},
	},
})

var ItemSetTwilightProphecy = core.NewItemSet(core.ItemSet{
	Name: "Twilight Prophecy",
	Bonuses: map[int32]core.ApplyEffect{
		// You may cast Flash Heal while in Shadowform.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// Increases your critical strike chance with spells and attacks by 2%.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.MeleeCrit: 2 * core.CritRatingPerCritChance,
				stats.SpellCrit: 2 * core.CritRatingPerCritChance,
			})
		},
		// Mind Blast critical strikes reduce the duration of your next Mind Flay by 50% while increasing its total damage by 50%.
		6: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()

			damageMultiplier := 1.50
			durationDivisor := time.Duration(2)

			buffAura := priest.GetOrRegisterAura(core.Aura{
				Label:    "Melting Faces",
				ActionID: core.ActionID{SpellID: 456549},
				Duration: core.NeverExpires,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spells := range priest.MindFlay {
						for _, spell := range spells {
							if spell != nil {
								spell.DamageMultiplier *= damageMultiplier
								for _, dot := range spell.Dots() {
									if dot != nil {
										dot.TickLength /= durationDivisor
									}
								}
							}
						}
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spells := range priest.MindFlay {
						for _, spell := range spells {
							if spell != nil {
								spell.DamageMultiplier /= damageMultiplier
								for _, dot := range spell.Dots() {
									if dot != nil {
										dot.TickLength *= durationDivisor
									}
								}
							}
						}
					}
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.SpellCode == SpellCode_PriestMindFlay {
						aura.Deactivate(sim)
					}
				},
			})

			priest.GetOrRegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Priest - Shadow 6P Bonus",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_PriestMindBlast && result.DidCrit() {
						buffAura.Activate(sim)
					}
				},
			})
		},
	},
})
