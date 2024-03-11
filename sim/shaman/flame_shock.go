package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const FlameShockRanks = 6

var FlameShockSpellId = [FlameShockRanks + 1]int32{0, 8050, 8052, 8053, 10447, 10448, 29228}
var FlameShockBaseDamage = [FlameShockRanks + 1]float64{0, 25, 51, 95, 164, 245, 292}
var FlameShockBaseDotDamage = [FlameShockRanks + 1]float64{0, 28, 48, 96, 168, 256, 320}
var FlameShockBaseSpellCoef = [FlameShockRanks + 1]float64{0, .134, .198, .214, .214, .214, .214}
var FlameShockDotSpellCoef = [FlameShockRanks + 1]float64{0, .063, .093, .1, .1, .1, .1}
var FlameShockManaCost = [FlameShockRanks + 1]float64{0, 55, 95, 160, 250, 345, 410}
var FlameShockLevel = [FlameShockRanks + 1]int{0, 10, 18, 28, 40, 52, 60}

func (shaman *Shaman) registerFlameShockSpell(shockTimer *core.Timer) {
	shaman.FlameShock = make([]*core.Spell, FlameShockRanks+1)

	for rank := 1; rank <= FlameShockRanks; rank++ {
		config := shaman.newFlameShockSpellConfig(rank, shockTimer)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.FlameShock[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newFlameShockSpellConfig(rank int, shockTimer *core.Timer) core.SpellConfig {
	numTicks := 4
	tickDuration := time.Second * 3

	spellId := FlameShockSpellId[rank]
	baseDamage := FlameShockBaseDamage[rank]
	baseDotDamage := FlameShockBaseDotDamage[rank] / float64(numTicks)
	baseSpellCoeff := FlameShockBaseSpellCoef[rank]
	dotSpellCoeff := FlameShockDotSpellCoef[rank]
	manaCost := FlameShockManaCost[rank]
	level := FlameShockLevel[rank]

	if shaman.Ranged().ID == TotemInvigoratingFlame {
		manaCost -= 10
	}

	spell := shaman.newShockSpellConfig(
		core.ActionID{SpellID: spellId},
		core.SpellSchoolFire,
		manaCost,
		shockTimer,
	)

	spell.SpellCode = SpellCode_ShamanFlameShock
	spell.RequiredLevel = level
	spell.Rank = rank

	spell.Cast.IgnoreHaste = true

	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := baseDamage + baseSpellCoeff*spell.SpellDamage()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		if result.Landed() {
			spell.Dot(target).NumberOfTicks = int32(numTicks)
			spell.Dot(target).Apply(sim)

			if shaman.HasRune(proto.ShamanRune_RuneLegsAncestralGuidance) {
				shaman.lastFlameShockTarget = target
			}

			if shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge) && sim.RandomFloat("Power Surge Proc") < ShamanPowerSurgeProcChance {
				shaman.PowerSurgeAura.Activate(sim)
			}
		}
		spell.DealDamage(sim, result)
	}

	spell.Dot = core.DotConfig{
		Aura: core.Aura{
			Label: fmt.Sprintf("Flame Shock (Rank %d)", rank),
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if shaman.HasRune(proto.ShamanRune_RuneHandsLavaBurst) {
					shaman.LavaBurst.BonusCritRating += 100 * core.CritRatingPerCritChance
					if shaman.HasRune(proto.ShamanRune_RuneChestOverload) {
						shaman.LavaBurstOverload.BonusCritRating += 100 * core.CritRatingPerCritChance
					}
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if shaman.HasRune(proto.ShamanRune_RuneHandsLavaBurst) {
					shaman.LavaBurst.BonusCritRating -= 100 * core.CritRatingPerCritChance
					if shaman.HasRune(proto.ShamanRune_RuneChestOverload) {
						shaman.LavaBurstOverload.BonusCritRating -= 100 * core.CritRatingPerCritChance
					}
				}
			},
		},

		NumberOfTicks: int32(numTicks),
		TickLength:    tickDuration,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = baseDotDamage + dotSpellCoeff*dot.Spell.SpellDamage()
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
		},

		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)

			if shaman.HasRune(proto.ShamanRune_RuneHandsMoltenBlast) && result.Landed() && sim.RandomFloat("Molten Blast Reset") < ShamanMoltenBlastResetChance {
				shaman.MoltenBlast.CD.Reset()
			}

			if shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge) && sim.RandomFloat("Power Surge Proc") < ShamanPowerSurgeProcChance {
				shaman.PowerSurgeAura.Activate(sim)
			}
		},
	}

	return spell
}
