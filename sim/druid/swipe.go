package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

type SwipeRankInfo struct {
	id     int32
	level  int32
	damage float64
}

var swipeSpells = []SwipeRankInfo{
	{
		id:     779,
		level:  16,
		damage: 18.0,
	},
	{
		id:     780,
		level:  24,
		damage: 25.0,
	},
	{
		id:     769,
		level:  34,
		damage: 36.0,
	},
	{
		id:     9754,
		level:  44,
		damage: 60.0,
	},
	{
		id:     9908,
		level:  54,
		damage: 83.0,
	},
}

func (druid *Druid) registerSwipeBearSpell() {
	// Add highest available rank for level.
	for rank := len(swipeSpells) - 1; rank >= 0; rank-- {
		if druid.Level >= swipeSpells[rank].level {
			config := druid.newSwipeBearSpellConfig(swipeSpells[rank])
			druid.SwipeBear = druid.RegisterSpell(Bear, config)
			break
		}
	}
}

// See https://www.wowhead.com/classic/spell=436895/s03-tuning-and-overrides-passive-druid
// Modifies Threat +101%:
const SwipeThreatMultiplier = 3.5

func (druid *Druid) newSwipeBearSpellConfig(swipeRank SwipeRankInfo) core.SpellConfig {
	hasImprovedSwipeRune := druid.HasRune(proto.DruidRune_RuneCloakImprovedSwipe)

	baseDamage := swipeRank.damage

	targetCount := core.TernaryInt32(hasImprovedSwipeRune, 10, 3)
	numHits := min(targetCount, druid.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_DruidSwipeBear,
		ActionID:       core.ActionID{SpellID: swipeRank.id},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 20,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: SwipeThreatMultiplier,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := baseDamage + .1*spell.MeleeAttackPower()
			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, damage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)

			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	}
}

func (druid *Druid) registerSwipeCatSpell() {
	if !druid.HasRune(proto.DruidRune_RuneCloakImprovedSwipe) {
		return
	}

	druid.SwipeCat = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 411128},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: ClassSpellMask_DruidSwipeCat,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagOmen | SpellFlagBuilder,

		EnergyCost: core.EnergyCostOptions{
			Cost: 50,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 2.5,
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
