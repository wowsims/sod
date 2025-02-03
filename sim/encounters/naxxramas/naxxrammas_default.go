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
	core.AddPresetEncounter("Level 60", []string{
		bossPrefix + "/Level 60",
	})
}


type DefaultNaxxAI struct {
	Target         *core.Target
	authorityFrozenWastesStacks int32
	authorityFrozenWastesAura *core.Aura
}

func NewDefaultNaxxAI() core.AIFactory {
	return func() core.TargetAI {
		return &DefaultNaxxAI{}
	}
}

func (ai *DefaultNaxxAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target

	ai.authorityFrozenWastesAura = ai.registerAuthorityOfTheFrozenWastesAura(ai.authorityFrozenWastesStacks)
}

func (ai *DefaultNaxxAI) registerAuthorityOfTheFrozenWastesAura(stacks int32) *core.Aura {
	charactertarget := &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
		
	return core.MakePermanent(charactertarget.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1218283},
		Label:     "Authority of the Frozen Wastes",
		MaxStacks: 4,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
			aura.SetStacks(sim, stacks)
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			aura.Unit.PseudoStats.DodgeReduction += 0.04 * float64(newStacks-oldStacks)

			for _, target := range sim.Encounter.TargetUnits {
				for _, at := range target.AttackTables[aura.Unit.UnitIndex] {
					at.BaseMissChance -= 0.01 * float64(newStacks-oldStacks)
				}
			}
		},
	}))
}

func (ai *DefaultNaxxAI) Reset(*core.Simulation) {
}

func (ai *DefaultNaxxAI) ExecuteCustomRotation(sim *core.Simulation) {
	target := ai.Target.CurrentTarget

	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	}

	if !ai.authorityFrozenWastesAura.IsActive() {
		ai.authorityFrozenWastesAura.Activate(sim)
	}
}
