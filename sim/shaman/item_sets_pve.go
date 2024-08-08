package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
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
//                            SoD Phase 3 Item Sets
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
					if shaman.PowerSurgeAura != nil && spell.SpellCode == SpellCode_ShamanLightningBolt && spell.ProcMask.Matches(core.ProcMaskSpellDamage) && result.DidCrit() && sim.Proc(.35, "Power Surge") {
						shaman.PowerSurgeAura.Activate(sim)
					}
				},
			})
		},
		// Lava Burst now also refreshes the duration of Flame Shock on your target back to 12 sec.
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.GetOrRegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Shaman - Elemental 6P Bonus",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					// Assuming only direct casts (not overloads) proc this for now
					if spell.SpellCode == SpellCode_ShamanLavaBurst && spell.Flags.Matches(core.SpellFlagAPL) && result.Landed() {
						for _, spell := range shaman.FlameShock {
							if spell != nil {
								if dot := spell.Dot(result.Target); dot.IsActive() {
									dot.NumberOfTicks = 4
									dot.Rollover(sim)
								}
							}
						}
					}
				},
			})
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
			// Implemented in talents.go
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
			shaman.makeFlurryConsumptionTrigger()

			core.MakePermanent(shaman.GetOrRegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Shaman - Tank 2P Bonus",
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Outcome.Matches(core.OutcomeParry) || result.Outcome.Matches(core.OutcomeDodge) || result.Outcome.Matches(core.OutcomeBlock) {
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
