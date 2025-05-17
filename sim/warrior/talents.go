package warrior

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warrior *Warrior) ToughnessArmorMultiplier() float64 {
	return 1.0 + 0.02*float64(warrior.Talents.Toughness)
}

func (warrior *Warrior) ApplyTalents() {
	warrior.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(warrior.Talents.Cruelty))
	warrior.ApplyEquipScaling(stats.Armor, warrior.ToughnessArmorMultiplier())
	warrior.AddStat(stats.Defense, 2*float64(warrior.Talents.Anticipation))
	warrior.AddStat(stats.Parry, 1*float64(warrior.Talents.Deflection))

	warrior.applyAngerManagement()
	warrior.applyDeepWounds()
	warrior.applyOneHandedWeaponSpecialization()
	warrior.applyTwoHandedWeaponSpecialization()
	warrior.applyWeaponSpecializations()
	warrior.applyUnbridledWrath()
	warrior.applyDualWieldSpecialization()
	warrior.applyEnrage()
	warrior.applyFlurry()
	warrior.applyShieldSpecialization()
	warrior.registerDeathWishCD()
	warrior.registerSweepingStrikesCD()
	warrior.registerLastStandCD()
}

func (warrior *Warrior) applyAngerManagement() {
	if !warrior.Talents.AngerManagement {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12296})

	warrior.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Second * 3,
			OnAction: func(sim *core.Simulation) {
				warrior.AddRage(sim, 1, rageMetrics)
				warrior.LastAMTick = sim.CurrentTime
			},
		})
	})
}

func (warrior *Warrior) applyTwoHandedWeaponSpecialization() {
	if warrior.Talents.TwoHandedWeaponSpecialization == 0 || warrior.MainHand().HandType != proto.HandType_HandTypeTwoHand {
		return
	}

	multiplier := 1 + 0.01*float64(warrior.Talents.TwoHandedWeaponSpecialization)
	warrior.OnSpellRegistered(func(spell *core.Spell) {
		if spell.BonusCoefficient > 0 {
			spell.ApplyMultiplicativeDamageBonus(multiplier)
		}
	})
}

func (warrior *Warrior) applyOneHandedWeaponSpecialization() {
	if warrior.Talents.OneHandedWeaponSpecialization == 0 || warrior.MainHand().HandType == proto.HandType_HandTypeTwoHand {
		return
	}

	multiplier := 1 + 0.02*float64(warrior.Talents.OneHandedWeaponSpecialization)
	warrior.OnSpellRegistered(func(spell *core.Spell) {
		if spell.BonusCoefficient > 0 {
			spell.ApplyMultiplicativeDamageBonus(multiplier)
		}
	})
}

func (warrior *Warrior) applyWeaponSpecializations() {
	if ss := warrior.Talents.SwordSpecialization; ss > 0 {
		if mask := warrior.GetProcMaskForTypes(proto.WeaponType_WeaponTypeSword); mask != core.ProcMaskUnknown {
			warrior.registerSwordSpecialization(mask)
		}
	}

	if as := warrior.Talents.AxeSpecialization; as > 0 {
		// the default character panel displays critical strike chance for main hand only
		switch warrior.GetProcMaskForTypes(proto.WeaponType_WeaponTypeAxe) {
		case core.ProcMaskMelee:
			warrior.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance*float64(as))
		case core.ProcMaskMeleeMH:
			warrior.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance*float64(as))
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating -= 1 * core.CritRatingPerCritChance * float64(as)
				}
			})
		case core.ProcMaskMeleeOH:
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating += 1 * core.CritRatingPerCritChance * float64(as)
				}
			})
		}
	}

	if ps := warrior.Talents.PolearmSpecialization; ps > 0 {
		// the default character panel displays critical strike chance for main hand only
		switch warrior.GetProcMaskForTypes(proto.WeaponType_WeaponTypePolearm) {
		case core.ProcMaskMelee:
			warrior.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance*float64(ps))
		case core.ProcMaskMeleeMH:
			warrior.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance*float64(ps))
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating -= 1 * core.CritRatingPerCritChance * float64(ps)
				}
			})
		case core.ProcMaskMeleeOH:
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating += 1 * core.CritRatingPerCritChance * float64(ps)
				}
			})
		}
	}

}

