package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerMainGaucheSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneMainGauche) {
		return
	}

	rogue.MainGauche = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.RogueRune_RuneMainGauche)},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeOHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | core.SpellFlagAPL | SpellFlagColdBlooded,
		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier([]float64{35, 32, 30}[rogue.Talents.ImprovedSinisterStrike]),
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 20,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		// TODO: Add Parry aura for tanking
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			// Insert Parry aura activation here

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
