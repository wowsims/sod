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
	retOptions := options.GetRetributionPaladin().Options

	pal := paladin.NewPaladin(character, options, retOptions)

	ret := &RetributionPaladin{
		Paladin:     pal,
		primarySeal: retOptions.PrimarySeal,
		IsUsingDivineStormStopAttack: retOptions.IsUsingDivineStormStopAttack,
		IsUsingJudgementStopAttack: retOptions.IsUsingJudgementStopAttack,
		IsUsingCrusaderStrikeStopAttack: retOptions.IsUsingCrusaderStrikeStopAttack,
	}

	ret.EnableAutoAttacks(ret, core.AutoAttackOptions{
		MainHand:       ret.WeaponFromMainHand(),
		AutoSwingMelee: true,
	})

	return ret
}

type RetributionPaladin struct {
	*paladin.Paladin

	primarySeal proto.PaladinSeal
	IsUsingDivineStormStopAttack bool
	IsUsingJudgementStopAttack bool
	IsUsingCrusaderStrikeStopAttack bool
}

func (ret *RetributionPaladin) GetPaladin() *paladin.Paladin {
	return ret.Paladin
}

func (ret *RetributionPaladin) Initialize() {
	ret.Paladin.Initialize()
}

func (ret *RetributionPaladin) Reset(_ *core.Simulation) {
	ret.Paladin.ResetCurrentPaladinAura()
	ret.Paladin.ResetPrimarySeal(ret.primarySeal)
}
