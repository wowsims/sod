package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const DrainSoulRanks = 4

func (warlock *Warlock) getDrainSoulBaseConfig(rank int) core.SpellConfig {
	hasSoulSiphonRune := warlock.HasRune(proto.WarlockRune_RuneCloakSoulSiphon)

	baseNumTicks := int32(5)
	numTicks := core.TernaryInt32(hasSoulSiphonRune, 15, baseNumTicks)
	tickLength := time.Second * time.Duration(core.TernaryInt32(hasSoulSiphonRune, 1, 3))

	spellId := [DrainSoulRanks + 1]int32{0, 1120, 8288, 8289, 11675}[rank]
	spellCoeff := [DrainSoulRanks + 1]float64{0, 0.063, 0.1, 0.1, 0.1}[rank]
	baseDamage := [DrainSoulRanks + 1]float64{0, 55, 155, 295, 455}[rank] / float64(baseNumTicks)
	manaCost := [DrainSoulRanks + 1]float64{0, 55, 125, 210, 290}[rank]
	level := [DrainSoulRanks + 1]int{0, 10, 24, 38, 52}[rank]

	return core.SpellConfig{
		SpellCode:   SpellCode_WarlockDrainSoul,
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagChanneled | core.SpellFlagResetAttackSwing | WarlockFlagAffliction | WarlockFlagHaunt,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DrainSoul-" + warlock.Label + strconv.Itoa(rank),
			},
			NumberOfTicks:    numTicks,
			TickLength:       tickLength,
			BonusCoefficient: spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)

				if hasSoulSiphonRune {
					dot.SnapshotAttackerMultiplier *= warlock.calcSoulSiphonMultiplier(target, sim.IsExecutePhase20())
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {

				dot := spell.Dot(target)
				dot.Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}
}

func (warlock *Warlock) registerDrainSoulSpell() {
	warlock.DrainSoul = make([]*core.Spell, 0)
	for rank := 1; rank <= DrainSoulRanks; rank++ {
		config := warlock.getDrainSoulBaseConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.DrainSoul = append(warlock.DrainSoul, warlock.GetOrRegisterSpell(config))
		}
	}
}
