package core

import (
	"math"
	"testing"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/simsignals"
	"github.com/wowsims/sod/sim/core/stats"
)

func Test_PartialResistsVsPlayer(t *testing.T) {
	attacker := &Unit{
		Type:  EnemyUnit,
		Level: 63,
		stats: stats.Stats{},
	}
	defender := &Unit{
		Type:  PlayerUnit,
		Level: 60,
		stats: stats.Stats{},
	}

	attackTable := NewAttackTable(attacker, defender, nil)

	sim := NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{},
		Encounter:  &proto.Encounter{},
		Raid:       &proto.Raid{},
	}, simsignals.CreateSignals())

	spell := &Spell{
		SpellSchool: SpellSchoolFire,
	}

	for resist := 0; resist < 5_000; resist += 1 {
		defender.stats[stats.FireResistance] = float64(resist)

		threshold00, threshold25, threshold50 := attackTable.Defender.partialResistRollThresholds(spell, attackTable.Attacker, false)
		thresholds := [4]float64{threshold00, threshold25, threshold50, 0.0}

		var cumulativeChance float64
		var resultingAr float64
		for bin, th := range thresholds {
			chance := max(min(1.0-th-cumulativeChance, 1.0), 0.0)
			resultingAr += chance * 0.25 * float64(bin)
			cumulativeChance += chance
			if cumulativeChance >= 1 {
				break
			}
		}

		resistanceScore := attackTable.Defender.resistCoeff(spell, attackTable.Attacker, false, false)
		expectedAr := 0.75*resistanceScore - 3.0/16.0*max(0.0, resistanceScore-2.0/3.0)

		if math.Abs(resultingAr-expectedAr) > 1e-2 {
			t.Errorf("resist = %d, thresholds = (%.2f, %.2f, %.2f), resultingAr = %.2f%%, expectedAr = %.2f%%", resist, threshold00, threshold25, threshold50, resultingAr*100, expectedAr*100)
			return
		}

		const n = 10_000

		outcomes := make(map[HitOutcome]int, n)
		var totalDamage float64
		for iter := 0; iter < n; iter++ {
			result := SpellResult{
				Outcome: OutcomeHit,
				Damage:  1000,
			}

			result.applyResistances(sim, spell, false, attackTable)

			outcomes[result.Outcome]++
			totalDamage += result.Damage
		}

		if math.Abs(expectedAr-(1-totalDamage/float64(1000*n))) > 0.01 {
			t.Logf("after %d iterations, resist = %d, ar = %.2f%% vs. damage lost = %.2f%%, outcomes = %v\n", n, resist, expectedAr*100, 100-100*totalDamage/float64(1000*n), outcomes)
		}
	}
}

func GetChancesAndMitFromThresholds(t0 float64, t25 float64, t50 float64) (float64, float64, float64, float64, float64) {
	chance0 := 1 - t0
	chance25 := t0 - t25
	chance50 := t25 - t50
	chance75 := t50
	avgResist := chance25*0.25 + chance50*0.50 + chance75*0.75
	return avgResist, chance0, chance25, chance50, chance75
}

func CloseEnough(f1 float64, f2 float64, eps float64) bool {
	return math.Abs(f1-f2) < eps
}

