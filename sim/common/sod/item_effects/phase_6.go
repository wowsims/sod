package item_effects

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/sod/sim/common/sod"
	"github.com/wowsims/sod/sim/common/vanilla"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	RazorspikeBattleplate    = 233492
	LeggingsOfImmersion      = 233505
	RingOfSwarmingThought    = 233507
	RobesOfTheBattleguard    = 233575
	RazorspikeShoulderplates = 233793
	RazorspikeHeadcage       = 233795
	RazorbrambleShoulderpads = 233804
	RazorbrambleCowl         = 233808
	RazorbrambleLeathers     = 233813
	VampiricCowl             = 233826
	VampiricShawl            = 233833
	VampiricRobe             = 233837
	TunedForceReactiveDisk   = 233988
	LodestoneOfRetaliation   = 233992

	// Obsidian Weapons
	ObsidianChampion    = 233490
	ObsidianReaver      = 233491
	ObsidianDestroyer   = 233796
	ObsidianStormhammer = 233797
	ObsidianSageblade   = 233798
	ObsidianDefender    = 233801
	ObsidianHeartseeker = 234428
	ObsidianShotgun     = 234434

	// Brood of Nozdormu Reputations Rings
	SignetRingBronzeDominatorR5   = 234034
	SignetRingBronzeDominatorR4   = 234030
	SignetRingBronzeDominatorR3   = 234026
	SignetRingBronzeDominatorR2   = 234021
	SignetRingBronzeDominatorR1   = 234017
	SignetRingBronzeInvokerR5     = 234032
	SignetRingBronzeInvokerR4     = 234028
	SignetRingBronzeInvokerR3     = 234024
	SignetRingBronzeInvokerR2     = 234020
	SignetRingBronzeInvokerR1     = 234016
	SignetRingBronzeFlamekeeperR5 = 234964
	SignetRingBronzeFlamekeeperR4 = 234965
	SignetRingBronzeFlamekeeperR3 = 234966
	SignetRingBronzeFlamekeeperR2 = 234967
	SignetRingBronzeFlamekeeperR1 = 234968
	SignetRingBronzePreserverR5   = 234033
	SignetRingBronzePreserverR4   = 234029
	SignetRingBronzePreserverR3   = 234025
	SignetRingBronzePreserverR2   = 234023
	SignetRingBronzePreserverR1   = 234019
	SignetRingBronzeProtectorR5   = 234035
	SignetRingBronzeProtectorR4   = 234031
	SignetRingBronzeProtectorR3   = 234027
	SignetRingBronzeProtectorR2   = 234022
	SignetRingBronzeProtectorR1   = 234018
	SignetRingBronzeSubjugatorR5  = 234436
	SignetRingBronzeSubjugatorR4  = 234437
	SignetRingBronzeSubjugatorR3  = 234438
	SignetRingBronzeSubjugatorR2  = 234439
	SignetRingBronzeSubjugatorR1  = 234440
	SignetRingBronzeConquerorR5   = 234202
	SignetRingBronzeConquerorR4   = 234201
	SignetRingBronzeConquerorR3   = 234200
	SignetRingBronzeConquerorR2   = 234199
	SignetRingBronzeConquerorR1   = 234198
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=233490/obsidian-champion
	// Chance on hit: Heal self for 270 to 450 and Increases Strength by 120 for 30 sec.
	// Chance on hit: Increases damage done by 20 and attack speed by 5% for 15 sec.
	// Equip: Chance on direct spell damage or melee attack to deal 325 to 500 Fire damage to your target.
	// 		  Deals extra damage to Silithid creatures in Silithus on almost every hit.
	// 		  (2.1s cooldown)
	core.NewItemEffect(ObsidianChampion, func(agent core.Agent) {
		character := agent.GetCharacter()

		ObsidianEdgedAura(ObsidianChampion, agent)

		strengthAura := vanilla.StrengthOfTheChampionAura(character)
		enrageAura := vanilla.EnrageAura446327(character)

		ppm := 1.0 // Estimated based on data from WoW Armaments Discord
		strDpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(ObsidianChampion, ppm, 0)
		enrageDpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(ObsidianChampion, ppm, 0)

		strengthProcTrigger := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Obsidian Champion (Strength)",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			DPM:               strDpm,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				strengthAura.Activate(sim)
			},
		})

		enrageProcTrigger := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Obsidian Champion (Enrage)",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			DPM:               enrageDpm,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				enrageAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterProc(ObsidianChampion, strengthProcTrigger)
		character.ItemSwap.RegisterProc(ObsidianChampion, enrageProcTrigger)

	})

	// https://www.wowhead.com/classic/item=233801/obsidian-defender
	// Equip: Chance on direct spell damage or melee attack to deal 325 to 500 Fire damage to your target.
	// 		  Deals extra damage to Silithid creatures in Silithus on almost every hit.
	// 		  (2.1s cooldown)
	core.NewItemEffect(ObsidianDefender, func(agent core.Agent) {
		ObsidianEdgedAura(ObsidianDefender, agent)
	})

	// https://www.wowhead.com/classic/item=233796/obsidian-destroyer
	// Chance on hit: Stuns target for 3 sec.
	// Equip: Chance on direct spell damage or melee attack to deal 325 to 500 Fire damage to your target.
	// 		  Deals extra damage to Silithid creatures in Silithus on almost every hit.
	// 		  (2.1s cooldown)
	core.NewItemEffect(ObsidianDestroyer, func(agent core.Agent) {
		ObsidianEdgedAura(ObsidianDestroyer, agent)
	})

	// https://www.wowhead.com/classic/item=234428/obsidian-heartseeker
	// Equip: Chance on direct spell damage or melee attack to deal 325 to 500 Fire damage to your target.
	// 		  Deals extra damage to Silithid creatures in Silithus on almost every hit.
	// 		  (2.1s cooldown)
	core.NewItemEffect(ObsidianHeartseeker, func(agent core.Agent) {
		ObsidianEdgedAura(ObsidianHeartseeker, agent)
	})

	// https://www.wowhead.com/classic/item=233491/obsidian-reaver
	// Equip: Chance on direct spell damage or melee attack to deal 325 to 500 Fire damage to your target.
	// 		  Deals extra damage to Silithid creatures in Silithus on almost every hit.
	// 		  (2.1s cooldown)
	core.NewItemEffect(ObsidianReaver, func(agent core.Agent) {
		ObsidianEdgedAura(ObsidianReaver, agent)
	})

	// https://www.wowhead.com/classic/item=233798/obsidian-sageblade
	// Equip: Chance on direct spell damage or melee attack to deal 325 to 500 Fire damage to your target.
	// 		  Deals extra damage to Silithid creatures in Silithus on almost every hit.
	// 		  (2.1s cooldown)
	core.NewItemEffect(ObsidianSageblade, func(agent core.Agent) {
		ObsidianEdgedAura(ObsidianSageblade, agent)
	})

	core.NewItemEffect(ObsidianShotgun, func(agent core.Agent) {
		ObsidianEdgedAura(ObsidianShotgun, agent)
	})

	// https://www.wowhead.com/classic/item=233797/obsidian-stormhammer
	// Chance on hit: Blasts up to 3 targets for 105 to 145 Nature damage.
	// Equip: Chance on direct spell damage or melee attack to deal 325 to 500 Fire damage to your target.
	// 		  Deals extra damage to Silithid creatures in Silithus on almost every hit.
	// 		  (2.1s cooldown)
	core.NewItemEffect(ObsidianStormhammer, func(agent core.Agent) {
		sod.StormhammerChainLightningProcAura(agent, ObsidianStormhammer)
		ObsidianEdgedAura(ObsidianStormhammer, agent)
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=233992/lodestone-of-retaliation
	// Equip: When struck in combat inflicts 80 Nature damage to the attacker.
	// Causes twice as much threat as damage dealt.
	core.NewItemEffect(LodestoneOfRetaliation, func(agent core.Agent) {
		thornsNatureDamageEffect(agent, LodestoneOfRetaliation, "Lodestone of Retaliation", 80)
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Rings
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=233507/ring-of-swarming-thought
	// Your damaging non-periodic spells increase your spell damage by 15 for 15 sec.
	// If the target is player controlled, gain 100 spell penetration for 15 sec instead.
	// (15s cooldown)
	core.NewItemEffect(RingOfSwarmingThought, func(agent core.Agent) {
		character := agent.GetCharacter()
		buffAura := character.NewTemporaryStatsAura("Swarming Thoughts", core.ActionID{SpellID: 474148}, stats.Stats{stats.SpellDamage: 15}, time.Second*15)

		aura := core.MakeProcTriggerAura(&agent.GetCharacter().Unit, core.ProcTrigger{
			Name:             "Swarming Thoughts Trigger",
			Callback:         core.CallbackOnSpellHitDealt,
			Outcome:          core.OutcomeLanded,
			ProcMask:         core.ProcMaskSpellDamage,
			CanProcFromProcs: true,
			ProcChance:       1.00,
			ICD:              time.Second * 15,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterProc(RingOfSwarmingThought, aura)
	})

	// https://www.wowhead.com/classic/item=234198/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeConquerorR5, func(agent core.Agent) {
		TimewornStrikeAura(agent, SignetRingBronzeConquerorR5)
	})
	core.NewItemEffect(SignetRingBronzeConquerorR4, func(agent core.Agent) {
		TimewornStrikeAura(agent, SignetRingBronzeConquerorR4)
	})
	core.NewItemEffect(SignetRingBronzeConquerorR3, func(agent core.Agent) {
		TimewornStrikeAura(agent, SignetRingBronzeConquerorR3)
	})
	core.NewItemEffect(SignetRingBronzeConquerorR2, func(agent core.Agent) {
		TimewornStrikeAura(agent, SignetRingBronzeConquerorR2)
	})
	core.NewItemEffect(SignetRingBronzeConquerorR1, func(agent core.Agent) {
		TimewornStrikeAura(agent, SignetRingBronzeConquerorR1)
	})

	// https://www.wowhead.com/classic/item=234034/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeDominatorR5, func(agent core.Agent) {
		TimewornStrikeAura(agent, SignetRingBronzeDominatorR5)
	})
	core.NewItemEffect(SignetRingBronzeDominatorR4, func(agent core.Agent) {
		TimewornStrikeAura(agent, SignetRingBronzeDominatorR4)
	})
	core.NewItemEffect(SignetRingBronzeDominatorR3, func(agent core.Agent) {
		TimewornStrikeAura(agent, SignetRingBronzeDominatorR3)
	})
	core.NewItemEffect(SignetRingBronzeDominatorR2, func(agent core.Agent) {
		TimewornStrikeAura(agent, SignetRingBronzeDominatorR2)
	})
	core.NewItemEffect(SignetRingBronzeDominatorR1, func(agent core.Agent) {
		TimewornStrikeAura(agent, SignetRingBronzeDominatorR1)
	})

	// https://www.wowhead.com/classic/item=234964/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeFlamekeeperR5, func(agent core.Agent) {
		TimewornPyromancyAura(agent, SignetRingBronzeFlamekeeperR5)
	})
	core.NewItemEffect(SignetRingBronzeFlamekeeperR4, func(agent core.Agent) {
		TimewornPyromancyAura(agent, SignetRingBronzeFlamekeeperR4)
	})
	core.NewItemEffect(SignetRingBronzeFlamekeeperR3, func(agent core.Agent) {
		TimewornPyromancyAura(agent, SignetRingBronzeFlamekeeperR3)
	})
	core.NewItemEffect(SignetRingBronzeFlamekeeperR2, func(agent core.Agent) {
		TimewornPyromancyAura(agent, SignetRingBronzeFlamekeeperR2)
	})
	core.NewItemEffect(SignetRingBronzeFlamekeeperR1, func(agent core.Agent) {
		TimewornPyromancyAura(agent, SignetRingBronzeFlamekeeperR1)
	})

	// https://www.wowhead.com/classic/item=234032/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeInvokerR5, func(agent core.Agent) {
		TimewornSpellAura(agent, SignetRingBronzeInvokerR5)
	})
	core.NewItemEffect(SignetRingBronzeInvokerR4, func(agent core.Agent) {
		TimewornSpellAura(agent, SignetRingBronzeInvokerR4)
	})
	core.NewItemEffect(SignetRingBronzeInvokerR3, func(agent core.Agent) {
		TimewornSpellAura(agent, SignetRingBronzeInvokerR3)
	})
	core.NewItemEffect(SignetRingBronzeInvokerR2, func(agent core.Agent) {
		TimewornSpellAura(agent, SignetRingBronzeInvokerR2)
	})
	core.NewItemEffect(SignetRingBronzeInvokerR1, func(agent core.Agent) {
		TimewornSpellAura(agent, SignetRingBronzeInvokerR1)
	})

	// https://www.wowhead.com/classic/item=234032/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzePreserverR5, func(agent core.Agent) {
		TimewornHealing(agent, SignetRingBronzePreserverR5)
	})
	core.NewItemEffect(SignetRingBronzePreserverR4, func(agent core.Agent) {
		TimewornHealing(agent, SignetRingBronzePreserverR4)
	})
	core.NewItemEffect(SignetRingBronzePreserverR3, func(agent core.Agent) {
		TimewornHealing(agent, SignetRingBronzePreserverR3)
	})
	core.NewItemEffect(SignetRingBronzePreserverR2, func(agent core.Agent) {
		TimewornHealing(agent, SignetRingBronzePreserverR2)
	})
	core.NewItemEffect(SignetRingBronzePreserverR1, func(agent core.Agent) {
		TimewornHealing(agent, SignetRingBronzePreserverR1)
	})

	// https://www.wowhead.com/classic/item=234035/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeProtectorR5, func(agent core.Agent) {
		TimewornExpertiseAura(agent, SignetRingBronzeProtectorR5)
	})
	core.NewItemEffect(SignetRingBronzeProtectorR4, func(agent core.Agent) {
		TimewornExpertiseAura(agent, SignetRingBronzeProtectorR4)
	})
	core.NewItemEffect(SignetRingBronzeProtectorR3, func(agent core.Agent) {
		TimewornExpertiseAura(agent, SignetRingBronzeProtectorR3)
	})
	core.NewItemEffect(SignetRingBronzeProtectorR2, func(agent core.Agent) {
		TimewornExpertiseAura(agent, SignetRingBronzeProtectorR2)
	})
	core.NewItemEffect(SignetRingBronzeProtectorR1, func(agent core.Agent) {
		TimewornExpertiseAura(agent, SignetRingBronzeProtectorR1)
	})

	// https://www.wowhead.com/classic/item=234436/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeSubjugatorR5, func(agent core.Agent) {
		TimewornDecayAura(agent, SignetRingBronzeSubjugatorR5)
	})
	core.NewItemEffect(SignetRingBronzeSubjugatorR4, func(agent core.Agent) {
		TimewornDecayAura(agent, SignetRingBronzeSubjugatorR4)
	})
	core.NewItemEffect(SignetRingBronzeSubjugatorR3, func(agent core.Agent) {
		TimewornDecayAura(agent, SignetRingBronzeSubjugatorR3)
	})
	core.NewItemEffect(SignetRingBronzeSubjugatorR2, func(agent core.Agent) {
		TimewornDecayAura(agent, SignetRingBronzeSubjugatorR2)
	})
	core.NewItemEffect(SignetRingBronzeSubjugatorR1, func(agent core.Agent) {
		TimewornDecayAura(agent, SignetRingBronzeSubjugatorR1)
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Other
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=233505/leggings-of-immersion
	// Your damaging non-periodic Nature spell critical strikes have a chance to grant you 40 increased spell damage for 10 sec.
	// (Proc chance: 50%, 5s cooldown)
	core.NewItemEffect(LeggingsOfImmersion, func(agent core.Agent) {
		character := agent.GetCharacter()
		buffAura := character.NewTemporaryStatsAura("Immersed", core.ActionID{SpellID: 474126}, stats.Stats{stats.SpellDamage: 40}, time.Second*10)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 5,
		}

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Immersed Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellDamage,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.SpellSchool.Matches(core.SpellSchoolNature) && result.DidCrit() && icd.IsReady(sim) && sim.Proc(0.50, "Immersed Proc") {
					buffAura.Activate(sim)
					icd.Use(sim)
				}
			},
		})

		character.ItemSwap.RegisterProc(LeggingsOfImmersion, triggerAura)
	})

	// https://www.wowhead.com/classic/item=233808/razorbramble-cowl
	// Equip: When struck in combat inflicts 100 Nature damage to the attacker.
	// Causes twice as much threat as damage dealt.
	core.NewItemEffect(RazorbrambleCowl, func(agent core.Agent) {
		thornsNatureDamageEffect(agent, RazorbrambleCowl, "Razorbramble Cowl", 100)
	})

	// https://www.wowhead.com/classic/item=233813/razorbramble-leathers
	// Equip: When struck in combat inflicts 100 Nature damage to the attacker.
	// Causes twice as much threat as damage dealt.
	core.NewItemEffect(RazorbrambleLeathers, func(agent core.Agent) {
		thornsNatureDamageEffect(agent, RazorbrambleLeathers, "Razorbramble Leathers", 100)
	})

	// https://www.wowhead.com/classic/item=233804/razorbramble-shoulderpads
	// Equip: When struck in combat inflicts 80 Nature damage to the attacker.
	// Causes twice as much threat as damage dealt.
	core.NewItemEffect(RazorbrambleShoulderpads, func(agent core.Agent) {
		thornsNatureDamageEffect(agent, RazorbrambleShoulderpads, "Razorbramble Shoulderpads", 80)
	})

	// https://www.wowhead.com/classic/item=233492/razorspike-battleplate
	// Equip: When struck in combat inflicts 100 Nature damage to the attacker.
	// Causes twice as much threat as damage dealt.
	core.NewItemEffect(RazorspikeBattleplate, func(agent core.Agent) {
		thornsNatureDamageEffect(agent, RazorspikeBattleplate, "Razorspike Battleplate", 100)
	})

	// https://www.wowhead.com/classic/item=233795/razorspike-headcage
	// Equip: When struck in combat inflicts 100 Nature damage to the attacker.
	// Causes twice as much threat as damage dealt.
	core.NewItemEffect(RazorspikeHeadcage, func(agent core.Agent) {
		thornsNatureDamageEffect(agent, RazorspikeHeadcage, "Razorspike Headcage", 100)
	})

	// https://www.wowhead.com/classic/item=233793/razorspike-shoulderplates
	// Equip: When struck in combat inflicts 80 Nature damage to the attacker.
	// Causes twice as much threat as damage dealt.
	core.NewItemEffect(RazorspikeShoulderplates, func(agent core.Agent) {
		thornsNatureDamageEffect(agent, RazorspikeShoulderplates, "Razorspike Shoulderplates", 80)
	})

	// https://www.wowhead.com/classic/item=233575/robes-of-the-battleguard
	// Your damaging non-periodic spells increase your spell damage by 20 for 15 sec.
	// If the target is player controlled, gain 120 spell penetration for 15 sec instead.
	// (15s cooldown)
	core.NewItemEffect(RobesOfTheBattleguard, func(agent core.Agent) {
		character := agent.GetCharacter()
		buffAura := character.NewTemporaryStatsAura("Resolve of the Battleguard", core.ActionID{SpellID: 1213241}, stats.Stats{stats.SpellDamage: 20}, time.Second*15)

		triggerAura := core.MakeProcTriggerAura(&agent.GetCharacter().Unit, core.ProcTrigger{
			Name:             "Resolve of the Battleguard Trigger",
			Callback:         core.CallbackOnSpellHitDealt,
			Outcome:          core.OutcomeLanded,
			ProcMask:         core.ProcMaskSpellDamage,
			CanProcFromProcs: true,
			ProcChance:       1.00,
			ICD:              time.Second * 15,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterProc(RobesOfTheBattleguard, triggerAura)
	})

	// https://www.wowhead.com/classic/item=233988/tuned-force-reactive-disk
	// Equip: When the shield blocks it releases an electrical charge that damages all nearby enemies. (1s cooldown)
	// Use: Charge up the energy within the shield for 3 sec to deal 450 to 750 Nature damage to all nearby enemies. After use, the shield needs 10 sec to recharge. (2 Min Cooldown)
	core.NewItemEffect(TunedForceReactiveDisk, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: TunedForceReactiveDisk},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, 35, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		procTriggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Force Reactive Disk",
			Callback: core.CallbackOnSpellHitTaken,
			ProcMask: core.ProcMaskMelee,
			Outcome:  core.OutcomeBlock,
			ICD:      time.Second,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, spell.Unit)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1213967},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				sim.AddPendingAction(&core.PendingAction{
					NextActionAt: sim.CurrentTime + time.Second*3,
					OnAction: func(sim *core.Simulation) {
						for _, aoeTarget := range sim.Encounter.TargetUnits {
							spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(450, 750), spell.OutcomeMagicHitAndCrit)
						}

						spell.CD.Use(sim)

						procTriggerAura.Deactivate(sim)
						sim.AddPendingAction(&core.PendingAction{
							NextActionAt: sim.CurrentTime + time.Second*10,
							OnAction: func(sim *core.Simulation) {
								procTriggerAura.Activate(sim)
							},
						})
					},
				})
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	// https://www.wowhead.com/classic/item=233826/vampiric-cowl
	// Equip: When struck in combat has a 20% chance of stealing 50 life from target enemy. (Proc chance: 20%)
	core.NewItemEffect(VampiricCowl, func(agent core.Agent) {
		vampiricDrainLifeEffect(agent, VampiricCowl, "Vampiric Cowl", 50)
	})

	// https://www.wowhead.com/classic/item=233833/vampiric-shawl
	// Equip: When struck in combat has a 20% chance of stealing 40 life from target enemy. (Proc chance: 20%)
	core.NewItemEffect(VampiricShawl, func(agent core.Agent) {
		vampiricDrainLifeEffect(agent, VampiricShawl, "Vampiric Shawl", 40)
	})

	// https://www.wowhead.com/classic/item=233837/vampiric-robe
	// Equip: When struck in combat has a 20% chance of stealing 50 life from target enemy. (Proc chance: 20%)
	core.NewItemEffect(VampiricRobe, func(agent core.Agent) {
		vampiricDrainLifeEffect(agent, VampiricRobe, "Vampiric Robe", 50)
	})

	core.AddEffectsToTest = true
}

