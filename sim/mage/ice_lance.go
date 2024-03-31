package mage

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: Classic review ice lance numbers on live
func (mage *Mage) registerIceLanceSpell() {
	if !mage.HasRune(proto.MageRune_RuneHandsIceLance) {
		return
	}

	baseDamageLow := mage.baseRuneAbilityDamage() * .55
	baseDamageHigh := mage.baseRuneAbilityDamage() * .65
	spellCoeff := .143
	manaCost := .08

	hasFingersOfFrostRune := mage.HasRune(proto.MageRune_RuneChestFingersOfFrost)

	mage.IceLance = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: int32(proto.MageRune_RuneHandsIceLance)},
		SpellSchool:  core.SpellSchoolFrost,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | core.SpellFlagAPL,
		MissileSpeed: 38,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			if hasFingersOfFrostRune && mage.FingersOfFrostAura.IsActive() {
				baseDamage *= 3
			}
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
