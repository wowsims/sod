package mage

import (
	"math"
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetSorcerersRegalia = core.NewItemSet(core.ItemSet{
	Name: "Sorcerer's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Your spellcasts have a 6% chance to energize you for 300 mana.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450527}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Mana Proc on Cast - Magister's Regalia",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
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

var ItemSetArcanistMoment = core.NewItemSet(core.ItemSet{
	Name: "Arcanist Moment",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Temporal Beacons last 20% longer.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// Increases all chronomantic healing you deal by 10%.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// Each time you heal a target with Regeneration, the remaining cooldown on Rewind Time is reduced by 1 sec.
		6: func(agent core.Agent) {
			// Nothing to do
		},
	},
})

var ItemSetArcanistInsight = core.NewItemSet(core.ItemSet{
	Name: "Arcanist Insight",
	Bonuses: map[int32]core.ApplyEffect{
		// You are immune to all damage while channeling Evocation.
		2: func(agent core.Agent) {
			// May important later but for now nothing to do
		},
		// You gain 1% increased damage for 15 sec each time you cast a spell from a different school of magic.
		4: func(agent core.Agent) {
			// TODO: This is all a bit of an assumption about how this may work without having more information.
			// We may need to rework it as we get more information
			mage := agent.(MageAgent).GetMage()

			damageMultiplierPerSchool := 1.01
			auraDuration := time.Second * 15

			arcaneAura := mage.RegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Mage - Damage 4P Bonus (Arcane)",
				ActionID: core.ActionID{SpellID: 456398}.WithTag(int32(stats.SchoolIndexArcane)),
				Duration: auraDuration,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					mage.PseudoStats.DamageDealtMultiplier *= damageMultiplierPerSchool
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					mage.PseudoStats.DamageDealtMultiplier /= damageMultiplierPerSchool
				},
			})

			fireAura := mage.RegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Mage - Damage 4P Bonus (Fire)",
				ActionID: core.ActionID{SpellID: 456398}.WithTag(int32(stats.SchoolIndexFire)),
				Duration: auraDuration,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					mage.PseudoStats.DamageDealtMultiplier *= damageMultiplierPerSchool
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					mage.PseudoStats.DamageDealtMultiplier /= damageMultiplierPerSchool
				},
			})

			frostAura := mage.RegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Mage - Damage 4P Bonus (Frost)",
				ActionID: core.ActionID{SpellID: 456398}.WithTag(int32(stats.SchoolIndexFrost)),
				Duration: auraDuration,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					mage.PseudoStats.DamageDealtMultiplier *= damageMultiplierPerSchool
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					mage.PseudoStats.DamageDealtMultiplier /= damageMultiplierPerSchool
				},
			})

			core.MakePermanent(mage.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Mage - Damage 4P Bonus",
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
						return
					}
					if spell.SpellSchool.Matches(core.SpellSchoolArcane) {
						arcaneAura.Activate(sim)
					}
					if spell.SpellSchool.Matches(core.SpellSchoolFire) {
						fireAura.Activate(sim)
					}
					if spell.SpellSchool.Matches(core.SpellSchoolFrost) {
						frostAura.Activate(sim)
					}
				},
			}))
		},
		// Mage Armor increases your mana regeneration while casting by an additional 15%. Molten Armor increases your spell damage and healing by 18. Ice Armor grants 20% increased chance to trigger Fingers of Frost.
		6: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()

			bonusFoFProcChance := .20
			bonusSpiritRegenRateCasting := .15
			bonusSpellPower := 18.0

			core.MakePermanent(mage.RegisterAura(core.Aura{
				ActionID:   core.ActionID{SpellID: 456402},
				Label:      "S03 - Item - T1 - Mage - Damage 6P Bonus",
				BuildPhase: core.CharacterBuildPhaseBuffs,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					switch mage.Options.Armor {
					case proto.Mage_Options_IceArmor:
						mage.FingersOfFrostProcChance += bonusFoFProcChance
					case proto.Mage_Options_MageArmor:
						mage.PseudoStats.SpiritRegenRateCasting += bonusSpiritRegenRateCasting
					case proto.Mage_Options_MoltenArmor:
						if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
							mage.AddStat(stats.SpellPower, bonusSpellPower)
						} else {
							mage.AddStatDynamic(sim, stats.SpellPower, bonusSpellPower)
						}
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					switch mage.Options.Armor {
					case proto.Mage_Options_IceArmor:
						mage.FingersOfFrostProcChance -= bonusFoFProcChance
					case proto.Mage_Options_MageArmor:
						mage.PseudoStats.SpiritRegenRateCasting -= bonusSpiritRegenRateCasting
					case proto.Mage_Options_MoltenArmor:
						if mage.Options.Armor == proto.Mage_Options_MoltenArmor && aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
							mage.AddStat(stats.SpellPower, -bonusSpellPower)
						} else {
							mage.AddStatDynamic(sim, stats.SpellPower, -bonusSpellPower)
						}
					}
				},
			}))
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 5 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetNetherwindInsight = core.NewItemSet(core.ItemSet{
	Name: "Netherwind Insight",
	Bonuses: map[int32]core.ApplyEffect{
		// Decreases the threat generated by your Fire spells by 20%.
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Mage - Damage 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range mage.Spellbook {
						if spell.SpellSchool == core.SpellSchoolFire && spell.Flags.Matches(SpellFlagMage) {
							spell.ThreatMultiplier *= .80
						}
					}
				},
			})
		},
		// Your Pyroblast deals 20% increased damage to targets afflicted with your Fireball's periodic effect.
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Mage - Damage 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					fireballSpells := core.FilterSlice(mage.Fireball, func(spell *core.Spell) bool { return spell != nil })

					for _, spell := range mage.Pyroblast {
						if spell == nil {
							continue
						}

						oldApplyEffects := spell.ApplyEffects
						spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
							multiplier := 1.0

							for _, spell := range fireballSpells {
								if spell.Dot(target).IsActive() {
									multiplier *= 1.20
									break
								}
							}

							spell.DamageMultiplier *= multiplier
							oldApplyEffects(sim, target, spell)
							spell.DamageMultiplier /= multiplier
						}
					}
				},
			})
		},
		// Your Fireball's periodic effect gains increased damage over its duration equal to 20% of its impact damage.
		6: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.MaintainFireballDoT = true
			core.MakePermanent(mage.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 467399},
				Label:    "S03 - Item - T2 - Mage - Damage 6P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_MageFireball && result.Landed() {
						mage.BonusFireballDoTAmount += result.Damage * 1.00 / float64(spell.Dot(result.Target).NumberOfTicks)
					}
				},
			}))
		},
	},
})

