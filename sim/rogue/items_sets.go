package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetBloodCorruptedLeathers = core.NewItemSet(core.ItemSet{
	Name: "Blood Corrupted Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()

			procAuras := rogue.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
				return target.GetOrRegisterAura(core.Aura{
					Label:     "Blood Corruption Proc",
					ActionID:  core.ActionID{SpellID: 449927},
					Duration:  time.Second * 15,
					MaxStacks: 30,

					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						aura.SetStacks(sim, aura.MaxStacks)

						for si := stats.SchoolIndexPhysical; si < stats.SchoolLen; si++ {
							aura.Unit.PseudoStats.SchoolBonusDamageTaken[si] += 7
						}
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						for si := stats.SchoolIndexPhysical; si < stats.SchoolLen; si++ {
							aura.Unit.PseudoStats.SchoolBonusDamageTaken[si] -= 7
						}
					},
					OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if result.Landed() && spell.ProcMask.Matches(core.ProcMaskDirect) {
							aura.RemoveStack(sim)
						}
					},
				})
			})

			handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == rogue.Backstab || spell == rogue.SinisterStrike {
					procAuras.Get(result.Target).Activate(sim)
				}
			}

			core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{SpellID: 449919},
				Name:       "Blood Corruption",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskDirect,
				Outcome:    core.OutcomeLanded,
				ProcChance: 1,
				Handler:    handler,
			})
		},
	},
})
