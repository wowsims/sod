package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const RainOfFireRanks = 4

func (warlock *Warlock) getRainOfFireBaseConfig(rank int) core.SpellConfig {
	hasLakeOfFireRune := warlock.HasRune(proto.WarlockRune_RuneChestLakeOfFire)

	spellId := [RainOfFireRanks + 1]int32{0, 5740, 6219, 11677, 11678}[rank]
	spellCoeff := [RainOfFireRanks + 1]float64{0, 0.083, 0.083, 0.083, 0.083}[rank]
	baseDamage := [RainOfFireRanks + 1]float64{0, 42, 92, 155, 226}[rank]
	manaCost := [RainOfFireRanks + 1]float64{0, 295, 605, 885, 1185}[rank]
	level := [RainOfFireRanks + 1]int{0, 20, 34, 46, 58}[rank]

	flags := core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagDestruction
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

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
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
	warlock.RainOfFire = make([]*core.Spell, 0)
	for rank := 1; rank <= RainOfFireRanks; rank++ {
		config := warlock.getRainOfFireBaseConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.RainOfFire = append(warlock.RainOfFire, warlock.GetOrRegisterSpell(config))
		}
	}
}
