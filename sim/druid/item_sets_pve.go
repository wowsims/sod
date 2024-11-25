package druid

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

var ItemSetLostWorshippersArmor = core.NewItemSet(core.ItemSet{
	Name: "Lost Worshipper's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.OnSpellRegistered(func(spell *core.Spell) {
				if spell.SpellCode == SpellCode_DruidWrath || spell.SpellCode == SpellCode_DruidStarfire {
					spell.BonusCritRating += 3 * core.CritRatingPerCritChance
				}
			})
		},
	},
})

var ItemSetCoagulateBloodguardsLeathers = core.NewItemSet(core.ItemSet{
	Name: "Coagulate Bloodguard's Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Strength, 10)
		},
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			// Power Shredder
			procAura := druid.GetOrRegisterAura(core.Aura{
				Label:    "Power Shredder Proc",
				ActionID: core.ActionID{SpellID: 449925},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					druid.CatForm.Cost.Multiplier -= 30
					//druid.BearForm.CostMultiplier -= 0.3
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					druid.CatForm.Cost.Multiplier += 30
					//druid.BearForm.CostMultiplier += 0.3
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == druid.CatForm.Spell /* || spell == druid.BearForm.Spell */ {
						aura.Deactivate(sim)
					}
				},
			})

			core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
				Name:     "Power Shredder",
				Callback: core.CallbackOnCastComplete,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_DruidShred {
						procAura.Activate(sim)
					}
				},
			})

			// Precise Claws should be implemented in the bear form spells when those get added back
			// Adds 2% hit while in bear/dire bear forms
		},
	},
})

var ItemSetExiledProphetsRaiment = core.NewItemSet(core.ItemSet{
	Name: "Exiled Prophet's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 4)
		},
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			// TODO: Not tested because Druid doesn't have healing spells implemented at the moment
			if druid.HasRune(proto.DruidRune_RuneFeetDreamstate) {
				core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
					Name:       "Exiled Dreamer",
					Callback:   core.CallbackOnHealDealt,
					ProcMask:   core.ProcMaskSpellHealing,
					Outcome:    core.OutcomeCrit,
					ProcChance: 0.5,
					Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
						druid.DreamstateManaRegenAura.Activate(sim)
					},
				})
			}
		},
	},
})

var ItemSetEmeraldWatcherVestments = core.NewItemSet(core.ItemSet{
	Name: "Emerald Watcher Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 12)
		},
	},
})

var ItemSetEmeraldDreamkeeperGarb = core.NewItemSet(core.ItemSet{
	Name: "Emerald Dreamkeeper Garb",
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

var ItemSetFeralheartRaiment = core.NewItemSet(core.ItemSet{
	Name: "Feralheart Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
				stats.SpellDamage:       23,
				stats.HealingPower:      44,
			})
		},
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450608}
			manaMetrics := c.NewManaMetrics(actionID)
			energyMetrics := c.NewEnergyMetrics(actionID)
			rageMetrics := c.NewRageMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Druid Energize Trigger - Wildheart Raiment (Mana)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.02,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					c.AddMana(sim, 300, manaMetrics)
				},
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Druid Energize Trigger - Wildheart Raiment (Energy)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMeleeWhiteHit,
				ProcChance: 0.06,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasEnergyBar() {
						c.AddEnergy(sim, 40, energyMetrics)
					}
				},
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Druid Energize Trigger - Wildheart Raiment (Rage)",
				Callback:   core.CallbackOnSpellHitTaken,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 0.03,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasRageBar() {
						c.AddRage(sim, 10, rageMetrics)
					}
				},
			})
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
	},
})

