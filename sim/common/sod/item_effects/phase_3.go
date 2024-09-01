package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	// Ordered by ID
	FistOfTheForsaken          = 220578
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
	SerpentsStrikerSlow        = 224409
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
	//                                 Rings
	///////////////////////////////////////////////////////////////////////////

	core.NewItemEffect(RoarOfTheDream, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Roar of the Dream", core.ActionID{SpellID: 446706}, stats.Stats{stats.SpellDamage: 66}, time.Second*10)
		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 446705},
			Name:       "Roar of the Dream Trigger",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellOrProc,
			ProcChance: 0.05,
			Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	core.NewItemEffect(AtalaiBloodRitualMedallion, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{ItemID: AtalaiBloodRitualMedallion}

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
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

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

		actionID := core.ActionID{ItemID: AtalaiBloodRitualBadge}
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
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

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

		actionID := core.ActionID{ItemID: AtalaiBloodRitualCharm}
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
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

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
			Duration: time.Second * 3,
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
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee | core.ProcMaskRanged,
			PPM:      7.0, // Estimate from log
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

		tickSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 449288},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(50, 100), spell.OutcomeMagicCrit)
			},
		})

		// Sandstorm lasts for 10 seconds and moves in an outward spiral. It seems to be able to hit the same boss target
		// multiple times during this duration, especially depending on size and positioning.
		// On Hakkar seems to on average hit anywhere from 1-3 times
		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 446388},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			Dot: core.DotConfig{
				IsAOE: true,
				Aura: core.Aura{
					Label: "Sandstorm Hit",
				},
				NumberOfTicks: 1,
				TickLength:    time.Second * 1,

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					for _, aoeTarget := range sim.Encounter.TargetUnits {
						tickSpell.Cast(sim, aoeTarget)
					}
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dot := spell.AOEDot()
				dot.NumberOfTicks = int32(sim.Roll(1, 3))
				dot.Apply(sim)
				dot.TickOnce(sim)
			},
		})

		// Custom ICD so it can be shared by both proc triggers
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 8,
		}

		handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if !icd.IsReady(sim) {
				return
			}

			icd.Use(sim)
			procSpell.Cast(sim, character.CurrentTarget)
		}

		hitAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{SpellID: 446389},
			Name:     "Sandstorm Spell Hit",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee | core.ProcMaskRanged,
			PPM:      10.0, // Estimate from log
			Handler:  handler,
		})
		hitAura.Icd = &icd

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 446389},
			Name:       "Sandstorm Spell Cast",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.30,
			Handler:    handler,
		})
	})

	core.NewSimpleStatOffensiveTrinketEffect(RoarOfTheGuardian, stats.Stats{stats.AttackPower: 70, stats.RangedAttackPower: 70}, time.Second*20, time.Minute*5)

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	itemhelpers.CreateWeaponProcSpell(FistOfTheForsaken, "Fist of the Forsaken", 7.0, func(character *core.Character) *core.Spell {
		actionID := core.ActionID{SpellID: 446317}
		healthMetrics := character.NewHealthMetrics(actionID)

		return character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: .20,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealDamage(sim, target, 39, spell.OutcomeAlwaysHit)
				character.GainHealth(sim, result.Damage, healthMetrics)
			},
		})
	})

	core.NewItemEffect(DragonsCry, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(DragonsCry)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Emerald Dragon Whelp Proc",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: procMask,
			PPM:      1.0, // Reported by armaments discord
			ICD:      time.Minute * 1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				for _, petAgent := range character.PetAgents {
					if whelp, ok := petAgent.(*guardians.EmeraldDragonWhelp); ok {
						whelp.EnableWithTimeout(sim, whelp, time.Second*15)
						break
					}
				}
			},
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
					character.AutoAttacks.ExtraMHAttackProc(sim , 1, core.ActionID{SpellID: 220588}, spell)
				}
			},
		})
	})

	serpentsStrikerEffect := func(character *core.Character) *core.Spell {
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
				if target.Level <= 55 {
					spell.CalcAndDealDamage(sim, target, 50, spell.OutcomeMagicHitAndCrit)

					procAuras.Get(target).Activate(sim)
				}
			},
		})
	}
	itemhelpers.CreateWeaponProcSpell(SerpentsStriker, "Serpent's Striker", 5.0, serpentsStrikerEffect)
	itemhelpers.CreateWeaponProcSpell(SerpentsStrikerSlow, "Serpent's Striker", 5.0, serpentsStrikerEffect)

	core.NewItemEffect(BloodthirstCrossbow, func(agent core.Agent) {
		character := agent.GetCharacter()

		healthMetrics := character.NewHealthMetrics(core.ActionID{SpellID: 446725})

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
			ProcChance: .10,
			Handler:    handler,
		})
	})

	core.NewItemEffect(ScalebaneGreataxe, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			character.PseudoStats.MobTypeAttackPower += 117
		}
	})

	core.AddEffectsToTest = true
}
