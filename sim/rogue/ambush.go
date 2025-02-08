package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerAmbushSpell() {
	flatDamageBonus := map[int32]float64{
		25: 28,
		40: 50,
		50: 92,
		60: 116,
	}[rogue.Level]

	spellID := map[int32]int32{
		25: 8676,
		40: 8725,
		50: 11268,
		60: 11269,
	}[rogue.Level]

	// waylay := rogue.HasRune(proto.RogueRune_RuneWaylay)
	hasCutthroatRune := rogue.HasRune(proto.RogueRune_RuneCutthroat)

	rogue.Ambush = rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueAmbush,
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          rogue.builderFlags(),

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
			if hasCutthroatRune && (rogue.CutthroatProcAura.IsActive() || rogue.IsStealthed()) {
				return true
			}
			return !rogue.PseudoStats.InFrontOfTarget && rogue.IsStealthed()
		},

		BonusCritRating:  15 * core.CritRatingPerCritChance * float64(rogue.Talents.ImprovedAmbush),
		DamageMultiplier: 2.5,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := (flatDamageBonus + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()))

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			if hasCutthroatRune && rogue.CutthroatProcAura.IsActive() {
				rogue.CutthroatProcAura.Deactivate(sim)
			}
			if result.Landed() {
				rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())
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
