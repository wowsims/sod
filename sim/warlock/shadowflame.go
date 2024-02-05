package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) registerShadowflameSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneBootsShadowflame) {
		return
	}

	level := float64(warlock.GetCharacter().Level)
	baseSpellCoeff := 0.715
	dotSpellCoeff := 0.022

	baseCalc := (6.568597 + 0.672028*level + 0.031721*level*level)
	baseDamage := baseCalc * 0.64
	dotDamage := baseCalc * 0.24

	shadowMasteryMulti := 1 + 0.02*float64(warlock.Talents.ShadowMastery)
	emberstormMulti := 1 + 0.02*float64(warlock.Talents.Emberstorm)

	numHits := warlock.Env.GetNumTargets()
	results := make([]*core.SpellResult, numHits)

	warlock.Shadowflame = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 426320},
		SpellSchool: core.SpellSchoolShadow,
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

		BonusCritRating:          float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Ruin, 1, 0)),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Shadowflame",
			},

			NumberOfTicks: 4,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// Use fire school for dot modifiers
				dot.Spell.SpellSchool = core.SpellSchoolFire
				dot.Spell.SchoolIndex = stats.SchoolIndexFire

				dot.SnapshotBaseDamage = dotDamage + dotSpellCoeff*dot.Spell.SpellPower()
				dot.SnapshotBaseDamage *= emberstormMulti

				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])

				// Revert to shadow school
				dot.Spell.SpellSchool = core.SpellSchoolShadow
				dot.Spell.SchoolIndex = stats.SchoolIndexShadow
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Use fire school for dot modifiers
				dot.Spell.SpellSchool = core.SpellSchoolFire
				dot.Spell.SchoolIndex = stats.SchoolIndexFire

				hasLof := false
				if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
					hasLof = true
					target.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 1.4
				}

				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)

				if hasLof {
					target.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= 1.4
				}

				// Revert to shadow school
				dot.Spell.SpellSchool = core.SpellSchoolShadow
				dot.Spell.SchoolIndex = stats.SchoolIndexShadow
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage = baseDamage + baseSpellCoeff*spell.SpellPower()
			baseDamage *= shadowMasteryMulti

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])

				if results[hitIndex].Landed() {
					spell.Dot(curTarget).Apply(sim)
				}

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}
