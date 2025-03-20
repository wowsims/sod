package hunter

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (hunter *Hunter) ApplyRunes() {
	hunter.applyShoulderRuneEffect()

	if hunter.HasRune(proto.HunterRune_RuneChestLoneWolf) && hunter.pet == nil {
		hunter.PseudoStats.DamageDealtMultiplier *= 1.30
	}

	if hunter.HasRune(proto.HunterRune_RuneChestBeastmastery) && hunter.pet != nil {
		hunter.pet.PseudoStats.DamageDealtMultiplierAdditive += 0.15

		core.MakePermanent(hunter.RegisterAura(core.Aura{
			Label: "Beastmastery Rune Focus",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				if hunter.pet != nil {
					hunter.pet.AddFocusRegenMultiplier(0.50)
				}
			},
		}))
	}

	if hunter.HasRune(proto.HunterRune_RuneBootsDualWieldSpecialization) {
		hunter.AutoAttacks.OHConfig().DamageMultiplier *= 1.60
	}

	hunter.applyCatlikeReflexes()
	hunter.applySniperTraining()
	hunter.applyCobraStrikes()
	hunter.applyExposeWeakness()
	// hunter.applyInvigoration()
	hunter.applyLockAndLoad()
	hunter.applyRaptorFury()
	hunter.applyCobraSlayer()
	hunter.applyHitAndRun()
	hunter.applyMasterMarksman()
	hunter.applyImprovedVolley()
	hunter.applyTNT()
	hunter.applyResourcefulness()
}

func (hunter *Hunter) applyShoulderRuneEffect() {
	if hunter.Equipment.Shoulders().Rune == int32(proto.HunterRune_HunterRuneNone) {
		return
	}

	switch hunter.Equipment.Shoulders().Rune {
	// Melee
	case int32(proto.HunterRune_RuneShouldersHuntsman):
		hunter.applyT1Melee2PBonus()
	case int32(proto.HunterRune_RuneShouldersRetaliator):
		hunter.applyT1Melee6PBonus()
	case int32(proto.HunterRune_RuneShouldersEchoes):
		hunter.applyT2Melee2PBonus()
	case int32(proto.HunterRune_RuneShouldersLethalLasher):
		hunter.applyT2Melee4PBonus()
	case int32(proto.HunterRune_RuneShouldersKineticist):
		hunter.applyT2Melee6PBonus()
	case int32(proto.HunterRune_RuneShouldersStrategist):
		hunter.applyTAQMelee2PBonus()
	case int32(proto.HunterRune_RuneShouldersDeadlyStriker):
		hunter.applyTAQMelee4PBonus()

	// Ranged
	case int32(proto.HunterRune_RuneShouldersPreyseeker):
		hunter.applyT1Ranged4PBonus()
	case int32(proto.HunterRune_RuneShouldersSharpshooter):
		hunter.applyT1Ranged6PBonus()
	case int32(proto.HunterRune_RuneShouldersHazardHarrier):
		hunter.applyT2Ranged2PBonus()
	case int32(proto.HunterRune_RuneShouldersAlternator):
		hunter.applyT2Ranged4PBonus()
	case int32(proto.HunterRune_RuneShouldersToxinologist):
		hunter.applyT2Ranged6PBonus()
	case int32(proto.HunterRune_RuneShouldersBountyHunter):
		hunter.applyTAQRanged2PBonus()
	case int32(proto.HunterRune_RuneShouldersTrickShooter):
		hunter.applyTAQRanged4PBonus()

	// Beastmaster
	case int32(proto.HunterRune_RuneShouldersBeastTender):
		hunter.applyZGBeastmaster3PBonus()
	case int32(proto.HunterRune_RuneShouldersHoundMaster):
		hunter.applyZGBeastmaster5PBonus()
	case int32(proto.HunterRune_RuneshouldersAlphaTamer):
		hunter.applyRAQBeastmastery5PBonus()
	}
}

// TODO: 2024-06-13 - Rune seemingly replaced with Wyvern Strike
// func (hunter *Hunter) applyInvigoration() {
// 	if !hunter.HasRune(proto.HunterRune_RuneBootsInvigoration) || hunter.pet == nil {
// 		return
// 	}

// 	procSpellId := core.ActionID{SpellID: 437999}
// 	metrics := hunter.NewManaMetrics(procSpellId)
// 	procSpell := hunter.RegisterSpell(core.SpellConfig{
// 		ActionID:    procSpellId,
// 		SpellSchool: core.SpellSchoolNature,
// 		ApplyEffects: func(sim *core.Simulation, u *core.Unit, spell *core.Spell) {
// 			hunter.AddMana(sim, hunter.MaxMana()*0.05, metrics)
// 		},
// 	})

// 	core.MakePermanent(hunter.pet.GetOrRegisterAura(core.Aura{
// 		Label: "Invigoration",
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
// 				return
// 			}

// 			if !result.DidCrit() {
// 				return
// 			}

