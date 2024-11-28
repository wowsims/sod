package core

import (
	"fmt"
	"math"

	"github.com/wowsims/sod/sim/core/stats"
)

// This function should do 3 things:
//  1. Set the Outcome of the hit effect.
//  2. Update spell outcome metrics.
//  3. Modify the damage if necessary.
type OutcomeApplier func(sim *Simulation, result *SpellResult, attackTable *AttackTable)

func (spell *Spell) OutcomeAlwaysHit(_ *Simulation, result *SpellResult, _ *AttackTable) {
	result.Outcome = OutcomeHit
	spell.SpellMetrics[result.Target.UnitIndex].Hits++
}

// Hit without Hits++ counter
func (spell *Spell) OutcomeAlwaysHitNoHitCounter(_ *Simulation, result *SpellResult, _ *AttackTable) {
	result.Outcome = OutcomeHit
}

func (spell *Spell) OutcomeAlwaysMiss(_ *Simulation, result *SpellResult, _ *AttackTable) {
	result.Outcome = OutcomeMiss
	result.Damage = 0
	spell.SpellMetrics[result.Target.UnitIndex].Misses++
}

func (dot *Dot) OutcomeTick(_ *Simulation, result *SpellResult, _ *AttackTable) {
	isPartialResist := result.DidResist()
	result.Outcome = OutcomeHit
	dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
	if isPartialResist {
		dot.Spell.SpellMetrics[result.Target.UnitIndex].ResistedTicks++
	}
}

func (dot *Dot) OutcomeTickPhysicalCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	isPartialResist := result.DidResist()

	if dot.Spell.PhysicalCritCheck(sim, attackTable) {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritMultiplier(attackTable)
		dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
		if isPartialResist {
			dot.Spell.SpellMetrics[result.Target.UnitIndex].ResistedCritTicks++
		}
	} else {
		result.Outcome = OutcomeHit
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
		if isPartialResist {
			dot.Spell.SpellMetrics[result.Target.UnitIndex].ResistedTicks++
		}
	}
}

func (dot *Dot) OutcomeSnapshotCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	isPartialResist := result.DidResist()

	if sim.RandomFloat("Snapshot Crit Roll") < dot.SnapshotCritChance {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritMultiplier(attackTable)
		dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
		if isPartialResist {
			dot.Spell.SpellMetrics[result.Target.UnitIndex].ResistedCritTicks++
		}
	} else {
		result.Outcome = OutcomeHit
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
		if isPartialResist {
			dot.Spell.SpellMetrics[result.Target.UnitIndex].ResistedTicks++
		}
	}
}

func (dot *Dot) OutcomeMagicHitAndSnapshotCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.MagicHitCheck(sim, attackTable) {
		isPartialResist := result.DidResist()

		if sim.RandomFloat("Snapshot Crit Roll") < dot.SnapshotCritChance {
			result.Outcome = OutcomeCrit
			result.Damage *= dot.Spell.CritMultiplier(attackTable)
			dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
			if isPartialResist {
				dot.Spell.SpellMetrics[result.Target.UnitIndex].ResistedCritTicks++
			}
		} else {
			result.Outcome = OutcomeHit
			dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
			if isPartialResist {
				dot.Spell.SpellMetrics[result.Target.UnitIndex].ResistedTicks++
			}
		}
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}

func (spell *Spell) OutcomeMagicHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMagicHitAndCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMagicHitAndCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMagicHitAndCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMagicHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	isPartialResist := result.DidResist()
	if spell.MagicHitCheck(sim, attackTable) {
		if spell.MagicCritCheck(sim, result.Target) {
			result.Outcome = OutcomeCrit
			result.Damage *= spell.CritMultiplier(attackTable)
			if countHits {
				spell.SpellMetrics[result.Target.UnitIndex].Crits++
				if isPartialResist {
					spell.SpellMetrics[result.Target.UnitIndex].ResistedCrits++
				}
			}
		} else {
			result.Outcome = OutcomeHit
			if countHits {
				spell.SpellMetrics[result.Target.UnitIndex].Hits++
				if isPartialResist {
					spell.SpellMetrics[result.Target.UnitIndex].ResistedHits++
				}
			}
		}
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Misses++
		}
	}
}

