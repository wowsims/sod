package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warrior *Warrior) ApplyRunes() {
	// Head
	warrior.applyVigilance()
	warrior.applyEndlessRage()
	warrior.applyShieldMastery()
	warrior.applyTasteForBlood()

	// Shoulders
	warrior.applyShoulderRuneEffect()

	// Cloak
	warrior.applySuddenDeath()
	warrior.applyFreshMeat()
	warrior.registerShockwaveSpell()

	// Chest
	warrior.applyFlagellation()
	warrior.registerRagingBlow()
	warrior.applyBloodFrenzy()

	// Bracers
	warrior.registerRampage()
	warrior.applySwordAndBoard()
	warrior.applyWreckingCrew()

	// Gloves
	warrior.applySingleMindedFury()
	warrior.registerQuickStrike()

	// Waist
	warrior.applyFocusedRage()
	warrior.applyBloodSurge()

	// Pants
	warrior.applyFrenziedAssault()
	warrior.applyConsumedByRage()
	// Furious Thunder implemented in thunder_clap.go

	// Boots
	// Gladiator implemented on stances.go
}

func (warrior *Warrior) applyShoulderRuneEffect() {
	if warrior.Equipment.Shoulders().Rune == int32(proto.WarriorRune_WarriorRuneNone) {
		return
	}

	switch warrior.Equipment.Shoulders().Rune {
	// Damage
	case int32(proto.WarriorRune_RuneShouldersTactician):
		warrior.applyT1Damage2PBonus()
	case int32(proto.WarriorRune_RuneShouldersWarVeteran):
		warrior.applyT1Damage4PBonus()
	case int32(proto.WarriorRune_RuneShouldersBattleForecaster):
		warrior.applyT1Damage6PBonus()
	case int32(proto.WarriorRune_RuneShouldersBloodseeker):
		warrior.applyT2Damage2PBonus()
	case int32(proto.WarriorRune_RuneShouldersTitan):
		warrior.applyT2Damage4PBonus()
	case int32(proto.WarriorRune_RuneShouldersDestroyer):
		warrior.applyT2Damage6PBonus()
	case int32(proto.WarriorRune_RuneShouldersDeathbound):
		warrior.applyTAQDamage2PBonus()
	case int32(proto.WarriorRune_RuneShouldersSanguinist):
		warrior.applyTAQDamage4PBonus()

	// Tank
	case int32(proto.WarriorRune_RuneShouldersSavage):
		warrior.applyT1Tank4PBonus()
	case int32(proto.WarriorRune_RuneShouldersEnmityWarrior):
		warrior.applyT1Tank6PBonus()
	case int32(proto.WarriorRune_RuneShouldersDeflective):
		warrior.applyT2Protection2PBonus()
	case int32(proto.WarriorRune_RuneShouldersRevenger):
		warrior.applyT2Protection4PBonus()
	case int32(proto.WarriorRune_RuneShouldersIncessant):
		warrior.applyT2Protection6PBonus()
	case int32(proto.WarriorRune_RuneShouldersThunderbringer):
		warrior.applyTAQTank2PBonus()
	case int32(proto.WarriorRune_RuneShouldersSentinel):
		warrior.applyTAQTank4PBonus()
	case int32(proto.WarriorRune_RuneShouldersAftershock):
		warrior.applyRAQTank3PBonus()

	// Gladiator
	case int32(proto.WarriorRune_RuneShouldersSouthpaw):
		warrior.applyZGGladiator3PBonus()
	case int32(proto.WarriorRune_RuneShouldersGladiator):
		warrior.applyZGGladiator5PBonus()
	}
}

func (warrior *Warrior) applyVigilance() {
	if !warrior.HasRune(proto.WarriorRune_RuneVigilance) {
		return
	}

	warrior.PseudoStats.ThreatMultiplier *= 1.1
}

func (warrior *Warrior) applyEndlessRage() {
	if !warrior.HasRune(proto.WarriorRune_RuneEndlessRage) {
		return
	}

	warrior.AddDamageDealtRageMultiplier(1.25)
}

func (warrior *Warrior) applyShieldMastery() {
	if !warrior.HasRune(proto.WarriorRune_RuneShieldMastery) {
		return
	}

	buffAura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Mastery Buff",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.1
		},
	})

	// Use a hidden aura to periodically verify that the player is still using a 2H weapon mid-sim, for example if Item Swapping
	warrior.RegisterAura(core.Aura{
		Label:    "Shield Mastery Dummy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:          time.Second * 2,
				TickImmediately: true,
				OnAction: func(sim *core.Simulation) {
					if !warrior.PseudoStats.CanBlock {
						buffAura.Deactivate(sim)
						return
					}

					buffAura.Activate(sim)
				},
			})
		},
	})
}

// You gain Rage from Physical damage taken as if you were wearing no armor.
func (warrior *Warrior) applyFlagellation() {
	if !warrior.HasRune(proto.WarriorRune_RuneFlagellation) {
		return
	}

	// TODO: Rage gain from hits
}

