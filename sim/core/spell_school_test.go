package core

import (
	"fmt"
	"math"
	"testing"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/simsignals"
	"github.com/wowsims/sod/sim/core/stats"
)

func BenchmarkMultiSchoolMultipliers(b *testing.B) {
	school := SpellSchoolFrost | SpellSchoolShadow

	multipliers := [stats.SchoolLen]float64{
		stats.SchoolIndexNone:     1,
		stats.SchoolIndexPhysical: 1,
		stats.SchoolIndexArcane:   1.1,
		stats.SchoolIndexFire:     1.1 * 1.15,
		stats.SchoolIndexFrost:    1.1 * 1.15,
		stats.SchoolIndexHoly:     1,
		stats.SchoolIndexNature:   1.2,
		stats.SchoolIndexShadow:   1.1 * 1.2,
	}

	indexes := []stats.SchoolIndex{stats.SchoolIndexFrost, stats.SchoolIndexShadow}

	mymax := func(a, b float64) float64 {
		if a < b {
			return b
		}
		return a
	}

	var dontOptimizeAway float64

	b.Run("index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var m float64
			for _, index := range indexes {
				m = mymax(m, multipliers[index])
			}
			dontOptimizeAway += m
		}
	})

	b.Run("school", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var m float64
			if school.Matches(SpellSchoolNone) {
				m = mymax(m, multipliers[stats.SchoolIndexNone])
			}
			if school.Matches(SpellSchoolPhysical) {
				m = mymax(m, multipliers[stats.SchoolIndexPhysical])
			}
			if school.Matches(SpellSchoolArcane) {
				m = mymax(m, multipliers[stats.SchoolIndexArcane])
			}
			if school.Matches(SpellSchoolFire) {
				m = mymax(m, multipliers[stats.SchoolIndexFire])
			}
			if school.Matches(SpellSchoolFrost) {
				m = mymax(m, multipliers[stats.SchoolIndexFrost])
			}
			if school.Matches(SpellSchoolHoly) {
				m = mymax(m, multipliers[stats.SchoolIndexHoly])
			}
			if school.Matches(SpellSchoolNature) {
				m = mymax(m, multipliers[stats.SchoolIndexNature])
			}
			if school.Matches(SpellSchoolShadow) {
				m = mymax(m, multipliers[stats.SchoolIndexShadow])
			}
			dontOptimizeAway += m
		}
	})

}

func Test_MultiSchoolIndexMapping(t *testing.T) {
	for si := stats.SchoolIndexPhysical; si < stats.SchoolLen; si++ {
		school := SpellSchoolFromIndex(si)
		if school == 0 {
			t.Errorf("No spell school for school index %d defined!", si)
			return
		}
	}
}

