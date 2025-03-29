package rogue

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (rogue *Rogue) ApplyTalents() {
	rogue.applyImprovedSinisterStrike()
	rogue.applyImprovedBackstab()
	rogue.applyImprovedEviscerate()
	rogue.applyRuthlessness()
	rogue.applyMurder()
	rogue.applyVilePoisons()
	rogue.applyRelentlessStrikes()
	rogue.applySealFate()
	rogue.applyWeaponSpecializations()
	rogue.applyWeaponExpertise()
	rogue.applyInitiative()
	rogue.applyAggression()
	rogue.applyOpportunity()
	rogue.applySerratedBlades()

	rogue.AddStat(stats.Dodge, 1*float64(rogue.Talents.LightningReflexes))
	rogue.AddStat(stats.Parry, 1*float64(rogue.Talents.Deflection))
	rogue.AddStat(stats.MeleeCrit, 1*float64(rogue.Talents.Malice))
	rogue.AddStat(stats.MeleeHit, 1*float64(rogue.Talents.Precision))
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
	rogue.applyRiposte()
}

func (rogue *Rogue) applyImprovedEviscerate() {
	if rogue.Talents.ImprovedEviscerate == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_RogueEviscerate,
		Kind:      core.SpellMod_DamageDone_Flat,
		IntValue:  5 * int64(rogue.Talents.ImprovedEviscerate),
	})
}

func (rogue *Rogue) applyImprovedSinisterStrike() {
	if rogue.Talents.ImprovedSinisterStrike == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_RogueSinisterStrikeDependent,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -[]int64{0, 3, 5}[rogue.Talents.ImprovedSinisterStrike],
	})
}

func (rogue *Rogue) applyImprovedBackstab() {
	if rogue.Talents.ImprovedBackstab == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		ClassMask:  ClassSpellMask_RogueBackstabDependent,
		Kind:       core.SpellMod_BonusCrit_Flat,
		FloatValue: 10 * float64(rogue.Talents.ImprovedBackstab),
	})
}

func (rogue *Rogue) applyOpportunity() {
	if rogue.Talents.Opportunity == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_RogueBackstabDependent | ClassSpellMask_RogueAmbush | ClassSpellMask_RogueMutilateHit,
		Kind:      core.SpellMod_ImpactDamageDone_Flat,
		IntValue:  4 * int64(rogue.Talents.Opportunity),
	})

	rogue.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_RogueGarrote,
		Kind:      core.SpellMod_PeriodicDamageDone_Flat,
		IntValue:  4 * int64(rogue.Talents.Opportunity),
	})
}

func (rogue *Rogue) applySerratedBlades() {
	if rogue.Talents.SerratedBlades == 0 {
		return
	}

	// TODO: Test the Armor reduction amount
	rogue.AddStat(stats.ArmorPenetration, float64(5/3*rogue.Talents.SerratedBlades*rogue.Level))

	rogue.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_RogueRuptureDependent,
		Kind:      core.SpellMod_PeriodicDamageDone_Flat,
		IntValue:  10 * int64(rogue.Talents.SerratedBlades),
	})
}

// dwsMultiplier returns the offhand damage multiplier
func (rogue *Rogue) dwsMultiplier() float64 {
	return 1 + 0.1*float64(rogue.Talents.DualWieldSpecialization)
}

func (rogue *Rogue) applyRuthlessness() {
	if rogue.Talents.Ruthlessness == 0 {
		return
	}

	procChance := 0.2 * float64(rogue.Talents.Ruthlessness)
	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 14161})
	rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
		if sim.Proc(procChance, "Ruthlessness") {
			rogue.AddComboPointsIgnoreTarget(sim, 1, cpMetrics)
		}
	})
}

// Murder talent / Draught of the Sands (these can stack)
func (rogue *Rogue) applyMurder() {
	murderMobTypes := []proto.MobType{proto.MobType_MobTypeHumanoid, proto.MobType_MobTypeGiant, proto.MobType_MobTypeBeast, proto.MobType_MobTypeDragonkin}
	murderTargets := core.FilterSlice(rogue.Env.Encounter.Targets, func(t *core.Target) bool { return slices.Contains(murderMobTypes, t.MobType) })

	if rogue.Talents.Murder > 0 && len(murderTargets) > 0 {
		rogue.Env.RegisterPostFinalizeEffect(func() {
			multiplier := []float64{1, 1.01, 1.02}[rogue.Talents.Murder]

			for _, t := range murderTargets {
				for _, at := range rogue.AttackTables[t.UnitIndex] {
					at.DamageDealtMultiplier *= multiplier
					at.CritMultiplier *= multiplier
				}
			}
		})
	}

	if rogue.Consumes.MiscConsumes != nil && rogue.Consumes.MiscConsumes.DraughtOfTheSands {
		rogue.Env.RegisterPostFinalizeEffect(func() {
			multiplier := 1.02
			for _, t := range rogue.Env.Encounter.Targets {
				for _, at := range rogue.AttackTables[t.UnitIndex] {
					at.DamageDealtMultiplier *= multiplier
					at.CritMultiplier *= multiplier
				}
			}
		})
	}
}

func (rogue *Rogue) applyVilePoisons() {
	if rogue.Talents.VilePoisons == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_RoguePoisonDependent,
		Kind:      core.SpellMod_DamageDone_Flat,
		IntValue:  4 * int64(rogue.Talents.VilePoisons),
	})
}

func (rogue *Rogue) applyAggression() {
	if rogue.Talents.Aggression == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_RogueSinisterStrikeDependent,
		Kind:      core.SpellMod_ImpactDamageDone_Flat,
		IntValue:  2 * int64(rogue.Talents.Aggression),
	})
}

func (rogue *Rogue) applyRelentlessStrikes() {
	if !rogue.Talents.RelentlessStrikes {
		return
	}

	cpMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 14179})
	rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
		if sim.Proc(0.2*float64(comboPoints), "RelentlessStrikes") {
			rogue.AddEnergy(sim, 25, cpMetrics)
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
				rogue.AddComboPoints(sim, 1, result.Target, cpMetrics)
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
						rogue.AddComboPoints(sim, 1, result.Target, cpMetrics)
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

func (rogue *Rogue) registerBladeFlurryCD() {
	if !rogue.Talents.BladeFlurry {
		return
	}

	// TODO verify that this double dips from damage modifiers

	var curDmg float64
	bfHit := rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 22482},
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
		ActionID: core.ActionID{SpellID: 13877},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, 1.2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, 1/1.2)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bfEligible := true

			//Checks for FoK Offhand and 2P TAQ Set Piece Extra Hits.
			if (spell.ActionID.SpellID == 409240 && spell.ActionID.Tag == 2) || spell.ActionID.SpellID == 1213754 {
				bfEligible = false
			}

			if sim.GetNumTargets() < 2 {
				return
			}

			if result.Damage == 0 || !spell.ProcMask.Matches(core.ProcMaskMelee) || !bfEligible {
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
		ClassSpellMask: ClassSpellMask_RogueBladeFlurry,
		ActionID:       core.ActionID{SpellID: 13877},
		Flags:          core.SpellFlagAPL,

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
			rogue.ApplyEnergyTickMultiplier(1.0)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ApplyEnergyTickMultiplier(-1.0)
		},
	})

	rogue.AdrenalineRush = rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueAdrenalineRush,
		ActionID:       AdrenalineRushActionID,
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
		Spell:    rogue.AdrenalineRush,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityBloodlust,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return rogue.CurrentEnergy() <= 45.0
		},
	})
}

func (rogue *Rogue) lethality() float64 {
	return 0.06 * float64(rogue.Talents.Lethality)
}
