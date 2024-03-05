package core

import (
	"fmt"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type SpellSchool byte

const (
	SpellSchoolNone     SpellSchool = 0
	SpellSchoolPhysical SpellSchool = 1 << iota
	SpellSchoolArcane
	SpellSchoolFire
	SpellSchoolFrost
	SpellSchoolHoly
	SpellSchoolNature
	SpellSchoolShadow

	SpellSchoolMagic = SpellSchoolArcane | SpellSchoolFire | SpellSchoolFrost | SpellSchoolHoly | SpellSchoolNature | SpellSchoolShadow

	SpellSchoolArcaneFire     = SpellSchoolArcane | SpellSchoolFire
	SpellSchoolArcaneFrost    = SpellSchoolArcane | SpellSchoolFrost
	SpellSchoolFireFrost      = SpellSchoolFire | SpellSchoolFrost
	SpellSchoolFireShadow     = SpellSchoolFire | SpellSchoolShadow
	SpellSchoolFrostShadow    = SpellSchoolFrost | SpellSchoolShadow
	SpellSchoolPhysicalNature = SpellSchoolPhysical | SpellSchoolNature
	SpellSchoolNatureShadow   = SpellSchoolNature | SpellSchoolShadow
	SpellSchoolPhysicalShadow = SpellSchoolPhysical | SpellSchoolShadow
)

// Get associated school mask for a school index.
// Keep in sync with stats.SchoolIndex
var schoolIndexToSchoolMask = [stats.SchoolLen]SpellSchool{
	SpellSchoolNone,
	SpellSchoolPhysical,
	SpellSchoolArcane,
	SpellSchoolFire,
	SpellSchoolFrost,
	SpellSchoolHoly,
	SpellSchoolNature,
	SpellSchoolShadow,
	SpellSchoolArcaneFire,
	SpellSchoolArcaneFrost,
	SpellSchoolFireFrost,
	SpellSchoolFireShadow,
	SpellSchoolFrostShadow,
	SpellSchoolPhysicalNature,
	SpellSchoolNatureShadow,
	SpellSchoolPhysicalShadow,
}

var schoolMaskToIndex = func() map[SpellSchool]stats.SchoolIndex {
	mti := map[SpellSchool]stats.SchoolIndex{}
	for i := stats.SchoolIndexNone; i < stats.SchoolLen; i++ {
		mti[schoolIndexToSchoolMask[i]] = i
	}
	return mti
}()

// LUT for base school indices a (multi)school is made of.
var schoolIndexToIndices = func() [stats.SchoolLen][]stats.SchoolIndex {
	arr := [stats.SchoolLen][]stats.SchoolIndex{}

	for schoolIndex := stats.SchoolIndexNone; schoolIndex < stats.SchoolLen; schoolIndex++ {
		multiMask := SpellSchoolFromIndex(schoolIndex)
		indexArr := []stats.SchoolIndex{}
		for baseSchoolIndex := stats.SchoolIndexNone; baseSchoolIndex < stats.SchoolIndexMultiSchoolStart; baseSchoolIndex++ {
			schoolFlag := SpellSchoolFromIndex(baseSchoolIndex)
			if multiMask.Matches(schoolFlag) {
				indexArr = append(indexArr, baseSchoolIndex)
			}
		}
		arr[schoolIndex] = indexArr
	}

	return arr
}()

// Get base school indices of the spell.
// If spell is a single school the array will just contain that school's index.
// If spell is multi school it will include all school indices the multi school is made of.
// TODO MS: Move onto spell?
func (spell *Spell) GetSchoolBaseIndices() []stats.SchoolIndex {
	return schoolIndexToIndices[spell.SchoolIndex]
}

// Check if school index is a multi-school.
func IsMultiSchoolIndex(schoolIndex stats.SchoolIndex) bool {
	return schoolIndex >= stats.SchoolIndexMultiSchoolStart
}

// Get spell school mask from school index.
func SpellSchoolFromIndex(schoolIndex stats.SchoolIndex) SpellSchool {
	return schoolIndexToSchoolMask[schoolIndex]
}

// Returns whether there is any overlap between the given masks.
func (ss SpellSchool) Matches(other SpellSchool) bool {
	return (ss & other) != 0
}

// Get school index from school mask. Will error if mask is for an undefined multi-school.
// This involves a map lookup. Do not use in hot path code.
func (ss SpellSchool) GetSchoolIndex() stats.SchoolIndex {
	idx, ok := schoolMaskToIndex[ss]
	if !ok {
		panic(fmt.Sprintf("No school index defined for schoolmask %d! You may need to define a new multi-school.", ss))
	}
	return idx
}

// LUT for resistance stat indices used by each multischool.
var schoolIndexToResistanceStats = func() [stats.SchoolLen][]stats.Stat {
	resistances := map[SpellSchool]stats.Stat{
		SpellSchoolPhysical: stats.Armor, // This is technically physical resistance
		SpellSchoolArcane:   stats.ArcaneResistance,
		SpellSchoolFire:     stats.FireResistance,
		SpellSchoolFrost:    stats.FrostResistance,
		SpellSchoolNature:   stats.NatureResistance,
		SpellSchoolShadow:   stats.ShadowResistance,
	}

	arr := [stats.SchoolLen][]stats.Stat{}

	for schoolIndex := stats.SchoolIndexPhysical; schoolIndex < stats.SchoolLen; schoolIndex++ {
		schoolMask := SpellSchoolFromIndex(schoolIndex)
		resiArr := []stats.Stat{}
		for resiSchool, resiStat := range resistances {
			if schoolMask.Matches(resiSchool) {
				resiArr = append(resiArr, resiStat)
			}
		}
		arr[schoolIndex] = resiArr
	}

	return arr
}()

// Get array of resistance stat indices for a (multi)school.
// Physical school uses Armor as stat index!
func GetSchoolResistanceStats(schoolIndex stats.SchoolIndex) []stats.Stat {
	return schoolIndexToResistanceStats[schoolIndex]
}

func SpellSchoolFromProto(p proto.SpellSchool) SpellSchool {
	switch p {
	case proto.SpellSchool_SpellSchoolPhysical:
		return SpellSchoolPhysical
	case proto.SpellSchool_SpellSchoolArcane:
		return SpellSchoolArcane
	case proto.SpellSchool_SpellSchoolFire:
		return SpellSchoolFire
	case proto.SpellSchool_SpellSchoolFrost:
		return SpellSchoolFrost
	case proto.SpellSchool_SpellSchoolHoly:
		return SpellSchoolHoly
	case proto.SpellSchool_SpellSchoolNature:
		return SpellSchoolNature
	case proto.SpellSchool_SpellSchoolShadow:
		return SpellSchoolShadow
	default:
		return SpellSchoolPhysical
	}
}

// Recalculate multipliers used for given multi school for unit and target.
// This needs to happen each time before a multi school spell enters its hit and damage calculations.
//
// Note: This is an overall highly unoptimized approach and should probably change if multi-school
// spells ever become a major part of all spells used. In that case recalculation should be
// hooked to change of the base school modifiers by e.g. implementing unit.ModifySchoolXxxxModifier() functions,
// to then update the affected multi schools as needed.
// Doing that would add overhead to all school modifier updates, which doesn't seem worth
// it in the context of SoD as of writing this.
func (spell *Spell) MultiSchoolUpdateModifiers(target *Unit) {
	spell.MultiSchoolUpdateDamageDealtMod()
	target.MultiSchoolUpdateDamageTakenMod(spell)
}

// Recalculate damage done modifier for multi-school.
// Also see spell.RecalculateMultiSchoolModifiers()
func (spell *Spell) MultiSchoolUpdateDamageDealtMod() {
	schoolIndex := spell.SchoolIndex
	unit := spell.Unit

	if !IsMultiSchoolIndex(schoolIndex) {
		return
	}

	maxDealt := 0.0

	for _, baseSchoolIndex := range spell.GetSchoolBaseIndices() {
		dealtMult := unit.PseudoStats.SchoolDamageDealtMultiplier[baseSchoolIndex]
		if dealtMult > maxDealt {
			maxDealt = dealtMult
		}
	}

	unit.PseudoStats.SchoolDamageDealtMultiplier[schoolIndex] = maxDealt
}

// Recalculate damage taken modifier for multi-school.
// Also see spell.RecalculateMultiSchoolModifiers()
func (unit *Unit) MultiSchoolUpdateDamageTakenMod(spell *Spell) {
	if !IsMultiSchoolIndex(spell.SchoolIndex) {
		return
	}

	maxTaken := 0.0
	maxTakenCrit := 0.0

	for _, baseSchoolIndex := range spell.GetSchoolBaseIndices() {
		takenMult := unit.PseudoStats.SchoolDamageTakenMultiplier[baseSchoolIndex]
		if takenMult > maxTaken {
			maxTaken = takenMult
		}

		takenCritMult := unit.PseudoStats.SchoolCritTakenMultiplier[baseSchoolIndex]
		if takenCritMult > maxTakenCrit {
			maxTakenCrit = takenCritMult
		}
	}

	unit.PseudoStats.SchoolDamageTakenMultiplier[spell.SchoolIndex] = maxTaken
	unit.PseudoStats.SchoolCritTakenMultiplier[spell.SchoolIndex] = maxTakenCrit
}