var ItemSetNetherwindMoment = core.NewItemSet(core.ItemSet{
	Name: "Netherwind Moment",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Arcane Missiles refunds 10% of its base mana cost each time it deals damage.
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			actionID := core.ActionID{SpellID: 467401}
			manaMetrics := mage.NewManaMetrics(actionID)
			core.MakePermanent(mage.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Mage - Healer 2P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_MageArcaneMissilesTick && result.Landed() {
						mage.AddMana(sim, mage.ArcaneMissiles[spell.Rank].Cost.BaseCost*0.1, manaMetrics)
					}
				},
			}))
		},
		// Arcane Blast gains a 10% additional change to trigger Missile Barrage, and Missile Barrage now affects Regeneration the same way it affects Arcane Missiles.
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Mage - Healer 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					mage.ArcaneBlastMissileBarrageChance += .10
				},
			})
		},
		// Your Temporal Beacons caused by Mass Regeneration now last 30 sec.
		6: func(agent core.Agent) {
		},
	},
})

var ItemSetIllusionistsAttire = core.NewItemSet(core.ItemSet{
	Name: "Illusionist's Attire",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage done by Frost spells and effects by up to 14.
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.AddStat(stats.FrostPower, 14)
		},
		// Increases the chance to trigger your Fingers of Frost rune by an additional 15%.
		3: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.RegisterAura(core.Aura{
				Label: "S03 - Item - ZG - Mage - Frost 3P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					mage.FingersOfFrostProcChance += .15
				},
			})
		},
		// Increases damage done by your Frostbolt spell by 65%.
		5: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.RegisterAura(core.Aura{
				Label: "S03 - Item - ZG - Mage - Frost 5P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range mage.Frostbolt {
						if spell != nil {
							spell.DamageMultiplier *= 1.65
						}
					}

					if mage.SpellfrostBolt != nil {
						mage.SpellfrostBolt.DamageMultiplier *= 1.65
					}
				},
			})
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 6 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetEnigmaInsight = core.NewItemSet(core.ItemSet{
	Name: "Enigma Insight",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Fire Blast now also causes your next Fire spell to gain 50% increased critical strike chance for 10 sec.
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()

			var affectedSpells []*core.Spell

			buffAura := mage.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1213317},
				Label:    "Fire Blast",
				Duration: time.Second * 10,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					affectedSpells = core.FilterSlice(mage.Spellbook, func(spell *core.Spell) bool {
						return spell.Flags.Matches(SpellFlagMage) && spell.SpellSchool.Matches(core.SpellSchoolFire) && !spell.Flags.Matches(core.SpellFlagPassiveSpell)
					})
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range affectedSpells {
						spell.BonusCritRating += 50 * core.SpellCritRatingPerCritChance
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range affectedSpells {
						spell.BonusCritRating -= 50 * core.SpellCritRatingPerCritChance
					}
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
					if aura.RemainingDuration(sim) == aura.Duration {
						return
					}

					if !slices.Contains(affectedSpells, spell) {
						return
					}

					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt: sim.CurrentTime + core.SpellBatchWindow,
						OnAction: func(sim *core.Simulation) {
							if aura.IsActive() {
								aura.Deactivate(sim)
							}
						},
					})
				},
			})

			core.MakePermanent(mage.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Mage - Fire 2P Bonus",
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.SpellCode == SpellCode_MageFireBlast {
						buffAura.Activate(sim)
					}
				},
			}))
		},
		// Increases the damage done by your Ignite talent by 20%.
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			if mage.Talents.Ignite == 0 {
				return
			}

			mage.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Mage - Fire 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					mage.Ignite.DamageMultiplier *= 1.20
				},
			})
		},
	},
})

