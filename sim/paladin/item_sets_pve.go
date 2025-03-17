package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Soul Related Set Bonus IDs
const (
	// Prot
	PaladinT1Prot2P  = 456536
	PaladinT1Prot4P  = 456538
	PaladinT1Prot6P  = 456541
	PaladinT2Prot2P  = 467531
	PaladinT2Prot4P  = 467532
	PaladinT2Prot6P  = 467536
	PaladinT3Prot2P  = 1219205
	PaladinT3Prot4P  = 1219225
	PaladinT3Prot6P  = 1219226
	PaladinTAQProt2P = 1213410
	PaladinTAQProt4P = 1213413
	PaladinTSEProt2P = 1226467
	PaladinTSEProt4P = 1226477
	PaladinTSEProt6P = 1226479

	// Holy
	PaladinT1Holy2P  = 456488
	PaladinT1Holy4P  = 457323
	PaladinT1Holy6P  = 456492
	PaladinT2Holy2P  = 467506
	PaladinT2Holy4P  = 467507
	PaladinT2Holy6P  = 467513
	PaladinTAQHoly2P = 1213349
	PaladinTAQHoly4P = 1213353
	PaladinTSEHoly2P = 1226452
	PaladinTSEHoly4P = 1226454
	PaladinTSEHoly6P = 1226459

	// Ret
	PaladinT1Ret2P  = 456494
	PaladinT1Ret4P  = 456489
	PaladinT1Ret6P  = 456533
	PaladinT2Ret2P  = 467518
	PaladinT2Ret4P  = 467526
	PaladinT2Ret6P  = 467529
	PaladinT3Ret2P  = 1219189
	PaladinT3Ret4P  = 1219191
	PaladinT3Ret6P  = 1219193
	PaladinTAQRet2P = 1213397
	PaladinTAQRet4P = 1213406
	PaladinTSERet2P = 1226460
	PaladinTSERet4P = 1226462
	PaladinTSERet6P = 1226463

	// Other
	PaladinZG2P  = 468401
	PaladinZG3P  = 468428
	PaladinZG5P  = 468431
	PaladinRAQ3P = 1213467
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 3 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetObsessedProphetsPlate = core.NewItemSet(core.ItemSet{
	Name: "Obsessed Prophet's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance)
			c.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.PseudoStats.SchoolBonusCritChance[stats.SchoolIndexHoly] += 3 * core.SpellCritRatingPerCritChance
		},
	},
})

var _ = core.NewItemSet(core.ItemSet{
	Name: "Emerald Encrusted Battleplate",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 22)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetSoulforgeArmor = core.NewItemSet(core.ItemSet{
	Name: "Soulforge Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power and up to 40 increased healing from spells.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
				stats.HealingPower:      40,
			})
		},
		// 6% chance on melee autoattack and 4% chance on spellcast to increase your damage and healing done by magical spells and effects by up to 95 for 10 sec.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450625}

			procAura := c.NewTemporaryStatsAura("Crusader's Wrath", core.ActionID{SpellID: 27499}, stats.Stats{stats.SpellPower: 95}, time.Second*10)
			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Item - Crusader's Wrath Proc - Lightforge Armor (Melee Auto)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMeleeWhiteHit,
				ProcChance: 0.06,
				Handler:    handler,
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Item - Crusader's Wrath Proc - Lightforge Armor (Spell Cast)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.04,
				Handler:    handler,
			})
		},
		// +8 All Resistances.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// +200 Armor.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
	},
})

var ItemSetLawbringerMercy = core.NewItemSet(core.ItemSet{
	Name: "Lawbringer Mercy",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT1Holy2P()
		},
		4: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT1Holy4P()
		},
		6: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT1Holy6P()
		},
	},
})

var ItemSetLawbringerRadiance = core.NewItemSet(core.ItemSet{
	Name: "Lawbringer Radiance",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT1Ret2P()
		},
		4: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT1Ret4P()
		},
		6: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT1Ret6P()
		},
	},
})

var ItemSetLawbringerWill = core.NewItemSet(core.ItemSet{
	Name: "Lawbringer Will",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT1Prot2P()
		},
		4: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT1Prot4P()
		},
		6: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT1Prot6P()
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 5 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetFreethinkersArmor = core.NewItemSet(core.ItemSet{
	Name: "Freethinker's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) { agent.(PaladinAgent).GetPaladin().applyPaladinZG2P() },
		3: func(agent core.Agent) { agent.(PaladinAgent).GetPaladin().applyPaladinZG3P() },
		5: func(agent core.Agent) { agent.(PaladinAgent).GetPaladin().applyPaladinZG5P() },
	},
})

var ItemSetMercifulJudgement = core.NewItemSet(core.ItemSet{
	Name: "Merciful Judgement",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT2Holy2P()
		},
		4: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT2Holy4P()
		},
		6: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT2Holy6P()
		},
	},
})

