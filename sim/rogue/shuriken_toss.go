package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ShurikenTossSpellID = int32(proto.RogueRune_RuneShurikenToss)

func (rogue *Rogue) registerShurikenTossSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneShurikenToss) {
		return
	}
	results := make([]*core.SpellResult, 5)

	numHits := min(5, rogue.Env.GetNumTargets())

	rogue.ShurikenToss = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: ShurikenTossSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeRanged,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   30,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.MeleeAttackPower() * 0.15

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				results[hitIndex] = spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeRangedHitAndCrit)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			if results[0].Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			}
		},
	})
}
