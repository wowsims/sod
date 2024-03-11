package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: 0.5s windup animation before activation
func (rogue *Rogue) registerQuickDrawSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneQuickDraw) {
		return
	}

	// Quick Draw applies a 50% slow, but bosses are immune

	rogue.QuickDraw = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 398196},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       SpellFlagBuilder | core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		EnergyCost: core.EnergyCostOptions{
			Cost:   []float64{25, 22, 20}[rogue.Talents.ImprovedSinisterStrike],
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 6,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.Ranged().RangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeThrown ||
				rogue.Ranged().RangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeUnknown
		},
		CastType: proto.CastType_CastTypeRanged,

		DamageMultiplier: 1,
		CritMultiplier:   rogue.RangedCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.RangedWeaponDamage(sim, spell.RangedAttackPower(target))

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
