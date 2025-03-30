package rogue

import (
	"fmt"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
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

// 2P
// While Just a Flesh Wound is not active, your Backstab, Sinister Strike, Saber Slash, and Mutilate deal 10% increased damage per your active Poison or Bleed effect
// afflicting the target, up to a maximum increase of 30%
func (rogue *Rogue) applyScarletEnclaveDamage2PBonus() {
	if rogue.HasRune(proto.RogueRune_RuneJustAFleshWound) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Rogue - Damage 2P Bonus"
	if rogue.HasAura(label) {
		return
	}

	totalBleedsAndPoisons := rogue.PoisonsActive + rogue.BleedsActive
	spellsModifiedBySetBonus := ClassSpellMask_RogueBackstab | ClassSpellMask_RogueSinisterStrike | ClassSpellMask_RogueSaberSlash | ClassSpellMask_RogueMutilate

	damageMod := rogue.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  spellsModifiedBySetBonus,
		FloatValue: 1.0,
	})

	core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
		ActionID:       core.ActionID{SpellID: 1226843},
		Name:           label,
		Callback:       core.CallbackOnApplyEffects,
		ProcChance:     0.1,
		ClassSpellMask: spellsModifiedBySetBonus,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			fmt.Println("Bleeds: ", rogue.BleedsActive)
			fmt.Println("Poisons: ", rogue.PoisonsActive)

			// Only apply the damage mod up to 3 times for the 30% bonus maximum
			fmt.Println("2P proc triggered")
			damageMod.UpdateFloatValue(1 + 0.10*float64(min(3, totalBleedsAndPoisons)))
			damageMod.Activate()
		},
	})
}

// 4P
// Your Poison and autoattack critical strikes have a 10% chance to grant you a combo point. (Proc chance: 10%)
func (rogue *Rogue) applyScarletEnclaveDamage4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Rogue - Damage 4P Bonus"
	if rogue.HasAura(label) {
		return
	}

	// Combo! spell that adds the combo point for the set bonus
	comboPointMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 1226869})

	// Add a combo point for poison tick crit
	core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: ClassSpellMask_RogueDeadlyPoisonTick | ClassSpellMask_RogueOccultPoisonTick | ClassSpellMask_RogueInstantPoison,
		Callback:       core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeCrit,
		ProcChance:     0.10,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rogue.AddComboPoints(sim, 1, rogue.CurrentTarget, comboPointMetrics)
		},
	})

	// Add a combo point for autoattack crit
	core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
		Name:       label,
		Callback:   core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.10,
		ProcMask:   core.ProcMaskMeleeWhiteHit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rogue.AddComboPoints(sim, 1, rogue.CurrentTarget, comboPointMetrics)
		},
	})
}

// 6P
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

// Damaging finishing moves have a 20% chance per combo point to restore 20 energy.
func (rogue *Rogue) ApplyFallenRegalityRogueBonus(aura *core.Aura) {
	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 1232184})
	aura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
		rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
			if spell.ProcMask != core.ProcMaskEmpty && sim.Proc(0.20*float64(comboPoints), "Fallen Regality Proc") {
				rogue.AddEnergy(sim, 20, energyMetrics)
			}
		})
	})
}
