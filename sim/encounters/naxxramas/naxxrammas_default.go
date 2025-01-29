package naxxramas

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func addNaxx60(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        213336, // TODO:
			Name:      "Level 60",
			Level:     63,
			MobType:   proto.MobType_MobTypeUndead,
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
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
	})
	core.AddPresetEncounter("Level 60", []string{
		bossPrefix + "/Level 60",
	})
}
