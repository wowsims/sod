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

// LUT for normal school indices a multischool is comprised of.
var multiSchoolIndexToIndicies = func() [stats.SchoolLen][]stats.SchoolIndex {
	arr := [stats.SchoolLen][]stats.SchoolIndex{}

	for multiIndex := stats.SchoolIndexMultiSchoolStart; multiIndex < stats.SchoolLen; multiIndex++ {
		multiMask := SpellSchoolFromIndex(multiIndex)
		indexArr := []stats.SchoolIndex{}
		for schoolIndex := stats.SchoolIndexNone; schoolIndex < stats.SchoolIndexMultiSchoolStart; schoolIndex++ {
			schoolFlag := SpellSchoolFromIndex(schoolIndex)
			if multiMask.Matches(schoolFlag) {
				indexArr = append(indexArr, schoolIndex)
			}
		}
		arr[multiIndex] = indexArr
	}

	return arr
}()

// Get base school indicies a multi school is comprised of.
func GetMultiSchoolBaseIndices(schoolIndex stats.SchoolIndex) []stats.SchoolIndex {
	return multiSchoolIndexToIndicies[schoolIndex]
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
		SpellSchoolArcane: stats.ArcaneResistance,
		SpellSchoolFire:   stats.FireResistance,
		SpellSchoolFrost:  stats.FrostResistance,
		SpellSchoolNature: stats.NatureResistance,
		SpellSchoolShadow: stats.ShadowResistance,
	}

	arr := [stats.SchoolLen][]stats.Stat{}

	for schoolIndex := stats.SchoolIndexMultiSchoolStart; schoolIndex < stats.SchoolLen; schoolIndex++ {
		msMask := SpellSchoolFromIndex(schoolIndex)
		resiArr := []stats.Stat{}
		for resiSchool, resiStat := range resistances {
			if msMask.Matches(resiSchool) {
				resiArr = append(resiArr, resiStat)
			}
		}
		arr[schoolIndex] = resiArr
	}

	return arr
}()

// Get array of resistance stat indicies for a multi-school.
// Do not use with normal school indicies! See stats.SchoolIndexMultiSchoolStart
func GetMultiSchoolResistanceStats(schoolIndex stats.SchoolIndex) []stats.Stat {
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
	schoolIndex := spell.SchoolIndex
	unit := spell.Unit

	if !IsMultiSchoolIndex(schoolIndex) {
		return
	}

	maxDealt := 0.0
	maxTaken := 0.0
	maxTakenCrit := 0.0

	for _, baseSchoolIndex := range GetMultiSchoolBaseIndices(schoolIndex) {
		dealtMult := unit.PseudoStats.SchoolDamageDealtMultiplier[baseSchoolIndex]
		if dealtMult > maxDealt {
			maxDealt = dealtMult
		}

		takenMult := target.PseudoStats.SchoolDamageTakenMultiplier[baseSchoolIndex]
		if takenMult > maxTaken {
			maxTaken = takenMult
		}

		takenCritMult := target.PseudoStats.SchoolCritTakenMultiplier[baseSchoolIndex]
		if takenCritMult > maxTakenCrit {
			maxTakenCrit = takenCritMult
		}
	}

	unit.PseudoStats.SchoolDamageDealtMultiplier[schoolIndex] = maxDealt
	target.PseudoStats.SchoolDamageTakenMultiplier[schoolIndex] = maxTaken
	target.PseudoStats.SchoolCritTakenMultiplier[schoolIndex] = maxTakenCrit
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

	for _, baseSchoolIndex := range GetMultiSchoolBaseIndices(schoolIndex) {
		dealtMult := unit.PseudoStats.SchoolDamageDealtMultiplier[baseSchoolIndex]
		if dealtMult > maxDealt {
			maxDealt = dealtMult
		}
	}

	unit.PseudoStats.SchoolDamageDealtMultiplier[schoolIndex] = maxDealt
}

// Recalculate damage taken modifier for multi-school.
// Also see spell.RecalculateMultiSchoolModifiers()
func (unit *Unit) MultiSchoolUpdateDamageTakenMod(schoolIndex stats.SchoolIndex) {
	if !IsMultiSchoolIndex(schoolIndex) {
		return
	}

	maxTaken := 0.0
	maxTakenCrit := 0.0

	for _, baseSchoolIndex := range GetMultiSchoolBaseIndices(schoolIndex) {
		takenMult := unit.PseudoStats.SchoolDamageTakenMultiplier[baseSchoolIndex]
		if takenMult > maxTaken {
			maxTaken = takenMult
		}

		takenCritMult := unit.PseudoStats.SchoolCritTakenMultiplier[baseSchoolIndex]
		if takenCritMult > maxTakenCrit {
			maxTakenCrit = takenCritMult
		}
	}

	unit.PseudoStats.SchoolDamageTakenMultiplier[schoolIndex] = maxTaken
	unit.PseudoStats.SchoolCritTakenMultiplier[schoolIndex] = maxTakenCrit
}
