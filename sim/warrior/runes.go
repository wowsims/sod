package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warrior *Warrior) ApplyRunes() {
	// Head
	warrior.applyEndlessRage()
	warrior.applyShieldMastery()
	warrior.applyTasteForBlood()

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
	warrior.applyFuriousThunder()

	// Boots
	// Gladiator implemented on stances.go
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
				Period: time.Second * 2,
				OnAction: func(sim *core.Simulation) {
					if warrior.OffHand().WeaponType != proto.WeaponType_WeaponTypeShield {
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
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Rend.StanceMask |= BerserkerStance
		},
	}))
}

func (warrior *Warrior) applyFrenziedAssault() {
	if !warrior.HasRune(proto.WarriorRune_RuneFrenziedAssault) {
		return
	}

	buffAura := warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 431046},
		Label:    "Frenzied Assault",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMeleeSpeed(sim, 1.2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMeleeSpeed(sim, 1/1.2)
		},
	})

	// Use a dummy aura to periodically verify that the player is still using a 2H weapon mid-sim, for example if Item Swapping
	warrior.RegisterAura(core.Aura{
		Label:    "Frenzied Assault Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 2,
				OnAction: func(sim *core.Simulation) {
					if warrior.MainHand().HandType != proto.HandType_HandTypeTwoHand {
						buffAura.Deactivate(sim)
						return
					}

					buffAura.Activate(sim)
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

func (warrior *Warrior) applyFuriousThunder() {
	if !warrior.HasRune(proto.WarriorRune_RuneFuriousThunder) {
		return
	}

	warrior.ThunderClap.StanceMask = AnyStance
}

func (warrior *Warrior) applyFocusedRage() {
	warrior.FocusedRageDiscount = core.TernaryFloat64(warrior.HasRune(proto.WarriorRune_RuneFocusedRage), 3.0, 0)
}

func (warrior *Warrior) applyBloodSurge() {
	if !warrior.HasRune(proto.WarriorRune_RuneBloodSurge) {
		return
	}

	warrior.BloodSurgeAura = warrior.RegisterAura(core.Aura{
		Label:    "Blood Surge Proc",
		ActionID: core.ActionID{SpellID: 413399},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if warrior.Slam != nil {
				warrior.Slam.DefaultCast.CastTime = 0
				warrior.Slam.CostMultiplier -= 1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if warrior.Slam != nil {
				warrior.Slam.DefaultCast.CastTime = 1500 * time.Millisecond
				warrior.Slam.CostMultiplier += 1
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if warrior.Slam != nil && spell.SpellCode == SpellCode_WarriorSlamOH { // removed even if slam doesn't land
				aura.Deactivate(sim)
			}
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Blood Surge",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if !spell.Flags.Matches(SpellFlagBloodSurge) {
				return
			}

			if sim.Proc(0.3, "Blood Surge") {
				warrior.BloodSurgeAura.Activate(sim)
			}
		},
	})
}

func (warrior *Warrior) applyTasteForBlood() {
	if !warrior.HasRune(proto.WarriorRune_RuneTasteForBlood) {
		return
	}

	icd := core.Cooldown{
		Timer:    warrior.NewTimer(),
		Duration: time.Millisecond * 580,
	}

	warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 426969},
		Label:    "Taste for Blood",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode == SpellCode_WarriorRend && icd.IsReady(sim) {
				icd.Use(sim)
				warrior.OverpowerAura.Duration = time.Second * 9
				warrior.OverpowerAura.Activate(sim)
				warrior.OverpowerAura.Duration = time.Second * 5
			}
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
			if result.Landed() && spell.SpellCode == SpellCode_WarriorExecute { // removed only when landed
				warrior.AddRage(sim, minRageKept, rageMetrics)
				aura.Deactivate(sim)
			}
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Sudden Death Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if sim.Proc(procChance, "Sudden Death") {
				warrior.SuddenDeathAura.Activate(sim)
			}
		},
	})
}

// Damaging a target with Bloodthirst has a 100% chance the first time and a 10% chance each subsequent time to
// Enrage you (activating abilities which requiring being Enraged), and cause you to deal 10% increased Physical damage for 12 sec.
func (warrior *Warrior) applyFreshMeat() {
	if !warrior.HasRune(proto.WarriorRune_RuneFreshMeat) {
		return
	}

	hasBloodthirstTalent := warrior.Talents.Bloodthirst

	damagedUnits := make(map[int32]bool)

	warrior.FreshMeatEnrageAura = warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 14201},
		Label:    "Enrage (Fresh Meat)",
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.1
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Fresh Meat Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			damagedUnits = make(map[int32]bool)
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !hasBloodthirstTalent || spell.SpellCode != SpellCode_WarriorBloodthirst {
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

	affectedSpells := core.FilterSlice(
		[]*WarriorSpell{warrior.MortalStrike, warrior.Bloodthirst, warrior.ShieldSlam},
		func(spell *WarriorSpell) bool { return spell != nil },
	)

	warrior.WreckingCrewEnrageAura = warrior.RegisterAura(core.Aura{
		Label:    "Enrage (Wrecking Crew)",
		ActionID: core.ActionID{SpellID: 427066},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplier *= 1.1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplier /= 1.1
			}
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Wrecking Crew Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			warrior.WreckingCrewEnrageAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) applySwordAndBoard() {
	if !warrior.HasRune(proto.WarriorRune_RuneSwordAndBoard) || !warrior.Talents.ShieldSlam {
		return
	}

	sabAura := warrior.GetOrRegisterAura(core.Aura{
		Label:    "Sword And Board",
		ActionID: core.ActionID{SpellID: int32(proto.WarriorRune_RuneSwordAndBoard)},
		Duration: 5 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.ShieldSlam.CostMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.ShieldSlam.CostMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode == SpellCode_WarriorShieldSlam {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(warrior.GetOrRegisterAura(core.Aura{
		Label: "Sword And Board Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell.SpellCode != SpellCode_WarriorRevenge && spell.SpellCode != SpellCode_WarriorDevastate {
				return
			}

			if sim.Proc(0.3, "Sword And Board") {
				sabAura.Activate(sim)
				warrior.ShieldSlam.CD.Reset()
			}
		},
	}))
}

// While dual-wielding, your movement speed is increased by 10% and you gain 2% attack speed each time your melee auto-attack strikes the same target as your previous auto-attack, stacking up to 5 times.
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
			warrior.MultiplyAttackSpeed(sim, 1/(1+.02*float64(oldStacks)))
			warrior.MultiplyAttackSpeed(sim, 1+.02*float64(newStacks))
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Single-Minded Fury Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
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
