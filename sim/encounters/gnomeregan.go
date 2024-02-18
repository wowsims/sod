package encounters

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func addGnomereganMechanical(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        218537, // TODO:
			Name:      "Level 40 Mechanical",
			Level:     42,
			MobType:   proto.MobType_MobTypeMechanical,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      279_345, // Electrocutioner 6000 health
				stats.Armor:       4000,    // Approx average armor of Gnomeregan bosses
				stats.AttackPower: 574,     // TODO:
				// TODO: Resistances
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2,      // TODO:
			MinBaseDamage:    1000,   // TODO:
			DamageSpread:     0.3333, // TODO:
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
		AI: NewGnomereganMechanicalAI(),
	})
	core.AddPresetEncounter("Level 40 Mechanical", []string{
		bossPrefix + "/Level 40 Mechanical",
	})
}

type GnomereganMechanicalAI struct {
	Target *core.Target
}

func NewGnomereganMechanicalAI() core.AIFactory {
	return func() core.TargetAI {
		return &GnomereganMechanicalAI{}
	}
}

func (ai *GnomereganMechanicalAI) Initialize(target *core.Target, _ *proto.Target) {
	target.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier = .8
	target.Unit.PseudoStats.PoisonDamageTakenMultiplier = .8

	ai.Target = target
}

func (ai *GnomereganMechanicalAI) Reset(*core.Simulation) {
}

func (ai *GnomereganMechanicalAI) ExecuteCustomRotation(sim *core.Simulation) {
}
