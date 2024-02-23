package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) registerSteadyShotSpell() {
	if !hunter.HasRune(proto.HunterRune_RuneBeltSteadyShot) {
		return
	}

	hasCobraStrikes := hunter.pet != nil && hunter.HasRune(proto.HunterRune_RuneChestCobraStrikes)

	manaCostMultiplier := 1 - 0.02*float64(hunter.Talents.Efficiency)
	if hunter.HasRune(proto.HunterRune_RuneChestMasterMarksman) {
		manaCostMultiplier -= 0.25
	}

	hunter.GetOrRegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 437123},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		CastType:     proto.CastType_CastTypeRanged,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.05,
			Multiplier: manaCostMultiplier,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.RangedSwingSpeed())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget >= 8
		},

		BonusCritRating:          0,
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         0.6,
		CritMultiplier:           hunter.critMultiplier(true, hunter.CurrentTarget),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := hunter.AutoAttacks.Ranged().CalculateWeaponDamage(sim, spell.RangedAttackPower(target)) +
				hunter.AmmoDamageBonus +
				spell.BonusWeaponDamage()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)

				if hasCobraStrikes && result.DidCrit() {
					hunter.CobraStrikesAura.Activate(sim)
					hunter.CobraStrikesAura.SetStacks(sim, 2)
				}
			})
		},
	})
}
