package protection

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/paladin"
)

func RegisterProtectionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionPaladin{},
		proto.Spec_SpecProtectionPaladin,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewProtectionPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ProtectionPaladin) // I don't really understand this line
			if !ok {
				panic("Invalid spec value for Protection Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewProtectionPaladin(character *core.Character, options *proto.Player) *ProtectionPaladin {
	protOptions := options.GetProtectionPaladin().Options

	pal := paladin.NewPaladin(character, options, protOptions)

	prot := &ProtectionPaladin{
		Paladin:                         pal,
		primarySeal:                     protOptions.PrimarySeal,
		IsUsingDivineStormStopAttack:    protOptions.IsUsingDivineStormStopAttack,
		IsUsingJudgementStopAttack:      protOptions.IsUsingJudgementStopAttack,
		IsUsingCrusaderStrikeStopAttack: protOptions.IsUsingCrusaderStrikeStopAttack,
	}

	prot.EnableAutoAttacks(prot, core.AutoAttackOptions{
		MainHand:       prot.WeaponFromMainHand(),
		AutoSwingMelee: true,
	})

	return prot
}

type ProtectionPaladin struct {
	*paladin.Paladin

	primarySeal                     proto.PaladinSeal
	IsUsingDivineStormStopAttack    bool
	IsUsingJudgementStopAttack      bool
	IsUsingCrusaderStrikeStopAttack bool
}

func (prot *ProtectionPaladin) GetPaladin() *paladin.Paladin {
	return prot.Paladin
}

func (prot *ProtectionPaladin) Initialize() {
	prot.Paladin.Initialize()
}

func (prot *ProtectionPaladin) Reset(_ *core.Simulation) {
	prot.Paladin.ResetCurrentPaladinAura()
	prot.Paladin.ResetPrimarySeal(prot.primarySeal)
}
