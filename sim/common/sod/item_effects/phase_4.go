package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const (
	// Ordered by ID
	HandOfJustice = 227989
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(HandOfJustice, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.AutoSwingMelee {
			return
		}

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 2,
		}

		character.RegisterAura(core.Aura{
			Label:    "Hand of Justice",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if sim.RandomFloat("HandOfJustice") < 0.02 {
					icd.Use(sim)
					aura.Unit.AutoAttacks.ExtraMHAttack(sim, 1, core.ActionID{SpellID: 15600})
				}
			},
		})
	})

	core.AddEffectsToTest = true
}
