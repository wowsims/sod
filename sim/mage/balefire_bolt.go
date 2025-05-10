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
	maxStacks := 5
	stackMultiplier := 0.20
	stackModifier := int32(25)

	damageMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_MageBalefireBolt | ClassSpellMask_MagePyroblast,
	})

	statDeps := make([]*stats.StatDependency, maxStacks+1) // 5 stacks + zero conditions
	for i := 0; i < maxStacks+1; i++ {
		statDeps[i] = mage.NewDynamicMultiplyStat(stats.Spirit, 1.0-stackMultiplier*float64(i))
	}

	mage.BalefireAura = mage.RegisterAura(core.Aura{
		Label:     "Balefire Bolt (Stacks)",
		ActionID:  core.ActionID{SpellID: int32(proto.MageRune_RuneBracersBalefireBolt)}.WithTag(1),
		Duration:  buffDuration,
		MaxStacks: int32(maxStacks),
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			damageMod.UpdateIntValue(int64(stackModifier * newStacks))
			aura.Unit.DisableDynamicStatDep(sim, statDeps[oldStacks])
			aura.Unit.EnableDynamicStatDep(sim, statDeps[newStacks])
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_MageBalefireBolt) && aura.GetStacks() == 5 {
				mage.RemoveHealth(sim, mage.CurrentHealth(), mage.DamageTakenHealthMetrics)
				if sim.Log != nil {
					mage.Log(sim, "YOU DIED")
				}

				sim.Cleanup()
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
	})

	mage.BalefireBolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: int32(proto.MageRune_RuneBracersBalefireBolt)},
		ClassSpellMask: ClassSpellMask_MageBalefireBolt,
		SpellSchool:    core.SpellSchoolArcane | core.SpellSchoolFire | core.SpellSchoolFrost,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,

		MissileSpeed: 24,

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

			mage.BalefireAura.Activate(sim)
			mage.BalefireAura.AddStack(sim)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},

		RelatedSelfBuff: mage.BalefireAura,
	})
}
