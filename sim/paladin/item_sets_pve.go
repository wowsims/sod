package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Soul Related Set Bonus IDs
const (
	PaladinT1Prot2P  = 456536
	PaladinT1Prot4P  = 456538
	PaladinT1Prot6P  = 456541
	PaladinT2Prot2P  = 467531
	PaladinT2Prot4P  = 467532
	PaladinT2Prot6P  = 467536
	PaladinTAQProt2P = 1213410
	PaladinTAQProt4P = 1213413
	PaladinT1Holy2P  = 456488
	PaladinT1Holy4P  = 457323
	PaladinT1Holy6P  = 456492
	PaladinT2Holy2P  = 467506
	PaladinT2Holy4P  = 467507
	PaladinT2Holy6P  = 467513
	PaladinTAQHoly2P = 1213349
	PaladinTAQHoly4P = 1213353
	PaladinT1Ret2P   = 456494
	PaladinT1Ret4P   = 456489
	PaladinT1Ret6P   = 456533
	PaladinT2Ret2P   = 467518
	PaladinT2Ret4P   = 467526
	PaladinT2Ret6P   = 467529
	PaladinT3Ret2P   = 1219189
	PaladinT3Ret4P   = 1219191
	PaladinT3Ret6P   = 1219193
	PaladinTAQRet2P  = 1213397
	PaladinTAQRet4P  = 1213406
	PaladinZG2P      = 468401
	PaladinZG3P      = 468428
	PaladinZG5P      = 468431
	PaladinRAQ3P     = 1213467
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

	if duplicateBonusCheckAndCreate(paladin, PaladinT1Prot2P, bonusLabel) {
		return
	}

	// (2) Set: Increases the block value of your shield by 30.
	paladin.AddStat(stats.BlockValue, 30)
}

func (paladin *Paladin) applyPaladinT1Prot4P() {
	bonusLabel := "S03 - Item - T1 - Paladin - Protection 4P Bonus"

	if paladin.HasAura(bonusLabel) {
		return
	}

	// (4) Set: Heal for 189 to 211 when you Block. (ICD: 3.5s)
	// Note: The heal does not scale with healing/spell power, but can crit.
	actionID := core.ActionID{SpellID: 456540}

	bastionOfLight := paladin.RegisterSpell(core.SpellConfig{
		ActionID:         actionID,
		SpellSchool:      core.SpellSchoolHoly,
		DefenseType:      core.DefenseTypeMagic,
		ProcMask:         core.ProcMaskSpellHealing,
		Flags:            core.SpellFlagHelpful,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHeal := sim.Roll(189, 211)
			spell.CalcAndDealHealing(sim, target, baseHeal, spell.OutcomeHealingCrit)
		},
	})

	handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		bastionOfLight.Cast(sim, result.Target)
	}

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:       bonusLabel,
		ActionID:   core.ActionID{SpellID: PaladinT1Prot4P},
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeBlock,
		ProcChance: 1.0,
		ICD:        time.Millisecond * 3500,
		Handler:    handler,
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
			auras := paladin.holyShieldAura
			procs := paladin.holyShieldProc
			blockBonus := 30.0 * core.BlockRatingPerBlockChance

			for i, values := range HolyShieldValues {

				if paladin.Level < values.level {
					break
				}

				damage := values.damage

				// Holy Shield's damage is increased by 80% of shield block value.
				procs[i].ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					sbv := paladin.BlockValue() * 0.8
					// Reminder: Holy Shield can crit, but does not miss.
					spell.CalcAndDealDamage(sim, target, (damage + sbv), spell.OutcomeMagicCrit)
				}

				// Holy Shield aura no longer has stacks...
				auras[i].MaxStacks = 0

				// ...and does not set stacks on gain...
				auras[i].OnGain = func(aura *core.Aura, sim *core.Simulation) {
					paladin.AddStatDynamic(sim, stats.Block, blockBonus)
				}

				// ...or remove stacks on block.
				auras[i].OnSpellHitTaken = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidBlock() {
						procs[i].Cast(sim, spell.Unit)
					}
				}
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

	blockBonus := 40.0 * core.BlockRatingPerBlockChance

	paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT2Prot2P},
		OnInit: func(_ *core.Aura, _ *core.Simulation) {
			for i, hsAura := range paladin.holyShieldAura {
				if paladin.Level < HolyShieldValues[i].level {
					break
				}
				oldOnGain := hsAura.OnGain
				hsAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
					oldOnGain(aura, sim)
					paladin.AddStatDynamic(sim, stats.Block, blockBonus)
				}
				oldOnExpire := hsAura.OnExpire
				hsAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
					oldOnExpire(aura, sim)
					paladin.AddStatDynamic(sim, stats.Block, -blockBonus)
				}
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

	paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT2Prot4P},
		OnInit: func(_ *core.Aura, _ *core.Simulation) {
			for i, hsAura := range paladin.holyShieldAura {
				if hsAura == nil || paladin.Level < HolyShieldValues[i].level {
					break
				}
				oldOnGain := hsAura.OnGain
				oldOnExpire := hsAura.OnExpire

				hsAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
					oldOnGain(aura, sim)
					paladin.PseudoStats.DamageTakenMultiplier *= 0.9
				}
				hsAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
					oldOnExpire(aura, sim)
					paladin.PseudoStats.DamageTakenMultiplier /= 0.9
				}
			}
		},
	})
}

