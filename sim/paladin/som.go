package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// Seal of Martyrdom is a spell consisting of:
// - A judgement that deals 70% weapon damage that is not normalised. Cannot miss or be dodged/blocked/parried.
// - An on-hit 100% chance proc that deals 40% *normalised* weapon damage.
// Both the on-hit and judgement are subject to weapon specialization talent modifiers as
// they both target melee defense.

func (paladin *Paladin) registerSealOfMartyrdomSpellAndAura() {
	if !paladin.HasRune(proto.PaladinRune_RuneChestSealofMartyrdom) {
		return
	}

	multiplier := 1.0 + 0.03*float64(paladin.Talents.ImprovedSealOfRighteousness)

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 407803}, // Judgement of Righteousness.
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

		DamageMultiplier: 0.85 * paladin.getWeaponSpecializationModifier() * multiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	onSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: 407799},
		SpellSchool:   core.SpellSchoolHoly,
		DefenseType:   core.DefenseTypeMelee,
		ProcMask:      core.ProcMaskMeleeMHSpecial | core.ProcMaskSuppressedExtraAttackAura,
		Flags:         core.SpellFlagMeleeMetrics,
		RequiredLevel: 1,

		DamageMultiplier: 0.5 * paladin.getWeaponSpecializationModifier() * multiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	auraActionID := core.ActionID{SpellID: 407798}
	paladin.SealOfMartyrdomAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Martyrdom",
		Tag:      "Seal",
		ActionID: auraActionID,
		Duration: SealDuration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				onSwingProc.Cast(sim, result.Target)
			}
			if spell.ProcMask.Matches(core.ProcMaskProc) {
				onSwingProc.Cast(sim, result.Target)
			}
		},
	})
	// Necessary because of the mix of % base mana cost and flat reduction on the libram
	manaCost := paladin.BaseMana * 0.04
	if paladin.Ranged().ID == LibramOfBenediction {
		manaCost -= 10
	}
	aura := paladin.SealOfMartyrdomAura
	paladin.SealOfMartyrdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1.0 - (float64(paladin.Talents.Benediction) * 0.03),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			paladin.ApplySeal(aura, onJudgementProc, sim)
		},
	})
}
