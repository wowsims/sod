package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	DiamondFlask         = 20130
	Exsanguinar          = 216497
	SuzerainDefender     = 224280
	GrileksCharmOFMight  = 231286
	RageOfMugamba        = 231350
	BandOfTheDreadnaught = 236022
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(BandOfTheDreadnaught, func(agent core.Agent) {
		character := agent.GetCharacter()
		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Band of the Dreadnaught Trigger",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskMelee,
			ProcChance: 0.02,
			ICD:        time.Millisecond * 200,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				character.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 1223010}, spell)
			},
		})
	})

	core.NewItemEffect(DiamondFlask, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := character.NewTemporaryStatsAura("Diamond Flask", core.ActionID{SpellID: 24427}, stats.Stats{stats.Strength: 75}, time.Second*60)

		triggerSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 24427},
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 6,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 60,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    triggerSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(Exsanguinar, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionId := core.ActionID{SpellID: 436332}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionId,
			SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			DamageMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Exsanguination",
				},
				TickLength:    2 * time.Second,
				NumberOfTicks: 15,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 5, isRollover)
				},

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					// Has no DefenseType, also haven't seen a miss in logs.
					result := spell.CalcAndDealDamage(sim, aoeTarget, 65, spell.OutcomeAlwaysHit)
					if result.Landed() {
						spell.Dot(aoeTarget).Apply(sim)
					}
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(GrileksCharmOFMight, func(agent core.Agent) {
		warrior := agent.(WarriorAgent).GetWarrior()
		actionID := core.ActionID{ItemID: GrileksCharmOFMight}
		rageMetrics := warrior.NewRageMetrics(actionID)
		duration := time.Second * 20

		aura := warrior.NewTemporaryStatsAura("Gri'lek's Guard", actionID, stats.Stats{stats.BlockValue: 200}, duration)

		spell := warrior.Character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warrior.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    warrior.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				warrior.AddRage(sim, 30, rageMetrics)
				aura.Activate(sim)
			},
		})

		warrior.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	core.NewItemEffect(RageOfMugamba, func(agent core.Agent) {
		warrior := agent.(WarriorAgent).GetWarrior()
		if !warrior.Talents.ShieldSlam {
			return
		}

		core.MakePermanent(warrior.RegisterAura(core.Aura{
			Label: "Reduced Shield Slam Cost (Rage of Mugamba)",
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Flat,
			ClassMask: ClassSpellMask_WarriorShieldSlam,
			IntValue:  -5,
		}))
	})

	core.NewItemEffect(SuzerainDefender, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: SuzerainDefender}

		// Store a reference in case the unit switches targets since we don't have a great way to do this right now
		fightingDragonkin := false
		rageOfSuzerain := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469025},
			Label:    "Rage of the Suzerain",
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.30
					fightingDragonkin = true
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if fightingDragonkin {
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.30
					fightingDragonkin = false
				}
			},
		})

		defenseOfDragonflights := character.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Defense of the Dragonflights",
			Duration: time.Second * 5,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.SchoolDamageTakenMultiplier.MultiplyMagicSchools(0.50)
				rageOfSuzerain.Activate(sim)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.SchoolDamageTakenMultiplier.MultiplyMagicSchools(2)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 1,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				defenseOfDragonflights.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	core.AddEffectsToTest = true
}

// Your Bloodthirst, Mortal Strike, Shield Slam, Heroic Strike, and Cleave critical strikes set the duration of your Rend on the target to 21 sec.
func (warrior *Warrior) ApplyQueensfallWarriorEffect(aura *core.Aura) {
	aura.AttachProcTrigger(core.ProcTrigger{
		Name:     "Queensfall Trigger - Warrior",
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeCrit,
		ClassSpellMask: ClassSpellMask_WarriorBloodthirst | ClassSpellMask_WarriorMortalStrike | ClassSpellMask_WarriorShieldSlam |
			ClassSpellMask_WarriorHeroicStrike | ClassSpellMask_WarriorCleave,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if dot := warrior.Rend.Dot(result.Target); dot.IsActive() {
				dot.NumberOfTicks = int32(21 / dot.TickLength.Seconds())
				dot.RecomputeAuraDuration()
				dot.Rollover(sim)
			}
		},
	})
}

// Striking a higher level enemy applies a stack of Coup, increasing their damage taken from your next Execute by 10% per stack, stacking up to 20 times. At 20 stacks, Execute may be cast regardless of the target's health.
func (warrior *Warrior) ApplyRegicideWarriorEffect(aura *core.Aura) {
	// Coup debuff array
	debuffAuras := warrior.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
		aura := unit.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 1231417},
			Label:     "Coup",
			MaxStacks: 20,
			Duration:  core.NeverExpires,
		})

		return aura
	})

	exeDamageMod := warrior.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: ClassSpellMask_WarriorExecute,
	})

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Coup Execute Damage",
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Coup Trigger - Warrior",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: ClassSpellMask_WarriorExecute,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				debuffAuras[result.Target.Index].SetStacks(sim, 0)
			}
		},
	}))

	// Apply the Coup debuff to the target hit by melee abilities
	aura.AttachProcTrigger(core.ProcTrigger{
		Name:     "Regicide Trigger - Warrior",
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		ProcMask: core.ProcMaskMelee,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			debuffAuras[result.Target.Index].Activate(sim)
			debuffAuras[result.Target.Index].AddStack(sim)
			exeDamageMod.UpdateFloatValue(1 + float64(debuffAuras[result.Target.Index].GetStacks())*0.10)
			exeDamageMod.Activate()
		},
	})
}
