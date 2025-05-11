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
		ClassSpellMask: ClassSpellMask_WarlockChaosBolt,
		ActionID:       core.ActionID{SpellID: 403629},
		SpellSchool:    core.SpellSchoolArcane | core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolNature | core.SpellSchoolShadow,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagDestruction,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.07,
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

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
		},
	})
}