func (spell *Spell) OutcomeMagicCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMagicCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMagicCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMagicCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMagicCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	isPartialResist := result.DidResist()

	if spell.MagicCritCheck(sim, result.Target) {
		result.Outcome = OutcomeCrit
		result.Damage *= spell.CritMultiplier(attackTable)
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
			if isPartialResist {
				spell.SpellMetrics[result.Target.UnitIndex].ResistedCrits++
			}
		}
	} else {
		result.Outcome = OutcomeHit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Hits++
			if isPartialResist {
				spell.SpellMetrics[result.Target.UnitIndex].ResistedHits++
			}
		}
	}
}

func (spell *Spell) OutcomeHealing(_ *Simulation, result *SpellResult, _ *AttackTable) {
	spell.outcomeHealing(nil, result, nil, true)
}
func (spell *Spell) OutcomeHealingNoHitCounter(_ *Simulation, result *SpellResult, _ *AttackTable) {
	spell.outcomeHealing(nil, result, nil, false)
}
func (spell *Spell) outcomeHealing(_ *Simulation, result *SpellResult, _ *AttackTable, countHits bool) {
	isPartialResist := result.DidResist()
	result.Outcome = OutcomeHit
	if countHits {
		spell.SpellMetrics[result.Target.UnitIndex].Hits++
		if isPartialResist {
			spell.SpellMetrics[result.Target.UnitIndex].ResistedHits++
		}
	}
}

func (spell *Spell) OutcomeHealingCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeHealingCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeHealingCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeHealingCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeHealingCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	isPartialResist := result.DidResist()

	if spell.HealingCritCheck(sim) {
		result.Outcome = OutcomeCrit
		result.Damage *= spell.CritMultiplier(attackTable)
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
			if isPartialResist {
				spell.SpellMetrics[result.Target.UnitIndex].ResistedCrits++
			}
		}
	} else {
		result.Outcome = OutcomeHit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Hits++
			if isPartialResist {
				spell.SpellMetrics[result.Target.UnitIndex].ResistedHits++
			}
		}
	}
}

func (spell *Spell) OutcomeTickMagicHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.MagicHitCheck(sim, attackTable) {
		result.Outcome = OutcomeHit
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
	}
}

func (spell *Spell) OutcomeMagicHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMagicHit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMagicHitNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMagicHit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMagicHit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	if spell.MagicHitCheck(sim, attackTable) {
		isPartialResist := result.DidResist()
		result.Outcome = OutcomeHit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Hits++
			if isPartialResist {
				spell.SpellMetrics[result.Target.UnitIndex].ResistedHits++
			}
		}
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Misses++
		}
	}
}

