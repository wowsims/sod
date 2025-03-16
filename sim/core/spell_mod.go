package core

import (
	"strconv"
	"time"
)

/*
SpellMod implementation.
*/

type SpellModConfig struct {
	ClassMask         uint64
	Kind              SpellModType
	School            SpellSchool
	SpellFlags        SpellFlag
	SpellFlagsExclude SpellFlag
	DefenseType       DefenseType
	ProcMask          ProcMask
	IntValue          int64
	TimeValue         time.Duration
	FloatValue        float64
	KeyValue          string
	ApplyCustom       SpellModApply
	RemoveCustom      SpellModRemove
}

type SpellMod struct {
	ClassMask         uint64
	Kind              SpellModType
	School            SpellSchool
	SpellFlags        SpellFlag
	SpellFlagsExclude SpellFlag
	DefenseType       DefenseType
	ProcMask          ProcMask
	floatValue        float64
	intValue          int64
	timeValue         time.Duration
	keyValue          string
	Apply             SpellModApply
	Remove            SpellModRemove
	IsActive          bool
	AffectedSpells    []*Spell
}

type SpellModApply func(mod *SpellMod, spell *Spell)
type SpellModRemove func(mod *SpellMod, spell *Spell)
type SpellModFunctions struct {
	Apply  SpellModApply
	Remove SpellModRemove
}

func buildMod(unit *Unit, config SpellModConfig) *SpellMod {
	functions := spellModMap[config.Kind]
	if functions == nil {
		panic("SpellMod " + strconv.Itoa(int(config.Kind)) + " not implemented")
	}

	var applyFn SpellModApply
	var removeFn SpellModRemove

	if config.Kind == SpellMod_Custom {
		if (config.ApplyCustom == nil) || (config.RemoveCustom == nil) {
			panic("ApplyCustom and RemoveCustom are mandatory fields for SpellMod_Custom")
		}

		applyFn = config.ApplyCustom
		removeFn = config.RemoveCustom
	} else {
		applyFn = functions.Apply
		removeFn = functions.Remove
	}

	mod := &SpellMod{
		ClassMask:         config.ClassMask,
		Kind:              config.Kind,
		School:            config.School,
		SpellFlags:        config.SpellFlags,
		SpellFlagsExclude: config.SpellFlagsExclude,
		DefenseType:       config.DefenseType,
		ProcMask:          config.ProcMask,
		floatValue:        config.FloatValue,
		intValue:          config.IntValue,
		timeValue:         config.TimeValue,
		keyValue:          config.KeyValue,
		Apply:             applyFn,
		Remove:            removeFn,
		IsActive:          false,
	}

	unit.OnSpellRegistered(func(spell *Spell) {
		if shouldApply(spell, mod) {
			mod.AffectedSpells = append(mod.AffectedSpells, spell)

			if mod.IsActive {
				mod.Apply(mod, spell)
			}
		}
	})

	return mod
}

func (unit *Unit) AddStaticMod(config SpellModConfig) {
	mod := buildMod(unit, config)
	mod.Activate()
}

func (unit *Unit) AddDynamicMod(config SpellModConfig) *SpellMod {
	return buildMod(unit, config)
}

func shouldApply(spell *Spell, mod *SpellMod) bool {

	if spell.Flags.Matches(SpellFlagNoSpellMods) {
		return false
	}

	if mod.DefenseType > 0 && spell.DefenseType != mod.DefenseType {
		return false
	}

	if mod.SpellFlags > 0 && !spell.Flags.Matches(mod.SpellFlags) {
		return false
	}

	if mod.SpellFlagsExclude > 0 && spell.Flags.Matches(mod.SpellFlagsExclude) {
		return false
	}

	if mod.ClassMask > 0 && !spell.Matches(mod.ClassMask) {
		return false
	}

	if mod.School > 0 && !spell.SpellSchool.Matches(mod.School) {
		return false
	}

	if mod.ProcMask > 0 && !spell.ProcMask.Matches(mod.ProcMask) {
		return false
	}

	return true
}

func (mod *SpellMod) UpdateIntValue(value int64) {
	if mod.IsActive {
		mod.Deactivate()
		mod.intValue = value
		mod.Activate()
	} else {
		mod.intValue = value
	}
}

func (mod *SpellMod) UpdateTimeValue(value time.Duration) {
	if mod.IsActive {
		mod.Deactivate()
		mod.timeValue = value
		mod.Activate()
	} else {
		mod.timeValue = value
	}
}

