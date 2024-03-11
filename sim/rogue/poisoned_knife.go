package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerPoisonedKnife() {
	rogue.PoisonedKnife = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.RogueRune_RunePoisonedKnife)},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeOHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   []float64{25, 22, 20}[rogue.Talents.ImprovedSinisterStrike],
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.HasOHWeapon()
		},
		CastType: proto.CastType_CastTypeRanged,

		DamageMultiplier: 1 +
			0.02*float64(rogue.Talents.Aggression)*rogue.dwsMultiplier(),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,
		BonusHitRating:   100 * core.MeleeHitRatingPerHitChance,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()

			// Cannot Miss, Dodge, or Parry
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
