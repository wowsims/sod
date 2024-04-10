package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	// Ordered by ID
	ShortswordOfVengeance    = 754
	Bloodrazor               = 809
	HammerOfTheNorthernWind  = 810
	Shadowblade              = 2163
	Gutwrencher              = 5616
	HanzoSword               = 8190
	BloodletterScalpel       = 9511
	TheHandOfAntusul         = 9639
	GryphonRidersStormhammer = 9651
	Firebreather             = 10797
	VilerendSlicer           = 11603
	HookfangShanker          = 11635
	LinkensSwordOfMastery    = 11902
	SearingNeedle            = 12531
	SerpentSlicer            = 13035
	JoonhosMercy             = 17054
	ThrashBlade              = 17705
	SatyrsLash               = 17752
	CobraFangClaw            = 220588
	SerpentsStriker          = 220589
	AtalaiBloodRitualCharm   = 220634
	ScalebaneGreataxe        = 220965
	DarkmoonCardDecay        = 221307
	DarkmoonCardOvergrowth   = 221308
	DarkmoonCardSandstorm    = 221309
	RoarOfTheDream           = 221440
	RoarOfTheGuardian        = 221442
	BloodthirstCrossbow      = 221451
	FistOfStone              = 223524
	BladeOfEternalDarkness   = 223964
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

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealDamage(sim, target, 40, spell.OutcomeMagicHitAndCrit)
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

		handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			decayProcSpell.Cast(sim, character.CurrentTarget)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{SpellID: 446392},
			Name:     "DMC Decay Spell Hit",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskMelee | core.ProcMaskRanged,
			PPM:      5.0, // Placeholder proc value
			ICD:      time.Millisecond * 200,
			Handler:  handler,
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 450110},
			Name:       "DMC Decay Spell Cast",
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

	itemhelpers.CreateWeaponProcSpell(Bloodrazor, "Bloodrazor", 1.0, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 17504},
			SpellSchool:      core.SpellSchoolPhysical,
			DefenseType:      core.DefenseTypeMelee,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
			Dot: core.DotConfig{
				NumberOfTicks: 10,
				TickLength:    time.Second * 3,
				Aura: core.Aura{
					Label: "Rend",
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.Spell.CalcAndDealPeriodicDamage(sim, target, 12, dot.OutcomeTick)
				},
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(HanzoSword, "Hanzo Sword", 1.0, 16405, core.SpellSchoolPhysical, 75, 0, 0, core.DefenseTypeMelee)

	itemhelpers.CreateWeaponProcDamage(HammerOfTheNorthernWind, "Hammer of the Northern Wind", 3.5, 13439, core.SpellSchoolFrost, 20, 10, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcDamage(Shadowblade, "Shadow Blade", 1.0, 18138, core.SpellSchoolShadow, 110, 30, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcSpell(HookfangShanker, "Hookfang Shanker", 1.0, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 13526},
			SpellSchool:      core.SpellSchoolNature,
			DefenseType:      core.DefenseTypeMagic,
			ProcMask:         core.ProcMaskEmpty,
			Flags:            core.SpellFlagPoison,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
			Dot: core.DotConfig{
				NumberOfTicks: 10,
				TickLength:    time.Second * 3,
				Aura: core.Aura{
					Label: "Corrosive Poison",
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.Armor: -50})
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.Armor: 50})
					},
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.Spell.CalcAndDealPeriodicDamage(sim, target, 7, dot.OutcomeTick)
				},
			},
		})
	})

	itemhelpers.CreateWeaponProcSpell(Firebreather, "Firebreather", 1.0, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 16413},
			SpellSchool:      core.SpellSchoolFire,
			DefenseType:      core.DefenseTypeMagic,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealDamage(sim, target, 70, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
			Dot: core.DotConfig{
				NumberOfTicks: 3,
				TickLength:    time.Second * 2,
				Aura: core.Aura{
					Label: "Fireball",
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.Spell.CalcAndDealPeriodicDamage(sim, target, 3, dot.OutcomeTick)
				},
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(VilerendSlicer, "Vilerend Slicer", 1.0, 16405, core.SpellSchoolPhysical, 75, 0, 0, core.DefenseTypeMelee)

	itemhelpers.CreateWeaponProcSpell(ThrashBlade, "Thrash Blade", 1.0, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 21919},
			SpellSchool:      core.SpellSchoolPhysical,
			DefenseType:      core.DefenseTypeMelee,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				character.AutoAttacks.ExtraMHAttack(sim)
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(JoonhosMercy, "Joonho's Mercy", 1.0, 20883, core.SpellSchoolArcane, 70, 0, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcDamage(SatyrsLash, "Satyr's Lash", 1.0, 18205, core.SpellSchoolShadow, 55, 30, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcSpell(FistOfStone, "Fist of Stone", 1.0, func(character *core.Character) *core.Spell {
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 21951})

		return character.RegisterSpell(core.SpellConfig{
			ActionID: manaMetrics.ActionID,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				character.AddMana(sim, 50, manaMetrics)
			},
		})
	})

	itemhelpers.CreateWeaponProcSpell(SerpentSlicer, "Serpent Slicer", 1.0, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 17511},
			SpellSchool:      core.SpellSchoolNature,
			DefenseType:      core.DefenseTypeMagic,
			ProcMask:         core.ProcMaskEmpty,
			Flags:            core.SpellFlagPoison,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
			Dot: core.DotConfig{
				NumberOfTicks: 10,
				TickLength:    time.Second * 2,
				Aura: core.Aura{
					Label: "Poison",
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.Spell.CalcAndDealPeriodicDamage(sim, target, 8, dot.OutcomeTick)
				},
			},
		})
	})

	// TODO Searing Needle adds an "Apply Aura: Mod Damage Done (Fire): 10" aura to the /target/, buffing it; not currently modelled
	itemhelpers.CreateWeaponProcDamage(SearingNeedle, "Searing Needle", 1.0, 16454, core.SpellSchoolFire, 60, 0, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcSpell(TheHandOfAntusul, "The Hand of Antu'sul", 1.0, func(character *core.Character) *core.Spell {
		debuffAuras := character.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
			aura := unit.GetOrRegisterAura(core.Aura{
				Label:    "ThunderClap-Antu'sul",
				ActionID: core.ActionID{SpellID: 13532},
				Duration: time.Second * 10,
			})
			core.AtkSpeedReductionEffect(aura, 1.11)
			return aura
		})

		results := make([]*core.SpellResult, min(4, character.Env.GetNumTargets()))

		return character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 13532},
			SpellSchool:      core.SpellSchoolNature,
			DefenseType:      core.DefenseTypeMagic,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for idx := range results {
					results[idx] = spell.CalcDamage(sim, target, 7, spell.OutcomeMagicHitAndCrit)
					target = character.Env.NextTargetUnit(target)
				}
				for _, result := range results {
					spell.DealDamage(sim, result)
					if result.Landed() {
						debuffAuras.Get(result.Target).Activate(sim)
					}
				}
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(ShortswordOfVengeance, "Shortsword of Vengeance", 1.0, 13519, core.SpellSchoolHoly, 30, 0, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcDamage(GryphonRidersStormhammer, "Gryphon Rider's Stormhammer", 1.0, 18081, core.SpellSchoolNature, 91, 34, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcSpell(Gutwrencher, "Gutwrencher", 1.0, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 16406},
			SpellSchool:      core.SpellSchoolPhysical,
			DefenseType:      core.DefenseTypeMelee,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
			Dot: core.DotConfig{
				NumberOfTicks: 10,
				TickLength:    time.Second * 3,
				Aura: core.Aura{
					Label: "Rend",
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.Spell.CalcAndDealPeriodicDamage(sim, target, 8, dot.OutcomeTick)
				},
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(BloodletterScalpel, "Bloodletter Scalpel", 1.0, 18081, core.SpellSchoolPhysical, 60, 10, 0, core.DefenseTypeMelee)

	itemhelpers.CreateWeaponProcDamage(LinkensSwordOfMastery, "Linken's Sword of Mastery", 1.0, 18089, core.SpellSchoolNature, 45, 30, 0, core.DefenseTypeMagic)

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

	core.NewItemEffect(ScalebaneGreataxe, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			character.AddStats(stats.Stats{
				stats.AttackPower:       93,
				stats.RangedAttackPower: 93,
			})
		}
	})

	core.AddEffectsToTest = true
}
