package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerBackstabSpell() {
	flatDamageBonus := map[int32]float64{
		25: 48,
		40: 90,
		50: 135,
		// TODO: AQ
		60: 210,
		// 60: 225,
	}[rogue.Level]

	spellID := map[int32]int32{
		25: 2590,
		40: 8721,
		50: 11279,
		// TODO: AQ
		60: 11281,
		// 60: 25300
	}[rogue.Level]

	// waylay := rogue.HasRune(proto.RogueRune_RuneWaylay)
	hasCutthroatRune := rogue.HasRune(proto.RogueRune_RuneCutthroat)
	hasSlaughterRune := rogue.HasRune(proto.RogueRune_RuneSlaughterFromTheShadows)

	damageMultiplier := 1.5 * []float64{1, 1.04, 1.08, 1.12, 1.16, 1.2}[rogue.Talents.Opportunity]
	energyCost := 60.0

	if hasSlaughterRune {
		damageMultiplier *= SlaughterFromTheShadowsDamageMultiplier
		energyCost -= SlaughterFromTheShadowsCostReduction
	}

	rogue.Backstab = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       rogue.builderFlags(),

		EnergyCost: core.EnergyCostOptions{
			Cost:   energyCost,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if !rogue.HasDagger(core.MainHand) {
				return false
			}
			return hasCutthroatRune || !rogue.PseudoStats.InFrontOfTarget
		},

		BonusCritRating: 10 * core.CritRatingPerCritChance * float64(rogue.Talents.ImprovedBackstab),

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := (flatDamageBonus + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()))
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				if hasCutthroatRune {
					rogue.rollCutthroat(sim)
				}
				/** Currently does not apply to bosses due to being a slow
				if waylay {
					rogue.WaylayAuras.Get(target).Activate(sim)
				} */
			} else {
				spell.IssueRefund(sim)
			}
		},

		RelatedAuras: []core.AuraArray{rogue.WaylayAuras},
	})
}