// 			procSpell.Cast(sim, result.Target)
// 		},
// 	}))
// }

func (hunter *Hunter) applyMasterMarksman() {
	if hunter.HasRune(proto.HunterRune_RuneChestMasterMarksman) {
		hunter.AddStat(stats.MeleeCrit, 5*core.CritRatingPerCritChance)
		hunter.AddStat(stats.SpellCrit, 5*core.SpellCritRatingPerCritChance)

		hunter.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Matches(ClassSpellMask_HunterShots) && spell.Cost != nil {
				spell.Cost.Multiplier -= 25
			}
		})
	}
}

func (hunter *Hunter) applyExposeWeakness() {
	if !hunter.HasRune(proto.HunterRune_RuneBeltExposeWeakness) {
		return
	}

	apBonus := hunter.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 0.4)
	apRangedBonus := hunter.NewDynamicStatDependency(stats.Agility, stats.RangedAttackPower, 0.4)

	procAura := hunter.GetOrRegisterAura(core.Aura{
		Label:    "Expose Weakness Proc",
		ActionID: core.ActionID{SpellID: 409507},
		Duration: time.Second * 7,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.EnableDynamicStatDep(sim, apBonus)
			hunter.EnableDynamicStatDep(sim, apRangedBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.DisableDynamicStatDep(sim, apBonus)
			hunter.DisableDynamicStatDep(sim, apRangedBonus)
		},
	})

	core.MakePermanent(hunter.GetOrRegisterAura(core.Aura{
		Label: "Expose Weakness",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}

			if !result.DidCrit() {
				return
			}

			procAura.Activate(sim)
		},
	}))
}

func (hunter *Hunter) applySniperTraining() {
	if !hunter.HasRune(proto.HunterRune_RuneLegsSniperTraining) {
		return
	}

	hunter.SniperTrainingAura = hunter.GetOrRegisterAura(core.Aura{
		Label:     "Sniper Training",
		ActionID:  core.ActionID{SpellID: 415399},
		Duration:  time.Second * 6,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			statDelta := float64(newStacks - oldStacks)
			for _, spell := range aura.Unit.Spellbook {
				if spell.ProcMask.Matches(core.ProcMaskRangedSpecial) || spell.Matches(ClassSpellMask_HunterChimeraSerpent) {
					spell.BonusCritRating += statDelta * 2 * core.CritRatingPerCritChance
				}

			}
		},
	})

	aura := hunter.SniperTrainingAura
	uptime := hunter.Options.SniperTrainingUptime
	chancePerTick := core.TernaryFloat64(uptime == 1, 1, 1.0-math.Pow(1-uptime, 1))

	lastMoved := false
	aura.Unit.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Second,
			OnAction: func(sim *core.Simulation) {
				if sim.Proc(chancePerTick, "FixedAura") {
					// Gain stack every second after 2 seconds
					if !lastMoved {
						aura.Activate(sim)
						aura.AddStack(sim)
					} else {
						lastMoved = false
					}
				} else {
					// Lose stack every second moving
					if aura.IsActive() {
						aura.RemoveStack(sim)
					}
					lastMoved = true
				}
			},
		})

		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period:   0,
			NumTicks: 1,
			OnAction: func(sim *core.Simulation) {
				if sim.Proc(chancePerTick, "FixedAura") {
					aura.Activate(sim)
					aura.SetStacks(sim, 5)
				} else {
					lastMoved = true
				}
			},
		})
	})
}

func (hunter *Hunter) applyCobraStrikes() {
	if !hunter.HasRune(proto.HunterRune_RuneChestCobraStrikes) || hunter.pet == nil {
		return
	}

	hunter.CobraStrikesAura = hunter.pet.GetOrRegisterAura(core.Aura{
		Label:     "Cobra Strikes",
		ActionID:  core.ActionID{SpellID: 425714},
		Duration:  time.Second * 30,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range aura.Unit.Spellbook {
				if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
					spell.BonusCritRating += 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range aura.Unit.Spellbook {
				if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
					spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				aura.RemoveStack(sim)
			}
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Cobra Strikes Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() && (spell.Matches(ClassSpellMask_HunterShots|ClassSpellMask_HunterMongooseBite) || spell.Flags.Matches(SpellFlagStrike)) {
				hunter.CobraStrikesAura.Activate(sim)
				hunter.CobraStrikesAura.SetStacks(sim, 2)
			}
		},
	}))
}

func (hunter *Hunter) applyLockAndLoad() {
	if !hunter.HasRune(proto.HunterRune_RuneHelmLockAndLoad) {
		return
	}

	lockAndLoadMetrics := hunter.Metrics.NewResourceMetrics(core.ActionID{SpellID: 415413}, proto.ResourceType_ResourceTypeMana)

	hunter.LockAndLoadAura = hunter.GetOrRegisterAura(core.Aura{
		Label:    "Lock And Load",
		ActionID: core.ActionID{SpellID: 415413},
		Duration: time.Second * 20,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_HunterShots) {
				aura.Deactivate(sim)
				hunter.AddMana(sim, spell.CurCast.Cost, lockAndLoadMetrics)

				if spell.CD.Timer != nil {
					spell.CD.Reset()
				}
			}
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Lock And Load Trigger",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_HunterTraps) {
				hunter.LockAndLoadAura.Activate(sim)
			}
		},
	}))
}

