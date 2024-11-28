package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// See https://www.wowhead.com/classic/spell=436895/s03-tuning-and-overrides-passive-druid
// Modifies Threat +101%:
const SwipeThreatMultiplier = 3.5

func (druid *Druid) registerSwipeBearSpell() {
	hasImprovedSwipeRune := druid.HasRune(proto.DruidRune_RuneCloakImprovedSwipe)
	baseMultiplier := 1.0

	baseDamage := 83 + .1*druid.GetStat(stats.AttackPower)

	rageCost := 20 - float64(druid.Talents.Ferocity)
	targetCount := core.TernaryInt32(hasImprovedSwipeRune, 10, 3)
	numHits := min(targetCount, druid.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	switch druid.Ranged().ID {
	case IdolOfBrutality:
		rageCost -= 3
	case IdolOfUrsinPower:
		baseMultiplier += .03
	}

	druid.SwipeBear = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 9908},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 20 - float64(druid.Talents.Ferocity),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: baseMultiplier + 0.1*float64(druid.Talents.SavageFury),
		ThreatMultiplier: SwipeThreatMultiplier,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}

func (druid *Druid) registerSwipeCatSpell() {
	if !druid.HasRune(proto.DruidRune_RuneCloakImprovedSwipe) {
		return
	}

	weaponMulti := 2.5

	druid.SwipeCat = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 411128},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagOmen | SpellFlagBuilder,

		EnergyCost: core.EnergyCostOptions{
			Cost: 50 - float64(druid.Talents.Ferocity),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: (1 + 0.1*float64(druid.Talents.SavageFury)) * weaponMulti,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			aoeTarget := target
			for i := 0; i < len(sim.Encounter.TargetUnits); i++ {
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				if i == 0 && result.Landed() {
					druid.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())
				}
				aoeTarget = sim.Environment.NextTargetUnit(aoeTarget)
			}
		},
	})
}

func (druid *Druid) CurrentSwipeCatCost() float64 {
	return druid.SwipeCat.Cost.GetCurrentCost()
}