func getResiStatForSchool(schoolIndex stats.SchoolIndex) stats.Stat {
	switch schoolIndex {
	case stats.SchoolIndexPhysical:
		return stats.Armor
	case stats.SchoolIndexArcane:
		return stats.ArcaneResistance
	case stats.SchoolIndexFire:
		return stats.FireResistance
	case stats.SchoolIndexFrost:
		return stats.FrostResistance
	case stats.SchoolIndexNature:
		return stats.NatureResistance
	case stats.SchoolIndexShadow:
		return stats.ShadowResistance
	default:
		return stats.Armor
	}
}

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

	spell := &Spell{}

	for schoolIndex1 := stats.SchoolIndexArcane; schoolIndex1 < stats.SchoolLen; schoolIndex1++ {
		for schoolIndex2 := stats.SchoolIndexArcane; schoolIndex2 < stats.SchoolLen; schoolIndex2++ {
			if schoolIndex2 == schoolIndex1 {
				continue
			}

			schoolMask := SpellSchoolFromIndex(schoolIndex1) | SpellSchoolFromIndex(schoolIndex2)
			spell.SpellSchool = schoolMask
			spell.SchoolIndex = schoolMask.GetSchoolIndex()
			spell.SchoolBaseIndices = schoolMask.GetBaseIndices()

			baseIndices := spell.SchoolBaseIndices
			const lowestValue float64 = 50.0
			var lowestSchool stats.SchoolIndex
			isHoly := false

			for rev := 0; rev < 2; rev++ {
				if spell.SpellSchool.Matches(SpellSchoolHoly) {
					isHoly = true
				} else {
					if rev != 0 {
						indicesLen := len(baseIndices)
						lowestSchool = baseIndices[indicesLen-1]
						j := indicesLen - 1
						for i := 0; i < indicesLen; i++ {
							if baseIndices[j] != stats.SchoolIndexHoly {
								defender.stats[getResiStatForSchool(baseIndices[j])] = lowestValue + 25.0*float64(i)
							}
							j--
						}
					} else {
						lowestSchool = baseIndices[0]
						for i, baseIndex := range baseIndices {
							if baseIndex != stats.SchoolIndexHoly {
								defender.stats[getResiStatForSchool(baseIndex)] = lowestValue + 25.0*float64(i)
							}
						}
					}
				}

				resistance := 0.0

				if !isHoly {
					resistance = lowestValue

					// Make sure setup is right
					var lowestFound stats.SchoolIndex
					lowestValFound := 99999.0
					for _, checkIndex := range baseIndices {
						if checkIndex == stats.SchoolIndexHoly {
							lowestFound = checkIndex
							lowestValFound = 0
						} else {
							stat := getResiStatForSchool(checkIndex)
							if defender.GetStat(stat) < lowestValFound {
								lowestValFound = defender.GetStat(stat)
								lowestFound = checkIndex
							}
						}
					}
					if lowestFound != lowestSchool || lowestValFound != resistance {
						t.Errorf("Expected resist %d to be lowest with %f, but found %d at %f to be lowest resist!", lowestSchool, resistance, lowestFound, lowestValFound)
						return
					}
				}

				// Expected values
				resistanceCap := float64(attacker.Level * 5)
				levelBased := float64(max(defender.Level-attacker.Level, 0)) * 0.02
				expectedCoef := min(1, resistance/resistanceCap+levelBased*1/0.75)
				expectedAvgMitigation := expectedCoef*0.75 - 3.0/16.0*max(0, expectedCoef-2.0/3.0)

				// Check if coef is correct to begin with
				resistCoef := defender.resistCoeff(spell, attacker, false, false)
				if math.Abs(resistCoef-expectedCoef) > 0.001 {
					t.Errorf("Resist coef is %.3f but expected %.3f at resistance %f", resistCoef, expectedCoef, resistance)
					return
				}

				// Check breakpoints
				threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell, spell.Flags.Matches(SpellFlagPureDot))
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
}

func Test_MultiSchoolResistanceArmor(t *testing.T) {
	ss := SpellSchoolPhysical | SpellSchoolFire
	spell := &Spell{
		SpellSchool:       ss,
		SchoolIndex:       ss.GetSchoolIndex(),
		SchoolBaseIndices: ss.GetBaseIndices(),
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
	}, simsignals.CreateSignals())

	// Armor 100, resistances 0 => should use resistance
	defender.AddStat(stats.Armor, 100)

	mult, outcome := spell.ResistanceMultiplier(sim, false, attackTable)
	if outcome == OutcomeEmpty && mult < 1 {
		t.Errorf("Expected partial or full hit with armor %f and resistance %f", defender.GetStat(stats.Armor), defender.GetStat(stats.FireResistance))
		return
	}

	// Armor 100, resistances 300 => should use armor
	defender.AddStat(stats.FireResistance, 300)

	mult, outcome = spell.ResistanceMultiplier(sim, false, attackTable)
	if outcome != OutcomeEmpty || mult == 1 {
		t.Errorf("Expected empty outcome and mult < 1 with armor %f and resistance %f", defender.GetStat(stats.Armor), defender.GetStat(stats.FireResistance))
		return
	}

	// Armor 400, resistances 300 => should use resistance (no non-partial hits)
	defender.AddStat(stats.Armor, 300)

	_, outcome = spell.ResistanceMultiplier(sim, false, attackTable)
	if (outcome & OutcomePartial) == 0 {
		t.Errorf("Expected partial hit with armor %f and resistance %f", defender.GetStat(stats.Armor), defender.GetStat(stats.FireResistance))
		return
	}

	_, outcome = spell.ResistanceMultiplier(sim, true, attackTable)
	if (outcome & OutcomePartial) == 0 {
		t.Errorf("Expected partial hit for periodic with armor %f and resistance %f", defender.GetStat(stats.Armor), defender.GetStat(stats.FireResistance))
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
		Unit: caster,
	}

	for schoolIndex1 := stats.SchoolIndexArcane; schoolIndex1 < stats.SchoolLen; schoolIndex1++ {
		for schoolIndex2 := stats.SchoolIndexArcane; schoolIndex2 < stats.SchoolLen; schoolIndex2++ {
			if schoolIndex2 == schoolIndex1 {
				continue
			}

			schoolMask := SpellSchoolFromIndex(schoolIndex1) | SpellSchoolFromIndex(schoolIndex2)
			spell.SpellSchool = schoolMask
			spell.SchoolIndex = schoolMask.GetSchoolIndex()
			spell.SchoolBaseIndices = schoolMask.GetBaseIndices()

			baseIndices := spell.SchoolBaseIndices
			const highestValue float64 = 555.0
			var highestStat stats.Stat

			// Note: If powerStat == stats.SpellPower, then that means physical school, which is PseudoStats.BonusDamage!

			for rev := 0; rev < 2; rev++ {
				if rev != 0 {
					indexLen := len(baseIndices)
					highestStat = stats.ArcanePower + stats.Stat(baseIndices[indexLen-1]) - 2

					for i := indexLen - 1; i >= 0; i-- {
						powerStat := stats.ArcanePower + stats.Stat(baseIndices[i]) - 2
						if powerStat == stats.SpellPower {
							caster.PseudoStats.BonusPhysicalDamage = highestValue - 25.0*float64(indexLen-1-i)
						} else {
							caster.stats[powerStat] = highestValue - 25.0*float64(indexLen-1-i)
						}
					}
				} else {
					highestStat = stats.ArcanePower + stats.Stat(baseIndices[0]) - 2
					for i, baseIndex := range baseIndices {
						powerStat := stats.ArcanePower + stats.Stat(baseIndex) - 2
						if powerStat == stats.SpellPower {
							caster.PseudoStats.BonusPhysicalDamage = highestValue - 25.0*float64(i)
						} else {
							caster.stats[powerStat] = highestValue - 25.0*float64(i)
						}
					}
				}

				// Make sure setup is right
				var highestFound stats.Stat
				highestValFound := 0.0
				for _, baseIndex := range baseIndices {
					powerStat := stats.ArcanePower + stats.Stat(baseIndex) - 2
					if powerStat == stats.SpellPower {
						if caster.PseudoStats.BonusPhysicalDamage > highestValFound {
							highestValFound = caster.PseudoStats.BonusPhysicalDamage
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
					t.Errorf("Expected power %d to be highest with %f, but found %d at %f to be highest school power for school %d!", highestStat, highestValue, highestFound, highestValFound, schoolMask)
					return
				}

				power := spell.GetBonusDamage()
				if power != highestValue {
					t.Errorf("Expected %f to be highest power value found, but got %f for school %d!", highestValue, power, schoolMask)
					return
				}
			}
		}
	}
}

const highestMult int = 5

func SchoolMultiplierArrayHelper[T stats.SchoolValueArrayValues](t *testing.T, caster *Unit, target *Unit, multArray *stats.SchoolValueArray[T],
	testFunc func(spell *Spell, schoolMask SpellSchool, highest T) (bool, string)) {

	var highest T = T(highestMult)

	spell := &Spell{
		baseDamageMultiplier:                1,
		baseDamageMultiplierAdditivePct:     0,
		damageMultiplier:                    1,
		damageMultiplierAdditivePct:         0,
		impactDamageMultiplierAdditivePct:   0,
		periodicDamageMultiplierAdditivePct: 0,
		Unit:                                caster,
	}

	spell.updateImpactDamageMultiplier()
	spell.updatePeriodicDamageMultiplier()

	for schoolIndex1 := stats.SchoolIndexPhysical; schoolIndex1 < stats.SchoolLen; schoolIndex1++ {
		for schoolIndex2 := stats.SchoolIndexPhysical; schoolIndex2 < stats.SchoolLen; schoolIndex2++ {
			if schoolIndex2 == schoolIndex1 {
				continue
			}

			schoolMask := SpellSchoolFromIndex(schoolIndex1) | SpellSchoolFromIndex(schoolIndex2)
			spell.SpellSchool = schoolMask
			spell.SchoolIndex = schoolMask.GetSchoolIndex()
			spell.SchoolBaseIndices = schoolMask.GetBaseIndices()

			for i, baseIndex := range spell.SchoolBaseIndices {
				multArray[baseIndex] = highest - T(i/2)
			}

			ok, errMsg := testFunc(spell, schoolMask, highest)

			if !ok {
				t.Error(errMsg)
				return
			}
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
			func(spell *Spell, schoolMask SpellSchool, highest float64) (bool, string) {
				mult := spell.AttackerDamageMultiplier(attackTable, false)
				if mult != highest {
					return false, fmt.Sprintf("Damage dealt multiplier for school %d returned %f, expected %f!", schoolMask, mult, highest)
				}
				return true, ""
			})
	})

	t.Run("DamageTaken", func(t *testing.T) {
		SchoolMultiplierArrayHelper(t, caster, target, &target.PseudoStats.SchoolDamageTakenMultiplier,
			func(spell *Spell, schoolMask SpellSchool, highest float64) (bool, string) {
				mult := spell.TargetDamageMultiplier(attackTable, false)
				if mult != highest {
					return false, fmt.Sprintf("Damage taken multiplier for school %d returned %f, expected %f!", schoolMask, mult, highest)
				}
				return true, ""
			})
	})

	t.Run("CritChanceTaken", func(t *testing.T) {
		SchoolMultiplierArrayHelper(t, caster, target, &target.PseudoStats.SchoolCritTakenChance,
			func(spell *Spell, schoolMask SpellSchool, highest float64) (bool, string) {
				critChance := spell.SpellCritChance(target)
				if critChance != highest {
					return false, fmt.Sprintf("Crit chance taken for school %d returned %f, expected %f!", schoolMask, critChance, highest)
				}
				return true, ""
			})
	})

	t.Run("CostMultiplier", func(t *testing.T) {
		SchoolMultiplierArrayHelper(t, caster, target, &caster.PseudoStats.SchoolCostMultiplier,
			func(spell *Spell, schoolMask SpellSchool, highest int32) (bool, string) {
				costMod := spell.Unit.GetSchoolCostModifier(spell)
				if costMod != highest {
					return false, fmt.Sprintf("Cost mod for school %d returned %d, expected %d!", schoolMask, costMod, highest)
				}
				return true, ""
			})
	})
}
