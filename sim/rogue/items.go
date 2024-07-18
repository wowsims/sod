package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	BloodSpatteredStilletto = 216522
	ShadowflameSword        = 228143
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(BloodSpatteredStilletto, func(agent core.Agent) {
		character := agent.GetCharacter()

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 436477},
			SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			DamageMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					//Confirmed always hits through logs
					spell.CalcAndDealDamage(sim, aoeTarget, 140, spell.OutcomeAlwaysHit)
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// https://www.wowhead.com/classic/item=228143/shadowflame-sword
	core.NewItemEffect(ShadowflameSword, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()

		if !rogue.Talents.BladeFlurry {
			return
		}

		// This is treated as a buff, NOT a debuff in-game
		procAura := rogue.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 461252},
			Label:    "Shadowflame Fury",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				for _, target := range sim.Encounter.TargetUnits {
					target.AddStatDynamic(sim, stats.Armor, -2000)
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				for _, target := range sim.Encounter.TargetUnits {
					target.AddStatDynamic(sim, stats.Armor, 2000)
				}
			},
		})

		core.MakePermanent(rogue.RegisterAura(core.Aura{
			Label: "Shadowflame Fury Trigger",
			OnCastComplete: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.SpellCode == SpellCode_RogueBladeFlurry {
					procAura.Duration = rogue.BladeFlurryAura.Duration
					procAura.Activate(sim)
				}
			},
		}))
	})

	core.AddEffectsToTest = true
}
