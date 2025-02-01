package warrior

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) registerExecuteSpell() {
	hasSuddenDeathRune := warrior.HasRune(proto.WarriorRune_RuneSuddenDeath)

	flatDamage := map[int32]float64{
		25: 125,
		40: 325,
		50: 450,
		60: 600,
	}[warrior.Level]

	convertedRageDamage := map[int32]float64{
		25: 3,
		40: 9,
		50: 12,
		60: 15,
	}[warrior.Level]

	spellID := map[int32]int32{
		25: 5308,
		40: 20660,
		50: 20661,
		60: 20662,
	}[warrior.Level]

	var rageMetrics *core.ResourceMetrics
	warrior.Execute = warrior.RegisterSpell(BattleStance|BerserkerStance, core.SpellConfig{
		ClassSpellMask: ClassSpellMask_WarriorExecute,
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagPassiveSpell | SpellFlagOffensive,

		RageCost: core.RageCostOptions{
			Cost:   15 - []float64{0, 2, 5}[warrior.Talents.ImprovedExecute],
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20() || (hasSuddenDeathRune && warrior.SuddenDeathAura.IsActive())
		},

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1.25,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			extraRage := spell.Unit.CurrentRage()
			warrior.SpendRage(sim, extraRage, rageMetrics)
			// We must count this rage event if the spell itself cost 0,
			// otherwise we could end up with 0 events even though rage was spent.
			if spell.Cost.GetCurrentCost() > 0 {
				rageMetrics.Events--
			}

			baseDamage := flatDamage + convertedRageDamage*(extraRage)

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
	rageMetrics = warrior.Execute.Cost.SpellCostFunctions.(*core.RageCost).ResourceMetrics
}
