package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) applyAncestralGuidance() {
	if !shaman.HasRune(proto.ShamanRune_RuneLegsAncestralGuidance) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsAncestralGuidance)}
	duration := time.Second * 10
	cooldown := time.Minute * 2

	damageConversion := .25  // 25%
	healingConversion := 1.0 // 100%
	numHealedAllies := int32(3)

	agDamageSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 409337}, // AG Damage has its own Spell ID
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagIgnoreResists,

		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: shaman.ShamanThreatMultiplier(1),
	})

	agHealSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 409333}, // AG Damage has its own Spell ID
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,
	})

	agAura := shaman.RegisterAura(core.Aura{
		Label:    "Ancestral Guidance",
		ActionID: actionID,
		Duration: duration,

		OnDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == agDamageSpell {
				return
			}

			targets := sim.Environment.Raid.GetFirstNPlayersOrPets(numHealedAllies)

			for hitIndex := int32(0); hitIndex < int32(len(targets)); hitIndex++ {
				target := targets[hitIndex]
				baseHealing := result.Damage * damageConversion
				agHealSpell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
			}
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == agHealSpell || shaman.lastFlameShockTarget == nil {
				return
			}

			baseDamage := result.Damage * healingConversion
			agDamageSpell.CalcAndDealDamage(sim, shaman.lastFlameShockTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})

	agCDSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       core.SpellFlagAPL | core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			agAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: agCDSpell,
		Type:  core.CooldownTypeDPS | core.CooldownTypeSurvival,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentManaPercent() <= 0.2
		},
	})
}
