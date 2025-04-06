package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (mage *Mage) ApplyTalents() {
	mage.applyArcaneTalents()
	mage.applyFireTalents()
	mage.applyFrostTalents()
}

func (mage *Mage) applyArcaneTalents() {
	mage.applyArcaneConcentration()
	mage.registerPresenceOfMindCD()
	mage.registerArcanePowerCD()

	// Arcane Subtlety
	if mage.Talents.ArcaneSubtlety > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_Threat_Pct,
			ClassMask:  ClassSpellMask_MageAll,
			School:     core.SpellSchoolArcane,
			FloatValue: 1 - .20*float64(mage.Talents.ArcaneSubtlety),
		})
	}

	// Arcane Focus
	if mage.Talents.ArcaneFocus > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusHit_Flat,
			ClassMask:  ClassSpellMask_MageAll,
			School:     core.SpellSchoolArcane,
			FloatValue: 2 * float64(mage.Talents.ArcaneFocus) * core.SpellHitRatingPerHitChance,
		})
	}

	// Magic Absorption
	if mage.Talents.MagicAbsorption > 0 {
		magicAbsorptionBonus := 2 * float64(mage.Talents.MagicAbsorption)
		mage.AddResistances(magicAbsorptionBonus)
	}

	// Arcane Meditation
	mage.PseudoStats.SpiritRegenRateCasting += 0.05 * float64(mage.Talents.ArcaneMeditation)

	if mage.Talents.ArcaneMind > 0 {
		mage.MultiplyStat(stats.Mana, 1.0+0.02*float64(mage.Talents.ArcaneMind))
	}

	// Arcane Instability
	if mage.Talents.ArcaneInstability > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_MageAll,
			IntValue:  int64(1 * mage.Talents.ArcaneInstability),
		})

		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Flat,
			ClassMask:  ClassSpellMask_MageAll,
			FloatValue: 1 * float64(mage.Talents.ArcaneInstability) * core.SpellCritRatingPerCritChance,
		})
	}
}

func (mage *Mage) applyFireTalents() {
	mage.applyIgnite()
	mage.applyImprovedFireBlast()
	mage.applyIncinerate()
	mage.applyImprovedScorch()
	mage.applyMasterOfElements()

	mage.registerCombustionCD()

	// Burning Soul
	if mage.Talents.BurningSoul > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_Threat_Pct,
			ClassMask:  ClassSpellMask_MageAll,
			School:     core.SpellSchoolFire,
			FloatValue: 1 - .15*float64(mage.Talents.BurningSoul),
		})
	}

	// Critical Mass
	if mage.Talents.CriticalMass > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Flat,
			ClassMask:  ClassSpellMask_MageAll,
			School:     core.SpellSchoolFire,
			FloatValue: 2 * float64(mage.Talents.CriticalMass) * core.SpellCritRatingPerCritChance,
		})
	}

	// Fire Power
	if mage.Talents.FirePower > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_MageAll ^ ClassSpellMask_MageIgnite,
			School:    core.SpellSchoolFire,
			IntValue:  int64(2 * mage.Talents.FirePower),
		})
	}
}

func (mage *Mage) applyFrostTalents() {
	mage.registerColdSnapCD()
	mage.registerIceBarrierSpell()
	mage.applyFrostbite()
	mage.applyImprovedBlizzard()
	mage.applyWintersChill()

	// Elemental Precision
	if mage.Talents.ElementalPrecision > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusHit_Flat,
			ClassMask:  ClassSpellMask_MageAll,
			School:     core.SpellSchoolFire | core.SpellSchoolFrost,
			FloatValue: 2 * float64(mage.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		})
	}

	// Ice Shards
	if mage.Talents.IceShards > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_CritDamageBonus_Flat,
			ClassMask:  ClassSpellMask_MageAll,
			School:     core.SpellSchoolFrost,
			FloatValue: .20 * float64(mage.Talents.IceShards),
		})
	}

	// Piercing Ice
	if mage.Talents.PiercingIce > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_MageAll,
			School:    core.SpellSchoolFrost,
			IntValue:  int64(2 * mage.Talents.PiercingIce),
		})
	}

	if mage.Talents.ArcticReach > 0 {
		rangeModifier := 1 + 0.10*float64(mage.Talents.ArcticReach)
		mage.OnSpellRegistered(func(spell *core.Spell) {
			// TODO: Set max range on these and remove range check
			if spell.Matches(ClassSpellMask_MageFrostbolt|ClassSpellMask_MageBlizzard) && spell.MaxRange > 0 {
				spell.MaxRange *= rangeModifier
			} else if spell.Matches(ClassSpellMask_MageConeOfCold) {
				spell.MinRange *= rangeModifier
			}
		})
	}

	// Frost Channeling
	if mage.Talents.FrostChanneling > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_Threat_Pct,
			ClassMask:  ClassSpellMask_MageAll,
			School:     core.SpellSchoolFrost,
			FloatValue: 1 - .10*float64(mage.Talents.FrostChanneling),
		})

		mage.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Pct,
			ClassMask: ClassSpellMask_MageAll,
			School:    core.SpellSchoolFrost,
			IntValue:  -5 * int64(mage.Talents.FrostChanneling),
		})
	}

	mage.applyShatter()
}