var ItemSetRadiantJudgement = core.NewItemSet(core.ItemSet{
	Name: "Radiant Judgement",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT2Ret2P()
		},
		4: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT2Ret4P()
		},
		6: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT2Ret6P()
		},
	},
})

var ItemSetWilfullJudgement = core.NewItemSet(core.ItemSet{
	Name: "Wilfull Judgement",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT2Prot2P()
		},
		4: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT2Prot4P()
		},
		6: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinT2Prot6P()
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 6 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetAvengersWill = core.NewItemSet(core.ItemSet{
	Name: "Avenger's Will",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinTAQProt2P()
		},
		4: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinTAQProt4P()
		},
	},
})

var ItemSetAvengersRadiance = core.NewItemSet(core.ItemSet{
	Name: "Avenger's Radiance",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinTAQRet2P()
		},
		4: func(agent core.Agent) {
			agent.(PaladinAgent).GetPaladin().applyPaladinTAQRet4P()
		},
	},
})

var ItemSetBattlegearOfEternalJustice = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Eternal Justice",
	Bonuses: map[int32]core.ApplyEffect{
		// Crusader Strike now unleashes the judgement effect of your seals, but does not consume the seal
		3: func(agent core.Agent) { agent.(PaladinAgent).GetPaladin().applyPaladinRAQ3P() },
	},
})

func (paladin *Paladin) applyPaladinT1Prot2P() {
	bonusLabel := "S03 - Item - T1 - Paladin - Protection 2P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	// (2) Set: Increases the block value of your shield by 30.
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT1Prot2P},
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachBuildPhaseStatBuff(stats.BlockValue, 30))
}

func (paladin *Paladin) applyPaladinT1Prot4P() {
	bonusLabel := "S03 - Item - T1 - Paladin - Protection 4P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	// (4) Set: Heal for 189 to 211 when you Block. (ICD: 3.5s)
	// Note: The heal does not scale with healing/spell power, but can crit.
	bastionOfLight := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 456540},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHeal := sim.Roll(189, 211)
			spell.CalcAndDealHealing(sim, target, baseHeal, spell.OutcomeHealingCrit)
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:       bonusLabel,
		ActionID:   core.ActionID{SpellID: PaladinT1Prot4P},
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeBlock,
		ProcChance: 1.0,
		ICD:        time.Millisecond * 3500,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bastionOfLight.Cast(sim, result.Target)
		},
	})
}

func (paladin *Paladin) applyPaladinT1Prot6P() {
	bonusLabel := "S03 - Item - T1 - Paladin - Protection 6P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT1Prot6P},
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for i, values := range HolyShieldValues {
				if paladin.Level < values.level {
					break
				}

				// Holy Shield's damage is increased by 80% of shield block value.
				paladin.holyShieldExtraDamage = func(_ *core.Simulation, paladin *Paladin) float64 {
					return paladin.BlockValue() * 0.8
				}

				// Holy Shield aura no longer has stacks and does not set stacks on gain or remove stacks on block.
				// Setting MaxStacks to 0 disables this behavior in holy_shield.go
				paladin.holyShieldAura[i].MaxStacks = 0
			}
		},
	})
}

func (paladin *Paladin) applyPaladinT2Prot2P() {
	bonusLabel := "S03 - Item - T2 - Paladin - Protection 2P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	//Increases the bonus chance to block from Holy Shield by 10%
	if !paladin.Talents.HolyShield {
		return
	}

	blockBonus := 10.0 * core.BlockRatingPerBlockChance

	paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT2Prot2P},
		OnInit: func(_ *core.Aura, _ *core.Simulation) {
			for i, hsAura := range paladin.holyShieldAura {
				if paladin.Level < HolyShieldValues[i].level {
					break
				}
				hsAura.AttachStatBuff(stats.Block, blockBonus)
			}
		},
	})
}

func (paladin *Paladin) applyPaladinT2Prot4P() {
	bonusLabel := "S03 - Item - T2 - Paladin - Protection 4P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	//You take 10% reduced damage while Holy Shield is active.
	if !paladin.Talents.HolyShield {
		return
	}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT2Prot4P},
		OnInit: func(_ *core.Aura, _ *core.Simulation) {
			for i, hsAura := range paladin.holyShieldAura {
				if hsAura == nil || paladin.Level < HolyShieldValues[i].level {
					break
				}

				hsAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
					paladin.PseudoStats.DamageTakenMultiplier *= 0.9
				}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
					paladin.PseudoStats.DamageTakenMultiplier /= 0.9
				})
			}
		},
	}))
}

