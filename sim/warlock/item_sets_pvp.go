package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetChampionsThreads = core.NewItemSet(core.ItemSet{
	Name: "Champion's Threads",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Reduces the casting time of your Immolate spell by 0.2 sec.
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.GetOrRegisterAura(core.Aura{
				Label: "Immolate Cast Time Reduction",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.Immolate {
						spell.DefaultCast.CastTime -= time.Millisecond * 200
					}
				},
			})
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

var ItemSetLieutenantCommandersThreads = core.NewItemSet(core.ItemSet{
	Name: "Lieutenant Commander's Threads",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Reduces the casting time of your Immolate spell by 0.2 sec.
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.GetOrRegisterAura(core.Aura{
				Label:    "Immolate Cast Time Reduction",
				ActionID: core.ActionID{SpellID: 23047},
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range warlock.Immolate {
						spell.DefaultCast.CastTime -= time.Millisecond * 200
					}
				},
			})
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

// /////////////////////////////////////////////////////////////////////////
//
//	SoD Phase 5 Item Sets
//
// /////////////////////////////////////////////////////////////////////////
func (warlock *Warlock) applyPhase5PvP3PBonus() {
	label := "Immolate Cast Time Reduction"
	if warlock.HasAura(label) {
		return
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		ClassMask: ClassSpellMask_WarlockImmolate,
		TimeValue: -time.Millisecond * 200,
	}))
}

var ItemSetWarlordsThreads = core.NewItemSet(core.ItemSet{
	Name: "Warlord's Threads",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Reduces the casting time of your Immolate spell by 0.2 sec.
		3: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyPhase5PvP3PBonus()
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
	},
})

var ItemSetFieldMarshalsThreads = core.NewItemSet(core.ItemSet{
	Name: "Field Marshal's Threads",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Stamina.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Reduces the casting time of your Immolate spell by 0.2 sec.
		3: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyPhase5PvP3PBonus()
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
	},
})
