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

	// TODO MS: Only use armor if armor DR is < expected resist. Use fast approximation?
	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		// All physical dots (Bleeds) ignore armor.
		if isPeriodic && !spell.Flags.Matches(SpellFlagApplyArmorReduction) {
			return 1, OutcomeEmpty
		}

		// Physical resistance (armor).
		return attackTable.GetArmorDamageModifier(spell), OutcomeEmpty
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
	resistanceCap := float64(unit.Level * 5)

	resistance := -attacker.stats[stats.SpellPenetration]

	// TODO MS: test this for multi-school spells!
	switch schoolIndex {
	case stats.SchoolIndexArcane:
		resistance += unit.GetStat(stats.ArcaneResistance)
	case stats.SchoolIndexFire:
		resistance += unit.GetStat(stats.FireResistance)
	case stats.SchoolIndexFrost:
		resistance += unit.GetStat(stats.FrostResistance)
	case stats.SchoolIndexHoly:
		break // Holy resistance doesn't exist.
	case stats.SchoolIndexNature:
		resistance += unit.GetStat(stats.NatureResistance)
	case stats.SchoolIndexShadow:
		resistance += unit.GetStat(stats.ShadowResistance)
	default:
		// Multi school, choose lowest resistance available.
		resisForMultiSchool := GetMultiSchoolResistanceStats(schoolIndex)
		resistance = 1000.0
		for _, resiStat := range resisForMultiSchool {
			resiVal := unit.GetStat(resiStat)
			if resiVal < resistance {
				resistance = resiVal
			}
		}
	}

	resistance = max(0, resistance)

	levelBasedResist := 0.0
	if !binary {
		levelBasedResist = unit.levelBasedResist(attacker)
	}

	// Pre-TBC all dots that dont have an initial damage component
	// use a 1/10 of the resistance score
	if pureDot {
		resistance /= 10
	}

	return min(resistanceCap, resistance)/resistanceCap + levelBasedResist
}

func (unit *Unit) levelBasedResist(attacker *Unit) float64 {
	if unit.Type == EnemyUnit && unit.Level > attacker.Level {
		return AverageMagicPartialResistPerLevelMultiplier * float64(unit.Level-attacker.Level)
	}
	return 0
}

func (unit *Unit) binaryHitChance(schoolIndex stats.SchoolIndex, attacker *Unit) float64 {
	resistCoeff := unit.resistCoeff(schoolIndex, attacker, true, false)
	return 1 - 0.75*resistCoeff
}

// Roll threshold for each type of partial resist.
func (unit *Unit) partialResistRollThresholds(schoolIndex stats.SchoolIndex, attacker *Unit, pureDot bool) (float64, float64, float64) {
	resistCoeff := unit.resistCoeff(schoolIndex, attacker, false, pureDot)

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