func (paladin *Paladin) applyPaladinT2Prot6P() {
	bonusLabel := "S03 - Item - T2 - Paladin - Protection 6P Bonus"

	if paladin.HasAura(bonusLabel) {
		return
	}

	// Your Reckoning Talent now has a 20% chance per talent point to trigger when
	// you block.
	if paladin.Talents.Reckoning == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 20178} // Reckoning proc ID
	procChance := 0.2 * float64(paladin.Talents.Reckoning)

	handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		paladin.AutoAttacks.ExtraMHAttack(sim, 1, actionID, spell.ActionID)
	}

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:       bonusLabel,
		ActionID:   core.ActionID{SpellID: PaladinT2Prot6P},
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeBlock,
		ProcChance: procChance,
		Handler:    handler,
	})
}

func (paladin *Paladin) applyPaladinTAQProt2P() {
	// Empty Function (Not Implemented)
}

func (paladin *Paladin) applyPaladinTAQProt4P() {
	// Empty Function (Not Implemented)
}

func (paladin *Paladin) applyPaladinT1Holy2P() {
	//(Not Implemented)
}

func (paladin *Paladin) applyPaladinT1Holy4P() {

	bonusLabel := "S03 - Item - T1 - Paladin - Holy 4P Bonus"

	if duplicateBonusCheckAndCreate(paladin, PaladinT1Holy4P, bonusLabel) {
		return
	}

	paladin.AddStat(stats.MeleeCrit, 2)
	paladin.AddStat(stats.SpellCrit, 2)
}

func (paladin *Paladin) applyPaladinT1Holy6P() {
	//(Not Implemented)
}

func (paladin *Paladin) applyPaladinT2Holy2P() {

	bonusLabel := "S03 - Item - T2 - Paladin - Holy 2P Bonus"

	if duplicateBonusCheckAndCreate(paladin, PaladinT2Holy2P, bonusLabel) {
		return
	}

	//Increases critical strike chance of holy shock spell by 5%
	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_PaladinHolyShock) {
			spell.BonusCritRating += 5.0
		}
	})
}

func (paladin *Paladin) applyPaladinT2Holy4P() {

	bonusLabel := "S03 - Item - T2 - Paladin - Holy 4P Bonus"

	if duplicateBonusCheckAndCreate(paladin, PaladinT2Holy4P, bonusLabel) {
		return
	}

	//Increases damage done by your Consecration spell by 50%
	paladin.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinConsecration,
		IntValue:  50,
	})
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

	if duplicateBonusCheckAndCreate(paladin, PaladinT1Ret4P, bonusLabel) {
		return
	}

	paladin.AddStat(stats.MeleeCrit, 2)
	paladin.AddStat(stats.SpellCrit, 2)
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
	t2Judgement6pcAura := paladin.GetOrRegisterAura(core.Aura{
		Label:     "Swift Judgement",
		ActionID:  core.ActionID{SpellID: 467530},
		Duration:  time.Second * 8,
		MaxStacks: 5,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= (1.0 + (float64(oldStacks) * 0.01))
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= (1.0 + (float64(newStacks) * 0.01))
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    bonusLabel,
		ActionID: core.ActionID{SpellID: PaladinT2Ret6P},
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			originalApplyEffects := paladin.judgement.ApplyEffects

			// Wrap the apply Judgement ApplyEffects with more Effects
			paladin.judgement.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				originalApplyEffects(sim, target, spell)
				// 6 pieces: Your Judgement grants 1% increased Holy damage for 8 sec, stacking up to 5 times.
				t2Judgement6pcAura.Activate(sim)
				t2Judgement6pcAura.AddStack(sim)
			}
		},
	})
}

func (paladin *Paladin) applyPaladinTAQRet2P() {

	bonusLabel := "S03 - Item - TAQ - Paladin - Retribution 2P Bonus"

	if duplicateBonusCheckAndCreate(paladin, PaladinTAQRet2P, bonusLabel) {
		return
	}

	if !paladin.hasRune(proto.PaladinRune_RuneHandsCrusaderStrike) {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinCrusaderStrike,
		IntValue:  50,
	})
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
			damageMod.UpdateIntValue(int64(40 * newStacks))
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
	// No soul provides this bonus
	paladin.AddStats(stats.Stats{
		stats.HolyPower: 14,
	})
}

func (paladin *Paladin) applyPaladinZG3P() {
	bonusLabel := "S03 - Item - ZG - Paladin - Caster 3P Bonus"

	if duplicateBonusCheckAndCreate(paladin, PaladinZG3P, bonusLabel) {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinHolyShock,
		IntValue:  50,
	})
}

func (paladin *Paladin) applyPaladinZG5P() {
	bonusLabel := "S03 - Item - ZG - Paladin - Caster 5P Bonus"

	if paladin.HasAura(bonusLabel) {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinExorcism,
		IntValue:  50,
	})

	paladin.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_PaladinExorcism,
		TimeValue: -time.Second * 3,
	})
}

func (paladin *Paladin) applyPaladinRAQ3P() {
	bonusLabel := "S03 - Item - RAQ - Paladin - Retribution 3P Bonus"

	if paladin.HasAura(bonusLabel) {
		return
	}

	aura := core.Aura{
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
	}

	paladin.RegisterAura(aura)
}

func duplicateBonusCheckAndCreate(agent core.Agent, bonusID int32, bonusString string) bool {
	paladin := agent.(PaladinAgent).GetPaladin()

	if paladin.HasAura(bonusString) {
		return true // Do not apply bonus aura more than once (Due to Soul)
	}

	paladin.RegisterAura(core.Aura{
		Label:    bonusString,
		ActionID: core.ActionID{SpellID: bonusID},
	})

	return false
}
