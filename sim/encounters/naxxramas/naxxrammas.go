package naxxramas

import (
	"github.com/wowsims/sod/sim/core"
	"fmt"
)

func Register() {
	addNaxx60("SoD/Naxxrammas")
	addPatchwerk("SoD/Naxxrammas")
	addThaddius("SoD/Naxxrammas")
	addLoatheb("SoD/Naxxrammas")
}


type NaxxramasEncounter struct {
		authorityFrozenWastesStacks   *int32
		authorityFrozenWastesAura     *core.Aura
}

func (naxxEncounter *NaxxramasEncounter) registerAuthorityOfTheFrozenWastesAura(target *core.Target, stacks int32) *core.Aura {
	charactertarget := target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	fmt.Println("Stacks:", stacks)
	return core.MakePermanent(charactertarget.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1218283},
		Label:     "Authority of the Frozen Wastes",
		MaxStacks: 4,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
			aura.SetStacks(sim, stacks)
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			aura.Unit.PseudoStats.DodgeReduction += 0.04 * float64(newStacks-oldStacks)

			for _, target := range sim.Encounter.TargetUnits {
				for _, at := range target.AttackTables[aura.Unit.UnitIndex] {
					at.BaseMissChance -= 0.01 * float64(newStacks-oldStacks)
				}
			}
		},
	}))
}

