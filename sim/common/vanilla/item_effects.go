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
	FlurryAxe                = 871
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
	LinkensSwordOfMastery    = 11902
	SearingNeedle            = 12531
	Felstriker               = 12590
	PipsSkinner              = 12709
	SerpentSlicer            = 13035
	SealOfTheDawn            = 13209
	JoonhosMercy             = 17054
	ThrashBlade              = 17705
	SatyrsLash               = 17752
	MarkOfTheChosen          = 17774
	Thunderfury              = 19019
	Nightfall                = 19169
	DarkmoonCardHeroism      = 19287
	DarkmoonCardBlueDragon   = 19288
	DarkmoonCardMaelstrom    = 19289
	RuneOfTheDawn            = 19812
	ScarabBrooch             = 21625
	MarkOfTheChampionPhys    = 23206
	MarkOfTheChampionSpell   = 23207
	ReavingNightfall         = 227843
	BurstOfKnowledge         = 227972
	HandOfJustice            = 227989 // 11815
	Ironfoe                  = 227991 // 11684
	FiendishMachete          = 228056 // 18310
	TalismanOfEphemeralPower = 228255 // 18820
	EssenceOfThePureFlame    = 228293 // 18815
	EmpyreanDemolisher       = 228397 // 17112
	DraconicInfusedEmblem    = 228678 // 22268
	HandOfJustice2           = 228722 // TODO: Unsure why there's a second version of this item

	// SulfurasHandOfRagnaros = 17182
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

	itemhelpers.CreateWeaponProcSpell(FlurryAxe, "Flurry Axe", 1.0, func(character *core.Character) *core.Spell {
		return character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 18797},
			SpellSchool:      core.SpellSchoolPhysical,
			DefenseType:      core.DefenseTypeMelee,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				character.AutoAttacks.ExtraMHAttack(sim, 1, core.ActionID{SpellID: 18797})
			},
		})
	})

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

	itemhelpers.CreateWeaponProcDamage(LinkensSwordOfMastery, "Linken's Sword of Mastery", 1.0, 18089, core.SpellSchoolNature, 45, 30, 0, core.DefenseTypeMagic)

	// TODO Searing Needle adds an "Apply Aura: Mod Damage Done (Fire): 10" aura to the /target/, buffing it; not currently modelled
	itemhelpers.CreateWeaponProcDamage(SearingNeedle, "Searing Needle", 1.0, 16454, core.SpellSchoolFire, 60, 0, 0, core.DefenseTypeMagic)

	core.NewItemEffect(PipsSkinner, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.MobTypeAttackPower += 45
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
				character.AutoAttacks.ExtraMHAttack(sim, 1, core.ActionID{SpellID: 21919})
			},
		})
	})

	itemhelpers.CreateWeaponProcSpell(Ironfoe, "Ironfoe", 1.0, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 15494},
			SpellSchool:      core.SpellSchoolPhysical,
			DefenseType:      core.DefenseTypeMelee,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				character.AutoAttacks.ExtraMHAttack(sim, 2, core.ActionID{SpellID: 15494})
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(SatyrsLash, "Satyr's Lash", 1.0, 18205, core.SpellSchoolShadow, 55, 30, 0, core.DefenseTypeMagic)

	core.NewItemEffect(FiendishMachete, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeElemental {
			character.PseudoStats.MobTypeAttackPower += 36
		}
	})

	core.NewItemEffect(Thunderfury, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(Thunderfury)
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

	nightfallEffect := func(agent core.Agent) {
		character := agent.GetCharacter()

		procAuras := character.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				Label:    "Spell Vulnerability (Nightfall)",
				ActionID: core.ActionID{SpellID: 23605},
				Duration: time.Second * 5,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexArcane] *= 1.15
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexFire] *= 1.15
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexFrost] *= 1.15
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexHoly] *= 1.15
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexNature] *= 1.15
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexShadow] *= 1.15
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexArcane] /= 1.15
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexFire] /= 1.15
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexFrost] /= 1.15
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexHoly] /= 1.15
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexNature] /= 1.15
					aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexShadow] /= 1.15
				},
			})
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Nightfall Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee,
			PPM:      2,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procAuras.Get(result.Target).Activate(sim)
			},
		})
	}
	core.NewItemEffect(Nightfall, nightfallEffect)
	core.NewItemEffect(ReavingNightfall, nightfallEffect)

	core.NewItemEffect(Felstriker, func(agent core.Agent) {
		character := agent.GetCharacter()

		effectAura := character.NewTemporaryStatsAura("Felstriker", core.ActionID{SpellID: 16551}, stats.Stats{stats.MeleeCrit: 100 * core.CritRatingPerCritChance, stats.MeleeHit: 100 * core.MeleeHitRatingPerHitChance}, time.Second*3)
		procMask := character.GetProcMaskForItem(Felstriker)
		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Felstriker Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: procMask,
			PPM:      1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				effectAura.Activate(sim)
			},
		})
	})

	itemhelpers.CreateWeaponProcAura(EmpyreanDemolisher, "Empyrean Demolisher", 1.0, func(character *core.Character) *core.Aura {
		return character.GetOrRegisterAura(core.Aura{
			Label:    "Empyrean Demolisher Haste Aura",
			ActionID: core.ActionID{SpellID: 21165},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.2)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1/1.2)
			},
		})
	})

	// https://www.wowhead.com/classic/item=17182/sulfuras-hand-of-ragnaros
	// Chance on hit: Hurls a fiery ball that causes 273 to 333 Fire damage and an additional 75 damage over 10 sec.
	// Equip: Deals 5 Fire damage to anyone who strikes you with a melee attack.
	// core.NewItemEffect(SulfurasHandOfRagnaros, func(agent core.Agent) {
	// 	character := agent.GetCharacter()

	// 	immolationActionID := core.ActionID{SpellID: 21142}

	// 	immolationSpell := character.RegisterSpell(core.SpellConfig{
	// 		ActionID:    immolationActionID,
	// 		SpellSchool: core.SpellSchoolFire,
	// 		ProcMask:    core.ProcMaskEmpty,

	// 		DamageMultiplier: 1,
	// 		ThreatMultiplier: 1,

	// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 			spell.CalcAndDealDamage(sim, target, 5, spell.OutcomeMagicHit)
	// 		},
	// 	})

	// 	character.RegisterAura(core.Aura{
	// 		ActionID: immolationActionID,
	// 		Label:    "Immolation (Hand of Ragnaros)",
	// 		Duration: core.NeverExpires,
	// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
	// 			aura.Activate(sim)
	// 		},
	// 		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 			if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
	// 				immolationSpell.Cast(sim, spell.Unit)
	// 			}
	// 		},
	// 	})

	// 	fireballSpell := character.RegisterSpell(core.SpellConfig{
	// 		ActionID:    core.ActionID{SpellID: 21162},
	// 		SpellSchool: core.SpellSchoolFire,
	// 		DefenseType: core.DefenseTypeMagic,
	// 		ProcMask:    core.ProcMaskEmpty,

	// 		DamageMultiplier: 1,
	// 		ThreatMultiplier: 1,

	// 		Dot: core.DotConfig{
	// 			Aura: core.Aura{
	// 				Label: "Fireball (Hand of Ragnaros)",
	// 			},
	// 			TickLength:    2 * time.Second,
	// 			NumberOfTicks: 5,

	// 			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
	// 				dot.Snapshot(target, 15, isRollover)
	// 			},

	// 			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	// 			},
	// 		},

	// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 			result := spell.CalcAndDealDamage(sim, target, sim.Roll(273, 333), spell.OutcomeMagicHitAndCrit)
	// 			if result.Landed() {
	// 				spell.Dot(target).Apply(sim)
	// 			}
	// 		},
	// 	})

	// 	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
	// 		Name:     "Hand of Ragnaros Trigger",
	// 		Callback: core.CallbackOnSpellHitDealt,
	// 		Outcome:  core.OutcomeLanded,
	// 		ProcMask: core.ProcMaskMelee,
	// 		PPM:      1, // Estimated based on data from WoW Armaments Discord
	// 		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 			fireballSpell.Cast(sim, result.Target)
	// 		},
	// 	})
	// })

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=13209/seal-of-the-dawn
	// Equip: +81 Attack Power when fighting Undead.
	core.NewItemEffect(SealOfTheDawn, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			character.AddStat(stats.AttackPower, 81)
			character.AddStat(stats.AttackPower, 81)
		}
	})

	// https://www.wowhead.com/classic/item=19812/rune-of-the-dawn
	// Equip: Increases damage done to Undead by magical spells and effects by up to 48.
	core.NewItemEffect(RuneOfTheDawn, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			character.AddStat(stats.SpellDamage, 48)
		}
	})

	core.NewItemEffect(HandOfJustice, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.AutoSwingMelee {
			return
		}

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 2,
		}

		character.RegisterAura(core.Aura{
			Label:    "Hand of Justice",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if sim.RandomFloat("HandOfJustice") < 0.02 {
					icd.Use(sim)
					aura.Unit.AutoAttacks.ExtraMHAttack(sim, 1, core.ActionID{SpellID: 15600})
				}
			},
		})
	})

	// https://www.wowhead.com/classic/item=19287/darkmoon-card-heroism
	// Equip: Sometimes heals bearer of 120 to 180 damage when damaging an enemy in melee.
	core.NewItemEffect(DarkmoonCardHeroism, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 23689}
		healthMetrics := character.NewHealthMetrics(actionID)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskEmpty,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				character.GainHealth(sim, sim.Roll(120, 180), healthMetrics)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Heroism Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee,
			PPM:      2,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, spell.Unit)
			},
		})
	})

	// https://www.wowhead.com/classic/item=19288/darkmoon-card-blue-dragon
	// Equip: 2% chance on successful spellcast to allow 100% of your Mana regeneration to continue while casting for 15 sec. (Proc chance: 2%)
	core.NewItemEffect(DarkmoonCardBlueDragon, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 23684}

		procAura := character.RegisterAura(core.Aura{
			Label:    "Aura of the Blue Dragon",
			ActionID: actionID,
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SpiritRegenRateCasting += 1
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SpiritRegenRateCasting -= 1
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Aura of the Blue Dragon Trigger",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
			ProcChance: .02,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})

	// https://www.wowhead.com/classic/item=19289/darkmoon-card-maelstrom
	// Equip: Chance to strike your melee target with lightning for 200 to 300 Nature damage.
	core.NewItemEffect(DarkmoonCardMaelstrom, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 23687}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(200, 300), spell.OutcomeMagicHitAndCrit)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Lightning Strike Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee,
			PPM:      1.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})
	})

	// https://www.wowhead.com/classic/item=227972/burst-of-knowledge
	// Use: Reduces mana cost of all spells by 100 for 10 sec. (5 Min Cooldown)
	core.NewItemEffect(BurstOfKnowledge, func(agent core.Agent) {
		character := agent.GetCharacter()

		aura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: BurstOfKnowledge},
			Label:    "Burst of Knowledge",
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.CostMultiplier -= 1
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.CostMultiplier += 1
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: BurstOfKnowledge},
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 5,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeMana,
			Spell: spell,
		})
	})

	// https://www.wowhead.com/classic/item=228255/talisman-of-ephemeral-power
	// Use: Increases damage and healing done by magical spells and effects by up to 184 for 15 sec. (1 Min, 30 Sec Cooldown)
	core.NewSimpleStatOffensiveTrinketEffect(TalismanOfEphemeralPower, stats.Stats{stats.SpellPower: 184}, time.Second*15, time.Second*90)

	// https://www.wowhead.com/classic/item=228293/essence-of-the-pure-flame
	// Equip: When struck in combat inflicts 50 Fire damage to the attacker.
	core.NewItemEffect(EssenceOfThePureFlame, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 461694},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 50, spell.OutcomeAlwaysHit)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Fiery Aura Proc",
			Callback: core.CallbackOnSpellHitTaken,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee, // TODO: Unsure if this means melee attacks or all attacks
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, spell.Unit)
			},
		})
	})

	// https://www.wowhead.com/classic/item=228678/draconic-infused-emblem
	// Use: Increases your spell damage by up to 100 and your healing by up to 190 for 15 sec. (1 Min, 30 Sec Cooldown)
	core.NewSimpleStatOffensiveTrinketEffect(DraconicInfusedEmblem, stats.Stats{stats.SpellDamage: 100, stats.HealingPower: 190}, time.Second*15, time.Second*75)

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
			character.PseudoStats.MobTypeAttackPower += 150
		}
	})

	core.NewItemEffect(MarkOfTheChampionSpell, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead || character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			character.PseudoStats.MobTypeSpellPower += 85
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