// https://www.wowhead.com/classic/spell=1214125/obsidian-edged
// Chance on direct spell damage or melee attack to deal 325 to 500 Fire damage to your target.
// Deals extra damage to Silithid creatures in Silithus on almost every hit.
// (2.1s cooldown)
func ObsidianEdgedAura(itemID int32, agent core.Agent) {
	character := agent.GetCharacter()

	damageSpell := character.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1213941},
		SpellSchool: core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sim.Roll(325, 500), spell.OutcomeMagicHitAndCrit)
		},
	})

	dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(itemID, 0, 0.05)

	meleeProcTrigger := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:     "Obsidian Edged Proc Physical",
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		DPM:      dpm,
		ICD:      time.Millisecond * 2100,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			damageSpell.Cast(sim, result.Target)
		},
	})

	spellProcTrigger := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       "Obsidian Edged Proc Spell",
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskSpellDamage,
		ProcChance: 0.05,
		ICD:        time.Millisecond * 2100,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			damageSpell.Cast(sim, result.Target)
		},
	})

	character.ItemSwap.RegisterProc(itemID, meleeProcTrigger)
	character.ItemSwap.RegisterProc(itemID, spellProcTrigger)
}

// https://www.wowhead.com/classic/spell=1214155/timeworn-decay
// Increases the damage dealt by all of your damage over time spells by 3% per piece of Timeworn armor equipped.
func TimewornDecayAura(agent core.Agent, itemID int32) {
	character := agent.GetCharacter()
	label := "Timeworn Decay Aura"

	if character.PseudoStats.TimewornBonus == 0 || character.HasAura(label) {
		return
	}

	aura := core.MakePermanent(character.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PeriodicDamageDone_Flat,
		FloatValue: 0.03 * float64(character.PseudoStats.TimewornBonus),
	}))

	character.ItemSwap.RegisterProc(itemID, aura)
}

