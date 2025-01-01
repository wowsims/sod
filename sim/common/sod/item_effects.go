package sod

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

// Ordered by ID
const ()

func init() {
	core.AddEffectsToTest = false

	// ! Please keep items ordered alphabetically within a given category !

	core.AddEffectsToTest = true
}

// For Automatic Crowd Pummeler and Druid's Catnip
func RegisterFiftyPercentHasteBuffCD(character *core.Character, actionID core.ActionID) {
	aura := character.GetOrRegisterAura(core.Aura{
		Label:    "Haste",
		ActionID: core.ActionID{SpellID: 13494},
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			character.MultiplyAttackSpeed(sim, 1.5)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			character.MultiplyAttackSpeed(sim, 1.0/1.5)
		},
	})

	spell := character.GetOrRegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    character.GetFiftyPercentHasteBuffCD(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aura.Activate(sim)
		},
	})

	character.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}

func StormhammerChainLightningProcAura(agent core.Agent) {
	character := agent.GetCharacter()

	maxHits := int(min(3, character.Env.GetNumTargets()))
	procSpell := character.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 463946},
		SpellSchool:      core.SpellSchoolNature,
		DefenseType:      core.DefenseTypeMagic,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		BonusCoefficient: 0.1,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for numHits := 0; numHits < maxHits; numHits++ {
				spell.CalcAndDealDamage(sim, target, sim.Roll(105, 145), spell.OutcomeMagicHitAndCrit)
				target = character.Env.NextTargetUnit(target)
			}
		},
	})

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:              "Chain Lightning (Skyrider's Masterwork Stormhammer Melee)",
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		ProcMask:          core.ProcMaskMelee,
		SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
		PPM:               4, // Someone in the armemnts Discord tested it out to 4 PPM
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			procSpell.Cast(sim, result.Target)
		},
	})

	icd := core.Cooldown{
		Timer:    character.NewTimer(),
		Duration: time.Millisecond * 100,
	}
	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       "Chain Lightning (Skyrider's Masterwork Stormhammer Spell)",
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskSpellDamage,
		ProcChance: .1,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			if icd.IsReady(sim) {
				procSpell.Cast(sim, result.Target)
				icd.Use(sim)
			}
		},
	})
}
