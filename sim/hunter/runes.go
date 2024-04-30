package hunter

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (hunter *Hunter) ApplyRunes() {
	if hunter.HasRune(proto.HunterRune_RuneChestHeartOfTheLion) {
		statMultiply := 1.1
		hunter.MultiplyStat(stats.Strength, statMultiply)
		hunter.MultiplyStat(stats.Stamina, statMultiply)
		hunter.MultiplyStat(stats.Agility, statMultiply)
		hunter.MultiplyStat(stats.Intellect, statMultiply)
		hunter.MultiplyStat(stats.Spirit, statMultiply)
	}

	if hunter.HasRune(proto.HunterRune_RuneChestMasterMarksman) {
		hunter.AddStat(stats.MeleeCrit, 5*core.CritRatingPerCritChance)
		hunter.AddStat(stats.SpellCrit, 5*core.SpellCritRatingPerCritChance)
	}

	if hunter.HasRune(proto.HunterRune_RuneChestLoneWolf) && hunter.pet == nil {
		hunter.PseudoStats.DamageDealtMultiplier *= 1.3
	}

	if hunter.HasRune(proto.HunterRune_RuneHandsBeastmastery) && hunter.pet != nil {
		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.1 // Changed to 10% per April 30 nerf
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
	hunter.applyInvigoration()
	hunter.applyLockAndLoad()
	hunter.applyRaptorFury()
}

func (hunter *Hunter) applyInvigoration() {
	if !hunter.HasRune(proto.HunterRune_RuneBootsInvigoration) || hunter.pet == nil {
		return
	}

	procSpellId := core.ActionID{SpellID: 437999}
	metrics := hunter.NewManaMetrics(procSpellId)
	procSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID:    procSpellId,
		SpellSchool: core.SpellSchoolNature,
		ApplyEffects: func(sim *core.Simulation, u *core.Unit, spell *core.Spell) {
			hunter.AddMana(sim, hunter.MaxMana()*0.05, metrics)
		},
	})

	core.MakePermanent(hunter.pet.GetOrRegisterAura(core.Aura{
		Label: "Invigoration",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
				return
			}

			if !result.DidCrit() {
				return
			}

			procSpell.Cast(sim, result.Target)
		},
	}))
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
				if spell.ProcMask.Matches(core.ProcMaskRangedSpecial) {
					spell.BonusCritRating += statDelta * 2 * core.CritRatingPerCritChance
				}
				// Chimera - Serpent double dips this bonus and has ProcMaskEmpty so just add 20 here
				if spell.ActionID.SpellID == 409493 {
					spell.BonusCritRating += statDelta * 4 * core.CritRatingPerCritChance
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
			if spell.ProcMask.Matches(core.ProcMaskRangedSpecial) && spell.Flags.Matches(core.SpellFlagMeleeMetrics) {
				aura.Deactivate(sim)
				hunter.AddMana(sim, spell.CurCast.Cost, lockAndLoadMetrics)

				if spell.CD.Timer != nil {
					spell.CD.Reset()
				}
			}
		},
	})
}

func (hunter *Hunter) applyRaptorFury() {
	if !hunter.HasRune(proto.HunterRune_RuneBracersRaptorFury) {
		return
	}

	hunter.RaptorFuryAura = hunter.GetOrRegisterAura(core.Aura{
		Label:     "Raptor Fury Buff",
		ActionID:  core.ActionID{SpellID: int32(proto.HunterRune_RuneBracersRaptorFury)},
		Duration:  time.Second * 15,
		MaxStacks: 5,
	})
}

func (hunter *Hunter) tntDamageMultiplier() float64 {
	if hunter.HasRune(proto.HunterRune_RuneBracersTNT) {
		return 1.1
	}
	return 1.0
}
