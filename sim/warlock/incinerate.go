package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) registerIncinerateSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneBracerIncinerate) {
		return
	}
	spellCoeff := 0.714
	baseLowDamage := warlock.baseRuneAbilityDamage() * 2.22
	baseHighDamage := warlock.baseRuneAbilityDamage() * 2.58

	warlock.IncinerateAura = warlock.RegisterAura(core.Aura{
		Label:    "Incinerate Aura",
		ActionID: core.ActionID{SpellID: int32(proto.WarlockRune_RuneBracerIncinerate)},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.25
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= 1.25
		},
	})

	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		SpellCode:    SpellCode_WarlockIncinerate,
		ActionID:     core.ActionID{SpellID: 412758},
		SpellSchool:  core.SpellSchoolFire,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagAPL | core.SpellFlagResetAttackSwing | core.SpellFlagBinary | WarlockFlagDestruction,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.14,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2250,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage = sim.Roll(baseLowDamage, baseHighDamage)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			warlock.IncinerateAura.Activate(sim)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
