package warrior

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
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
		// You gain 10 Rage when you change stances.
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 457652})
			core.MakePermanent(warrior.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 457652}, // Intentionally exposing for stance-dancing APL conditions
				Label:    "S03 - Item - T1 - Warrior - Damage 2P Bonus Trigger",
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if slices.Contains(StanceCodes, spell.SpellCode) {
						warrior.AddRage(sim, 10, rageMetrics)
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
					warrior.PseudoStats.SchoolDamageDealtMultiplier[core.SpellSchoolPhysical] *= 1.10
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warrior.PseudoStats.SchoolDamageDealtMultiplier[core.SpellSchoolPhysical] /= 1.10
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
