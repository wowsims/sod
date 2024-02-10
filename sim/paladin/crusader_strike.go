package paladin

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (paladin *Paladin) registerCrusaderStrikeSpell() {
	if !paladin.HasRune(proto.PaladinRune_RuneHandsCrusaderStrike) {
		fmt.Println(
			"vetoing CS", proto.PaladinRune_RuneHandsCrusaderStrike,
		)
		return
	}
	fmt.Println(
		"CS rune found, registering CS",
	)
	// bonusDmg := core.TernaryFloat64(paladin.Ranged().ID == 31033, 36, 0) + // Libram of Righteous Power
	// 	core.TernaryFloat64(paladin.Ranged().ID == 40191, 79, 0) // Libram of Radiance

	// jowAuras := paladin.NewEnemyAuraArray(core.JudgementOfWisdomAura)
	// jolAuras := paladin.NewEnemyAuraArray(core.JudgementOfLightAura)

	paladin.CrusaderStrike = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 407676}, // 35395  407676
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		// ManaCost: core.ManaCostOptions{
		// 	BaseCost: 0.05,
		// },
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

		// BonusCritRating: core.TernaryFloat64(paladin.HasSetBonus(ItemSetAegisBattlegear, 4), 10, 0) * core.CritRatingPerCritChance,
		BonusCritRating: core.CritRatingPerCritChance,
		// DamageMultiplierAdditive: 1 +
		// 	paladin.getTalentSanctityOfBattleBonus() +
		// 	paladin.getTalentTheArtOfWarBonus() +
		// 	paladin.getItemSetGladiatorsVindicationBonusGloves(),
		DamageMultiplier: 0.75,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			// jowAura := jowAuras.Get(target)
			// if jowAura.IsActive() {
			// 	jowAura.Refresh(sim)
			// }

			// jolAura := jolAuras.Get(target)
			// if jolAura.IsActive() {
			// 	jolAura.Refresh(sim)
			// }
		},
	})
}
