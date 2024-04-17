package retribution

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/paladin"
)

func RegisterRetributionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_RetributionPaladin{},
		proto.Spec_SpecRetributionPaladin,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewRetributionPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_RetributionPaladin) // I don't really understand this line
			if !ok {
				panic("Invalid spec value for Retribution Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewRetributionPaladin(character *core.Character, options *proto.Player) *RetributionPaladin {
	retOptions := options.GetRetributionPaladin()

	pal := paladin.NewPaladin(character, options.TalentsString)

	ret := &RetributionPaladin{
		Paladin:     pal,
		PrimarySeal: retOptions.Options.PrimarySeal,
	}

	ret.EnableAutoAttacks(ret, core.AutoAttackOptions{
		MainHand:       ret.WeaponFromMainHand(),
		AutoSwingMelee: true,
	})

	return ret
}

type RetributionPaladin struct {
	*paladin.Paladin

	PrimarySeal proto.PaladinSeal
}

func (ret *RetributionPaladin) GetPaladin() *paladin.Paladin {
	return ret.Paladin
}

func (ret *RetributionPaladin) Initialize() {
	ret.Paladin.Initialize()
}

func (ret *RetributionPaladin) Reset(sim *core.Simulation) {
	ret.Paladin.Reset(sim)
	ret.CurrentSeal = nil

	// Set the primary seal for APL actions.
	ret.PrimarySealSpell = ret.Paladin.GetMaxRankSeal(ret.PrimarySeal)
}