func ResistanceCheck(t *testing.T, isDoT bool) {
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

	schoolMask := SpellSchoolFromIndex(stats.SchoolIndexNature)
	spell := &Spell{
		SpellSchool:       schoolMask,
		SchoolIndex:       schoolMask.GetSchoolIndex(),
		SchoolBaseIndices: schoolMask.GetBaseIndices(),
	}

	if isDoT {
		spell.Flags |= SpellFlagPureDot
	}

	maxResist := float64(attacker.Level) * 5.0

	// Check if coef is 0.08 (from +3 level based resist) at 0 res
	defender.stats[stats.NatureResistance] = 0
	coef := defender.resistCoeff(spell, attacker, false, isDoT)
	if coef != 0.08 {
		t.Errorf("Resist coef is %.3f at 0 resistance, but should be 0.08!", coef)
		return
	}

	// Check known value
	defender.stats[stats.NatureResistance] = 200
	expectedMitigation := 0.545
	expectedChances := []float64{0, 0.18, 0.46, 0.36}
	if isDoT {
		expectedMitigation = 0.11
		expectedChances = []float64{0.67, 0.24, 0.08, 0.01}
	}
	threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell, spell.Flags.Matches(SpellFlagPureDot))
	avgResist, chance0, chance25, chance50, chance75 := GetChancesAndMitFromThresholds(threshold00, threshold25, threshold50)
	if !CloseEnough(avgResist, expectedMitigation, 0.01) {
		t.Errorf("Avg mitigation %.3f at 200 resistance, but should be %.3f!", avgResist, expectedMitigation)
		return
	}
	if !CloseEnough(chance0, expectedChances[0], 0.01) ||
		!CloseEnough(chance25, expectedChances[1], 0.01) ||
		!CloseEnough(chance50, expectedChances[2], 0.01) ||
		!CloseEnough(chance75, expectedChances[3], 0.01) {
		t.Errorf("Bucket chances do not match known values at 200 resistance. Known %v, returned %v!", expectedChances, []float64{chance0, chance25, chance50, chance75})
		return
	}

	// Check various resistance values
	for resist := 0.0; resist < maxResist; resist += 1.0 {
		defender.stats[stats.NatureResistance] = resist

		resistanceUsed := resist
		if isDoT {
			resistanceUsed /= 10
		}

		// Expected values
		resistanceCap := float64(attacker.Level * 5)
		levelBased := float64(max(defender.Level-attacker.Level, 0)) * 0.02
		expectedCoef := min(1, resistanceUsed/resistanceCap+levelBased*1/0.75)
		expectedAvgMitigation := expectedCoef*0.75 - 3.0/16.0*max(0, expectedCoef-2.0/3.0)

		// Check if coef is correct to begin with
		resistCoef := defender.resistCoeff(spell, attacker, false, isDoT)
		if math.Abs(resistCoef-expectedCoef) > 0.001 {
			t.Errorf("Resist coef is %.3f but expected %.3f at resistance %f", resistCoef, expectedCoef, resistanceUsed)
			return
		}

		// Check breakpoints
		threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell, spell.Flags.Matches(SpellFlagPureDot))
		avgResist, _, _, _, _ := GetChancesAndMitFromThresholds(threshold00, threshold25, threshold50)
		if math.Abs(avgResist-expectedAvgMitigation) > 0.005 {
			t.Errorf("resist = %.2f, thresholds = %f, resultingAr = %.2f%%, expectedAr = %.2f%%", resistanceUsed, threshold00, avgResist, expectedAvgMitigation)
			return
		}
	}
}

func Test_ResistsVsBoss(t *testing.T) {
	t.Run("Direct", func(t *testing.T) { ResistanceCheck(t, false) })
	t.Run("DoT", func(t *testing.T) { ResistanceCheck(t, true) })
}

func Test_ResistBinary(t *testing.T) {
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

	schoolMask := SpellSchoolFromIndex(stats.SchoolIndexNature)
	spell := &Spell{
		Flags:             SpellFlagBinary,
		SpellSchool:       schoolMask,
		SchoolIndex:       schoolMask.GetSchoolIndex(),
		SchoolBaseIndices: schoolMask.GetBaseIndices(),
	}

	// Check if coef is 0.0 at 0 resistance, binary spells do not get level based resistance!
	defender.stats[stats.NatureResistance] = 0
	coef := defender.resistCoeff(spell, attacker, true, false)
	if coef != 0.0 {
		t.Errorf("Resist coef is %.3f at 0 resistance for binary spell, but should be 0.0!", coef)
		return
	}

	// Should not partial resist
	dmgMult, outcome := spell.ResistanceMultiplier(nil, false, attackTable)
	if dmgMult != 1 || outcome != OutcomeEmpty {
		t.Errorf("ResistanceMultiplier for binary spell did not return mult=1 and empty outcome, got %.3f and outcome %d!", dmgMult, outcome)
		return
	}

	// Hit chance
	tests := [][]float64{
		{0.0, 1},
		{100.0, 0.75},
		{200.0, 0.5},
		{300.0, 0.25},
	}
	for _, test := range tests {
		resistance := test[0]
		defender.stats[stats.NatureResistance] = resistance
		expectedResult := test[1]
		result := attackTable.GetBinaryHitChance(spell)
		if !CloseEnough(result, expectedResult, 0.000001) {
			t.Errorf("Binary hit chance result at %.0f resistance was %.3f, expected %.3f!", resistance, result, expectedResult)
			return
		}
	}
}
