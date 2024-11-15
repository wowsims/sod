package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetDreadHuntersChain = core.NewItemSet(core.ItemSet{
	Name: "Dread Hunter's Chain",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
			c.AddStat(stats.RangedAttackPower, 20)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetBeastmasterArmor = core.NewItemSet(core.ItemSet{
	Name: "Beastmaster Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Your melee and ranged autoattacks have a 6% chance to energize you for 300 mana.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450577}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Mana Proc on Cast - Beaststalker Armor",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskWhiteHit,
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

var ItemSetGiantstalkerProwess = core.NewItemSet(core.ItemSet{
	Name: "Giantstalker Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Mongoose Bite also reduces its target's chance to Dodge by 1% and increases your chance to hit by 1% for 30 sec.
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()

			procBonus := stats.Stats{
				stats.SpellHit: 1,
				stats.MeleeHit: 1,
			}

			stalkerAura := hunter.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 458403},
				Label:    "Stalker",
				Duration: time.Second * 30,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatsDynamic(sim, procBonus)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatsDynamic(sim, procBonus.Invert())
				},
			})

			debuffAuras := hunter.NewEnemyAuraArray(core.MeleeHunterDodgeReductionAura)
			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Hunter - Melee 2P Bonus Trigger",
				OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_HunterMongooseBite && result.Landed() {
						debuffAuras.Get(result.Target).Activate(sim)
						stalkerAura.Activate(sim)
					}
				},
			}))
		},
		// While tracking a creature type, you deal 3% increased damage to that creature type.
		// Unsure if this stacks with the Pursuit 4p
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			// Just adding 3% damage to assume the hunter is tracking their target's type
			c.PseudoStats.DamageDealtMultiplier *= 1.03
		},
		// Mongoose Bite also activates for 5 sec whenever your target Parries or Blocks or when your melee attack misses.
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Hunter - Melee 6P Bonus Trigger",
				OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ProcMask.Matches(core.ProcMaskMelee) && (result.Outcome == core.OutcomeMiss || result.Outcome == core.OutcomeBlock || result.Outcome == core.OutcomeParry) {
						hunter.DefensiveState.Activate(sim)
					}
				},
			}))
		},
	},
})

var ItemSetGiantstalkerPursuit = core.NewItemSet(core.ItemSet{
	Name: "Giantstalker Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		// You generate 100% more threat for 8 sec after using Distracting Shot.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// While tracking a creature type, you deal 3% increased damage to that creature type.
		// Unsure if this stacks with the Prowess 4p
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			// Just adding 3% damage to assume the hunter is tracking their target's type
			c.PseudoStats.DamageDealtMultiplier *= 1.03
		},
		// Your next Shot ability within 12 sec after Aimed Shot deals 20% more damage.
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()

			if !hunter.Talents.AimedShot {
				return
			}

			procAura := hunter.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 456379},
				Label:    "S03 - Item - T1 - Hunter - Ranged 6P Bonus",
				Duration: time.Second * 12,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range hunter.Shots {
						if spell != nil {
							spell.DamageMultiplier *= 1.20
						}
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range hunter.Shots {
						if spell != nil {
							spell.DamageMultiplier /= 1.20
						}
					}
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if !spell.Flags.Matches(SpellFlagShot) || (aura.RemainingDuration(sim) == aura.Duration && spell.SpellCode == SpellCode_HunterAimedShot) {
						return
					}

					aura.Deactivate(sim)
				},
			})

			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Hunter - Ranged 6P Bonus Trigger",
				OnCastComplete: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.SpellCode == SpellCode_HunterAimedShot {
						procAura.Activate(sim)
					}
				},
			}))
		},
	},
})

var ItemSetDragonstalkerProwess = core.NewItemSet(core.ItemSet{
	Name: "Dragonstalker's Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		// Raptor Strike increases the damage done by your next other melee ability within 5 sec by 20%.
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()

			affectedSpells := make(map[*core.Spell]bool)

			procAura := hunter.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 467331},
				Label:    "Clever Strikes",
				Duration: time.Second * 5,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range hunter.MeleeSpells {
						if spell.SpellCode != SpellCode_HunterRaptorStrikeHit && spell.SpellCode != SpellCode_HunterRaptorStrike && spell.SpellCode != SpellCode_HunterWingClip {
							affectedSpells[spell] = true
						}
					}
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for spell := range affectedSpells {
						spell.DamageMultiplier *= 1.20
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for spell := range affectedSpells {
						spell.DamageMultiplier /= 1.20
					}
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if !affectedSpells[spell] {
						return
					}

					aura.Deactivate(sim)
				},
			})

			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Hunter - Melee 2P Bonus Trigger",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_HunterRaptorStrikeHit {
						procAura.Activate(sim)
					}
				},
			}))
		},
		// Increases damage dealt by your main hand weapon with Raptor Strike and Wyvern Strike by 20%.
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()

			hunter.OnSpellRegistered(func(spell *core.Spell) {
				if spell.SpellCode == SpellCode_HunterWyvernStrike || (spell.SpellCode == SpellCode_HunterRaptorStrikeHit && spell.ProcMask.Matches(core.ProcMaskMeleeMHSpecial)) {
					spell.DamageMultiplier *= 1.20
				}
			})
		},
		// Your periodic damage has a 5% chance to reset the cooldown on one of your Strike abilities. The Strike with the longest remaining cooldown is always chosen.
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label:    "S03 - Item - T2 - Hunter - Melee 6P Bonus Trigger",
				ActionID: core.ActionID{SpellID: 467334},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if sim.Proc(0.05, "T2 Melee 6PC Strike Reset") {
						maxSpell := hunter.RaptorStrike

						for _, strike := range hunter.Strikes {
							if strike.TimeToReady(sim) > maxSpell.TimeToReady(sim) {
								maxSpell = strike
							}
						}

						maxSpell.CD.Reset()
						aura.Activate(sim) // used for metrics
					}
				},
			}))
		},
	},
})

