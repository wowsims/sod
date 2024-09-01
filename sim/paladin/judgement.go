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
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.06,
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
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Seal of Command requires this spell to act as its intermediary dummy,
			// rolling on the spell hit table. If it succeeds, the actual Judgement of Command rolls on the
			// melee special attack crit/hit table, necessitating two discrete spells.
			// All other judgements are cast directly.

			// Phase 1-3
			//if paladin.currentJudgement.SpellCode == SpellCode_PaladinJudgementOfCommand {
			//	spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			//} else {
			//	paladin.currentJudgement.Cast(sim, target)
			//}
			// paladin.currentSeal.Deactivate(sim)

			// Phase 4 - (Not tied to T1 6pc bonus) - Judge all Seals (2 possible without 6pc, 2 or 3 with 6pc TBD)
			for _, sealAura := range paladin.aurasSoC {
				if sealAura.IsActive() {
					spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
					sealAura.Deactivate(sim)
				}
			}

			for i, sealAura := range paladin.aurasSotC {
				if sealAura.IsActive() {
					paladin.spellsJotC[i].Cast(sim, target)
					sealAura.Deactivate(sim)
				}
			}

			for i, sealAura := range paladin.aurasSoR {
				if sealAura.IsActive() {
					paladin.spellsJoR[i].Cast(sim, target)
					sealAura.Deactivate(sim)
				}
			}

			if paladin.auraSoM.IsActive() {
				paladin.spellJoM.Cast(sim, target)
				paladin.auraSoM.Deactivate(sim)
			}

		},
	})
}
