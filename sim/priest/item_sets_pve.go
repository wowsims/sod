package priest

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
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

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 5 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetTwilightOfTranscendence = core.NewItemSet(core.ItemSet{
	Name: "Twilight of Transcendence",
	Bonuses: map[int32]core.ApplyEffect{
		// Reduces the cooldown of your Shadow Word: Death spell by 6 sec.
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			if !priest.HasRune(proto.PriestRune_RuneHandsShadowWordDeath) {
				return
			}

			priest.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Priest - Shadow 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					priest.ShadowWordDeath.CD.Duration -= time.Second * 6
				},
			})
		},
		// Your Shadow Word: Pain has a 0.5% chance per talent point in Spirit Tap to trigger your Spirit Tap talent when it deals damage,
		// or a 20% chance per talent point when a target dies with your Shadow Word: Pain active.
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			if priest.Talents.SpiritTap == 0 {
				return
			}

			procChance := .005 * float64(priest.Talents.SpiritTap)

			core.MakePermanent(priest.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Priest - Shadow 4P Bonus",
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_PriestShadowWordPain && sim.Proc(procChance, "Proc Spirit Tap") {
						priest.SpiritTapAura.Activate(sim)
					}
				},
			}))
		},
		// While Spirit Tap is active, your Periodic damage spells deal 20% more damage.
		6: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			if priest.Talents.SpiritTap == 0 {
				return
			}

			priestDots := []*core.Dot{}
			// Mind sear ticks are separate spells, not a DoT
			mindSearticks := []*core.Spell{}

			priest.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Priest - Shadow 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range priest.Spellbook {
						if !spell.Flags.Matches(SpellFlagPriest) {
							continue
						}

						if dots := spell.Dots(); len(dots) > 0 || spell.Flags.Matches(core.SpellFlagPureDot|core.SpellFlagChanneled) {
							priestDots = append(
								priestDots,
								core.FilterSlice(dots, func(dot *core.Dot) bool { return dot != nil })...,
							)
						}
					}

					mindSearticks = core.FilterSlice(priest.MindSearTicks, func(spell *core.Spell) bool { return spell != nil })

					oldOnGain := priest.SpiritTapAura.OnGain
					priest.SpiritTapAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
						oldOnGain(aura, sim)

						for _, dot := range priestDots {
							dot.DamageMultiplier *= 1.20
						}

						for _, spell := range mindSearticks {
							spell.DamageMultiplier *= 1.20
						}
					}

					oldOnExpire := priest.SpiritTapAura.OnExpire
					priest.SpiritTapAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
						oldOnExpire(aura, sim)

						for _, dot := range priestDots {
							dot.DamageMultiplier /= 1.20
						}

						for _, spell := range mindSearticks {
							spell.DamageMultiplier /= 1.20
						}
					}
				},
			})
		},
	},
})

var ItemSetDawnOfTranscendence = core.NewItemSet(core.ItemSet{
	Name: "Dawn of Transcendence",
	Bonuses: map[int32]core.ApplyEffect{
		// Allows 15% of your Mana regeneration to continue while casting.
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.PseudoStats.SpiritRegenRateCasting += .15
		},
		// Your periodic healing has a 2% chance to make your next spell with a casting time less than 10 seconds an instant cast spell.
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()

			affectedSpells := []*core.Spell{}

			buffAura := priest.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 467543},
				Label:    "Deliverance",
				Duration: core.NeverExpires, // TODO: Verify duration
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					affectedSpells = core.FilterSlice(priest.Spellbook, func(spell *core.Spell) bool {
						return spell.Flags.Matches(SpellFlagPriest) && spell.DefaultCast.CastTime < time.Second*10
					})
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range affectedSpells {
						spell.CastTimeMultiplier -= 1
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range affectedSpells {
						spell.CastTimeMultiplier += 1
					}
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if slices.Contains(affectedSpells, spell) {
						aura.Deactivate(sim)
					}
				},
			})

			core.MakePermanent(priest.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Priest - Healer 4P Bonus",
				OnPeriodicHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ProcMask.Matches(core.ProcMaskSpellHealing) && sim.Proc(.02, "Proc Deliverance") {
						buffAura.Activate(sim)
					}
				},
			}))
		},
		// Circle of Healing and Penance also place a heal over time effect on their targets that heals for 25% as much over 15 sec.
		6: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()

			hasCoHRune := priest.HasRune(proto.PriestRune_RuneHandsCircleOfHealing)
			hasPenanceRune := priest.HasRune(proto.PriestRune_RuneHandsPenance)
			if !hasCoHRune && !hasPenanceRune {
				return
			}

			priest.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Priest - Healer 6P Bonus",
				// TODO: How is this implemented in-game?
			})
		},
	},
})

var ItemSetConfessorsRaiment = core.NewItemSet(core.ItemSet{
	Name: "Confessor's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases healing done by up to 22 and damage done by up to 7 for all magical spells and effects.
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.AddStats(stats.Stats{
				stats.HealingPower: 22,
				stats.SpellDamage:  7,
			})
		},
		// Reduces the cooldown of your Penance spell by 6 sec.
		3: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			if !priest.HasRune(proto.PriestRune_RuneHandsPenance) {
				return
			}

			priest.RegisterAura(core.Aura{
				Label: "S03 - Item - ZG - Priest - Discipline 3P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					priest.Penance.CD.Duration -= time.Second * 6
					priest.PenanceHeal.CD.Duration -= time.Second * 6
				},
			})
		},
		// Increases the damage absorbed by your Power Word: Shield spell by 20%.
		5: func(agent core.Agent) {
		},
	},
})
