package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
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
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

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

			thisJudgement := 0 // Seals Judged this Judgement

			// Phase 4 - (Not tied to T1 6pc bonus) - Judge all Seals (2 possible without 6pc, 2 or 3 with 6pc TBD)
			for _, sealAura := range paladin.aurasSoC {
				if sealAura.IsActive() {
					spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
					if !paladin.t2Judgement2pc {
						sealAura.Deactivate(sim)
					}
					thisJudgement |= 1
				}
			}

			for i, sealAura := range paladin.aurasSotC {
				if sealAura.IsActive() {
					paladin.spellsJotC[i].Cast(sim, target)
					if !paladin.t2Judgement2pc {
						sealAura.Deactivate(sim)
					}
					thisJudgement |= (1 << 1)
				}
			}

			for i, sealAura := range paladin.aurasSoR {
				if sealAura.IsActive() {
					paladin.spellsJoR[i].Cast(sim, target)
					if !paladin.t2Judgement2pc {
						sealAura.Deactivate(sim)
					}
					thisJudgement |= (1 << 2)
				}
			}

			if paladin.auraSoM.IsActive() {
				paladin.spellsJoM[0].Cast(sim, target)
				if !paladin.t2Judgement2pc {
					paladin.auraSoM.Deactivate(sim)
				}
				thisJudgement |= (1 << 3)
			}

			if paladin.t2Judgement4pc && thisJudgement != paladin.lastJudgement {
				// 4 pieces: The cooldown on your Judgement is instantly reset if used on a different Seal than your last Judgement.
				paladin.judgement.CD.Reset()
				paladin.lastJudgement = thisJudgement
			}

			if paladin.t2Judgement6pcAura != nil {
				// 6 pieces: Your Judgement grants 1% increased Holy damage for 8 sec, stacking up to 5 times.
				if sim.Log != nil {
					paladin.Log(sim, "Test Message! 12345")
				}
				paladin.t2Judgement6pcAura.Activate(sim)
				paladin.t2Judgement6pcAura.AddStack(sim)
			}
		},
	})
}

func (paladin *Paladin) enableT2Judgement6pc() {
	if paladin.t2Judgement6pcAura != nil {
		return
	}

	paladin.t2Judgement6pcAura = paladin.GetOrRegisterAura(core.Aura{
		Label:     "Swift Judgement",
		ActionID:  core.ActionID{SpellID: 467530},
		Duration:  time.Second * 8,
		MaxStacks: 5,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.t2Judgement6pc {
				aura.Activate(sim)
				aura.SetStacks(sim, 0)
			}
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= (1.0 + (float64(oldStacks) * 0.01))
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= (1.0 + (float64(newStacks) * 0.01))
		},
	})
}

func (paladin *Paladin) enableT2Judgement2pc() {
	if paladin.t2Judgement2pcApplied {
		return
	}

	for _, judgeSpells := range paladin.allJudgeSpells {
		for _, judgeRankSpell := range judgeSpells {
			judgeRankSpell.DamageMultiplier *= 1.2
		}
	}

	paladin.t2Judgement2pcApplied = true
}
