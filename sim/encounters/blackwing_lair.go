package encounters

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func addVaelastraszTheCorrupt(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        13020, // Vanilla Vaelastrasz The Corrupt - no ID for SoD yet?
			Name:      "Blackwing Lair Vaelastrasz The Corrupt",
			Level:     63,
			MobType:   proto.MobType_MobTypeDragonkin,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      4_130_000, // Approx Vaelastrasz HP w/ Black Difficulty
				stats.Armor:       3731,      // TODO:
				stats.AttackPower: 805,       // TODO: Unknown attack power
				// TODO: Resistances
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2,     // TODO: Very slow attack interrupted by spells
			MinBaseDamage:    5000,  // TODO: Minimum unmitigated damage on reviewed log
			DamageSpread:     0.333, // TODO:
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs: []*proto.TargetInput{
				{
					Label:       "Time Burning Adrenaline Received",
					Tooltip:     "How long into the fight Burning Adrenaline is cast on the player. First cast is 20s (Select 0 to never receive)",
					InputType:   proto.InputType_Number,
					NumberValue: 0,
				},
			},
		},
		AI: NewVaelastraszTheCorruptAI(),
	})
	core.AddPresetEncounter("Blackwing Lair Vaelastrasz the Corrupt", []string{
		bossPrefix + "/Blackwing Lair Vaelastrasz the Corrupt",
	})
}

type VaelastraszTheCorruptAI struct {
	Target                     *core.Target
	essenceOfTheRedSpell       *core.Spell
	burningAdrenalineSpell     *core.Spell
	burningAdrenalineTankSpell *core.Spell
	burningAdrenalineTime      float64
	fireNovaSpell              *core.Spell
	flameBreathSpell           *core.Spell
	cleaveSpell                *core.Spell
	canAct                     bool
}

func NewVaelastraszTheCorruptAI() core.AIFactory {
	return func() core.TargetAI {
		return &VaelastraszTheCorruptAI{}
	}
}

func (ai *VaelastraszTheCorruptAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
	ai.burningAdrenalineTime = config.TargetInputs[0].NumberValue
	ai.registerSpells()
	ai.canAct = true
}

