package warlock

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
			c.AddResistances(8)
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
		// Your Nightfall talent has a 4% increased chance to trigger.
		// Incinerate has a 4% chance to trigger the Warlock’s Decimation.
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			warlock6pt1Aura := warlock.RegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Warlock - Damage 6P Bonus",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnGain: func(_ *core.Aura, _ *core.Simulation) {
					warlock.nightfallProcChance += 0.04
				},
				OnExpire: func(_ *core.Aura, _ *core.Simulation) {
					warlock.nightfallProcChance -= 0.04
				},
			})

			if !warlock.HasRune(proto.WarlockRune_RuneBracerIncinerate) || !warlock.HasRune(proto.WarlockRune_RuneCloakDecimation) {
				return
			}

			warlock6pt1Aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.SpellCode == SpellCode_WarlockIncinerate && result.Landed() && sim.Proc(.04, "T1 6P Incinerate Proc") {
					warlock.DecimationAura.Activate(sim)
				}
			}
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

			soulFireCastTime := SoulFireCastTime
			procAura := warlock.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457643},
				Label:    "Soul Fire!",
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.SoulFire {
						spell.DefaultCast.CastTime = 0
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.SoulFire {
						spell.DefaultCast.CastTime = soulFireCastTime
					}
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.SpellCode == SpellCode_WarlockSoulFire {
						aura.Deactivate(sim)
					}
				},
			})

			icd := core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Millisecond * 100,
			}

			warlock.RegisterAura(core.Aura{
				Label:    "S03 - Item - T1 - Warlock - Tank 6P Bonus",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
					soulFireCastTime = warlock.SoulFire[0].DefaultCast.CastTime
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && spell.SpellCode == SpellCode_WarlockShadowCleave && icd.IsReady(sim) && sim.Proc(0.2, "Soul Fire! Proc") {
						procAura.Activate(sim)
						icd.Use(sim)
					}
				},
			})
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 5 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetCorruptedNemesis = core.NewItemSet(core.ItemSet{
	Name: "Corrupted Nemesis",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases the damage of your periodic spells and Felguard pet by 10%
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			warlock.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Warlock - Damage 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.Spellbook {
						if spell.Flags.Matches(SpellFlagWarlock) && len(spell.Dots()) > 0 {
							spell.DamageMultiplier *= 1.10
						}
					}

					if warlock.HasRune(proto.WarlockRune_RuneBracerSummonFelguard) {
						warlock.Felguard.PseudoStats.DamageDealtMultiplier *= 1.10
					}
				},
			})
		},
		// Periodic damage from your Shadowflame, Unstable Affliction, and Curse of Agony spells and damage done by your Felguard have a 4% chance to grant the Shadow Trance effect.
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			procChance := 0.04

			affectedSpellCodes := []int32{SpellCode_WarlockCurseOfAgony, SpellCode_WarlockShadowflame, SpellCode_WarlockUnstableAffliction}
			core.MakePermanent(warlock.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Warlock - Damage 4P Bonus",
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if slices.Contains(affectedSpellCodes, spell.SpellCode) && sim.Proc(procChance, "Proc Shadow Trance") {
						warlock.ShadowTranceAura.Activate(sim)
					}
				},
			}))

			if !warlock.HasRune(proto.WarlockRune_RuneBracerSummonFelguard) {
				return
			}

			core.MakePermanent(warlock.Felguard.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Warlock - Damage 4P Bonus - Felguard Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && sim.Proc(procChance, "Proc Shadow Trance") {
						warlock.ShadowTranceAura.Activate(sim)
					}
				},
			}))
		},
		// Shadowbolt deals 10% increased damage for each of your effects afflicting the target, up to a maximum of 30%.
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			warlock.shadowBoltActiveEffectMultiplierPer = .10
			warlock.shadowBoltActiveEffectMultiplierMax = 1.30
		},
	},
})

var ItemSetDemoniacsThreads = core.NewItemSet(core.ItemSet{
	Name: "Demoniac's Threads",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 12.
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.AddStat(stats.SpellPower, 12)
		},
		// Increases the Attack Power and Spell Damage your Demon pet gains from your attributes by 20%.
		3: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			core.MakePermanent(warlock.RegisterAura(core.Aura{
				Label: "S03 - Item - ZG - Warlock - Demonology 3P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, pet := range warlock.BasePets {
						oldStatInheritance := pet.GetStatInheritance()
						pet.UpdateStatInheritance(
							func(ownerStats stats.Stats) stats.Stats {
								s := oldStatInheritance(ownerStats)
								s[stats.AttackPower] *= 1.20
								s[stats.SpellPower] *= 1.20
								return s
							},
						)
					}
				},
			}))
		},
		// Increases the benefits of your Master Demonologist talent by 50%.
		5: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.masterDemonologistBonus += .50
		},
	},
})
