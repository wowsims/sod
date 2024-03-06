package core

import (
	"fmt"
	"math"
	"testing"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func Test_MultiSchoolResistance(t *testing.T) {
	attacker := &Unit{
		Type:  PlayerUnit,
		Level: 60,
		stats: stats.Stats{},
	}

	defender := &Unit{
		Type:  EnemyUnit,
		Level: attacker.Level + 3,
		stats: stats.Stats{},
	}

	attackTable := NewAttackTable(attacker, defender, nil)

	for schoolIndex := stats.SchoolIndexArcane; schoolIndex < stats.SchoolLen; schoolIndex++ {
		spellSchool := SpellSchoolFromIndex(schoolIndex)

		spell := &Spell{
			SpellSchool: spellSchool,
			SchoolIndex: spellSchool.GetSchoolIndex(),
		}

		resistanceStats := GetSchoolResistanceStats(spell.SchoolIndex)
		const lowestValue float64 = 50.0
		var lowestStat stats.Stat
		isHoly := false

		for rev := 0; rev < 2; rev++ {
			if spell.SpellSchool.Matches(SpellSchoolHoly) {
				isHoly = true
			} else {
				if rev != 0 {
					resiLen := len(resistanceStats)
					lowestStat = resistanceStats[resiLen-1]
					j := resiLen - 1
					for i := 0; i < resiLen; i++ {
						defender.stats[resistanceStats[j]] = lowestValue + 25.0*float64(i)
						j--
					}
				} else {
					lowestStat = resistanceStats[0]
					for i, stat := range resistanceStats {
						defender.stats[stat] = lowestValue + 25.0*float64(i)
					}
				}
			}

			resistance := 0.0

			if !isHoly {
				resistance = lowestValue

				// Make sure setup is right
				var lowestFound stats.Stat
				lowestValFound := 99999.0
				for _, checkStat := range resistanceStats {
					if defender.GetStat(checkStat) < lowestValFound {
						lowestValFound = defender.GetStat(checkStat)
						lowestFound = checkStat
					}
				}
				if lowestFound != lowestStat || lowestValFound != resistance {
					t.Errorf("Expected resist %d to be lowest with %f, but found %d at %f to be lowest resist!", lowestStat, resistance, lowestFound, lowestValFound)
					return
				}
			}

			// Expected values
			resistanceCap := float64(attacker.Level * 5)
			levelBased := float64(max(defender.Level-attacker.Level, 0)) * 0.02
			expectedCoef := min(1, resistance/resistanceCap+levelBased*1/0.75)
			expectedAvgMitigation := expectedCoef*0.75 - 3.0/16.0*max(0, expectedCoef-2.0/3.0)

			// Check if coef is correct to begin with
			resistCoef := defender.resistCoeff(spell.SchoolIndex, attacker, false, false)
			if math.Abs(resistCoef-expectedCoef) > 0.001 {
				t.Errorf("Resist coef is %.3f but expected %.3f at resistance %f", resistCoef, expectedCoef, resistance)
				return
			}

			// Check breakpoints
			threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell.SchoolIndex, spell.Flags.Matches(SpellFlagPureDot))
			chance25 := threshold00 - threshold25
			chance50 := threshold25 - threshold50
			chance75 := threshold50
			avgResist := chance25*0.25 + chance50*0.50 + chance75*0.75
			if math.Abs(avgResist-expectedAvgMitigation) > 0.005 {
				t.Errorf("resist = %.2f, thresholds = %f, resultingAr = %.2f%%, expectedAr = %.2f%%", resistance, threshold00, avgResist, expectedAvgMitigation)
				return
			}
		}
	}
}

