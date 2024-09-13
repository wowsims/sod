package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"time"
)

func (paladin *Paladin) registerHammerOfTheRighteous() {
	if !paladin.hasRune(proto.PaladinRune_RuneWristHammerOfTheRighteous) {
		return
	}

	// Phase 4: Hammer of the Righteous damage reduced by 50% but threat increased by 2X.
	// https://www.wowhead.com/classic/news/development-notes-for-phase-4-ptr-season-of-discovery-new-runes-class-changes-3428960
	results := make([]*core.SpellResult, min(3, paladin.Env.GetNumTargets()))

	paladin.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.PaladinRune_RuneWristHammerOfTheRighteous)},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIgnoreResists,

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
			return paladin.MainHand().HandType == proto.HandType_HandTypeOneHand
		},
		DamageMultiplier: 3,
		ThreatMultiplier: 2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weapon := paladin.AutoAttacks.MH()
			baseDamage := weapon.CalculateAverageWeaponDamage(spell.MeleeAttackPower()) / weapon.SwingSpeed

			for idx := range results {
				// Hammer of the Righteous does not miss, but can crit and be blocked.
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedCritOnly)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}
