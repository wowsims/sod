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
	canHitOffhand := hasConsumedByRageRune && warrior.AutoAttacks.IsDualWielding
	if canHitOffhand {
		warrior.WhirlwindOH = warrior.newWhirlwindHitSpell(false)
	}

	warrior.Whirlwind = warrior.RegisterSpell(BerserkerStance, core.SpellConfig{
		ClassSpellMask: ClassSpellMask_WarriorWhirlwind,
		ActionID:       core.ActionID{SpellID: 1680},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | SpellFlagOffensive,

		RageCost: core.RageCostOptions{
			Cost: 25,
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
				if canHitOffhand && warrior.IsEnraged() {
					warrior.WhirlwindOH.Cast(sim, aoeTarget)
				}
			}
		},
	})
}

func (warrior *Warrior) newWhirlwindHitSpell(isMH bool) *WarriorSpell {
	damageFunc := warrior.MHNormalizedWeaponDamage
	if !isMH {
		damageFunc = warrior.OHNormalizedWeaponDamage
	}

	return warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ClassSpellMask: core.Ternary(isMH, ClassSpellMask_WarriorWhirlwindMH, ClassSpellMask_WarriorWhirlwindOH),
		ActionID:       core.ActionID{SpellID: 1680}.WithTag(int32(core.Ternary(isMH, 1, 2))),
		SpellSchool:    core.SpellSchoolPhysical,
		CastType:       core.Ternary(isMH, proto.CastType_CastTypeMainHand, proto.CastType_CastTypeOffHand),
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.Ternary(isMH, core.ProcMaskMeleeMHSpecial, core.ProcMaskMeleeOHSpecial),
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: core.Ternary(isMH, 1.0, warrior.AutoAttacks.OHConfig().DamageMultiplier),
		ThreatMultiplier: 1.25,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := damageFunc(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}
