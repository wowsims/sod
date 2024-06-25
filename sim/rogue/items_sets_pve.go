package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var _ = core.NewItemSet(core.ItemSet{
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
				return target.RegisterAura(core.Aura{
					Label:     "Blood Corruption",
					ActionID:  core.ActionID{SpellID: 449927},
					Duration:  time.Second * 15,
					MaxStacks: 30,

					OnGain: func(aura *core.Aura, sim *core.Simulation) {
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

			core.MakePermanent(rogue.RegisterAura(core.Aura{
				Label:    "Blood Corrupting" + rogue.Label,
				ActionID: core.ActionID{SpellID: 449928},
				OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRangedSpecial) {
						return
					}

					switch spell {
					case rogue.Backstab, rogue.Mutilate, rogue.SinisterStrike, rogue.SaberSlash, rogue.Shiv, rogue.PoisonedKnife, rogue.MainGauche, rogue.QuickDraw:
						aura := procAuras.Get(result.Target)
						aura.Activate(sim)
						aura.SetStacks(sim, aura.MaxStacks)
					}
				},
			}))
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////
