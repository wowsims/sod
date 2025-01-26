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

func (paladin *Paladin) registerDivineStorm() {
	if !paladin.hasRune(proto.PaladinRune_RuneChestDivineStorm) {
		return
	}

	results := make([]*core.SpellResult, min(4, paladin.Env.GetNumTargets()))

	healthMetrics := paladin.NewHealthMetrics(core.ActionID{SpellID: int32(proto.PaladinRune_RuneChestDivineStorm)})

	divineStormSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       healthMetrics.ActionID,
		ClassSpellMask: ClassSpellMask_PaladinDivineStorm,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlag_RV,

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

			for idx := range results {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + core.SpellBatchWindow,
					OnAction: func(s *core.Simulation) {
						spell.DealDamage(sim, result)
						paladin.GainHealth(sim, result.Damage*0.25, healthMetrics)
					},
				})
			}

		},
	})

	paladin.divineStorm = divineStormSpell
}
