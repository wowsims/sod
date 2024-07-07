package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var _ = core.NewItemSet(core.ItemSet{
	Name: "Blood Corrupted Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()

			procAuras := rogue.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
				return target.RegisterAura(core.Aura{
					Label:     "Blood Corruption",
					ActionID:  core.ActionID{SpellID: 449927},
					Duration:  time.Second * 15,
					MaxStacks: 30,

					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						for si := stats.SchoolIndexPhysical; si < stats.SchoolLen; si++ {
							aura.Unit.PseudoStats.SchoolBonusDamageTaken[si] += 7
						}
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						for si := stats.SchoolIndexPhysical; si < stats.SchoolLen; si++ {
							aura.Unit.PseudoStats.SchoolBonusDamageTaken[si] -= 7
						}
					},
					OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if result.Landed() && spell.ProcMask.Matches(core.ProcMaskDirect) {
							aura.RemoveStack(sim)
						}
					},
				})
			})

			core.MakePermanent(rogue.RegisterAura(core.Aura{
				Label:    "Blood Corrupting" + rogue.Label,
				ActionID: core.ActionID{SpellID: 449928},
				OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRangedSpecial) {
						return
					}

					switch spell {
					case rogue.Backstab, rogue.Mutilate, rogue.SinisterStrike, rogue.SaberSlash, rogue.Shiv, rogue.PoisonedKnife, rogue.MainGauche, rogue.QuickDraw:
						aura := procAuras.Get(result.Target)
						aura.Activate(sim)
						aura.SetStacks(sim, aura.MaxStacks)
					}
				},
			}))
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetDarkmantleArmor = core.NewItemSet(core.ItemSet{
	Name: "Darkmantle Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Chance on melee attack to restore 35 energy.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 27787}
			energyMetrics := c.NewEnergyMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID: actionID,
				Name:     "Rogue Armor Energize",
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
				ProcMask: core.ProcMaskMeleeWhiteHit,
				PPM:      1,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasEnergyBar() {
						c.AddEnergy(sim, 35, energyMetrics)
					}
				},
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

var ItemSetNightSlayerThrill = core.NewItemSet(core.ItemSet{
	Name: "Nightslayer Thrill",
	Bonuses: map[int32]core.ApplyEffect{
		// Feint also grants Avoidance for 6 sec, reducing all damage taken from area of effect attacks from non-players by 50%
		//2: Not yet implemented
		//},
		// Increases the critical strike damage bonus of your Poisons by 100%.
		4: func(agent core.Agent) {
			rogue := agent.GetCharacter()
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.Flags.Matches(SpellFlagRoguePoison) {
					spell.CritDamageBonus += 1
				}
			})
		},
		// Your finishing moves have a 5% chance per combo point to make your next ability cost no energy.
		//https://www.wowhead.com/classic/spell=457342/clearcasting
		6: 	func(agent core.Agent) {		
			rogue := agent.(RogueAgent).GetRogue()			
				
			aura := rogue.RegisterAura(core.Aura{
				Label:    "Clearcasting (S03 - Item - T1 - Rogue - Damage 6P Bonus)",
				ActionID: core.ActionID{SpellID: 457342},
				Duration: time.Second * 15,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.CostMultiplier -= 1
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.CostMultiplier += 1
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if aura.RemainingDuration(sim) == aura.Duration {
						return
					}
					if spell.DefaultCast.Cost == 0 {
						return
					}
					aura.Deactivate(sim)
				},
			})

			rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
				if sim.Proc(.05*float64(comboPoints), "Clearcasting (S03 - Item - T1 - Rogue - Damage 6P Bonus)") {
					aura.Activate(sim)
				}
			})
		},
	},
})