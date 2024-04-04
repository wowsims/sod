package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) registerFlankingStrikeSpell() {
	if !hunter.HasRune(proto.HunterRune_RuneLegsFlankingStrike) {
		return
	}

	hasCatlikeReflexes := hunter.HasRune(proto.HunterRune_RuneHelmCatlikeReflexes)
	cooldownModifier := 1.0
	if hasCatlikeReflexes {
		cooldownModifier *= 0.5
	}

	hunter.FlankingStrikeAura = hunter.GetOrRegisterAura(core.Aura{
		Label:     "Flanking Strike Buff",
		ActionID:  core.ActionID{SpellID: 415320},
		MaxStacks: 3,
		Duration:  time.Second * 10,
	})

	if hunter.pet != nil {
		hunter.pet.FlankingStrike = hunter.pet.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 415320},
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

			DamageMultiplier: 0.45,
			BonusCoefficient: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			},
		})
	}

	hunter.FlankingStrike = hunter.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 415320},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.015,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * time.Duration(30*cooldownModifier),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= 5
		},

		DamageMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if hunter.pet != nil {
				hunter.pet.FlankingStrike.Cast(sim, hunter.pet.CurrentTarget)
			}

			hunter.FlankingStrikeAura.Activate(sim)
			hunter.FlankingStrikeAura.AddStack(sim)
		},
	})
}
