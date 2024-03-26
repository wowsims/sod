package core

import (
	"github.com/wowsims/sod/sim/core/stats"
)

// Could be in constants.go, but they won't be used anywhere else
// C values are divided by 100 so that we are working with 1% = 0.01
const Diminish_k_Druid = 0.972
const Diminish_k_Nondruid = 0.956
const Diminish_Cd_Druid = 116.890707 / 100
const Diminish_Cd_Nondruid = 88.129021 / 100
const Diminish_Cp = 47.003525 / 100
const Diminish_Cm = 16.000 / 100
const Diminish_kCd_Druid = (Diminish_k_Druid * Diminish_Cd_Druid)
const Diminish_kCd_Nondruid = (Diminish_k_Nondruid * Diminish_Cd_Nondruid)
const Diminish_kCp = (Diminish_k_Nondruid * Diminish_Cp)
const Diminish_kCm_Druid = (Diminish_k_Druid * Diminish_Cm)
const Diminish_kCm_Nondruid = (Diminish_k_Nondruid * Diminish_Cm)

// Non-diminishing sources are added separately in spell outcome funcs
// TODO: Move these functions somewhere better

func (unit *Unit) GetDodgeChance() float64 {

	// undiminished Dodge % = D

	return unit.stats[stats.Dodge]/100 +
		unit.stats[stats.Defense]*DefenseRatingToChanceReduction
}

func (unit *Unit) GetParryChance() float64 {

	// undiminished Parry % = P

	return unit.stats[stats.Parry]/100 +
		unit.stats[stats.Defense]*DefenseRatingToChanceReduction

}

func (unit *Unit) GetMissChance() float64 {

	// undiminished Miss % = M

	return unit.stats[stats.Defense] * DefenseRatingToChanceReduction
}
