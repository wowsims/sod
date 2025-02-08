package naxxramas

<<<<<<< HEAD:sim/encounters/naxxramas/naxxramas.go
import (
	"github.com/wowsims/sod/sim/core"
)

func Register() {
	addNaxx60("SoD/Naxxramas")
	addPatchwerk("SoD/Naxxramas")
	addThaddius("SoD/Naxxramas")
	addLoatheb("SoD/Naxxramas")
}

type NaxxramasEncounter struct {
	authorityFrozenWastesStacks int32
	authorityFrozenWastesAura   *core.Aura
}

func (naxxEncounter *NaxxramasEncounter) registerAuthorityOfTheFrozenWastesAura(target *core.Target, stacks int32) {
	charactertarget := target.Env.Raid.Parties[0].Players[0].GetCharacter()

	core.MakePermanent(charactertarget.GetOrRegisterAura(core.Aura{
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
=======
func Register() {
	addPatchwerk25("Naxxrammas 25")
	addKelThuzad25("Naxxrammas 25")
	addThaddius25("Naxxrammas 25")
	addLoatheb25("Naxxrammas 25")

	// TODO: Figure out why this isn't pickable
	//addPatchwerk10("Naxxrammas")
>>>>>>> master:sim/encounters/naxxramas/naxxrammas.go
}