// https://www.wowhead.com/classic/spell=1213407/timeworn-expertise
// Reduces the chance for your attacks to be dodged or parried by 1% per piece of Timeworn armor equipped.
func TimewornExpertiseAura(agent core.Agent, itemID int32) {
	character := agent.GetCharacter()
	label := "Timeworn Expertise Aura"
	if character.PseudoStats.TimewornBonus == 0 || character.HasAura(label) {
		return
	}

	stats := stats.Stats{stats.Expertise: float64(character.PseudoStats.TimewornBonus) * core.ExpertiseRatingPerExpertiseChance}

	aura := core.MakePermanent(character.RegisterAura(core.Aura{
		ActionID:   core.ActionID{SpellID: 1214218},
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachStatsBuff(stats))

	character.ItemSwap.RegisterProc(itemID, aura)
}

// https://www.wowhead.com/classic/spell=1213405/timeworn-healing
// Increases the effectiveness of your healing and shielding spells by 2% per piece of Timeworn armor equipped.
func TimewornHealing(agent core.Agent, itemID int32) {
	character := agent.GetCharacter()
	label := "Timeworn Healing Aura"
	if character.PseudoStats.TimewornBonus == 0 || character.HasAura(label) {
		return
	}

	healShieldMultiplier := 1 + 0.02*float64(character.PseudoStats.TimewornBonus)

	aura := core.MakePermanent(character.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213405},
		Label:    label,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			character.PseudoStats.HealingDealtMultiplier *= healShieldMultiplier
			character.PseudoStats.ShieldDealtMultiplier *= healShieldMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			character.PseudoStats.HealingDealtMultiplier /= healShieldMultiplier
			character.PseudoStats.ShieldDealtMultiplier /= healShieldMultiplier
		},
	}))

	character.ItemSwap.RegisterProc(itemID, aura)
}