func (mod *SpellMod) UpdateFloatValue(value float64) {
	if mod.IsActive {
		mod.Deactivate()
		mod.floatValue = value
		mod.Activate()
	} else {
		mod.floatValue = value
	}
}

func (mod *SpellMod) GetIntValue() int64 {
	return mod.intValue
}

func (mod *SpellMod) GetFloatValue() float64 {
	return mod.floatValue
}

func (mod *SpellMod) GetTimeValue() time.Duration {
	return mod.timeValue
}

func (mod *SpellMod) Activate() {
	if mod.IsActive {
		return
	}

	for _, spell := range mod.AffectedSpells {
		mod.Apply(mod, spell)
	}

	mod.IsActive = true
}

func (mod *SpellMod) Deactivate() {
	if !mod.IsActive {
		return
	}

	for _, spell := range mod.AffectedSpells {
		mod.Remove(mod, spell)
	}

	mod.IsActive = false
}

// Mod implmentations
type SpellModType uint32

const (
	// Will multiply the spell.DamageDoneMultiplier. +5% = 0.05
	// Uses FloatValue
	SpellMod_DamageDone_Pct SpellModType = 1 << iota

	// Will add the value spell.DamageMultiplierAdditive
	// Uses IntValue
	SpellMod_DamageDone_Flat

	// Will add the value spell.BaseDamageMultiplierAdditive
	// Uses IntValue
	SpellMod_BaseDamageDone_Flat

	// Will add the value spell.PeriodicDamageMultiplierAdditive
	// Uses IntValue
	SpellMod_PeriodicDamageDone_Flat

	// Will add the value spell.ImpactDamageMultiplierAdditive
	// Uses IntValue
	SpellMod_ImpactDamageDone_Flat

	// Will add the value spell.CritDamageBonus
	// Uses FloatValue
	SpellMod_CritDamageBonus_Flat

	// Will reduce spell.DefaultCast.Cost by % amount. -5% = -0.05
	// Uses IntValue
	SpellMod_PowerCost_Pct

	// Increases or decreases spell.DefaultCast.Cost by flat amount
	// Uses IntValue
	SpellMod_PowerCost_Flat

	// Will add time.Duration to spell.CD.FlatModifier
	// Uses TimeValue
	SpellMod_Cooldown_Flat

	// Increases or decreases spell.CD.Multiplier by flat amount. 50% = 50
	// Apply Aura: Modifies Cooldown (11)
	// -X%
	// Uses IntValue
	SpellMod_Cooldown_Multi_Flat

	// Add/subtract BonusCritRating. +1% = 1.0
	// Uses: FloatValue
	SpellMod_BonusCrit_Flat

	// Add/subtract BonusHitRating. +1% = 1.0
	// Uses: FloatValue
	SpellMod_BonusHit_Flat

	// Will add / substract % amount from the cast time multiplier.
	// Uses: FloatValue
	SpellMod_CastTime_Pct

	// Will add / substract time from the cast time.
	// Uses: TimeValue
	SpellMod_CastTime_Flat

	// Add/subtract to the dots max ticks
	// Uses: IntValue
	SpellMod_DotNumberOfTicks_Flat

	// Add/substract to the base tick frequency
	// Uses: TimeValue
	SpellMod_DotTickLength_Flat

	// Increases or decreases the base tick frequency by % amount. +50% = 0.5
	// Uses: FloatValue
	SpellMod_DotTickLength_Pct

	// Add/subtract to the casts gcd
	// Uses: TimeValue
	SpellMod_GlobalCooldown_Flat

	// Add/subtract bonus coefficient
	// Uses: FloatValue
	SpellMod_BonusCoeffecient_Flat

	// Enables casting while moving
	SpellMod_AllowCastWhileMoving

	// Add/subtract bonus spell power
	// Uses: FloatValue
	SpellMod_BonusDamage_Flat

	// Add/subtract bonus expertise rating
	// Uses: FloatValue
	SpellMod_BonusExpertise_Rating

	// Increases or decreases spell.ThreatMultiplier by flat amount
	// Uses: FloatValue
	SpellMod_Threat_Flat

	// Increases or decreases the spell.ThreatMultiplier by % amount. +50% = 0.5
	// Uses: FloatValue
	SpellMod_Threat_Pct

	// Increases or decreases the spell.FlatThreatBonus by flat amount.
	// Uses: FloatValue
	SpellMod_BonusThreat_Flat

	// Add/subtract duration for associated debuff
	// Uses: KeyValue, TimeValue
	SpellMod_DebuffDuration_Flat

	// Add/subtract duration for associated self-buff
	// Uses: TimeValue
	SpellMod_BuffDuration_Flat

	// User-defined implementation
	// Uses: ApplyCustom | RemoveCustom
	SpellMod_Custom
)

