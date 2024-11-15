package hunter

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (hunter *Hunter) ApplyRunes() {
	if hunter.HasRune(proto.HunterRune_RuneChestLoneWolf) && hunter.pet == nil {
		hunter.PseudoStats.DamageDealtMultiplier *= 1.30
	}

	if hunter.HasRune(proto.HunterRune_RuneChestBeastmastery) && hunter.pet != nil {
		// https://www.wowhead.com/classic/news/class-tuning-incoming-hunter-shaman-warlock-season-of-discovery-339072?webhook
		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.1
		core.MakePermanent(hunter.RegisterAura(core.Aura{
			Label: "Beastmastery Rune Focus",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				if hunter.pet != nil {
					hunter.pet.AddFocusRegenMultiplier(1.50)
				}
			},
		}))
	}

	if hunter.HasRune(proto.HunterRune_RuneBootsDualWieldSpecialization) {
		hunter.AutoAttacks.OHConfig().DamageMultiplier *= 1.5
	}

	if hunter.HasRune(proto.HunterRune_RuneHelmCatlikeReflexes) {
		hunter.AddStat(stats.Dodge, 20*core.DodgeRatingPerDodgeChance)
		if hunter.pet != nil {
			hunter.pet.AddStat(stats.Dodge, 9*core.DodgeRatingPerDodgeChance)
		}
	}

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
			if spell.Flags.Matches(SpellFlagShot) {
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
				if spell.ProcMask.Matches(core.ProcMaskRangedSpecial) || spell.SpellCode == SpellCode_HunterChimeraSerpent {
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
			if result.DidCrit() && (spell.Flags.Matches(SpellFlagShot|SpellFlagStrike) || spell.SpellCode == SpellCode_HunterMongooseBite) {
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

	// icd := core.Cooldown{
	// 	Timer:    hunter.NewTimer(),
	// 	Duration: time.Second * 8,
	// }

	hunter.LockAndLoadAura = hunter.GetOrRegisterAura(core.Aura{
		Label:    "Lock And Load",
		ActionID: core.ActionID{SpellID: 415413},
		Duration: time.Second * 20,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Flags.Matches(SpellFlagShot) {
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
			if spell.Flags.Matches(SpellFlagTrap) {
				spell.WaitTravelTime(sim, func(s *core.Simulation) {
					// if icd.IsReady(sim) {
					// 	icd.Use(sim)
					hunter.LockAndLoadAura.Activate(sim)
					// }
				})
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

	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagTrap) || spell.SpellCode == SpellCode_HunterExplosiveShot {
			spell.DamageMultiplier *= 1.10
		}
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
		if spell.Flags.Matches(SpellFlagTrap) {
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

func (hunter *Hunter) applyImprovedVolley() {
	if !hunter.HasRune(proto.HunterRune_RuneCloakImprovedVolley) && hunter.Volley != nil {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label:    "Improved Volley",
		ActionID: core.ActionID{SpellID: 440520},
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// The 3% rAP scaling and manacost reduction is applied inside the volley spell config itself
			hunter.Volley.DamageMultiplier *= 2
		},
	})
}
