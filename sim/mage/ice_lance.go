package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: Classic review ice lance numbers on live
func (mage *Mage) registerIceLanceSpell() {
	if !mage.HasRune(proto.MageRune_RuneHandsIceLance) {
		return
	}

	hasWintersChillTalent := mage.Talents.WintersChill > 0

	baseDamageLow := mage.baseRuneAbilityDamage() * 0.55
	baseDamageHigh := mage.baseRuneAbilityDamage() * 0.65
	spellCoeff := 0.429
	manaCost := 0.08

	damageModPct := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  ClassSpellMask_MageIceLance,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 1,
	})

	mage.IceLance = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: int32(proto.MageRune_RuneHandsIceLance)},
		ClassSpellMask: ClassSpellMask_MageIceLance,
		SpellSchool:    core.SpellSchoolFrost,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,

		MissileSpeed: 38,
		MetricSplits: 11, // Possible 8 total stacks

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if !hasWintersChillTalent {
					return
				}

				if glaciateAura := mage.GlaciateAuras.Get(mage.CurrentTarget); glaciateAura != nil {
					spell.SetMetricsSplit(glaciateAura.GetStacks())
				}
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)

			damageModPct.UpdateFloatValue(core.TernaryFloat64(mage.isTargetFrozen(target), 4, 1))
			damageModPct.Activate()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	if !hasWintersChillTalent {
		return
	}

	mage.GlaciateAuras = mage.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
		return unit.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 1218345},
			Label:     "Glaciate",
			Duration:  time.Second * 15,
			MaxStacks: 5,
		})
	})

	iceLanceDamageMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  ClassSpellMask_MageIceLance,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 1.0,
	})

	// For some reason they set up Deep Freeze as base points
	deepFreezeDamageMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_MageDeepFreeze,
		Kind:      core.SpellMod_BaseDamageDone_Flat,
	})

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Glaciate",
		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_MageIceLance | ClassSpellMask_MageDeepFreeze) {
				multiplier := 1.0
				if glaciateAura := mage.GlaciateAuras.Get(target); glaciateAura.IsActive() {
					multiplier += 0.20 * float64(glaciateAura.GetStacks())
				}

				iceLanceDamageMod.UpdateFloatValue(multiplier)
				iceLanceDamageMod.Activate()
				deepFreezeDamageMod.UpdateIntValue(int64(multiplier * 100))
				deepFreezeDamageMod.Activate()
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Note: Glaciate is not removed when using Deep Freeze even though it gains damage
			if glaciateAura := mage.GlaciateAuras.Get(result.Target); spell.Matches(ClassSpellMask_MageIceLance) && result.Landed() && glaciateAura.IsActive() {
				glaciateAura.Deactivate(sim)
			}
		},
	}))

	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:             "Glaciate Trigger",
		ClassSpellMask:   ClassSpellMask_MageAll ^ ClassSpellMask_MageIceLance,
		Callback:         core.CallbackOnSpellHitDealt,
		Outcome:          core.OutcomeLanded,
		SpellSchool:      core.SpellSchoolFrost,
		Harmful:          true,
		CanProcFromProcs: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			glaciateAura := mage.GlaciateAuras.Get(result.Target)
			glaciateAura.Activate(sim)
			glaciateAura.AddStack(sim)
		},
	})
}
