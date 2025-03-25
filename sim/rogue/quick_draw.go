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

	ammoBonusDamage := map[int32]float64{
		25: 7.5,
		40: 13,
		50: 15,
		60: 20,
	}[rogue.Level]
	normalizedAmmoBonusDamage := ammoBonusDamage / 2.8

	// Quick Draw applies a 50% slow, but bosses are immune

	rogue.QuickDraw = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: int32(proto.RogueRune_RuneQuickDraw)},
		ClassSpellMask: ClassSpellMask_RogueQuickdraw,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeRanged,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          rogue.builderFlags(),
		CastType:       proto.CastType_CastTypeRanged,
		MaxRange:       20,

		MissileSpeed: 40,

		EnergyCost: core.EnergyCostOptions{
			Cost:   25,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Millisecond * 500,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 10,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.Ranged().RangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeThrown ||
				rogue.Ranged().RangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeUnknown
		},

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := rogue.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target, false)) +
				normalizedAmmoBonusDamage

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)

				if result.Landed() {
					rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())
				} else {
					spell.IssueRefund(sim)
				}
			})
		},
	})
}
