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
			SwingSpeed:       4,      // TODO: Very slow attack interrupted by spells
			MinBaseDamage:    5200,   // TODO: Minimum unmitigated damage on reviewed log
			DamageSpread:     0.333, // TODO:
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
		AI: NewVaelastraszTheCorruptAI(),
	})
	core.AddPresetEncounter("Blackwing Lair Vaelastrasz The Corrupt", []string{
		bossPrefix + "/Blackwing Lair Vaelastrasz The Corrupt",
	})
}

type VaelastraszTheCorruptAI struct {
	Target *core.Target
	essenceOfTheRedAura *core.Aura
	essenceOfTheRedSpell *core.Spell
}

func NewVaelastraszTheCorruptAI() core.AIFactory {
	return func() core.TargetAI {
		return &VaelastraszTheCorruptAI{}
	}
}

func (ai *VaelastraszTheCorruptAI) Initialize(target *core.Target, _ *proto.Target) {
	ai.Target = target
	ai.registerSpells()
}

func (ai *VaelastraszTheCorruptAI) registerSpells() {
	essenceOfTheRedActionID := core.ActionID{SpellID: 23513}
	
	ai.essenceOfTheRedAura = ai.Target.GetOrRegisterAura(core.Aura{
		Label:     "Essence of the Red",
		ActionID:  essenceOfTheRedActionID,
		Duration:  time.Minute * 4,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			
		},
	})

	ai.essenceOfTheRedSpell = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID: essenceOfTheRedActionID,
		ProcMask: core.ProcMaskEmpty,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Minute * 4,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				ai.essenceOfTheRedAura.Activate(sim)
		},
	})
}

func (ai *VaelastraszTheCorruptAI) Reset(*core.Simulation) {
}

func (ai *VaelastraszTheCorruptAI) ExecuteCustomRotation(sim *core.Simulation) {
//	target := ai.Target.CurrentTarget
//	if target == nil {
		target := &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
//	}
	if ai.essenceOfTheRedSpell.CanCast(sim, target) {
		ai.essenceOfTheRedSpell.Cast(sim, target)
		return
	}
}
