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
	has6pcT3 := paladin.HasSetBonus(ItemSetRedemptionBulwark, 6)

	actionID := core.ActionID{SpellID: core.TernaryInt32(hasHoR, int32(horRune), 25780)}

	rfThreatMultiplier := 1.6 + core.TernaryFloat64(hasHoR, 0.2, 0.0)
	// Improved Righteous Fury is multiplicative.
	rfThreatMultiplier *= 1.0 + []float64{0.0, 0.16, 0.33, 0.5}[paladin.Talents.ImprovedRighteousFury]

	paladin.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_Threat_Pct,
		School:     core.SpellSchoolHoly,
		FloatValue: rfThreatMultiplier,
	})

	auraConfig := core.MakePermanent(&core.Aura{Label: "Righteous Fury", ActionID: actionID})

	// Passive effects granted by Hand of Reckoning rune; only active if Righteous Fury is on.
	if hasHoR {

		// Damage which takes you below 35% health is reduced by 20% (DR component of WotLK's Ardent Defender)
		rfDamageReduction := 0.2

		handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			incomingDamage := result.Damage
			if (paladin.CurrentHealth()-incomingDamage)/paladin.MaxHealth() <= 0.35 {
				damageReduction := rfDamageReduction

				if has6pcT3 && spell.Unit.MobType == proto.MobType_MobTypeUndead {
					damageReduction = 0.5
				}

				result.Damage -= (paladin.MaxHealth()*0.35 - (paladin.CurrentHealth() - incomingDamage)) * damageReduction
				if sim.Log != nil {
					paladin.Log(sim, "Righteous Fury absorbs %d damage", int32(incomingDamage-result.Damage))
				}
			}
		}

		paladin.AddDynamicDamageTakenModifier(handler)

		// Gives you mana when healed by other friendly targets' spells equal to 25% of the amount healed.
		horManaMetrics := paladin.NewManaMetrics(actionID)

		auraConfig.OnHealTaken = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.IsOtherAction(proto.OtherAction_OtherActionHealingModel) {
				manaGained := result.Damage * 0.25
				paladin.AddMana(sim, manaGained, horManaMetrics)
			}
		}
	}

	paladin.righteousFuryAura = paladin.RegisterAura(*auraConfig)
}
