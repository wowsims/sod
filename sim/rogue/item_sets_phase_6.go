package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var ItemSetEmblemsofVeiledShadows = core.NewItemSet(core.ItemSet{
	Name: "Emblems of Veiled Shadows",
	Bonuses: map[int32]core.ApplyEffect{
		// 3 pieces: Your finishing moves cost 50% less Energy.
		3: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.RegisterAura(core.Aura{
				Label: "S03 - Item - RAQ - Rogue - Damage 3P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, finisher := range rogue.Finishers {
						finisher.Cost.Multiplier -= 50
					}
				},
			})
		},
	},
})

var ItemSetDeathdealersThrill = core.NewItemSet(core.ItemSet{
	Name: "Deathdealer's Thrill",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases Saber Slash and Sinister Strike damage by 20%
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Rogue - Damage 2P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					if rogue.SaberSlash != nil {
						rogue.SaberSlash.DamageMultiplierAdditive += 0.20
					}
					if rogue.Mutilate != nil {
						rogue.MutilateMH.DamageMultiplierAdditive += 0.20
						rogue.MutilateOH.DamageMultiplierAdditive += 0.20
					}
				},
			})
		},
		// Reduces the cooldown on Adrenaline Rush by 4 minutes.
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			if !rogue.Talents.AdrenalineRush {
				return
			}
			rogue.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Rogue - Damage 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					rogue.AdrenalineRush.CD.Duration -= time.Second * 240
				},
			})
		},
	},
})

var ItemSetDeathdealersBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Deathdealer's Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Main Gauche now strikes 1 additional nearby target and also causes your Sinister Strike to strike 1 additional nearby target.
		// These additional strikes are not duplicated by Blade Flurry.
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			if !rogue.HasRune(proto.RogueRune_RuneMainGauche) {
				return
			}

			if rogue.Env.GetNumTargets() == 1 {
				return
			}

			var curDmg float64

			cleaveHit := rogue.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 1213754},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, curDmg, spell.OutcomeAlwaysHit)
				},
			})

			cleaveAura := rogue.RegisterAura(core.Aura{
				Label:    "2P Cleave Buff",
				Duration: time.Second * 10,
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && (spell.SpellCode == SpellCode_RogueSinisterStrike) {
						curDmg = result.Damage / result.ResistanceMultiplier
						cleaveHit.Cast(sim, rogue.Env.NextTargetUnit(result.Target))
						cleaveHit.SpellMetrics[result.Target.UnitIndex].Casts--
					}
				},
			})

			core.MakePermanent(rogue.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Rogue - Tank 2P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && spell.SpellCode == SpellCode_RogueMainGauche {
						cleaveAura.Activate(sim)
						curDmg = result.Damage / result.ResistanceMultiplier
						cleaveHit.Cast(sim, rogue.Env.NextTargetUnit(result.Target))
						cleaveHit.SpellMetrics[result.Target.UnitIndex].Casts--
					}
				},
			}))

		},
		// While active, your Main Gauche also causes you to heal for 10% of all damage done by Sinister Strike.
		// Any excess healing becomes a Blood Barrier, absorbing damage up to 20% of your maximum health.
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			if !rogue.HasRune(proto.RogueRune_RuneMainGauche) {
				return
			}
			healthMetrics := rogue.NewHealthMetrics(core.ActionID{SpellID: 11294})
			healAmount := 0.0
			shieldAmount := 0.0
			currentShield := 0.0

			var shieldSpell *core.Spell

			shieldSpell = rogue.GetOrRegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 1213761},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskSpellHealing,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Shield: core.ShieldConfig{
					SelfOnly: true,
					Aura: core.Aura{
						Label:    "Blood Barrier",
						ActionID: core.ActionID{SpellID: 1213762},
						Duration: time.Second * 15,
						OnReset: func(aura *core.Aura, sim *core.Simulation) {
							shieldAmount = 0.0
							currentShield = 0.0
						},
						OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
							if currentShield <= 0 || result.Damage <= 0 {
								return
							}

							damageReduced := min(result.Damage, currentShield)
							currentShield -= damageReduced

							rogue.GainHealth(sim, damageReduced, shieldSpell.HealthMetrics(result.Target))
							if currentShield <= 0 {
								shieldSpell.SelfShield().Deactivate(sim)
							}
						},
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					if currentShield < rogue.MaxHealth()*0.2 {
						shieldAmount = min(shieldAmount, rogue.MaxHealth()*0.2-currentShield)
						currentShield += shieldAmount
						spell.SelfShield().Apply(sim, shieldAmount)
					}
				},
			})

			activeAura := core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
				Name:     "Main Gauche - Blood Barrier",
				ActionID: core.ActionID{SpellID: 1213762},
				Callback: core.CallbackOnSpellHitDealt,
				Duration: time.Second * 15,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && spell.SpellCode == SpellCode_RogueSinisterStrike {
						healAmount = result.Damage * 0.15
						if rogue.CurrentHealth() < rogue.MaxHealth() {
							rogue.GainHealth(sim, healAmount, healthMetrics)
						} else {
							shieldAmount = healAmount
							shieldSpell.Cast(sim, result.Target)
						}

					}
				},
			})

			core.MakePermanent(rogue.RegisterAura(core.Aura{
				Label: "S03 - Item - TAQ - Rogue - Tank 4P Bonus",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && spell.SpellCode == SpellCode_RogueMainGauche {
						activeAura.Activate(sim)
					}
				},
			}))
		},
	},
})
