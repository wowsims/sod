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

	// Physical x Other
	SpellSchoolSpellstrike  = SpellSchoolPhysical | SpellSchoolArcane
	SpellSchoolFlamestrike  = SpellSchoolPhysical | SpellSchoolFire
	SpellSchoolFroststrike  = SpellSchoolPhysical | SpellSchoolFrost
	SpellSchoolHolystrike   = SpellSchoolPhysical | SpellSchoolHoly
	SpellSchoolStormstrike  = SpellSchoolPhysical | SpellSchoolNature
	SpellSchoolShadowstrike = SpellSchoolPhysical | SpellSchoolShadow

	// Arcane x Other
	SpellSchoolSpellfire   = SpellSchoolArcane | SpellSchoolFire
	SpellSchoolSpellFrost  = SpellSchoolArcane | SpellSchoolFrost
	SpellSchoolDivine      = SpellSchoolArcane | SpellSchoolHoly
	SpellSchoolAstral      = SpellSchoolArcane | SpellSchoolNature
	SpellSchoolSpellShadow = SpellSchoolArcane | SpellSchoolShadow

	// Fire x Other
	SpellSchoolFrostfire   = SpellSchoolFire | SpellSchoolFrost
	SpellSchoolRadiant     = SpellSchoolFire | SpellSchoolHoly
	SpellSchoolVolcanic    = SpellSchoolFire | SpellSchoolNature
	SpellSchoolShadowflame = SpellSchoolFire | SpellSchoolShadow

	// Frost x Other
	SpellSchoolHolyfrost   = SpellSchoolFrost | SpellSchoolHoly
	SpellSchoolFroststorm  = SpellSchoolFrost | SpellSchoolNature
	SpellSchoolShadowfrost = SpellSchoolFrost | SpellSchoolShadow

	// Holy x Other
	SpellSchoolHolystorm = SpellSchoolHoly | SpellSchoolNature
	SpellSchoolTwilight  = SpellSchoolHoly | SpellSchoolShadow

	// Nature x Other
	SpellSchoolPlague = SpellSchoolNature | SpellSchoolShadow

	SpellSchoolElemental = SpellSchoolFire | SpellSchoolFrost | SpellSchoolNature

	SpellSchoolAttack = SpellSchoolPhysical |
		SpellSchoolSpellstrike | SpellSchoolFlamestrike | SpellSchoolFroststrike | SpellSchoolHolystrike | SpellSchoolStormstrike | SpellSchoolShadowstrike

	SpellSchoolMagic = SpellSchoolArcane | SpellSchoolFire | SpellSchoolFrost | SpellSchoolHoly | SpellSchoolNature | SpellSchoolShadow
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

	// Physical x Other
	SpellSchoolSpellstrike,
	SpellSchoolFlamestrike,
	SpellSchoolFroststrike,
	SpellSchoolHolystrike,
	SpellSchoolStormstrike,
	SpellSchoolShadowstrike,

	// Arcane x Other
	SpellSchoolSpellfire,
	SpellSchoolSpellFrost,
	SpellSchoolDivine,
	SpellSchoolAstral,
	SpellSchoolSpellShadow,

	// Fire x Other
	SpellSchoolFrostfire,
	SpellSchoolRadiant,
	SpellSchoolVolcanic,
	SpellSchoolShadowflame,

	// Frost x Other
	SpellSchoolHolyfrost,
	SpellSchoolFroststorm,
	SpellSchoolShadowfrost,

	// Holy x Other
	SpellSchoolHolystorm,
	SpellSchoolTwilight,

	// Nature x Other
	SpellSchoolPlague,

	SpellSchoolElemental,
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
		for baseSchoolIndex := stats.SchoolIndexNone; baseSchoolIndex < stats.PrimarySchoolLen; baseSchoolIndex++ {
			schoolFlag := SpellSchoolFromIndex(baseSchoolIndex)
			if multiMask.Matches(schoolFlag) {
				indexArr = append(indexArr, baseSchoolIndex)
			}
		}
		arr[schoolIndex] = indexArr
	}

	return arr
}()

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

// Get array of resistance stat indices for a (multi)school.
// Physical school uses Armor as stat index!
func GetSchoolResistanceStats(schoolIndex stats.SchoolIndex) []stats.Stat {
	return schoolIndexToResistanceStats[schoolIndex]
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

// Sets school index, mask and base indices.
func (spell *Spell) SetSchool(schoolIndex stats.SchoolIndex) {
	spell.SchoolIndex = schoolIndex
	spell.SpellSchool = SpellSchoolFromIndex(schoolIndex)
	spell.SchoolBaseIndices = schoolIndexToIndices[schoolIndex]
	spell.IsMultischool = schoolIndex.IsMultiSchool()
}

func selectMaxMultInSchoolArray(spell *Spell, array *[stats.PrimarySchoolLen]float64) float64 {
	high := 0.0
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
	if !spell.IsMultischool {
		return unit.PseudoStats.SchoolDamageDealtMultiplier[spell.SchoolIndex]
	}
	return selectMaxMultInSchoolArray(spell, &unit.PseudoStats.SchoolDamageDealtMultiplier)
}

// Get school damage taken multiplier.
// Returns highest multiplier if spell is multi school.
func (unit *Unit) GetSchoolDamageTakenMultiplier(spell *Spell) float64 {
	if !spell.IsMultischool {
		return unit.PseudoStats.SchoolDamageTakenMultiplier[spell.SchoolIndex]
	}
	return selectMaxMultInSchoolArray(spell, &unit.PseudoStats.SchoolDamageTakenMultiplier)
}

// Get school crit taken multiplier.
// Returns highest multiplier if spell is multi school.
func (unit *Unit) GetCritTakenMultiplier(spell *Spell) float64 {
	if !spell.IsMultischool {
		return unit.PseudoStats.SchoolCritTakenMultiplier[spell.SchoolIndex]
	}
	return selectMaxMultInSchoolArray(spell, &unit.PseudoStats.SchoolCritTakenMultiplier)
}
