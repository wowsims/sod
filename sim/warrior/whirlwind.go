package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) registerWhirlwindSpell() {
	if warrior.Level < 36 {
		return
	}

	hasConsumedByRageRune := warrior.HasRune(proto.WarriorRune_RuneConsumedByRage)

	warrior.WhirlwindMH = warrior.newWhirlwindHitSpell(true)
	if hasConsumedByRageRune {
		warrior.WhirlwindOH = warrior.newWhirlwindHitSpell(false)
	}

	warrior.Whirlwind = warrior.RegisterSpell(BerserkerStance, core.SpellConfig{
		SpellCode:   SpellCode_WarriorWhirlwind,
		ActionID:    core.ActionID{SpellID: 1680},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagAPL | SpellFlagBloodSurge,

		RageCost: core.RageCostOptions{
			Cost: 25 - warrior.FocusedRageDiscount,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				warrior.WhirlwindMH.Cast(sim, aoeTarget)
				if warrior.AutoAttacks.IsDualWielding && warrior.WhirlwindOH != nil && warrior.IsEnraged() {
					warrior.WhirlwindOH.Cast(sim, aoeTarget)
				}
			}
		},
	})
}

func (warrior *Warrior) newWhirlwindHitSpell(isMH bool) *WarriorSpell {
	procMask := core.ProcMaskMeleeSpecial
	damageFunc := warrior.MHNormalizedWeaponDamage
	if !isMH {
		procMask = core.ProcMaskMeleeOHSpecial
		damageFunc = warrior.OHNormalizedWeaponDamage
	}

	return warrior.RegisterSpell(AnyStance, core.SpellConfig{
		SpellCode:   core.Ternary(isMH, SpellCode_WarriorWhirlwindMH, SpellCode_WarriorWhirlwindOH),
		ActionID:    core.ActionID{SpellID: 1680}.WithTag(int32(core.Ternary(isMH, 1, 2))),
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1.25,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := damageFunc(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}
