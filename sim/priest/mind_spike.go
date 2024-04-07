package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (priest *Priest) registerMindSpikeSpell() {
	if !priest.HasRune(proto.PriestRune_RuneWaistMindSpike) {
		return
	}

	priest.MindSpike = priest.GetOrRegisterSpell(priest.newMindSpikeSpellConfig())
}

func (priest *Priest) newMindSpikeSpellConfig() core.SpellConfig {
	// 2024-02-22 tuning 10% buff
	baseDamageLow := priest.baseRuneAbilityDamage() * 1.11 * 1.1
	baseDamageHigh := priest.baseRuneAbilityDamage() * 1.29 * 1.1
	spellCoeff := .429
	manaCost := .06
	castTime := time.Millisecond * 1500

	priest.MindSpikeAuras = priest.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		return priest.newMindSpikeAura(unit)
	})

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.PriestRune_RuneWaistMindSpike)},
		SpellSchool: core.SpellSchoolShadowfrost,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: castTime,
				GCD:      core.GCDDefault,
			},
		},

		BonusHitRating:  priest.shadowHitModifier(),
		BonusCritRating: priest.forceOfWillCritRating(),

		DamageMultiplier: priest.forceOfWillDamageModifier() * priest.darknessDamageModifier(),
		ThreatMultiplier: priest.shadowThreatModifier(),
		BonusCoefficient: spellCoeff,

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			damage := (baseDamageLow + baseDamageHigh) / 2
			spell.DamageMultiplier *= priest.MindBlastModifier
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeExpectedMagicHitAndCrit)
			spell.DamageMultiplier /= priest.MindBlastModifier
			return result
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)

			spell.DamageMultiplier *= priest.MindBlastModifier
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DamageMultiplier /= priest.MindBlastModifier

			if result.Landed() {
				priest.AddShadowWeavingStack(sim, target)
				priest.MindSpikeAuras.Get(target).Activate(sim)
				priest.MindSpikeAuras.Get(target).AddStack(sim)
			}

			spell.DealDamage(sim, result)
		},
	}
}

func (priest *Priest) newMindSpikeAura(unit *core.Unit) *core.Aura {
	return unit.RegisterAura(core.Aura{
		Label:     "Mind Spike",
		ActionID:  core.ActionID{SpellID: int32(proto.PriestRune_RuneWaistMindSpike)},
		Duration:  time.Second * 10,
		MaxStacks: 3,
	})
}