var ItemSetCenarionEclipse = core.NewItemSet(core.ItemSet{
	Name: "Cenarion Eclipse",
	Bonuses: map[int32]core.ApplyEffect{
		// Damage dealt by Thorns increased by 100% and duration increased by 200%.
		2: func(agent core.Agent) {
			// TODO: Thorns
		},
		// Increases your chance to hit with spells and attacks by 3%.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.MeleeHit: 3,
				stats.SpellHit: 3,
			})
		},
		// Reduces the cooldown on Starfall by 50%.
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			if !druid.HasRune(proto.DruidRune_RuneCloakStarfall) {
				return
			}

			druid.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Druid - Balance 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					druid.Starfall.CD.Duration /= 2
				},
			})
		},
	},
})

var ItemSetCenarionCunning = core.NewItemSet(core.ItemSet{
	Name: "Cenarion Cunning",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Faerie Fire and Faerie Fire (Feral) also increase the chance for all attacks to hit that target by 1% for 40 sec.
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.ImprovedFaerieFireAuras = druid.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
				return core.ImprovedFaerieFireAura(target)
			})

			core.MakePermanent(druid.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Druid - Feral 2P Bonus Trigger",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if (spell.SpellCode == SpellCode_DruidFaerieFire || spell.SpellCode == SpellCode_DruidFaerieFireFeral) && result.Landed() {
						druid.ImprovedFaerieFireAuras.Get(result.Target).Activate(sim)
					}
				},
			}))
		},
		// Periodic damage from your Rake and Rip can now be critical strikes.
		4: func(agent core.Agent) {
			// Implemented in rake.go and rip.go
		},
		// Your Rip and Ferocious Bite have a 20% chance per combo point spent to refresh the duration of Savage Roar back to its initial value.
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			if !druid.HasRune(proto.DruidRune_RuneLegsSavageRoar) {
				return
			}

			// Explicitly creating this aura for APL tracking
			core.MakePermanent(druid.RegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Druid - Feral 6P Bonus",
				ActionID: core.ActionID{SpellID: 455873},
			}))

			druid.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
				if spell == druid.SavageRoar.Spell || !druid.SavageRoarAura.IsActive() {
					return
				}

				if sim.Proc(.2*float64(comboPoints), "S03 - Item - T1 - Druid - Feral 6P Bonus") {
					druid.SavageRoarAura.Refresh(sim)
				}
			})
		},
	},
})

var ItemSetCenarionRage = core.NewItemSet(core.ItemSet{
	Name: "Cenarion Rage",
	Bonuses: map[int32]core.ApplyEffect{
		// You may cast Rebirth and Innervate while in Bear Form or Dire Bear Form.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// Reduces the cooldown of Enrage by 30 sec and it no longer reduces your armor.
		4: func(agent core.Agent) {
			// TODO: Enrage
		},
		// Bear Form and Dire Bear Form increase all threat you generate by an additional 20%, and Cower now removes all your threat against the target but has a 20 sec longer cooldown.
		6: func(agent core.Agent) {
			// TODO: Bear, Dire Bear forms
		},
	},
})

