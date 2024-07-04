package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (mage *Mage) registerFrostfireBoltSpell() {
	if !mage.HasRune(proto.MageRune_RuneBeltFrostfireBolt) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.MageRune_RuneBeltFrostfireBolt)}
	// TODO: Tuning 2024-07-03 FFB base damage reduced by 16%
	// Need to update values once wowhead data is updated
	baseDamageLow := mage.baseRuneAbilityDamage() * 3.87 * .84
	baseDamageHigh := mage.baseRuneAbilityDamage() * 4.51 * .84
	baseDotDamage := mage.baseRuneAbilityDamage() * .08
	spellCoeff := 1.0
	castTime := time.Second * 3
	manaCost := .14

	numTicks := int32(3)
	tickLength := time.Second * 3

	mage.FrostfireBolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellCode:    SpellCode_MageFrostfireBolt,
		SpellSchool:  core.SpellSchoolFrost | core.SpellSchoolFire,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | SpellFlagChillSpell | core.SpellFlagBinary | core.SpellFlagAPL,
		MissileSpeed: 28,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: castTime,
				GCD:      core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Frostfire Bolt",
				ActionID: actionID.WithTag(1),
			},
			NumberOfTicks: numTicks,
			TickLength:    tickLength,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)

				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			})
		},
	})
}
