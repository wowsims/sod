package encounters

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func addSunkenTempleDragonkin(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        218571, // SoD Shade of Eranikus
			Name:      "Sunken Temple Dragonkin Boss",
			Level:     52,
			MobType:   proto.MobType_MobTypeDragonkin,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      1_450_000, // Approx Shdae of Eranikus health
				stats.Armor:       3700,      // TODO:
				stats.AttackPower: 574,       // TODO:
				// TODO: Resistances
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2,      // TODO:
			MinBaseDamage:    3000,   // TODO:
			DamageSpread:     0.3333, // TODO:
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
		AI: NewSunkenTempleDragonkinAI(),
	})
	core.AddPresetEncounter("Sunken Temple Dragonkin Boss", []string{
		bossPrefix + "/Sunken Temple Dragonkin Boss",
	})
}

type SunkenTempleDragonkinAI struct {
	Target *core.Target
}

func NewSunkenTempleDragonkinAI() core.AIFactory {
	return func() core.TargetAI {
		return &SunkenTempleDragonkinAI{}
	}
}

func (ai *SunkenTempleDragonkinAI) Initialize(target *core.Target, _ *proto.Target) {
	ai.Target = target
}

func (ai *SunkenTempleDragonkinAI) Reset(*core.Simulation) {
}

func (ai *SunkenTempleDragonkinAI) ExecuteCustomRotation(_ *core.Simulation) {
}
