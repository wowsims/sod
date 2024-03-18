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

	hasDeadlyBrew := rogue.HasRune(proto.RogueRune_RuneDeadlyBrew)

	rogue.PoisonedKnife = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.RogueRune_RunePoisonedKnife)},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
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

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: []float64{1, 1.02, 1.04, 1.06}[rogue.Talents.Aggression] * rogue.dwsMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()

			// Cannot Miss, Dodge, or Parry as per spell flags
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				// 100% application of OH poison (except for 1%? It can resist extremely rarely)
				switch rogue.Consumes.OffHandImbue {
				case proto.WeaponImbue_InstantPoison:
					rogue.InstantPoison[ShivProc].Cast(sim, target)
				case proto.WeaponImbue_DeadlyPoison:
					rogue.DeadlyPoison[ShivProc].Cast(sim, target)
				case proto.WeaponImbue_WoundPoison:
					rogue.WoundPoison[ShivProc].Cast(sim, target)
				default:
					if hasDeadlyBrew {
						rogue.InstantPoison[NormalProc].Cast(sim, target)
					}
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
