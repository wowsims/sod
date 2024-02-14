package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const actionIDCS = 407676

func (paladin *Paladin) registerCrusaderStrikeSpell() {
	if !paladin.HasRune(proto.PaladinRune_RuneHandsCrusaderStrike) {
		return
	}
	actionID := core.ActionID{SpellID: 407676}
	manaMetrics := paladin.NewManaMetrics(actionID)

	paladin.CrusaderStrike = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
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
		BonusCritRating:  core.CritRatingPerCritChance,
		DamageMultiplier: 1.0,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.75*spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			paladin.AddMana(sim, 0.05*paladin.MaxMana(), manaMetrics)
		},
	})
}
