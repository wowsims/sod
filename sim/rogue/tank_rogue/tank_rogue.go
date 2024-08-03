package tankrogue

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
	"github.com/wowsims/sod/sim/rogue"
)

func RegisterTankRogue() {
	core.RegisterAgentFactory(
		proto.Player_TankRogue{},
		proto.Spec_SpecTankRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewTankRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_TankRogue)
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
		Rogue: rogue.NewRogue(character, options, options.GetTankRogue().Options),
	}

	tank.AddStatDependency(stats.Agility, stats.Dodge, core.DodgePerAgiAtLevel[character.Class][int(tank.Level)])
	tank.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	tank.PseudoStats.InFrontOfTarget = true

	return tank
}

func (tank *TankRogue) GetRogue() *rogue.Rogue {
	return tank.Rogue
}

func (tank *TankRogue) Initialize() {
	tank.Rogue.Initialize()

	// Initialize tank related CDs/auras here
	tank.Rogue.RegisterEvasionSpell()
}

func (tank *TankRogue) Reset(sim *core.Simulation) {
	tank.Rogue.Reset(sim)
	// Aura/CD reset here
}