const RaptorFuryPerStackDamageMultiplier = 0.15

func (hunter *Hunter) raptorFuryDamageMultiplier() float64 {
	stacks := hunter.RaptorFuryAura.GetStacks()
	if stacks == 0 {
		return 1
	}

	return 1 + RaptorFuryPerStackDamageMultiplier*float64(stacks)
}

func (hunter *Hunter) applyRaptorFury() {
	if !hunter.HasRune(proto.HunterRune_RuneBracersRaptorFury) {
		return
	}

	hunter.RaptorFuryAura = hunter.GetOrRegisterAura(core.Aura{
		Label:     "Raptor Fury Buff",
		ActionID:  core.ActionID{SpellID: int32(proto.HunterRune_RuneBracersRaptorFury)},
		Duration:  time.Second * 30,
		MaxStacks: 5,
	})
}

func (hunter *Hunter) applyCobraSlayer() {
	if !hunter.HasRune(proto.HunterRune_RuneHandsCobraSlayer) {
		return
	}

	hunter.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: int32(proto.HunterRune_RuneHandsCobraSlayer)},
		Label:     "Cobra Slayer",
		Duration:  core.NeverExpires,
		MaxStacks: 20,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if result.DidDodge() {
				aura.SetStacks(sim, 1)
				hunter.DefensiveState.Activate(sim)
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && result.Outcome.Matches(core.OutcomeLanded) && sim.Proc((float64(aura.GetStacks())*0.10), "Cobra Slayer") {
				aura.SetStacks(sim, 1)
				hunter.DefensiveState.Activate(sim)
				return
			}

			aura.AddStack(sim)
		},
	})
}

func (hunter *Hunter) applyTNT() {
	if !hunter.HasRune(proto.HunterRune_RuneBracersTNT) {
		return
	}
	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterExplosiveShot | ClassSpellMask_HunterTraps,
		IntValue:  10,
	})
}

func (hunter *Hunter) tntDamageFlatBonus() float64 {
	if hunter.HasRune(proto.HunterRune_RuneBracersTNT) {
		return math.Max(hunter.GetStat(stats.AttackPower), hunter.GetStat(stats.RangedAttackPower)) * 0.25
	}
	return 0.0
}

func (hunter *Hunter) trapRange() float64 {
	if hunter.HasRune(proto.HunterRune_RuneBootsTrapLauncher) {
		return 35
	}
	return 5
}

func (hunter *Hunter) applyResourcefulness() {
	if !hunter.HasRune(proto.HunterRune_RuneCloakResourcefulness) {
		return
	}

	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_HunterTraps) {
			spell.Cost.BaseCost = 0
			spell.CD.Duration = spell.CD.Duration / 100 * 60
		}
	})
}

func (hunter *Hunter) applyHitAndRun() {
	if hunter.HasRune(proto.HunterRune_RuneCloakHitAndRun) {
		hunter.HitAndRunAura = hunter.RegisterAura(core.Aura{
			Label:    "Hit And Run",
			ActionID: core.ActionID{SpellID: 440533},
			Duration: time.Second * 15,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				hunter.Unit.AddMoveSpeedModifier(&aura.ActionID, 1.30)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				hunter.Unit.RemoveMoveSpeedModifier(&aura.ActionID)
			},
		})
	}
}

func (hunter *Hunter) applyCatlikeReflexes() {
	if !hunter.HasRune(proto.HunterRune_RuneHelmCatlikeReflexes) {
		return
	}
	label := "Catlike Reflexes"

	hunter.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_HunterFlankingStrike,
		Kind:      core.SpellMod_Cooldown_Multi_Flat,
		IntValue:  -50,
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.AddBuildPhaseStatDynamic(sim, stats.Dodge, 20*core.DodgeRatingPerDodgeChance)
			if hunter.pet != nil {
				hunter.pet.AddBuildPhaseStatDynamic(sim, stats.Dodge, 9*core.DodgeRatingPerDodgeChance)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.AddBuildPhaseStatDynamic(sim, stats.Dodge, -20*core.DodgeRatingPerDodgeChance)
			if hunter.pet != nil {
				hunter.pet.AddBuildPhaseStatDynamic(sim, stats.Dodge, -9*core.DodgeRatingPerDodgeChance)
			}
		},
	}))
}

func (hunter *Hunter) applyImprovedVolley() {
	if !hunter.HasRune(proto.HunterRune_RuneCloakImprovedVolley) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label:    "Improved Volley",
		ActionID: core.ActionID{SpellID: 440520},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterVolley,
		// The 3% rAP scaling and manacost reduction is applied inside the volley spell config itself
		IntValue: 100,
	}))
}
