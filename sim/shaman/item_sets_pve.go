package shaman

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 2 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetElectromanticStormbringer = core.NewItemSet(core.ItemSet{
	Name: "Electromantic Stormbringer's Chain",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 12)
		},
		3: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.OnSpellRegistered(func(spell *core.Spell) {
				if spell.SpellCode == SpellCode_ShamanLightningBolt {
					spell.DefaultCast.CastTime -= time.Millisecond * 100
				}
			})
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var OstracizedBerserksBattlemail = core.NewItemSet(core.ItemSet{
	Name: "Ostracized Berserker's Battlemail",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()

			procAura := c.GetOrRegisterAura(core.Aura{
				Label:     "Fiery Strength Proc",
				ActionID:  core.ActionID{SpellID: 449932},
				Duration:  time.Second * 12,
				MaxStacks: 10,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					statsDelta := float64(newStacks-oldStacks) * 5.0
					aura.Unit.AddStatDynamic(sim, stats.AttackPower, statsDelta)
				},
			})

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:     "Fiery Strength",
				Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
				Outcome:  core.OutcomeLanded,
				ProcMask: core.ProcMaskDirect,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellSchool.Matches(core.SpellSchoolFire) {
						procAura.Activate(sim)
						procAura.AddStack(sim)
					}
				},
			})
		},
	},
})

var ItemSetEmeraldChainmail = core.NewItemSet(core.ItemSet{
	Name: "Emerald Chainmail",
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

var ItemSetEmeraldScalemail = core.NewItemSet(core.ItemSet{
	Name: "Emerald Scalemail",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
		},
	},
})

var ItemSetEmeraldLadenChain = core.NewItemSet(core.ItemSet{
	Name: "Emerald Laden Chain",
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

var ItemSetTheFiveThunders = core.NewItemSet(core.ItemSet{
	Name: "The Five Thunders",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power, up to 23 increased damage from spells, and up to 44 increased healing from spells.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
				stats.SpellDamage:       23,
				stats.HealingPower:      44,
			})
		},
		// 6% chance on mainhand autoattack and 4% chance on spellcast to increase your damage and healing done by magical spells and effects by up to 95 for 10 sec.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()

			procAura := c.NewTemporaryStatsAura("The Furious Storm", core.ActionID{SpellID: 27775}, stats.Stats{stats.SpellPower: 95}, time.Second*10)
			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "Item - The Furious Storm Proc (MH Auto)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMeleeMHAuto,
				ProcChance: 0.06,
				Handler:    handler,
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "Item - The Furious Storm Proc (Spell Cast)",
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

var ItemSetEarthfuryEruption = core.NewItemSet(core.ItemSet{
	Name: "Earthfury Eruption",
	Bonuses: map[int32]core.ApplyEffect{
		// The radius of your totems that affect friendly targets is increased to 40 yd.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// Your Lightning Bolt critical strikes have a 35% chance to reset the cooldown on Lava Burst and Chain Lightning and make the next Lava Burst, Chain Heal, or Chain Lightning within 10 sec instant.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.GetOrRegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Shaman - Elemental 4P Bonus",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_ShamanLightningBolt && spell.ProcMask.Matches(core.ProcMaskSpellDamage) && result.DidCrit() && sim.Proc(.35, "Power Surge") {
						shaman.PowerSurgeDamageAura.Activate(sim)
					}
				},
			})
		},
		// Lava Burst now also refreshes the duration of Flame Shock on your target back to 12 sec.
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			core.MakePermanent(shaman.GetOrRegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Shaman - Elemental 6P Bonus",
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.SpellCode == SpellCode_ShamanLavaBurst {
						for _, spell := range shaman.FlameShock {
							if spell == nil {
								continue
							}

							if dot := spell.Dot(shaman.CurrentTarget); dot.IsActive() {
								dot.NumberOfTicks = dot.OriginalNumberOfTicks
								dot.Rollover(sim)
							}
						}
					}
				},
			}))
		},
	},
})

var ItemSetEarthfuryRelief = core.NewItemSet(core.ItemSet{
	Name: "Earthfury Relief",
	Bonuses: map[int32]core.ApplyEffect{
		// The radius of your totems that affect friendly targets is increased to 40 yd.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// After casting your Healing Wave, Lesser Healing Wave, or Riptide spell, gives you a 25% chance to gain Mana equal to 35% of the base cost of the spell.
		4: func(agent core.Agent) {
			// Not implementing for now
		},
		// Your Healing Wave will now jump to additional nearby targets. Each jump reduces the effectiveness of the heal by 80%, and the spell will jump to up to 2 additional targets.
		6: func(agent core.Agent) {
			// Not implementing for now
		},
	},
})

