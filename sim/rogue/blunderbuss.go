package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: 10 yd range
func (rogue *Rogue) registerBlunderbussSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneBlunderbuss) {
		return
	}

	results := make([]*core.SpellResult, min(4, rogue.Env.GetNumTargets()))

	rogue.Blunderbuss = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 436564},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeRanged,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   20,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 15,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := rogue.baseRuneAbilityDamage() + spell.MeleeAttackPower()*0.48
			baseDamageVariable := sim.Roll(baseDamage*192, baseDamage*288)

			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, baseDamageVariable, spell.OutcomeRangedHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}
