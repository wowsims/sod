package warden

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/shaman"
)

func RegisterWardenShaman() {
	core.RegisterAgentFactory(
		proto.Player_WardenShaman{},
		proto.Spec_SpecWardenShaman,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewWardenShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_WardenShaman)
			if !ok {
				panic("Invalid spec value for Warden Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

type WardenShaman struct {
	*shaman.Shaman

	Options *proto.WardenShaman_Options
}

func NewWardenShaman(character *core.Character, options *proto.Player) *WardenShaman {
	warden := &WardenShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString),
	}

	// Enable Auto Attacks for this spec
	warden.EnableAutoAttacks(warden, core.AutoAttackOptions{
		MainHand:       warden.WeaponFromMainHand(),
		OffHand:        warden.WeaponFromOffHand(),
		AutoSwingMelee: true,
	})

	return warden
}

func (warden *WardenShaman) GetShaman() *shaman.Shaman {
	return warden.Shaman
}

func (warden *WardenShaman) Initialize() {
	warden.Shaman.Initialize()
}

func (warden *WardenShaman) Reset(sim *core.Simulation) {
	warden.Shaman.Reset(sim)
	warden.Shaman.PseudoStats.Stunned = false
}