var spellModMap = map[SpellModType]*SpellModFunctions{
	SpellMod_DamageDone_Pct: {
		Apply:  applyDamageDonePercent,
		Remove: removeDamageDonePercent,
	},

	SpellMod_DamageDone_Flat: {
		Apply:  applyDamageDoneAdd,
		Remove: removeDamageDoneAdd,
	},

	SpellMod_BaseDamageDone_Flat: {
		Apply:  applyBaseDamageDoneAdd,
		Remove: removeBaseDamageDoneAdd,
	},

	SpellMod_PeriodicDamageDone_Flat: {
		Apply:  applyPeriodicDamageDoneAdd,
		Remove: removePeriodicDamageDoneAdd,
	},

	SpellMod_ImpactDamageDone_Flat: {
		Apply:  applyImpactDamageDoneAdd,
		Remove: removeImpactDamageDoneAdd,
	},

	SpellMod_CritDamageBonus_Flat: {
		Apply:  applyCritDamageBonusAdd,
		Remove: removeCritDamageBonusAdd,
	},

	SpellMod_PowerCost_Pct: {
		Apply:  applyPowerCostPercent,
		Remove: removePowerCostPercent,
	},

	SpellMod_PowerCost_Flat: {
		Apply:  applyPowerCostFlat,
		Remove: removePowerCostFlat,
	},

	SpellMod_Cooldown_Flat: {
		Apply:  applyCooldownFlat,
		Remove: removeCooldownFlat,
	},

	SpellMod_Cooldown_Multi_Flat: {
		Apply:  applyCooldownMultiplierFlat,
		Remove: removeCooldownMultiplierFlat,
	},

	SpellMod_CastTime_Pct: {
		Apply:  applyCastTimePercent,
		Remove: removeCastTimePercent,
	},

	SpellMod_CastTime_Flat: {
		Apply:  applyCastTimeFlat,
		Remove: removeCastTimeFlat,
	},

	SpellMod_BonusCrit_Flat: {
		Apply:  applyBonusCritFlat,
		Remove: removeBonusCritFlat,
	},

	SpellMod_BonusHit_Flat: {
		Apply:  applyBonusHitFlat,
		Remove: removeBonusHitFlat,
	},

	SpellMod_DotNumberOfTicks_Flat: {
		Apply:  applyDotNumberOfTicks,
		Remove: removeDotNumberOfTicks,
	},

	SpellMod_DotTickLength_Flat: {
		Apply:  applyDotTickLengthFlat,
		Remove: removeDotTickLengthFlat,
	},

	SpellMod_DotTickLength_Pct: {
		Apply:  applyDotTickLengthPercent,
		Remove: removeDotTickLengthPercent,
	},

	SpellMod_GlobalCooldown_Flat: {
		Apply:  applyGlobalCooldownFlat,
		Remove: removeGlobalCooldownFlat,
	},

	SpellMod_BonusCoeffecient_Flat: {
		Apply:  applyBonusCoefficientFlat,
		Remove: removeBonusCoefficientFlat,
	},

	SpellMod_BonusDamage_Flat: {
		Apply:  applyBonusDamageFlat,
		Remove: removeBonusDamageFlat,
	},

	SpellMod_Threat_Flat: {
		Apply:  applyThreatFlat,
		Remove: removeThreatFlat,
	},

	SpellMod_Threat_Pct: {
		Apply:  applyThreatPct,
		Remove: removeThreatPct,
	},

	SpellMod_BonusThreat_Flat: {
		Apply:  applyBonusThreatFlat,
		Remove: removeBonusThreatFlat,
	},

	SpellMod_Custom: {
		// Doesn't have dedicated Apply/Remove functions as ApplyCustom/RemoveCustom is handled in buildMod()
	},
}

func applyDamageDonePercent(mod *SpellMod, spell *Spell) {
	spell.ApplyMultiplicativeDamageBonus(mod.floatValue)
}

