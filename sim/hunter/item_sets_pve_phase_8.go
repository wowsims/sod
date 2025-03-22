package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// TODO: UNCOMMENT WHEN TIER SETS ARE ADDED TO THE SIM
// var ItemSetDawnstalkerProwess = core.NewItemSet(core.ItemSet{
// 	Name: "Dawnstalker Prowess",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			hunter := agent.(HunterAgent).GetHunter()
// 			hunter.applyScarletEnclaveMelee2PBonus()
// 		},
// 		4: func(agent core.Agent) {
// 			hunter := agent.(HunterAgent).GetHunter()
// 			hunter.applyScarletEnclaveMelee4PBonus()
// 		},
// 		6: func(agent core.Agent) {
// 			hunter := agent.(HunterAgent).GetHunter()
// 			hunter.applyScarletEnclaveMelee6PBonus()
// 		},
// 	},
// })

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
			hasDebuff := hunter.SerpentSting.Dot(result.Target).IsActive() || hunter.WyvernStrike.Dot(result.Target).IsActive()
			damageMod.UpdateFloatValue(core.TernaryFloat64(hasDebuff, 0.20, 0.0))
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
	label := "S03 - Item - Scarlet Enclave - Hunter - Melee 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	if !hunter.HasRune(proto.HunterRune_RuneBracersRaptorFury) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			hunter.BonusRaptorFuryDamageMultiplier = 0.10
		},
	}))
}

// TODO: UNCOMMENT WHEN TIER SETS ARE ADDED TO THE SIM
// var ItemSetDawnstalkerArmor = core.NewItemSet(core.ItemSet{
// 	Name: "Dawnstalker Armor",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			hunter := agent.(HunterAgent).GetHunter()
// 			hunter.applyScarletEnclaveRanged2PBonus()
// 		},
// 		4: func(agent core.Agent) {
// 			hunter := agent.(HunterAgent).GetHunter()
// 			hunter.applyScarletEnclaveRanged4PBonus()
// 		},
// 		6: func(agent core.Agent) {
// 			hunter := agent.(HunterAgent).GetHunter()
// 			hunter.applyScarletEnclaveRanged6PBonus()
// 		},
// 	},
// })

// Your Shots deal 20% increased damage to targets afflicted with your Serpent Sting.
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
			damageMod.UpdateIntValue(core.TernaryInt64(hunter.SerpentSting.Dot(result.Target).IsActive(), 20, 0))
			damageMod.Activate()
		},
	})
}

// Reduces the cooldown on your Chimera Shot, Explosive Shot, and Aimed Shot abilities by 1.5 sec and reduces the cooldown on your Kill Shot ability by 3sec.
func (hunter *Hunter) applyScarletEnclaveRanged4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Hunter - Ranged 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	apBonus := hunter.NewDynamicMultiplyStat(stats.AttackPower, 1.2)
	apRangedBonus := hunter.NewDynamicMultiplyStat(stats.RangedAttackPower, 1.2)

	procAura := hunter.GetOrRegisterAura(core.Aura{
		Label:    "Wicked Shot",
		ActionID: core.ActionID{SpellID: 1226136},
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.EnableDynamicStatDep(sim, apBonus)
			hunter.EnableDynamicStatDep(sim, apRangedBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.DisableDynamicStatDep(sim, apBonus)
			hunter.DisableDynamicStatDep(sim, apRangedBonus)
		},
	})

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

	damageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterMultiShot,
		IntValue:  100,
	})

	multishotAura := hunter.RegisterAura(core.Aura{
		Label:    "Multi-Shot +100% Damage",
		Duration: time.Minute * 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
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
		Name:           label,
		ClassSpellMask: ClassSpellMask_HunterChimeraShot | ClassSpellMask_HunterKillShot,
		Callback:       core.CallbackOnApplyEffects,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			multishotAura.Activate(sim)
		},
	})
}
