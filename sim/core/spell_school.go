package core

import (
	"slices"

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

var spellSchoolsOrdered = []SpellSchool{
	SpellSchoolNone,
	SpellSchoolPhysical,
	SpellSchoolArcane,
	SpellSchoolFire,
	SpellSchoolFrost,
	SpellSchoolHoly,
	SpellSchoolNature,
	SpellSchoolShadow,
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

// Get school index from school mask. Only used for spell registration.
func (ss SpellSchool) GetSchoolIndex() stats.SchoolIndex {
	return stats.SchoolIndex(slices.Index(spellSchoolsOrdered[:], ss))
}

// Get school damage done multiplier.
// Returns highest multiplier if spell is multi school.
func (unit *Unit) GetSchoolDamageDoneMultiplier(spell *Spell) float64 {
	if spell.SchoolIndex != stats.SchoolIndexMultiSchool {
		return unit.PseudoStats.SchoolDamageDealtMultiplier[spell.SchoolIndex]
	}
	return spell.getMultiSchoolMultiplier(unit.PseudoStats.SchoolDamageDealtMultiplier)
}

// Get school damage taken multiplier.
// Returns highest multiplier if spell is multi school.
func (unit *Unit) GetSchoolDamageTakenMultiplier(spell *Spell) float64 {
	if spell.SchoolIndex != stats.SchoolIndexMultiSchool {
		return unit.PseudoStats.SchoolDamageTakenMultiplier[spell.SchoolIndex]
	}
	return spell.getMultiSchoolMultiplier(unit.PseudoStats.SchoolDamageTakenMultiplier)
}

func (spell *Spell) getMultiSchoolMultiplier(multipliers [stats.SchoolLen]float64) float64 {
	var multiplier float64
	for idx, ss := range spellSchoolsOrdered {
		if spell.SpellSchool.Matches(ss) {
			multiplier = max(multiplier, multipliers[idx])
		}
	}
	return multiplier
}
