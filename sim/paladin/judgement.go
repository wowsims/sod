package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) canJudgement(sim *core.Simulation) bool {
	return paladin.CurrentSeal != nil && paladin.CurrentSeal.IsActive() && paladin.Judgement.IsReady(sim)
}

func (paladin *Paladin) registerJudgementSpell() {
	// jowAuras := paladin.NewEnemyAuraArray(core.JudgementOfWisdomAura)

	paladin.Judgement = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20271},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty, // can proc TaJ itself and from seal
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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Primary Judgements cannot crit or be dodged, parried, or blocked-- only miss. (Unless target is a hunter.)
			// jow := jowAuras.Get(target)
			// if jow.IsActive() {
			// 	jow.Refresh(sim)
			// } else {
			// 	jow.Activate(sim)
			// }
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeRangedHit)
		},

		// RelatedAuras: []core.AuraArray{jowAuras},
	})
}
