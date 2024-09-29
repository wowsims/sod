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

	ssProcSpell := make([]*core.Spell, 10)
	for i := 1; i <= 9; i++ {
		ssProcSpell[i] = hunter.chimeraShotSerpentStingSpell(i)
	}

	hunter.ChimeraShot = hunter.RegisterSpell(core.SpellConfig{
		SpellCode:    SpellCode_HunterChimeraShot,
		ActionID:     core.ActionID{SpellID: 409433},
		SpellSchool:  core.SpellSchoolNature,
		DefenseType:  core.DefenseTypeRanged,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists | core.SpellFlagAPL | SpellFlagShot,
		CastType:     proto.CastType_CastTypeRanged,
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
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget >= core.MinRangedAttackDistance
		},

		CritDamageBonus: hunter.mortalShots(),

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
						ssProcSpell[hunter.SerpentSting.Rank].Cast(sim, target)
						hunter.SerpentSting.Dot(target).Rollover(sim)
					}
				}
			})
		},
	})
}
