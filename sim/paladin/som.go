package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// Seal of Martyrdom is a spell consisting of:
// - A judgement that deals 85% weapon damage that is not normalised. Cannot miss or be dodged/blocked/parried.
// - An on-hit 100% chance proc that deals 50% *normalised* weapon damage.
// Both the on-hit and judgement are subject to weapon specialization talent modifiers as
// they both target melee defense.

func (paladin *Paladin) registerSealOfMartyrdom() {
	manaMetrics := paladin.NewManaMetrics(core.ActionID{SpellID: 407802}) // SoM's mana restore

	judgeSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 407803},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlag_RV,

		DamageMultiplier: 0.85 * paladin.getWeaponSpecializationModifier() * paladin.improvedSoR(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	procSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: 407799},
		SpellSchool:   core.SpellSchoolHoly,
		DefenseType:   core.DefenseTypeMelee,
		ProcMask:      core.ProcMaskMeleeMHSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagSuppressWeaponProcs,
		RequiredLevel: 1,

		DamageMultiplier: 0.5 * paladin.getWeaponSpecializationModifier() * paladin.improvedSoR(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			// damages the paladin for 10% of rawDamage, then adds 133% of that for everyone in the raid
			paladin.AddMana(sim, result.RawDamage()*0.1*1.33, manaMetrics)
		},
	})

	aura := paladin.RegisterAura(core.Aura{
		Label:    "Seal of Martyrdom" + paladin.Label,
		ActionID: core.ActionID{SpellID: int32(proto.PaladinRune_RuneUtilitySealOfMartyrdom)},
		Duration: time.Second * 30,

		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit | core.ProcMaskProc) {
				procSpell.Cast(sim, result.Target)
			}
		},
	})

	paladin.auraSoM = aura

	paladin.sealOfMartyrdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    aura.ActionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost: paladin.BaseMana*0.04 - paladin.getLibramSealCostReduction(),
			Multiplier: paladin.benediction(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.applySeal(aura, judgeSpell, sim)
		},
	})

	paladin.spellJoM = judgeSpell
}
