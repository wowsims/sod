package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetObsessedProphetsPlate = core.NewItemSet(core.ItemSet{
	Name: "Obsessed Prophet's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance)
			c.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.PseudoStats.SchoolBonusCritChance[stats.SchoolIndexHoly] += 3 * core.SpellCritRatingPerCritChance
		},
	},
})

var _ = core.NewItemSet(core.ItemSet{
	Name: "Emerald Encrusted Battleplate",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 22)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetSoulforgeArmor = core.NewItemSet(core.ItemSet{
	Name: "Soulforge Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power and up to 40 increased healing from spells.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
				stats.HealingPower:      40,
			})
		},
		// 6% chance on melee autoattack and 4% chance on spellcast to increase your damage and healing done by magical spells and effects by up to 95 for 10 sec.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450625}

			procAura := c.NewTemporaryStatsAura("Crusader's Wrath", core.ActionID{SpellID: 27499}, stats.Stats{stats.SpellPower: 95}, time.Second*10)
			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Item - Crusader's Wrath Proc - Lightforge Armor (Melee Auto)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMeleeWhiteHit,
				ProcChance: 0.06,
				Handler:    handler,
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Item - Crusader's Wrath Proc - Lightforge Armor (Spell Cast)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.04,
				Handler:    handler,
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

var ItemSetLawbringerRadiance = core.NewItemSet(core.ItemSet{
	Name: "Lawbringer Radiance",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// No need to model
			//(2) Set : Your Judgement of Light and Judgement of Wisdom also grant the effects of Judgement of the Crusader.
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 2)
			character.AddStat(stats.SpellCrit, 2)
		},
		6: func(agent core.Agent) {
			// Implemented in Paladin.go
			paladin := agent.(PaladinAgent).GetPaladin()
			core.MakePermanent(paladin.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Paladin - Retribution 6P Bonus",
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					paladin.lingerDuration = time.Second * 6
					paladin.enableMultiJudge = true
				},
			}))
		},
	},
})

var ItemSetLawbringerWill = core.NewItemSet(core.ItemSet{
	Name: "Lawbringer Will",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// (2) Set: Increases the block value of your shield by 30.
			character := agent.GetCharacter()
			character.AddStat(stats.BlockValue, 30)
		},
		4: func(agent core.Agent) {
			// (4) Set: Heal for 189 to 211 when you Block. (ICD: 3.5s)
			// Note: The heal does not scale with healing/spell power, but can crit.
			paladin := agent.(PaladinAgent).GetPaladin()
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 456540}

			bastionOfLight := paladin.RegisterSpell(core.SpellConfig{
				ActionID:         actionID,
				SpellSchool:      core.SpellSchoolHoly,
				DefenseType:      core.DefenseTypeMagic,
				ProcMask:         core.ProcMaskSpellHealing,
				Flags:            core.SpellFlagHelpful,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseHeal := sim.Roll(189, 211)
					spell.CalcAndDealHealing(sim, target, baseHeal, spell.OutcomeHealingCrit)
				},
			})

			handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				bastionOfLight.Cast(sim, result.Target)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Item - T1 - Paladin - Protection 4P Bonus",
				Callback:   core.CallbackOnSpellHitTaken,
				Outcome:    core.OutcomeBlock,
				ProcChance: 1.0,
				ICD:        time.Millisecond * 3500,
				Handler:    handler,
			})
		},
		6: func(agent core.Agent) {

			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Paladin - Protection 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					auras := paladin.holyShieldAura
					procs := paladin.holyShieldProc
					blockBonus := 30.0 * core.BlockRatingPerBlockChance

					for i, values := range HolyShieldValues {

						if paladin.Level < values.level {
							break
						}

						damage := values.damage

						// Holy Shield's damage is increased by 80% of shield block value.
						procs[i].ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
							sbv := paladin.BlockValue() * 0.8
							// Reminder: Holy Shield can crit, but does not miss.
							spell.CalcAndDealDamage(sim, target, (damage + sbv), spell.OutcomeMagicCrit)
						}

						// Holy Shield aura no longer has stacks...
						auras[i].MaxStacks = 0

						// ...and does not set stacks on gain...
						auras[i].OnGain = func(aura *core.Aura, sim *core.Simulation) {
							paladin.AddStatDynamic(sim, stats.Block, blockBonus)
						}

						// ...or remove stacks on block.
						auras[i].OnSpellHitTaken = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
							if result.DidBlock() {
								procs[i].Cast(sim, spell.Unit)
							}
						}
					}
				},
			})
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 5 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetFreethinkersArmor = core.NewItemSet(core.ItemSet{
	Name: "Freethinker's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.HolyPower: 14,
			})
		},
		3: func(agent core.Agent) {
			// Increases damage done by your holy shock spell by 50%
			paladin := agent.GetCharacter()
			paladin.OnSpellRegistered(func(spell *core.Spell) {
				if spell.SpellCode == SpellCode_PaladinHolyShock {
					//Damage multiplier is Additive with Infusion of Light rather than multiplicitive
					spell.DamageMultiplier += 0.5
				}
			})
		},
		5: func(agent core.Agent) {
			// Reduce cooldown of Exorcism by 3 seconds
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.RegisterAura(core.Aura{
				Label: "S03 - Item - ZG - Paladin - Caster 5P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range paladin.exorcism {
						spell.CD.Duration -= time.Second * 3
						spell.DamageMultiplier *= 1.50
					}
				},
			})
		},
	},
})

