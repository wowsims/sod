package tank

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/druid"
)

func RegisterFeralTankDruid() {
	core.RegisterAgentFactory(
		proto.Player_FeralTankDruid{},
		proto.Spec_SpecFeralTankDruid,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFeralTankDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FeralTankDruid)
			if !ok {
				panic("Invalid spec value for Feral Tank Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFeralTankDruid(character *core.Character, options *proto.Player) *FeralTankDruid {
	tankOptions := options.GetFeralTankDruid()
	selfBuffs := druid.SelfBuffs{}

	bear := &FeralTankDruid{
		Druid:   druid.New(character, druid.Bear, selfBuffs, options.TalentsString),
		Options: tankOptions.Options,
	}

	bear.SelfBuffs.InnervateTarget = &proto.UnitReference{}
	if tankOptions.Options.InnervateTarget != nil {
		bear.SelfBuffs.InnervateTarget = tankOptions.Options.InnervateTarget
	}

	bear.EnableRageBar(core.RageBarOptions{
		StartingRage:          bear.Options.StartingRage,
		DamageDealtMultiplier: 1,
		DamageTakenMultiplier: 1,
	})

	bear.EnableAutoAttacks(bear, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand:       bear.GetBearWeapon(),
		AutoSwingMelee: true,
		ReplaceMHSwing: bear.TryMaul,
	})
	bear.ReplaceBearMHFunc = bear.TryMaul

	healingModel := options.HealingModel
	if healingModel != nil {
		if healingModel.InspirationUptime > 0.0 {
			core.ApplyInspiration(bear.GetCharacter(), healingModel.InspirationUptime)
		}
	}

	bear.PseudoStats.FeralCombatEnabled = true
	bear.PseudoStats.InFrontOfTarget = true

	return bear
}

type FeralTankDruid struct {
	*druid.Druid

	Options *proto.FeralTankDruid_Options
}

func (bear *FeralTankDruid) GetDruid() *druid.Druid {
	return bear.Druid
}

func (bear *FeralTankDruid) Initialize() {
	bear.Druid.Initialize()
	queuedRealismICD := &core.Cooldown{
		Timer:    bear.NewTimer(),
		Duration: core.SpellBatchWindow * 10,
	}
	bear.RegisterFeralTankSpells(queuedRealismICD)
}

func (bear *FeralTankDruid) Reset(sim *core.Simulation) {
	bear.Druid.Reset(sim)
	bear.Druid.CancelShapeshift(sim)
	bear.BearFormAura.Activate(sim)
	bear.Druid.PseudoStats.Stunned = false
}
