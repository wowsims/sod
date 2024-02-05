package warlock

import (
	"math"

	"github.com/wowsims/sod/sim/core"
)

func (warlock *Warlock) getDarkPactConfig(rank int) core.SpellConfig {
	spellId := [4]int32{0, 18220, 18937, 18938}[rank]
	manaRestore := [4]float64{0, 150, 200, 250}[rank]
	level := [4]int{0, 0, 50, 60}[rank]

	actionID := core.ActionID{SpellID: spellId}
	manaMetrics := warlock.NewManaMetrics(actionID)
	petManaMetrics := warlock.Pet.NewManaMetrics(actionID)

	return core.SpellConfig{
		ActionID:      actionID,
		SpellSchool:   core.SpellSchoolShadow,
		ProcMask:      core.ProcMaskEmpty,
		Flags:         core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		FlatThreatBonus: 80,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			maxDrain := manaRestore
			actualDrain := math.Min(maxDrain, warlock.Pet.CurrentMana())

			warlock.Pet.SpendMana(sim, actualDrain, petManaMetrics)
			warlock.AddMana(sim, actualDrain, manaMetrics)
		},
	}
}

func (warlock *Warlock) registerDarkPactSpell() {
	if !warlock.Talents.DarkPact {
		return
	}

	maxRank := 3

	for i := 1; i <= maxRank; i++ {
		config := warlock.getDarkPactConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.DarkPact = warlock.GetOrRegisterSpell(config)
		}
	}
}
