package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const (
	LivingBombBaseNumTicks   = 4
	LivingBombBaseTickLength = time.Second * 3
)

func (mage *Mage) registerLivingBombSpell() {
	if !mage.HasRune(proto.MageRune_RuneHandsLivingBomb) {
		return
	}

	hasImprovedScorchTalent := mage.Talents.ImprovedScorch > 0

	actionID := core.ActionID{SpellID: int32(proto.MageRune_RuneHandsLivingBomb)}
	baseDotDamage := mage.baseRuneAbilityDamage() * .85
	baseExplosionDamage := mage.baseRuneAbilityDamage() * 1.71
	dotCoeff := .20
	explosionCoeff := .40
	manaCost := .22

	livingBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		ClassSpellMask: ClassSpellMask_MageLivingBombExplosion,
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: explosionCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, baseExplosionDamage, spell.OutcomeMagicCrit)

			// Unlike the normal Scorch application, Living Bomb explosions are guaranteed to apply Improved Scorch as long as they have at least 1 point talented.
			if hasImprovedScorchTalent {
				impScorchAura := mage.ImprovedScorchAuras.Get(target)
				impScorchAura.Activate(sim)
				impScorchAura.AddStack(sim)
			}
		},
	})

	mage.LivingBomb = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_MageLivingBomb,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagPureDot,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Living Bomb (DoT)",
			},

			NumberOfTicks:    LivingBombBaseNumTicks,
			TickLength:       LivingBombBaseTickLength,
			BonusCoefficient: dotCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)

				// Unlike the normal Scorch application, Living Bomb ticks are guaranteed to apply Improved Scorch as long as they have at least 1 point talented.
				if hasImprovedScorchTalent {
					impScorchAura := mage.ImprovedScorchAuras.Get(target)
					impScorchAura.Activate(sim)
					impScorchAura.AddStack(sim)
				}

				if !dot.IsActive() {
					for _, aoeTarget := range sim.Encounter.TargetUnits {
						livingBombExplosionSpell.Cast(sim, aoeTarget)
					}
				}
			},
		},

		BonusCritRating: 2 * float64(mage.Talents.Incinerate),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).ApplyOrReset(sim)
			}
		},
	})
}
