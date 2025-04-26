package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ShadowflameCastTime = time.Second * 2

func (warlock *Warlock) registerShadowflameSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneBootsShadowflame) {
		return
	}

	warlock.Shadowflame = warlock.RegisterSpell(warlock.getShadowflameConfig())
}

func (warlock *Warlock) getShadowflameConfig() core.SpellConfig {
	hasHauntRune := warlock.HasRune(proto.WarlockRune_RuneHandsHaunt)
	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)

	numTicks := int32(5)

	baseSpellCoeff := 0.20
	dotSpellCoeff := 0.13
	baseDamage := warlock.baseRuneAbilityDamage() * 2.26
	dotDamage := warlock.baseRuneAbilityDamage() * 3.2 / float64(numTicks)

	tickLength := time.Second * 3

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_WarlockShadowflame,
		ActionID:       core.ActionID{SpellID: 426320},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolShadow,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagAffliction | WarlockFlagDestruction | WarlockFlagHaunt,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.27,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Shadowflame" + warlock.Label,
			},

			NumberOfTicks:    numTicks,
			TickLength:       tickLength,
			BonusCoefficient: dotSpellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if hasPandemicRune {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				}
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: baseSpellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damageMultiplier := 1.0
			if hasHauntRune && warlock.HauntDebuffAuras.Get(target).IsActive() {
				// @Lucenia: Haunt incorrectly applies to the impact damage of the spell even in-game/
				// This was fixed in Phase 7
				damageMultiplier = hauntMultiplier(spell, warlock.AttackTables[target.UnitIndex][proto.CastType_CastTypeMainHand])
			}

			spell.ApplyMultiplicativeDamageBonus(1 / damageMultiplier)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.ApplyMultiplicativeDamageBonus(damageMultiplier)

			if result.Landed() {
				// Shadowflame and Immolate are exclusive
				immoDot := warlock.getActiveImmolateSpell(target)
				if immoDot != nil {
					immoDot.Dot(target).Deactivate(sim)
				}

				spell.Dot(target).ApplyOrReset(sim)
			}
		},
	}
}
