package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetBenevolentProphetsVestments = core.NewItemSet(core.ItemSet{
	Name: "Benevolent Prophet's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		// Restores 4 mana per 5 sec.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 4)
		},
		// Your Holy damage spells cause you to gain 60 increased damage and healing power for 15 sec.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()

			procAura := c.NewTemporaryStatsAura("Faith and Magic Proc", core.ActionID{SpellID: 449923}, stats.Stats{stats.SpellPower: 60}, time.Second*15)

			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if spell.SpellSchool.Matches(core.SpellSchoolHoly) {
					procAura.Activate(sim)
				}
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "Faith and Magic",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage,
				ProcChance: 1,
				Handler:    handler,
			})
		},
	},
})