var ItemSetEarthfuryImpact = core.NewItemSet(core.ItemSet{
	Name: "Earthfury Impact",
	Bonuses: map[int32]core.ApplyEffect{
		// The radius of your totems that affect friendly targets is increased to 40 yd.
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
		// Your Flurry talent grants an additional 10% increase to your attack speed.
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Shaman - Enhancement 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					shaman.bonusFlurrySpeed += .10
				},
			})
		},
	},
})

var ItemSetEarthfuryResolve = core.NewItemSet(core.ItemSet{
	Name: "Earthfury Resolve",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases your attack speed by 30% for your next 3 swings after you parry, dodge, or block.
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			flurryAura := shaman.makeFlurryAura(5)
			// The consumption trigger may not exist if the Shaman doesn't talent into Flurry
			shaman.makeFlurryConsumptionTrigger(flurryAura)

			core.MakePermanent(shaman.GetOrRegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Shaman - Tank 2P Bonus",
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidParry() || result.DidDodge() || result.DidBlock() {
						flurryAura.Activate(sim)
						flurryAura.SetStacks(sim, 3)
					}
				},
			}))
		},
		// Your parries and dodges also activate your Shield Mastery rune ability.
		4: func(agent core.Agent) {
			// Implemented in runes.go
		},
		// Your Stoneskin Totem also reduces Physical damage taken by 5% and your Windwall Totem also reduces Magical damage taken by 5%.
		6: func(agent core.Agent) {
			// Debuffs implemented in core/buffs.go, activated with a raid buff setting or in earth_totems.go and air_totems.go
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetEruptionOfTheTenStorms = core.NewItemSet(core.ItemSet{
	Name: "Eruption of the Ten Storms",
	Bonuses: map[int32]core.ApplyEffect{
		// Your spell critical strikes now have a 100% chance trigger your Elemental Focus talent.
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			if !shaman.Talents.ElementalFocus {
				return
			}

			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Elemental 2P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if shaman.isShamanDamagingSpell(spell) && result.DidCrit() {
						shaman.ClearcastingAura.Activate(sim)
						shaman.ClearcastingAura.SetStacks(sim, shaman.ClearcastingAura.MaxStacks)
					}
				},
			}))
		},
		// Loyal Beta from your Spirit of the Alpha ability now also increases Fire, Frost, and Nature damage by 5%.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			if !shaman.HasRune(proto.ShamanRune_RuneFeetSpiritOfTheAlpha) {
				return
			}

			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Elemental 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					oldOnGain := shaman.LoyalBetaAura.OnGain
					shaman.LoyalBetaAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
						oldOnGain(aura, sim)
						shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1.05
						shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.05
						shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1.05
					}
				},
			}))
		},
		// While Clearcasting is active, you deal 15% more non-Physical damage.
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			if !shaman.Talents.ElementalFocus {
				return
			}

			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Elemental 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					oldOnGain := shaman.ClearcastingAura.OnGain
					shaman.ClearcastingAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
						oldOnGain(aura, sim)
						shaman.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1.15)
					}

					oldOnExpire := shaman.ClearcastingAura.OnExpire
					shaman.ClearcastingAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
						oldOnExpire(aura, sim)
						shaman.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1 / 1.15)
					}
				},
			}))
		},
	},
})

var ItemSetResolveOfTheTenStorms = core.NewItemSet(core.ItemSet{
	Name: "Resolve of the Ten Storms",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Flame Shock also grants 30% increased chance to Block for 5 sec or until you Block an attack.
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			shieldBlockAura := shaman.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 467891},
				Label:    "Shield Block",
				Duration: time.Second * 5,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					shaman.AddStatDynamic(sim, stats.Block, 30*core.BlockRatingPerBlockChance)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					shaman.AddStatDynamic(sim, stats.Block, -30*core.BlockRatingPerBlockChance)
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					if result.DidBlock() {
						aura.Deactivate(sim)
					}
				},
			})

			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Tank 2P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_ShamanFlameShock {
						shieldBlockAura.Activate(sim)
					}
				},
			}))
		},
		// Each time you Block, your Block amount is increased by 10% of your Spell Damage for 6 sec, stacking up to 3 times.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			statDeps := []*stats.StatDependency{
				nil,
				shaman.NewDynamicMultiplyStat(stats.BlockValue, 1.10),
				shaman.NewDynamicMultiplyStat(stats.BlockValue, 1.20),
				shaman.NewDynamicMultiplyStat(stats.BlockValue, 1.30),
			}

			// Couldn't find a separate spell for this
			blockAura := shaman.RegisterAura(core.Aura{
				ActionID:  core.ActionID{SpellID: 467909},
				Label:     "S03 - Item - T2 - Shaman - Tank 4P Bonus",
				Duration:  time.Second * 6,
				MaxStacks: 3,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					if oldStacks != 0 {
						shaman.DisableDynamicStatDep(sim, statDeps[oldStacks])
					}
					if newStacks != 0 {
						shaman.EnableDynamicStatDep(sim, statDeps[newStacks])
					}
				},
			})

			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Tank 4P Bonus Trigger",
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidBlock() {
						blockAura.Activate(sim)
						blockAura.AddStack(sim)
					}
				},
			}))
		},
		// Each time you Block an attack, you have a 50% chance to trigger your Maelstrom Weapon rune.
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			if !shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
				return
			}

			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Tank 6P Bonus",
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidBlock() && sim.Proc(0.50, "T2 6P Proc Maelstrom Weapon") {
						shaman.MaelstromWeaponAura.Activate(sim)
						shaman.MaelstromWeaponAura.AddStack(sim)
					}
				},
			}))
		},
	},
})

