package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerEnvenom() {
	if !rogue.HasRune(proto.RogueRune_RuneEnvenom) {
		return
	}

	baseAbilityDamage := rogue.baseRuneAbilityDamage()

	cutToTheChase := rogue.HasRune(proto.RogueRune_RuneCutToTheChase)

	rogue.EnvenomAura = rogue.RegisterAura(core.Aura{
		Label:    "Envenom",
		ActionID: core.ActionID{SpellID: int32(proto.RogueRune_RuneEnvenom)},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.instantPoisonProcChanceBonus += 0.75
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.instantPoisonProcChanceBonus -= 0.75
		},
	})

	rogue.Envenom = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: int32(proto.RogueRune_RuneEnvenom)},
		SpellSchool:  core.SpellSchoolNature,
		DefenseType:  core.DefenseTypeMelee,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        rogue.finisherFlags() | SpellFlagColdBlooded | core.SpellFlagIgnoreResists | core.SpellFlagPoison,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:   35,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0 && rogue.deadlyPoisonTick.Dot(target).IsActive()
		},

		DamageMultiplier: rogue.getPoisonDamageMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			comboPoints := rogue.ComboPoints()
			// - the aura is active even if the attack fails to land
			// - the aura is applied before the hit effect
			// See: https://github.com/where-fore/rogue-wotlk/issues/32
			// Still true in SoD
			rogue.EnvenomAura.Duration = rogue.EnvenomDuration(rogue.ComboPoints())
			rogue.EnvenomAura.Activate(sim)

			dp := rogue.deadlyPoisonTick.Dot(target)
			// - base damage is scaled by consumed doses (<= comboPoints)
			// - apRatio is independent of consumed doses (== comboPoints)
			// - Spell power is 1:1 at all ranks and cp
			consumed := min(dp.GetStacks(), comboPoints)
			baseDamage := baseAbilityDamage*float64(consumed)*0.8 + 0.072*float64(comboPoints)*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.SpendComboPoints(sim, spell)
				if cutToTheChase {
					rogue.ApplyCutToTheChase(sim)
				}
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (rogue *Rogue) EnvenomDuration(comboPoints int32) time.Duration {
	return time.Second * (1 + time.Duration(comboPoints))
}
