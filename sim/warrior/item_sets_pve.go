package warrior

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetBattlegearOfValor = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Heroism",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Chance on melee attack to heal you for 88 to 132 and energize you for 10 Rage
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450587}
			healthMetrics := c.NewHealthMetrics(core.ActionID{SpellID: 450589})
			rageMetrics := c.NewRageMetrics(core.ActionID{SpellID: 450589})

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID: actionID,
				Name:     "S03 - Warrior Armor Heal Trigger - Battlegear of Valor",
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
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

var ItemSetUnstoppableMight = core.NewItemSet(core.ItemSet{
	Name: "Unstoppable Might",
	Bonuses: map[int32]core.ApplyEffect{
		// After changing stances, your next offensive ability's rage cost is reduced by 10.
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()

			var affectedSpells []*core.Spell
			tacticianAura := warrior.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 464241},
				Label:    "Tactician",
				Duration: time.Second * 10,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warrior.Spellbook {
						if spell.Cost != nil && spell.Cost.CostType() == core.CostTypeRage && !spell.Flags.Matches(core.SpellFlagHelpful) {
							affectedSpells = append(affectedSpells, spell)
						}
					}
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range affectedSpells {
						spell.Cost.FlatModifier -= 10
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range affectedSpells {
						spell.Cost.FlatModifier += 10
					}
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if slices.Contains(affectedSpells, spell) {
						aura.Deactivate(sim)
					}
				},
			})

			core.MakePermanent(warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Warrior - Damage 2P Bonus Trigger",
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if slices.Contains(StanceCodes, spell.SpellCode) {
						tacticianAura.Activate(sim)
					}
				},
			}))
		},
		// For 5 sec after leaving a stance, you can use abilities requiring that stance as if you were still in that stance.
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()

			battleStanceAura := warrior.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457706},
				Label:    "Echoes of Battle Stance",
				Duration: time.Second * 5,
			})
			defStanceAura := warrior.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457699},
				Label:    "Echoes of Defensive Stance",
				Duration: time.Second * 5,
			})
			berserkStanceAura := warrior.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457708},
				Label:    "Echoes of Berserker Stance",
				Duration: time.Second * 5,
			})
			gladStanceAura := warrior.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457819},
				Label:    "Echoes of Gladiator Stance",
				Duration: time.Second * 5,
			})

			// We're assuming these will be exclusive but TBD
			warrior.newStanceOverrideExclusiveEffect(BattleStance, battleStanceAura)
			warrior.newStanceOverrideExclusiveEffect(DefensiveStance, defStanceAura)
			warrior.newStanceOverrideExclusiveEffect(BerserkerStance, berserkStanceAura)
			warrior.newStanceOverrideExclusiveEffect(AnyStance, gladStanceAura)

			core.MakePermanent(warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Warrior - Damage 4P Bonus Trigger",
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if slices.Contains(StanceCodes, spell.SpellCode) {
						switch warrior.PreviousStance {
						case BattleStance:
							battleStanceAura.Activate(sim)
						case DefensiveStance:
							defStanceAura.Activate(sim)
						case BerserkerStance:
							berserkStanceAura.Activate(sim)
						case GladiatorStance:
							gladStanceAura.Activate(sim)
						}
					}
				},
			}))
		},
		// For the first 10 sec after activating a stance, you can gain an additional benefit:
		// Battle Stance/Gladiator Stance: 10% increased damage done.
		// Berserker Stance: 10% increased critical strike chance.
		// Defensive Stance: 10% reduced Physical damage taken.
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()

			battleAura := warrior.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457816},
				Label:    "Battle Forecast",
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.10
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.10
				},
			})
			defenseAura := warrior.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457814},
				Label:    "Defense Forecast",
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warrior.PseudoStats.DamageTakenMultiplier *= 0.90
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warrior.PseudoStats.DamageTakenMultiplier /= 0.90
				},
			})
			berserkAura := warrior.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457817},
				Label:    "Berserker Forecast",
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warrior.AddStatDynamic(sim, stats.MeleeCrit, 10*core.CritRatingPerCritChance)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warrior.AddStatDynamic(sim, stats.MeleeCrit, -10*core.CritRatingPerCritChance)
				},
			})

			core.MakePermanent(warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Warrior - Damage 6P Bonus Trigger",
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					switch spell.SpellCode {
					case SpellCode_WarriorStanceBattle:
						battleAura.Activate(sim)
					case SpellCode_WarriorStanceGladiator:
						battleAura.Activate(sim)
					case SpellCode_WarriorStanceDefensive:
						defenseAura.Activate(sim)
					case SpellCode_WarriorStanceBerserker:
						berserkAura.Activate(sim)
					}
				},
			}))
		},
	},
})

