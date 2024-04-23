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
		Flags:       core.SpellFlagAPL,

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
			return paladin.CurrentSeal.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Seal of Command requires this spell to act as its intermediary dummy,
			// rolling on the spell hit table. If it succeeds, the actual Judgement of Command rolls on the
			// melee special attack crit/hit table, necessitating two discrete spells.
			// All other judgements are cast directly.
			if paladin.CurrentJudgement.SpellCode == SpellCode_PaladinJudgementOfCommand {
				spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			} else {
				paladin.CurrentJudgement.Cast(sim, target)
			}

			paladin.CurrentSealExpiration = sim.CurrentTime
			paladin.CurrentSeal.Deactivate(sim)
		},
	})
}
