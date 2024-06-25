package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetNightmareProphetsGarb = core.NewItemSet(core.ItemSet{
	Name: "Nightmare Prophet's Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			warlock.shadowSparkAura = warlock.GetOrRegisterAura(core.Aura{
				Label:     "Shadow Spark Proc",
				ActionID:  core.ActionID{SpellID: 450013},
				Duration:  time.Second * 12,
				MaxStacks: 2,
			})

			core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
				Label: "Shadow Spark",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellCode == SpellCode_WarlockShadowCleave && result.Landed() {
						warlock.shadowSparkAura.Activate(sim)
						warlock.shadowSparkAura.AddStack(sim)
					}
				},
			}))
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetDeathmistRaiment = core.NewItemSet(core.ItemSet{
	Name: "Deathmist Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Your melee autoattacks and spellcasts have a 6% chance to heal you for 270 to 330 health.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450585}
			manaMetrics := c.NewManaMetrics(core.ActionID{SpellID: 450583})

			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if c.HasManaBar() {
					c.AddMana(sim, sim.Roll(270, 300), manaMetrics)
				}
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Heal Proc on Cast - Dreadmist Raiment (Melee Auto)",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskWhiteHit,
				ProcChance: 0.06,
				Handler:    handler,
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Heal Proc on Cast - Dreadmist Raiment (Spell Cast)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.06,
				Handler:    handler,
			})
		},
		// +8 All Resistances.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.ArcaneResistance: 8,
				stats.FireResistance:   8,
				stats.FrostResistance:  8,
				stats.NatureResistance: 8,
				stats.ShadowResistance: 8,
			})
		},
		// +200 Armor.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
	},
})
