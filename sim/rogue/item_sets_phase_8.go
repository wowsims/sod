package rogue

import (
	"github.com/wowsims/sod/sim/core"
)

var ItemSetDuskwraithArmor = core.NewItemSet(core.ItemSet{
	Name: "Duskwraith Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyScarletEnclaveDamage2PBonus()
		},
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyScarletEnclaveDamage4PBonus()
		},
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyScarletEnclaveDamage6PBonus()
		},
	},
})

// While Just a Flesh Wound is not active, your Backstab, Sinister Strike, Saber Slash, and Mutilate deal 10% increased damage per your active Poison or Bleed effect
// afflicting the target, up to a maximum increase of 30%
func (rogue *Rogue) applyScarletEnclaveDamage2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Rogue - Damage 2P Bonus"
	if rogue.HasAura(label) {
		return
	}

	healthMetrics := rogue.NewHealthMetrics(core.ActionID{SpellID: 1219261})

	core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: ClassSpellMask_RogueDeadlyPoisonTick | ClassSpellMask_RogueOccultPoisonTick | ClassSpellMask_RogueInstantPoison,
		Callback:       core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rogue.GainHealth(sim, result.Damage*0.05, healthMetrics)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_RogueAmbush | ClassSpellMask_RogueInstantPoison,
		IntValue:  20,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_RogueBackstab,
		IntValue:  10,
	})
}

// Your Poison and autoattack critical strikes have a 10% chance to grant you a combo point. (Proc chance: 10%)
func (rogue *Rogue) applyScarletEnclaveDamage4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Rogue - Damage 4P Bonus"
	if rogue.HasAura(label) {
		return
	}

	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 1219288})

	core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
		Name:        label,
		SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolNature,
		Callback:    core.CallbackOnPeriodicDamageDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rogue.AddEnergy(sim, 1, energyMetrics)
		},
	})
}

// Increases Ambush, Eviscerate, Crimson Tempest, and Envenom damage by 50%
func (rogue *Rogue) applyScarletEnclaveDamage6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Rogue - Damage 6P Bonus"
	if rogue.HasAura(label) {
		return
	}

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_RogueAmbush | ClassSpellMask_RogueEviscerate | ClassSpellMask_RogueCrimsonTempest | ClassSpellMask_RogueEnvenom,
		IntValue:  50,
	}))
}
