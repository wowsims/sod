package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (paladin *Paladin) registerSealOfMartyrdomSpellAndAura() {

	if !paladin.HasRune(proto.PaladinRune_RuneChestSealofMartyrdom) {
		return
	}

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 407803}, // Judgement of Righteousness.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) * 0.7
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	onSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 407799},
		SpellSchool:      core.SpellSchoolHoly,
		ProcMask:         core.ProcMaskMeleeSpecial, // This needs figured out properly
		Flags:            core.SpellFlagMeleeMetrics,
		RequiredLevel:    1,
		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   paladin.MeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) * 0.4
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	// TODO: SoM also procs on a variety of weapon proc effects like Talwar's shadow bolt.
	// There is a brief icd preventing SoM and weapon procs chaining back and forth, they
	// can proc each other exactly once.
	// When the rates are figured out
	icd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Millisecond * 5,
	}

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
				if icd.IsReady(sim) {
					icd.Use(sim)
				}
				onSwingProc.Cast(sim, result.Target)
			}
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				onJudgementProc.Cast(sim, result.Target)
			}
		},
	})

	aura := paladin.SealOfMartyrdomAura
	paladin.SealOfMartyrdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   0.04,
			Multiplier: 1.0 - (float64(paladin.Talents.Benediction) * 0.03),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