var ItemSetImpactOfTheTenStorms = core.NewItemSet(core.ItemSet{
	Name: "Impact of the Ten Storms",
	Bonuses: map[int32]core.ApplyEffect{
		// Your chance to trigger Static Shock is increased by 12% (6% while dual-wielding)
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			if !shaman.HasRune(proto.ShamanRune_RuneBracersStaticShock) {
				return
			}

			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Enhancement 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					shaman.staticSHocksProcChance += 0.06
				},
			}))
		},
		// Main-hand Stormstrike now deals 50% more damage.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			if !shaman.Talents.Stormstrike {
				return
			}

			shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Enhancement 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					shaman.StormstrikeMH.DamageMultiplier *= 1.50
				},
			})
		},
		// Your Lightning Shield now gains a charge each time you hit a target with Lightning Bolt or Chain Lightning, up to a maximum of 9 charges.
		// In addition, your Lightning Shield can now deal critical damage.
		// Note: Only works with Static Shock
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			if !shaman.HasRune(proto.ShamanRune_RuneBracersStaticShock) {
				return
			}

			affectedSpellCodes := []int32{SpellCode_ShamanLightningBolt, SpellCode_ShamanChainLightning}
			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Enhancement 6P Bonus",
				OnInit: func(t26pAura *core.Aura, sim *core.Simulation) {
					for _, aura := range shaman.LightningShieldAuras {
						if aura == nil {
							continue
						}

						oldOnGain := aura.OnGain
						aura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
							oldOnGain(aura, sim)
							t26pAura.Activate(sim)
						}

						oldOnExpire := aura.OnExpire
						aura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
							oldOnExpire(aura, sim)
							t26pAura.Deactivate(sim)
						}
					}

					shaman.lightningShieldCanCrit = true
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					// Tested and it doesn't proc from overloads
					if slices.Contains(affectedSpellCodes, spell.SpellCode) && spell.ActionID.Tag != CastTagOverload && result.Landed() {
						shaman.ActiveShieldAura.AddStack(sim)
					}
				},
			}))
		},
	},
})

var ItemSetReliefOfTheTenStorms = core.NewItemSet(core.ItemSet{
	Name: "Relief of the Ten Storms",
	Bonuses: map[int32]core.ApplyEffect{
		// Your damaging and healing critical strikes now have a 100% chance to trigger your Water Shield, but do not consume a charge or trigger its cooldown.
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Restoration 2P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ProcMask.Matches(core.ProcMaskSpellDamage) && result.DidCrit() {
						shaman.WaterShieldRestore.Cast(sim, aura.Unit)
					}
				},
			}))
		},
		// Your Chain Lightning now also heals the target of your Earth Shield for 100% of the damage done.
		4: func(agent core.Agent) {
			// TODO: Implement Earth Shield
			// shaman := agent.(ShamanAgent).GetShaman()

			// core.MakePermanent(shaman.RegisterAura())
		},
		// Increases the healing of Chain Heal and the damage of Chain Lightning by 20%.
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Restoration 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					spells := core.Flatten([][]*core.Spell{
						shaman.ChainHeal,
						shaman.ChainHealOverload,
						shaman.ChainLightning,
						shaman.ChainLightningOverload,
					})

					for _, spell := range spells {
						if spell != nil {
							spell.DamageMultiplier *= 1.20
						}
					}
				},
			})
		},
	},
})

var ItemSetAugursRegalia = core.NewItemSet(core.ItemSet{
	Name: "Augur's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +7.
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.AddStat(stats.Defense, 7)
		},
		// Increases your chance to block attacks with a shield by 10%.
		3: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.AddStat(stats.Block, 10*core.BlockRatingPerBlockChance)
		},
		// Increases the chance to trigger your Power Surge rune by an additional 5%.
		5: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			if !shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge) {
				return
			}

			shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - ZG - Shaman - Tank 5P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					shaman.powerSurgeProcChance += .05
				},
			})
		},
	},
})
