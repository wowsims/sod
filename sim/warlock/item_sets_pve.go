package warlock

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetNightmareProphetsGarb = core.NewItemSet(core.ItemSet{
	Name: "Nightmare Prophet's Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			warlock.shadowSparkAura = warlock.GetOrRegisterAura(core.Aura{
				Label:     "Shadow Spark Proc",
				ActionID:  core.ActionID{SpellID: 450013},
				Duration:  time.Second * 12,
				MaxStacks: 2,
			})

			core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
				Label: "Shadow Spark",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_WarlockShadowCleave && result.Landed() {
						warlock.shadowSparkAura.Activate(sim)
						warlock.shadowSparkAura.AddStack(sim)
					}
				},
			}))
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetDeathmistRaiment = core.NewItemSet(core.ItemSet{
	Name: "Deathmist Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Your melee autoattacks and spellcasts have a 6% chance to heal you for 270 to 330 health.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			manaMetrics := c.NewManaMetrics(core.ActionID{SpellID: 450583})

			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if c.HasManaBar() {
					c.AddMana(sim, sim.Roll(270, 300), manaMetrics)
				}
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Heal Proc on Cast - Dreadmist Raiment (Melee Auto)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskWhiteHit,
				ProcChance: 0.06,
				Handler:    handler,
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Heal Proc on Cast - Dreadmist Raiment (Spell Cast)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.06,
				Handler:    handler,
			})
		},
		// +8 All Resistances.
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
		// +200 Armor.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
	},
})

var ItemSetCorruptedFelheart = core.NewItemSet(core.ItemSet{
	Name: "Corrupted Felheart",
	Bonuses: map[int32]core.ApplyEffect{
		// Lifetap generates 50% more mana and 100% less threat.
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.RegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Warlock - Damage 2P Bonus",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.LifeTap {
						spell.DamageMultiplier *= 1.5
						spell.ThreatMultiplier *= -1
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.LifeTap {
						spell.DamageMultiplier /= 1.5
						spell.ThreatMultiplier *= -1
					}
				},
			})
		},
		// Increases your critical strike chance with spells and attacks by 2%.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.MeleeCrit: 2 * core.CritRatingPerCritChance,
				stats.SpellCrit: 2 * core.CritRatingPerCritChance,
			})
		},
		// Your Nightfall talent has a 4% increased chance to trigger. Your Immolate periodic damage has a 4% chance to grant Fire Trance, reducing the cast time of your next Incinerate or Immolate by 100%.
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			// Nightfall aspect implemented in talents.go

			affectedSpellCodes := []int32{SpellCode_WarlockImmolate, SpellCode_WarlockShadowflame, SpellCode_WarlockIncinerate}
			fireTranceAura := warlock.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457558},
				Label:    "Fire Trance",
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.Immolate {
						spell.CastTimeMultiplier -= 1
					}
					if warlock.Shadowflame != nil {
						warlock.Shadowflame.CastTimeMultiplier -= 1
					}
					if warlock.Incinerate != nil {
						warlock.Incinerate.CastTimeMultiplier -= 1
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.Immolate {
						spell.CastTimeMultiplier += 1
					}
					if warlock.Shadowflame != nil {
						warlock.Shadowflame.CastTimeMultiplier += 1
					}
					if warlock.Incinerate != nil {
						warlock.Incinerate.CastTimeMultiplier += 1
					}
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if slices.Contains(affectedSpellCodes, spell.SpellCode) && spell.CurCast.CastTime == 0 {
						aura.Deactivate(sim)
					}
				},
			})

			warlock.RegisterAura(core.Aura{
				Label:    "Fire Trance Trigger",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if (spell.SpellCode == SpellCode_WarlockImmolate || spell.SpellCode == SpellCode_WarlockShadowflame) && sim.Proc(.04, "Fire Trance") {
						fireTranceAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetWickedFelheart = core.NewItemSet(core.ItemSet{
	Name: "Wicked Felheart",
	Bonuses: map[int32]core.ApplyEffect{
		// Banish is now instant cast, and can be cast on yourself while you are a Demon. You cannot Banish yourself while you have Forbearance, and doing so will give you Forbearance for 1 min.
		2: func(agent core.Agent) {
			// TODO: Banish not implemented
		},
		// Each time you take damage, you and your pet gain mana equal to the damage taken, up to a maximum of 420 mana per event. Can only occur once every few seconds.
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			actionID := core.ActionID{SpellID: 457572}
			icd := core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Millisecond * 3500,
			}
			manaMetrics := warlock.NewManaMetrics(actionID)
			for _, pet := range warlock.BasePets {
				pet.T1Tank4PManaMetrics = pet.NewManaMetrics(actionID)
			}
			warlock.RegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Warlock - Tank 4P Bonus",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if icd.IsReady(sim) {
						restoreAmount := min(result.Damage, 420)
						warlock.AddMana(sim, restoreAmount, manaMetrics)
						if warlock.ActivePet != nil {
							warlock.ActivePet.AddMana(sim, restoreAmount, warlock.ActivePet.T1Tank4PManaMetrics)
						}
					}
				},
			})
		},
		// Your Shadow Cleave hits have a 20% chance to grant you a Soul Shard, reset the cooldown on Soul Fire, and make your next Soul Fire within 10 sec instant.
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			procAura := warlock.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457643},
				Label:    "Soul Fire!",
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.SoulFire {
						spell.CastTimeMultiplier -= 1
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.SoulFire {
						spell.CastTimeMultiplier += 1
					}
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.SpellCode == SpellCode_WarlockSoulFire {
						aura.Deactivate(sim)
					}
				},
			})

			warlock.RegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Warlock - Tank 6P Bonus",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && spell.SpellCode == SpellCode_WarlockShadowCleave {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})
