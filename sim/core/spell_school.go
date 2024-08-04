package core

import (
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
}

// Get spell school mask from school index.
func SpellSchoolFromIndex(schoolIndex stats.SchoolIndex) SpellSchool {
	return schoolIndexToSchoolMask[schoolIndex]
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

// Returns whether there is any overlap between the given masks.
func (ss SpellSchool) Matches(other SpellSchool) bool {
	return (ss & other) != 0
}

// Get school index from school mask.
func (ss SpellSchool) GetSchoolIndex() stats.SchoolIndex {
	switch ss {
	case SpellSchoolNone:
		return stats.SchoolIndexNone
	case SpellSchoolPhysical:
		return stats.SchoolIndexPhysical
	case SpellSchoolArcane:
		return stats.SchoolIndexArcane
	case SpellSchoolFire:
		return stats.SchoolIndexFire
	case SpellSchoolFrost:
		return stats.SchoolIndexFrost
	case SpellSchoolHoly:
		return stats.SchoolIndexHoly
	case SpellSchoolNature:
		return stats.SchoolIndexNature
	case SpellSchoolShadow:
		return stats.SchoolIndexShadow
	default:
		return stats.SchoolIndexMultischool
	}
}

func (schoolMask SpellSchool) GetBaseIndices() []stats.SchoolIndex {
	indexArr := []stats.SchoolIndex{}
	for baseSchoolIndex := stats.SchoolIndexNone; baseSchoolIndex < stats.SchoolLen; baseSchoolIndex++ {
		schoolFlag := SpellSchoolFromIndex(baseSchoolIndex)
		if schoolMask.Matches(schoolFlag) {
			indexArr = append(indexArr, baseSchoolIndex)
		}
	}
	return indexArr
}

func selectMaxMultInSchoolArray[T stats.SchoolValueArrayValues](spell *Spell, array *stats.SchoolValueArray[T]) T {
	var high T = 0
	for _, baseIndex := range spell.SchoolBaseIndices {
		mult := array[baseIndex]
		if mult > high {
			high = mult
		}
	}
	return high
}

// Get school damage done multiplier.
// Returns highest multiplier if spell is multi school.
func (unit *Unit) GetSchoolDamageDoneMultiplier(spell *Spell) float64 {
	if !spell.SchoolIndex.IsMultiSchool() {
		return unit.PseudoStats.SchoolDamageDealtMultiplier[spell.SchoolIndex]
	}
	return selectMaxMultInSchoolArray(spell, &unit.PseudoStats.SchoolDamageDealtMultiplier)
}

// Get school damage taken multiplier.
// Returns highest multiplier if spell is multi school.
func (unit *Unit) GetSchoolDamageTakenMultiplier(spell *Spell) float64 {
	if !spell.SchoolIndex.IsMultiSchool() {
		return unit.PseudoStats.SchoolDamageTakenMultiplier[spell.SchoolIndex]
	}
	return selectMaxMultInSchoolArray(spell, &unit.PseudoStats.SchoolDamageTakenMultiplier)
}

// Returns highest if spell is multi school.
func (unit *Unit) GetSchoolCritTakenChance(spell *Spell) float64 {
	if !spell.SchoolIndex.IsMultiSchool() {
		return unit.PseudoStats.SchoolCritTakenChance[spell.SchoolIndex]
	}
	return selectMaxMultInSchoolArray(spell, &unit.PseudoStats.SchoolCritTakenChance)
}

// Returns highest if spell is multi school.
func (unit *Unit) GetSchoolBonusDamageTaken(spell *Spell) float64 {
	if !spell.SchoolIndex.IsMultiSchool() {
		return unit.PseudoStats.SchoolBonusDamageTaken[spell.SchoolIndex]
	}
	return selectMaxMultInSchoolArray(spell, &unit.PseudoStats.SchoolBonusDamageTaken)
}

// Get school bonus hit chance
// Returns highest multiplier if spell is multi school.
func (unit *Unit) GetSchoolBonusHitChance(spell *Spell) float64 {
	if !spell.SchoolIndex.IsMultiSchool() {
		return unit.PseudoStats.SchoolBonusHitChance[spell.SchoolIndex]
	}
	return selectMaxMultInSchoolArray(spell, &unit.PseudoStats.SchoolBonusHitChance)
}

// Get cost modifier for school
// Returns highest mod if spell is multi school.
func (unit *Unit) GetSchoolCostModifier(spell *Spell) int32 {
	if !spell.SchoolIndex.IsMultiSchool() {
		return unit.PseudoStats.SchoolCostMultiplier[spell.SchoolIndex]
	}
	return selectMaxMultInSchoolArray(spell, &unit.PseudoStats.SchoolCostMultiplier)
}
