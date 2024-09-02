package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const (
	Truthbearer2H  = 229749
	Truthbearer1H  = 229806
	Stormwrath     = 231387
	LightningsCell = 231784	
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=231387/stormwrath-sanctified-shortblade-of-the-galefinder
	// Equip: Damaging non-periodic spells have a chance to blast up to 3 targets for 181 to 229.
	// (Proc chance: 10%, 100ms cooldown)
	core.NewItemEffect(Stormwrath, func(agent core.Agent) {
		character := agent.GetCharacter()

		maxHits := int(min(3, character.Env.GetNumTargets()))
		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 468670},
			SpellSchool:      core.SpellSchoolNature,
			DefenseType:      core.DefenseTypeMagic,
			ProcMask:         core.ProcMaskEmpty,
			BonusCoefficient: 0.15,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for numHits := 0; numHits < maxHits; numHits++ {
					spell.CalcAndDealDamage(sim, target, sim.Roll(180, 230), spell.OutcomeMagicHitAndCrit)
					target = character.Env.NextTargetUnit(target)
				}
			},
		})

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Millisecond * 100,
		}
		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Chain Lightning (Stormwrath)",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: .10,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				if !icd.IsReady(sim) {
					return
				}
				procSpell.Cast(sim, result.Target)
				icd.Use(sim)
			},
		})
	})

	// https://www.wowhead.com/classic/item=229806/truthbearer
	// Chance on hit: Increases damage done by 15 and attack speed by 30% for 8 sec.
	// TODO: Proc rate assumed and needs testing
	core.NewItemEffect(Truthbearer2H, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(Truthbearer2H)
		crusadersZealAura := crusadersZealAura465414(character)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Truthbearer (Crusaders Zeal)",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          procMask,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			PPM:               1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				crusadersZealAura.Activate(sim)
			},
		})
	})

	// https://www.wowhead.com/classic/item=229749/truthbearer
	// Chance on hit: Increases damage done by 15 and attack speed by 30% for 8 sec.
	// TODO: Proc rate assumed and needs testing
	core.NewItemEffect(Truthbearer1H, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(Truthbearer1H)
		crusadersZealAura := crusadersZealAura465414(character)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Truthbearer (Crusaders Zeal)",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          procMask,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			PPM:               1, 
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				crusadersZealAura.Activate(sim)
			},
		})
	})

	

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=231784/lightnings-cell
	// You gain a charge of Gathering Storm each time you cause a damaging spell critical strike.
	// When you reach 3 charges of Gathering Storm, they will release, firing an Unleashed Storm for 277 to 323 damage.
	// Gathering Storm cannot be gained more often than once every 2.5 sec. (2.5s cooldown)
	core.NewItemEffect(LightningsCell, func(agent core.Agent) {
		character := agent.GetCharacter()

		unleashedStormSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 468782},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(277, 323), spell.OutcomeMagicHitAndCrit)
			},
		})

		chargeAura := character.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 468780},
			Label:     "Lightning's Cell",
			Duration:  core.NeverExpires,
			MaxStacks: 3,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				if aura.GetStacks() == aura.MaxStacks {
					unleashedStormSpell.Cast(sim, aura.Unit.CurrentTarget)
					aura.Deactivate(sim)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Lightning's Cell Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeCrit,
			ProcMask: core.ProcMaskSpellDamage,
			ICD:      time.Millisecond * 2500,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				chargeAura.Activate(sim)
				chargeAura.AddStack(sim)
			},
		})
	})

	core.AddEffectsToTest = true
}

	// https://www.wowhead.com/classic/spell=465414/crusaders-zeal
	// Used by:
	// - https://www.wowhead.com/classic/item=229806/truthbearer and
	// - https://www.wowhead.com/classic/item=229749/truthbearer
	func crusadersZealAura465414(character *core.Character) *core.Aura {
		return character.GetOrRegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 465414},
			Label:    "Crusaders Zeal (465414)",
			Duration: time.Second * 8,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.BonusDamage += 15
				character.MultiplyAttackSpeed(sim, 1.30)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.BonusDamage -= 15
				character.MultiplyAttackSpeed(sim, 1/1.30)
			},
		})
	}