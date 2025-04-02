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
			Name:             "Roar of the Dream Trigger",
			Callback:         core.CallbackOnCastComplete,
			ProcMask:         core.ProcMaskSpellDamage,
			CanProcFromProcs: true,
			ProcChance:       0.05,
			Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=220632/atalai-blood-ritual-medallion
	// Increases your melee and ranged damage by 20 for 20 sec. Every time you hit a target, this bonus is reduced by 1.
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
				character.PseudoStats.BonusPhysicalDamage += float64(newStacks - oldStacks)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					aura.RemoveStack(sim)
				}
			},
		})

		triggerSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
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

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell: triggerSpell,
			Type:  core.CooldownTypeDPS,
		})
	})

	// https://www.wowhead.com/classic/item=220633/atalai-blood-ritual-badge
	// Increases your armor by 1000 and defense skill by 20 for 20 sec.
	// Every time you take melee or ranged damage, this bonus is reduced by 100 armor and 2 defense.
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

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell: triggerSpell,
			Type:  core.CooldownTypeSurvival,
		})
	})

	// https://www.wowhead.com/classic/item=220634/atalai-blood-ritual-charm
	// Increases your spell damage by up to 96 and your healing by up to 192 for 20 sec.
	// Every time you cast a spell, the bonus is reduced by 8 spell damage and 16 healing.
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

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
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
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

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
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagBinary,

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
			Name:              "DMC Decay Spell Hit",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskMelee | core.ProcMaskRanged,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			PPM:               7.0, // Estimate from log
			Handler:           handler,
		})
		hitAura.Icd = &icd

		castAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "DMC Decay Spell Cast",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.35,
			Handler:    handler,
		})
		castAura.Icd = &icd

		character.ItemSwap.RegisterProc(DarkmoonCardDecay, hitAura)
		character.ItemSwap.RegisterProc(DarkmoonCardDecay, castAura)
	})

	core.NewItemEffect(DarkmoonCardSandstorm, func(agent core.Agent) {
		character := agent.GetCharacter()

		tickSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 449288},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

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
			Name:              "Sandstorm Spell Hit",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskMelee | core.ProcMaskRanged,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			PPM:               10.0, // Estimate from log
			Handler:           handler,
		})
		hitAura.Icd = &icd

		castAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Sandstorm Spell Cast",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.30,
			Handler:    handler,
		})

		character.ItemSwap.RegisterProc(DarkmoonCardSandstorm, hitAura)
		character.ItemSwap.RegisterProc(DarkmoonCardSandstorm, castAura)
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

		dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(DragonsCry, 1.0, 0)
		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Emerald Dragon Whelp Proc",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			DPM:               dpm, // Reported by armaments discord
			ICD:               time.Minute * 1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				for _, petAgent := range character.PetAgents {
					if whelp, ok := petAgent.(*guardians.EmeraldDragonWhelp); ok {
						whelp.EnableWithTimeout(sim, whelp, time.Second*15)
						break
					}
				}
			},
		})

		character.ItemSwap.RegisterProc(DragonsCry, triggerAura)
	})

	core.NewItemEffect(CobraFangClaw, func(agent core.Agent) {
		character := agent.GetCharacter()

		dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(CobraFangClaw, 1.0, 0)

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:         "Cobra Fang Claw Extra Attack",
			Duration:     core.NeverExpires,
			Outcome:      core.OutcomeLanded,
			Callback:     core.CallbackOnSpellHitDealt,
			DPM:          dpm,
			DPMProcCheck: core.DPMProc,
			Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				character.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 220588}, spell)
			},
		})

		character.ItemSwap.RegisterProc(CobraFangClaw, triggerAura)
	})

	serpentsStrikerEffect := func(character *core.Character) *core.Spell {
		procAuras := character.NewEnemyAuraArray(core.SerpentsStrikerFistDebuffAura)

		return character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 447894},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagPoison | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

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

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Bloodthirst Crossbow Proc Aura",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskRanged,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				thirstForBlood.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterProc(BloodthirstCrossbow, triggerAura)
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

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Engulfing Shadows",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: .10,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Damage > 0 {
					procSpell.Cast(sim, character.CurrentTarget)
				}
			},
		})

		character.ItemSwap.RegisterProc(BladeOfEternalDarkness, triggerAura)
	})

	// https://www.wowhead.com/classic/item=220965/scalebane-greataxe
	// Equip: +117 Attack Power when fighting Dragonkin.
	core.NewMobTypeAttackPowerEffect(ScalebaneGreataxe, []proto.MobType{proto.MobType_MobTypeDragonkin}, 117)

	core.AddEffectsToTest = true
}