func (warrior *Warrior) applyBloodFrenzy() {
	if !warrior.HasRune(proto.WarriorRune_RuneBloodFrenzy) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Blood Frenzy Dummy",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Rend.StanceMask |= BerserkerStance
		},
	}))
}

func (warrior *Warrior) applyFrenziedAssault() {
	if !warrior.HasRune(proto.WarriorRune_RuneFrenziedAssault) {
		return
	}

	actionID := core.ActionID{SpellID: 431046}
	rageMetrics := warrior.NewRageMetrics(actionID)

	buffAura := warrior.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Frenzied Assault",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMeleeSpeed(sim, 1.40)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMeleeSpeed(sim, 1/1.40)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskWhiteHit) && result.Landed() {
				warrior.AddRage(sim, core.Ternary(result.DidCrit(), 4.0, 2.0), rageMetrics)
			}
		},
	})

	// Use a dummy aura to periodically verify that the player is still using a 2H weapon mid-sim, for example if Item Swapping
	warrior.RegisterAura(core.Aura{
		Label:    "Frenzied Assault Dummy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:          time.Second * 2,
				TickImmediately: true,
				OnAction: func(sim *core.Simulation) {
					if warrior.MainHand().HandType != proto.HandType_HandTypeTwoHand {
						buffAura.Deactivate(sim)
						return
					}

					if !buffAura.IsActive() {
						buffAura.Activate(sim)
					}
				},
			})
		},
	})
}

// Enrages you (activating abilities which require being Enraged) for 12 sec  after you exceed 60 Rage.
// In addition, Whirlwind also strikes with off-hand melee weapons while you are Enraged
func (warrior *Warrior) applyConsumedByRage() {
	if !warrior.HasRune(proto.WarriorRune_RuneConsumedByRage) {
		return
	}

	warrior.ConsumedByRageAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Enrage (Consumed by Rage)",
		ActionID: core.ActionID{SpellID: 425415},
		Duration: time.Second * 12,
	})

	warrior.ConsumedByRageAura.NewExclusiveEffect("Enrage", true, core.ExclusiveEffect{Priority: 0})

	warrior.RegisterAura(core.Aura{
		Label:    "Consumed By Rage Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnRageChange: func(aura *core.Aura, sim *core.Simulation, metrics *core.ResourceMetrics) {
			// Refunding rage should not enable CBR
			if warrior.CurrentRage() < 60 || metrics.ActionID.OtherID == proto.OtherAction_OtherActionRefund {
				return
			}

			warrior.ConsumedByRageAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) applyFocusedRage() {
	if !warrior.HasRune(proto.WarriorRune_RuneFocusedRage) {
		return
	}

	warrior.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagOffensive) && spell.Cost != nil {
			spell.Cost.FlatModifier -= 3
		}
	})
}

func (warrior *Warrior) applyBloodSurge() {
	if !warrior.HasRune(proto.WarriorRune_RuneBloodSurge) {
		return
	}

	warrior.BloodSurgeAura = warrior.RegisterAura(core.Aura{
		Label:    "Blood Surge Proc",
		ActionID: core.ActionID{SpellID: 413399},
		Duration: time.Second * 15,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// removed even if slam doesn't land
			if spell.Matches(ClassSpellMask_WarriorSlamMH) {
				aura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  ClassSpellMask_WarriorSlam,
		FloatValue: -1,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: ClassSpellMask_WarriorSlam,
		IntValue:  -100,
	})

	warrior.bloodSurgeClassMask = ClassSpellMask_WarriorHeroicStrike | ClassSpellMask_WarriorWhirlwind | ClassSpellMask_WarriorBloodthirst | ClassSpellMask_WarriorQuickStrike

	procTrigger := core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:       "Blood Surge",
		ProcChance: 0.3,
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Can be updated dynamically by DPS 2pT4
			if spell.Matches(warrior.bloodSurgeClassMask) {
				warrior.BloodSurgeAura.Activate(sim)
			}
		},
	})

	procTrigger.OnInit = func(aura *core.Aura, sim *core.Simulation) {
		if warrior.Slam == nil {
			aura.Deactivate(sim)
			return
		}
	}

}

func (warrior *Warrior) applyTasteForBlood() {
	if !warrior.HasRune(proto.WarriorRune_RuneTasteForBlood) {
		return
	}

	warrior.TasteForBloodAura = warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 426969},
		Label:    "Taste for Blood",
		Duration: time.Second * 9,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_WarriorOverpower) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           "Taste for Blood Trigger",
		ClassSpellMask: ClassSpellMask_WarriorRend,
		ICD:            time.Millisecond * 5800,
		Callback:       core.CallbackOnPeriodicDamageDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.TasteForBloodAura.Activate(sim)
		},
	})
}

