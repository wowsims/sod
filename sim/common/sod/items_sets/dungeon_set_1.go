package item_sets

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetWildheartRaiment = core.NewItemSet(core.ItemSet{
	Name: "Wildheart Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// +40 Attack Power, up to 23 increased damage from spells, and up to 44 increased healing from spells.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
				stats.SpellDamage:       23,
				stats.HealingPower:      44,
			})
		},
		// 2% chance on spellcast to energize you for 300 mana, 6% chance on dealing a melee autoattack to energize you for 40 Energy, and 3% chance on being hit by a melee attack to energize you for 10 Rage.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450608}
			manaMetrics := c.NewManaMetrics(actionID)
			energyMetrics := c.NewEnergyMetrics(actionID)
			rageMetrics := c.NewRageMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Druid Energize Trigger - Wildheart Raiment (Mana)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.02,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasManaBar() {
						c.AddMana(sim, 300, manaMetrics)
					}
				},
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Druid Energize Trigger - Wildheart Raiment (Energy)",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeWhiteHit,
				ProcChance: 0.06,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasEnergyBar() {
						c.AddEnergy(sim, 40, energyMetrics)
					}
				},
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
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
		// +8 All Resistances.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
		},
	},
})

var ItemSetBeaststalkerArmor = core.NewItemSet(core.ItemSet{
	Name: "Beaststalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// +40 Attack Power.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Your melee and ranged autoattacks have a 6% chance to energize you for 300 mana.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450577}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Mana Proc on Cast - Beaststalker Armor",
				Callback:   core.CallbackOnSpellHitDealt,
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
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
		},
	},
})

var ItemSetMagistersRegalia = core.NewItemSet(core.ItemSet{
	Name: "Magister's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Your spellcasts have a 6% chance to energize you for 300 mana.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450527}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Mana Proc on Cast - Magister's Regalia",
				Callback:   core.CallbackOnSpellHitDealt,
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
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
		},
	},
})

var ItemSetLightforgeArmor = core.NewItemSet(core.ItemSet{
	Name: "Lightforge Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// +40 Attack Power and up to 40 increased healing from spells.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
				stats.HealingPower:      40,
			})
		},
		// 6% chance on melee autoattack and 4% chance on spellcast to increase your damage and healing done by magical spells and effects by up to 95 for 10 sec.
		6: func(agent core.Agent) {
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
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
		},
	},
})

var ItemSetVestmentsOfTheDevout = core.NewItemSet(core.ItemSet{
	Name: "Vestments of the Devout",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Your spellcasts have a 6% chance to energize you for 300 mana.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450576}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
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
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
		},
	},
})

var ItemSetShadowcraftArmor = core.NewItemSet(core.ItemSet{
	Name: "Shadowcraft Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// +40 Attack Power.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Chance on melee attack to restore 35 energy.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 27787}
			energyMetrics := c.NewEnergyMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID: actionID,
				Name:     "Rogue Armor Energize",
				Callback: core.CallbackOnSpellHitDealt,
				ProcMask: core.ProcMaskMelee,
				PPM:      1,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasEnergyBar() {
						c.AddEnergy(sim, 35, energyMetrics)
					}
				},
			})
		},
		// +8 All Resistances.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
		},
	},
})

var ItemSetTheElements = core.NewItemSet(core.ItemSet{
	Name: "The Elements",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// +40 Attack Power, up to 23 increased damage from spells, and up to 44 increased healing from spells.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
				stats.SpellDamage:       23,
				stats.HealingPower:      44,
			})
		},
		// 6% chance on mainhand autoattack and 4% chance on spellcast to increase your damage and healing done by magical spells and effects by up to 95 for 10 sec.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450626}

			procAura := c.NewTemporaryStatsAura("The Furious Storm", core.ActionID{SpellID: 27775}, stats.Stats{stats.SpellPower: 95}, time.Second*10)
			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Item - The Furious Storm Proc (MH Auto)",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeMHAuto,
				ProcChance: 0.06,
				Handler:    handler,
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Item - The Furious Storm Proc (Spell Cast)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.04,
				Handler:    handler,
			})
		},

		// +8 All Resistances.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
		},
	},
})

var ItemSetDreadmistRaiment = core.NewItemSet(core.ItemSet{
	Name: "Dreadmist Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Your melee autoattacks and spellcasts have a 6% chance to heal you for 270 to 330 health.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450585}
			healthMetrics := c.NewHealthMetrics(core.ActionID{SpellID: 450583})

			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				c.GainHealth(sim, sim.Roll(270, 300), healthMetrics)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Heal Proc on Cast - Dreadmist Raiment (Melee Auto)",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskWhiteHit,
				ProcChance: 0.06,
				Handler:    handler,
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Heal Proc on Cast - Dreadmist Raiment (Spell Cast)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.06,
				Handler:    handler,
			})
		},
		// +8 All Resistances.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
		},
	},
})

var ItemSetBattlegearOfValor = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Valor",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// +40 Attack Power.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Chance on melee attack to heal you for 88 to 132 and energize you for 10 Rage
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450587}
			healthMetrics := c.NewHealthMetrics(core.ActionID{SpellID: 450589})
			rageMetrics := c.NewRageMetrics(core.ActionID{SpellID: 450589})

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID: actionID,
				Name:     "S03 - Warrior Armor Heal Trigger - Battlegear of Valor",
				Callback: core.CallbackOnSpellHitDealt,
				ProcMask: core.ProcMaskMelee,
				PPM:      1,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					c.GainHealth(sim, sim.Roll(88, 132), healthMetrics)
					if c.HasRageBar() {
						c.AddRage(sim, 10, rageMetrics)
					}
				},
			})
		},
		// +8 All Resistances.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
		},
	},
})
