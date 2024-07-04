package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SulfurasHandOfRagnaros = 227683 // 17182
	DukesDomain            = 227915
	AccursedChalice        = 228078
	GerminatingPoisonseed  = 228081
	GloamingTreeheart      = 228083
	WoodcarvedMoonstalker  = 228089
	TheMoltenCore          = 228122
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=227683/sulfuras-hand-of-ragnaros
	// Chance on hit: Hurls a fiery ball that causes 273 to 333 Fire damage and purges the target's soul, increasing Fire and Holy damage taken by up to 30 and dealing an additional 75 damage over 10 sec.
	// Equip: 20% chance to deal 25 Fire damage to all nearby enemies when you are struck by a melee attack. (Proc chance: 20%)
	core.NewItemEffect(SulfurasHandOfRagnaros, func(agent core.Agent) {
		character := agent.GetCharacter()

		immolationSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 460335},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			BonusCoefficient: .025,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 25, spell.OutcomeAlwaysHit)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Immolation (Hand of Ragnaros)",
			Callback:   core.CallbackOnSpellHitTaken,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskMelee,
			ProcChance: .20,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					immolationSpell.Cast(sim, aoeTarget)
				}
			},
		})

		debuffAuras := character.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
			return unit.GetOrRegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 460338},
				Label:    "Purged by Fire",
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexFire] += 30
					unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexHoly] += 30
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexFire] -= 30
					unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexHoly] -= 30
				},
			})
		})

		purgedByFireSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 460338},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Purged By Fire",
				},
				TickLength:    2 * time.Second,
				NumberOfTicks: 5,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 15, isRollover)
					debuffAuras.Get(target).Activate(sim)
				},

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealDamage(sim, target, sim.Roll(273, 333), spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Purged by Fire Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee,
			PPM:      1, // Estimated based on data from WoW Armaments Discord
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				purgedByFireSpell.Cast(sim, result.Target)
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=227915/dukes-domain
	// Use: Expand the Duke's Domain, increasing the Fire Resistance of those who reside within by 50. Lasts for 15 sec. (1 Min, 30 Sec Cooldown)
	// TODO: Raid-wide effect if we ever do raid sim
	core.NewSimpleStatDefensiveTrinketEffect(DukesDomain, stats.Stats{stats.FireResistance: 50}, time.Second*15, time.Second*90)

	// https://www.wowhead.com/classic/item=228078/accursed-chalice
	// Use: Increases your Strength by 80.  Effect lasts for 20 sec. (2 Min Cooldown)
	core.NewSimpleStatOffensiveTrinketEffect(AccursedChalice, stats.Stats{stats.Strength: 80}, time.Second*20, time.Minute*2)

	// https://www.wowhead.com/classic/item=228081/germinating-poisonseed
	// Use: Increases your Nature Damage by up to 115.  Effect lasts for 20 sec. (2 Min Cooldown)
	core.NewSimpleStatOffensiveTrinketEffect(GerminatingPoisonseed, stats.Stats{stats.NaturePower: 115}, time.Second*20, time.Minute*2)

	// https://www.wowhead.com/classic/item=228083/gloaming-treeheart
	// Use: Increases your Nature Resistance by 90.  Effect lasts for 30 sec. (3 Min Cooldown)
	core.NewSimpleStatDefensiveTrinketEffect(GloamingTreeheart, stats.Stats{stats.NatureResistance: 90}, time.Second*30, time.Minute*3)

	// https://www.wowhead.com/classic/item=228089/woodcarved-moonstalker
	// Use: Increases your Strength by 60.  Effect lasts for 15 sec. (1 Min, 30 Sec Cooldown)
	core.NewItemEffect(WoodcarvedMoonstalker, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.PseudoStats.BonusDamage += 4

		aura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: WoodcarvedMoonstalker},
			Label:    "Woodcarved Moonstalker",
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.Strength, 60)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.Strength, -60)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: WoodcarvedMoonstalker},
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 90,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	// https://www.wowhead.com/classic/item=228122/the-molten-core
	// Equip: Inflicts 20 Fire damage to nearby enemies every 2 sec.
	core.NewItemEffect(TheMoltenCore, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 461228}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			BonusCoefficient: .045,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 20, spell.OutcomeAlwaysHit)
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "The Molten Core Trigger",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:   time.Second * 2,
					Priority: core.ActionPriorityDOT, // High prio
					OnAction: func(sim *core.Simulation) {
						for _, aoeTarget := range sim.Encounter.TargetUnits {
							spell.Cast(sim, aoeTarget)
						}
					},
				})
			},
		})
	})

	core.AddEffectsToTest = true
}
