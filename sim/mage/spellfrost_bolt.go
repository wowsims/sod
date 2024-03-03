package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (mage *Mage) registerSpellfrostBolt() {
	if !mage.HasRune(proto.MageRune_RuneBeltSpellfrostBolt) {
		return
	}

	level := float64(mage.Level)
	baseCalc := 13.828124 + 0.018012*level + 0.044141*level*level
	baseDamageLow := baseCalc * 2.03
	baseDamageHigh := baseCalc * 2.37
	spellCoeff := .714
	castTime := time.Millisecond * 2500
	manaCost := .12

	mage.SpellfrostBolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:  core.ActionID{SpellID: int32(proto.MageRune_RuneBeltSpellfrostBolt)},
		SpellCode: SpellCode_MageSpellfrostBolt,
		// TODO: Multi-school spells
		SpellSchool:  core.SpellSchoolArcane, // | core.SpellSchoolFrost
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
		CritMultiplier:   mage.MageCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