func removeDamageDonePercent(mod *SpellMod, spell *Spell) {
	spell.ApplyMultiplicativeDamageBonus(1 / mod.floatValue)
}

func applyDamageDoneAdd(mod *SpellMod, spell *Spell) {
	spell.ApplyAdditiveDamageBonus(mod.intValue)
}

func removeDamageDoneAdd(mod *SpellMod, spell *Spell) {
	spell.ApplyAdditiveDamageBonus(-mod.intValue)
}

func applyBaseDamageDoneAdd(mod *SpellMod, spell *Spell) {
	spell.ApplyAdditiveBaseDamageBonus(mod.intValue)
}

func removeBaseDamageDoneAdd(mod *SpellMod, spell *Spell) {
	spell.ApplyAdditiveBaseDamageBonus(-mod.intValue)
}

func applyPeriodicDamageDoneAdd(mod *SpellMod, spell *Spell) {
	if len(spell.Dots()) > 0 {
		spell.ApplyAdditivePeriodicDamageBonus(mod.intValue)
	}
}

func removePeriodicDamageDoneAdd(mod *SpellMod, spell *Spell) {
	if len(spell.Dots()) > 0 {
		spell.ApplyAdditivePeriodicDamageBonus(-mod.intValue)
	}
}

func applyImpactDamageDoneAdd(mod *SpellMod, spell *Spell) {
	spell.ApplyAdditiveImpactDamageBonus(mod.intValue)
}

func removeImpactDamageDoneAdd(mod *SpellMod, spell *Spell) {
	spell.ApplyAdditiveImpactDamageBonus(-mod.intValue)
}

func applyCritDamageBonusAdd(mod *SpellMod, spell *Spell) {
	spell.CritDamageBonus += mod.floatValue
}

func removeCritDamageBonusAdd(mod *SpellMod, spell *Spell) {
	spell.CritDamageBonus -= mod.floatValue
}

func applyPowerCostPercent(mod *SpellMod, spell *Spell) {
	if spell.Cost != nil {
		spell.Cost.Multiplier += int32(mod.intValue)
	}
}

func removePowerCostPercent(mod *SpellMod, spell *Spell) {
	if spell.Cost != nil {
		spell.Cost.Multiplier -= int32(mod.intValue)
	}
}

func applyPowerCostFlat(mod *SpellMod, spell *Spell) {
	if spell.Cost != nil {
		spell.Cost.FlatModifier += int32(mod.intValue)
	}
}

func removePowerCostFlat(mod *SpellMod, spell *Spell) {
	if spell.Cost != nil {
		spell.Cost.FlatModifier -= int32(mod.intValue)
	}
}

func applyCooldownFlat(mod *SpellMod, spell *Spell) {
	spell.CD.ApplyFlatCooldownMod(mod.timeValue)
}

func removeCooldownFlat(mod *SpellMod, spell *Spell) {
	spell.CD.ApplyFlatCooldownMod(-mod.timeValue)
}

func applyCooldownMultiplierFlat(mod *SpellMod, spell *Spell) {
	spell.CD.ApplyFlatPercentCooldownMod(mod.intValue)
}

func removeCooldownMultiplierFlat(mod *SpellMod, spell *Spell) {
	spell.CD.ApplyFlatPercentCooldownMod(-mod.intValue)
}

func applyCastTimePercent(mod *SpellMod, spell *Spell) {
	if spell.DefaultCast.CastTime > 0 {
		spell.CastTimeMultiplier += mod.floatValue
	}
}

func removeCastTimePercent(mod *SpellMod, spell *Spell) {
	if spell.DefaultCast.CastTime > 0 {
		spell.CastTimeMultiplier -= mod.floatValue
	}
}

func applyCastTimeFlat(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.CastTime += mod.timeValue
}

func removeCastTimeFlat(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.CastTime -= mod.timeValue
}

func applyBonusCritFlat(mod *SpellMod, spell *Spell) {
	spell.BonusCritRating += mod.floatValue
}

func removeBonusCritFlat(mod *SpellMod, spell *Spell) {
	spell.BonusCritRating -= mod.floatValue
}

func applyBonusHitFlat(mod *SpellMod, spell *Spell) {
	spell.BonusHitRating += mod.floatValue
}

func removeBonusHitFlat(mod *SpellMod, spell *Spell) {
	spell.BonusHitRating -= mod.floatValue
}

