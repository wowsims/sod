package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) applyUnfairAdvantage() {
	if !rogue.HasRune(proto.RogueRune_RuneUnfairAdvantage) {
		return
	}
	comboMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 432274})

	unfairAdvantage := rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 432274},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1.5, //added in AQ patch
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: Verify this should be normalized as the spell has 2 effects
			// one being normalized with 0 BasePoints and one being not normalized with 100 base points
			damage := rogue.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})

	icd := core.Cooldown{
		Timer:    rogue.NewTimer(),
		Duration: time.Second,
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Unfair Advantage Trigger",
		ActionID: core.ActionID{SpellID: int32(proto.RogueRune_RuneUnfairAdvantage)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeDodge|core.OutcomeParry) && icd.IsReady(sim) {
				unfairAdvantage.Cast(sim, rogue.CurrentTarget)
				rogue.AddComboPoints(sim, 1, rogue.CurrentTarget, comboMetrics)
				icd.Use(sim)
			}
		},
	})
}
