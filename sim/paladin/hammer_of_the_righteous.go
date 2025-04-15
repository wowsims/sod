package paladin

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (paladin *Paladin) registerHammerOfTheRighteous() {
	if !paladin.hasRune(proto.PaladinRune_RuneWristHammerOfTheRighteous) {
		return
	}

	// Phase 4: Hammer of the Righteous damage reduced by 50% but threat increased by 2X.
	// https://www.wowhead.com/classic/news/development-notes-for-phase-4-ptr-season-of-discovery-new-runes-class-changes-3428960
	results := make([]*core.SpellResult, min(3, paladin.Env.GetNumTargets()))

	paladin.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: int32(proto.PaladinRune_RuneWristHammerOfTheRighteous)},
		SpellSchool:    core.SpellSchoolHoly,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIgnoreResists | core.SpellFlagBatchStartAttackMacro,
		ClassSpellMask: ClassSpellMask_PaladinHammerOfTheRighteous,
		MissileSpeed:   35,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.06,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: false,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return slices.Contains([]proto.HandType{proto.HandType_HandTypeMainHand, proto.HandType_HandTypeOneHand}, paladin.MainHand().HandType) &&
				paladin.MainHand().WeaponType != proto.WeaponType_WeaponTypeUnknown
		},
		DamageMultiplier: 1,
		ThreatMultiplier: 2, // verified with TinyThreat in game

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weapon := paladin.AutoAttacks.MH()
			baseDamage := 3.0 * (weapon.CalculateAverageWeaponDamage(spell.MeleeAttackPower()) / weapon.SwingSpeed)

			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				for _, result := range results {
					spell.DealDamage(sim, result)
				}
			})
		},
	})
}
