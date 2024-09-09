package item_sets

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                                 Cloth
///////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////
//                                 Other
///////////////////////////////////////////////////////////////////////////

var ItemSetMajorMojoInfusion = core.NewItemSet(core.ItemSet{
	Name: "Major Mojo Infusion",
	Bonuses: map[int32]core.ApplyEffect{
		// +30 Attack Power.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStats(stats.Stats{
				stats.AttackPower:       30,
				stats.RangedAttackPower: 30,
			})
		},
	},
})

var ItemSetOverlordsResolution = core.NewItemSet(core.ItemSet{
	Name: "Overlord's Resolution",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +8.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 8)
		},
	},
})

var ItemSetPrayerOfThePrimal = core.NewItemSet(core.ItemSet{
	Name: "Prayer of the Primal",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases healing done by up to 33 and damage done by up to 11 for all magical spells and effects.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStats(stats.Stats{
				stats.HealingPower: 33,
				stats.SpellDamage:  11,
			})
		},
	},
})

var ItemSetPrimalBlessing = core.NewItemSet(core.ItemSet{
	Name: "Primal Blessing",
	Bonuses: map[int32]core.ApplyEffect{
		// Grants a small chance when ranged or melee damage is dealt to infuse the wielder with a blessing from the Primal Gods.
		// Ranged and melee attack power increased by 300 for 12 sec.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			aura := character.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 467742},
				Label:    "Primal Blessing",
				Duration: time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					character.AddStatsDynamic(sim, stats.Stats{
						stats.AttackPower:       300,
						stats.RangedAttackPower: 300,
					})
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					character.AddStatsDynamic(sim, stats.Stats{
						stats.AttackPower:       -300,
						stats.RangedAttackPower: -300,
					})
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Primal Blessing Trigger",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.05,
				ICD:        time.Second * 72,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					aura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetTwinBladesofHakkari = core.NewItemSet(core.ItemSet{
	Name: "The Twin Blades of Hakkari",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases Swords +3
		// 2% chance on melee hit to gain 1 extra attack.  (1%, 100ms cooldown)
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SwordsSkill += 3
			if !character.AutoAttacks.AutoSwingMelee {
				return
			}

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:              "Twin Blades of the Hakkari",
				Callback:          core.CallbackOnSpellHitDealt,
				Outcome:           core.OutcomeLanded,
				ProcMask:          core.ProcMaskMelee,
				SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
				ProcChance:        0.02,
				ICD:               time.Millisecond * 100,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					spell.Unit.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 468255}, spell)
				},
			})
		},
	},
})

var ItemSetZanzilsConcentration = core.NewItemSet(core.ItemSet{
	Name: "Zanzil's Concentration",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 6.
		// Improves your chance to hit with all spells and attacks by 1%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStats(stats.Stats{
				stats.SpellPower: 6,
				stats.SpellHit:   1 * core.SpellHitRatingPerHitChance,
				stats.MeleeHit:   1 * core.MeleeHitRatingPerHitChance,
			})
		},
	},
})
