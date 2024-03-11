package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var MutilateSpellID int32 = 399956

func (rogue *Rogue) newMutilateHitSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: 399960}
	procMask := core.ProcMaskMeleeMHSpecial
	if !isMH {
		actionID = core.ActionID{SpellID: 399961}
		procMask = core.ProcMaskMeleeOHSpecial
	}

	// waylay := rogue.HasRune(proto.RogueRune_RuneWaylay)

	flatDamageBonus := rogue.RuneAbilityBaseDamage()

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID.WithTag(int32(core.Ternary(isMH, 1, 2))),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder | SpellFlagColdBlooded,

		BonusCritRating: 10 * core.CritRatingPerCritChance * float64(rogue.Talents.ImprovedBackstab),

		DamageMultiplier: 1 *
			core.TernaryFloat64(isMH, 1, rogue.dwsMultiplier()) *
			[]float64{1, 1.04, 1.08, 1.12, 1.16, 1.2}[rogue.Talents.Opportunity],
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO check whether Mutilate profits from "+Damage Done" or "+Damage Taken" auras (not the case in TBC)
			var baseDamage float64
			if isMH {
				baseDamage = flatDamageBonus + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				baseDamage = flatDamageBonus + spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			}
			// TODO: Add support for all poison effects (such as chipped bite proc)
			if target.GetAuraByID(rogue.DeadlyPoison[0].ActionID).IsActive() || rogue.woundPoisonDebuffAuras.Get(target).IsActive() {
				baseDamage *= 1.2
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})
}

func (rogue *Rogue) registerMutilateSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneMutilate) {
		return
	}

	// Requires Daggers (2 of them)
	if !rogue.HasDagger(core.MainHand) || !rogue.HasDagger(core.OffHand) {
		return
	}

	rogue.MutilateMH = rogue.newMutilateHitSpell(true)
	rogue.MutilateOH = rogue.newMutilateHitSpell(false)

	rogue.Mutilate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: MutilateSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   40,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit) // Miss/Dodge/Parry/Hit
			if result.Landed() {
				rogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())

				/** Disable until it works on bosses
				if waylay {
					rogue.WaylayAuras.Get(target).Activate(sim)
				} */
				rogue.MutilateMH.Cast(sim, target)
				rogue.MutilateOH.Cast(sim, target)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
