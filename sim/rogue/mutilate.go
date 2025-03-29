package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) newMutilateHitSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: int32(proto.RogueRune_RuneMutilate)}
	castType := proto.CastType_CastTypeMainHand
	procMask := core.ProcMaskMeleeMHSpecial
	if !isMH {
		castType = proto.CastType_CastTypeOffHand
		procMask = core.ProcMaskMeleeOHSpecial
	}

	// waylay := rogue.HasRune(proto.RogueRune_RuneWaylay)

	flatDamageBonus := rogue.baseRuneAbilityDamage()

	return rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueMutilateHit,
		ActionID:       actionID.WithTag(int32(core.Ternary(isMH, 1, 2))),
		SpellSchool:    core.SpellSchoolPhysical,
		CastType:       castType,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       procMask,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | SpellFlagBuilder | SpellFlagColdBlooded | SpellFlagCarnage,

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: 1 * core.TernaryFloat64(isMH, 1, rogue.dwsMultiplier()),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = flatDamageBonus*0.8 + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())*0.8
			} else {
				baseDamage = flatDamageBonus*0.8 + spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())*0.8
			}

			// TODO: Add support for all poison effects (such as chipped bite proc), if they apply ;)
			oldMultiplier := spell.GetDamageMultiplier()
			if rogue.deadlyPoisonTick.Dot(target).IsActive() || (rogue.occultPoisonTick != nil && rogue.occultPoisonTick.Dot(target).IsActive()) || rogue.woundPoisonDebuffAuras.Get(target).IsActive() {
				spell.ApplyMultiplicativeDamageBonus(1.2)
			}
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
			spell.SetMultiplicativeDamageBonus(oldMultiplier)
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
		ActionID:       core.ActionID{SpellID: int32(proto.RogueRune_RuneMutilate)},
		ClassSpellMask: ClassSpellMask_RogueMutilate,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       5,

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
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				rogue.AddComboPoints(sim, 2, target, spell.ComboPointMetrics())

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
