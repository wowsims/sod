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

			castTimeMod := warlock.AddDynamicMod(core.SpellModConfig{
				Kind:      core.SpellMod_CastTime_Pct,
				ClassMask: ClassSpellMask_WarlockImmolate | ClassSpellMask_WarlockShadowflame,
			})

			buffAura := warlock.GetOrRegisterAura(core.Aura{
				Label:     "Shadow Spark Proc",
				ActionID:  core.ActionID{SpellID: 450013},
				Duration:  time.Second * 12,
				MaxStacks: 2,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					castTimeMod.UpdateFloatValue(-0.50 * float64(newStacks))
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					castTimeMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					castTimeMod.Deactivate()
				},
			})

			core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
				Label: "Shadow Spark",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.Matches(ClassSpellMask_WarlockShadowCleave) && result.Landed() {
						buffAura.Activate(sim)
						buffAura.AddStack(sim)
					}
				},
			}))
		},
	},
})