func (paladin *Paladin) applyPaladinT2Prot6P() {
	bonusLabel := "S03 - Item - T2 - Paladin - Protection 6P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	// Your Reckoning Talent now has a 20% chance per talent point to trigger when you block.
	if paladin.Talents.Reckoning == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 20178} // Reckoning proc ID
	procChance := 0.2 * float64(paladin.Talents.Reckoning)

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:       bonusLabel,
		ActionID:   core.ActionID{SpellID: PaladinT2Prot6P},
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeBlock,
		ProcChance: procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.AutoAttacks.ExtraMHAttack(sim, 1, actionID, spell.ActionID)
		},
	})
}

func (paladin *Paladin) applyPaladinTAQProt2P() {
	bonusLabel := "S03 - Item - TAQ - Paladin - Protection 2P Bonus"
	if paladin.HasAura(bonusLabel) || paladin.Options.PersonalBlessing != proto.Blessings_BlessingOfSanctuary {
		return
	}

	statDeps := []*stats.StatDependency{
		paladin.NewDynamicMultiplyStat(stats.Stamina, 1.10),
		paladin.NewDynamicMultiplyStat(stats.Agility, 1.10),
		paladin.NewDynamicMultiplyStat(stats.Strength, 1.10),
		paladin.NewDynamicMultiplyStat(stats.Intellect, 1.10),
		paladin.NewDynamicMultiplyStat(stats.Spirit, 1.10),
	}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		ActionID:   core.ActionID{SpellID: PaladinTAQProt2P},
		Label:      bonusLabel,
		BuildPhase: core.CharacterBuildPhaseBuffs,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, dep := range statDeps {
				aura.Unit.EnableBuildPhaseStatDep(sim, dep)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, dep := range statDeps {
				aura.Unit.DisableBuildPhaseStatDep(sim, dep)
			}
		},
	}))
}

func (paladin *Paladin) applyPaladinTAQProt4P() {
	// Empty Function (Not Implemented)
}

func (paladin *Paladin) applyPaladinT1Holy2P() {
	//(Not Implemented)
}

func (paladin *Paladin) applyPaladinT1Holy4P() {
	bonusLabel := "S03 - Item - T1 - Paladin - Holy 4P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	bonusStats := stats.Stats{
		stats.SpellCrit: 2 * core.SpellCritRatingPerCritChance,
		stats.MeleeCrit: 2 * core.CritRatingPerCritChance,
	}
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT1Holy4P},
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachBuildPhaseStatsBuff(bonusStats))
}

func (paladin *Paladin) applyPaladinT1Holy6P() {
	//(Not Implemented)
}

func (paladin *Paladin) applyPaladinT2Holy2P() {
	bonusLabel := "S03 - Item - T2 - Paladin - Holy 2P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	//Increases critical strike chance of holy shock spell by 5%
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT2Holy2P},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Flat,
		ClassMask:  ClassSpellMask_PaladinHolyShock,
		FloatValue: 5.0 * core.SpellCritRatingPerCritChance,
	}))
}

func (paladin *Paladin) applyPaladinT2Holy4P() {
	bonusLabel := "S03 - Item - T2 - Paladin - Holy 4P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	//Increases damage done by your Consecration spell by 50%
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT2Holy4P},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinConsecration,
		IntValue:  50,
	}))
}

func (paladin *Paladin) applyPaladinT2Holy6P() {
	// While you are not your Beacon of Light target, your Beacon of Light target is also healed by 100% of the damage you deal
	// with Consecration, Exorcism, Holy Shock, Holy Wrath, and Hammer of Wrath
	// No need to Sim
}

func (paladin *Paladin) applyPaladinTAQHoly2P() {
	//(Not Implemented)
}

func (paladin *Paladin) applyPaladinTAQHoly4P() {
	//(Not Implemented)
}

func (paladin *Paladin) applyPaladinT1Ret2P() {
	// No need to model
	//(2) Set : Your Judgement of Light and Judgement of Wisdom also grant the effects of Judgement of the Crusader.
}

func (paladin *Paladin) applyPaladinT1Ret4P() {
	bonusLabel := "S03 - Item - T1 - Paladin - Retribution 4P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	bonusStats := stats.Stats{
		stats.SpellCrit: 2 * core.SpellCritRatingPerCritChance,
		stats.MeleeCrit: 2 * core.CritRatingPerCritChance,
	}
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT1Ret4P},
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachBuildPhaseStatsBuff(bonusStats))
}

func (paladin *Paladin) applyPaladinT1Ret6P() {
	bonusLabel := "S03 - Item - T1 - Paladin - Retribution 6P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT1Ret6P},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			paladin.lingerDuration = time.Second * 6
			paladin.enableMultiJudge = true // Implemented in Paladin.go
		},
	}))
}