var ItemSetEnigmaMoment = core.NewItemSet(core.ItemSet{
	Name: "Enigma Moment",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Arcane Blast increases damage and healing done by an additional 10% per stack.
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Mage - Arcane 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					mage.ArcaneBlastDamageMultiplier += 0.10
				},
			})
		},
		// Your Mana Shield, Fire Ward, and Frost Ward absorb 50% more damage and also place a Temporal Beacon on the target for 30 sec.
		4: func(agent core.Agent) {
			// Nothing to do
		},
	},
})

var ItemSetTrappingsOfVaultedSecrets = core.NewItemSet(core.ItemSet{
	Name: "Trappings of Vaulted Secrets",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Fireball, Frostfire Bolt, and Balefire Bolt spells gain 5% increased damage for each of your Fire effects on your target, up to a maximum increased of 20%.
		3: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()

			mage.RegisterAura(core.Aura{
				Label: "S03 - Item - RAQ - Mage - Fire 3P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					dotSpells := core.FilterSlice(mage.Spellbook, func(spell *core.Spell) bool {
						return spell.Flags.Matches(SpellFlagMage) && spell.SpellSchool.Matches(core.SpellSchoolFire) && len(spell.Dots()) > 0
					})

					affectedSpells := core.FilterSlice(
						core.Flatten(
							[][]*core.Spell{
								mage.Fireball,
								{mage.FrostfireBolt},
								{mage.BalefireBolt},
							},
						), func(spell *core.Spell) bool { return spell != nil },
					)

					for _, spell := range affectedSpells {
						oldApplyEffects := spell.ApplyEffects
						spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
							multiplier := 1.0

							for _, spell := range dotSpells {
								if spell.Dot(target).IsActive() {
									multiplier += 0.05
								}
							}

							multiplier = math.Min(1.20, multiplier)

							// TODO: Additive or Multiplicative?
							spell.DamageMultiplier *= multiplier
							oldApplyEffects(sim, target, spell)
							spell.DamageMultiplier /= multiplier
						}
					}
				},
			})
		},
	},
})
