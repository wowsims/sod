package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (mage *Mage) registerDeepFreezeSpell() {
	if !mage.HasRune(proto.MageRune_RuneHelmDeepFreeze) {
		return
	}

	level := float64(mage.Level)
	baseCalc := (13.828124 + 0.018012*level + 0.044141*level*level)
	baseDamageLow := baseCalc * 4.62
	baseDamageHigh := baseCalc * 5.38
	spellCoeff := 2.5
	cooldown := time.Second * 30
	manaCost := .12

	hasFingersOfFrostRune := mage.HasRune(proto.MageRune_RuneChestFingersOfFrost)

	mage.DeepFreeze = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.MageRune_RuneHelmDeepFreeze)},
		SpellSchool: core.SpellSchoolFrost,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hasFingersOfFrostRune && mage.FingersOfFrostAura.IsActive()
		},
	})
}