func Test_MultiSchoolResistanceArmor(t *testing.T) {
	spell := &Spell{
		SpellSchool: SpellSchoolFlamestrike,
		SchoolIndex: SpellSchoolFlamestrike.GetSchoolIndex(),
	}

	attacker := &Unit{
		Type:  PlayerUnit,
		Level: 60,
		stats: stats.Stats{},
	}

	defender := &Unit{
		Type:        EnemyUnit,
		Level:       63,
		stats:       stats.Stats{},
		PseudoStats: stats.NewPseudoStats(),
	}

	attackTable := NewAttackTable(attacker, defender, nil)

	sim := NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{},
		Encounter:  &proto.Encounter{},
		Raid:       &proto.Raid{},
	})

	defender.AddStat(stats.Armor, 100)
	defender.AddStat(stats.FireResistance, 200)

	mult, outcome := spell.ResistanceMultiplier(sim, false, attackTable)

	if outcome != OutcomeEmpty || mult == 1 {
		t.Errorf("Expected empty outcome with armor %f and resistance %f", defender.GetStat(stats.Armor), defender.GetStat(stats.FireResistance))
		return
	}

	defender.AddStat(stats.Armor, 200)
	mult, outcome = spell.ResistanceMultiplier(sim, false, attackTable)
	isPartial := (outcome & OutcomePartial) != OutcomePartial
	isFullHit := outcome == OutcomeEmpty && mult == 1
	if !isFullHit && !isPartial {
		t.Errorf("Expected empty outcome with armor %f and resistance %f", defender.GetStat(stats.Armor), defender.GetStat(stats.FireResistance))
		return
	}
}

func Test_MultiSchoolSpellPower(t *testing.T) {
	caster := &Unit{
		Type:        PlayerUnit,
		Level:       60,
		stats:       stats.Stats{},
		PseudoStats: stats.NewPseudoStats(),
	}

	spell := &Spell{
		SpellSchool: SpellSchoolNone,
		SchoolIndex: SpellSchoolNone.GetSchoolIndex(),
		Unit:        caster,
	}

	for schoolIndex := stats.SchoolIndexArcane; schoolIndex < stats.SchoolLen; schoolIndex++ {
		spellSchool := SpellSchoolFromIndex(schoolIndex)

		spell.SchoolIndex = schoolIndex
		spell.SchoolIndex = spellSchool.GetSchoolIndex()

		baseIndeces := spell.GetSchoolBaseIndices()
		const highestValue float64 = 555.0
		var highestStat stats.Stat

		// Note: If powerStat == stats.SpellPower, then that means physical school, which is PseudoStats.BonusDamage!

		for rev := 0; rev < 2; rev++ {
			if rev != 0 {
				indexLen := len(baseIndeces)
				highestStat = stats.ArcanePower + stats.Stat(baseIndeces[indexLen-1]) - 2

				for i := indexLen - 1; i >= 0; i-- {
					powerStat := stats.ArcanePower + stats.Stat(baseIndeces[i]) - 2
					if powerStat == stats.SpellPower {
						caster.PseudoStats.BonusDamage = highestValue - 25.0*float64(indexLen-1-i)
					} else {
						caster.stats[powerStat] = highestValue - 25.0*float64(indexLen-1-i)
					}
				}
			} else {
				highestStat = stats.ArcanePower + stats.Stat(baseIndeces[0]) - 2
				for i, baseIndex := range baseIndeces {
					powerStat := stats.ArcanePower + stats.Stat(baseIndex) - 2
					if powerStat == stats.SpellPower {
						caster.PseudoStats.BonusDamage = highestValue - 25.0*float64(i)
					} else {
						caster.stats[powerStat] = highestValue - 25.0*float64(i)
					}
				}
			}

			// Make sure setup is right
			var highestFound stats.Stat
			highestValFound := 0.0
			for _, baseIndex := range baseIndeces {
				powerStat := stats.ArcanePower + stats.Stat(baseIndex) - 2
				if powerStat == stats.SpellPower {
					if caster.PseudoStats.BonusDamage > highestValFound {
						highestValFound = caster.PseudoStats.BonusDamage
						highestFound = powerStat
					}
				} else {
					if caster.GetStat(powerStat) > highestValFound {
						highestValFound = caster.GetStat(powerStat)
						highestFound = powerStat
					}
				}
			}
			if highestFound != highestStat || highestValFound != highestValue {
				t.Errorf("Expected power %d to be highest with %f, but found %d at %f to be highest school power for index %d!", highestStat, highestValue, highestFound, highestValFound, schoolIndex)
				return
			}

			power := spell.SpellSchoolPower()
			if power != highestValue {
				t.Errorf("Expected %f to be highest power value found, but got %f for index %d!", highestValue, power, schoolIndex)
				return
			}
		}
	}
}

