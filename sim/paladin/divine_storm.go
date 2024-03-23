package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// Divine Storm is a non-ap normalised instant attack that has a weapon damage % modifier with a 1.1 coefficient.
// It does this damage to up to 4 targets in range.
// DS also heals up to 3 party or raid members for 25% of the total damage caused. This has implications for prot
// paladin threat, so we'll implement this as a heal to the casting paladin for now.

func (paladin *Paladin) registerDivineStormSpell() {
	if !paladin.HasRune(proto.PaladinRune_RuneChestDivineStorm) {
		return
	}
	numHits := min(4, paladin.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	actionID := core.ActionID{SpellID: 407778}
	healthMetrics := paladin.NewHealthMetrics(actionID)

	paladin.DivineStorm = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.12,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 1.1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			totalDamageDealt := 0.0
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()

				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				totalDamageDealt += results[hitIndex].Damage
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
			paladin.GainHealth(sim, totalDamageDealt*0.25, healthMetrics)
		},
	})
}
