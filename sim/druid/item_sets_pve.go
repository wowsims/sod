package druid

import (
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
					druid.CatForm.CostValues.Multiplier -= 30
					//druid.BearForm.CostMultiplier -= 0.3
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					druid.CatForm.CostValues.Multiplier += 30
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
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
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
			// Implemented in starfall.go
		},
	},
})

var ItemSetCenarionCunning = core.NewItemSet(core.ItemSet{
	Name: "Cenarion Cunning",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Faerie Fire and Faerie Fire (Feral) also increase the chance for all attacks to hit that target by 1% for 40 sec.
		2: func(agent core.Agent) {
			// Implemented in faerie_fire.go
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