func (spell *Spell) OutcomeMeleeWhite(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWhite(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeWhiteNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWhite(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeWhite(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	glanceRoll := sim.RandomFloat("White Hit Glancing Penalty")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMiss(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableGlance(spell, attackTable, roll, &chance, glanceRoll, countHits) &&
			!result.applyAttackTableBlock(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableCrit(spell, attackTable, roll, &chance, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	} else {
		if !result.applyAttackTableMiss(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableGlance(spell, attackTable, roll, &chance, glanceRoll, countHits) &&
			!result.applyAttackTableCrit(spell, attackTable, roll, &chance, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) OutcomeMeleeSpecialHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialHit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeSpecialHitNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialHit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialHit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) OutcomeMeleeSpecialHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialHitAndCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeSpecialHitAndCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialHitAndCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance, countHits) {
			if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
				result.applyAttackTableBlock(spell, attackTable, roll, &chance, countHits)
			} else {
				if !result.applyAttackTableBlock(spell, attackTable, roll, &chance, countHits) {
					result.applyAttackTableHit(spell, countHits)
				}
			}
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) OutcomeMeleeWeaponSpecialHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWeaponSpecialHitAndCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeWeaponSpecialHitAndCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWeaponSpecialHitAndCrit(sim, result, attackTable, false)
}

// Like OutcomeMeleeSpecialHitAndCrit, but blocks prevent crits (all weapon damage based attacks).
func (spell *Spell) outcomeMeleeWeaponSpecialHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	if spell.Unit.PseudoStats.InFrontOfTarget {
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableBlock(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	} else {
		spell.outcomeMeleeSpecialHitAndCrit(sim, result, attackTable, countHits)
	}
}

// Outcome for counted melee abilities matching:
// ✓ Miss
// ✓ Block
// ✓ Dodge
// ✓ Parry
// X Crit
func (spell *Spell) OutcomeMeleeWeaponSpecialNoCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWeaponSpecialNoCrit(sim, result, attackTable, true)
}

// Outcome for un-counted melee abilities matching:
// ✓ Miss
// ✓ Block
// ✓ Dodge
// ✓ Parry
// X Crit
func (spell *Spell) OutcomeMeleeWeaponSpecialNoCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWeaponSpecialNoCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeWeaponSpecialNoCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableBlock(spell, attackTable, roll, &chance, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

// Outcome for counted melee abilities matching:
// ✓ Miss
// ✓ Block
// X Dodge
// X Parry
// ✓ Crit
func (spell *Spell) OutcomeMeleeSpecialNoDodgeParry(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialNoDodgeParry(sim, result, attackTable, true)
}

// Outcome for counted melee abilities matching:
// ✓ Miss
// ✓ Block
// X Dodge
// X Parry
// ✓ Crit
func (spell *Spell) OutcomeMeleeSpecialNoDodgeParryNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialNoDodgeParry(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialNoDodgeParry(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
		!result.applyAttackTableBlock(spell, attackTable, roll, &chance, countHits) &&
		!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}
}

// Outcome for counted melee abilities matching:
// ✓ Miss
// X Block
// X Dodge
// X Parry
// ✓ Crit
func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParry(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialNoBlockDodgeParry(sim, result, attackTable, true)
}

// Outcome for un-counted melee abilities matching:
// ✓ Miss
// X Block
// X Dodge
// X Parry
// ✓ Crit
func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParryNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialNoBlockDodgeParry(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialNoBlockDodgeParry(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
		!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}
}

// Outcome for counted melee abilities matching:
// ✓ Miss
// X Block
// X Dodge
// X Parry
// X Crit
func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParryNoCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialNoBlockDodgeParryNoCrit(sim, result, attackTable, true)
}

// Outcome for un-counted melee abilities matching:
// ✓ Miss
// X Block
// X Dodge
// X Parry
// X Crit
func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParryNoCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialNoBlockDodgeParryNoCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialNoBlockDodgeParryNoCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}
}

// Outcome for counted melee abilities matching:
// X Miss
// X Block
// X Dodge
// X Parry
// ✓ Crit
func (spell *Spell) OutcomeMeleeSpecialCritOnly(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialCritOnly(sim, result, attackTable, true)
}

// Outcome for un-counted melee abilities matching:
// X Miss
// X Block
// X Dodge
// X Parry
// ✓ Crit
func (spell *Spell) OutcomeMeleeSpecialCritOnlyNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialCritOnly(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialCritOnly(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	if !result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}
}

func (spell *Spell) OutcomeRangedHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeRangedHitNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeRangedHit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}
}

func (spell *Spell) OutcomeRangedHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHitAndCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeRangedHitAndCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHitAndCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeRangedHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if spell.Unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) {
			if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
				result.applyAttackTableBlock(spell, attackTable, roll, &chance, countHits)
			} else {
				if !result.applyAttackTableBlock(spell, attackTable, roll, &chance, countHits) {
					result.applyAttackTableHit(spell, countHits)
				}
			}
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (dot *Dot) OutcomeRangedHitAndCritSnapshot(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	dot.outcomeRangedHitAndCritSnapshot(sim, result, attackTable, true)
}
func (dot *Dot) OutcomeRangedHitAndCritSnapshotNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	dot.outcomeRangedHitAndCritSnapshot(sim, result, attackTable, false)
}
func (dot *Dot) outcomeRangedHitAndCritSnapshot(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if dot.Spell.Unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(dot.Spell, attackTable, roll, &chance, countHits) {
			if result.applyAttackTableCritSeparateRollSnapshot(sim, dot, attackTable, countHits) {
				result.applyAttackTableBlock(dot.Spell, attackTable, roll, &chance, countHits)
			} else {
				if !result.applyAttackTableBlock(dot.Spell, attackTable, roll, &chance, countHits) {
					result.applyAttackTableHit(dot.Spell, countHits)
				}
			}
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(dot.Spell, attackTable, roll, &chance, countHits) &&
			!result.applyAttackTableCritSeparateRollSnapshot(sim, dot, attackTable, countHits) {
			result.applyAttackTableHit(dot.Spell, countHits)
		}
	}
}

