package sod

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
)

func init() {
	core.AddEffectsToTest = false

	// Proc effects. Keep these in order by item ID.

	// Shawarmageddon
	core.NewItemEffect(213105, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 434488}

		fireStrike := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 434488},
			SpellSchool:      core.SpellSchoolFire,
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 7.0, spell.OutcomeMagicHitAndCrit)
			},
		})

		spicyAura := character.RegisterAura(core.Aura{
			Label:    "Spicy!",
			ActionID: actionID,
			Duration: time.Second * 30,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.04)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1/1.04)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if result.Landed() {
					fireStrike.Cast(sim, spell.Unit.CurrentTarget)
				}
			},
		})

		spicy := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Cast: core.CastConfig{
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				spicyAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spicy,
			Type:  core.CooldownTypeDPS,
		})
	})

	// Mekkatorque's Arcano-Shredder
	itemhelpers.CreateWeaponProcSpell(213409, "Mekkatorque", 5.0, func(character *core.Character) *core.Spell {
		procAuras := character.NewEnemyAuraArray(core.MekkatorqueFistDebuffAura)

		return character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 434841},
			SpellSchool: core.SpellSchoolArcane,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 30+0.05*spell.SpellDamage(), spell.OutcomeMagicHitAndCrit)
				procAuras.Get(target).Activate(sim)
			},
		})
	})

	core.AddEffectsToTest = true
}
