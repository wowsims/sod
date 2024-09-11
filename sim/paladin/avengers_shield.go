package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"time"
)

func (paladin *Paladin) registerAvengersShield() {
	if !paladin.hasRune(proto.PaladinRune_RuneLegsAvengersShield) {
		return
	}

	// Avenger's Shield hits up to 3 targets. It cannot miss or be resisted.
	results := make([]*core.SpellResult, min(3, paladin.Env.GetNumTargets()))
	lowDamage := 366 * paladin.baseRuneAbilityDamage() / 100
	highDamage := 448 * paladin.baseRuneAbilityDamage() / 100

	paladin.GetOrRegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: int32(proto.PaladinRune_RuneLegsAvengersShield)},
		SpellCode:    SpellCode_PaladinAvengersShield,
		SpellSchool:  core.SpellSchoolHoly,
		DefenseType:  core.DefenseTypeMelee, // Crits as if melee for 200%
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagIgnoreResists | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed: 35, // Verified from game files using WoW tools.
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.26,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 15,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield
		},
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.091, // for spell damage; we add the AP bonus manually

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			apBonus := 0.091 * spell.MeleeAttackPower()
			for idx := range results {
				baseDamage := sim.Roll(lowDamage, highDamage) + apBonus
				// Avenger's Shield cannot miss and uses magic critical _chance_.
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
				target = sim.Environment.NextTargetUnit(target)
			}
			// Avenger's Shield bounces from target 1 > target 2 > target 3 at MissileSpeed.
			// We approximate it by assuming targets are standing ~3 yds apart from each other.
			// The damage for each target is therefore scheduled to arrive at:
			// T1 = (TravelTime from player; by default 5 yard max melee range)
			// T2 = T1 + (3 yd TravelTime)
			// T3 = T2 + (3 yd TravelTime)
			baseTravelTime := spell.TravelTime()
			interTargetTravelTime := int(float64(time.Second) * 3.0 / spell.MissileSpeed)
			for i, result := range results {
				delay := time.Duration(interTargetTravelTime * i)
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + baseTravelTime + delay,
					OnAction: func(s *core.Simulation) {
						spell.DealDamage(sim, result)
					},
				})
			}
		},
	})
}
