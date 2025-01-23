package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (mage *Mage) registerArcaneSurgeSpell() {
	if !mage.HasRune(proto.MageRune_RuneLegsArcaneSurge) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.MageRune_RuneLegsArcaneSurge)}
	baseDamageLow := mage.baseRuneAbilityDamage() * 2.26
	baseDamageHigh := mage.baseRuneAbilityDamage() * 2.64
	spellCoeff := .429
	cooldown := time.Minute * 2
	auraDuration := time.Second * 8
	damageMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_MageArcaneSurge,
	})

	manaMetrics := mage.NewManaMetrics(actionID)

	manaAura := mage.GetOrRegisterAura(core.Aura{
		Label:    "Arcane Surge",
		ActionID: actionID,
		Duration: auraDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpiritRegenMultiplier *= 3
			mage.PseudoStats.ForceFullSpiritRegen = true
			mage.UpdateManaRegenRates()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpiritRegenMultiplier /= 3
			mage.PseudoStats.ForceFullSpiritRegen = false
			mage.UpdateManaRegenRates()
		},
	})

	mage.ArcaneSurge = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_MageArcaneSurge,
		SpellSchool:    core.SpellSchoolArcane,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost: 0.0, // Drains remaining mana so we have to use ModifyCast
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamageLow, baseDamageHigh)
			// Damage increased based on remaining mana up to 300%
			damageMod.UpdateFloatValue(mage.CurrentManaPercent() * 3)
			damageMod.Activate()
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			damageMod.Deactivate()
			// Because of the 0 base mana cost we have to create resource metrics
			mage.SpendMana(sim, mage.CurrentMana(), manaMetrics)
			manaAura.Activate(sim)
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.CurrentMana() > 0
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.ArcaneSurge,
		Type:  core.CooldownTypeDPS,
	})
}
