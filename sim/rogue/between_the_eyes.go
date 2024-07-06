package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerBetweenTheEyes() {
	if !rogue.HasRune(proto.RogueRune_RuneBetweenTheEyes) {
		return
	}

	flatDamage := rogue.baseRuneAbilityDamage()
	comboDamageBonus := rogue.baseRuneAbilityDamageCombo()

	rogue.BetweenTheEyes = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: int32(proto.RogueRune_RuneBetweenTheEyes)},
		SpellSchool:  core.SpellSchoolPhysical,
		DefenseType:  core.DefenseTypeRanged,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        rogue.finisherFlags(),
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:   35,
			Refund: 0,
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
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			comboPoints := rogue.ComboPoints()
			flatBaseDamage := flatDamage + comboDamageBonus*float64(comboPoints)
			variableDamage := sim.Roll(flatBaseDamage*0.53, flatBaseDamage*0.81)

			// TODO: test combo point AP scaling. Also, does BTE use Melee or Ranged Attack Power?
			baseDamage := variableDamage +
				0.03*float64(comboPoints)*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			if result.Landed() {
				rogue.SpendComboPoints(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
