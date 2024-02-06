package elemental

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/shaman"
)

func RegisterElementalShaman() {
	core.RegisterAgentFactory(
		proto.Player_ElementalShaman{},
		proto.Spec_SpecElementalShaman,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewElementalShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ElementalShaman)
			if !ok {
				panic("Invalid spec value for Elemental Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewElementalShaman(character *core.Character, options *proto.Player) *ElementalShaman {
	eleOptions := options.GetElementalShaman()

	selfBuffs := shaman.SelfBuffs{
		Shield: eleOptions.Options.Shield,
	}

	totems := &proto.ShamanTotems{}
	if eleOptions.Options.Totems != nil {
		totems = eleOptions.Options.Totems
		totems.UseFireMcd = true // Control fire totems as MCD.
	}

	ele := &ElementalShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, totems, selfBuffs),
	}

	// if mh := ele.GetMHWeapon(); mh != nil {
	// 	ele.ApplyFlametongueImbueToItem(mh, false)
	// }

	// if oh := ele.GetOHWeapon(); oh != nil {
	// 	ele.ApplyFlametongueImbueToItem(oh, false)
	// }

	return ele
}

type ElementalShaman struct {
	*shaman.Shaman
}

func (eleShaman *ElementalShaman) GetShaman() *shaman.Shaman {
	return eleShaman.Shaman
}

func (eleShaman *ElementalShaman) Reset(sim *core.Simulation) {
	eleShaman.Shaman.Reset(sim)
}
