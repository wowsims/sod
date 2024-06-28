package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getRainOfFireBaseConfig(rank int) core.SpellConfig {
	hasLakeOfFireRune := warlock.HasRune(proto.WarlockRune_RuneChestLakeOfFire)

	spellId := [5]int32{0, 5740, 6219, 11677, 11678}[rank]
	spellCoeff := [5]float64{0, 0.083, 0.083, 0.083, 0.083}[rank]
	baseDamage := [5]float64{0, 42, 92, 155, 226}[rank]
	manaCost := [5]float64{0, 295, 605, 885, 1185}[rank]
	level := [5]int{0, 20, 34, 46, 58}[rank]

	flags := core.SpellFlagAPL | core.SpellFlagResetAttackSwing
	if !hasLakeOfFireRune {
		flags |= core.SpellFlagChanneled
	}

	config := core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         flags,
		RequiredLevel: level,
		Rank:          rank,

		BonusCritRating: float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,

		CritDamageBonus: warlock.ruin(),

		DamageMultiplierAdditive: 1 + 0.02*float64(warlock.Talents.Emberstorm),
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "RainOfFire-" + warlock.Label + strconv.Itoa(rank),
			},
			NumberOfTicks:    4,
			TickLength:       time.Second * 2,
			BonusCoefficient: spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
				}

			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	}

	if hasLakeOfFireRune {
		config.Cast.CD = core.Cooldown{
			Timer:    warlock.NewTimer(),
			Duration: time.Second * 8,
		}
	}

	return config
}

func (warlock *Warlock) registerRainOfFireSpell() {
	maxRank := 4

	for i := 1; i <= maxRank; i++ {
		config := warlock.getRainOfFireBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.RainOfFire = warlock.GetOrRegisterSpell(config)
		}
	}
}
