package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"time"
)

// Seal of Martyrdom is a spell consisting of:
// - A judgement that deals 70% weapon damage that is not normalised. Cannot miss or be dodged/blocked/parried.
// - An on-hit 100% chance proc that deals 40% *normalised* weapon damage.
// Both the on-hit and judgement are subject to weapon specialization talent modifiers as
// they both target melee defense.

func (paladin *Paladin) registerSealOfMartyrdom() {
	if !paladin.HasRune(proto.PaladinRune_RuneChestSealOfMartyrdom) {
		return
	}

	impSoRModifier := 1.0 + 0.03*float64(paladin.Talents.ImprovedSealOfRighteousness)

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 407803}, // Judgement of Martyrdom
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplier: 0.85 * paladin.getWeaponSpecializationModifier() * impSoRModifier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	onSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: 407799}, // Seal of Martyrdom
		SpellSchool:   core.SpellSchoolHoly,
		DefenseType:   core.DefenseTypeMelee,
		ProcMask:      core.ProcMaskMeleeMHSpecial,
		Flags:         core.SpellFlagMeleeMetrics,
		RequiredLevel: 1,

		DamageMultiplier: 0.5 * paladin.getWeaponSpecializationModifier() * impSoRModifier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	paladin.SealOfMartyrdomAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Martyrdom",
		ActionID: core.ActionID{SpellID: int32(proto.PaladinRune_RuneChestSealOfMartyrdom)},
		Duration: time.Second * 30,

		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit | core.ProcMaskProc) {
				onSwingProc.Cast(sim, result.Target)
			}
		},
	})

	aura := paladin.SealOfMartyrdomAura

	paladin.SealOfMartyrdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    aura.ActionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   paladin.BaseMana*0.04 - paladin.GetLibramSealCostReduction(),
			Multiplier: 1 - 0.03*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.ApplySeal(aura, onJudgementProc, sim)
		},
	})
}
