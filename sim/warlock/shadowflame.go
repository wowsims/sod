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

	hasInvocationRune := warlock.HasRune(proto.WarlockRune_RuneBeltInvocation)
	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)

	baseSpellCoeff := 0.20
	dotSpellCoeff := 0.13
	baseDamage := warlock.baseRuneAbilityDamage() * 2.26
	dotDamage := warlock.baseRuneAbilityDamage() * 0.61

	numTicks := int32(5)
	tickLength := time.Second * 3

	warlock.Shadowflame = warlock.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_WarlockShadowflame,
		ActionID:    core.ActionID{SpellID: 426320},
		SpellSchool: core.SpellSchoolFire | core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagHaunt | WarlockFlagAffliction | WarlockFlagDestruction,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.27,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: ShadowflameCastTime,
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
				var result *core.SpellResult
				if hasPandemicRune {
					// We add the crit damage bonus and remove it after the call to not affect the initial damage portion of the spell
					dot.Spell.CritDamageBonus += 1
					result = dot.CalcSnapshotDamage(sim, target, dot.OutcomeTickSnapshotCrit)
					dot.Spell.CritDamageBonus -= 1
				} else {
					result = dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
				}
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: baseSpellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			oldMultiplier := spell.DamageMultiplier
			spell.DamageMultiplier *= 1 + warlock.improvedImmolateBonus()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DamageMultiplier = oldMultiplier

			if result.Landed() {
				// Shadowflame and Immolate are exclusive
				immoDot := warlock.getActiveImmolateSpell(target)
				if immoDot != nil {
					immoDot.Dot(target).Deactivate(sim)
				}

				if hasInvocationRune && spell.Dot(target).IsActive() {
					warlock.InvocationRefresh(sim, spell.Dot(target))
				}

				spell.Dot(target).Apply(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
