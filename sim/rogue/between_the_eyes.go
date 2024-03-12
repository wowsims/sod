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
	flatDamage := rogue.RuneAbilityBaseDamage()
	comboDamageBonus := rogue.RuneAbilityDamagePerCombo()
	// TODO: Update damage variance
	damageVariance := map[int32]float64{
		25: 20,
		40: 44,
		50: 68,
		60: 108,
	}[rogue.Level]

	rogue.BetweenTheEyes = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: int32(proto.RogueRune_RuneBetweenTheEyes)},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | rogue.finisherFlags() | core.SpellFlagAPL,
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
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			comboPoints := rogue.ComboPoints()
			flatBaseDamage := flatDamage + comboDamageBonus*float64(comboPoints)

			// TODO: test combo point AP scaling
			baseDamage := sim.Roll(flatBaseDamage, flatBaseDamage+damageVariance) +
				0.03*float64(comboPoints)*spell.MeleeAttackPower() +
				spell.BonusWeaponDamage()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			if result.Landed() {
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
