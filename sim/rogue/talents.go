package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (rogue *Rogue) ApplyTalents() {
	rogue.applyMurder()
	rogue.applySealFate()
	rogue.applyWeaponSpecializations()
	rogue.applyWeaponExpertise()
	rogue.applyInitiative()

	rogue.AddStat(stats.Dodge, core.DodgeRatingPerDodgeChance*1*float64(rogue.Talents.LightningReflexes))
	rogue.AddStat(stats.Parry, core.ParryRatingPerParryChance*1*float64(rogue.Talents.Deflection))
	rogue.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(rogue.Talents.Malice))
	rogue.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*1*float64(rogue.Talents.Precision))
	rogue.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*1*float64(rogue.Talents.Precision))
	// TODO: Test the Armor reduction amount
	rogue.AddStat(stats.ArmorPenetration, float64(5/3*rogue.Talents.SerratedBlades*rogue.Level))
	rogue.AutoAttacks.OHConfig().DamageMultiplier *= rogue.dwsMultiplier()

	if rogue.Talents.Deadliness > 0 {
		rogue.MultiplyStat(stats.AttackPower, 1.0+0.02*float64(rogue.Talents.Deadliness))
	}

	rogue.registerColdBloodCD()
	rogue.registerBladeFlurryCD()
	rogue.registerAdrenalineRushCD()
	rogue.registerPreparationCD()
	rogue.registerPremeditation()
	rogue.registerGhostlyStrikeSpell()
}

// dwsMultiplier returns the offhand damage multiplier
func (rogue *Rogue) dwsMultiplier() float64 {
	return 1 + 0.1*float64(rogue.Talents.DualWieldSpecialization)
}

func (rogue *Rogue) makeFinishingMoveEffectApplier() func(sim *core.Simulation, numPoints int32) {
	ruthlessnessMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 14161})
	relentlessStrikesMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 14179})

	return func(sim *core.Simulation, numPoints int32) {
		if t := rogue.Talents.Ruthlessness; t > 0 {
			if sim.RandomFloat("Ruthlessness") < 0.2*float64(t) {
				rogue.AddComboPoints(sim, 1, ruthlessnessMetrics)
			}
		}
		if t := rogue.Talents.RelentlessStrikes; t {
			if sim.RandomFloat("RelentlessStrikes") < 0.2*float64(numPoints) {
				rogue.AddEnergy(sim, 25, relentlessStrikesMetrics)
			}
		}
	}
}

// Murder talent
func (rogue *Rogue) applyMurder() {
	if rogue.Talents.Murder == 0 {
		return
	}

	multiplier := []float64{1, 1.01, 1.02}[rogue.Talents.Murder]

	// TODO Murder, Monster Slaying, Humanoid Slaying, and Beast Slaying (Troll) all affect critical strike damage as well

	// post finalize, since attack tables need to be setup
	rogue.Env.RegisterPostFinalizeEffect(func() {
		for _, t := range rogue.Env.Encounter.Targets {
			switch t.MobType {
			case proto.MobType_MobTypeHumanoid, proto.MobType_MobTypeGiant, proto.MobType_MobTypeBeast, proto.MobType_MobTypeDragonkin:
				for _, at := range rogue.AttackTables[t.UnitIndex] {
					at.DamageDealtMultiplier *= multiplier
				}
			}
		}
	})
}

// Cold Blood talent
func (rogue *Rogue) registerColdBloodCD() {
	if !rogue.Talents.ColdBlood {
		return
	}

	actionID := core.ActionID{SpellID: 14177}

	coldBloodAura := rogue.RegisterAura(core.Aura{
		Label:    "Cold Blood",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range rogue.Spellbook {
				if spell.Flags.Matches(SpellFlagColdBlooded) {
					spell.BonusCritRating += 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range rogue.Spellbook {
				if spell.Flags.Matches(SpellFlagColdBlooded) {
					spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// deactivate after use, but for MutilateMH, so MutilateOH is cold-blooded as well
			if spell.Flags.Matches(SpellFlagColdBlooded) && spell != rogue.MutilateMH {
				aura.Deactivate(sim)
			}
		},
	})

	rogue.ColdBlood = rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			coldBloodAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.ColdBlood,
		Type:  core.CooldownTypeDPS,
	})
}

// Seal Fate talent
func (rogue *Rogue) applySealFate() {
	if rogue.Talents.SealFate == 0 {
		return
	}

	procChance := 0.2 * float64(rogue.Talents.SealFate)
	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 14195})

	icd := core.Cooldown{
		Timer:    rogue.NewTimer(),
		Duration: 500 * time.Millisecond,
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Seal Fate",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagBuilder) {
				return
			}

			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if icd.IsReady(sim) && sim.Proc(procChance, "Seal Fate") {
				rogue.AddComboPoints(sim, 1, cpMetrics)
				icd.Use(sim)
			}
		},
	})
}

