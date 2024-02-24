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
			character.AddStat(stats.RangedAttackPower, 30)
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

	// Electrocutioner's Needle
	core.NewItemEffect(213286, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(213286)
		ppmm := character.AutoAttacks.NewPPMManager(6.5, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 434839},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 25+0.05*spell.SpellDamage(), spell.OutcomeMagicHitAndCrit)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Electrocutioner's Needle Proc Aura",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && ppmm.Proc(sim, spell.ProcMask, "Electrocutioner's Needle Proc") {
					procSpell.Cast(sim, result.Target)
				}
			},
		})
	})

	// Supercharged Headchopper
	core.NewItemEffect(213296, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(213296)
		ppmm := character.AutoAttacks.NewPPMManager(1.5, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 434842},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(80, 100) + 0.1*spell.SpellDamage()
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Supercharged Headchopper Proc Aura",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && ppmm.Proc(sim, spell.ProcMask, "Supercharged Headchopper Proc") {
					procSpell.Cast(sim, result.Target)
				}
			},
		})
	})

	// Toxic Revenger II
	core.NewItemEffect(213291, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(213291)
		ppmm := character.AutoAttacks.NewPPMManager(3.0, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 435169},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Toxic Revenger II Poison Cloud",
				},
				TickLength:    5 * time.Second,
				NumberOfTicks: 3,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = 30
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType]
					dot.SnapshotCritChance = 0
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				},

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					result := spell.CalcOutcome(sim, aoeTarget, spell.OutcomeMagicHit)
					if result.Landed() {
						spell.Dot(aoeTarget).Apply(sim)
					}
				}
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Toxic Revenger II Proc Aura",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && ppmm.Proc(sim, spell.ProcMask, "Toxic Revenger II Proc") {
					procSpell.Cast(sim, result.Target)
				}
			},
		})
	})

	core.AddEffectsToTest = true
}
