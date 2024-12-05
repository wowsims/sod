package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (druid *Druid) registerMangleBearSpell() {
	if !druid.HasRune(proto.DruidRune_RuneHandsMangle) {
		return
	}

	idolMultiplier := 1.0
	rageCostReduction := float64(druid.Talents.Ferocity)

	switch druid.Ranged().ID {
	case IdolOfUrsinPower:
		idolMultiplier += .03
	case IdolOfBrutality:
		rageCostReduction += 3
	}

	mangleAuras := druid.NewEnemyAuraArray(core.MangleAura)
	apProcAura := core.DefendersResolveAttackPower(druid.GetCharacter())

	druid.MangleBear = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 407995},
		SpellCode:   SpellCode_DruidMangleBear,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   15 - rageCostReduction,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		// TODO: Berserk 3 target mangle cleave - Saeyon

		DamageMultiplier: 1.6 * (1 + 0.1*float64(druid.Talents.SavageFury)) * idolMultiplier,
		ThreatMultiplier: 1.5,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			dotBonusCrit := 0.0
			if druid.LacerateBleed.Dot(target).GetStacks() > 0 {
				dotBonusCrit = druid.FuryOfStormrageCritRatingBonus
			}

			spell.BonusCritRating += dotBonusCrit
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.BonusCritRating -= dotBonusCrit

			if result.Landed() {
				mangleAuras.Get(target).Activate(sim)
				if stacks := int32(druid.GetStat(stats.Defense)); stacks > 0 {
					apProcAura.Activate(sim)
					if apProcAura.GetStacks() != stacks {
						apProcAura.SetStacks(sim, stacks)
					}
				}
			} else {
				spell.IssueRefund(sim)
			}

			if druid.BerserkAura.IsActive() {
				spell.CD.Reset()
			}
		},

		RelatedAuras: []core.AuraArray{mangleAuras},
	})
}

func (druid *Druid) registerMangleCatSpell() {
	if !druid.HasRune(proto.DruidRune_RuneHandsMangle) {
		return
	}

	hasGoreRune := druid.HasRune(proto.DruidRune_RuneHelmGore)

	weaponMulti := 2.7
	energyCost := 40 - float64(druid.Talents.Ferocity)

	mangleAuras := druid.NewEnemyAuraArray(core.MangleAura)
	druid.MangleCat = druid.RegisterSpell(Cat, core.SpellConfig{
		SpellCode:   SpellCode_DruidMangleCat,
		ActionID:    core.ActionID{SpellID: 409828},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagOmen | SpellFlagBuilder,

		EnergyCost: core.EnergyCostOptions{
			Cost:   energyCost,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplierAdditive: (1 + 0.1*float64(druid.Talents.SavageFury)) * weaponMulti,
		ThreatMultiplier:         1,
		BonusCoefficient:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())
				mangleAuras.Get(target).Activate(sim)

				if hasGoreRune {
					druid.rollGoreCatReset(sim)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},

		RelatedAuras: []core.AuraArray{mangleAuras},
	})
}

func (druid *Druid) CurrentMangleCatCost() float64 {
	return druid.MangleCat.Cost.GetCurrentCost()
}

func (druid *Druid) IsMangle(spell *core.Spell) bool {
	if druid.MangleBear != nil && druid.MangleBear.IsEqual(spell) {
		return true
	} else if druid.MangleCat != nil && druid.MangleCat.IsEqual(spell) {
		return true
	}
	return false
}
