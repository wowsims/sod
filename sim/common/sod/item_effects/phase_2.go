package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Cloth
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Leather
	///////////////////////////////////////////////////////////////////////////

	// Machinist's Gloves
	core.NewItemEffect(213319, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeMechanical {
			character.AddStat(stats.AttackPower, 30)
		}
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Mail
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Plate
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// Miniaturized Combustion Chamber
	core.NewItemEffect(213347, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 435167}
		manaMetrics := character.NewManaMetrics(actionID)

		manaRoll := 0.0
		dmgRoll := 0.0

		regChannel := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagChanneled | core.SpellFlagAPL,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 30,
				},
			},

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Miniaturized Combustion Chamber",
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						character.AutoAttacks.CancelAutoSwing(sim)
						manaRoll = sim.RollWithLabel(1, 150, "Miniaturized Combustion Chamber")
						dmgRoll = sim.RollWithLabel(1, 150, "Miniaturized Combustion Chamber")
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						character.AutoAttacks.EnableAutoSwing(sim)
					},
				},
				SelfOnly:      true,
				NumberOfTicks: 10,
				TickLength:    time.Second,

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					character.AddMana(sim, manaRoll, manaMetrics)
					character.RemoveHealth(sim, dmgRoll)
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				spell.SelfHot().Apply(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    regChannel,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeMana,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// Automatic Crowd Pummeler
	core.NewItemEffect(210741, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 13494}

		hasteAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Haste",
			ActionID: actionID,
			Duration: time.Second * 30,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.5)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.0/1.5)
			},
		})

		hasteSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				hasteAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    hasteSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Bloodbark Crusher
	core.NewItemEffect(216499, func(agent core.Agent) {
		character := agent.GetCharacter()
		auraActionID := core.ActionID{SpellID: 436482}
		numHits := min(3, character.Env.GetNumTargets())

		triggeredDmgSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 436481},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			// TODO: "Causes additional threat" in Tooltip, no clue what the multiplier is.
			ThreatMultiplier: 1,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					// The spell is affected by phys school mods because it's nature + physical school.
					dmgWithPhysMods := 5 * character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
					spell.CalcAndDealDamage(sim, curTarget, dmgWithPhysMods, spell.OutcomeMagicHitAndCrit)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		})

		mainAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Bloodbark Cleave",
			ActionID: auraActionID,
			Duration: 20 * time.Second,

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask&core.ProcMaskMelee != 0 {
					triggeredDmgSpell.Cast(sim, result.Target)
					return
				}
			},
		})

		mainSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: auraActionID,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				mainAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    mainSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.AddEffectsToTest = true
}