// Your melee hits have a 10% chance to grant Sudden Death. Sudden Death allows one use of Execute regardless of the target's health state.
// When Execute is enabled by Sudden Death, you will retain 10 rage after using Execute.
func (warrior *Warrior) applySuddenDeath() {
	if !warrior.HasRune(proto.WarriorRune_RuneSuddenDeath) {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: int32(proto.WarriorRune_RuneSuddenDeath)})

	minRageKept := 10.0
	procChance := 0.10

	warrior.SuddenDeathAura = warrior.RegisterAura(core.Aura{
		Label:    "Sudden Death",
		ActionID: core.ActionID{SpellID: 440114},
		Duration: time.Second * 10,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.Matches(ClassSpellMask_WarriorExecute) { // removed only when landed
				warrior.AddRage(sim, minRageKept, rageMetrics)
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:       "Sudden Death Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.SuddenDeathAura.Activate(sim)
		},
	})
}

// Damaging a target with Bloodthirst has a 100% chance the first time and a 10% chance each subsequent time to
// Enrage you (activating abilities which requiring being Enraged), and cause you to deal 10% increased Physical damage for 12 sec.
func (warrior *Warrior) applyFreshMeat() {
	if !warrior.HasRune(proto.WarriorRune_RuneFreshMeat) {
		return
	}

	warrior.FreshMeatEnrageAura = warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 14201},
		Label:    "Enrage (Fresh Meat)",
		Duration: time.Second * 12,
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical], 1.1,
	)

	var damagedUnits map[int32]bool
	affectedSpellClassMasks := ClassSpellMask_WarriorBloodthirst | ClassSpellMask_WarriorMortalStrike | ClassSpellMask_WarriorShieldSlam

	warrior.RegisterAura(core.Aura{
		Label:    "Fresh Meat Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			damagedUnits = make(map[int32]bool)
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(affectedSpellClassMasks) {
				return
			}

			procChance := 0.10
			if !damagedUnits[result.Target.UnitIndex] {
				procChance = 1.00
				damagedUnits[result.Target.UnitIndex] = true
			}

			if sim.Proc(procChance, "Fresh Meat") {
				warrior.FreshMeatEnrageAura.Activate(sim)
			}
		},
	})
}

// Your melee critical hits Enrage you (activating abilities which require being Enraged), and increase Mortal Strike, Bloodthirst, and Shield Slam damage by 10% for 12 sec.
func (warrior *Warrior) applyWreckingCrew() {
	if !warrior.HasRune(proto.WarriorRune_RuneWreckingCrew) {
		return
	}

	warrior.WreckingCrewEnrageAura = warrior.RegisterAura(core.Aura{
		Label:    "Enrage (Wrecking Crew)",
		ActionID: core.ActionID{SpellID: 427066},
		Duration: time.Second * 6,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  ClassSpellMask_WarriorMortalStrike | ClassSpellMask_WarriorBloodthirst | ClassSpellMask_WarriorShieldSlam,
		FloatValue: 1.1,
	})

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:     "Wrecking Crew Trigger",
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeCrit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.WreckingCrewEnrageAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) applySwordAndBoard() {
	if !warrior.HasRune(proto.WarriorRune_RuneSwordAndBoard) || !warrior.Talents.ShieldSlam {
		return
	}

	sabAura := warrior.RegisterAura(core.Aura{
		Label:    "Sword And Board",
		ActionID: core.ActionID{SpellID: int32(proto.WarriorRune_RuneSwordAndBoard)},
		Duration: 5 * time.Second,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_WarriorShieldSlam) {
				aura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: ClassSpellMask_WarriorShieldSlam,
		IntValue:  -100,
	})

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           "Sword And Board Trigger",
		ClassSpellMask: ClassSpellMask_WarriorRevenge | ClassSpellMask_WarriorDevastate,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ProcChance:     0.3,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			sabAura.Activate(sim)
			warrior.ShieldSlam.CD.Reset()
		},
	})
}

// While dual-wielding, your movement speed is increased by 10% and you gain 3% attack speed each time your melee auto-attack strikes the same target as your previous auto-attack, stacking up to 5 times.
// Lasts 10 sec or until your auto-attack strikes a different target.
func (warrior *Warrior) applySingleMindedFury() {
	if !warrior.HasRune(proto.WarriorRune_RuneSingleMindedFury) {
		return
	}

	buffAura := warrior.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 461470},
		Label:     "Single-Minded Fury Attack Speed",
		Duration:  time.Second * 10,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			warrior.MultiplyAttackSpeed(sim, 1/(1+0.04*float64(oldStacks)))
			warrior.MultiplyAttackSpeed(sim, 1+0.04*float64(newStacks))
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Single-Minded Fury Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			warrior.lastMeleeAutoTarget = nil
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !warrior.AutoAttacks.IsDualWielding || !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if warrior.lastMeleeAutoTarget == nil || warrior.lastMeleeAutoTarget != result.Target {
				warrior.lastMeleeAutoTarget = result.Target
				buffAura.Deactivate(sim)
				return
			}

			buffAura.Activate(sim)
			buffAura.AddStack(sim)
		},
	})
}
