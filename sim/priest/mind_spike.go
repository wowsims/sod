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
	level := float64(priest.GetCharacter().Level)
	baseDamage := (9.456667 + 0.635108*level + 0.039063*level*level)
	// 2024-02-22 tuning 10% buff
	baseDamageLow := baseDamage * 1.11 * 1.1
	baseDamageHigh := baseDamage * 1.29 * 1.1
	spellCoeff := .429
	manaCost := .06
	castTime := time.Millisecond * 1500

	priest.MindSpikeAuras = priest.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		return priest.newMindSpikeAura(unit)
	})

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.PriestRune_RuneWaistMindSpike)},
		SpellSchool: core.SpellSchoolFrost | core.SpellSchoolShadow,
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

		BonusHitRating:   priest.shadowHitModifier(),
		DamageMultiplier: priest.forceOfWillDamageModifier() * priest.darknessDamageModifier(),
		BonusCritRating:  priest.forceOfWillCritRating(),
		CritMultiplier:   priest.DefaultSpellCritMultiplier(),
		ThreatMultiplier: priest.shadowThreatModifier(),

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			damage := (baseDamageLow+baseDamageHigh)/2 + spellCoeff*spell.SpellDamage()*priest.MindBlastModifier
			return spell.CalcDamage(sim, target, damage, spell.OutcomeExpectedMagicHitAndCrit)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

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