func (ai *VaelastraszTheCorruptAI) registerSpells() {
	essenceOfTheRedActionID := core.ActionID{SpellID: 23513}
	burningAdrenalineActionID := core.ActionID{SpellID: 367987}
	burningAdrenalineTankActionID := core.ActionID{SpellID: 469261}
	fireNovaActionID := core.ActionID{SpellID: 23462}
	flameBreathActionID := core.ActionID{SpellID: 23461}
	cleaveActionID := core.ActionID{SpellID: 19983}

	target := &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit

	essenceOfTheRedManaMetrics := target.NewManaMetrics(essenceOfTheRedActionID)
	essenceOfTheRedEnergyMetrics := target.NewEnergyMetrics(essenceOfTheRedActionID)
	essenceOfTheRedRageMetrics := target.NewRageMetrics(essenceOfTheRedActionID)

	ai.essenceOfTheRedSpell = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID: essenceOfTheRedActionID,
		ProcMask: core.ProcMaskEmpty,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Minute * 4,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Essemce of the Red",
			},
			NumberOfTicks: 240,
			TickLength:    time.Second * 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if target.HasManaBar() {
					target.AddMana(sim, 500, essenceOfTheRedManaMetrics)
				}
				if target.HasEnergyBar() {
					target.AddEnergy(sim, 50, essenceOfTheRedEnergyMetrics)
				}
				if target.HasRageBar() {
					target.AddRage(sim, 20, essenceOfTheRedRageMetrics)
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	multiplierBonus := 1.0
	inverseMultiplierBonus := 1.0

	ai.burningAdrenalineSpell = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID: burningAdrenalineActionID,
		ProcMask: core.ProcMaskEmpty,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Minute * 4,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Burning Adrenaline",
				MaxStacks: 100,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
					inverseMultiplierBonus = 1 / (1.0 + float64(oldStacks)*0.1)
					target.MultiplyMeleeSpeed(sim, inverseMultiplierBonus)
					target.MultiplyCastSpeed(inverseMultiplierBonus)
					target.PseudoStats.DamageDealtMultiplier *= inverseMultiplierBonus
					multiplierBonus = 1.0 + float64(newStacks)*0.1
					target.MultiplyMeleeSpeed(sim, multiplierBonus)
					target.MultiplyCastSpeed(multiplierBonus)
					target.PseudoStats.DamageDealtMultiplier *= multiplierBonus
				},
			},
			NumberOfTicks: 240,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.AddStack(sim)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
			spell.Dot(target).AddStack(sim)
		},
	})

	ai.burningAdrenalineTankSpell = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID:         burningAdrenalineTankActionID,
		SpellSchool:      core.SpellSchoolFire,
		DefenseType:      core.DefenseTypeMagic,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagIgnoreResists,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Second * 70,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Burning Adrenaline (Tank)",
				MaxStacks: 16,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
					inverseMultiplierBonus = 1 / (1.0 + float64(oldStacks)*0.1)
					target.MultiplyMeleeSpeed(sim, inverseMultiplierBonus)
					target.MultiplyCastSpeed(inverseMultiplierBonus)
					target.PseudoStats.DamageDealtMultiplier *= inverseMultiplierBonus
					multiplierBonus = 1.0 + float64(newStacks)*0.1
					target.MultiplyMeleeSpeed(sim, multiplierBonus)
					target.MultiplyCastSpeed(multiplierBonus)
					target.PseudoStats.DamageDealtMultiplier *= multiplierBonus
					if newStacks == 16 {
						target.RemoveHealth(sim, target.CurrentHealth())
						if sim.Log != nil {
							target.Log(sim, "Burning Adrenaline (Tank) Death")
						}
						sim.Cleanup()
					}
				},
			},
			NumberOfTicks: 16,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				burningAdrenalineTickDamage := target.MaxHealth() * 0.01 * float64(dot.GetStacks())
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, burningAdrenalineTickDamage, dot.OutcomeTick)
				if !(dot.GetStacks() == 0) {
					dot.AddStack(sim)
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
			spell.Dot(target).AddStack(sim)
		},
	})

	ai.cleaveSpell = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID:         cleaveActionID,
		SpellSchool:      core.SpellSchoolPhysical,
		DefenseType:      core.DefenseTypeMelee,
		ProcMask:         core.ProcMaskEmpty,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Second * 8,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(5000, 8000)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		},
	})

	ai.fireNovaSpell = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID:         fireNovaActionID,
		SpellSchool:      core.SpellSchoolFire,
		DefenseType:      core.DefenseTypeMagic,
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Second * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(555.0, 645.0)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
		},
	})

	flameBreathTickDamage := 0.0

	ai.flameBreathSpell = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID:         flameBreathActionID,
		SpellSchool:      core.SpellSchoolFire,
		DefenseType:      core.DefenseTypeMagic,
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Millisecond * 3240, // Next server tick after cast complete
				CastTime: time.Millisecond * 2000,
			},
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Second * 9,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, false)
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Flame Breath",
				MaxStacks: 100,
				ActionID:  flameBreathActionID,
				Duration:  time.Second * 15,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					flameBreathTickDamage = 0.0
				},
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, flameBreathTickDamage, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(3500, 4500)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHit)

			if result.Landed() {
				flameBreathTickDamage += sim.Roll(938, 1062)
				spell.Dot(target).Activate(sim)
				spell.Dot(target).AddStack(sim)
			}
			spell.DealDamage(sim, result)
		},
	})

}

func (ai *VaelastraszTheCorruptAI) Reset(*core.Simulation) {
}

const BossGCD = time.Millisecond * 1600

func (ai *VaelastraszTheCorruptAI) ExecuteCustomRotation(sim *core.Simulation) {
	if !ai.canAct {
		ai.Target.WaitUntil(sim, sim.CurrentTime+BossGCD)
		return
	}

	target := ai.Target.CurrentTarget
	isTank := true
	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
		isTank = false
	}

	if ai.essenceOfTheRedSpell.CanCast(sim, target) {
		ai.essenceOfTheRedSpell.Cast(sim, target)
		ai.Target.WaitUntil(sim, sim.CurrentTime+BossGCD)
		return
	}

	if isTank {
		if ai.burningAdrenalineTankSpell.CanCast(sim, target) && (sim.CurrentTime >= (time.Second * 10)) {
			ai.burningAdrenalineTankSpell.Cast(sim, target)
		}
		if ai.fireNovaSpell.CanCast(sim, target) {
			ai.fireNovaSpell.Cast(sim, target)
		}
		if ai.flameBreathSpell.CanCast(sim, target) {
			ai.flameBreathSpell.Cast(sim, target)
			return
		}
		if ai.cleaveSpell.CanCast(sim, target) {
			ai.cleaveSpell.Cast(sim, target)
			return
		}
	} else {
		//Do not cast if player input is 0 for BA Time
		if ai.burningAdrenalineTime == 0 {
			return
		} else if ai.burningAdrenalineSpell.CanCast(sim, target) && (sim.CurrentTime >= (time.Second * time.Duration(ai.burningAdrenalineTime))) {
			ai.burningAdrenalineSpell.Cast(sim, target)
		}
	}
}