var ItemSetCenarionBounty = core.NewItemSet(core.ItemSet{
	Name: "Cenarion Bounty",
	Bonuses: map[int32]core.ApplyEffect{
		// When you cast Innervate on another player, it is also cast on you.
		2: func(agent core.Agent) {
			// TODO: Would need to rework innervate to make this work
		},
		// Casting your Healing Touch or Nourish spells gives you a 25% chance to gain Mana equal to 35% of the base cost of the spell.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// Reduces the cooldown on Tranquility by 100% and increases its healing by 100%.
		6: func(agent core.Agent) {
			// Nothing to do
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 5 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetEclipseOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Eclipse of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases the damage done and damage radius of Starfall's stars and Hurricane by 25%.
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Druid - Balance 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					affectedSpells := core.FilterSlice(
						core.Flatten([][]*DruidSpell{
							druid.Hurricane,
							{druid.Starfall, druid.StarfallTick, druid.StarfallSplash},
						}), func(spell *DruidSpell) bool { return spell != nil },
					)

					for _, spell := range affectedSpells {
						spell.DamageMultiplierAdditive += 0.25
					}
				},
			})
		},
		// Your Wrath casts have a 10% chance to summon a stand of 3 Treants to attack your target for until cancelled.
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			affectedSpellCodes := []int32{SpellCode_DruidWrath, SpellCode_DruidStarsurge}
			core.MakePermanent(druid.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Druid - Balance 4P Bonus",
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if slices.Contains(affectedSpellCodes, spell.SpellCode) && !druid.t26pcTreants.IsActive() && sim.Proc(0.10, "Summon Treants") {
						druid.t26pcTreants.EnableWithTimeout(sim, druid.t26pcTreants, time.Second*15)
					}
				},
			}))
		},
		// Your Wrath critical strikes have a 30% chance to make your next Starfire instant cast.
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			starfires := []*DruidSpell{}
			buffAura := druid.RegisterAura(core.Aura{
				ActionID:  core.ActionID{SpellID: 467088},
				Label:     "Astral Power",
				Duration:  time.Second * 15,
				MaxStacks: 3,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range druid.Starfire {
						if spell != nil {
							starfires = append(starfires, spell)
						}
					}
				},
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					for _, spell := range starfires {
						spell.DamageMultiplierAdditive -= 0.10 * float64(oldStacks)
						spell.DamageMultiplierAdditive += 0.10 * float64(newStacks)
					}
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.SpellCode == SpellCode_DruidStarfire {
						aura.Deactivate(sim)
					}
				},
			})

			procSpellCodes := []int32{SpellCode_DruidWrath, SpellCode_DruidStarsurge}
			core.MakePermanent(druid.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Druid - Balance 6P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if slices.Contains(procSpellCodes, spell.SpellCode) && result.DidCrit() && sim.Proc(0.50, "Astral Power") {
						buffAura.Activate(sim)
						buffAura.AddStack(sim)
					}
				},
			}))
		},
	},
})

var ItemSetCunningOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Cunning of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases the duration of Rake by 6 sec and its periodic damage by 50%.
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 467207},
				Label:    "S03 - Item - T2- Druid - Feral 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, dot := range druid.Rake.Dots() {
						if dot != nil {
							dot.NumberOfTicks += 2
							dot.RecomputeAuraDuration()
							oldOnSnapshot := dot.OnSnapshot
							dot.OnSnapshot = func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
								oldOnSnapshot(sim, target, dot, isRollover)
								dot.SnapshotAttackerMultiplier *= 1.50
							}
						}
					}
				},
			})
		},
		// Your critical strike chance is increased by 15% while Tiger's Fury is active.
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.RegisterAura(core.Aura{
				Label: "S03 - Item - T2- Druid - Feral 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					oldOnGain := druid.TigersFuryAura.OnGain
					druid.TigersFuryAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
						oldOnGain(aura, sim)
						druid.AddStatsDynamic(sim, stats.Stats{stats.MeleeCrit: 15 * core.CritRatingPerCritChance})
					}
					oldOnExpire := druid.TigersFuryAura.OnExpire
					druid.TigersFuryAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
						oldOnExpire(aura, sim)
						druid.AddStatsDynamic(sim, stats.Stats{stats.MeleeCrit: -15 * core.CritRatingPerCritChance})
					}
				},
			})
		},
		// Your Shred and Mangle(Cat) abilities deal 10% increased damage per your Bleed effect on the target, up to a maximum of 20% increase.
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.RegisterAura(core.Aura{
				Label: "S03 - Item - T2- Druid - Feral 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					bleedSpells := []*DruidSpell{druid.Rake, druid.Rip}
					for _, spell := range []*DruidSpell{druid.Shred, druid.MangleCat, druid.FerociousBite} {
						if spell == nil {
							continue
						}

						oldApplyEffects := spell.ApplyEffects
						spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
							modifier := 0.0
							for _, dotSpell := range bleedSpells {
								if dotSpell.Dot(target).IsActive() {
									modifier += 0.10
								}
							}
							spell.DamageMultiplierAdditive += modifier
							oldApplyEffects(sim, target, spell)
							spell.DamageMultiplierAdditive -= modifier
						}
					}
				},
			})
		},
	},
})

var ItemSetFuryOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Fury of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		// Swipe(Bear) also causes your Maul to hit 1 additional target for the next 6 sec.
		2: func(agent core.Agent) {
		},
		// Your Mangle(Bear), Swipe(Bear), Maul, and Lacerate abilities gain 5% increased critical strike chance against targets afflicted by your Lacerate.
		4: func(agent core.Agent) {
		},
		// Your Swipe now spreads your Lacerate from your primary target to other targets it strikes.
		6: func(agent core.Agent) {
		},
	},
})

var ItemSetBountyOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Bounty of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		// Your healing spell critical strikes trigger the Dreamstate effect, granting you 50% of your mana regeneration while casting for 8 sec.
		2: func(agent core.Agent) {
		},
		// Your non-periodic spell critical strikes reduce the casting time of your next Healing Touch, Regrowth, or Nourish spell by 0.5 sec.
		4: func(agent core.Agent) {
		},
		// Increases healing from Wild Growth by 10%. In addition, Wild Growth can now be used in Moonkin Form, and its healing is increased by an additional 50% in that form.
		6: func(agent core.Agent) {
		},
	},
})

var ItemSetHaruspexsGarb = core.NewItemSet(core.ItemSet{
	Name: "Haruspex's Garb",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 12.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 12)
		},
		// Reduces the cast time and global cooldown of Starfire by 0.5 sec.
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.RegisterAura(core.Aura{
				Label: "S03 - Item - ZG - Druid - Balance 3P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range druid.Starfire {
						if spell != nil {
							spell.DefaultCast.CastTime -= time.Millisecond * 500
							spell.DefaultCast.GCD -= time.Millisecond * 500
						}
					}
				},
			})
		},
		// Increases the critical strike chance of Wrath by 10%.
		5: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.RegisterAura(core.Aura{
				Label: "S03 - Item - ZG - Druid - Balance 5P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range druid.Wrath {
						if spell != nil {
							spell.BonusCritRating += 10 * core.SpellCritRatingPerCritChance
						}
					}
				},
			})
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 6 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetGenesisEclipse = core.NewItemSet(core.ItemSet{
	Name: "Genesis Eclipse",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Nature's Grace talent gains 1 additional charge each time it triggers.
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Druid - Balance 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					druid.NaturesGraceProcAura.MaxStacks += 1
				},
			})
		},
		// Increases the critical strike damage bonus of your Starfire, Starsurge, and Wrath by 60%.
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Druid - Balance 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					affectedSpells := core.FilterSlice(
						core.Flatten(
							[][]*DruidSpell{
								druid.Wrath,
								druid.Starfire,
								{druid.Starsurge},
							},
						),
						func(spell *DruidSpell) bool { return spell != nil },
					)

					for _, spell := range affectedSpells {
						spell.CritDamageBonus += 0.60
					}
				},
			})
		},
	},
})

var ItemSetGenesisCunning = core.NewItemSet(core.ItemSet{
	Name: "Genesis Cunning",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Shred no longer has a positional requirement, but deals 20% more damage if you are behind the target.
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1213171}, // Intentionally exposing it to the APL
				Label:    "S03 - Item - TAQ - Druid - Feral 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					druid.ShredPositionOverride = true
					if !druid.PseudoStats.InFrontOfTarget {
						druid.Shred.DamageMultiplierAdditive += 0.20
					}
				},
			})
		},
		// Your Mangle, Shred, and Ferocious Bite critical strikes cause your target to Bleed for 30% of the damage done over the next 4 sec sec.
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			// This is the spell used for the bleed proc.
			// https://www.wowhead.com/classic/spell=1213176/tooth-and-claw
			toothAndClawSpell := druid.RegisterSpell(Any, core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 1213176},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BonusCoefficient: 1,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label: "Tooth and Claw",
					},
					NumberOfTicks: 2,
					TickLength:    time.Second * 2,
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.Dot(target).ApplyOrRefresh(sim)
					spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
				},
			})

			core.MakePermanent(druid.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1213174},
				Label:    "S03 - Item - TAQ - Druid - Feral 4P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Outcome.Matches(core.OutcomeCrit) || !(spell == druid.Shred.Spell || spell == druid.MangleCat.Spell || spell == druid.FerociousBite.Spell) {
						return
					}

					dot := toothAndClawSpell.Dot(result.Target)
					dotDamage := result.Damage * 0.3
					if dot.IsActive() {
						dotDamage += dot.SnapshotBaseDamage * float64(dot.MaxTicksRemaining())
					}
					dot.SnapshotBaseDamage = dotDamage / float64(dot.NumberOfTicks)
					dot.SnapshotAttackerMultiplier = 1

					toothAndClawSpell.Cast(sim, result.Target)
				},
			}))
		},
	},
})