func (mage *Mage) applyArcaneConcentration() {
	if mage.Talents.ArcaneConcentration == 0 {
		return
	}

	procChance := 0.02 * float64(mage.Talents.ArcaneConcentration)

	mage.ClearcastingAura = mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 12577},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolCostMultiplier.AddToMagicSchools(-100)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolCostMultiplier.AddToMagicSchools(100)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}

			if !spell.Matches(ClassSpellMask_MageAll) || spell.Cost == nil {
				return
			}
			aura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:           "Arcane Concentration",
		ClassSpellMask: ClassSpellMask_MageAll,
		ProcMask:       core.ProcMaskSpellDamage,
		ProcChance:     procChance,
		Outcome:        core.OutcomeLanded,
		Callback:       core.CallbackOnSpellHitDealt,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return spell.Cost != nil
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			mage.ClearcastingAura.Activate(sim)
		},
	})
}

func (mage *Mage) registerPresenceOfMindCD() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	actionID := core.ActionID{SpellID: 12043}
	cooldown := time.Second * 180

	classSpellMasks := uint64(ClassSpellMask_MageAll)
	castTimeMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  ClassSpellMask_MageAll,
		FloatValue: -1,
	})

	pomAura := mage.RegisterAura(core.Aura{
		Label:    "Presence of Mind",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
			mage.PresenceOfMind.CD.Use(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(classSpellMasks) && spell.DefaultCast.CastTime > 0 {
				aura.Deactivate(sim)
			}
		},
	})

	mage.PresenceOfMind = mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: cooldown,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.GCD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			pomAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.PresenceOfMind,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerArcanePowerCD() {
	if !mage.Talents.ArcanePower {
		return
	}

	actionID := core.ActionID{SpellID: 12042}

	costMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: ClassSpellMask_MageAll,
		IntValue:  30,
	})
	damageMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_MageAll,
		IntValue:  30,
	})

	buffAura := mage.RegisterAura(core.Aura{
		Label:    "Arcane Power",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
			costMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
			costMod.Deactivate()
		},
	})

	core.RegisterPercentDamageModifierEffect(buffAura, 1.3)

	mage.ArcanePower = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_MageArcanePower,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 180,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buffAura.Activate(sim)
		},
		RelatedSelfBuff: buffAura,
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.ArcanePower,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) applyImprovedFireBlast() {
	if mage.Talents.ImprovedFireBlast == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_MageFireBlast,
		TimeValue: -500 * time.Millisecond * time.Duration(mage.Talents.ImprovedFireBlast),
	})
}

func (mage *Mage) applyIncinerate() {
	if mage.Talents.Incinerate == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Flat,
		ClassMask:  ClassSpellMask_MageScorch | ClassSpellMask_MageFireBlast | ClassSpellMask_MageLivingBombExplosion,
		FloatValue: 2 * float64(mage.Talents.Incinerate) * core.SpellCritRatingPerCritChance,
	})
}

func (mage *Mage) applyImprovedScorch() {
	if mage.Talents.ImprovedScorch == 0 {
		return
	}

	mage.ImprovedScorchAuras = mage.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		return core.ImprovedScorchAura(unit)
	})
}

func (mage *Mage) applyMasterOfElements() {
	if mage.Talents.MasterOfElements == 0 {
		return
	}

	refundCoeff := 0.1 * float64(mage.Talents.MasterOfElements)
	manaMetrics := mage.NewManaMetrics(core.ActionID{SpellID: 29076})

	mage.RegisterAura(core.Aura{
		Label:    "Master of Elements",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}
			if spell.CurCast.Cost == 0 {
				return
			}
			if result.DidCrit() {
				mage.AddMana(sim, spell.Cost.BaseCost*refundCoeff, manaMetrics)
			}
		},
	})
}

