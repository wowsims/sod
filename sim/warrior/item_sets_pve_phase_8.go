package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetLightbreakersWarplate = core.NewItemSet(core.ItemSet{
	Name: "Lightbreaker's Warplate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyScarletEnclaveDamage2PBonus()
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyScarletEnclaveDamage4PBonus()
		},
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyScarletEnclaveDamage6PBonus()
		},
	},
})

// Increases Heroic Strike, Cleave, and Quick Strike damage by 20%.
// Your Cleave strikes 1 additional target and can trigger Blood Surge.
func (warrior *Warrior) applyScarletEnclaveDamage2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Warrior - Damage 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warrior.CleaveTargetCount += 1
			warrior.bloodSurgeClassMask |= ClassSpellMask_WarriorCleave
		},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_WarriorHeroicStrike | ClassSpellMask_WarriorCleave | ClassSpellMask_WarriorQuickStrike,
		IntValue:  20,
	})
}

// Each time you hit a target with Whirlwind, Heroic Strike, Quick Strike, or Cleave, the damage of your next Slam is increased by 20%, stacking up to 5 times.
// If you are wielding a two-handed weapon, you will gain 2 stacks each time.
func (warrior *Warrior) applyScarletEnclaveDamage4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Warrior - Damage 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	buffClassMask := ClassSpellMask_WarriorSlamMH | ClassSpellMask_WarriorSlamOH

	damageMod := warrior.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: buffClassMask,
	})

	buffAura := warrior.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1227232},
		Label:     label + " Stacking Buff", // TODO: Find real buff
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			damageMod.UpdateIntValue(20 * int64(newStacks))
			damageMod.Activate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(buffClassMask) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           label,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: ClassSpellMask_WarriorWhirlwindMH | ClassSpellMask_WarriorWhirlwindOH | ClassSpellMask_WarriorHeroicStrike | ClassSpellMask_WarriorQuickStrike | ClassSpellMask_WarriorCleave,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buffAura.Activate(sim)
			buffAura.AddStack(sim)

			if warrior.MainHand().HandType == proto.HandType_HandTypeTwoHand {
				buffAura.AddStack(sim)
			}
		},
	})
}

// Each time Deep Wounds deals damage, it reduces the remaining cooldown on your Whirlwind by 3 sec.
// Whirlwind deals 50% increased damage to targets afflicted with your Deep Wounds.
func (warrior *Warrior) applyScarletEnclaveDamage6PBonus() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Warrior - Damage 6P Bonus"
	if warrior.HasAura(label) {
		return
	}

	damageMod := warrior.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  ClassSpellMask_WarriorWhirlwind,
		FloatValue: 1.0,
	})

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_WarriorDeepWounds) {
				warrior.Whirlwind.CD.ModifyRemainingCooldown(sim, -3*time.Second)
			}
		},
		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_WarriorWhirlwind) {
				damageMod.UpdateFloatValue(core.TernaryFloat64(warrior.DeepWounds.Dot(target).IsActive(), 1.5, 1.0))
				damageMod.Activate()
			}
		},
	}))
}

var ItemSetLightbreakersBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Lightbreaker's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyScarletEnclaveProtection2PBonus()
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyScarletEnclaveProtection4PBonus()
		},
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyScarletEnclaveProtection6PBonus()
		},
	},
})

// Your Shockwave deals 100% increased damage and its cooldown is reduced by 2.0 sec each time you hit a target with Heroic Strike, Quick Strike, or Cleave.
func (warrior *Warrior) applyScarletEnclaveProtection2PBonus() {
	if !warrior.HasRune(proto.WarriorRune_RuneShockwave) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Warrior - Protection 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           label,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: ClassSpellMask_WarriorHeroicStrike | ClassSpellMask_WarriorQuickStrike | ClassSpellMask_WarriorCleave,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.Shockwave.CD.ModifyRemainingCooldown(sim, -2*time.Second)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_WarriorShockwave,
		IntValue:  100,
	})
}

// Your Recklessness, Retaliation, and Shield Wall abilities no longer share a cooldown.
// Additionally, your Recklessness ability lasts 15 sec longer, and while it is active you gain 50% of your Defense Skill over 300 as Strength.
func (warrior *Warrior) applyScarletEnclaveProtection4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Warrior - Protection 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	var snapshottedDefense float64
	buffAura := warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1227242}, // TODO: Find real spell
		Label:    label + " Strength buff",
		Duration: DefaultRecklessnessDuration + time.Second*15,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			snapshottedDefense = 0
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			snapshottedDefense = warrior.GetStat(stats.Defense) * 0.50
			warrior.AddStatDynamic(sim, stats.Strength, snapshottedDefense)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Strength, snapshottedDefense)
		},
	})

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Recklessness.SharedCD.Duration = 0
			warrior.ShieldWall.SharedCD.Duration = 0
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_WarriorRecklesness) {
				buffAura.Activate(sim)
			}
		},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_BuffDuration_Flat,
		ClassMask: ClassSpellMask_WarriorRecklesness,
		TimeValue: time.Second * 15,
	})
}

// Gladiator Stance no longer reduces your Armor or Threat, and instead increases threat by 30%.
// In addition, each time your Revenge, Devastate, or Shield Slam hits, the damage done by your next Whirlwind or Execute is increased by 20%, stacking up to 5 times.
func (warrior *Warrior) applyScarletEnclaveProtection6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Warrior - Protection 6P Bonus"
	if warrior.HasAura(label) {
		return
	}

	classMask := ClassSpellMask_WarriorWhirlwindMH | ClassSpellMask_WarriorWhirlwindOH | ClassSpellMask_WarriorExecute

	damageMod := warrior.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: classMask,
	})

	buffAura := warrior.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1227245}, // TODO: Find real spell
		Label:     label + " Stacking Buff",
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			damageMod.UpdateIntValue(20 * int64(newStacks))
			damageMod.Activate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(classMask) {
				aura.Deactivate(sim)
			}
		},
	})

	aura := core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           label,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: ClassSpellMask_WarriorRevenge | ClassSpellMask_WarriorDevastate | ClassSpellMask_WarriorShieldSlam,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buffAura.Activate(sim)
			buffAura.AddStack(sim)
		},
	})

	if warrior.HasRune(proto.WarriorRune_RuneGladiatorStance) {
		aura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
			warrior.gladiatorStanceThreatMultiplier = 1.30
			warrior.gladiatorStanceArmorMultiplier = 1
		})
	}
}
