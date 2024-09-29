package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	Heartstriker                 = 230253
	DrakeTalonCleaver            = 230271 // 19353
	ClawOfChromaggus             = 230794
	JekliksCrusher               = 230911
	ZulianSlicer                 = 230930
	WillOfArlokk                 = 230939
	HaldberdOfSmiting            = 230991
	TigulesHarpoon               = 231272
	GrileksCarver                = 231273
	GrileksGrinder               = 231274
	PitchforkOfMadness           = 231277
	Stormwrath                   = 231387
	WrathOfWray                  = 231779
	LightningsCell               = 231784
	Windstriker                  = 231817
	GrileksCarverBloodied        = 231846
	GrileksGrinderBloodied       = 231847
	TigulesHarpoonBloodied       = 231849
	WillOfArlokkBloodied         = 231850
	JekliksCrusherBloodied       = 231861
	PitchforkOfMadnessBloodied   = 231864
	HaldberdOfSmitingBloodied    = 231870
	ZulianSlicerBloodied         = 231876
	ClawOfChromaggusShadowflame  = 232557
	DrakeTalonCleaverShadowflame = 232562
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=230794/claw-of-chromaggus
	// Your offensive spellcasts increase the spell damage of a random school of magic by 50 for 10 sec. (9.5s cooldown)
	core.NewItemEffect(ClawOfChromaggus, func(agent core.Agent) {
		ClawOfChromaggusEffect(agent.GetCharacter())
	})
	// https://www.wowhead.com/classic/item=232557/claw-of-chromaggus
	core.NewItemEffect(ClawOfChromaggusShadowflame, func(agent core.Agent) {
		ClawOfChromaggusEffect(agent.GetCharacter())
	})

	// https://www.wowhead.com/classic/item=230271/drake-talon-cleaver
	// Chance on hit: Delivers a fatal wound for 300 damage.
	// Original proc rate 1.0 increased to approximately 1.60 in SoD phase 5
	itemhelpers.CreateWeaponCoHProcDamage(DrakeTalonCleaver, "Drake Talon Cleaver", 1.0, 467167, core.SpellSchoolPhysical, 300, 0, 0.0, core.DefenseTypeMelee) // TBD confirm 1 ppm in SoD
	// https://www.wowhead.com/classic/item=232562/drake-talon-cleaver
	itemhelpers.CreateWeaponCoHProcDamage(DrakeTalonCleaverShadowflame, "Drake Talon Cleaver", 1.0, 467167, core.SpellSchoolPhysical, 300, 0, 0.0, core.DefenseTypeMelee) // TBD confirm 1 ppm in SoD

	// https://www.wowhead.com/classic/item=231273/grileks-carver
	// +141 Attack Power when fighting Dragonkin.
	core.NewItemEffect(GrileksCarver, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			character.PseudoStats.MobTypeAttackPower += 141
		}
	})
	core.NewItemEffect(GrileksCarverBloodied, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			character.PseudoStats.MobTypeAttackPower += 141
		}
	})

	// https://www.wowhead.com/classic/item=231274/grileks-grinder
	// +60 Attack Power when fighting Dragonkin.
	core.NewItemEffect(GrileksGrinder, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			character.PseudoStats.MobTypeAttackPower += 60
		}
	})
	core.NewItemEffect(GrileksGrinderBloodied, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			character.PseudoStats.MobTypeAttackPower += 60
		}
	})

	// https://www.wowhead.com/classic/item=230991/halberd-of-smiting
	// Equip: Chance to decapitate the target on a melee swing, causing 452 to 676 damage.
	itemhelpers.CreateWeaponEquipProcDamage(HaldberdOfSmiting, "Halberd of Smiting", 2.1, 467819, core.SpellSchoolPhysical, 452, 224, 0.0, core.DefenseTypeMelee)         // Works as phantom strike
	itemhelpers.CreateWeaponEquipProcDamage(HaldberdOfSmitingBloodied, "Halberd of Smiting", 2.1, 467819, core.SpellSchoolPhysical, 452, 224, 0.0, core.DefenseTypeMelee) // Works as phantom strike

	// https://www.wowhead.com/classic/item=230911/jekliks-crusher
	// Chance on hit: Wounds the target for 200 to 220 damage.
	// Original proc rate 4.0 lowered to 1.5 in SoD phase 5
	itemhelpers.CreateWeaponCoHProcDamage(JekliksCrusher, "Jeklik's Crusher", 1.5, 467642, core.SpellSchoolPhysical, 200, 20, 0.0, core.DefenseTypeMelee)
	itemhelpers.CreateWeaponCoHProcDamage(JekliksCrusherBloodied, "Jeklik's Crusher", 1.5, 467642, core.SpellSchoolPhysical, 200, 20, 0.0, core.DefenseTypeMelee)

	// https://www.wowhead.com/classic/item=231277/pitchfork-of-madness
	// +141 Attack Power when fighting Demons.
	core.NewItemEffect(PitchforkOfMadness, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			character.PseudoStats.MobTypeAttackPower += 141
		}
	})
	core.NewItemEffect(PitchforkOfMadnessBloodied, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			character.PseudoStats.MobTypeAttackPower += 141
		}
	})

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

	// https://www.wowhead.com/classic/item=231272/tigules-harpoon
	// +99 Attack Power when fighting Beasts.
	core.NewItemEffect(TigulesHarpoon, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.MobTypeAttackPower += 99
		}
	})
	core.NewItemEffect(TigulesHarpoonBloodied, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.MobTypeAttackPower += 99
		}
	})

	// https://www.wowhead.com/classic/item=230939/will-of-arlokk
	// Use: Calls forth a charmed snake to worship you, increasing your Spirit by 200 for 20 sec. (2 Min Cooldown)
	core.NewItemEffect(WillOfArlokk, func(agent core.Agent) {
		character := agent.GetCharacter()
		makeWillOfWarlookOnUseEffect(character, WillOfArlokk)
	})
	core.NewItemEffect(WillOfArlokkBloodied, func(agent core.Agent) {
		character := agent.GetCharacter()
		makeWillOfWarlookOnUseEffect(character, WillOfArlokkBloodied)
	})

	// https://www.wowhead.com/classic/item=231817/windstriker
	// Chance on hit: All attacks are guaranteed to land and will be critical strikes for the next 3 sec.
	core.NewItemEffect(Windstriker, func(agent core.Agent) {
		character := agent.GetCharacter()

		effectAura := character.NewTemporaryStatsAura("Felstriker", core.ActionID{SpellID: 16551}, stats.Stats{stats.MeleeCrit: 100 * core.CritRatingPerCritChance, stats.MeleeHit: 100 * core.MeleeHitRatingPerHitChance}, time.Second*3)
		procMask := character.GetProcMaskForItem(Windstriker)
		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Felstriker Trigger",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          procMask,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			PPM:               0.6,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				effectAura.Activate(sim)
			},
		})
	})

	// https://www.wowhead.com/classic/item=230930/zulian-slicer
	// Chance on hit: Slices the enemy for 72 to 96 Nature damage.
	itemhelpers.CreateWeaponCoHProcDamage(ZulianSlicer, "Zulian Slicer", 1.2, 467738, core.SpellSchoolNature, 72, 24, 0.35, core.DefenseTypeMelee)
	itemhelpers.CreateWeaponCoHProcDamage(ZulianSlicerBloodied, "Zulian Slicer", 1.2, 467738, core.SpellSchoolNature, 72, 24, 0.35, core.DefenseTypeMelee)

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=231784/lightnings-cell
	// You gain a charge of Gathering Storm each time you cause a damaging spell critical strike.
	// When you reach 3 charges of Gathering Storm, they will release, firing an Unleashed Storm for 277 to 323 damage.
	// Gathering Storm cannot be gained more often than once every 2 sec. (2s cooldown)
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
			ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc, // Procs on procs
			ICD:      time.Millisecond * 2000,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				chargeAura.Activate(sim)
				chargeAura.AddStack(sim)
			},
		})
	})

	// https://www.wowhead.com/classic/item=230253/hearstriker
	// Equip: 2% chance on ranged hit to gain 1 extra attack. (Proc chance: 1%, 1s cooldown) // obviously something wrong here lol
	core.NewItemEffect(Heartstriker, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.AutoSwingRanged {
			return
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Heartstrike",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskRanged,
			ProcChance:        0.02,
			ICD:               time.Second * 1,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,

			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				spell.Unit.AutoAttacks.ExtraRangedAttack(sim, 1, core.ActionID{SpellID: 461164}, spell.ActionID)
			},
		})
	})

	core.NewSimpleStatOffensiveTrinketEffect(WrathOfWray, stats.Stats{stats.Strength: 92}, time.Second*20, time.Minute*2)

	core.AddEffectsToTest = true
}

