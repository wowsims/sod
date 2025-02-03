package priest

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const HolyFireRanks = 8

var HolyFireSpellId = [HolyFireRanks + 1]int32{0, 14914, 15262, 15263, 15264, 15265, 15266, 15267, 15261}
var HolyFireBaseDamage = [HolyFireRanks + 1][]float64{{0}, {84, 104}, {97, 122}, {144, 178}, {173, 218}, {219, 273}, {259, 328}, {323, 406}, {355, 449}}
var HolyFireDotDamage = [HolyFireRanks + 1]float64{0, 30, 40, 55, 65, 85, 100, 125, 145}
var HolyFireSpellCoef = [HolyFireRanks + 1]float64{0, 0.123, 0.271, 0.554, 0.714, 0.714, 0.714, 0.714, 0.714}
var HolyFireManaCost = [HolyFireRanks + 1]float64{0, 85, 95, 125, 145, 170, 200, 230, 255}
var HolyFireLevel = [HolyFireRanks + 1]int{0, 20, 24, 30, 36, 42, 48, 54, 60}

func (priest *Priest) registerHolyFire() {
	priest.HolyFire = make([]*core.Spell, HolyFireRanks+1)

	for rank := 1; rank <= HolyFireRanks; rank++ {
		config := priest.getHolyFireConfig(rank)

		if config.RequiredLevel <= int(priest.Level) {
			priest.HolyFire[rank] = priest.GetOrRegisterSpell(config)
		}
	}
}

func (priest *Priest) getHolyFireConfig(rank int) core.SpellConfig {
	ticks := int32(5)

	spellId := HolyFireSpellId[rank]
	baseDamageLow := HolyFireBaseDamage[rank][0]
	baseDamageHigh := HolyFireBaseDamage[rank][1]
	dotDamage := HolyFireDotDamage[rank] / float64(ticks)
	manaCost := HolyFireManaCost[rank]
	level := HolyFireLevel[rank]

	directCoeff := 0.75
	dotCoeff := 0.05
	castTime := time.Millisecond * 3500

	hasDespairRune := priest.HasRune(proto.PriestRune_RuneBracersDespair)

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_PriestHolyFire,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolHoly,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime - time.Millisecond*100*time.Duration(priest.Talents.DivineFury),
			},
		},

		BonusCoefficient: directCoeff,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("Holy Fire (Rank %d)", rank),
			},

			NumberOfTicks:    ticks,
			TickLength:       time.Second * 2,
			BonusCoefficient: dotCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dotDamage, isRollover)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if hasDespairRune {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	}
}
