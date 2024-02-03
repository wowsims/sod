package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) registerSlamSpell() {
	if warrior.Level < 30 {
		return
	}

	var castTime time.Duration
	var cooldown time.Duration

	if warrior.HasRune(proto.WarriorRune_RunePreciseTiming) {
		castTime = 0
		cooldown = 6 * time.Second
	} else {
		castTime = time.Millisecond*1500 - time.Millisecond*500*time.Duration(warrior.Talents.ImprovedSlam)
		cooldown = 0
	}

	flatDamageBonus := map[int32]float64{
		40: 43,
		50: 68,
		60: 87,
	}[warrior.Level]

	spellID := map[int32]int32{
		40: 8820,
		50: 11604,
		60: 11605,
	}[warrior.Level]

	warrior.Slam = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   15 - warrior.FocusedRageDiscount,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldown,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if cast.CastTime > 0 {
					warrior.AutoAttacks.DelayMeleeBy(sim, cast.CastTime)
				}
			},
		},

		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		FlatThreatBonus:  140,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatDamageBonus +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
