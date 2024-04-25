package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (druid *Druid) registerShredSpell() {
	damageMultiplier := 2.25

	flatDamageBonus := map[int32]float64{
		25: 24,
		40: 44,
		50: 64,
		60: 80,
	}[druid.Level]

	hasGoreRune := druid.HasRune(proto.DruidRune_RuneHelmGore)
	hasElunesFires := druid.HasRune(proto.DruidRune_RuneBracersElunesFires)

	if druid.Ranged().ID == IdolOfTheDream {
		damageMultiplier *= 1.02
		flatDamageBonus *= 1.02
	}

	// cp. https://www.wowhead.com/classic/spell=436895/s03-tuning-and-overrides-passive-druid
	damageMultiplier += 0.75 // multiplier +75%
	flatDamageBonus *= 0.75  // base damage -25%

	druid.Shred = druid.RegisterSpell(Cat, core.SpellConfig{
		SpellCode: SpellCode_DruidShred,
		ActionID: core.ActionID{SpellID: map[int32]int32{
			25: 5221,
			40: 8992,
			50: 9829,
			60: 9830,
		}[druid.Level]},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   60 - 6*float64(druid.Talents.ImprovedShred),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !druid.PseudoStats.InFrontOfTarget
		},

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatDamageBonus +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			modifier := 1.0
			if druid.BleedCategories.Get(target).AnyActive() {
				modifier = 1.3
			}

			spell.DamageMultiplier *= modifier
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.DamageMultiplier /= modifier

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())

				if hasGoreRune {
					druid.rollGoreCatReset(sim)
				}

				if hasElunesFires {
					druid.tryElunesFiresRipExtension(sim, target)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := flatDamageBonus + spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower())

			modifier := 1.0
			if druid.BleedCategories.Get(target).AnyActive() {
				modifier = 1.3
			}

			spell.DamageMultiplier *= modifier
			baseres := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			spell.DamageMultiplier /= modifier

			attackTable := spell.Unit.AttackTables[target.UnitIndex][spell.CastType]
			critChance := spell.PhysicalCritChance(attackTable)
			critMod := critChance * (spell.CritMultiplier(attackTable) - 1)

			baseres.Damage *= 1 + critMod

			return baseres
		},
	})
}

func (druid *Druid) CanShred() bool {
	return !druid.PseudoStats.InFrontOfTarget && druid.CurrentEnergy() >= druid.CurrentShredCost()
}

func (druid *Druid) CurrentShredCost() float64 {
	return druid.Shred.ApplyCostModifiers(druid.Shred.DefaultCast.Cost)
}