// https://www.wowhead.com/classic/spell=1215404/timeworn-pyromancy
// While Metamorphosis or Way of Earth is active, increases the effectiveness of your Fire damage spells by 3% per piece of Timeworn armor equipped.
func TimewornPyromancyAura(agent core.Agent, itemID int32) {
	character := agent.GetCharacter()
	label := "Timeworn Pyromancy Aura"
	if character.PseudoStats.TimewornBonus == 0 || character.HasAura(label) {
		return
	}

	// Just applying this rune if the user has Meta or WoE
	if !character.HasRuneById(int32(proto.WarlockRune_RuneHandsMetamorphosis)) && !character.HasRuneById(int32(proto.ShamanRune_RuneLegsWayOfEarth)) {
		return
	}

	fireMultiplier := 1 + 0.03*float64(character.PseudoStats.TimewornBonus)

	aura := core.MakePermanent(character.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1215404},
		Label:    label,
	}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire], fireMultiplier))

	character.ItemSwap.RegisterProc(itemID, aura)
}

// https://www.wowhead.com/classic/spell=1213398/timeworn-spell
// Increases the casting speed of your spells by 2% per piece of Timeworn armor equipped.
func TimewornSpellAura(agent core.Agent, itemID int32) {
	character := agent.GetCharacter()
	label := "Timeworn Spell Aura"
	if character.PseudoStats.TimewornBonus == 0 || character.HasAura(label) {
		return
	}

	castSpeedMultiplier := math.Pow((1 + 0.02), float64(character.PseudoStats.TimewornBonus))

	aura := core.MakePermanent(character.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213398},
		Label:    label,
	}).AttachMultiplyCastSpeed(&character.Unit, castSpeedMultiplier))

	character.ItemSwap.RegisterProc(itemID, aura)
}

