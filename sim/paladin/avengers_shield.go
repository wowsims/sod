package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// The Arcane Shot of a Hunter, and the Hammer of Wrath and Avenger's Shield talent of a
// Protection-specced Paladin, are resolved as ranged attacks that do non-physical damage.
// They can miss—rather than be "fully resisted"—and they do double damage on a crit.
// The only difference is that if a mob target is higher level than the player attacker,
// or if the target has any resistance to the school of magic used by the attack, the same
// check is made to see if the damage is partially resisted as would happen from a spell.
// https://wowwiki-archive.fandom.com/wiki/Attack_table#Magic-damage_ranged_special_attacks
func (paladin *Paladin) registerAvengersShield() {
	if !paladin.hasRune(proto.PaladinRune_RuneLegsAvengersShield) {
		return
	}

	hasLibramOfAvenging := func() bool {
		return paladin.Ranged().ID == LibramOfAvenging
	}

	hasDefenseSpecRune := paladin.HasRuneById(int32(proto.RingRune_RuneRingDefenseSpecialization))

	// Avenger's Shield hits up to 3 targets. It cannot miss or be resisted.
	numTargets := min(3, int(paladin.Env.GetNumTargets()))

	lowDamage := 366 * paladin.baseRuneAbilityDamage() / 100
	highDamage := 448 * paladin.baseRuneAbilityDamage() / 100

	paladin.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: int32(proto.PaladinRune_RuneLegsAvengersShield)},
		ClassSpellMask: ClassSpellMask_PaladinAvengersShield,
		SpellSchool:    core.SpellSchoolHoly,
		DefenseType:    core.DefenseTypeRanged, // Crits as if melee for 200%
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagBinary | core.SpellFlagBatchStartAttackMacro,
		CastType:       proto.CastType_CastTypeRanged,
		MissileSpeed:   35, // Verified from game files using WoW tools.
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
		BonusHitRating:   core.TernaryFloat64(hasDefenseSpecRune, 3.0*core.MeleeHitRatingPerHitChance, 0),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			apBonus := 0.091 * spell.MeleeAttackPower()
			baseTravelTime := spell.TravelTime()
			if hasLibramOfAvenging() {
				// Libram of Avenging causes Avenger's Shield to be single target, but it
				// hits the target twice. The second projectile fires after a fixed 1.5s delay.
				firstHit := spell.CalcDamage(sim, target, sim.Roll(lowDamage, highDamage)+apBonus, spell.OutcomeRangedHitAndCrit)
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, firstHit)
				})

				secondHit := spell.CalcDamage(sim, target, sim.Roll(lowDamage, highDamage)+apBonus, spell.OutcomeRangedHitAndCrit)
				timeToSecondHit := baseTravelTime + time.Millisecond*1500
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + timeToSecondHit,
					OnAction: func(s *core.Simulation) {
						spell.DealDamage(sim, secondHit)
					},
				})

			} else {
				interTargetTravelTime := int(float64(time.Second) * 3.0 / spell.MissileSpeed)
				for i := 0; i < numTargets; i++ {
					// Avenger's Shield bounces from target 1 > target 2 > target 3 at MissileSpeed.
					// We approximate it by assuming targets are standing ~3 yds apart from each other.
					// The damage for each target is therefore scheduled to arrive at:
					// T1 = (TravelTime from player; by default 5 yard max melee range)
					// T2 = T1 + (3 yd TravelTime)
					// T3 = T2 + (3 yd TravelTime)
					baseDamage := sim.Roll(lowDamage, highDamage) + apBonus
					delay := time.Duration(interTargetTravelTime * i)
					nextTarget := target // create new ref for delayed action evaluation
					result := spell.CalcDamage(sim, nextTarget, baseDamage, spell.OutcomeRangedHitAndCrit)

					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt: sim.CurrentTime + baseTravelTime + delay,
						OnAction: func(s *core.Simulation) {
							spell.DealDamage(sim, result)

						},
					})
					target = sim.NextTargetUnit(target)
				}
			}
		},
	})
}
