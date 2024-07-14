package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) registerShockwaveSpell() {
	if !warrior.HasRune(proto.WarriorRune_RuneShockwave) {
		return
	}

	apCoef := 0.50

	warrior.Shockwave = warrior.RegisterSpell(DefensiveStance, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.WarriorRune_RuneShockwave)},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeRanged,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 15 - warrior.FocusedRageDiscount,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: 20 * time.Second,
			},
		},

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := apCoef * spell.MeleeAttackPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				// Shockwave can miss and be blocked, but it can't be dodged or parried
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialNoDodgeParry)
			}
		},
	})
}
