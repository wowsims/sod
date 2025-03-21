package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// var ItemSetInquisitionWarplate = core.NewItemSet(core.ItemSet{
// 	ID:   1940,
// 	Name: "Inquisition Warplate",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		// Crusader Strike and Exorcism grant you Holy Power, increasing all Holy damage you deal by 20%, stacking up to 3 times.
// 		2: func(agent core.Agent) {
// 			paladin := agent.(PaladinAgent).GetPaladin()
// 			paladin.applyScarletEnclaveRetribution2PBonus()
// 		},
// 		// Divine Storm, Holy Shock, and Holy Wrath consume all your Holy Power, dealing 100% increased damage per Holy Power you have accumulated.
// 		4: func(agent core.Agent) {
// 			paladin := agent.(PaladinAgent).GetPaladin()
// 			paladin.applyScarletEnclaveRetribution4PBonus()
// 		},
// 		// Consuming Holy Power increases your Attack Power by 10% per Holy Power consumed for 10 sec.
// 		6: func(agent core.Agent) {
// 			paladin := agent.(PaladinAgent).GetPaladin()
// 			paladin.applyScarletEnclaveRetribution6PBonus()
// 		},
// 	},
// })

// Crusader Strike and Exorcism grant you Holy Power, increasing all Holy damage you deal by 20%, stacking up to 3 times.
func (paladin *Paladin) applyScarletEnclaveRetribution2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Paladin - Retribution 2P Bonus"
	if paladin.HasAura(label) {
		return
	}

	paladin.registerHolyPowerAura()

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: PaladinTSERet2P},
		Label:    label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(ClassSpellMask_PaladinExorcism | ClassSpellMask_PaladinCrusaderStrike) {
				return
			}

			if !result.Landed() {
				return
			}

			paladin.holyPowerAura.Activate(sim)
			paladin.holyPowerAura.AddStack(sim)
		},
	}))
}

func (paladin *Paladin) registerHolyPowerAura() {
	if paladin.holyPowerAura == nil {
		paladin.holyPowerAura = paladin.GetOrRegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 1226461},
			Label:     "Holy Power",
			MaxStacks: 3,
			Duration:  time.Second * 15,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= ((1.0 + (0.2 * float64(newStacks))) / (1.0 + (0.2 * float64(oldStacks))))
			},
		})
	}
}

// Divine Storm, Holy Shock, and Holy Wrath consume all your Holy Power, dealing 100% increased damage per Holy Power you have accumulated.
func (paladin *Paladin) applyScarletEnclaveRetribution4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Paladin - Retribution 4P Bonus"
	if paladin.HasAura(label) {
		return
	}

	paladin.registerHolyPowerAura()

	spenderMask := ClassSpellMask_PaladinDivineStorm | ClassSpellMask_PaladinHolyShock | ClassSpellMask_PaladinHolyWrath
	damageMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: spenderMask,
	})

	paladin.holyPowerAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		damageMod.Activate()
	}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		damageMod.Deactivate()
	}).ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
		damageMod.UpdateIntValue(100 * int64(newStacks))
	})

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: PaladinTSERet4P},
		Label:    label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(spenderMask) || !paladin.holyPowerAura.IsActive() {
				return
			}

			if !result.Landed() {
				return
			}

			if paladin.onHolyPowerSpent != nil {
				paladin.onHolyPowerSpent(sim, paladin.holyPowerAura.GetStacks())
			}

			paladin.holyPowerAura.Deactivate(sim)
		},
	}))
}

// Consuming Holy Power increases your Attack Power by 10% per Holy Power consumed for 10 sec.
func (paladin *Paladin) applyScarletEnclaveRetribution6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Paladin - Retribution 6P Bonus"
	if paladin.HasAura(label) {
		return
	}

	templarMultiplier := []*stats.StatDependency{
		paladin.NewDynamicMultiplyStat(stats.AttackPower, 1.0),
		paladin.NewDynamicMultiplyStat(stats.AttackPower, 1.1),
		paladin.NewDynamicMultiplyStat(stats.AttackPower, 1.2),
		paladin.NewDynamicMultiplyStat(stats.AttackPower, 1.3),
	}

	templarAura := paladin.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1226464},
		Label:     "Templar",
		Duration:  time.Second * 10,
		MaxStacks: 3,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			paladin.EnableDynamicStatDep(sim, templarMultiplier[0])
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			paladin.DisableDynamicStatDep(sim, templarMultiplier[oldStacks])
			paladin.EnableDynamicStatDep(sim, templarMultiplier[newStacks])
		},
	})

	paladin.onHolyPowerSpent = func(sim *core.Simulation, holyPower int32) {
		if holyPower > 0 {
			templarAura.Activate(sim)
			templarAura.SetStacks(sim, holyPower)
		}
	}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: PaladinTSERet6P},
		Label:    label,
	}))
}

