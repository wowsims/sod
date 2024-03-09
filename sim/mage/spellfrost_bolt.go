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
	// 2024-03-05 tuning SFB +50% base damage and same spell coeff as max rank Frostbolt
	baseCalc := (13.828124 + 0.018012*level + 0.044141*level*level)
	baseDamageLow := baseCalc * 3.04
	baseDamageHigh := baseCalc * 3.55
	spellCoeff := .814
	castTime := time.Millisecond * 2500
	manaCost := .12

	mage.SpellfrostBolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: int32(proto.MageRune_RuneBeltSpellfrostBolt)},
		SpellCode:    SpellCode_MageSpellfrostBolt,
		SpellSchool:  core.SpellSchoolSpellFrost,
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
