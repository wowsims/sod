package scarlet_enclave

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func addScarlet60(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        243269, // TODO:
			Name:      "Generic",
			Level:     63,
			MobType:   proto.MobType_MobTypeUnknown,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      127_393, // TODO:
				stats.Armor:       3731,    // TODO:
				stats.AttackPower: 805,     // TODO:
				stats.BlockValue:  46,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2,      // TODO:
			MinBaseDamage:    6000,   // TODO:
			DamageSpread:     0.3333, // TODO:
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
		},
		AI: NewDefaultScarletEnclaveAI(),
	})
	core.AddPresetEncounter("Generic", []string{
		bossPrefix + "/Generic",
	})
}

type DefaultScarletEnclaveAI struct {
	ScarletEnclaveEncounter
	Target *core.Target
}

func NewDefaultScarletEnclaveAI() core.AIFactory {
	return func() core.TargetAI {
		return &DefaultScarletEnclaveAI{}
	}
}

func (ai *DefaultScarletEnclaveAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
	ai.registerScarletDominionAura(ai.Target)
}

func (ai *DefaultScarletEnclaveAI) Reset(*core.Simulation) {
}

func (ai *DefaultScarletEnclaveAI) ExecuteCustomRotation(sim *core.Simulation) {
	target := ai.Target.CurrentTarget

	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	}
}