var ItemSetMercifulJudgement = core.NewItemSet(core.ItemSet{
	Name: "Merciful Judgement",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Increases critical strike chance of holy shock spell by 5%
			paladin := agent.GetCharacter()
			paladin.OnSpellRegistered(func(spell *core.Spell) {
				if spell.SpellCode == SpellCode_PaladinHolyShock {
					spell.BonusCritRating += 5.0
				}
			})
		},
		4: func(agent core.Agent) {
			//Increases damage done by your Consecration spell by 50%
			paladin := agent.GetCharacter()
			paladin.OnSpellRegistered(func(spell *core.Spell) {
				if spell.SpellCode == SpellCode_PaladinConsecration {
					spell.AOEDot().DamageMultiplier *= 1.5
				}
			})
		},
		6: func(agent core.Agent) {
			// While you are not your Beacon of Light target, your Beacon of Light target is also healed by 100% of the damage you deal
			// with Consecration, Exorcism, Holy Shock, Holy Wrath, and Hammer of Wrath
			// No need to Sim
		},
	},
})

var ItemSetRadiantJudgement = core.NewItemSet(core.ItemSet{
	Name: "Radiant Judgement",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// 2 pieces: Increases damage done by your damaging Judgements by 20% and your Judgements no longer consume your Seals on the target.
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Paladin - Retribution 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, judgeSpells := range paladin.allJudgeSpells {
						for _, judgeRankSpell := range judgeSpells {
							judgeRankSpell.DamageMultiplier *= 1.2
						}
					}

					paladin.consumeSealsOnJudge = false
				},
			})
		},
		4: func(agent core.Agent) {
			// 4 pieces: Reduces the cooldown on your Judgement ability by 5 seconds.
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Paladin - Retribution 4P Bonus",

				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					paladin.judgement.CD.Duration -= 5 * time.Second
					paladin.enableMultiJudge = false // Even though this is baseline in phase 5, we set it here to avoid breaking P4
				},
			})
		},
		6: func(agent core.Agent) {
			// 6 pieces: Your Judgement grants 1% increased Holy damage for 8 sec, stacking up to 5 times.
			paladin := agent.(PaladinAgent).GetPaladin()

			t2Judgement6pcAura := paladin.GetOrRegisterAura(core.Aura{
				Label:     "Swift Judgement",
				ActionID:  core.ActionID{SpellID: 467530},
				Duration:  time.Second * 8,
				MaxStacks: 5,

				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
					aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= (1.0 + (float64(oldStacks) * 0.01))
					aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= (1.0 + (float64(newStacks) * 0.01))
				},
			})

			paladin.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Paladin - Retribution 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					originalApplyEffects := paladin.judgement.ApplyEffects

					// Wrap the apply Judgement ApplyEffects with more Effects
					paladin.judgement.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						originalApplyEffects(sim, target, spell)
						// 6 pieces: Your Judgement grants 1% increased Holy damage for 8 sec, stacking up to 5 times.
						t2Judgement6pcAura.Activate(sim)
						t2Judgement6pcAura.AddStack(sim)
					}
				},
			})
		},
	},
})

