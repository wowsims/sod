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
			Id:        15952 , // TODO:
			Name:      "Generic",
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
			TargetInputs:     []*proto.TargetInput{
				{
					Label:       "Authority of The Frozen Wastes Stacks",
					Tooltip:     "Hard Modes Activated?",
					InputType:   proto.InputType_Enum,
					EnumValue:   0,
					EnumOptions: []string{
						"0", "1", "2", "3", "4",
					},
				},
			},
		},
		AI: NewDefaultNaxxAI(),
	})
	core.AddPresetEncounter("Generic", []string{
		bossPrefix + "/Generic",
	})
}


type DefaultNaxxAI struct {
	NaxxramasEncounter
	Target         *core.Target
}

func NewDefaultNaxxAI() core.AIFactory {
	return func() core.TargetAI {
		return &DefaultNaxxAI{}
	}
}

func (ai *DefaultNaxxAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
	ai.authorityFrozenWastesStacks = config.TargetInputs[0].EnumValue

	ai.registerAuthorityOfTheFrozenWastesAura(ai.Target, ai.authorityFrozenWastesStacks)
}


func (ai *DefaultNaxxAI) Reset(*core.Simulation) {
}

func (ai *DefaultNaxxAI) ExecuteCustomRotation(sim *core.Simulation) {
	target := ai.Target.CurrentTarget

	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	}
}
