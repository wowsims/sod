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

	Options *proto.TankWarrior_Options
}

func NewTankWarrior(character *core.Character, options *proto.Player) *TankWarrior {
	warOptions := options.GetTankWarrior()

	war := &TankWarrior{
		Warrior: warrior.NewWarrior(character, options.TalentsString, warrior.WarriorInputs{
			QueueDelay:     warOptions.Options.QueueDelay,
			Stance:         warOptions.Options.Stance,
			StanceSnapshot: warOptions.Options.StanceSnapshot,
		}),
		Options: warOptions.Options,
	}

	war.EnableRageBar(core.RageBarOptions{
		StartingRage:          warOptions.Options.StartingRage,
		DamageDealtMultiplier: 1,
		DamageTakenMultiplier: 1,
	})

	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(),
		OffHand:        war.WeaponFromOffHand(),
		AutoSwingMelee: true,
		ReplaceMHSwing: war.TryHSOrCleave,
	})

	return war
}

func (war *TankWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *TankWarrior) Initialize() {
	war.Warrior.Initialize()
	war.DefensiveStanceAura.BuildPhase = core.CharacterBuildPhaseTalents
}

func (war *TankWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
	war.Warrior.PseudoStats.Stunned = false
}