func (mage *Mage) registerCombustionCD() {
	if !mage.Talents.Combustion {
		return
	}

	hasOverheatRune := mage.HasRune(proto.MageRune_RuneCloakOverheat)

	actionID := core.ActionID{SpellID: 11129}
	cd := core.Cooldown{
		Timer:    mage.NewTimer(),
		Duration: time.Minute * 3,
	}

	critMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_BonusCrit_Flat,
		ClassMask: ClassSpellMask_MageAll,
		School:    core.SpellSchoolFire,
	})

	numCrits := 0
	critPerStack := 10.0 * core.SpellCritRatingPerCritChance

	mage.CombustionAura = mage.RegisterAura(core.Aura{
		Label:     "Combustion",
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			numCrits = 0
			critMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			cd.Use(sim)
			mage.UpdateMajorCooldowns()
			critMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			critMod.UpdateFloatValue(critPerStack * float64(newStacks))
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || numCrits >= 3 || !spell.SpellSchool.Matches(core.SpellSchoolFire) || !spell.Matches(ClassSpellMask_MageAll) {
				return
			}

			// Ignite, Living Bomb explosions, and Fire Blast with Overheart don't consume crit stacks
			if spell.Matches(ClassSpellMask_MageIgnite|ClassSpellMask_MageLivingBombExplosion) || (hasOverheatRune && spell.Matches(ClassSpellMask_MageFireBlast)) {
				return
			}

			// TODO: This wont work properly with flamestrike
			aura.AddStack(sim)

			if result.DidCrit() {
				numCrits++
				if numCrits == 3 {
					aura.Deactivate(sim)
				}
			}
		},
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Cast: core.CastConfig{
			CD: cd,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mage.CombustionAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.CombustionAura.Activate(sim)
			mage.CombustionAura.AddStack(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

// While bosses Immune Frostbite itself, procs still trigger Fingers of Frost
func (mage *Mage) applyFrostbite() {
	if mage.Talents.Frostbite == 0 {
		return
	}

	frostbiteSpell := mage.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_MageFrostbite,
		ActionID:       core.ActionID{SpellID: 12494},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellProc,
		Flags:          core.SpellFlagPassiveSpell,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			// For whatever reason frostbite always procs Fingers of Frost
			if mage.FingersOfFrostAura != nil {
				mage.FingersOfFrostAura.Activate(sim)
				mage.FingersOfFrostAura.SetStacks(sim, mage.FingersOfFrostAura.MaxStacks)
			}
		},
	})

	procChance := 0.05 * float64(mage.Talents.Frostbite)

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Frostbite Trigger",
	})).AttachProcTrigger(core.ProcTrigger{
		Name:       "Frostbite Trigger Direct",
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		SpellFlags: SpellFlagChillSpell,
		ProcChance: procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			frostbiteSpell.Cast(sim, result.Target)
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Frostbite Trigger Periodic",
		Callback:       core.CallbackOnPeriodicDamageDealt,
		ClassSpellMask: ClassSpellMask_MageBlizzard, // Only procs from Blizzard
		SpellFlags:     SpellFlagChillSpell,
		ProcChance:     procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			frostbiteSpell.Cast(sim, result.Target)
		},
	})
}

func (mage *Mage) registerColdSnapCD() {
	if !mage.Talents.ColdSnap {
		return
	}

	// Grab all frost spells with a CD > 0
	var affectedSpells = []*core.Spell{}
	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellSchool.Matches(core.SpellSchoolFrost) && spell.CD.Duration > 0 {
			affectedSpells = append(affectedSpells, spell)
		}
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 12472},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Duration(time.Minute * 5),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, spell := range affectedSpells {
				spell.CD.Reset()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) applyImprovedBlizzard() {
	if mage.Talents.ImprovedBlizzard == 0 {
		return
	}

	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_MageBlizzard) {
			spell.Flags |= SpellFlagChillSpell
		}
	})
}

func (mage *Mage) applyShatter() {
	mage.FrozenAuras = core.FilterSlice(
		mage.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
			return unit.RegisterAura(core.Aura{
				Label:    fmt.Sprintf("Shatter (%s)", mage.LogLabel()),
				Duration: core.NeverExpires,
			})
		}),
		func(aura *core.Aura) bool { return aura != nil },
	)

	mage.isTargetFrozen = func(target *core.Unit) bool {
		return mage.FrozenAuras.Get(target).IsActive()
	}

	if mage.Talents.Shatter == 0 {
		return
	}
	bonusCrit := 10 * float64(mage.Talents.Shatter) * core.SpellCritRatingPerCritChance

	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_MageAll) && spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
			oldApplyEffects := spell.ApplyEffects
			spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spellBonusCrit := 0.0

				if mage.isTargetFrozen(target) {
					spellBonusCrit += bonusCrit
				}

				spell.BonusCritRating += spellBonusCrit
				oldApplyEffects(sim, target, spell)
				spell.BonusCritRating -= spellBonusCrit
			}
		}
	})
}

func (mage *Mage) applyWintersChill() {
	if mage.Talents.WintersChill == 0 {
		return
	}

	procChance := float64(mage.Talents.WintersChill) * 0.2

	mage.WintersChillAuras = mage.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return core.WintersChillAura(target)
	})

	mage.Env.RegisterPreFinalizeEffect(func() {
		for _, spell := range mage.GetSpellsMatchingSchool(core.SpellSchoolFrost) {
			spell.RelatedAuras = append(spell.RelatedAuras, mage.WintersChillAuras)
		}
	})

	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:             "Winters Chill",
		Callback:         core.CallbackOnSpellHitDealt,
		Outcome:          core.OutcomeLanded,
		ClassSpellMask:   ClassSpellMask_MageAll,
		SpellSchool:      core.SpellSchoolFrost,
		ProcChance:       procChance,
		Harmful:          true,
		CanProcFromProcs: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			aura := mage.WintersChillAuras.Get(result.Target)
			aura.Activate(sim)
			aura.AddStack(sim)
		},
	})
}
