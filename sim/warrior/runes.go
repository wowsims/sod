package warrior

import (
	"slices"
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
			warrior.MultiplyMeleeSpeed(sim, 1.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMeleeSpeed(sim, 1/1.3)
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
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Slam.CastTimeMultiplier -= 1
			warrior.Slam.Cost.Multiplier -= 100
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Slam.CastTimeMultiplier += 1
			warrior.Slam.Cost.Multiplier += 100
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// removed even if slam doesn't land
			if spell.SpellCode == SpellCode_WarriorSlamMH {
				aura.Deactivate(sim)
			}
		},
	})

	affectedSpells := make(map[*core.Spell]bool)

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Blood Surge Trigger",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if warrior.Slam == nil {
				aura.Deactivate(sim)
				return
			}

			affectedSpells[warrior.HeroicStrike.Spell] = true
			affectedSpells[warrior.Whirlwind.Spell] = true

			if warrior.Bloodthirst != nil {
				affectedSpells[warrior.Bloodthirst.Spell] = true
			}

			if warrior.HasRune(proto.WarriorRune_RuneQuickStrike) {
				affectedSpells[warrior.QuickStrike.Spell] = true
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && affectedSpells[spell] && sim.Proc(0.3, "Blood Surge") {
				warrior.BloodSurgeAura.Activate(sim)
			}
		},
	}))
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
			if spell.SpellCode == SpellCode_WarriorOverpower {
				aura.Deactivate(sim)
			}
		},
	})

	icd := core.Cooldown{
		Timer:    warrior.NewTimer(),
		Duration: time.Millisecond * 5800,
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Taste for Blood Trigger",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode == SpellCode_WarriorRend && icd.IsReady(sim) {
				icd.Use(sim)
				warrior.TasteForBloodAura.Activate(sim)
			}
		},
	}))
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

	var damagedUnits map[int32]bool
	affectedSpellCodes := []int32{SpellCode_WarriorBloodthirst, SpellCode_WarriorMortalStrike, SpellCode_WarriorShieldSlam}

	warrior.RegisterAura(core.Aura{
		Label:    "Fresh Meat Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			damagedUnits = make(map[int32]bool)
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !slices.Contains(affectedSpellCodes, spell.SpellCode) {
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

	var affectedSpells []*WarriorSpell
	warrior.WreckingCrewEnrageAura = warrior.RegisterAura(core.Aura{
		Label:    "Enrage (Wrecking Crew)",
		ActionID: core.ActionID{SpellID: 427066},
		Duration: time.Second * 6,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(
				[]*WarriorSpell{warrior.MortalStrike, warrior.Bloodthirst, warrior.ShieldSlam},
				func(spell *WarriorSpell) bool { return spell != nil },
			)
		},
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
			if spell.ProcMask.Matches(core.ProcMaskMelee) && result.Outcome.Matches(core.OutcomeCrit) {
				warrior.WreckingCrewEnrageAura.Activate(sim)
			}
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
			warrior.ShieldSlam.Cost.Multiplier -= 100
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.ShieldSlam.Cost.Multiplier += 100
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
