package dps_hunter

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/hunter"
)

func RegisterDPSHunter() {
	core.RegisterAgentFactory(
		proto.Player_Hunter{},
		proto.Spec_SpecHunter,
		func(character *core.Character, options *proto.Player) core.Agent {
			return hunter.NewHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Hunter)
			if !ok {
				panic("Invalid spec value for Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewDPSHunter(character *core.Character, options *proto.Player) *DPSHunter {
	hunterOptions := options.GetHunter()

	hunter := &DPSHunter{
		Hunter:  hunter.NewHunter(character, options),
		Options: hunterOptions.Options,
	}

	return hunter
}

type DPSHunter struct {
	*hunter.Hunter

	Options *proto.Hunter_Options
}
