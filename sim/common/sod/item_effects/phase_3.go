package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/common/vanilla"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	// Ordered by ID
	DragonsCry                 = 220582
	CobraFangClaw              = 220588
	SerpentsStriker            = 220589
	AtalaiBloodRitualMedallion = 220632
	AtalaiBloodRitualBadge     = 220633
	AtalaiBloodRitualCharm     = 220634
	ScalebaneGreataxe          = 220965
	DarkmoonCardDecay          = 221307
	DarkmoonCardOvergrowth     = 221308
	DarkmoonCardSandstorm      = 221309
	RoarOfTheDream             = 221440
	RoarOfTheGuardian          = 221442
	BloodthirstCrossbow        = 221451
	FistOfStone                = 223524
	BladeOfEternalDarkness     = 223964
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

	core.NewItemEffect(RoarOfTheDream, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Roar of the Dream", core.ActionID{SpellID: 446705}, stats.Stats{stats.SpellDamage: 66}, time.Second*10)
		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 450110},
			Name:       "Roar of the Dream Trigger",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellOrProc,
			ProcChance: 0.05,
			Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(AtalaiBloodRitualMedallion, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 446289}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Relentless Strength",
			ActionID:  actionID,
			Duration:  time.Second * 20,
			MaxStacks: 20,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.SetStacks(sim, aura.MaxStacks)
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				character.PseudoStats.BonusDamage += float64(newStacks - oldStacks)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					aura.RemoveStack(sim)
				}
			},
		})

		triggerSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
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

	core.NewItemEffect(AtalaiBloodRitualBadge, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 446310}
		bonusPerStack := stats.Stats{
			stats.Armor:   100,
			stats.Defense: 2,
		}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Fragile Armor",
			ActionID:  actionID,
			Duration:  time.Second * 20,
			MaxStacks: 10,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.SetStacks(sim, aura.MaxStacks)
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				bonusStats := bonusPerStack.Multiply(float64(newStacks - oldStacks))
				character.AddStatsDynamic(sim, bonusStats)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				aura.RemoveStack(sim)
			},
		})

		triggerSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
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
			Type:     core.CooldownTypeSurvival,
		})
	})

	core.NewItemEffect(AtalaiBloodRitualCharm, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 446297}
		bonusPerStack := stats.Stats{
			stats.SpellDamage:  8,
			stats.HealingPower: 16,
		}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Unrestrained Power",
			ActionID:  actionID,
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
			ActionID: actionID,
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

		decayAura := character.GetOrRegisterAura(core.Aura{
			Label:     "DMC Decay",
			ActionID:  core.ActionID{SpellID: 446393},
			Duration:  core.NeverExpires,
			MaxStacks: 5,
		})

		decayStackedSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 446810},
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(225, 375), spell.OutcomeMagicHitAndCrit)
			},
		})

		decayProcSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 446393},
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealDamage(sim, target, 40, spell.OutcomeMagicHit)
				if result.Landed() {
					spell.CalcAndDealHealing(sim, &character.Unit, result.Damage, spell.OutcomeHealing)
					decayAura.Activate(sim)
					decayAura.AddStack(sim)
				}
				if decayAura.GetStacks() == 5 {
					decayStackedSpell.Cast(sim, target)
					decayAura.Deactivate(sim)
				}
			},
		})

		// Custom ICD so it can be shared by both proc triggers
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Millisecond * 200,
		}

		handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if !icd.IsReady(sim) {
				return
			}

			icd.Use(sim)
			decayProcSpell.Cast(sim, character.CurrentTarget)
		}

		hitAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{SpellID: 446392},
			Name:     "DMC Decay Spell Hit",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskMelee | core.ProcMaskRanged,
			PPM:      5.0, // Placeholder proc value
			Handler:  handler,
		})
		hitAura.Icd = &icd

		castAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 450110},
			Name:       "DMC Decay Spell Cast",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.35,
			Handler:    handler,
		})
		castAura.Icd = &icd
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
					spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(50, 100), spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 446389},
			Name:       "Sandstorm",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskDirect,
			ProcChance: 0.30,
			ICD:        time.Second * 5,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procSpell.Cast(sim, character.CurrentTarget)
			},
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

	core.NewItemEffect(DragonsCry, func(agent core.Agent) {
		vanilla.MakeEmeraldDragonWhelpTriggerAura(agent, DragonsCry)
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
			BonusCoefficient: 0.05,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 50, spell.OutcomeMagicHitAndCrit)

				procAuras.Get(target).Activate(sim)
			},
		})
	})

	core.NewItemEffect(BloodthirstCrossbow, func(agent core.Agent) {
		character := agent.GetCharacter()

		healthMetrics := character.NewManaMetrics(core.ActionID{SpellID: 446725})

		thirstForBlood := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         healthMetrics.ActionID,
			SpellSchool:      core.SpellSchoolShadow,
			DefenseType:      core.DefenseTypeMagic,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				// TODO this has a "HEALTH_LEECH" and a "HEAL" effect, so maybe it heals for leeched amount + heal?
				spell.CalcAndDealDamage(sim, target, 5, spell.OutcomeMagicHit)
				character.GainHealth(sim, 5, healthMetrics)
			},
		})

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Bloodthirst Crossbow Proc Aura",
			OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskRanged) {
					thirstForBlood.Cast(sim, result.Target)
				}
			},
		}))
	})

	itemhelpers.CreateWeaponProcSpell(FistOfStone, "Fist of Stone", 1.0, func(character *core.Character) *core.Spell {
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 21951})

		return character.RegisterSpell(core.SpellConfig{
			ActionID: manaMetrics.ActionID,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				character.AddMana(sim, 50, manaMetrics)
			},
		})
	})

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
			Outcome:    core.OutcomeLanded,
			ProcChance: .10,
			Handler:    handler,
		})
	})

	core.NewItemEffect(ScalebaneGreataxe, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			character.AddStat(stats.AttackPower, 93)
		}
	})

	core.AddEffectsToTest = true
}
