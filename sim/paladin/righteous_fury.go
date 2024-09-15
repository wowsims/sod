package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (paladin *Paladin) registerRighteousFury() {
	if !paladin.Options.RighteousFury {
		return
	}
	horRune := proto.PaladinRune_RuneHandsHandOfReckoning
	hasHoR := paladin.hasRune(horRune)
	actionID := core.ActionID{SpellID: core.TernaryInt32(hasHoR, int32(horRune), 25780)}

	rfThreatMultiplier := 0.6 + core.TernaryFloat64(hasHoR, 0.2, 0.0)
	// Improved Righteous Fury is multiplicative.
	rfThreatMultiplier *= 1.0 + []float64{0.0, 0.16, 0.33, 0.5}[paladin.Talents.ImprovedRighteousFury]

	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellSchool.Matches(core.SpellSchoolHoly) {
			spell.ThreatMultiplier *= 1.0 + rfThreatMultiplier
		}
	})

	if !hasHoR { // This is just a visual/UI indicator when we don't have HoR rune.
		paladin.RegisterAura(core.Aura{
			Label:    "Righetous Fury",
			ActionID: actionID,
			Duration: core.NeverExpires, // 30 minutes without HoR rune, but no need to model
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
		})
	} else {
		// Passive effects granted by Hand of Reckoning rune.

		// Damage which takes you below 35% health is reduced by 20% (DR component of WotLK's Ardent Defender)
		rfDamageReduction := 0.2

		handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			incomingDamage := result.Damage
			if (paladin.CurrentHealth()-incomingDamage)/paladin.MaxHealth() <= 0.35 {
				result.Damage -= (paladin.MaxHealth()*0.35 - (paladin.CurrentHealth() - incomingDamage)) * rfDamageReduction
				if sim.Log != nil {
					paladin.Log(sim, "Righteous Fury absorbs %d damage", int32(incomingDamage-result.Damage))
				}
			}
		}

		paladin.AddDynamicDamageTakenModifier(handler)

		// Gives you mana when healed by other friendly targets' spells equal to 25% of the amount healed.
		horManaMetrics := paladin.NewManaMetrics(actionID)

		paladin.RegisterAura(core.Aura{
			Label:    "Righteous Fury",
			ActionID: core.ActionID{SpellID: 407627},
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnHealTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.IsOtherAction(proto.OtherAction_OtherActionHealingModel) {
					manaGained := result.Damage * 0.25
					paladin.AddMana(sim, manaGained, horManaMetrics)
				}
			},
		})
	} // else
}