var ItemSetImmoveableMight = core.NewItemSet(core.ItemSet{
	Name: "Immoveable Might",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases the block value of your shield by 30.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.BlockValue, 30)
		},
		// You gain 1 extra Rage every time you take any damage or deal auto attack damage.
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.AddDamageDealtRageBonus(1)
			warrior.AddDamageTakenRageBonus(1)
		},
		// Increases all threat you generate in Defensive Stance by an additional 10% and increases all damage you deal in Gladiator Stance by 4%.
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			core.MakePermanent(warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T1 - Warrior - Tank 6P Bonus",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warrior.defensiveStanceThreatMultiplier *= 1.10
					warrior.gladiatorStanceDamageMultiplier *= 1.04
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warrior.defensiveStanceThreatMultiplier /= 1.10
					warrior.gladiatorStanceDamageMultiplier /= 1.04
				},
			}))
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 5 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetUnstoppableWrath = core.NewItemSet(core.ItemSet{
	Name: "Unstoppable Wrath",
	Bonuses: map[int32]core.ApplyEffect{
		// Overpower critical strikes refresh the duration of Rend on your target back to its maximum duration.
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			core.MakePermanent(warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Warrior - Damage 2P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_WarriorOverpower && result.DidCrit() {
						if dot := warrior.Rend.Dot(result.Target); dot.IsActive() {
							dot.Refresh(sim)
						}
					}
				},
			}))
		},
		// Your Whirlwind deals 10% more damage and can be used in all stances.
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Warrior - Damage 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					if warrior.Whirlwind != nil {
						warrior.Whirlwind.DamageMultiplier *= 1.10
						warrior.Whirlwind.StanceMask = AnyStance
					}
				},
			})
		},
		// Your Slam hits reset the remaining cooldown on your Mortal Strike, Bloodthirst, and Shield Slam abilities.
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			var affectedSpells []*WarriorSpell
			core.MakePermanent(warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Warrior - Damage 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range []*WarriorSpell{warrior.Bloodthirst, warrior.MortalStrike, warrior.ShieldSlam} {
						if spell != nil {
							affectedSpells = append(affectedSpells, spell)
						}
					}
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_WarriorSlamMH && result.Landed() {
						for _, spell := range affectedSpells {
							spell.CD.Reset()
						}
					}
				},
			}))
		},
	},
})

var ItemSetImmoveableWrath = core.NewItemSet(core.ItemSet{
	Name: "Immoveable Wrath",
	Bonuses: map[int32]core.ApplyEffect{
		// You gain 10 Rage every time you Parry or one of your attacks is Parried.
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()

			actionID := core.ActionID{SpellID: 468066}
			rageMetrics := warrior.NewRageMetrics(actionID)

			core.MakePermanent(warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Warrior - Protection 2P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ProcMask.Matches(core.ProcMaskMelee) && result.DidParry() {
						warrior.AddRage(sim, 10, rageMetrics)
					}
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidParry() {
						warrior.AddRage(sim, 10, rageMetrics)
					}
				},
			}))
		},
		// Revenge also grants you Flurry, increasing your attack speed by 30% for the next 3 swings.
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()

			flurryAura := warrior.makeFlurryAura(5)
			// The consumption trigger may not exist if the Shaman doesn't talent into Flurry
			warrior.makeFlurryConsumptionTrigger(flurryAura)

			core.MakePermanent(warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Warrior - Protection 4P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_WarriorRevenge && result.Landed() {
						flurryAura.Activate(sim)
						flurryAura.SetStacks(sim, 3)
					}
				},
			}))
		},
		// When your target Parries an attack, you instantly Retaliate for 200% weapon damage to that target.
		// Retaliate cannot be Dodged, Blocked, or Parried, but can only occur once every 30 sec per target.
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()

			retaliate := warrior.RegisterSpell(AnyStance, core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 468071},
				SpellSchool: core.SpellSchoolPhysical,
				DefenseType: core.DefenseTypeMelee,
				ProcMask:    core.ProcMaskMeleeMHSpecial, // Retaliate and Retaliation count as normal yellow hits that can proc things
				Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, warrior.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()), spell.OutcomeMeleeSpecialNoBlockDodgeParry)
				},
			})

			icds := warrior.NewEnemyICDArray(func(u *core.Unit) *core.Cooldown {
				return &core.Cooldown{
					Timer:    warrior.NewTimer(),
					Duration: time.Second * 30,
				}
			})

			core.MakePermanent(warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Warrior - Protection 6P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.DidParry() {
						return
					}

					if icd := icds.Get(result.Target); icd.IsReady(sim) {
						retaliate.Cast(sim, result.Target)
						icd.Use(sim)
					}
				},
			}))
		},
	},
})

var ItemSetVindicatorsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Vindicator's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +7.
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.AddStat(stats.Defense, 7)
		},
		// Reduces the cooldown on your Shield Slam ability by 2 sec.
		3: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			if !warrior.Talents.ShieldSlam {
				return
			}

			warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - ZG - Warrior - Gladiator 3P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					warrior.ShieldSlam.CD.Duration -= time.Second * 2
				},
			})
		},
		// Reduces the cooldown on your Bloodrage ability by 30 sec while you are in Gladiator Stance.
		5: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			if !warrior.HasRune(proto.WarriorRune_RuneGladiatorStance) {
				return
			}

			warrior.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Warrior - Protection 6P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					idx := slices.IndexFunc(warrior.GladiatorStanceAura.ExclusiveEffects, func(ee *core.ExclusiveEffect) bool {
						return ee.Category.Name == stanceEffectCategory
					})
					ee := warrior.GladiatorStanceAura.ExclusiveEffects[idx]
					oldOnGain := ee.OnGain
					ee.OnGain = func(ee *core.ExclusiveEffect, sim *core.Simulation) {
						oldOnGain(ee, sim)
						warrior.Bloodrage.CD.Duration -= time.Second * 30
					}

					oldOnExpire := ee.OnExpire
					ee.OnExpire = func(ee *core.ExclusiveEffect, sim *core.Simulation) {
						oldOnExpire(ee, sim)
						warrior.Bloodrage.CD.Duration += time.Second * 30
					}
				},
			})
		},
	},
})