var ItemSetDragonstalkerPursuit = core.NewItemSet(core.ItemSet{
	Name: "Dragonstalker's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Aimed Shot deals 20% more damage to targets afflicted by one of your trap effects.
		2: func(agent core.Agent) {
			// Implemented in aimed_shot.go
		},
		// Your damaging Shot abilities deal 10% increased damage if the previous damaging Shot used was different than the current one.
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()

			shotSpells := []*core.Spell{}
			procAura := hunter.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 467312},
				Label:    "S03 - Item - T2 - Hunter - Ranged 4P Bonus",
				Duration: time.Second * 12,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					shotSpells = core.FilterSlice(hunter.Shots, func(s *core.Spell) bool { return s != nil })
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range shotSpells {
						if spell.SpellCode != hunter.LastShot.SpellCode {
							spell.DamageMultiplier *= 1.10
						}
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range shotSpells {
						if spell.SpellCode != hunter.LastShot.SpellCode {
							spell.DamageMultiplier /= 1.10
						}
					}
				},
			})

			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Hunter - Ranged 4P Bonus Trigger",
				OnCastComplete: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.Flags.Matches(SpellFlagShot) {
						procAura.Deactivate(sim)
						hunter.LastShot = spell
						procAura.Activate(sim)
					}
				},
			}))
		},
		//  Your Serpent Sting damage is increased by 25% of your Attack Power over its normal duration.
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Hunter - Ranged 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					hunter.SerpentStingAPCoeff += 0.25
				},
			}))
		},
	},
})

var ItemSetPredatorArmor = core.NewItemSet(core.ItemSet{
	Name: "Predator's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
			c.AddStat(stats.RangedAttackPower, 20)
		},
		// Increases the Attack Power your Beast pet gains from your attributes by 20%.
		3: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			if hunter.pet == nil {
				return
			}

			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "Predator's Armor 3P",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					oldStatInheritance := hunter.pet.GetStatInheritance()
					hunter.pet.UpdateStatInheritance(
						func(ownerStats stats.Stats) stats.Stats {
							s := oldStatInheritance(ownerStats)
							s[stats.AttackPower] *= 1.20
							return s
						},
					)
				},
			}))
		},
		// Increases the Focus regeneration of your Beast pet by 20%.
		5: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			if hunter.pet == nil {
				return
			}

			hunter.RegisterAura(core.Aura{
				Label: "Predator's Armor 5P",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					hunter.pet.AddFocusRegenMultiplier(0.20)
				},
			})
		},
	},
})

var TrappingsOfTheUnseenPath = core.NewItemSet(core.ItemSet{
	Name: "Trappings of the Unseen Path",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases the Focus regeneration of your Beast pet by 100%.
		3: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			if hunter.pet == nil {
				return
			}

			hunter.RegisterAura(core.Aura{
				Label: "Trappings of the Unseen Path 3P",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					hunter.pet.AddFocusRegenMultiplier(1.00)
				},
			})
		},
	},
})

var StrikersProwess = core.NewItemSet(core.ItemSet{
	Name: "Striker's Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases Wyvern Strike DoT by 50%
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			if !hunter.HasRune(proto.HunterRune_RuneBootsWyvernStrike) {
				return
			}

			hunter.RegisterAura(core.Aura{
				Label: "Striker's Prowess 2P",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					hunter.WyvernStrike.DoTDamageMultiplier *= 1.50
				},
			})
		},
		// Increases the Impact Damage of Mongoose Bite and all Strikes by 10%
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.RegisterAura(core.Aura{
				Label: "Striker's Prowess 4P",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range hunter.Strikes {
						spell.ImpactDamageMultiplier *= 1.10
					}
					hunter.RaptorStrikeMH.ImpactDamageMultiplier *= 1.10
					hunter.RaptorStrikeOH.ImpactDamageMultiplier *= 1.10
					hunter.MongooseBite.ImpactDamageMultiplier *= 1.10
				},
			})
		},
	},
})

var StrikersPursuit = core.NewItemSet(core.ItemSet{
	Name: "Striker's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		// Kill Shot's remaining cooldown is reduced by 50% when used on targets between 20% and 50% health, and has no cooldown while your Rapid Fire is active
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			if !hunter.HasRune(proto.HunterRune_RuneLegsKillShot) {
				return
			}

			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "Striker's Pursuit 2P",
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.SpellCode != SpellCode_HunterKillShot {
						return
					}

					if hunter.HasActiveAura("Rapid Fire") {
						spell.CD.Reset()
					} else if sim.CurrentTime > sim.Encounter.Duration/2 {
						spell.CD.Set(sim.CurrentTime + spell.CD.TimeToReady(sim)/2)
					}
				},
			}))
		},
		// Increases Kill Shot damage by 50%
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			if !hunter.HasRune(proto.HunterRune_RuneLegsKillShot) {
				return
			}

			hunter.RegisterAura(core.Aura{
				Label: "Striker's Pursuit 4P",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					hunter.KillShot.DamageMultiplier *= 1.50
				},
			})
		},
	},
})