func (spell *Spell) OutcomeRangedHitAndCritNoBlock(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHitAndCritNoBlock(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeRangedHitAndCritNoBlockNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHitAndCritNoBlock(sim, result, attackTable, false)
}
func (spell *Spell) outcomeRangedHitAndCritNoBlock(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance, countHits) &&
		!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}
}

func (spell *Spell) OutcomeRangedCritOnly(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedCritOnly(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeRangedCritOnlyNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedCritOnly(sim, result, attackTable, false)
}
func (spell *Spell) outcomeRangedCritOnly(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	// Block already checks for this, but we can skip the RNG roll which is expensive.
	if spell.Unit.PseudoStats.InFrontOfTarget {
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
			result.applyAttackTableBlock(spell, attackTable, roll, &chance, countHits)
		} else {
			if !result.applyAttackTableBlock(spell, attackTable, roll, &chance, countHits) {
				result.applyAttackTableHit(spell, countHits)
			}
		}
	} else {
		if !result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) OutcomeEnemyMeleeWhite(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeEnemyMeleeWhite(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeEnemyMeleeWhiteNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeEnemyMeleeWhite(sim, result, attackTable, false)
}
func (spell *Spell) outcomeEnemyMeleeWhite(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("Enemy White Hit Table")
	chance := 0.0

	didHit := !result.applyEnemyAttackTableMiss(spell, attackTable, roll, &chance, countHits)
	if !result.Target.IsCasting(sim) {
		didHit = didHit &&
			!result.applyEnemyAttackTableDodge(spell, attackTable, roll, &chance, countHits) &&
			!result.applyEnemyAttackTableParry(spell, attackTable, roll, &chance, countHits) &&
			!result.applyEnemyAttackTableBlock(spell, attackTable, roll, &chance, countHits)
	}

	if didHit && !result.applyEnemyAttackTableCrit(spell, attackTable, roll, &chance, countHits) {
		if didHit && !result.applyEnemyAttackTableCrush(spell, attackTable, roll, &chance, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) fixedCritCheck(sim *Simulation, critChance float64) bool {
	return sim.RandomFloat("Fixed Crit Roll") < critChance
}

func (result *SpellResult) applyAttackTableMiss(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	missChance := attackTable.BaseMissChance - spell.PhysicalHitChance(attackTable)
	missChance = math.Round(missChance*1000) / 1000

	if spell.Unit.AutoAttacks.IsDualWielding && !spell.Unit.PseudoStats.DisableDWMissPenalty {
		missChance += 0.19
	}

	*chance = max(0, missChance)

	if roll < *chance {
		result.Outcome = OutcomeMiss
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Misses++
		}
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableMissNoDWPenalty(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	missChance := attackTable.BaseMissChance - spell.PhysicalHitChance(attackTable)
	*chance = max(0, missChance)

	if roll < *chance {
		result.Outcome = OutcomeMiss
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Misses++
		}
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableBlock(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	*chance += attackTable.BaseBlockChance
	if roll < *chance {
		isCrit := result.DidCrit()
		isPartialResist := result.DidResist()
		result.Outcome |= OutcomeBlock
		if countHits {
			if isCrit {
				spell.SpellMetrics[result.Target.UnitIndex].BlockedCrits++
				spell.SpellMetrics[result.Target.UnitIndex].Crits--
				if isPartialResist {
					spell.SpellMetrics[result.Target.UnitIndex].ResistedCrits--
				}
			} else {
				spell.SpellMetrics[result.Target.UnitIndex].Blocks++
			}
		}
		// Physical abilities tagged with "Completely Blocked" are fully blocked every time
		if spell.Flags.Matches(SpellFlagBinary) {
			result.Damage = 0
		} else {
			result.Damage = max(0, result.Damage-result.Target.BlockValue())
		}
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableDodge(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	// In SoD this works like crit or hit chance.
	expertiseDodgeReduction := attackTable.Attacker.stats[stats.Expertise] / 100

	*chance += max(0, attackTable.BaseDodgeChance-attackTable.Defender.PseudoStats.DodgeReduction-expertiseDodgeReduction)
	*chance = math.Round(*chance*1000) / 1000

	if roll < *chance {
		result.Outcome = OutcomeDodge
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Dodges++
		}
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableParry(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	// In SoD this works like crit or hit chance.
	expertiseParryReduction := attackTable.Attacker.stats[stats.Expertise] / 100

	*chance += max(0, attackTable.BaseParryChance-expertiseParryReduction)

	if roll < *chance {
		result.Outcome = OutcomeParry
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Parries++
		}
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableGlance(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, glanceRoll float64, countHits bool) bool {
	*chance += attackTable.BaseGlanceChance

	if roll < *chance {
		result.Outcome = OutcomeGlance
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Glances++
		}
		result.Damage *= attackTable.GlanceMultiplierMin + glanceRoll*(attackTable.GlanceMultiplierMax-attackTable.GlanceMultiplierMin)
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableCrit(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	*chance += spell.PhysicalCritChance(attackTable)

	if roll < *chance {
		isPartialResist := result.DidResist()
		result.Outcome = OutcomeCrit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
			if isPartialResist {
				spell.SpellMetrics[result.Target.UnitIndex].ResistedCrits++
			}
		}
		result.Damage *= spell.CritMultiplier(attackTable)
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableCritSeparateRoll(sim *Simulation, spell *Spell, attackTable *AttackTable, countHits bool) bool {
	if spell.PhysicalCritCheck(sim, attackTable) {
		isPartialResist := result.DidResist()
		result.Outcome = OutcomeCrit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
			if isPartialResist {
				spell.SpellMetrics[result.Target.UnitIndex].ResistedCrits++
			}
		}
		result.Damage *= spell.CritMultiplier(attackTable)
		return true
	}
	return false
}
func (result *SpellResult) applyAttackTableCritSeparateRollSnapshot(sim *Simulation, dot *Dot, attackTable *AttackTable, countHits bool) bool {
	if sim.RandomFloat("Physical Crit Roll") < dot.SnapshotCritChance {
		isPartialResist := result.DidResist()
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritMultiplier(attackTable)
		if countHits {
			dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
			if isPartialResist {
				dot.Spell.SpellMetrics[result.Target.UnitIndex].ResistedCritTicks++
			}
		}
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableHit(spell *Spell, countHits bool) {
	isPartialResist := result.DidResist()
	result.Outcome = OutcomeHit

	if countHits {
		spell.SpellMetrics[result.Target.UnitIndex].Hits++
		if isPartialResist {
			spell.SpellMetrics[result.Target.UnitIndex].ResistedHits++
		}
	}
}

func (result *SpellResult) applyEnemyAttackTableMiss(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	missChance := attackTable.BaseMissChance + spell.Unit.PseudoStats.IncreasedMissChance +
		result.Target.stats[stats.Defense]*DefenseRatingToChanceReduction
	if spell.Unit.AutoAttacks.IsDualWielding && !spell.Unit.PseudoStats.DisableDWMissPenalty {
		missChance += 0.19
	}
	*chance = max(0, missChance)

	if roll < *chance {
		result.Outcome = OutcomeMiss
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Misses++
		}
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableBlock(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	if !result.Target.PseudoStats.CanBlock || result.Target.PseudoStats.Stunned {
		return false
	}

	blockChance := attackTable.BaseBlockChance +
		result.Target.stats[stats.Block]/BlockRatingPerBlockChance/100
	*chance += max(0, blockChance)

	if roll < *chance {
		result.Outcome |= OutcomeBlock
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Blocks++
		}
		result.Damage = max(0, result.Damage-result.Target.BlockValue())
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableDodge(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	if result.Target.PseudoStats.Stunned {
		return false
	}

	dodgeChance := attackTable.BaseDodgeChance +
		result.Target.GetStat(stats.Dodge)/100
	*chance += max(0, dodgeChance)

	if roll < *chance {
		result.Outcome = OutcomeDodge
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Dodges++
		}
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableParry(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	if !result.Target.PseudoStats.CanParry || result.Target.PseudoStats.Stunned {
		return false
	}

	parryChance := attackTable.BaseParryChance +
		result.Target.GetStat(stats.Parry)/100
	*chance += max(0, parryChance)

	if roll < *chance {
		result.Outcome = OutcomeParry
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Parries++
		}
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableCrit(spell *Spell, at *AttackTable, roll float64, chance *float64, countHits bool) bool {
	// "Base Melee Crit" is set as part of AttackTable
	critChance := at.BaseCritChance + spell.BonusCritRating/100
	// Crit reduction from bonus Defense of target (Talent, Gear, etc)
	critChance -= result.Target.stats[stats.Defense] * DefenseRatingToChanceReduction
	// Crit chance reduction (Rune: Just a Flesh Wound, etc)
	critChance -= result.Target.PseudoStats.ReducedCritTakenChance
	*chance += max(0, critChance)

	if roll < *chance {
		isPartialResist := result.DidResist()
		result.Outcome = OutcomeCrit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
			if isPartialResist {
				spell.SpellMetrics[result.Target.UnitIndex].ResistedCrits++
			}
		}
		result.Damage *= 2
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableCrush(spell *Spell, at *AttackTable, roll float64, chance *float64, countHits bool) bool {
	if !at.Attacker.PseudoStats.CanCrush {
		return false
	}

	crushChance := at.BaseCrushChance
	*chance += max(0, crushChance)

	if roll < *chance {
		result.Outcome = OutcomeCrush
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crushes++
		}
		result.Damage *= 1.5
		return true
	}
	return false
}

func (spell *Spell) OutcomeExpectedTick(_ *Simulation, _ *SpellResult, _ *AttackTable) {
	// result.Damage *= 1
}
func (spell *Spell) OutcomeExpectedMagicAlwaysHit(_ *Simulation, _ *SpellResult, _ *AttackTable) {
	// result.Damage *= 1
}
func (spell *Spell) OutcomeExpectedMagicHit(_ *Simulation, result *SpellResult, attackTable *AttackTable) {
	averageMultiplier := 1.0
	averageMultiplier -= spell.SpellChanceToMiss(attackTable)

	result.Damage *= averageMultiplier
}

func (spell *Spell) OutcomeExpectedMagicCrit(_ *Simulation, result *SpellResult, attackTable *AttackTable) {
	averageMultiplier := 1.0
	averageMultiplier += spell.SpellCritChance(result.Target) * (spell.CritMultiplier(attackTable) - 1)

	result.Damage *= averageMultiplier
}

func (spell *Spell) OutcomeExpectedMagicHitAndCrit(_ *Simulation, result *SpellResult, attackTable *AttackTable) {
	averageMultiplier := 1.0
	averageMultiplier -= spell.SpellChanceToMiss(attackTable)
	averageMultiplier += averageMultiplier * spell.SpellCritChance(result.Target) * (spell.CritMultiplier(attackTable) - 1)

	result.Damage *= averageMultiplier
}

func (dot *Dot) OutcomeExpectedMagicSnapshotCrit(_ *Simulation, result *SpellResult, attackTable *AttackTable) {
	averageMultiplier := 1.0
	averageMultiplier += dot.SnapshotCritChance * (dot.Spell.CritMultiplier(attackTable) - 1)

	result.Damage *= averageMultiplier
}

// CritMultiplier() returns the damage multiplier for critical strikes, based on CritDamageBonus and DefenseType.
// https://web.archive.org/web/20081014064638/http://elitistjerks.com/f31/t12595-relentless_earthstorm_diamond_-_melee_only/p4/
// https://github.com/TheGroxEmpire/TBC_DPS_Warrior_Sim/issues/30
func (spell *Spell) CritMultiplier(at *AttackTable) float64 {
	switch spell.DefenseType {
	case DefenseTypeNone:
		panic(fmt.Sprintf("using CritMultiplier() for spellID %d which has no DefenseType", spell.SpellID))
	case DefenseTypeMagic:
		return 1 + (1.5*at.CritMultiplier-1)*spell.CritDamageBonus
	default:
		return 1 + (2.0*at.CritMultiplier*at.Attacker.PseudoStats.MeleeCritMultiplier-1)*spell.CritDamageBonus
	}
}
