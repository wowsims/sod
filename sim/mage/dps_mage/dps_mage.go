package dps_mage

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/mage"
)

func RegisterDPSMage() {
	core.RegisterAgentFactory(
		proto.Player_Mage{},
		proto.Spec_SpecMage,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDPSMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Mage)
			if !ok {
				panic("Invalid spec value for Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewDPSMage(character *core.Character, options *proto.Player) *DPSMage {
	mageOptions := options.GetMage()

	mage := &DPSMage{
		Mage:    mage.NewMage(character, options),
		Options: mageOptions.Options,
	}

	return mage
}

type DPSMage struct {
	*mage.Mage

	Options *proto.Mage_Options
}
