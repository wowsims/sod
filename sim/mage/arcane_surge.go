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

	level := float64(mage.GetCharacter().Level)
	actionID := core.ActionID{SpellID: int32(proto.MageRune_RuneLegsArcaneSurge)}
	baseCalc := (13.828124 + 0.018012*level + 0.044141*level*level)
	baseDamageLow := baseCalc * 2.26
	baseDamageHigh := baseCalc * 2.64
	spellCoeff := .429
	cooldown := time.Minute * 2
	auraDuration := time.Second * 8

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
		ActionID:    actionID,
		SpellCode:   SpellCode_MageArcaneSurge,
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

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
		CritMultiplier:   mage.MageCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			// Damage increased based on remaining mana up to 300%
			damage *= 1 + mage.CurrentManaPercent()*3
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			// Because of the 0 base mana cost we have to create resource metrics
			mage.SpendMana(sim, mage.CurrentMana(), mage.NewManaMetrics(actionID))
			manaAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.ArcaneSurge,
		Type:  core.CooldownTypeDPS,
	})
}
