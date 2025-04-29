package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) registerChimeraShotSpell() {
	if !hunter.HasRune(proto.HunterRune_RuneHandsChimeraShot) {
		return
	}

	hunter.SerpentStingChimeraShot = hunter.chimeraShotSerpentStingSpell(hunter.SerpentSting.Rank)

	hunter.ChimeraShot = hunter.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterChimeraShot,
		ActionID:       core.ActionID{SpellID: 409433},
		SpellSchool:    core.SpellSchoolNature,
		CastType:       proto.CastType_CastTypeRanged,
		DefenseType:    core.DefenseTypeRanged,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists | core.SpellFlagAPL,

		MinRange:     core.MinRangedAttackRange,
		MaxRange:     core.MaxRangedAttackRange,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.06,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1.35,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target, false)) +
				hunter.AmmoDamageBonus

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)

				if result.Landed() {
					if hunter.SerpentSting.Dot(target).IsActive() {
						hunter.SerpentStingChimeraShot.Cast(sim, target)
						hunter.SerpentSting.Dot(target).Rollover(sim)
					}
				}
			})
		},
	})
}
