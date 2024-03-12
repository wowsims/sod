package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerPoisonedKnife() {
	if !rogue.HasRune(proto.RogueRune_RunePoisonedKnife) {
		return
	}

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
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 6,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.HasOHWeapon() && rogue.DistanceFromTarget >= 8
		},
		CastType: proto.CastType_CastTypeRanged,

		DamageMultiplier: 1 +
			0.02*float64(rogue.Talents.Aggression)*rogue.dwsMultiplier(),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,
		// Cannot Miss
		BonusHitRating: 100 * core.MeleeHitRatingPerHitChance,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()

			// Cannot Miss, Dodge, or Parry as per spell flags
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				// 100% application of OH poison (except for 1%? It can resist extremely rarely)
				if rogue.getImbueProcMask(proto.WeaponImbue_InstantPoison).Matches(core.ProcMaskMeleeOH) {
					rogue.InstantPoison[ShivProc].Cast(sim, target)
					return
				} else if rogue.getImbueProcMask(proto.WeaponImbue_DeadlyPoison).Matches(core.ProcMaskMeleeOH) {
					rogue.DeadlyPoison[ShivProc].Cast(sim, target)
					return
				} else if rogue.getImbueProcMask(proto.WeaponImbue_WoundPoison).Matches(core.ProcMaskMeleeOH) {
					rogue.WoundPoison[ShivProc].Cast(sim, target)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
