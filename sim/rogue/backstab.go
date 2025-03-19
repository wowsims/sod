package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerBackstabSpell() {
	flatDamageBonus := map[int32]float64{
		25: 32,
		40: 60,
		50: 90,
		60: core.TernaryFloat64(core.IncludeAQ, 150, 140),
	}[rogue.Level]

	spellID := map[int32]int32{
		25: 2590,
		40: 8721,
		50: 11279,
		60: core.TernaryInt32(core.IncludeAQ, 25300, 11281),
	}[rogue.Level]

	// waylay := rogue.HasRune(proto.RogueRune_RuneWaylay)
	hasCutthroatRune := rogue.HasRune(proto.RogueRune_RuneCutthroat)

	rogue.Backstab = rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueBackstab,
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          rogue.builderFlags(),
		MaxRange:       5,

		EnergyCost: core.EnergyCostOptions{
			Cost:   60,
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

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: 1.5,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := (flatDamageBonus + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()))
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())
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