// https://www.wowhead.com/classic/spell=1213390/timeworn-strike
// Gives you 1% chance per piece of Timeworn armor equipped to get an extra attack on regular melee or ranged hit that deals 100% weapon damage.
// (100ms cooldown)
func TimewornStrikeAura(agent core.Agent, itemID int32) {
	character := agent.GetCharacter()
	if character.PseudoStats.TimewornBonus == 0 {
		return
	}

	procChance := float64(character.PseudoStats.TimewornBonus) * 0.01

	meleeAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       "Timeworn Strike Aura Melee",
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskMelee,
		ProcChance: procChance,
		ICD:        time.Millisecond * 200,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			character.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 1213381}, spell)
		},
	})

	rangedAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       "Timeworn Strike Aura Ranged",
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskRanged,
		ProcChance: procChance,
		ICD:        time.Millisecond * 200,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Extra ranged attacks do not reset the swing timer confirmed by Zirene
			character.AutoAttacks.StoreExtraRangedAttack(sim, 1, core.ActionID{SpellID: 1213381}, spell.ActionID)
		},
	})

	character.ItemSwap.RegisterProc(itemID, meleeAura)
	character.ItemSwap.RegisterProc(itemID, rangedAura)
}

func thornsNatureDamageEffect(agent core.Agent, itemID int32, itemName string, damage float64) {
	character := agent.GetCharacter()
	character.PseudoStats.ThornsDamage += damage

	procSpell := character.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{ItemID: itemID},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHit)
		},
	})

	aura := core.MakePermanent(character.GetOrRegisterAura(core.Aura{
		Label: fmt.Sprintf("Damage Shield Dmg +%f (%s)", damage, itemName),
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	}))

	character.ItemSwap.RegisterProc(itemID, aura)
}

func vampiricDrainLifeEffect(agent core.Agent, itemID int32, itemName string, damage float64) {
	character := agent.GetCharacter()

	actionID := core.ActionID{ItemID: itemID}
	healthMetrics := character.NewHealthMetrics(actionID)
	drainLifeSpell := character.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeAlwaysHit)
			character.GainHealth(sim, result.Damage, healthMetrics)
		},
	})

	triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       fmt.Sprintf("Drain Life Trigger (%s)", itemName),
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskMelee,
		ProcChance: 0.20,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			drainLifeSpell.Cast(sim, spell.Unit)
		},
	})

	character.ItemSwap.RegisterProc(itemID, triggerAura)
}
