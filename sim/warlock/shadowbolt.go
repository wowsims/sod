package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ShadowBoltRanks = 10

func (warlock *Warlock) getShadowBoltBaseConfig(rank int) core.SpellConfig {
	hasMarkOfChaosRune := warlock.HasRune(proto.WarlockRune_RuneBootsMarkOfChaos)

	spellCoeff := [ShadowBoltRanks + 1]float64{0, .14, .299, .56, .857, .857, .857, .857, .857, .857, .857}[rank]
	baseDamage := [ShadowBoltRanks + 1][]float64{{0}, {13, 18}, {26, 32}, {52, 61}, {92, 104}, {150, 170}, {213, 240}, {292, 327}, {373, 415}, {455, 507}, {482, 538}}[rank]
	spellId := [ShadowBoltRanks + 1]int32{0, 686, 695, 705, 1088, 1106, 7641, 11659, 11660, 11661, 25307}[rank]
	manaCost := [ShadowBoltRanks + 1]float64{0, 25, 40, 70, 110, 160, 210, 265, 315, 370, 380}[rank]
	level := [ShadowBoltRanks + 1]int{0, 1, 6, 12, 20, 28, 36, 44, 52, 60, 60}[rank]
	castTime := [ShadowBoltRanks + 1]int32{0, 1700, 2200, 2800, 3000, 3000, 3000, 3000, 3000, 3000, 3000}[rank]

	shadowboltVolley := warlock.HasRune(proto.WarlockRune_RuneHandsShadowBoltVolley)
	damageMulti := core.TernaryFloat64(shadowboltVolley, 0.95, 1.0)

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
				activeEffectMultiplier := 1.0

				if warlock.shadowBoltActiveEffectMultiplierPer > 0 && warlock.shadowBoltActiveEffectMultiplierMax > 0 {
					for _, spell := range warlock.DoTSpells {
						if spell.Dot(warlock.CurrentTarget).IsActive() {
							activeEffectMultiplier += warlock.shadowBoltActiveEffectMultiplierPer
						}
					}

					for _, spell := range warlock.DebuffSpells {
						if spell.RelatedAuras[0].Get(warlock.CurrentTarget).IsActive() {
							activeEffectMultiplier += warlock.shadowBoltActiveEffectMultiplierPer
						}
					}

					if hasMarkOfChaosRune && warlock.MarkOfChaosAuras.Get(warlock.CurrentTarget).IsActive() {
						activeEffectMultiplier += warlock.shadowBoltActiveEffectMultiplierPer
					}

					activeEffectMultiplier = min(warlock.shadowBoltActiveEffectMultiplierMax, activeEffectMultiplier)
				}

				spell.DamageMultiplier *= activeEffectMultiplier
				results[idx] = spell.CalcDamage(sim, target, sim.Roll(baseDamage[0], baseDamage[1]), spell.OutcomeMagicHitAndCrit)
				spell.DamageMultiplier /= activeEffectMultiplier

				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	}
}

func (warlock *Warlock) registerShadowBoltSpell() {
	warlock.ShadowBolt = make([]*core.Spell, 0)

	// TODO: AQ <=
	for rank := 1; rank < ShadowBoltRanks; rank++ {
		config := warlock.getShadowBoltBaseConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.ShadowBolt = append(warlock.ShadowBolt, warlock.GetOrRegisterSpell(config))
		}
	}
}
