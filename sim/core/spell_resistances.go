package core

import (
	"github.com/wowsims/sod/sim/core/stats"
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

		if spell.SchoolIndex == stats.SchoolIndexPhysical || MultiSchoolShouldUseArmor(spell.SchoolIndex, attackTable.Defender) {
			// Physical resistance (armor).
			return attackTable.GetArmorDamageModifier(spell), OutcomeEmpty
		}
	}

	// Magical resistance.
	if spell.Flags.Matches(SpellFlagBinary) {
		return 1, OutcomeEmpty
	}

	resistanceRoll := sim.RandomFloat("Partial Resist")

	threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell.SchoolIndex, spell.Flags.Matches(SpellFlagPureDot))
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

// Decide whether to use armor for physical multi school spells
//
// TODO: This is most likely not accurate. A short test showed that it seems
// to not simply use armor if it's lower, but the breakpoint appeared to be pretty
// close to the resistance value.
func MultiSchoolShouldUseArmor(schoolIndex stats.SchoolIndex, target *Unit) bool {
	resistance := 100000.0
	lowestStat := stats.Armor
	for _, resiStat := range GetSchoolResistanceStats(schoolIndex) {
		resiVal := target.GetStat(resiStat)
		if resiVal < resistance {
			resistance = resiVal
			lowestStat = resiStat
		}
	}
	return lowestStat == stats.Armor
}

func (at *AttackTable) GetArmorDamageModifier(spell *Spell) float64 {
	armorPenRating := at.Attacker.stats[stats.ArmorPenetration] + spell.BonusArmorPenRating
	defenderArmor := max(at.Defender.Armor()-armorPenRating, 0.0)
	return 1 - defenderArmor/(defenderArmor+400+85*float64(at.Attacker.Level))
}

func (at *AttackTable) GetPartialResistThresholds(si stats.SchoolIndex, pureDot bool) (float64, float64, float64) {
	return at.Defender.partialResistRollThresholds(si, at.Attacker, pureDot)
}

func (at *AttackTable) GetBinaryHitChance(si stats.SchoolIndex) float64 {
	return at.Defender.binaryHitChance(si, at.Attacker)
}

// All of the following calculations are based on this guide:
// https://royalgiraffe.github.io/resist-guide

func (unit *Unit) resistCoeff(schoolIndex stats.SchoolIndex, attacker *Unit, binary bool, pureDot bool) float64 {
	var resistance float64

	switch schoolIndex {
	case stats.SchoolIndexNone:
		return 0
	case stats.SchoolIndexPhysical:
		return 0
	case stats.SchoolIndexArcane:
		resistance = unit.GetStat(stats.ArcaneResistance)
	case stats.SchoolIndexFire:
		resistance = unit.GetStat(stats.FireResistance)
	case stats.SchoolIndexFrost:
		resistance = unit.GetStat(stats.FrostResistance)
	case stats.SchoolIndexHoly:
		resistance = 0 // Holy resistance doesn't exist.
	case stats.SchoolIndexNature:
		resistance = unit.GetStat(stats.NatureResistance)
	case stats.SchoolIndexShadow:
		resistance = unit.GetStat(stats.ShadowResistance)
	default:
		// Multi school: Choose the lowest resistance available.
		resistance = 1000.0
		if !SpellSchoolFromIndex(schoolIndex).Matches(SpellSchoolHoly) {
			for _, resiStat := range GetSchoolResistanceStats(schoolIndex) {
				resiVal := unit.GetStat(resiStat)
				if resiVal < resistance {
					resistance = resiVal
				}
			}
		} else {
			resistance = 0.0
		}
	}

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

func (unit *Unit) binaryHitChance(schoolIndex stats.SchoolIndex, attacker *Unit) float64 {
	resistCoeff := unit.resistCoeff(schoolIndex, attacker, true, false)
	return 1 - 0.75*resistCoeff
}

// Roll threshold for each type of partial resist.
func (unit *Unit) partialResistRollThresholds(schoolIndex stats.SchoolIndex, attacker *Unit, pureDot bool) (float64, float64, float64) {
	resistCoeff := unit.resistCoeff(schoolIndex, attacker, false, pureDot)

	// Based on the piecewise linear regression estimates at https://royalgiraffe.github.io/partial-resist-table.
	if val := resistCoeff * 3; val <= 1 {
		return 0.76 * val, 0.21 * val, 0.03 * val
	} else if val <= 2 {
		val -= 1
		return 0.76 + 0.24*val, 0.21 + 0.57*val, 0.03 + 0.19*val
	} else {
		val -= 2
		return 1, 0.78 + 0.18*val, 0.22 + 0.58*val
	}
}
