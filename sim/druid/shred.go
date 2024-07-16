package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// See https://www.wowhead.com/classic/spell=436895/s03-tuning-and-overrides-passive-druid
// Modifies Effect #1's Value -24%:
// Modifies Effect #2's Value +76:
const ShredBonusDmgMultiplier = 1.75 // multiplier +75%
const ShredFlatDmgMultiplier = .075  // multiplier -25%

func (druid *Druid) registerShredSpell() {
	damageMultiplier := 2.25 + ShredBonusDmgMultiplier

	flatDamageBonus := map[int32]float64{
		25: 24,
		40: 44,
		50: 64,
		60: 80,
	}[druid.Level] * ShredFlatDmgMultiplier

	hasGoreRune := druid.HasRune(proto.DruidRune_RuneHelmGore)

	if druid.Ranged().ID == IdolOfTheDream {
		damageMultiplier *= 1.02
		flatDamageBonus *= 1.02
	}

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
			baseDamage := flatDamageBonus + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			oldMultiplier := spell.DamageMultiplier
			if druid.BleedCategories.Get(target).AnyActive() {
				spell.DamageMultiplier *= 1.3
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.DamageMultiplier = oldMultiplier

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())

				if hasGoreRune {
					druid.rollGoreCatReset(sim)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := flatDamageBonus + spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower())

			oldMultiplier := spell.DamageMultiplier
			if druid.BleedCategories.Get(target).AnyActive() {
				spell.DamageMultiplier *= 1.3
			}

			baseres := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			spell.DamageMultiplier = oldMultiplier

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
