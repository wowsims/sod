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
		ActionID:    core.ActionID{SpellID: 398196},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeRanged,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       rogue.builderFlags(),

		MissileSpeed: 40,

		EnergyCost: core.EnergyCostOptions{
			Cost:   []float64{25, 22, 20}[rogue.Talents.ImprovedSinisterStrike],
			Refund: 0.8,
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

		DamageMultiplier: []float64{1, 1.02, 1.04, 1.06}[rogue.Talents.Aggression],
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := rogue.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target)) +
				normalizedAmmoBonusDamage + spell.BonusWeaponDamage()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)

				if result.Landed() {
					rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				} else {
					spell.IssueRefund(sim)
				}
			})
		},
	})
}
