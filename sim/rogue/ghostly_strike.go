package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core/stats"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerGhostlyStrikeSpell() {
	if !rogue.Talents.GhostlyStrike {
		return
	}

	ghostlyStrikeAura := rogue.NewTemporaryStatsAura("Ghostly Strike Buff", core.ActionID{SpellID: 14278}, stats.Stats{stats.Dodge: 15 * core.DodgeRatingPerDodgeChance}, time.Second*7)

	rogue.GhostlyStrike = rogue.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_RogueGhostlyStrike,
		ActionID:    ghostlyStrikeAura.ActionID,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       rogue.builderFlags(),
		EnergyCost: core.EnergyCostOptions{
			Cost:   40.0,
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

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: 1.25,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			ghostlyStrikeAura.Activate(sim)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
