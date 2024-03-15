package core

import (
	"github.com/wowsims/sod/sim/core/stats"
	"math"
)

func (result *SpellResult) applyResistances(sim *Simulation, spell *Spell, isPeriodic bool, attackTable *AttackTable) {
	resistanceMultiplier, outcome := spell.ResistanceMultiplier(sim, isPeriodic, attackTable)

	result.Damage *= resistanceMultiplier
	result.Outcome |= outcome

	result.ResistanceMultiplier = resistanceMultiplier
	result.PreOutcomeDamage = result.Damage
}

// Modifies damage based on Armor or Magic resistances, depending on the damage type.
func (spell *Spell) ResistanceMultiplier(sim *Simulation, isPeriodic bool, attackTable *AttackTable) (float64, HitOutcome) {
	if spell.Flags.Matches(SpellFlagIgnoreResists) {
		return 1, OutcomeEmpty
	}

	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		// All physical dots (Bleeds) ignore armor.
		if isPeriodic && !spell.Flags.Matches(SpellFlagApplyArmorReduction) {
			return 1, OutcomeEmpty
		}

		return attackTable.GetArmorDamageModifier(spell), OutcomeEmpty
	}

	// Magical resistance.
	if spell.Flags.Matches(SpellFlagBinary) {
		return 1, OutcomeEmpty
	}

	resistanceRoll := sim.RandomFloat("Partial Resist")

	threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell.SpellSchool, spell.Flags.Matches(SpellFlagPureDot))
	//if sim.Log != nil {
	//	sim.Log("Resist thresholds: %0.04f, %0.04f, %0.04f", threshold00, threshold25, threshold50)
	//}

	if resistanceRoll > threshold00 {
		// No partial resist.
		return 1, OutcomeEmpty
	} else if resistanceRoll > threshold25 {
		return 0.75, OutcomePartial1_4
	} else if resistanceRoll > threshold50 {
		return 0.5, OutcomePartial2_4
	} else {
		return 0.25, OutcomePartial3_4
	}
}

func (at *AttackTable) GetArmorDamageModifier(spell *Spell) float64 {
	armorPenRating := at.Attacker.stats[stats.ArmorPenetration] + spell.BonusArmorPenRating
	defenderArmor := max(at.Defender.Armor()-armorPenRating, 0.0)
	return 1 - defenderArmor/(defenderArmor+400+85*float64(at.Attacker.Level))
}

func (at *AttackTable) GetPartialResistThresholds(ss SpellSchool, pureDot bool) (float64, float64, float64) {
	return at.Defender.partialResistRollThresholds(ss, at.Attacker, pureDot)
}

func (at *AttackTable) GetBinaryHitChance(ss SpellSchool) float64 {
	return at.Defender.binaryHitChance(ss, at.Attacker)
}

// All the following calculations are based on this guide:
// https://royalgiraffe.github.io/resist-guide

func (unit *Unit) resistance(ss SpellSchool) float64 {
	switch ss {
	case SpellSchoolNone:
		return 0
	case SpellSchoolPhysical:
		return 0
	case SpellSchoolArcane:
		return unit.GetStat(stats.ArcaneResistance)
	case SpellSchoolFire:
		return unit.GetStat(stats.FireResistance)
	case SpellSchoolFrost:
		return unit.GetStat(stats.FrostResistance)
	case SpellSchoolHoly:
		return 0 // Holy resistance doesn't exist.
	case SpellSchoolNature:
		return unit.GetStat(stats.NatureResistance)
	case SpellSchoolShadow:
		return unit.GetStat(stats.ShadowResistance)
	default:
		return unit.multiResistance(ss) // bacteria rejoice
	}
}

func (unit *Unit) multiResistance(ss SpellSchool) float64 {
	resistance := math.MaxFloat64
	if ss.Matches(SpellSchoolArcane) {
		resistance = min(resistance, unit.GetStat(stats.ArcaneResistance))
	}
	if ss.Matches(SpellSchoolFire) {
		resistance = min(resistance, unit.GetStat(stats.FireResistance))
	}
	if ss.Matches(SpellSchoolFrost) {
		resistance = min(resistance, unit.GetStat(stats.FrostResistance))
	}
	if ss.Matches(SpellSchoolHoly) {
		resistance = min(resistance, 0)
	}
	if ss.Matches(SpellSchoolNature) {
		resistance = min(resistance, unit.GetStat(stats.NatureResistance))
	}
	if ss.Matches(SpellSchoolShadow) {
		resistance = min(resistance, unit.GetStat(stats.ShadowResistance))
	}
	return resistance
}

func (unit *Unit) resistCoeff(ss SpellSchool, attacker *Unit, binary bool, pureDot bool) float64 {
	resistance := unit.resistance(ss)

	resistance = max(0, resistance-attacker.stats[stats.SpellPenetration])

	resistanceCap := float64(attacker.Level * 5)
	resistanceCoef := resistance / resistanceCap

	// Pre-TBC all dots that don't have an initial damage component
	// use a 1/10 of the resistance score
	if pureDot {
		resistanceCoef /= 10
	}

	if !binary && unit.Type == EnemyUnit && unit.Level > attacker.Level {
		avgMitigationAdded := AverageMagicPartialResistPerLevelMultiplier * float64(unit.Level-attacker.Level)
		// coef is scaled 0 to 1, not 0 to 0.75
		resistanceCoef += avgMitigationAdded * 1 / 0.75
	}

	return min(1, resistanceCoef)
}

func (unit *Unit) binaryHitChance(ss SpellSchool, attacker *Unit) float64 {
	resistCoeff := unit.resistCoeff(ss, attacker, true, false)
	return 1 - 0.75*resistCoeff
}

// Roll threshold for each type of partial resist.
func (unit *Unit) partialResistRollThresholds(ss SpellSchool, attacker *Unit, pureDot bool) (float64, float64, float64) {
	resistCoeff := unit.resistCoeff(ss, attacker, false, pureDot)

	// Based on the piecewise linear regression estimates at https://royalgiraffe.github.io/partial-resist-table.
	//partialResistChance00 := piecewiseLinear3(resistCoeff, 1, 0.24, 0.00, 0.00)
	partialResistChance25 := piecewiseLinear3(resistCoeff, 0, 0.55, 0.22, 0.04)
	partialResistChance50 := piecewiseLinear3(resistCoeff, 0, 0.18, 0.56, 0.16)
	partialResistChance75 := piecewiseLinear3(resistCoeff, 0, 0.03, 0.22, 0.80)

	return partialResistChance25 + partialResistChance50 + partialResistChance75,
		partialResistChance50 + partialResistChance75,
		partialResistChance75
}

// Interpolation for a 3-part piecewise linear function (which all the partial resist equations use).
func piecewiseLinear3(val float64, p0 float64, p1 float64, p2 float64, p3 float64) float64 {
	if val < 1.0/3.0 {
		return interpolate(val*3, p0, p1)
	} else if val < 2.0/3.0 {
		return interpolate((val-1.0/3.0)*3, p1, p2)
	} else {
		return interpolate((val-2.0/3.0)*3, p2, p3)
	}
}

func interpolate(val float64, p0 float64, p1 float64) float64 {
	return p0*(1-val) + p1*val
}
