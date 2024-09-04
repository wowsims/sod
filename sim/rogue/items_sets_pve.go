package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var _ = core.NewItemSet(core.ItemSet{
	Name: "Blood Corrupted Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()

			procAuras := rogue.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
				return target.RegisterAura(core.Aura{
					Label:     "Blood Corruption",
					ActionID:  core.ActionID{SpellID: 449927},
					Duration:  time.Second * 15,
					MaxStacks: 30,

					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						for si := stats.SchoolIndexPhysical; si < stats.SchoolLen; si++ {
							aura.Unit.PseudoStats.SchoolBonusDamageTaken[si] += 7
						}
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						for si := stats.SchoolIndexPhysical; si < stats.SchoolLen; si++ {
							aura.Unit.PseudoStats.SchoolBonusDamageTaken[si] -= 7
						}
					},
					OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if result.Landed() && spell.ProcMask.Matches(core.ProcMaskDirect) {
							aura.RemoveStack(sim)
						}
					},
				})
			})

			core.MakePermanent(rogue.RegisterAura(core.Aura{
				Label:    "Blood Corrupting" + rogue.Label,
				ActionID: core.ActionID{SpellID: 449928},
				OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRangedSpecial) {
						return
					}

					switch spell {
					case rogue.Backstab, rogue.Mutilate, rogue.SinisterStrike, rogue.SaberSlash, rogue.Shiv, rogue.PoisonedKnife, rogue.MainGauche, rogue.QuickDraw:
						aura := procAuras.Get(result.Target)
						aura.Activate(sim)
						aura.SetStacks(sim, aura.MaxStacks)
					}
				},
			}))
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetDarkmantleArmor = core.NewItemSet(core.ItemSet{
	Name: "Darkmantle Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Chance on melee attack to restore 35 energy.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 27787}
			energyMetrics := c.NewEnergyMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID: actionID,
				Name:     "Rogue Armor Energize",
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
				ProcMask: core.ProcMaskMeleeWhiteHit,
				PPM:      1,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasEnergyBar() {
						c.AddEnergy(sim, 35, energyMetrics)
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

var ItemSetNightSlayerThrill = core.NewItemSet(core.ItemSet{
	Name: "Nightslayer Thrill",
	Bonuses: map[int32]core.ApplyEffect{
		// Feint also grants Avoidance for 6 sec, reducing all damage taken from area of effect attacks from non-players by 50%
		2: func(agent core.Agent) {
			// Not yet implemented
		},
		// Increases the critical strike damage bonus of your Poisons by 100%.
		4: func(agent core.Agent) {
			rogue := agent.GetCharacter()
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.Flags.Matches(SpellFlagRoguePoison) {
					spell.CritDamageBonus += 1
				}
			})
		},
		// Your finishing moves have a 5% chance per combo point to make your next ability cost no energy.
		//https://www.wowhead.com/classic/spell=457342/clearcasting
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()

			var affectedSpells []*core.Spell

			aura := rogue.RegisterAura(core.Aura{
				Label:    "Clearcasting (S03 - Item - T1 - Rogue - Damage 6P Bonus)",
				ActionID: core.ActionID{SpellID: 457342},
				Duration: time.Second * 15,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					affectedSpells = core.FilterSlice(
						rogue.Spellbook,
						func(spell *core.Spell) bool {
							return spell != nil && spell.Cost != nil && spell.Cost.CostType() == core.CostTypeEnergy
						},
					)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					core.Each(affectedSpells, func(spell *core.Spell) { spell.Cost.Multiplier -= 100 })
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					core.Each(affectedSpells, func(spell *core.Spell) { spell.Cost.Multiplier += 100 })
				},
				//Poisoned Knife is ignored
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if aura.RemainingDuration(sim) == aura.Duration || spell.DefaultCast.Cost == 0 || spell.ActionID.SpellID == 425012 {
						return
					}
					aura.Deactivate(sim)
				},
			})
			rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
				if sim.Proc(.05*float64(comboPoints), "Clearcasting (S03 - Item - T1 - Rogue - Damage 6P Bonus)") {
					aura.Activate(sim)
				}
			})
		},
	},
})

var ItemSetNightSlayerBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Nightslayer Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		// While Just a Flesh Wound and Blade Dance are active, Crimson Tempest, Blunderbuss, and Fan of Knives cost 20 less Energy and generate 100% increased threat.
		2: func(agent core.Agent) {
			// Implemented in individual rune sections
		},
		// Vanish now reduces all Magic damage you take by 50% for its duration, but it no longer grants Stealth or breaks movement impairing effects.  - 457437
		4: func(agent core.Agent) {
			// Implemented in Vanish.go
		},
		// Your finishing moves have a 20% chance per combo point to make you take 50% less Physical damage from the next melee attack that hits you within 10 sec.
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			damageTaken := 0.5

			aura := rogue.RegisterAura(core.Aura{
				Label:    "Resilient (S03 - Item - T1 - Rogue - Tank 6P Bonus)",
				ActionID: core.ActionID{SpellID: 457469},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					rogue.PseudoStats.DamageTakenMultiplier *= damageTaken
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					rogue.PseudoStats.DamageTakenMultiplier /= damageTaken
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ProcMask.Matches(core.ProcMaskMelee) && result.Outcome.Matches(core.OutcomeLanded) {
						aura.Deactivate(sim)
					}
				},
			})

			rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
				if sim.Proc(0.2*float64(comboPoints), "Resilient (S03 - Item - T1 - Rogue - Tank 6P Bonus)") {
					aura.Activate(sim)
				}
			})
		},
	},
})

var ItemSetBloodfangThrill = core.NewItemSet(core.ItemSet{
	Name: "Bloodfang Thrill",
	Bonuses: map[int32]core.ApplyEffect{
		// Your opening moves have a 100% chance to make your next ability cost no energy.
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()

			var affectedSpells []*core.Spell

			aura := rogue.RegisterAura(core.Aura{
				Label:    "Clearcasting (S03 - Item - T2 - Rogue - Damage 2P Bonus)",
				ActionID: core.ActionID{SpellID: 467735},
				Duration: time.Second * 15,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					affectedSpells = core.FilterSlice(
						rogue.Spellbook,
						func(spell *core.Spell) bool {
							return spell != nil && spell.Cost != nil && spell.Cost.CostType() == core.CostTypeEnergy
						},
					)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					core.Each(affectedSpells, func(spell *core.Spell) { spell.Cost.Multiplier -= 100 })
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					core.Each(affectedSpells, func(spell *core.Spell) { spell.Cost.Multiplier += 100 })
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if aura.RemainingDuration(sim) == aura.Duration || spell.DefaultCast.Cost == 0 {
						return
					}
					aura.Deactivate(sim)
				},
			})

			core.MakePermanent(rogue.RegisterAura(core.Aura{
				Label:    "S03 - Item - T2 - Rogue - Damage 2P Bonus",
				ActionID: core.ActionID{SpellID: 467708},
				OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRangedSpecial) {
						return
					}
					switch spell {
					case rogue.Ambush, rogue.Garrote:
						aura.Activate(sim)
					}
				},
			}))
		},
		// Increases main hand weapon damage by 10%
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
					spell.DamageMultiplier *= 1.10
				}
			})
		},
			// Reduces cooldown on vanish by 2.5 minutes
		6: func(agent core.Agent) {
			// Implemented in Vanish.go
		},
	},
})

var ItemSetBloodfangBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Bloodfang Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Rolling with the Punches now also activates every time you gain a combo point.
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			if !rogue.HasRune(proto.RogueRune_RuneRollingWithThePunches) {
				return
			}

//			rogue.OnComboPointsGained(func(sim *core.Simulation) {
//				if rogue.RollingWithThePunchesProcAura.GetStacks() == 5 {
//					rogue.RollingWithThePunchesProcAura.Refresh(sim)
//				} else if rogue.RollingWithThePunchesProcAura.IsActive() {
//					rogue.RollingWithThePunchesProcAura.AddStack(sim)
//				} else {
//					rogue.RollingWithThePunchesProcAura.Activate(sim)
//					rogue.RollingWithThePunchesProcAura.AddStack(sim)
//				}
//			})
		},
		// Your Rolling with the Punches also grants you 20% increased Armor from items per stack.
		4: func(agent core.Agent) {
			//Implemented in runes.go RollingwiththePunches function TO-DO
		},
			// The cooldown on your Main Gauche resets every time your target Dodges or Parries.
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			
			core.MakePermanent(rogue.RegisterAura(core.Aura{
				Label:    "S03 - Item - T2 - Rogue - Tank 6P Bonus",
				ActionID: core.ActionID{SpellID: 467803},
				OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidDodge() || result.DidParry() {
						rogue.MainGauche.CD.Reset()
					}
				},
			}))
		},
	},
})