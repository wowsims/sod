package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerShurikenTossSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneShurikenToss) {
		return
	}

	hasJustAFleshWound := rogue.HasRune(proto.RogueRune_RuneJustAFleshWound)

	results := make([]*core.SpellResult, min(5, rogue.Env.GetNumTargets()))

	rogue.ShurikenToss = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.RogueRune_RuneShurikenToss)},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeRanged,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       SpellFlagBuilder | SpellFlagCarnage | core.SpellFlagMeleeMetrics | core.SpellFlagAPL, // not affected by Cold Blood

		EnergyCost: core.EnergyCostOptions{
			Cost:   30,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 20,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: core.TernaryFloat64(hasJustAFleshWound, 1.5, 1),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.MeleeAttackPower() * 0.50
			var combopoints int32 = 0

			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
				combopoints++
			}

			if results[0].Landed() {
				rogue.AddComboPoints(sim, combopoints, spell.ComboPointMetrics())
			}
		},
	})
}
