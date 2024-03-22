package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// Crusader Strike is an ap-normalised instant attack that has a weapon damage % modifier with a 0.75 coefficient.
// It also returns 5% of the paladin's maximum mana when cast, regardless of the ability being negated.
// As of 27/02/24 it deals holy school damage, but otherwise behaves like a melee attack.

func (paladin *Paladin) registerCrusaderStrikeSpell() {
	if !paladin.HasRune(proto.PaladinRune_RuneHandsCrusaderStrike) {
		return
	}
	actionID := core.ActionID{SpellID: 407676}
	manaMetrics := paladin.NewManaMetrics(actionID)

	paladin.CrusaderStrike = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // cs is on phys gcd, which cannot be hasted
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 0.75 * paladin.getWeaponSpecializationModifier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			paladin.AddMana(sim, 0.05*paladin.MaxMana(), manaMetrics)
		},
	})
}
