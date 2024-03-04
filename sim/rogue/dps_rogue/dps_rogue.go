package dpsrogue

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/rogue"
)

// Register new DpsRogue
func RegisterDpsRogue() {
	core.RegisterAgentFactory(
		proto.Player_Rogue{},
		proto.Spec_SpecRogue,
		func(character *core.Character, spec *proto.Player) core.Agent {
			return NewDpsRogue(character, spec)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Rogue)
			if !ok {
				panic("Invalid spec value for Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

type DpsRogue struct {
	*rogue.Rogue
}

func NewDpsRogue(character *core.Character, options *proto.Player) *DpsRogue {
	rog := &DpsRogue{
		Rogue: rogue.NewRogue(character, options, options.GetRogue().Options),
	}

	return rog
}

func (rogue *DpsRogue) GetRogue() *rogue.Rogue {
	return rogue.Rogue
}

func (rogue *DpsRogue) Initialize() {
	rogue.Rogue.Initialize()

	// DPS related CD and talents here
}

func (rogue *DpsRogue) Reset(sim *core.Simulation) {
	rogue.Rogue.Reset(sim)
	// Aura resets here
}
