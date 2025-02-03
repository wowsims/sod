package naxxramas

import (
	"time"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func addPatchwerk(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        16028,
			Name:      "Patchwerk",
			Level:     63,
			MobType:   proto.MobType_MobTypeUndead,
			TankIndex:       0,

			Stats: stats.Stats{
				stats.Health:      16_950_147,
				stats.Armor:       3731,
				stats.AttackPower: 805,
				stats.BlockValue:  46,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       0.75,
			MinBaseDamage:    7166,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			DamageSpread:     0.1,
			TargetInputs: []*proto.TargetInput{
				{
					Label:       "Hateful Tank",
					Tooltip:     "Click to turn off auto attacks and activate hateful strikes.",
					InputType:   proto.InputType_Bool,
				},
				{
					Label:       "Percent of Hatefuls Taken",
					Tooltip:     "What `%` (0-100) of hateful strikes are you targeted by?",
					InputType:   proto.InputType_Number,
					NumberValue: 70.0,
				},
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
		AI: NewPatchwerkAI(),
	})
	core.AddPresetEncounter("Patchwerk", []string{
		bossPrefix + "/Patchwerk",
	})
}

type PatchwerkAI struct {
	// Unit references
	Target   *core.Target
	isHatefulTank  bool
	hatefulPercent float64
	authorityFrozenWastesStacks int32
	authorityFrozenWastesAura *core.Aura

	// Spells
	HatefulStrikePrimer *core.Spell
	HatefulStrike       *core.Spell
	Frenzy              *core.Spell
}

func NewPatchwerkAI() core.AIFactory {
	return func() core.TargetAI {
		return &PatchwerkAI{}
	}
}

func (ai *PatchwerkAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
	ai.isHatefulTank = config.TargetInputs[0].BoolValue
	ai.hatefulPercent = config.TargetInputs[1].NumberValue / 100.0

	ai.registerHatefulStrikePrimerSpell(ai.Target)
	ai.registerHatefulStrikeSpell(ai.Target)
	ai.registerFrenzySpell(ai.Target)
	ai.authorityFrozenWastesAura = ai.registerAuthorityOfTheFrozenWastesAura(ai.authorityFrozenWastesStacks)
}

func (ai *PatchwerkAI) registerAuthorityOfTheFrozenWastesAura(stacks int32) *core.Aura {
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

func (ai *PatchwerkAI) Reset(*core.Simulation) {
}

func (ai *PatchwerkAI) registerHatefulStrikePrimerSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 28307}

	ai.HatefulStrikePrimer = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskDirect,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Millisecond * 1200,
			},
		},

		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Spell does nothing but does proc things such as lightning shield.
		},
	})
}

func (ai *PatchwerkAI) registerHatefulStrikeSpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 28308}

	ai.HatefulStrike = target.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Millisecond * 1200,
			},
		},

		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO cannot crit
			baseDamage := sim.Roll(22000, 30000)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialNoCrit)
		},
	})
}

func (ai *PatchwerkAI) registerFrenzySpell(target *core.Target) {
	actionID := core.ActionID{SpellID: 28131}
	frenzyAura := target.GetOrRegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Frenzy",
		Duration: 5 * time.Minute,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.25
			aura.Unit.MultiplyMeleeSpeed(sim, 1.4)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.25
			aura.Unit.MultiplyMeleeSpeed(sim, 1.0/1.4)
		},
	})

	ai.Frenzy = target.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			frenzyAura.Activate(sim)
		},
	})
}

func (ai *PatchwerkAI) ExecuteCustomRotation(sim *core.Simulation) {
	target := ai.Target.CurrentTarget

	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	}

	if !ai.authorityFrozenWastesAura.IsActive() {
		ai.authorityFrozenWastesAura.Activate(sim)
	}

	if ai.HatefulStrikePrimer.IsReady(sim) {
		ai.HatefulStrikePrimer.Cast(sim, target)
	}
	
	if ai.Frenzy.IsReady(sim) && sim.GetRemainingDurationPercent() < 0.05 {
		ai.Frenzy.Cast(sim, target)
	}

	if ai.isHatefulTank {
		ai.Target.AutoAttacks.CancelAutoSwing(sim)
		 
		if ai.HatefulStrike.IsReady(sim) {
			if sim.Proc(ai.hatefulPercent, "Hateful Strike Target Chance") {
				ai.HatefulStrike.Cast(sim, target)
				return
			}
			ai.Target.WaitUntil(sim, sim.CurrentTime + 1200 * time.Millisecond)
		}
	}
}