var ItemSetGenesisBounty = core.NewItemSet(core.ItemSet{
	Name: "Genesis Bounty",
	Bonuses: map[int32]core.ApplyEffect{
		// Reduces the cooldown of your Rebirth and Innervate spells by 65%.
		2: func(agent core.Agent) {
		},
		// Your critical heals with Healing Touch, Regrowth, and Nourish instantly heal the target for another 50% of the healing they dealt.
		4: func(agent core.Agent) {
		},
	},
})

var ItemSetGenesisFury = core.NewItemSet(core.ItemSet{
	Name: "Genesis Fury",
	Bonuses: map[int32]core.ApplyEffect{
		// Each time you Dodge while in Dire Bear Form, you gain 10% increased damage on your next Mangle or Swipe, stacking up to 5 times.
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			druid.Tank2PieceAqProcAura = druid.RegisterAura(core.Aura{
				Label:     "Feral 2P Bonus Proc",
				ActionID:  core.ActionID{SpellID: 1213188},
				Duration:  time.Second * 10,
				MaxStacks: 5,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					if newStacks == 0 {
						druid.MangleBear.DamageMultiplierAdditive -= float64(oldStacks) * .1
						druid.SwipeBear.DamageMultiplierAdditive -= float64(oldStacks) * .1
					} else {
						druid.MangleBear.DamageMultiplierAdditive += .1
						druid.SwipeBear.DamageMultiplierAdditive += .1
					}
				},
			})

			druid.Tank2PieceAqAura = druid.RegisterAura(core.Aura{
				Label:           "S03 - Item - TAQ - Druid - Feral 2P Bonus",
				ActionID:        core.ActionID{SpellID: 1213188},
				ActionIDForProc: core.ActionID{SpellID: 1213190},
				Duration:        core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if druid.HasRageBar() && spell.ProcMask.Matches(core.ProcMaskMelee|core.ProcMaskRanged) && result.Outcome.Matches(core.OutcomeDodge) {
						druid.Tank2PieceAqProcAura.Activate(sim)
						druid.Tank2PieceAqProcAura.AddStack(sim)
					}
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellID == druid.MangleBear.SpellID || spell.SpellID == druid.SwipeBear.SpellID {
						druid.Tank2PieceAqProcAura.SetStacks(sim, 0)
					}
				},
			})
		},
		// Reduces the cooldown on Mangle (Bear) by 1.5 sec.
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			if !druid.HasRune(proto.DruidRune_RuneHandsMangle) {
				return
			}
			druid.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Druid - Feral 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					druid.MangleBear.CD.Duration -= (time.Second + 500*time.Millisecond)
				},
			})
		},
	},
})

var ItemSetSymbolsOfUnendingLife = core.NewItemSet(core.ItemSet{
	Name: "Symbols of Unending Life",
	Bonuses: map[int32]core.ApplyEffect{
		// Your melee attacks have 5% less chance to be Dodged or Parried.
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.AddStat(stats.Expertise, 5)
		},
	},
})
