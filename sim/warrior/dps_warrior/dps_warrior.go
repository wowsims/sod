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
	warOptions := options.GetWarrior()

	war := &DpsWarrior{
		Warrior: warrior.NewWarrior(character, options, warrior.WarriorInputs{
			StanceSnapshot: warOptions.Options.StanceSnapshot,
		}),
	}

	war.IsUsingRendStopAttack = warOptions.Options.IsUsingRendStopAttack
	war.IsUsingBloodthirstStopAttack = warOptions.Options.IsUsingBloodthirstStopAttack
	war.IsUsingQuickStrikeStopAttack = warOptions.Options.IsUsingQuickStrikeStopAttack
	war.IsUsingHamstringStopAttack = warOptions.Options.IsUsingHamstringStopAttack
	war.IsUsingWhirlwindStopAttack = warOptions.Options.IsUsingWhirlwindStopAttack
	war.IsUsingExecuteStopAttack = warOptions.Options.IsUsingExecuteStopAttack
	war.IsUsingOverpowerStopAttack = warOptions.Options.IsUsingOverpowerStopAttack
	war.IsUsingHeroicStrikeStopAttack = warOptions.Options.IsUsingHeroicStrikeStopAttack

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
