package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// var ItemSetLightbreakersWarplate = core.NewItemSet(core.ItemSet{
// 	Name: "Lightbreaker's Warplate",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			warrior := agent.(WarriorAgent).GetWarrior()
// 			warrior.applyScarletEnclaveDamage2PBonus()
// 		},
// 		4: func(agent core.Agent) {
// 			warrior := agent.(WarriorAgent).GetWarrior()
// 			warrior.applyScarletEnclaveDamage4PBonus()
// 		},
// 		6: func(agent core.Agent) {
// 			warrior := agent.(WarriorAgent).GetWarrior()
// 			warrior.applyScarletEnclaveDamage6PBonus()
// 		},
// 	},
// })

// Your Cleave strikes 1 additional target.
func (warrior *Warrior) applyScarletEnclaveDamage2PBonus() {
	if warrior.Env.GetNumTargets() < 3 {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Warrior - Damage 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warrior.cleaveTargetCount += 1
		},
	}))
}

// Your Mortal Strike, Bloodthirst, and Shield Slam deal 50% more damage.
func (warrior *Warrior) applyScarletEnclaveDamage4PBonus() {
	if !warrior.Talents.Bloodthirst && !warrior.Talents.MortalStrike && !warrior.Talents.ShieldSlam {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Warrior - Damage 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_WarriorMortalStrike | ClassSpellMask_WarriorBloodthirst | ClassSpellMask_WarriorShieldSlam,
		IntValue:  50,
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
				warrior.Whirlwind.ModifyRemainingCooldown(sim, -3*time.Second)
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

// var ItemSetLightbreakersBattlegear = core.NewItemSet(core.ItemSet{
// 	Name: "Lightbreaker's Battlegear",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			warrior := agent.(WarriorAgent).GetWarrior()
// 			warrior.applyScarletEnclaveProtection2PBonus()
// 		},
// 		4: func(agent core.Agent) {
// 			warrior := agent.(WarriorAgent).GetWarrior()
// 			warrior.applyScarletEnclaveProtection4PBonus()
// 		},
// 		6: func(agent core.Agent) {
// 			warrior := agent.(WarriorAgent).GetWarrior()
// 			warrior.applyScarletEnclaveProtection6PBonus()
// 		},
// 	},
// })

// Your Shockwave deals 100% increased damage.
func (warrior *Warrior) applyScarletEnclaveProtection2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Warrior - Protection 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_WarriorShockwave,
		IntValue:  100,
	})
}

// Increases the duration of your Recklessness by 15 sec.
func (warrior *Warrior) applyScarletEnclaveProtection4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Warrior - Protection 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_BuffDuration_Flat,
		ClassMask: ClassSpellMask_WarriorRecklesness,
		TimeValue: time.Second * 15,
	})
}

// Gladiator Stance no longer reduces your Armor or Threat.
func (warrior *Warrior) applyScarletEnclaveProtection6PBonus() {
	if !warrior.HasRune(proto.WarriorRune_RuneGladiatorStance) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Warrior - Protection 6P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warrior.gladiatorStanceThreatMultiplier = 1
			warrior.gladiatorStanceArmorMultiplier = 1
		},
	}))
}