func (warrior *Warrior) registerSwordSpecialization(procMask core.ProcMask) {
	icd := core.Cooldown{
		Timer:    warrior.NewTimer(),
		Duration: time.Millisecond * 200,
	}
	procChance := 0.01 * float64(warrior.Talents.SwordSpecialization)

	warrior.RegisterAura(core.Aura{
		Label:    "Sword Specialization",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}
			if !spell.ProcMask.Matches(procMask) {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			if sim.RandomFloat("Sword Specialization") < procChance {
				icd.Use(sim)
				warrior.AutoAttacks.ExtraMHAttack(sim, 1, core.ActionID{SpellID: 12815}, spell.ActionID)
			}
		},
	})
}

func (warrior *Warrior) applyUnbridledWrath() {
	if warrior.Talents.UnbridledWrath == 0 {
		return
	}

	procChance := 0.08 * float64(warrior.Talents.UnbridledWrath)

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12964})

	warrior.RegisterAura(core.Aura{
		Label:    "Unbridled Wrath",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && sim.RandomFloat("Unbrided Wrath") < procChance {
				warrior.AddRage(sim, 1, rageMetrics)
			}
		},
	})
}

func (warrior *Warrior) applyDualWieldSpecialization() {
	if warrior.Talents.DualWieldSpecialization == 0 {
		return
	}

	multiplier := 1 + 0.05*float64(warrior.Talents.DualWieldSpecialization)

	warrior.AutoAttacks.OHConfig().DamageMultiplier *= multiplier
}

func (warrior *Warrior) applyEnrage() {
	if warrior.Talents.Enrage == 0 {
		return
	}

	warrior.EnrageAura = warrior.GetOrRegisterAura(core.Aura{
		Label:     "Enrage",
		ActionID:  core.ActionID{SpellID: 13048},
		Duration:  time.Second * 12,
		MaxStacks: 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + 0.05*float64(warrior.Talents.Enrage)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1 + 0.05*float64(warrior.Talents.Enrage)
		},
	})

	warrior.EnrageAura.NewExclusiveEffect("Enrage", true, core.ExclusiveEffect{Priority: 5 * float64(warrior.Talents.Enrage)})

	warrior.RegisterAura(core.Aura{
		Label:    "Enrage Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !warrior.EnrageAura.IsActive() {
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMelee) {
				warrior.EnrageAura.RemoveStack(sim)
			}
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			warrior.EnrageAura.Activate(sim)
			if warrior.EnrageAura.IsActive() {
				warrior.EnrageAura.SetStacks(sim, 12)
			}
		},
	})
}

// func (warrior *Warrior) applyFlurry() {
// 	if warrior.Talents.Flurry == 0 {
// 		return
// 	}

// 	haste := []float64{1, 1.1, 1.15, 1.2, 1.25, 1.3}[warrior.Talents.Flurry]

// 	procAura := warrior.RegisterAura(core.Aura{
// 		Label:     "Flurry Proc",
// 		ActionID:  core.ActionID{SpellID: 12974},
// 		Duration:  core.NeverExpires,
// 		MaxStacks: 3,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			warrior.MultiplyMeleeSpeed(sim, haste)
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			warrior.MultiplyMeleeSpeed(sim, 1/haste)
// 		},
// 	})

// 	warrior.RegisterAura(core.Aura{
// 		Label:    "Flurry",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
// 				return
// 			}

// 			if result.Outcome.Matches(core.OutcomeCrit) {
// 				procAura.Activate(sim)
// 				procAura.SetStacks(sim, 3)
// 				return
// 			}

// 			// Remove a stack.
// 			if procAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
// 				procAura.RemoveStack(sim)
// 			}
// 		},
// 	})
// }

func (warrior *Warrior) applyFlurry() {
	if warrior.Talents.Flurry == 0 {
		return
	}

	talentAura := warrior.makeFlurryAura(warrior.Talents.Flurry)

	// This must be registered before the below trigger because in-game a crit weapon swing consumes a stack before the refresh, so you end up with:
	// 3 => 2
	// refresh
	// 2 => 3
	warrior.makeFlurryConsumptionTrigger(talentAura)

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Flurry Proc Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee) && result.Outcome.Matches(core.OutcomeCrit) {
				talentAura.Activate(sim)
				if talentAura.IsActive() {
					talentAura.SetStacks(sim, 3)
				}
				return
			}
		},
	}))
}

