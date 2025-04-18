package warrior

import (
	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) registerHamstringSpell() {
	damage := map[int32]float64{
		25: 5,
		40: 18,
		50: 18,
		60: 45,
	}[warrior.Level]

	spellID := map[int32]int32{
		25: 1715,
		40: 7372,
		50: 7372,
		60: 27584,
	}[warrior.Level]

	spell_level := map[int32]int32{
		25: 8,
		40: 32,
		50: 32,
		60: 54,
	}[warrior.Level]

	warrior.Hamstring = warrior.RegisterSpell(BattleStance|BerserkerStance, core.SpellConfig{
		ClassSpellMask: ClassSpellMask_WarriorHamstring,
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagBinary | SpellFlagOffensive,

		RageCost: core.RageCostOptions{
			Cost:   10,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1.25,
		FlatThreatBonus:  1.25 * 2 * float64(spell_level),
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
