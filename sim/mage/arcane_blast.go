package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: Classic verify Arcane Blast rune numbers
// https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044#news-post-336044
// https://www.wowhead.com/classic/spell=400574/arcane-blast
func (mage *Mage) registerArcaneBlastSpell() {
	if !mage.HasRune(proto.MageRune_RuneHandsArcaneBlast) {
		return
	}

	// Spells that consume Arcane Blast stacks
	affectedSpells := ClassSpellMask_MageArcaneBarrage | ClassSpellMask_MageArcaneExplosion | ClassSpellMask_MageArcaneSurge | ClassSpellMask_MageBalefireBolt | ClassSpellMask_MageSpellfrostBolt

	baseLowDamage := mage.baseRuneAbilityDamage() * 4.53
	baseHighDamage := mage.baseRuneAbilityDamage() * 5.27
	spellCoeff := .714
	castTime := time.Millisecond * 2500
	manaCost := .07

	mage.ArcaneBlastDamageMultiplier = 0.15

	abDamageModFlat := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_MageArcaneBlastAuraFlat,
		Kind:      core.SpellMod_DamageDone_Flat,
	})
	abDamageModPct := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  ClassSpellMask_MageArcaneBlastAuraPct,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 1,
	})

	abCostMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_MageArcaneBlast,
		Kind:      core.SpellMod_PowerCost_Pct,
	})

	mage.ArcaneBlastAura = mage.GetOrRegisterAura(core.Aura{
		Label:     "Arcane Blast Aura",
		ActionID:  core.ActionID{SpellID: 400573},
		Duration:  time.Second * 6,
		MaxStacks: 4,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			abCostMod.UpdateIntValue(175 * int64(newStacks))
			abDamageModFlat.UpdateIntValue(int64(mage.ArcaneBlastDamageMultiplier*100) * int64(newStacks))
			abDamageModPct.UpdateFloatValue((1 + mage.ArcaneBlastDamageMultiplier*float64(newStacks)))
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_MageAll) && spell.Matches(affectedSpells) && (mage.ArcaneTunnelingAura == nil || !mage.ArcaneTunnelingAura.IsActive()) {
				aura.Deactivate(sim)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			abCostMod.Activate()
			abDamageModFlat.Activate()
			abDamageModPct.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			abCostMod.Deactivate()
			abDamageModFlat.Deactivate()
			abDamageModPct.Deactivate()
		},
	})

	mage.ArcaneBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 400574},
		ClassSpellMask: ClassSpellMask_MageArcaneBlast,
		SpellSchool:    core.SpellSchoolArcane,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			mage.ArcaneBlastAura.Activate(sim)
			mage.ArcaneBlastAura.AddStack(sim)
		},
	})
}
