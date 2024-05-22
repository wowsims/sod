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
	_ = options.GetElementalShaman()

	ele := &ElementalShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString),
	}

	// Enable Auto Attacks for this spec
	ele.EnableAutoAttacks(ele, core.AutoAttackOptions{
		MainHand:       ele.WeaponFromMainHand(),
		OffHand:        ele.WeaponFromOffHand(),
		AutoSwingMelee: true,
	})

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
