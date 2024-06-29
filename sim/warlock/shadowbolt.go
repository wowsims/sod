package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getShadowBoltBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [11]float64{0, .14, .299, .56, .857, .857, .857, .857, .857, .857, .857}[rank]
	baseDamage := [11][]float64{{0}, {13, 18}, {26, 32}, {52, 61}, {92, 104}, {150, 170}, {213, 240}, {292, 327}, {373, 415}, {455, 507}, {482, 538}}[rank]
	spellId := [11]int32{0, 686, 695, 705, 1088, 1106, 7641, 11659, 11660, 11661, 25307}[rank]
	manaCost := [11]float64{0, 25, 40, 70, 110, 160, 210, 265, 315, 370, 380}[rank]
	level := [11]int{0, 1, 6, 12, 20, 28, 36, 44, 52, 60, 60}[rank]
	castTime := [11]int32{0, 1700, 2200, 2800, 3000, 3000, 3000, 3000, 3000, 3000, 3000}[rank]

	shadowboltVolley := warlock.HasRune(proto.WarlockRune_RuneHandsShadowBoltVolley)
	damageMulti := core.TernaryFloat64(shadowboltVolley, 0.8, 1.0)

	results := make([]*core.SpellResult, min(core.TernaryInt32(shadowboltVolley, 5, 1), warlock.Env.GetNumTargets()))

	return core.SpellConfig{
		SpellCode:     SpellCode_WarlockShadowBolt,
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagDestruction,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * time.Duration(castTime),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.MetamorphosisAura == nil || !warlock.MetamorphosisAura.IsActive()
		},

		DamageMultiplier: damageMulti,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				damage := sim.Roll(baseDamage[0], baseDamage[1])
				results[idx] = spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	}
}

func (warlock *Warlock) registerShadowBoltSpell() {
	maxRank := 10

	for i := 1; i <= maxRank; i++ {
		config := warlock.getShadowBoltBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.ShadowBolt = warlock.GetOrRegisterSpell(config)
		}
	}
}
