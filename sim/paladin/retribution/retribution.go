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

	// ret.PaladinAura = retOptions.Options.Aura

	ret.EnableAutoAttacks(ret, core.AutoAttackOptions{
		MainHand:       ret.WeaponFromMainHand(0), // Set crit multiplier later when we have targets.
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
	ret.PrimarySealSpell = ret.Paladin.SealOfRighteousness[1]

	switch ret.PrimarySeal {
	case proto.PaladinSeal_RighteousnessRank1:
		ret.PrimarySealSpell = ret.Paladin.SealOfRighteousness[1]
	case proto.PaladinSeal_RighteousnessRank2:
		ret.PrimarySealSpell = ret.Paladin.SealOfRighteousness[2]
	}
}

func (ret *RetributionPaladin) Reset(sim *core.Simulation) {
	ret.Paladin.Reset(sim)
	ret.CurrentSeal = nil

	// case proto.PaladinSeal_Vengeance:
	// 	ret.CurrentSeal = ret.SealOfVengeanceAura
	// 	ret.SealOfVengeanceAura.Activate(sim)
	// case proto.PaladinSeal_Command:
	// 	ret.CurrentSeal = ret.SealOfCommandAura
	// 	ret.SealOfCommandAura.Activate(sim)
	// case proto.PaladinSeal_Righteousness:
	// 	ret.CurrentSeal = ret.SealOfRighteousnessAura
	// 	ret.SealOfRighteousnessAura.Activate(sim)
}
