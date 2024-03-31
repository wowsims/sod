package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SignetOfBeasts          = 209823
	BloodlashBow            = 216516
	GurubashiPitFightersBow = 221450
)

func init() {
	core.NewItemEffect(SignetOfBeasts, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		if hunter.pet != nil {
			hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.01
		}
	})

	core.NewItemEffect(BloodlashBow, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		hunter.newBloodlashProcItem(50)
	})

	core.NewItemEffect(GurubashiPitFightersBow, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		hunter.newBloodlashProcItem(75)
	})
}

func (hunter *Hunter) newBloodlashProcItem(bonusStrength float64) {
	procAura := hunter.NewTemporaryStatsAura("Bloodlash", core.ActionID{SpellID: 436471}, stats.Stats{stats.Strength: bonusStrength}, time.Second*15)
	core.MakePermanent(hunter.GetOrRegisterAura(core.Aura{
		Label: "Bloodlash Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && sim.RandomFloat("Bloodlash") < .05 {
				procAura.Activate(sim)
			}
		},
	}))
}
