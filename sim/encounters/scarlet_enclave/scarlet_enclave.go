package scarlet_enclave

import (
	"github.com/wowsims/sod/sim/core"
)

func Register() {
	addScarlet60("SoD/Scarlet Enclave")
}

type ScarletEnclaveEncounter struct{}

func (scarletEnclaveEncounter *ScarletEnclaveEncounter) registerScarletDominionAura(target *core.Target) {
	charactertarget := target.Env.Raid.Parties[0].Players[0].GetCharacter()

	core.MakePermanent(charactertarget.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1232014},
		Label:    "Scarlet Dominion",
		// Have to use OnGain/OnExpire since AttackTables aren't reset, so the BaseMissChance carries over
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DodgeReduction += 0.2

			for _, target := range sim.Encounter.TargetUnits {
				for _, at := range target.AttackTables[aura.Unit.UnitIndex] {
					at.BaseMissChance -= 0.05
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DodgeReduction -= 0.2

			for _, target := range sim.Encounter.TargetUnits {
				for _, at := range target.AttackTables[aura.Unit.UnitIndex] {
					at.BaseMissChance += 0.05
				}
			}
		},
	}))
}
