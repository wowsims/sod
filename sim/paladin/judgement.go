package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) canJudgement(sim *core.Simulation, _ *core.Unit) bool {
	return paladin.CurrentSeal != nil && paladin.CurrentSeal.IsActive()
}

func (paladin *Paladin) registerJudgementSpell() {
	// Judgement functions as a dummy spell in vanilla.
	// It rolls on the spell hit table and can only miss or hit.
	// Individual seals have their own effects that this spell triggers,
	// that are handled in the implementations of the seal auras.
	paladin.Judgement = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20271},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagPrimaryJudgement | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.06,
			Multiplier: 1 - 0.03*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: (time.Second * 10) - (time.Second * time.Duration(paladin.Talents.ImprovedJudgement)),
			},
		},
		ExtraCastCondition: paladin.canJudgement,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// The judgement dummy spell only rolls spell hit in classic.
			// Subsequent judgement effects from seals have their own outcomes.
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			paladin.CurrentSealExpiration = sim.CurrentTime
			paladin.CurrentSeal.Deactivate(sim)
		},
	})
}