func makeWillOfWarlookOnUseEffect(character *core.Character, itemID int32) {
	actionID := core.ActionID{ItemID: itemID}

	buffAura := character.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Serpentine Spirit",
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			character.AddStatDynamic(sim, stats.Spirit, 200)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			character.AddStatDynamic(sim, stats.Spirit, -200)
		},
	})

	spell := character.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
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
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

/*
Your offensive spellcasts increase the spell damage of a random school of magic by 50 for 10 sec.
(9.5s cooldown)
The in-game implementation is different for every class based on what schools are considered useful.
Each class has a set of "good" schools based on certain parameters. These "good" schools have a higher chance to be procced.
For each proc the game rolls to determine which school buff will be given.
There's no buff for the Holy school, so there are 5 total schools that can be procced.

schoolChances should be a map of school indexes with the relative proc chance of that school for the given class
*/
func ClawOfChromaggusEffect(character *core.Character) {
	arcaneChance, fireChance, frostChance, natureChance, shadowChance := 0.20, 0.20, 0.20, 0.20, 0.20

	switch character.Class {
	case proto.Class_ClassDruid:
		// Assuming 25% Arcane, 25% Nature, 50% divided among the other 3
		arcaneChance, natureChance = 0.25, 0.25
		fireChance, frostChance, shadowChance = 0.50/3, 0.50/3, 0.50/3
	case proto.Class_ClassMage:
		// The weapon's effect for mage is specialized, based off of selected runes
		if character.HasRuneById(int32(proto.MageRune_RuneBeltFrostfireBolt)) {
			fireChance, frostChance = 0.25, 0.25
			arcaneChance, natureChance, shadowChance = 0.50/3, 0.50/3, 0.50/3
			// Never implemented differently for Spellfrost Bolt
			// } else if character.HasRuneById(int32(proto.MageRune_RuneBeltSpellfrostBolt)) {
			// 		arcaneChance, frostChance = 0.25, 0.25
			// 		fireChance, natureChance, shadowChance = 0.50/3, 0.50/3, 0.50/3
		} else if character.HasRuneById(int32(proto.MageRune_RuneBeltMissileBarrage)) {
			arcaneChance = 0.50
			fireChance, frostChance, natureChance, shadowChance = 0.125, 0.125, 0.125, 0.125
		}
	case proto.Class_ClassPriest:
		// Confirmed 50% proc chance for Shadow and the other 50% divided among the other 4 schools
		shadowChance = 0.50
		arcaneChance, fireChance, frostChance, natureChance = 0.125, 0.125, 0.125, 0.125
	case proto.Class_ClassShaman:
		// Assuming 25% Fire, 25% Nature, 50% divided among the other 3
		fireChance, natureChance = 0.25, 0.25
		arcaneChance, frostChance, shadowChance = 0.50/3, 0.50/3, 0.50/3
	case proto.Class_ClassWarlock:
		if character.HasRuneById(int32(proto.WarlockRune_RuneBracerIncinerate)) {
			// Confirmed 50% Fire, 50% divided among the other 4
			fireChance = 0.50
			arcaneChance, frostChance, natureChance, shadowChance = 0.125, 0.125, 0.125, 0.125
		} else {
			// Confirmed 50% Shadow, 50% divided among the other 4
			shadowChance = 0.50
			arcaneChance, fireChance, frostChance, natureChance = 0.125, 0.125, 0.125, 0.125
		}
	}

	duration := time.Second * 10

	arcaneAura := character.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467410},
		Label:    "Brood Boon: Bronze",
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.ArcanePower, 50)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.ArcanePower, -50)
		},
	})

	fireAura := character.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467414},
		Label:    "Brood Boon: Red",
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.FirePower, 50)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.FirePower, -50)
		},
	})

	frostAura := character.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467412},
		Label:    "Brood Boon: Blue",
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.FrostPower, 50)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.FrostPower, -50)
		},
	})

	natureAura := character.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467413},
		Label:    "Brood Boon: Green",
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.NaturePower, 50)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.NaturePower, -50)
		},
	})

	shadowAura := character.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467411},
		Label:    "Brood Boon: Black",
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.ShadowPower, 50)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.ShadowPower, -50)
		},
	})

	arcaneRangeMax := 0.0 + arcaneChance
	fireRangeMax := arcaneRangeMax + fireChance
	frostRangeMax := fireRangeMax + frostChance
	natureRangeMax := frostRangeMax + natureChance
	shadowRangeMax := natureRangeMax + shadowChance

	if shadowRangeMax > 1.0 {
		panic("Invalid school chances provided to Claw of Chromaggus effect.")
	}

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:     "Claw of the Chromatic Trigger",
		Callback: core.CallbackOnCastComplete,
		ProcMask: core.ProcMaskSpellDamage,
		ICD:      time.Millisecond * 9500,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			roll := sim.RandomFloat("Claw of Chromaggus")
			if roll < arcaneRangeMax {
				arcaneAura.Activate(sim)
			} else if roll < fireRangeMax {
				fireAura.Activate(sim)
			} else if roll < frostRangeMax {
				frostAura.Activate(sim)
			} else if roll < natureRangeMax {
				natureAura.Activate(sim)
			} else if roll < shadowRangeMax {
				shadowAura.Activate(sim)
			}
		},
	})
}