const highestMult float64 = 5.0

func SchoolMultiplierArrayHelper(t *testing.T, caster *Unit, target *Unit, multArray *[stats.SchoolLen]float64,
	testFunc func(spell *Spell, schoolIndex stats.SchoolIndex) (bool, string)) {

	spell := &Spell{
		SpellSchool:              SpellSchoolNone,
		SchoolIndex:              SpellSchoolNone.GetSchoolIndex(),
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		Unit:                     caster,
	}

	for schoolIndex := stats.SchoolIndexSpellstrike; schoolIndex < stats.SchoolLen; schoolIndex++ {
		spellSchool := SpellSchoolFromIndex(schoolIndex)

		spell.SchoolIndex = schoolIndex
		spell.SchoolIndex = spellSchool.GetSchoolIndex()

		baseIndeces := spell.GetSchoolBaseIndices()

		for i, baseIndex := range baseIndeces {
			multArray[baseIndex] = highestMult - float64(i)*0.5
		}

		ok, errMsg := testFunc(spell, schoolIndex)

		if !ok {
			t.Error(errMsg)
			return
		}
	}
}

func Test_MultiSchoolModifiers(t *testing.T) {
	caster := &Unit{
		Type:        PlayerUnit,
		Level:       60,
		stats:       stats.Stats{},
		PseudoStats: stats.NewPseudoStats(),
	}

	target := &Unit{
		Type:        EnemyUnit,
		Level:       63,
		stats:       stats.Stats{},
		PseudoStats: stats.NewPseudoStats(),
	}

	attackTable := NewAttackTable(caster, target, nil)

	t.Run("DamageDealt", func(t *testing.T) {
		SchoolMultiplierArrayHelper(t, caster, target, &caster.PseudoStats.SchoolDamageDealtMultiplier,
			func(spell *Spell, schoolIndex stats.SchoolIndex) (bool, string) {
				spell.MultiSchoolUpdateModifiers(target)
				mult := spell.AttackerDamageMultiplier(attackTable)
				if mult != highestMult {
					return false, fmt.Sprintf("Damage dealt multiplier for school %d returned %f, expected %f!", schoolIndex, mult, highestMult)
				}
				return true, ""
			})
	})

	t.Run("DamageTaken", func(t *testing.T) {
		SchoolMultiplierArrayHelper(t, caster, target, &target.PseudoStats.SchoolDamageTakenMultiplier,
			func(spell *Spell, schoolIndex stats.SchoolIndex) (bool, string) {
				spell.MultiSchoolUpdateModifiers(target)
				mult := spell.TargetDamageMultiplier(attackTable, false)
				if mult != highestMult {
					return false, fmt.Sprintf("Damage taken multiplier for school %d returned %f, expected %f!", schoolIndex, mult, highestMult)
				}
				return true, ""
			})
	})

	// TODO: Test for crit taken, it's currently not used anywhere.

	// t.Run("CritTaken", func(t *testing.T) {
	// 	SchoolMultiplierArrayHelper(t, caster, target, &target.PseudoStats.SchoolCritTakenMultiplier,
	// 		func(spell *Spell, schoolIndex stats.SchoolIndex) (bool, string) {
	// 			spell.MultiSchoolUpdateModifiers(target)
	// 			mult := spell.TargetDamageMultiplier(attackTable, false)
	// 			if mult != highestMult {
	// 				return false, fmt.Sprintf("Damage taken multiplier for school %d returned %f, expected %f!", schoolIndex, mult, highestMult)
	// 			}
	// 			return true, ""
	// 		})
	// })
}
