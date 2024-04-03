package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (druid *Druid) registerShredSpell() {
	shredDamageMultiplier := 3.0

	flatDamageBonus := map[int32]float64{
		25: 54.0,
		40: 99.0,
		50: 144.0,
		60: 180.0,
	}[druid.Level] / shredDamageMultiplier

	hasGoreRune := druid.HasRune(proto.DruidRune_RuneHelmGore)
	hasElunesFires := druid.HasRune(proto.DruidRune_RuneBracersElunesFires)

	druid.Shred = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID: core.ActionID{SpellID: map[int32]int32{
			25: 5221,
			40: 8992,
			50: 9829,
			60: 9830,
		}[druid.Level]},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

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

		DamageMultiplier: shredDamageMultiplier,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatDamageBonus +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			modifier := 1.0
			if druid.BleedCategories.Get(target).AnyActive() {
				modifier += .3
			}

			/*
				ripDot := druid.Rip.Dot(target)
				if druid.AssumeBleedActive || ripDot.IsActive() || druid.Rake.Dot(target).IsActive() || druid.Lacerate.Dot(target).IsActive() {
					modifier *= 1.0 + (0.04 * float64(druid.Talents.RendAndTear))
				}
			*/

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
				modifier += .3
			}

			/*
				if druid.AssumeBleedActive || druid.Rip.Dot(target).IsActive() || druid.Rake.Dot(target).IsActive() || druid.Lacerate.Dot(target).IsActive() {
					modifier *= 1.0 + (0.04 * float64(druid.Talents.RendAndTear))
				}
			*/

			spell.DamageMultiplier *= modifier
			baseres := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			spell.DamageMultiplier /= modifier

			attackTable := spell.Unit.AttackTables[target.UnitIndex][spell.CastType]
			critChance := spell.PhysicalCritChance(attackTable)
			critMod := (critChance * (spell.CritMultiplier(attackTable) - 1))

			baseres.Damage *= (1 + critMod)

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
