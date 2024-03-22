package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) registerShadowflameSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneBootsShadowflame) {
		return
	}

	level := float64(warlock.GetCharacter().Level)
	baseSpellCoeff := 0.0715
	dotSpellCoeff := 0.022

	baseCalc := (6.568597 + 0.672028*level + 0.031721*level*level)
	baseDamage := baseCalc * 0.64
	dotDamage := baseCalc * 0.24

	shadowMasteryMulti := 1 + 0.02*float64(warlock.Talents.ShadowMastery)
	emberstormMulti := 1 + 0.02*float64(warlock.Talents.Emberstorm)

	fireDot := warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 426325},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Shadowflame" + warlock.Label,
			},

			NumberOfTicks: 4,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = dotDamage + dotSpellCoeff*dot.Spell.SpellDamage()
				dot.SnapshotBaseDamage *= emberstormMulti

				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
				if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
					result.Damage *= 1.4
					result.Threat *= 1.4
				}
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			spell.Dot(target).Apply(sim)
		},
	})

	warlock.Shadowflame = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 426320},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.27,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: 15 * time.Second,
			},
		},

		BonusCritRating: float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,

		CritDamageBonus: warlock.ruin(),

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage = baseDamage + baseSpellCoeff*spell.SpellDamage()
			baseDamage *= shadowMasteryMulti

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					fireDot.Cast(sim, aoeTarget)
				}
			}
		},
	})
}
