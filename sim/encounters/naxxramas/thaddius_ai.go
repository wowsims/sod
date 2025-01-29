package naxxramas

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func addThaddius(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        15928,
			Name:      "Thaddius",
			Level:     63,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      39_520_129,
				stats.Armor:       3731,
				stats.AttackPower: 805,
				stats.BlockValue:  46,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2,
			MinBaseDamage:    6000,
			DamageSpread:     0.3333,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     []*proto.TargetInput{
				{
					Label:       "Stacks of Polarity Expected",
					Tooltip:     "How many stacks of polarity do you predict to have?  Max: ",
					InputType:   proto.InputType_Number,
					NumberValue: 0,
				},
			},
		},
		AI: NewThaddiusAI(),
	})
	core.AddPresetEncounter("Thaddius", []string{
		bossPrefix + "/Thaddius",
	})
}

type ThaddiusAI struct {
	Target         *core.Target
	ChainLightning *core.Spell
	Polarity       *core.Spell
	polarityStacks  float64
}

func NewThaddiusAI() core.AIFactory {
	return func() core.TargetAI {
		return &ThaddiusAI{}
	}
}

func (ai *ThaddiusAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
	ai.polarityStacks = config.TargetInputs[0].NumberValue

	ai.registerPolarity(ai.Target)
	ai.registerChainLightning(ai.Target)
}

const BossGCD = time.Millisecond * 1600

func (ai *ThaddiusAI) Reset(*core.Simulation) {
}

func (ai *ThaddiusAI) registerChainLightning(target *core.Target) {
	actionID := core.ActionID{SpellID: 28167}

	ai.ChainLightning = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		DefenseType:      core.DefenseTypeMagic,
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if sim.Proc(0.6, "Chain Lightning Target Chance") { // damage and target chance estimated from PTR
				baseDamage := sim.Roll(1850.0, 2250.0) 
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			}
		},
	})
}

func (ai *ThaddiusAI) registerPolarity(target *core.Target) {
	actionID := core.ActionID{SpellID: 28089}
	multiplierBonus := 1.0
	inverseMultiplierBonus := 1.0
	charactertarget := &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit

	polarityAura := charactertarget.RegisterAura(core.Aura{
		Label:    "Polarity Stacks",
		ActionID: core.ActionID{SpellID: 28059},
		Duration: time.Minute * 1,
		MaxStacks: 20,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			inverseMultiplierBonus = 1 / (1.0 + float64(oldStacks)*0.1)
			charactertarget.PseudoStats.DamageDealtMultiplier *= inverseMultiplierBonus
			multiplierBonus = 1.0 + float64(newStacks)*0.1
			charactertarget.PseudoStats.DamageDealtMultiplier *= multiplierBonus
		},
	})

	ai.Polarity = target.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		ProcMask: core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Millisecond * 4240, // Next server tick after cast complete
				CastTime: time.Millisecond * 3000,
			},
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 30,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, false)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			polarityAura.Deactivate(sim)
			polarityAura.Activate(sim)
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + time.Second*5,
				OnAction: func(sim *core.Simulation) {
					polarityAura.SetStacks(sim, int32(ai.polarityStacks)) // delay for stack activation
				},
			})
		},
	})
}

func (ai *ThaddiusAI) ExecuteCustomRotation(sim *core.Simulation) {
	target := ai.Target.CurrentTarget

	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	}

	if ai.Polarity.IsReady(sim) {
		ai.Polarity.Cast(sim, target)
		return
	}
	
	if ai.ChainLightning.IsReady(sim) {
		ai.Target.WaitUntil(sim, sim.CurrentTime+BossGCD)
		ai.ChainLightning.Cast(sim, target)
		return
	}
}