// Initiative talent
func (rogue *Rogue) applyInitiative() {
	if rogue.Talents.Initiative == 0 {
		return
	}

	procChance := 0.25 * float64(rogue.Talents.Initiative)
	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 13980})

	rogue.RegisterAura(core.Aura{
		Label:    "Initiative",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == rogue.Garrote || spell == rogue.Ambush {
				if result.Landed() {
					if sim.Proc(procChance, "Initiative") {
						rogue.AddComboPoints(sim, 1, cpMetrics)
					}
				}
			}
		},
	})
}

// Rogue weapon specialization talents. Bonus is shown if the main hand is specialized, but not if off hand only
func (rogue *Rogue) applyWeaponSpecializations() {
	// Sword specialization. Implemented in 'sword_specialization.go'
	if swordSpec := rogue.Talents.SwordSpecialization; swordSpec > 0 {
		if mask := rogue.GetProcMaskForTypes(proto.WeaponType_WeaponTypeSword); mask != core.ProcMaskUnknown {
			rogue.registerSwordSpecialization(mask)
		}
	}

	// Dagger Specialization
	if daggerSpec := rogue.Talents.DaggerSpecialization; daggerSpec > 0 {
		switch rogue.GetProcMaskForTypes(proto.WeaponType_WeaponTypeDagger) {
		case core.ProcMaskMelee:
			rogue.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*float64(daggerSpec))
		case core.ProcMaskMeleeMH:
			// the default character pane displays critical strike chance for main hand only
			rogue.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*float64(daggerSpec))
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating -= core.CritRatingPerCritChance * float64(daggerSpec)
				}
			})
		case core.ProcMaskMeleeOH:
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating += core.CritRatingPerCritChance * float64(daggerSpec)
				}
			})
		}
	}

	// Fist Weapon Specialization. Same as above but for fists
	if fistSpec := rogue.Talents.FistWeaponSpecialization; fistSpec > 0 {
		switch rogue.GetProcMaskForTypes(proto.WeaponType_WeaponTypeFist) {
		case core.ProcMaskMelee:
			rogue.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*float64(fistSpec))
		case core.ProcMaskMeleeMH:
			// the default character pane displays critical strike chance for main hand only
			rogue.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*float64(fistSpec))
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating -= core.CritRatingPerCritChance * float64(fistSpec)
				}
			})
		case core.ProcMaskMeleeOH:
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating += core.CritRatingPerCritChance * float64(fistSpec)
				}
			})
		}
	}

	// Mace Specialization. Offers weapon skill for Maces and RNG stun (not implemented for being useless on boss)
	if maceSpec := rogue.Talents.MaceSpecialization; maceSpec > 0 {
		if mask := rogue.GetProcMaskForTypes(proto.WeaponType_WeaponTypeMace); mask != core.ProcMaskUnknown {
			rogue.PseudoStats.MacesSkill += float64(maceSpec)
		}
	}
}

func (rogue *Rogue) applyWeaponExpertise() {
	if wepExpertise := rogue.Talents.WeaponExpertise; wepExpertise > 0 {
		wepBonus := []float64{0, 3, 5}
		rogue.PseudoStats.SwordsSkill += wepBonus[wepExpertise]
		rogue.PseudoStats.DaggersSkill += wepBonus[wepExpertise]
		rogue.PseudoStats.UnarmedSkill += wepBonus[wepExpertise]
	}
}

/** Wrath Energy gain talents
func (rogue *Rogue) applyCombatPotency() {
	if rogue.Talents.CombatPotency == 0 {
		return
	}

	const procChance = 0.2
	energyBonus := 3.0 * float64(rogue.Talents.CombatPotency)
	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 35553})

	rogue.RegisterAura(core.Aura{
		Label:    "Combat Potency",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// from 3.0.3 patch notes: "Combat Potency: Now only works with auto attacks"
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOHAuto) {
				return
			}

			if sim.RandomFloat("Combat Potency") < procChance {
				rogue.AddEnergy(sim, energyBonus, energyMetrics)
			}
		},
	})
}

func (rogue *Rogue) applyFocusedAttacks() {
	if rogue.Talents.FocusedAttacks == 0 {
		return
	}

	procChance := []float64{0, 0.33, 0.66, 1}[rogue.Talents.FocusedAttacks]
	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 51637})

	rogue.RegisterAura(core.Aura{
		Label:    "Focused Attacks",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.DidCrit() {
				return
			}
			// Fan of Knives OH hits do not trigger focused attacks
			if spell.ProcMask.Matches(core.ProcMaskMeleeOH) && spell.IsSpellAction(FanOfKnivesSpellID) {
				return
			}
			if sim.Proc(procChance, "Focused Attacks") {
				rogue.AddEnergy(sim, 2, energyMetrics)
			}
		},
	})
}*/