var ItemSetWilfullJudgement = core.NewItemSet(core.ItemSet{
	Name: "Wilfull Judgement",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Increases the bonus chance to block from Holy Shield by 10%
			paladin := agent.(PaladinAgent).GetPaladin()
			if !paladin.Talents.HolyShield {
				return
			}

			blockBonus := 40.0 * core.BlockRatingPerBlockChance
			numCharges := int32(4)

			paladin.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Paladin - Protection 2P Bonus",
				OnInit: func(_ *core.Aura, _ *core.Simulation) {
					for i, hsAura := range paladin.holyShieldAura {
						if paladin.Level < HolyShieldValues[i].level {
							break
						}
						hsAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
							aura.SetStacks(sim, numCharges)
							paladin.AddStatDynamic(sim, stats.Block, blockBonus)
						}
						hsAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
							paladin.AddStatDynamic(sim, stats.Block, -blockBonus)
						}
					}
				},
			})
		},
		4: func(agent core.Agent) {
			//You take 10% reduced damage while Holy Shield is active.
			paladin := agent.(PaladinAgent).GetPaladin()
			if !paladin.Talents.HolyShield {
				return
			}

			paladin.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Paladin - Protection 4P Bonus",
				OnInit: func(_ *core.Aura, _ *core.Simulation) {
					for i, hsAura := range paladin.holyShieldAura {
						if hsAura == nil || paladin.Level < HolyShieldValues[i].level {
							break
						}
						oldOnGain := hsAura.OnGain
						oldOnExpire := hsAura.OnExpire

						hsAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
							oldOnGain(aura, sim)
							paladin.PseudoStats.DamageTakenMultiplier *= 0.9
						}
						hsAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
							oldOnExpire(aura, sim)
							paladin.PseudoStats.DamageTakenMultiplier /= 0.9
						}
					}
				},
			})
		},
		6: func(agent core.Agent) {
			// Your Reckoning Talent now has a 20% chance per talent point to trigger when
			// you block.
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Talents.Reckoning == 0 {
				return
			}

			actionID := core.ActionID{SpellID: 20178} // Reckoning proc ID
			procChance := 0.2 * float64(paladin.Talents.Reckoning)

			handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				paladin.AutoAttacks.ExtraMHAttack(sim, 1, actionID, spell.ActionID)
			}

			core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
				Name:       "Item - T2 - Paladin - Protection 6P Bonus",
				Callback:   core.CallbackOnSpellHitTaken,
				Outcome:    core.OutcomeBlock,
				ProcChance: procChance,
				Handler:    handler,
			})
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 6 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetAvengersRadiance = core.NewItemSet(core.ItemSet{
	Name: "Avenger's Radiance",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.OnSpellRegistered(func(spell *core.Spell) {
				//"S03 - Item - TAQ - Paladin - Retribution 2P Bonus",
				if spell.SpellCode == SpellCode_PaladinCrusaderStrike {
					// 2 Set: Increases Crusader Strike Damage by 50%
					spell.DamageMultiplier *= 1.5
				}
			})
		},
		4: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.OnSpellRegistered(func(spell *core.Spell) {
				//"S03 - Item - TAQ - Paladin - Retribution 4P Bonus",
				if spell.SpellCode == SpellCode_PaladinHolyWrath || spell.SpellCode == SpellCode_PaladinConsecration || spell.SpellCode == SpellCode_PaladinExorcism || spell.SpellCode == SpellCode_PaladinHolyShock || spell.SpellCode == SpellCode_PaladinHammerOfWrath {
					// 4 Set: Increases the critical strike damage bonus of your Exorcism, Holy Wrath, Holy Shock, Hammer of Wrath, and Consecration by 60%.
					spell.CritDamageBonus += 0.6
				}
			})
		},
	},
})

var ItemSetBattlegearOfEternalJustice = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Eternal Justice",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.RegisterAura(core.Aura{
				Label: "S03 - Item - RAQ - Paladin - Retribution 3P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					originalApplyEffects := paladin.crusaderStrike.ApplyEffects
					extraApplyEffects := paladin.judgement.ApplyEffects

					// Wrap the apply Crusader Strike ApplyEffects with more Effects
					paladin.crusaderStrike.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						originalApplyEffects(sim, target, spell)
						// 3 pieces: Crusader Strike now unleashes the judgement effect of your seals, but does not consume the seal
						consumeSealsOnJudgeSaved := paladin.consumeSealsOnJudge // Save current value
						paladin.consumeSealsOnJudge = false                     // Set to not consume seals
						if paladin.currentSeal.IsActive() {
							extraApplyEffects(sim, target, paladin.judgement)
						}
						paladin.consumeSealsOnJudge = consumeSealsOnJudgeSaved // Restore saved value
					}
				},
			})
		},
	},
})
