package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: Classic warlock verify chaos bolt mechanics
func (warlock *Warlock) registerChaosBoltSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsChaosBolt) {
		return
	}

	spellCoeff := 0.714
	baseLowDamage := warlock.baseRuneAbilityDamage() * 5.22
	baseHighDamage := warlock.baseRuneAbilityDamage() * 6.62

	warlock.ChaosBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 403629},
		SpellSchool: core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.07,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		BonusCritRating: float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,

		CritDamageBonus: warlock.ruin(),

		DamageMultiplierAdditive: 1 + 0.02*float64(warlock.Talents.Emberstorm),
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + spellCoeff*spell.SpellDamage()

			if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
				baseDamage *= 1.4
			}

			// Assuming 100% hit for all target levels, numbers could be updated for level comparison later
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
		},
	})
}
