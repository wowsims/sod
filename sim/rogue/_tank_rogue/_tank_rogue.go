package tankrogue

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/rogue"
)

func RegisterTankRogue() {
	core.RegisterAgentFactory(
		proto.Player_Rogue{},
		proto.Spec_SpecTankRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewTankRogue(character, options)
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

type TankRogue struct {
	*rogue.Rogue
}

func NewTankRogue(character *core.Character, options *proto.Player) *TankRogue {
	tank := &TankRogue{
		Rogue: rogue.NewRogue(character, options, options.GetRogue().Options),
	}

	return tank
}

func (tank *TankRogue) GetRogue() *rogue.Rogue {
	return tank.Rogue
}

func (tank *TankRogue) Initialize() {
	tank.Rogue.Initialize()

	// Initialize tank related CDs/auras here
	// Evasion
}

func (tank *TankRogue) Reset(sim *core.Simulation) {
	tank.Rogue.Reset(sim)
	// Aura/CD reset here
}
