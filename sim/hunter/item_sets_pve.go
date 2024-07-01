package hunter

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetDreadHuntersChain = core.NewItemSet(core.ItemSet{
	Name: "Dread Hunter's Chain",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
			c.AddStat(stats.RangedAttackPower, 20)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddBonusRangedCritRating(2)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetBeastmasterArmor = core.NewItemSet(core.ItemSet{
	Name: "Beastmaster Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Your melee and ranged autoattacks have a 6% chance to energize you for 300 mana.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450577}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Mana Proc on Cast - Beaststalker Armor",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskWhiteHit,
				ProcChance: 0.06,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasManaBar() {
						c.AddMana(sim, 300, manaMetrics)
					}
				},
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
