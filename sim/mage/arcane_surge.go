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

	manaMetrics := mage.NewManaMetrics(actionID)

	manaAura := mage.GetOrRegisterAura(core.Aura{
		Label:    "Arcane Surge",
		ActionID: actionID,
		Duration: auraDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpiritRegenRateCasting += 1
			aura.Unit.PseudoStats.SpiritRegenMultiplier *= 3
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpiritRegenRateCasting -= 1
			aura.Unit.PseudoStats.SpiritRegenMultiplier /= 3
		},
	})

	mage.ArcaneSurge = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
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
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeExpectedMagicHitAndCrit)

			if result.Landed() {
				mage.SpendMana(sim, mage.CurrentMana(), manaMetrics)
				manaAura.Activate(sim)
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.ArcaneSurge,
		Type:  core.CooldownTypeDPS,
	})
}