func (paladin *Paladin) applyPaladinT2Ret2P() {
	bonusLabel := "S03 - Item - T2 - Paladin - Retribution 2P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	// 2 pieces: Increases damage done by your damaging Judgements by 20% and your Judgements no longer consume your Seals on the target.
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT2Ret2P},
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// Made to enable double judge 2025-02-01
			paladin.enableMultiJudge = true // Implemented in Paladin.go
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.consumeSealsOnJudge = false
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.consumeSealsOnJudge = true
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinJudgements,
		IntValue:  20,
	}))
}

func (paladin *Paladin) applyPaladinT2Ret4P() {
	bonusLabel := "S03 - Item - T2 - Paladin - Retribution 4P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	// 4 pieces: Reduces the cooldown on your Judgement ability by 5 seconds.
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT2Ret4P},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_PaladinJudgement,
		TimeValue: -5 * time.Second,
	}))
}

func (paladin *Paladin) applyPaladinT2Ret6P() {
	bonusLabel := "S03 - Item - T2 - Paladin - Retribution 6P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	// 6 pieces: Your Judgement grants 1% increased Holy damage for 8 sec, stacking up to 5 times.
	swiftJudgementAura := paladin.GetOrRegisterAura(core.Aura{
		Label:     "Swift Judgement",
		ActionID:  core.ActionID{SpellID: 467530},
		Duration:  time.Second * 8,
		MaxStacks: 5,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= (1.0 + (float64(oldStacks) * 0.01))
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= (1.0 + (float64(newStacks) * 0.01))
		},
	})

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT2Ret6P},
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			originalApplyEffects := paladin.judgement.ApplyEffects

			// Wrap the apply Judgement ApplyEffects with more Effects
			paladin.judgement.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				originalApplyEffects(sim, target, spell)

				// 6 pieces: Your Judgement grants 1% increased Holy damage for 8 sec, stacking up to 5 times.
				swiftJudgementAura.Activate(sim)
				swiftJudgementAura.AddStack(sim)
			}
		},
	}))
}

func (paladin *Paladin) applyPaladinTAQRet2P() {
	bonusLabel := "S03 - Item - TAQ - Paladin - Retribution 2P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinTAQRet2P},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinCrusaderStrike,
		IntValue:  50,
	}))
}

func (paladin *Paladin) applyPaladinTAQRet4P() {
	bonusLabel := "S03 - Item - TAQ - Paladin - Retribution 4P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	damageMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinExorcism,
	})

	buffAura := paladin.GetOrRegisterAura(core.Aura{
		Label:     "Excommunication",
		ActionID:  core.ActionID{SpellID: 1217927},
		Duration:  time.Second * 20,
		MaxStacks: 3,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			damageMod.UpdateIntValue(core.TernaryInt64(paladin.MainHand().HandType == proto.HandType_HandTypeTwoHand, 35, 0) * int64(newStacks))
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
	})

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinTAQRet4P},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			} else if spell.Matches(ClassSpellMask_PaladinExorcism) {
				buffAura.Deactivate(sim)
			}
		},
	}))
}

func (paladin *Paladin) applyPaladinZG2P() {
	bonusLabel := "S03 - Item - ZG - Paladin - Caster 2P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	// Increases damage done by Holy spells and effects by up to 14.
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinZG2P},
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachBuildPhaseStatBuff(stats.HolyPower, 14))
}

func (paladin *Paladin) applyPaladinZG3P() {
	bonusLabel := "S03 - Item - ZG - Paladin - Caster 3P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	// Increases damage done by your Holy Shock spell by 50%.
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinZG3P},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinHolyShock,
		IntValue:  50,
	}))
}

func (paladin *Paladin) applyPaladinZG5P() {
	bonusLabel := "S03 - Item - ZG - Paladin - Caster 5P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	// Reduces the cooldown of your Exorcism spell by 3 sec and increases its damage done by 50%.
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinZG3P},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinExorcism,
		IntValue:  50,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_PaladinExorcism,
		TimeValue: -time.Second * 3,
	}))
}

func (paladin *Paladin) applyPaladinRAQ3P() {
	bonusLabel := "S03 - Item - RAQ - Paladin - Retribution 3P Bonus"
	if paladin.HasAura(bonusLabel) {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinRAQ3P},
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if !paladin.hasRune(proto.PaladinRune_RuneHandsCrusaderStrike) {
				return
			}

			originalApplyEffects := paladin.crusaderStrike.ApplyEffects
			extraApplyEffects := paladin.judgement.ApplyEffects

			// Wrap the apply Crusader Strike ApplyEffects with more Effects
			paladin.crusaderStrike.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				originalApplyEffects(sim, target, spell)
				consumeSealsOnJudgeSaved := paladin.consumeSealsOnJudge // Save current value
				paladin.consumeSealsOnJudge = false                     // Set to not consume seals
				if paladin.currentSeal.IsActive() {
					extraApplyEffects(sim, target, paladin.judgement)
				}
				paladin.consumeSealsOnJudge = consumeSealsOnJudgeSaved // Restore saved value
			}
		},
	})
}
