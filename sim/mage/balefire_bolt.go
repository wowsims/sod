package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// https://www.wowhead.com/classic/spell=428878/balefire-bolt
func (mage *Mage) registerBalefireBoltSpell() {
	if !mage.HasRune(proto.MageRune_RuneBracersBalefireBolt) {
		return
	}

	baseDamageLow := mage.baseRuneAbilityDamage() * 2.8
	baseDamageHigh := mage.baseRuneAbilityDamage() * 4.2
	spellCoeff := .857
	castTime := time.Millisecond * 2500
	buffDuration := time.Second * 30
	manaCost := .20

	statDeps := make([]*stats.StatDependency, 11) // 10 stacks + zero conditions
	for i := 1; i < 11; i++ {
		statDeps[i] = mage.NewDynamicMultiplyStat(stats.Spirit, 1.0-.1*float64(i))
	}

	balefireAura := mage.RegisterAura(core.Aura{
		Label:     "Balefire Bolt (Stacks)",
		ActionID:  core.ActionID{SpellID: int32(proto.MageRune_RuneBracersBalefireBolt)}.WithTag(1),
		Duration:  buffDuration,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			mage.BalefireBolt.DamageMultiplierAdditive -= .2 * float64(oldStacks)
			mage.BalefireBolt.DamageMultiplierAdditive += .2 * float64(newStacks)

			if oldStacks != 0 {
				aura.Unit.DisableDynamicStatDep(sim, statDeps[oldStacks])
			}
			if newStacks != 0 {
				aura.Unit.EnableDynamicStatDep(sim, statDeps[newStacks])
			}

			if newStacks == 10 {
				mage.RemoveHealth(sim, mage.CurrentHealth())

				if sim.Log != nil {
					mage.Log(sim, "YOU DIED")
				}

				sim.Cleanup()
			}
		},
	})

	mage.BalefireBolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.MageRune_RuneBracersBalefireBolt)},
		SpellCode:   SpellCode_MageBalefireBolt,
		SpellSchool: core.SpellSchoolArcane | core.SpellSchoolFire | core.SpellSchoolFrost,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,
		// TODO: Verify missile speed
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

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			balefireAura.Activate(sim)
			balefireAura.AddStack(sim)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
