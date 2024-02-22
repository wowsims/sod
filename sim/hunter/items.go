package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	core.NewItemEffect(209823, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		if hunter.pet != nil {
			hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.01
		}
	})

	core.NewItemEffect(216516, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()

		procAura := hunter.NewTemporaryStatsAura("Bloodlash", core.ActionID{SpellID: 436471}, stats.Stats{stats.Strength: 50}, time.Second*15)
		ppm := hunter.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMeleeOrRanged)

		core.MakePermanent(hunter.GetOrRegisterAura(core.Aura{
			Label: "Bloodlash Trigger",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if !ppm.Proc(sim, spell.ProcMask, "Bloodlash Proc") {
					return
				}

				procAura.Activate(sim)
			},
		}))
	})
}
