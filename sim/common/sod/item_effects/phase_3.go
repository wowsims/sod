package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	// Ordered by ID
	CobraFangClaw          = 220588
	SerpentsStriker        = 220589
	AtalaiBloodRitualCharm = 220634
	DarkmoonCardDecay      = 221307
	DarkmoonCardOvergrowth = 221308
	DarkmoonCardSandstorm  = 221309
	RoarOfTheGuardian      = 221442
	BladeOfEternalDarkness = 223964
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Cloth
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Leather
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Mail
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Plate
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	core.NewItemEffect(AtalaiBloodRitualCharm, func(agent core.Agent) {
		character := agent.GetCharacter()

		bonusPerStack := stats.Stats{
			stats.SpellDamage:  8,
			stats.HealingPower: 16,
		}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Unrestrained Power",
			ActionID:  core.ActionID{SpellID: 446297},
			Duration:  time.Second * 20,
			MaxStacks: 12,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.SetStacks(sim, aura.MaxStacks)
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				bonusStats := bonusPerStack.Multiply(float64(newStacks - oldStacks))
				character.AddStatsDynamic(sim, bonusStats)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
					return
				}
				aura.RemoveStack(sim)
			},
		})

		triggerSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 446297},
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    triggerSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(DarkmoonCardDecay, func(agent core.Agent) {
		character := agent.GetCharacter()

		decayAuras := character.NewEnemyAuraArray(func(target *core.Unit, playerLevel int32) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				Label:     "Decay" + character.Label,
				ActionID:  core.ActionID{SpellID: 446393},
				Duration:  core.NeverExpires,
				MaxStacks: 5,
			})
		})

		decayStackedSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 446810},
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(150, 250), spell.OutcomeMagicHitAndCrit)
			},
		})

		decayProcSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 446393},
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				targetAura := decayAuras.Get(target)
				result := spell.CalcAndDealDamage(sim, target, 20, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					spell.CalcAndDealHealing(sim, &character.Unit, result.Damage, spell.OutcomeHealing)
					targetAura.Activate(sim)
					targetAura.AddStack(sim)
				}
				if targetAura.GetStacks() == 5 {
					decayStackedSpell.Cast(sim, target)
					targetAura.Deactivate(sim)
				}
			},
		})

		handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			decayProcSpell.Cast(sim, character.CurrentTarget)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{SpellID: 446392},
			Name:     "Decay Spell Hit",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskMelee | core.ProcMaskRanged,
			PPM:      5.0, // Placeholder proc value
			ICD:      time.Millisecond * 200,
			Handler:  handler,
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 450110},
			Name:       "Decay Spell Cast",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.35,
			ICD:        time.Millisecond * 200,
			Handler:    handler,
		})
	})

	core.NewItemEffect(DarkmoonCardSandstorm, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 446388},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(100, 200), spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			procSpell.Cast(sim, character.CurrentTarget)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{SpellID: 446389},
			Name:     "Sandstorm",
			Callback: core.CallbackOnCastComplete,
			ProcMask: core.ProcMaskDirect,

			// Placeholder proc value
			ProcChance: 0.025,

			Handler: handler,
		})
	})

	core.NewItemEffect(RoarOfTheGuardian, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Roar of the Guardian",
			ActionID: core.ActionID{SpellID: 446709},
			Duration: time.Second * 20,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.AttackPower, 70)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.AttackPower, -70)
			},
		})

		triggerSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 446709},
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 5,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    triggerSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	core.NewItemEffect(BladeOfEternalDarkness, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 27860}
		manaMetrics := character.NewManaMetrics(actionID)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 100, spell.OutcomeAlwaysHit)
				character.AddMana(sim, 100, manaMetrics)
			},
		})

		handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Damage > 0 {
				procSpell.Cast(sim, character.CurrentTarget)
			}
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 21978},
			Name:       "Engulfing Shadows",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: .10,
			Handler:    handler,
		})
	})

	core.NewItemEffect(CobraFangClaw, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(CobraFangClaw)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		character.RegisterAura(core.Aura{
			Label:    "Cobra Fang Claw Thrash",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}
				if !spell.ProcMask.Matches(procMask) {
					return
				}

				if ppmm.Proc(sim, procMask, "Cobra Fang Claw Extra Attack") {
					character.AutoAttacks.ExtraMHAttack(sim)
				}
			},
		})
	})

	// Serpent's Striker
	itemhelpers.CreateWeaponProcSpell(SerpentsStriker, "Serpent's Striker", 5.0, func(character *core.Character) *core.Spell {
		procAuras := character.NewEnemyAuraArray(core.SerpentsStrikerFistDebuffAura)

		return character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 447894},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagPoison,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 50+0.05*spell.SpellDamage(), spell.OutcomeMagicHitAndCrit)

				procAuras.Get(target).Activate(sim)
			},
		})
	})

	core.AddEffectsToTest = true
}
