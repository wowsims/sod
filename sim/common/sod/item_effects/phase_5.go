package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/vanilla"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	Heartstriker                    = 230253
	ClawOfChromaggus                = 230794
	WillOfArlokk                    = 230939
	NatPaglesFishTerminator         = 231016
	BlazefuryRetributer             = 231275
	Stormwrath                      = 231387
	WrathOfWray                     = 231779
	LightningsCell                  = 231784
	Windstriker                     = 231817
	NatPaglesFishTerminatorBloodied = 231848
	WillOfArlokkBloodied            = 231850
	BlazefuryRetributerBloodied     = 231862
	ClawOfChromaggusShadowflame     = 232557
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=231275/blazefury-retributer
	// Adds 2 fire damage to your melee attacks.
	core.NewItemEffect(BlazefuryRetributer, func(agent core.Agent) {
		vanilla.BlazefuryTriggerAura(agent.GetCharacter(), 468169, core.SpellSchoolFire, 2)
	})
	// https://www.wowhead.com/classic/item=231862/blazefury-retributer
	core.NewItemEffect(BlazefuryRetributerBloodied, func(agent core.Agent) {
		vanilla.BlazefuryTriggerAura(agent.GetCharacter(), 468169, core.SpellSchoolFire, 2)
	})

	// https://www.wowhead.com/classic/item=230794/claw-of-chromaggus
	// Your offensive spellcasts increase the spell damage of a random school of magic by 50 for 10 sec. (9.5s cooldown)
	core.NewItemEffect(ClawOfChromaggus, func(agent core.Agent) {
		ClawOfChromaggusEffect(agent.GetCharacter())
	})
	// https://www.wowhead.com/classic/item=232557/claw-of-chromaggus
	core.NewItemEffect(ClawOfChromaggusShadowflame, func(agent core.Agent) {
		ClawOfChromaggusEffect(agent.GetCharacter())
	})

	// https://www.wowhead.com/classic/item=231016/nat-pagles-fish-terminator
	// Chance on hit: Zap nearby enemies dealing 175 to 225 damage to them. Will affect up to 4 targets.
	core.NewItemEffect(NatPaglesFishTerminator, fishTerminatorEffect)
	// https://www.wowhead.com/classic/item=231848/nat-pagles-fish-terminator
	core.NewItemEffect(NatPaglesFishTerminatorBloodied, fishTerminatorEffect)

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
			Name:             "Lightning's Cell Trigger",
			Callback:         core.CallbackOnSpellHitDealt,
			Outcome:          core.OutcomeCrit,
			ProcMask:         core.ProcMaskSpellDamage,
			CanProcFromProcs: true,
			ICD:              time.Millisecond * 2000,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				chargeAura.Activate(sim)
				chargeAura.AddStack(sim)
			},
		})
	})

	// https://www.wowhead.com/classic/item=230253/hearstriker
	// Equip: 2% chance on ranged hit to gain 1 extra attack. (Proc chance: 2%, 1s cooldown)
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
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			ProcChance:        0.02,
			ICD:               time.Second * 1,
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
		} else if character.HasRuneById(int32(proto.MageRune_RuneBeltSpellfrostBolt)) {
			arcaneChance, frostChance = 0.25, 0.25
			fireChance, natureChance, shadowChance = 0.50/3, 0.50/3, 0.50/3
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
		Name:     "Claw of the Chromaggus Trigger",
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

func fishTerminatorEffect(agent core.Agent) {
	character := agent.GetCharacter()

	results := make([]*core.SpellResult, min(4, character.Env.GetNumTargets()))

	procSpell := character.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 467836},
		// Same school and defense type as Thunder Clap
		SpellSchool:      core.SpellSchoolPhysical,
		DefenseType:      core.DefenseTypeMagic,
		ProcMask:         core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, sim.Roll(175, 225), spell.OutcomeMagicHitAndCrit)
				target = character.Env.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:              "Fish Terminator Trigger",
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		ProcMask:          core.ProcMaskMeleeMH,
		SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
		PPM:               1.50, // 1.50 PPM tested on PTR
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			procSpell.Cast(sim, result.Target)
		},
	})
}
