package rogue

import (
	"time"

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
// While Just a Flesh Wound is not active, your Backstab, Sinister Strike, Saber Slash, and Mutilate deal 20% increased damage per your active Poison or Bleed effect
// afflicting the target, up to a maximum increase of 60%
func (rogue *Rogue) applyScarletEnclaveDamage2PBonus() {
	if rogue.HasRune(proto.RogueRune_RuneJustAFleshWound) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Rogue - Damage 2P Bonus"
	if rogue.HasAura(label) {
		return
	}

	spellsModifiedBySetBonus := ClassSpellMask_RogueBackstab | ClassSpellMask_RogueSinisterStrike | ClassSpellMask_RogueSaberSlash | ClassSpellMask_RogueMutilateHit

	damageMod := rogue.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  spellsModifiedBySetBonus,
		FloatValue: 1.0,
	})

	core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
		ActionID:       core.ActionID{SpellID: 1226843},
		Name:           label,
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: spellsModifiedBySetBonus,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			totalBleedsAndPoisons := rogue.PoisonsActive[rogue.CurrentTarget.UnitIndex] + rogue.BleedsActive[rogue.CurrentTarget.UnitIndex]

			// Only apply the damage mod up to 3 times for the 60% bonus maximum
			damageMod.UpdateFloatValue(1 + 0.20*float64(min(3, totalBleedsAndPoisons)))
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
	comboPointMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 1226868})

	procClassMasks := ClassSpellMask_RogueDeadlyPoisonTick | ClassSpellMask_RogueOccultPoisonTick | ClassSpellMask_RogueInstantPoison

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label:    label,
		ActionID: core.ActionID{SpellID: 1226869},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() && spell.Matches(procClassMasks) && sim.Proc(0.10, "Combo! proc") {
				rogue.AddComboPoints(sim, 1, rogue.CurrentTarget, comboPointMetrics)
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() && (spell.Matches(procClassMasks) || spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit)) && sim.Proc(0.10, "Combo! proc") {
				rogue.AddComboPoints(sim, 1, rogue.CurrentTarget, comboPointMetrics)
			}
		},
	}))
}

// 6P
// Increases Ambush, Eviscerate, Crimson Tempest, and Envenom damage by 50%
func (rogue *Rogue) applyScarletEnclaveDamage6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Rogue - Damage 6P Bonus"
	if rogue.HasAura(label) {
		return
	}

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label:    label,
		ActionID: core.ActionID{SpellID: 1226871},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_RogueAmbush | ClassSpellMask_RogueEviscerate | ClassSpellMask_RogueCrimsonTempestHit | ClassSpellMask_RogueEnvenom,
		IntValue:  50,
	}))
}

var ItemSetDuskwraithLeathers = core.NewItemSet(core.ItemSet{
	Name: "Duskwraith Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyScarletEnclaveTank2PBonus()
		},
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyScarletEnclaveTank4PBonus()
		},
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyScarletEnclaveTank6PBonus()
		},
	},
})

// Your stacks of Rolling with the Punches also increase all damage you deal by 1%.
func (rogue *Rogue) applyScarletEnclaveTank2PBonus() {

	if !rogue.HasRune(proto.RogueRune_RuneRollingWithThePunches) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Rogue - Tank 2P Bonus"
	if rogue.HasAura(label) {
		return
	}
	rogue.rollingWithThePunchesDamageMultiplier += 0.01

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.RollingWithThePunchesProcAura.ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				rogue.PseudoStats.DamageDealtMultiplier /= 1 + rogue.rollingWithThePunchesDamageMultiplier*float64(oldStacks)
				rogue.PseudoStats.DamageDealtMultiplier *= 1 + rogue.rollingWithThePunchesDamageMultiplier*float64(newStacks)
			})
		},
	})
}

// Your Blade Flurry now also strikes a third target and increases your attack speed by an additional 10%. In addition, each combo point you spend reduces the remaining cooldown on your Blade Flurry by 0.5 sec.
func (rogue *Rogue) applyScarletEnclaveTank4PBonus() {

	if !rogue.Talents.BladeFlurry {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Rogue - Tank 4P Bonus"

	if rogue.HasAura(label) {
		return
	}

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
				cdReduction := time.Millisecond * time.Duration(500) * time.Duration(comboPoints)
				rogue.BladeFlurry.CD.ModifyRemainingCooldown(sim, -cdReduction)
			})

			rogue.bladeFlurryAttackSpeedBonus += 0.1
			rogue.bladeFlurryTargetCount += 1

		},
	}))
}

// Your Rolling with the Punches now grants 2% more health and 1% more damage per stack. At 5 stacks, each time you Dodge or Parry you will gain 10 Energy.
func (rogue *Rogue) applyScarletEnclaveTank6PBonus() {

	if !rogue.HasRune(proto.RogueRune_RuneRollingWithThePunches) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Rogue - Tank 6P Bonus"

	if rogue.HasAura(label) {
		return
	}

	rogue.rollingWithThePunchesBonusHealthStackMultiplier += 0.02
	rogue.rollingWithThePunchesDamageMultiplier += 0.01

	metrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 1226957})

	energyProc := rogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 1226957},
		Flags:    core.SpellFlagNoLifecycleCallbacks,
		ApplyEffects: func(sim *core.Simulation, u *core.Unit, spell *core.Spell) {
			rogue.AddEnergy(sim, 10, metrics)
		},
	})

	energyAura := rogue.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226957},
		Label:    "Float Like a Butterfly, Sting Like a Bee",
		Duration: core.NeverExpires,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeDodge | core.OutcomeParry) {
				energyProc.Cast(sim, result.Target)
			}
		},
	})

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.RollingWithThePunchesProcAura.ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				if newStacks == 5 {
					energyAura.Activate(sim)
				} else if newStacks < 5 && oldStacks == 5 {
					energyAura.Deactivate(sim)
				}
			})
		},
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