func applyDotNumberOfTicks(mod *SpellMod, spell *Spell) {
	if spell.dots != nil {
		for _, dot := range spell.dots {
			if dot != nil {
				dot.NumberOfTicks += int32(mod.intValue)
				dot.RecomputeAuraDuration()
			}
		}
	}
	if spell.aoeDot != nil {
		spell.aoeDot.NumberOfTicks += int32(mod.intValue)
		spell.aoeDot.RecomputeAuraDuration()
	}
}

func removeDotNumberOfTicks(mod *SpellMod, spell *Spell) {
	if spell.dots != nil {
		for _, dot := range spell.dots {
			if dot != nil {
				dot.NumberOfTicks -= int32(mod.intValue)
				dot.RecomputeAuraDuration()
			}
		}
	}
	if spell.aoeDot != nil {
		spell.aoeDot.NumberOfTicks -= int32(mod.intValue)
		spell.aoeDot.RecomputeAuraDuration()
	}
}

func applyDotTickLengthFlat(mod *SpellMod, spell *Spell) {
	if spell.dots != nil {
		for _, dot := range spell.dots {
			if dot != nil {
				dot.TickLength += mod.timeValue
				dot.RecomputeAuraDuration()
			}
		}
	}
	if spell.aoeDot != nil {
		spell.aoeDot.TickLength += mod.timeValue
		spell.aoeDot.RecomputeAuraDuration()
	}
}

func removeDotTickLengthFlat(mod *SpellMod, spell *Spell) {
	if spell.dots != nil {
		for _, dot := range spell.dots {
			if dot != nil {
				dot.TickLength -= mod.timeValue
				dot.RecomputeAuraDuration()
			}
		}
	}
	if spell.aoeDot != nil {
		spell.aoeDot.TickLength -= mod.timeValue
		spell.aoeDot.RecomputeAuraDuration()
	}
}

func applyDotTickLengthPercent(mod *SpellMod, spell *Spell) {
	if spell.dots != nil {
		for _, dot := range spell.dots {
			if dot != nil {
				dot.TickLength = time.Duration(float64(dot.TickLength) * mod.floatValue)
				dot.RecomputeAuraDuration()
			}
		}
	}
	if spell.aoeDot != nil {
		spell.aoeDot.TickLength = time.Duration(float64(spell.aoeDot.TickLength) * mod.floatValue)
		spell.aoeDot.RecomputeAuraDuration()
	}
}

func removeDotTickLengthPercent(mod *SpellMod, spell *Spell) {
	if spell.dots != nil {
		for _, dot := range spell.dots {
			if dot != nil {
				dot.TickLength = time.Duration(float64(dot.TickLength) / mod.floatValue)
				dot.RecomputeAuraDuration()
			}
		}
	}
	if spell.aoeDot != nil {
		spell.aoeDot.TickLength = time.Duration(float64(spell.aoeDot.TickLength) / mod.floatValue)
		spell.aoeDot.RecomputeAuraDuration()
	}
}

func applyGlobalCooldownFlat(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.GCD += mod.timeValue
}

func removeGlobalCooldownFlat(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.GCD -= mod.timeValue
}

func applyBonusCoefficientFlat(mod *SpellMod, spell *Spell) {
	spell.BonusCoefficient += mod.floatValue
}

func removeBonusCoefficientFlat(mod *SpellMod, spell *Spell) {
	spell.BonusCoefficient -= mod.floatValue
}

func applyBonusDamageFlat(mod *SpellMod, spell *Spell) {
	spell.BonusDamage += mod.floatValue
}

func removeBonusDamageFlat(mod *SpellMod, spell *Spell) {
	spell.BonusDamage -= mod.floatValue
}

func applyThreatFlat(mod *SpellMod, spell *Spell) {
	spell.ThreatMultiplier += mod.floatValue
}

func removeThreatFlat(mod *SpellMod, spell *Spell) {
	spell.ThreatMultiplier -= mod.floatValue
}

func applyThreatPct(mod *SpellMod, spell *Spell) {
	spell.ThreatMultiplier *= mod.floatValue
}

func removeThreatPct(mod *SpellMod, spell *Spell) {
	spell.ThreatMultiplier /= mod.floatValue
}

func applyBonusThreatFlat(mod *SpellMod, spell *Spell) {
	spell.FlatThreatBonus += mod.floatValue
}

func removeBonusThreatFlat(mod *SpellMod, spell *Spell) {
	spell.FlatThreatBonus -= mod.floatValue
}
