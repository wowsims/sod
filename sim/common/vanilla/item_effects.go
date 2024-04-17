package vanilla

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	ShortswordOfVengeance    = 754
	FieryWarAxe              = 870
	Bloodrazor               = 809
	HammerOfTheNorthernWind  = 810
	Nightblade               = 1982
	Shadowblade              = 2163
	GutRipper                = 2164
	BowOfSearingArrows       = 2825
	Gutwrencher              = 5616
	Ravager                  = 7717
	HanzoSword               = 8190
	TheJackhammer            = 9423
	PendulumOfDoom           = 9425
	BloodletterScalpel       = 9511
	TheHandOfAntusul         = 9639
	GryphonRidersStormhammer = 9651
	Firebreather             = 10797
	VilerendSlicer           = 11603
	HookfangShanker          = 11635
	HandOfJustice            = 11815
	LinkensSwordOfMastery    = 11902
	SearingNeedle            = 12531
	PipsSkinner              = 12709
	SerpentSlicer            = 13035
	JoonhosMercy             = 17054
	ThrashBlade              = 17705
	SatyrsLash               = 17752
	MarkOfTheChosen          = 17774
	FiendishMachete          = 18310
	Thunderfury              = 19019
	ScarabBrooch             = 21625
	MarkOfTheChampionPhys    = 23206
	MarkOfTheChampionSpell   = 23207
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	itemhelpers.CreateWeaponProcDamage(ShortswordOfVengeance, "Shortsword of Vengeance", 1.0, 13519, core.SpellSchoolHoly, 30, 0, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcSpell(FieryWarAxe, "Fiery War Axe", 1.0, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 18796},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Fiery War Axe Fireball",
				},
				TickLength:    2 * time.Second,
				NumberOfTicks: 3,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 8, isRollover)
				},

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(155, 197)
				result := spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
		})
	})

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

	itemhelpers.CreateWeaponProcDamage(HammerOfTheNorthernWind, "Hammer of the Northern Wind", 3.5, 13439, core.SpellSchoolFrost, 20, 10, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcDamage(Nightblade, "Nightblade", 1.0, 18211, core.SpellSchoolShadow, 125, 150, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcDamage(Shadowblade, "Shadowblade", 1.0, 18138, core.SpellSchoolShadow, 110, 30, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcDamage(GutRipper, "Gut Ripper", 1.0, 18107, core.SpellSchoolPhysical, 95, 26, 0, core.DefenseTypeMelee)

	itemhelpers.CreateWeaponProcSpell(BowOfSearingArrows, "Bow of Searing Arrows", 3.35, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 29638},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeRanged,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(18, 26)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeRangedCritOnly)
			},
		})
	})

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

	itemhelpers.CreateWeaponProcAura(Ravager, "Ravager", 1.0, func(character *core.Character) *core.Aura {
		tickActionID := core.ActionID{SpellID: 9633}
		procActionID := core.ActionID{SpellID: 9632}
		auraActionID := core.ActionID{SpellID: 433801}

		ravegerBladestormTickSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    tickActionID,
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1,
			BonusCoefficient: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				damage := 5.0 + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMeleeSpecialHitAndCrit)
				}
			},
		})

		character.GetOrRegisterSpell(core.SpellConfig{
			SpellSchool: core.SpellSchoolPhysical,
			ActionID:    procActionID,
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Flags:       core.SpellFlagChanneled,
			Dot: core.DotConfig{
				IsAOE: true,
				Aura: core.Aura{
					Label: "Ravager Whirlwind",
				},
				NumberOfTicks:       3,
				TickLength:          time.Second * 3,
				AffectedByCastSpeed: false,
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					ravegerBladestormTickSpell.Cast(sim, target)
				},
			},
		})

		return character.GetOrRegisterAura(core.Aura{
			Label:    "Ravager Bladestorm",
			ActionID: auraActionID,
			Duration: time.Second * 9,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AutoAttacks.CancelAutoSwing(sim)
				dotSpell := character.GetSpell(procActionID)
				dotSpell.AOEDot().Apply(sim)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AutoAttacks.EnableAutoSwing(sim)
				dotSpell := character.GetSpell(procActionID)
				dotSpell.AOEDot().Cancel(sim)
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(HanzoSword, "Hanzo Sword", 1.0, 16405, core.SpellSchoolPhysical, 75, 0, 0, core.DefenseTypeMelee)

	itemhelpers.CreateWeaponProcAura(TheJackhammer, "The Jackhammer", 1.0, func(character *core.Character) *core.Aura {
		return character.GetOrRegisterAura(core.Aura{
			Label:    "The Jackhammer Haste Aura",
			ActionID: core.ActionID{SpellID: 13533},
			Duration: time.Second * 9,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.3)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1/1.3)
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(PendulumOfDoom, "Pendulum of Doom", 0.5, 10373, core.SpellSchoolPhysical, 250, 100, 0, core.DefenseTypeMelee)

	itemhelpers.CreateWeaponProcDamage(BloodletterScalpel, "Bloodletter Scalpel", 1.0, 18081, core.SpellSchoolPhysical, 60, 10, 0, core.DefenseTypeMelee)

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

	itemhelpers.CreateWeaponProcDamage(GryphonRidersStormhammer, "Gryphon Rider's Stormhammer", 1.0, 18081, core.SpellSchoolNature, 91, 34, 0, core.DefenseTypeMagic)

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

	core.NewItemEffect(HandOfJustice, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.AutoSwingMelee {
			return
		}

		var handOfJusticeSpell *core.Spell
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 2,
		}
		procChance := 0.013333

		character.RegisterAura(core.Aura{
			Label:    "Hand of Justice",
			Duration: core.NeverExpires,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				config := *character.AutoAttacks.MHConfig()
				config.ActionID = core.ActionID{ItemID: HandOfJustice}
				handOfJusticeSpell = character.GetOrRegisterSpell(config)
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// https://wotlk.wowhead.com/spell=15600/hand-of-justice, proc mask = 20.
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if sim.RandomFloat("HandOfJustice") > procChance {
					return
				}
				icd.Use(sim)

				aura.Unit.AutoAttacks.MaybeReplaceMHSwing(sim, handOfJusticeSpell).Cast(sim, result.Target)
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(LinkensSwordOfMastery, "Linken's Sword of Mastery", 1.0, 18089, core.SpellSchoolNature, 45, 30, 0, core.DefenseTypeMagic)

	// TODO Searing Needle adds an "Apply Aura: Mod Damage Done (Fire): 10" aura to the /target/, buffing it; not currently modelled
	itemhelpers.CreateWeaponProcDamage(SearingNeedle, "Searing Needle", 1.0, 16454, core.SpellSchoolFire, 60, 0, 0, core.DefenseTypeMagic)

	core.NewItemEffect(PipsSkinner, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.AddStat(stats.AttackPower, 45)
		}
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

	itemhelpers.CreateWeaponProcDamage(JoonhosMercy, "Joonho's Mercy", 1.0, 20883, core.SpellSchoolArcane, 70, 0, 0, core.DefenseTypeMagic)

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

	itemhelpers.CreateWeaponProcDamage(SatyrsLash, "Satyr's Lash", 1.0, 18205, core.SpellSchoolShadow, 55, 30, 0, core.DefenseTypeMagic)

	core.NewItemEffect(FiendishMachete, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeElemental {
			character.AddStat(stats.AttackPower, 36)
		}
	})

	core.NewItemEffect(Thunderfury, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(19019)
		ppmm := character.AutoAttacks.NewPPMManager(6.0, procMask)

		procActionID := core.ActionID{SpellID: 21992}

		singleTargetSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(1),
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 0.5,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 300, spell.OutcomeMagicHitAndCrit)
			},
		})

		debuffAuras := character.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				Label:    "Thunderfury",
				ActionID: procActionID,
				Duration: time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, -25)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, 25)
				},
			})
		})

		results := make([]*core.SpellResult, min(5, character.Env.GetNumTargets()))

		bounceSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(2),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			ThreatMultiplier: 1,
			FlatThreatBonus:  63,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for idx := range results {
					results[idx] = spell.CalcDamage(sim, target, 0, spell.OutcomeMagicHit)
					target = sim.Environment.NextTargetUnit(target)
				}
				for _, result := range results {
					if result.Landed() {
						debuffAuras[result.Target.Index].Activate(sim)
					}
					spell.DealDamage(sim, result)
				}
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "Thunderfury",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Thunderfury") {
					singleTargetSpell.Cast(sim, result.Target)
					bounceSpell.Cast(sim, result.Target)
				}
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	core.NewItemEffect(ScarabBrooch, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: ScarabBrooch}

		shieldSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 26470},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellHealing,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Shield: core.ShieldConfig{
				Aura: core.Aura{
					Label:    "Scarab Brooch Shield",
					Duration: time.Second * 30,
				},
			},
		})

		activeAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Persistent Shield",
			ActionID: core.ActionID{SpellID: 26467},
			Callback: core.CallbackOnHealDealt,
			Duration: time.Second * 30,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				shieldSpell.Shield(result.Target).Apply(sim, result.Damage*0.15)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				activeAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	core.NewItemEffect(MarkOfTheChampionPhys, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead || character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			character.AddStat(stats.AttackPower, 150)
		}
	})

	core.NewItemEffect(MarkOfTheChampionSpell, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead || character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			character.AddStat(stats.SpellDamage, 85)
		}
	})

	core.NewItemEffect(MarkOfTheChosen, func(agent core.Agent) {
		character := agent.GetCharacter()
		statIncrease := float64(25)
		markProcChance := 0.02

		procAura := character.RegisterAura(core.Aura{
			Label:    "Mark of the Chosen Effect",
			ActionID: core.ActionID{SpellID: 21970},
			Duration: time.Minute,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatsDynamic(sim, stats.Stats{
					stats.Stamina:   statIncrease,
					stats.Agility:   statIncrease,
					stats.Strength:  statIncrease,
					stats.Intellect: statIncrease,
					stats.Spirit:    statIncrease,
				})
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatsDynamic(sim, stats.Stats{
					stats.Stamina:   -statIncrease,
					stats.Agility:   -statIncrease,
					stats.Strength:  -statIncrease,
					stats.Intellect: -statIncrease,
					stats.Spirit:    -statIncrease,
				})
			},
		})

		core.MakePermanent(character.RegisterAura(core.Aura{
			Label:    "Mark of the Chosen",
			ActionID: core.ActionID{SpellID: 21969},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) && sim.RandomFloat("Mark of the Chosen") < markProcChance {
					procAura.Activate(sim)
				}
			},
		}))
	})

	core.AddEffectsToTest = true
}
