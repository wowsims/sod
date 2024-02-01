package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const FlameShockRanks = 6

// First entry is the base spell ID, second entry is the overload's spell ID
var FlameShockSpellId = [FlameShockRanks + 1]int32{0, 8050, 8052, 8053, 10447, 10448, 29228}
var FlameShockBaseDamage = [FlameShockRanks + 1]float64{0, 25, 51, 95, 164, 245, 292}
var FlameShockBaseDotDamage = [FlameShockRanks + 1]float64{0, 28, 48, 96, 168, 256, 320}
var FlameShockBaseSpellCoef = [FlameShockRanks + 1]float64{0, .134, .198, .214, .214, .214, .214}
var FlameShockDotSpellCoef = [FlameShockRanks + 1]float64{0, .134, .198, .214, .214, .214, .214}
var FlameShockManaCost = [FlameShockRanks + 1]float64{0, 55, 95, 160, 250, 345, 410}
var FlameShockLevel = [FlameShockRanks + 1]int{0, 10, 18, 28, 40, 52, 60}

func (shaman *Shaman) registerFlameShockSpell(shockTimer *core.Timer) {
	for i := 1; i <= FlameShockRanks; i++ {
		config := shaman.newFlameShockSpellConfig(shockTimer, i)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.FlameShock = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newFlameShockSpellConfig(shockTimer *core.Timer, rank int) core.SpellConfig {
	spellId := FlameShockSpellId[rank]
	baseDamage := FlameShockBaseDamage[rank]
	baseDotDamage := FlameShockBaseDotDamage[rank]
	baseSpellCoeff := FlameShockBaseSpellCoef[rank]
	dotSpellCoeff := FlameShockDotSpellCoef[rank]
	manaCost := FlameShockManaCost[rank]
	level := FlameShockLevel[rank]

	numTicks := 4
	tickDuration := time.Second * 3

	spell := shaman.newShockSpellConfig(
		core.ActionID{SpellID: spellId},
		core.SpellSchoolFire,
		manaCost,
		shockTimer,
	)

	spell.RequiredLevel = level
	spell.Rank = rank

	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		damage := baseDamage + baseSpellCoeff*spell.SpellPower()
		result := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		if result.Landed() {
			spell.Dot(target).NumberOfTicks = int32(numTicks)
			spell.Dot(target).Apply(sim)
		}
		spell.DealDamage(sim, result)
	}

	spell.Dot = core.DotConfig{
		Aura: core.Aura{
			Label: "FlameShock",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.LavaBurst.BonusCritRating += 100 * core.CritRatingPerCritChance
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.LavaBurst.BonusCritRating -= 100 * core.CritRatingPerCritChance
			},
		},
		NumberOfTicks: int32(numTicks),
		TickLength:    tickDuration,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = (baseDotDamage / float64(numTicks)) + dotSpellCoeff*dot.Spell.SpellPower()
			dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)

			if shaman.MoltenBlastAura != nil && result.Landed() && sim.RandomFloat("Molten Blast Reset") < shaman.MoltenBlastResetChance {
				shaman.MoltenBlast.CD.Reset()
			}
		},
	}

	return spell
}
