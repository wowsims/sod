package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetDawnstalkerProwess = core.NewItemSet(core.ItemSet{
	Name: "Dawnstalker Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyScarletEnclaveMelee2PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyScarletEnclaveMelee4PBonus()
		},
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyScarletEnclaveMelee6PBonus()
		},
	},
})

// Your Strikes and Mongoose Bite deal 20% increased damage to targets afflicted with your Serpent Sting or Wyvern Strike.
func (hunter *Hunter) applyScarletEnclaveMelee2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Hunter - Melee 2P Bonus"
	if hunter.HasAura(label) {
		return
	}

	damageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: ClassSpellMask_HunterStrikes | ClassSpellMask_HunterMongooseBite,
	})

	core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: ClassSpellMask_HunterStrikes | ClassSpellMask_HunterMongooseBite,
		Callback:       core.CallbackOnApplyEffects,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			hasSerpentSting := hunter.SerpentSting.Dot(result.Target).IsActive()
			if hunter.SoFSerpentSting != nil {
				for _, sting := range hunter.SoFSerpentSting {
					if sting != nil {
						hasSerpentSting = hasSerpentSting || sting.Dot(result.Target).IsActive()
					}
				}
			}

			hasWyvernStrike := hunter.WyvernStrike != nil && hunter.WyvernStrike.Dot(result.Target).IsActive()
			damageMod.UpdateFloatValue(core.TernaryFloat64(hasSerpentSting || hasWyvernStrike, 1.20, 1.0))
			damageMod.Activate()
		},
	})
}

// Your melee critical strikes increase your attack speed by 20% for 10 sec.
func (hunter *Hunter) applyScarletEnclaveMelee4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Hunter - Melee 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	attackSpeed := 1.20

	procAura := core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label:    "Wicked Fast",
		ActionID: core.ActionID{SpellID: 1226357},
		Duration: time.Second * 10,
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			hunter.MultiplyMeleeSpeed(sim, attackSpeed)
		},
		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			hunter.MultiplyMeleeSpeed(sim, 1/attackSpeed)
		},
	}))

	core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
		Name:     label,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeCrit,
		ProcMask: core.ProcMaskMelee,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAura.Activate(sim)
		},
	})
}

// Increases the bonus damage from Raptor Fury by an additional 10% per stack.
func (hunter *Hunter) applyScarletEnclaveMelee6PBonus() {
	if !hunter.HasRune(proto.HunterRune_RuneBracersRaptorFury) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Hunter - Melee 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			hunter.BonusRaptorFuryDamageMultiplier = 0.10
		},
	}))
}

var ItemSetDawnstalkerArmor = core.NewItemSet(core.ItemSet{
	Name: "Dawnstalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyScarletEnclaveRanged2PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyScarletEnclaveRanged4PBonus()
		},
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyScarletEnclaveRanged6PBonus()
		},
	},
})

// Your Shots deal 25% increased damage to targets afflicted with your Serpent Sting.
func (hunter *Hunter) applyScarletEnclaveRanged2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Hunter - Ranged 2P Bonus"
	if hunter.HasAura(label) {
		return
	}

	damageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: ClassSpellMask_HunterShots,
	})

	core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: ClassSpellMask_HunterShots,
		Callback:       core.CallbackOnApplyEffects,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			damageMod.UpdateFloatValue(core.TernaryFloat64(hunter.SerpentSting.Dot(result.Target).IsActive(), 1.25, 1.0))
			damageMod.Activate()
		},
	})
}

// Your ranged critical strikes increase your Attack Power by 30% for 10 sec.
func (hunter *Hunter) applyScarletEnclaveRanged4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Hunter - Ranged 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	apBonus := hunter.NewDynamicMultiplyStat(stats.AttackPower, 1.30)
	apRangedBonus := hunter.NewDynamicMultiplyStat(stats.RangedAttackPower, 1.30)

	procAura := hunter.GetOrRegisterAura(core.Aura{
		Label:    "Wicked Shot",
		ActionID: core.ActionID{SpellID: 1226136},
		Duration: time.Second * 10,
	}).AttachStatDependency(apBonus).AttachStatDependency(apRangedBonus)

	core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
		Name:     label,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeCrit,
		ProcMask: core.ProcMaskRanged,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAura.Activate(sim)
		},
	})
}

// Your Multi-Shot hits 2 additional targets, and your Kill Shot and Chimera Shot hits increase the damage done by your next Multi-Shot by 100%.
func (hunter *Hunter) applyScarletEnclaveRanged6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Hunter - Ranged 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	damMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterMultiShot,
	})

	multishotAura := hunter.RegisterAura(core.Aura{
		Label:     "Trick Shots",
		ActionID:  core.ActionID{SpellID: 1233451},
		Duration:  time.Minute * 5,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			damMod.UpdateIntValue(int64(100 * newStacks))
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_HunterMultiShot) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			hunter.MultiShotBonusTargets = 2
		},
	})).AttachProcTrigger(core.ProcTrigger{
		Name:             label,
		ClassSpellMask:   ClassSpellMask_HunterChimeraShot | ClassSpellMask_HunterKillShot,
		Callback:         core.CallbackOnApplyEffects,
		CanProcFromProcs: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			multishotAura.Activate(sim)
			multishotAura.AddStack(sim)
		},
	})
}

// Flanking Strike's damage buff is increased by an additional 2% per stack. When striking from behind, your target takes 150% increased damage from Flanking Strike.
func (hunter *Hunter) ApplyFallenRegalityHunterBonus(aura *core.Aura) {
	if !hunter.HasRune(proto.HunterRune_RuneLegsFlankingStrike) {
		return
	}

	flankingBuffDamageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:     core.SpellMod_DamageDone_Flat,
		ProcMask: core.ProcMaskMelee,
	})

	if !hunter.PseudoStats.InFrontOfTarget {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  ClassSpellMask_HunterFlankingStrike,
			FloatValue: 2.50,
		})

		if hunter.pet != nil {
			hunter.pet.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  ClassSpellMask_HunterPetFlankingStrike,
				FloatValue: 2.50,
			})
		}
	}

	aura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
		hunter.FlankingStrike.RelatedSelfBuff.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			flankingBuffDamageMod.Activate()
		}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
			flankingBuffDamageMod.Deactivate()
		}).ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			flankingBuffDamageMod.UpdateIntValue(int64(2 * newStacks))
		})
	})
}
