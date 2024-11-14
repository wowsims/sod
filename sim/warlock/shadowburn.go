package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const ShadowburnRanks = 6

func (warlock *Warlock) registerShadowBurnBaseConfig(rank int) core.SpellConfig {
	spellId := [ShadowburnRanks + 1]int32{0, 17877, 18867, 18868, 18869, 18870, 18871}[rank]
	baseDamage := [ShadowburnRanks + 1][]float64{{0}, {91, 104}, {123, 140}, {196, 221}, {274, 307}, {365, 408}, {462, 514}}[rank]
	manaCost := [ShadowburnRanks + 1]float64{0, 105, 130, 190, 245, 305, 365}[rank]
	level := [ShadowburnRanks + 1]int{0, 15, 24, 32, 40, 48, 56}[rank]

	spellCoeff := 0.429

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellCode:     SpellCode_WarlockShadowburn,
		SpellSchool:   core.SpellSchoolShadow,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing | core.SpellFlagBinary | WarlockFlagDestruction,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(15),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1])
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}
}

func (warlock *Warlock) registerShadowBurnSpell() {
	if !warlock.Talents.Shadowburn {
		return
	}

	warlock.Shadowburn = make([]*core.Spell, 0)
	for rank := 1; rank <= ShadowburnRanks; rank++ {
		config := warlock.registerShadowBurnBaseConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.Shadowburn = append(warlock.Shadowburn, warlock.GetOrRegisterSpell(config))
		}
	}
}
