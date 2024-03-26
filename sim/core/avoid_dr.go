package core

import (
	"github.com/wowsims/sod/sim/core/stats"
)

// Non-diminishing sources are added separately in spell outcome funcs
// TODO: Move these functions somewhere better

func (unit *Unit) GetDodgeChance() float64 {

	// undiminished Dodge % = D

	return unit.stats[stats.Dodge]
}

func (unit *Unit) GetParryChance() float64 {

	// undiminished Parry % = P

	return unit.stats[stats.Parry]

}

// TODO: Unused?
func (unit *Unit) GetMissChance() float64 {

	// undiminished Miss % = M

	return unit.stats[stats.Defense] * DefenseRatingToChanceReduction
}
