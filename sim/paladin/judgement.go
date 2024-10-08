package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerJudgement() {
	// Judgement functions as a dummy spell in vanilla.
	// It rolls on the spell hit table and can only miss or hit.
	// Individual seals have their own effects that this spell triggers,
	// that are handled in the implementations of the seal auras.
	paladin.judgement = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20271},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagCastTimeNoGCD,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.06,
			Multiplier: paladin.benediction(),
		},

		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * (10 - time.Duration(paladin.Talents.ImprovedJudgement)),
			},
		},
		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return paladin.currentSeal.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {

			// Phase 1-3
			//if paladin.currentJudgement.SpellCode == SpellCode_PaladinJudgementOfCommand {
			//	spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			//} else {
			//	paladin.currentJudgement.Cast(sim, target)
			//}
			// paladin.currentSeal.Deactivate(sim)

			// Phase 4 - (Double Judge not tied to T1 6pc bonus) - Judge all Seals (2 possible without 6pc, and with 6pc)
			// Random order of Judgements (Server Side - Indeterministic)

			// Phase 5 - (Double Judge now tied to T1 6pc bonus) - Judge all Seals (2 possible with 6pc)
			// Otherwise Judge Current Seal, or Previous Seal if two are active
			multipleSealsActive := false
			if paladin.prevSeal != nil && paladin.prevSeal.IsActive() {
				multipleSealsActive = true
			}

			if multipleSealsActive {
				paladin.castSpecificJudgement(sim, target, paladin.prevJudgement, paladin.prevSeal)

				if paladin.enableMultiJudge {
					paladin.castSpecificJudgement(sim, target, paladin.currentJudgement, paladin.currentSeal)
				}
			} else {
				paladin.castSpecificJudgement(sim, target, paladin.currentJudgement, paladin.currentSeal)
			}

		},
	})

	paladin.judgement.DefaultCast.GCD = 0
}

// Helper Function For casting Judgement
func (paladin *Paladin) castSpecificJudgement(sim *core.Simulation, target *core.Unit, judgementSpell *core.Spell, matchingSeal *core.Aura) {

	judgementSpell.Cast(sim, target)

	if paladin.consumeSealsOnJudge {
		matchingSeal.Deactivate(sim)
	}
}
