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
//                            SoD Phase 5 Item Sets
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
					shaman.StormstrikeMH.DamageMultiplier += 0.50
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
			if !shaman.HasRune(proto.ShamanRune_RuneHandsWaterShield) {
				return
			}

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
			shaman := agent.(ShamanAgent).GetShaman()
			if !shaman.HasRune(proto.ShamanRune_RuneLegsEarthShield) {
				return
			}

			// core.MakePermanent(shaman.RegisterAura())
		},
		// Increases the healing of Chain Heal and the damage of Chain Lightning by 20%.
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Shaman - Restoration 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					spells := core.FilterSlice(
						core.Flatten([][]*core.Spell{
							shaman.ChainHeal,
							shaman.ChainHealOverload,
							shaman.ChainLightning,
							shaman.ChainLightningOverload,
						}), func(spell *core.Spell) bool { return spell != nil },
					)

					for _, spell := range spells {
						spell.DamageMultiplierAdditive += 0.20
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

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 6 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetStormcallersEruption = core.NewItemSet(core.ItemSet{
	Name: "Stormcaller's Eruption",
	Bonuses: map[int32]core.ApplyEffect{
		// You have a 70% chance to avoid interruption caused by damage while casting Lightning Bolt, Chain Lightning, or Lava Burst, and a 10% increased chance to trigger your Elemental Focus talent.
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Shaman - Elemental 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					affectedPushbackSpells := core.FilterSlice(
						core.Flatten(
							[][]*core.Spell{
								shaman.LightningBolt,
								shaman.ChainLightning,
								{shaman.LavaBurst},
							},
						),
						func(spell *core.Spell) bool { return spell != nil },
					)

					for _, spell := range affectedPushbackSpells {
						spell.PushbackReduction += .70
					}

					shaman.elementalFocusProcChance += .10
				},
			})
		},
		// Increases the critical strike damage bonus of your Fire, Frost, and Nature spells by 60%.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.OnSpellRegistered(func(spell *core.Spell) {
				if (spell.Flags.Matches(SpellFlagShaman) || spell.Flags.Matches(SpellFlagTotem)) && spell.DefenseType == core.DefenseTypeMagic {
					spell.CritDamageBonus += 0.60
				}
			})
		},
	},
})

var ItemSetStormcallersResolve = core.NewItemSet(core.ItemSet{
	Name: "Stormcaller's Resolve",
	Bonuses: map[int32]core.ApplyEffect{
		// Damaging a target with Stormstrike, Lava Burst, or Molten Blast also reduces all damage you take by 10% for 10 sec.
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			affectedSpellCodess := []int32{SpellCode_ShamanStormstrike, SpellCode_ShamanLavaBurst, SpellCode_ShamanMoltenBlast}

			buffAura := shaman.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1213934},
				Label:    "Stormbraced",
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.90
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.90
				},
			})

			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Shaman - Elemental 2P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && slices.Contains(affectedSpellCodess, spell.SpellCode) {
						buffAura.Activate(sim)
					}
				},
			}))
		},
		// Your Spirit of the Alpha also increases your health by 10%, threat by 20%, and damage by 10% when cast on self.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			// Do a simple multiply stat because we assume that a tank shaman is using Alpha on theirself
			if shaman.HasRune(proto.ShamanRune_RuneFeetSpiritOfTheAlpha) && shaman.IsTanking() {
				shaman.PseudoStats.DamageDealtMultiplier *= 1.10
				shaman.PseudoStats.ThreatMultiplier *= 1.20
				shaman.MultiplyStat(stats.Health, 1.10)
			}
		},
	},
})

var ItemSetStormcallersRelief = core.NewItemSet(core.ItemSet{
	Name: "Stormcaller's Relief",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Riptide increases the amount healed by Chain Heal by an additional 25%.
		2: func(agent core.Agent) {
		},
		// Reduces the cast time of Chain Heal by 0.5 sec.
		4: func(agent core.Agent) {
		},
	},
})

var ItemSetStormcallersImpact = core.NewItemSet(core.ItemSet{
	Name: "Stormcaller's Impact",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases Stormstrike and Lava Lash damage by 50%. Stormstrike's damage is increased by an additional 50% when using a Two-handed weapon.
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Shaman - Enhancement 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					if shaman.StormstrikeMH != nil {
						shaman.StormstrikeMH.DamageMultiplierAdditive += core.TernaryFloat64(shaman.HasRune(proto.ShamanRune_RuneChestTwoHandedMastery), 1.00, 0.50)
					}

					if shaman.StormstrikeOH != nil {
						shaman.StormstrikeOH.DamageMultiplierAdditive += 0.50
					}

					if shaman.LavaLash != nil {
						shaman.LavaLash.DamageMultiplierAdditive += 0.50
					}
				},
			})
		},
		// Your Stormstrike, Lava Lash, and Lava Burst critical strikes cause your target to burn for 30% of the damage done over 4 sec.
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			// This is the spell used for the burn proc.
			// https://www.wowhead.com/classic/spell=1213915/burning
			burnSpell := shaman.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 1213915},
				SpellSchool: core.SpellSchoolFire,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BonusCoefficient: 1,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label: "Burning",
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

			core.MakePermanent(shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Shaman - Enhancement 4P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Outcome.Matches(core.OutcomeCrit) || !(spell == shaman.StormstrikeMH || spell == shaman.StormstrikeOH || spell == shaman.LavaLash || spell == shaman.LavaBurst) {
						return
					}

					dot := burnSpell.Dot(result.Target)
					dotDamage := result.Damage * 0.3
					if dot.IsActive() {
						dotDamage += dot.SnapshotBaseDamage * float64(dot.MaxTicksRemaining())
					}
					dot.SnapshotBaseDamage = dotDamage / float64(dot.NumberOfTicks)
					dot.SnapshotAttackerMultiplier = 1

					burnSpell.Cast(sim, result.Target)
				},
			}))
		},
	},
})

var ItemSetGiftOfTheGatheringStorm = core.NewItemSet(core.ItemSet{
	Name: "Gift of the Gathering Storm",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Lava Burst deals increased damage equal to its critical strike chance.
		3: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.RegisterAura(core.Aura{
				Label: "S03 - Item - RAQ - Shaman - Elemental 3P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					shaman.useLavaBurstCritScaling = true
				},
			})
		},
	},
})
