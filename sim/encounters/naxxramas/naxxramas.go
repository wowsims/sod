package naxxramas

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func Register() {
	addNaxx60("SoD/Naxxramas")
	addPatchwerk("SoD/Naxxramas")
	addThaddius("SoD/Naxxramas")
	addLoatheb("SoD/Naxxramas")
}

type NaxxramasEncounter struct {
	authorityFrozenWastesStacks int32
}

var NaxxramasDifficultyLevels = &proto.TargetInput{
	Label:       "Difficulty Level",
	Tooltip:     "Affects the Authority of the Frozen Wastes debuff for tanks.",
	InputType:   proto.InputType_Enum,
	EnumValue:   0,
	EnumOptions: []string{"Normal", "Hardmode 1/4", "Hardmode 2/4", "Hardmode 3/4", "Hardmode 4/4"},
}

func (naxxEncounter *NaxxramasEncounter) registerAuthorityOfTheFrozenWastesAura(target *core.Target, stacks int32) {
	charactertarget := target.Env.Raid.Parties[0].Players[0].GetCharacter()

	charactertarget.GetOrRegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1218283},
		Label:     "Authority of the Frozen Wastes",
		MaxStacks: 4,
		Duration:  core.NeverExpires,
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
	})
}