// These are separated out because of the T1 Shaman Tank 2P that can proc Flurry separately from the talent.
// It triggers the max-rank Flurry aura but with dodge, parry, or block.
func (warrior *Warrior) makeFlurryAura(points int32) *core.Aura {
	if points == 0 {
		return nil
	}

	spellID := []int32{12319, 12971, 12972, 12973, 12974}[points-1]
	attackSpeed := []float64{1.1, 1.15, 1.2, 1.25, 1.3}[points-1]
	label := fmt.Sprintf("Flurry Proc (%d)", spellID)

	if aura := warrior.GetAura(label); aura != nil {
		return aura
	}

	aura := warrior.GetOrRegisterAura(core.Aura{
		Label:     label,
		ActionID:  core.ActionID{SpellID: spellID},
		Duration:  core.NeverExpires,
		MaxStacks: 3,
	})

	aura.NewExclusiveEffect("Flurry", true, core.ExclusiveEffect{
		Priority: attackSpeed,
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			warrior.MultiplyMeleeSpeed(sim, attackSpeed)
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			warrior.MultiplyMeleeSpeed(sim, 1/attackSpeed)
		},
	})

	return aura
}

// With the Protection T2 4pc it's possible to have 2 different Flurry auras if using less than 5/5 points in Flurry.
// The two different buffs don't stack whatsoever. Instead the stronger aura takes precedence and each one is only refreshed by the corresponding triggers.
func (warrior *Warrior) makeFlurryConsumptionTrigger(flurryAura *core.Aura) *core.Aura {
	label := fmt.Sprintf("Flurry Consume Trigger - %d", flurryAura.ActionID.SpellID)
	if aura := warrior.GetAura(label); aura != nil {
		return aura
	}

	icd := core.Cooldown{
		Timer:    warrior.NewTimer(),
		Duration: time.Millisecond * 500,
	}

	return core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Remove a stack.
			if flurryAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && icd.IsReady(sim) {
				icd.Use(sim)
				flurryAura.RemoveStack(sim)
			}
		},
	}))
}

func (warrior *Warrior) applyShieldSpecialization() {
	if warrior.Talents.ShieldSpecialization == 0 {
		return
	}

	warrior.AddStat(stats.Block, core.BlockRatingPerBlockChance*1*float64(warrior.Talents.ShieldSpecialization))

	procChance := 0.2 * float64(warrior.Talents.ShieldSpecialization)
	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12727})

	warrior.RegisterAura(core.Aura{
		Label:    "Shield Specialization",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidBlock() {
				if sim.Proc(procChance, "Shield Specialization") {
					warrior.AddRage(sim, 1.0, rageMetrics)
				}
			}
		},
	})
}

func (warrior *Warrior) registerDeathWishCD() {
	if !warrior.Talents.DeathWish {
		return
	}

	actionID := core.ActionID{SpellID: 12328}

	deathWishAura := warrior.RegisterAura(core.Aura{
		Label:    "Death Wish",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.2
			warrior.PseudoStats.ArmorMultiplier *= 0.8
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.2
			warrior.PseudoStats.ArmorMultiplier /= 0.8
		},
	})
	core.RegisterPercentDamageModifierEffect(deathWishAura, 1.2)

	warrior.DeathWish = warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_WarriorDeathWish,
		Flags:          core.SpellFlagHelpful,
		RageCost: core.RageCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			IgnoreHaste: true,
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			deathWishAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.DeathWish.Spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (warrior *Warrior) registerLastStandCD() {
	if !warrior.Talents.LastStand {
		return
	}

	actionID := core.ActionID{SpellID: 12975}
	healthMetrics := warrior.NewHealthMetrics(actionID)

	var bonusHealth float64
	warrior.LastStandAura = warrior.RegisterAura(core.Aura{
		Label:    "Last Stand",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = warrior.MaxHealth() * 0.3
			warrior.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})
			warrior.GainHealth(sim, bonusHealth, healthMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})
		},
	})

	lastStandSpell := warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			warrior.LastStandAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: lastStandSpell.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}

func (warrior *Warrior) impale() float64 {
	return 0.1 * float64(warrior.Talents.Impale)
}