// var ItemSetInquisitionBulwark = core.NewItemSet(core.ItemSet{
// 	ID:   1942,
// 	Name: "Inquisition Bulwark",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		// Shield of Righteousness also increases your Block Value by 100% for 6 sec.
// 		2: func(agent core.Agent) {
// 			paladin := agent.(PaladinAgent).GetPaladin()
// 			paladin.applyScarletEnclaveProtection2PBonus()
// 		},
// 		// Shield of Righteousness deals percentage increased damage equal to your Block Chance.
// 		4: func(agent core.Agent) {
// 			paladin := agent.(PaladinAgent).GetPaladin()
// 			paladin.applyScarletEnclaveProtection4PBonus()
// 		},
// 		// Each time you Block, you heal for 50% of your Block Value.
// 		6: func(agent core.Agent) {
// 			paladin := agent.(PaladinAgent).GetPaladin()
// 			paladin.applyScarletEnclaveProtection6PBonus()
// 		},
// 	},
// })

// Shield of Righteousness also increases your Block Value by 100% for 6 sec.
func (paladin *Paladin) applyScarletEnclaveProtection2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Paladin - Protection 2P Bonus"
	if paladin.HasAura(label) {
		return
	}

	blockValueModifier := paladin.NewDynamicMultiplyStat(stats.BlockValue, 2.0)

	righteousShieldAura := paladin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226466},
		Label:    "Righteous Shield",
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.EnableDynamicStatDep(sim, blockValueModifier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.DisableDynamicStatDep(sim, blockValueModifier)
		},
	})

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: PaladinTSEProt2P},
		Label:    label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(ClassSpellMask_PaladinShieldOfRighteousness) {
				return
			}

			if !result.Landed() {
				return
			}

			righteousShieldAura.Activate(sim)
		},
	}))
}

// Shield of Righteousness deals percentage increased damage equal to your Block Chance.
func (paladin *Paladin) applyScarletEnclaveProtection4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Paladin - Protection 4P Bonus"
	if paladin.HasAura(label) {
		return
	}

	damageMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  ClassSpellMask_PaladinShieldOfRighteousness,
		FloatValue: 1.0,
	})

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: PaladinTSEProt4P},
		Label:    label,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if !spell.Matches(ClassSpellMask_PaladinShieldOfRighteousness) {
				return
			}

			damageMod.UpdateFloatValue(1.0 + (paladin.GetStat(stats.Block) / 100))
		},
	}))
}

// Each time you Block, you heal for 50% of your Block Value.
func (paladin *Paladin) applyScarletEnclaveProtection6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Paladin - Protection 6P Bonus"
	if paladin.HasAura(label) {
		return
	}

	righteousBlock := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1226480},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			blockValue := paladin.BlockValue()
			healAmount := blockValue * 0.5
			spell.CalcAndDealHealing(sim, target, healAmount, spell.OutcomeHealingCrit)
		},
	})

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: PaladinTSEProt6P},
		Label:    label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result == nil || !result.DidBlock() {
				return
			}

			righteousBlock.Cast(sim, &paladin.Unit)
		},
	}))
}

// var ItemSetInquisitionArmor = core.NewItemSet(core.ItemSet{
// 	ID:   1941,
// 	Name: "Inquisition Armor",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		// Lay on Hands also grants you 20% spell haste for 1 min.
// 		2: func(agent core.Agent) {
// 			paladin := agent.(PaladinAgent).GetPaladin()
// 			paladin.applyScarletEnclaveHoly2PBonus()
// 		},
// 		// Casting Holy Light, Flash of Light, or Divine Light on your Beacon of Light target causes you to gain 100% of the spell's base mana cost.
// 		4: func(agent core.Agent) {
// 		},
// 		// An additional 25% of your healing is transferred to your Beacon of Light target.
// 		6: func(agent core.Agent) {
// 		},
// 	},
// })

// // Lay on Hands also grants you 20% spell haste for 1 min.
// func (paladin *Paladin) applyScarletEnclaveHoly2PBonus() {
// 	label := "S03 - Item - Scarlet Enclave - Paladin - Holy 2P Bonus"
// 	if paladin.HasAura(label) {
// 		return
// 	}

// 	spellHasteModifier := paladin.NewDynamicMultiplyStat(stats.SpellHaste, 20.0)

// 	emergencyAura := paladin.RegisterAura(core.Aura{
// 		ActionID: core.ActionID{SpellID: 1226451},
// 		Label:    "Emergency",
// 		Duration: time.Second * 60,

// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			paladin.EnableDynamicStatDep(sim, spellHasteModifier)
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			paladin.DisableDynamicStatDep(sim, spellHasteModifier)
// 		},
// 	})

// 	core.MakePermanent(paladin.RegisterAura(core.Aura{
// 		ActionID: core.ActionID{SpellID: PaladinTSEHoly2P},
// 		Label:    label,
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if !spell.Matches(ClassSpellMask_PaladinLayOnHands) {
// 				return
// 			}

// 			emergencyAura.Activate(sim)
// 		},
// 	}))
// }
