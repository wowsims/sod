package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	DiamondFlask        = 20130
	Exsanguinar         = 216497
	SuzerainDefender    = 224280
	GrileksCharmOFMight = 231286
	RageOfMugamba       = 231350
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(DiamondFlask, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := character.NewTemporaryStatsAura("Diamond Flask", core.ActionID{SpellID: 24427}, stats.Stats{stats.Strength: 75}, time.Second*60)

		triggerSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 24427},
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 6,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    triggerSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(Exsanguinar, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionId := core.ActionID{SpellID: 436332}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionId,
			SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			DamageMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Exsanguination",
				},
				TickLength:    2 * time.Second,
				NumberOfTicks: 15,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 5, isRollover)
				},

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					// Has no DefenseType, also haven't seen a miss in logs.
					result := spell.CalcAndDealDamage(sim, aoeTarget, 65, spell.OutcomeAlwaysHit)
					if result.Landed() {
						spell.Dot(aoeTarget).Apply(sim)
					}
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(GrileksCharmOFMight, func(agent core.Agent) {
		warrior := agent.(WarriorAgent).GetWarrior()
		actionID := core.ActionID{ItemID: GrileksCharmOFMight}
		rageMetrics := warrior.NewRageMetrics(actionID)

		aura := warrior.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Gri'lek's Guard",
			Duration: time.Second * 20,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warrior.AddStatDynamic(sim, stats.BlockValue, 200)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warrior.AddStatDynamic(sim, stats.BlockValue, -200)
			},
		})

		spell := warrior.Character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warrior.NewTimer(),
					Duration: time.Minute * 2,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				warrior.AddRage(sim, 30, rageMetrics)
				aura.Activate(sim)
			},
		})

		warrior.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	core.NewItemEffect(RageOfMugamba, func(agent core.Agent) {
		warrior := agent.(WarriorAgent).GetWarrior()
		if !warrior.Talents.ShieldSlam {
			return
		}

		warrior.RegisterAura(core.Aura{
			Label: "Reduced Shield Slam Cost (Rage of Mugamba)",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				warrior.ShieldSlam.Cost.FlatModifier -= 5
			},
		})
	})

	core.NewItemEffect(SuzerainDefender, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: SuzerainDefender}

		// Store a reference in case the unit switches targets since we don't have a great way to do this right now
		fightingDragonkin := false
		rageOfSuzerain := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469025},
			Label:    "Rage of the Suzerain",
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.30
					fightingDragonkin = true
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if fightingDragonkin {
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.30
					fightingDragonkin = false
				}
			},
		})

		defenseOfDragonflights := character.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Defense of the Dragonflights",
			Duration: time.Second * 5,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.MultiplySchoolDamageTaken(0.50)
				rageOfSuzerain.Activate(sim)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.MultiplySchoolDamageTaken(2)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 1,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				defenseOfDragonflights.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	core.AddEffectsToTest = true
}
