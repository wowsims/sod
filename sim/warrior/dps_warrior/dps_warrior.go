package dpswarrior

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/warrior"
)

func RegisterDpsWarrior() {
	core.RegisterAgentFactory(
		proto.Player_Warrior{},
		proto.Spec_SpecWarrior,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDpsWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Warrior)
			if !ok {
				panic("Invalid spec value for Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type DpsWarrior struct {
	*warrior.Warrior
}

func NewDpsWarrior(character *core.Character, options *proto.Player) *DpsWarrior {
	warOptions := options.GetWarrior().Options

	war := warrior.NewWarrior(character, options, warOptions, warrior.WarriorInputs{
		StanceSnapshot: warOptions.StanceSnapshot,
	})

	dpsWar := &DpsWarrior{
		Warrior: war,
	}

	dpsWar.EnableRageBar(core.RageBarOptions{
		StartingRage:          warOptions.StartingRage,
		DamageDealtMultiplier: 1,
		DamageTakenMultiplier: 1,
	})

	dpsWar.EnableAutoAttacks(dpsWar, core.AutoAttackOptions{
		MainHand:       dpsWar.WeaponFromMainHand(),
		OffHand:        dpsWar.WeaponFromOffHand(),
		AutoSwingMelee: true,
		ReplaceMHSwing: dpsWar.TryHSOrCleave,
	})

	return dpsWar
}

func (war *DpsWarrior) OnGCDReady(sim *core.Simulation) {
}

func (war *DpsWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *DpsWarrior) Initialize() {
	war.Warrior.Initialize()
}

func (war *DpsWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
}
