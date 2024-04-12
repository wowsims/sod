package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (mage *Mage) registerSpellfrostBoltSpell() {
	if !mage.HasRune(proto.MageRune_RuneBeltSpellfrostBolt) {
		return
	}

	baseDamageLow := mage.baseRuneAbilityDamage() * 3.04
	baseDamageHigh := mage.baseRuneAbilityDamage() * 3.55
	spellCoeff := .814
	castTime := time.Millisecond * 2500
	manaCost := .12

	mage.SpellfrostBolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: int32(proto.MageRune_RuneBeltSpellfrostBolt)},
		SpellCode:    SpellCode_MageSpellfrostBolt,
		SpellSchool:  core.SpellSchoolArcane | core.SpellSchoolFrost,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | SpellFlagChillSpell | core.SpellFlagAPL,
		MissileSpeed: 28,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: castTime,
				GCD:      core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
