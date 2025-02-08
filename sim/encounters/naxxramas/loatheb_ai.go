package naxxramas

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func addLoatheb(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        16011,
			Name:      "Loatheb",
			Level:     63,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      26_286_324,
				stats.Armor:       3731,
				stats.AttackPower: 805,
				stats.BlockValue:  46,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       1.2,
			MinBaseDamage:    6229,
			DamageSpread:     0.3333,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs: []*proto.TargetInput{
				NaxxramasDifficultyLevels,
				{
					Label:       "Spore Assignment (1-10)",
					Tooltip:     "Which spore are you assigned to?",
					InputType:   proto.InputType_Number,
					NumberValue: 0,
				},
			},
		},
		AI: NewLoathebAI(),
	})
	core.AddPresetEncounter("Loatheb", []string{
		bossPrefix + "/Loatheb",
	})
}

type LoathebAI struct {
	NaxxramasEncounter
	Target          *core.Target
	SummonSpore     *core.Spell
	sporeAssignment float64
}

func NewLoathebAI() core.AIFactory {
	return func() core.TargetAI {
		return &LoathebAI{}
	}
}

func (ai *LoathebAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
	ai.authorityFrozenWastesStacks = config.TargetInputs[0].EnumValue
	ai.sporeAssignment = config.TargetInputs[1].NumberValue
	ai.registerSummonSpore(target)

	ai.registerAuthorityOfTheFrozenWastesAura(ai.Target, ai.authorityFrozenWastesStacks)
}

func (ai *LoathebAI) registerSummonSpore(target *core.Target) {
	actionID := core.ActionID{SpellID: 29234}
	charactertarget := &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit

	fungalBloomAura := charactertarget.RegisterAura(core.Aura{
		Label:    "Fungal Bloom",
		ActionID: core.ActionID{SpellID: 29232},
		Duration: time.Second * 90,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.MeleeCrit, 50*core.SpellCritRatingPerCritChance)
			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, 60*core.SpellCritRatingPerCritChance)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.MeleeCrit, -50*core.SpellCritRatingPerCritChance)
			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, -60*core.SpellCritRatingPerCritChance)
		},
	})

	ai.SummonSpore = target.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		ProcMask: core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Millisecond * 4240, // Next server tick after cast complete
			},
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 13,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if !fungalBloomAura.IsActive() {
				fungalBloomAura.Activate(sim)
			}
		},
	})
}

func (ai *LoathebAI) Reset(*core.Simulation) {
}

func (ai *LoathebAI) ExecuteCustomRotation(sim *core.Simulation) {
	target := ai.Target.CurrentTarget

	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	}

	if ai.SummonSpore.IsReady(sim) {
		if sim.CurrentTime > ((time.Duration(ai.sporeAssignment)*13)+4)*time.Second && ai.sporeAssignment != 0 {
			ai.SummonSpore.Cast(sim, target)
			return
		}
	}
}