var BladeFlurryActionID = core.ActionID{SpellID: 13877}
var BladeFlurryHitID = core.ActionID{SpellID: 22482}

func (rogue *Rogue) registerBladeFlurryCD() {
	if !rogue.Talents.BladeFlurry {
		return
	}

	// TODO verify that this double dips from damage modifiers

	var curDmg float64
	bfHit := rogue.RegisterSpell(core.SpellConfig{
		ActionID:    BladeFlurryHitID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty, // No proc mask, so it won't proc itself.
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, curDmg, spell.OutcomeAlwaysHit)
		},
	})

	rogue.BladeFlurryAura = rogue.RegisterAura(core.Aura{
		Label:    "Blade Flurry",
		ActionID: BladeFlurryActionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, 1.2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, 1/1.2)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.GetNumTargets() < 2 {
				return
			}
			if result.Damage == 0 || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			// Undo armor reduction to get the raw damage value.
			curDmg = result.Damage / result.ResistanceMultiplier

			bfHit.Cast(sim, rogue.Env.NextTargetUnit(result.Target))
			bfHit.SpellMetrics[result.Target.UnitIndex].Casts--
		},
	})

	cooldownDur := time.Minute * 2
	rogue.BladeFlurry = rogue.RegisterSpell(core.SpellConfig{
		ActionID: BladeFlurryActionID,
		Flags:    core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: cooldownDur,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			rogue.BladeFlurryAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.BladeFlurry,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if sim.GetRemainingDuration() > cooldownDur+time.Second*15 {
				// We'll have enough time to cast another BF, so use it immediately to make sure we get the 2nd one.
				return true
			}

			// Since this is our last BF, wait until we have SND / procs up.
			sndTimeRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)
			return sndTimeRemaining >= time.Second
		},
	})
}

var AdrenalineRushActionID = core.ActionID{SpellID: 13750}

func (rogue *Rogue) registerAdrenalineRushCD() {
	if !rogue.Talents.AdrenalineRush {
		return
	}

	rogue.AdrenalineRushAura = rogue.RegisterAura(core.Aura{
		Label:    "Adrenaline Rush",
		ActionID: AdrenalineRushActionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ResetEnergyTick(sim)
			rogue.ApplyEnergyTickMultiplier(1.0)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ResetEnergyTick(sim)
			rogue.ApplyEnergyTickMultiplier(-1.0)
		},
	})

	adrenalineRushSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID: AdrenalineRushActionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			rogue.AdrenalineRushAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    adrenalineRushSpell,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityBloodlust,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return rogue.CurrentEnergy() <= 45.0
		},
	})
}

/** Honor Among Thieves (Possible P2 rune)
func (rogue *Rogue) registerHonorAmongThieves() {
	// When anyone in your group critically hits with a damage or healing spell or ability,
	// you have a [33%/66%/100%] chance to gain a combo point on your current target.
	// This effect cannot occur more than once per second.
	if rogue.Talents.HonorAmongThieves == 0 {
		return
	}

	procChance := []float64{0, 0.33, 0.66, 1}[rogue.Talents.HonorAmongThieves]
	comboMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 51701})
	honorAmongThievesID := core.ActionID{SpellID: 51701}

	icd := core.Cooldown{
		Timer:    rogue.NewTimer(),
		Duration: time.Second,
	}

	maybeProc := func(sim *core.Simulation) {
		if icd.IsReady(sim) && sim.Proc(procChance, "honor of thieves") {
			rogue.AddComboPoints(sim, 1, comboMetrics)
			icd.Use(sim)
		}
	}

	rogue.HonorAmongThieves = rogue.RegisterAura(core.Aura{
		Label:    "Honor Among Thieves",
		ActionID: honorAmongThievesID,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			// In an ideal party, you'd probably get up to 6 ability crits/s (Rate = 600).
			//  Survival Hunters, Enhancement Shamans, and Assassination Rogues are particularly good.
			if rogue.Options.HonorOfThievesCritRate <= 0 {
				return
			}

			if rogue.Options.HonorOfThievesCritRate > 2000 {
				rogue.Options.HonorOfThievesCritRate = 2000 // limited, so performance doesn't suffer
			}

			rateToDuration := float64(time.Second) * 100 / float64(rogue.Options.HonorOfThievesCritRate)

			pa := &core.PendingAction{}
			pa.OnAction = func(sim *core.Simulation) {
				maybeProc(sim)
				pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomExpFloat("next party crit")*rateToDuration)
				sim.AddPendingAction(pa)
			}
			pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomExpFloat("next party crit")*rateToDuration)
			sim.AddPendingAction(pa)
		},
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() && !spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto|core.ProcMaskMeleeOHAuto|core.ProcMaskRangedAuto) {
				maybeProc(sim)
			}
		},
		OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
				maybeProc(sim)
			}
		},
	})
}*/
