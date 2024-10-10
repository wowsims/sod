package tankwarrior

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/warrior"
)

func RegisterTankWarrior() {
	core.RegisterAgentFactory(
		proto.Player_TankWarrior{},
		proto.Spec_SpecTankWarrior,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewTankWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_TankWarrior)
			if !ok {
				panic("Invalid spec value for Tank Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type TankWarrior struct {
	*warrior.Warrior
}

func NewTankWarrior(character *core.Character, options *proto.Player) *TankWarrior {
	warOptions := options.GetTankWarrior().Options

	war := warrior.NewWarrior(character, options, warOptions, warrior.WarriorInputs{
		StanceSnapshot: warOptions.StanceSnapshot,
	})

	tankWar := &TankWarrior{
		Warrior: war,
	}

	tankWar.EnableRageBar(core.RageBarOptions{
		StartingRage:          warOptions.StartingRage,
		DamageDealtMultiplier: 1,
		DamageTakenMultiplier: 1,
	})

	tankWar.EnableAutoAttacks(tankWar, core.AutoAttackOptions{
		MainHand:       tankWar.WeaponFromMainHand(),
		OffHand:        tankWar.WeaponFromOffHand(),
		AutoSwingMelee: true,
		ReplaceMHSwing: tankWar.TryHSOrCleave,
	})

	return tankWar
}

func (war *TankWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *TankWarrior) Initialize() {
	war.Warrior.Initialize()

	war.RegisterShieldWallCD()
	war.RegisterShieldBlockCD()
	war.DefensiveStanceAura.BuildPhase = core.CharacterBuildPhaseTalents
}

func (war *TankWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
	war.Warrior.PseudoStats.Stunned = false
}
