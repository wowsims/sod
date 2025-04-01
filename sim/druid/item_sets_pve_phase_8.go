package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetWaywatcherEclipse = core.NewItemSet(core.ItemSet{
	Name: "Waywatcher Eclipse",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyScarletEnclaveBalance2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyScarletEnclaveBalance4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyScarletEnclaveBalance6PBonus()
		},
	},
})

// Your Starfire deals 20% more damage to targets with your Moonfire, and your Wrath deals 20% more damage to targets with your Sunfire.
func (druid *Druid) applyScarletEnclaveBalance2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Druid - Balance 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	hasSunfireRune := druid.HasRune(proto.DruidRune_RuneHandsSunfire)

	starfireDamageMod := druid.AddDynamicMod(core.SpellModConfig{
		ClassMask:  ClassSpellMask_DruidStarfire,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 1.0,
	})
	wrathDamageMod := druid.AddDynamicMod(core.SpellModConfig{
		ClassMask:  ClassSpellMask_DruidWrath,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 1.0,
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_DruidStarfire) {
				starfireDamageMod.UpdateFloatValue(1)
				for _, spell := range druid.Moonfire {
					if spell.Dot(target).IsActive() {
						starfireDamageMod.Activate()
						starfireDamageMod.UpdateFloatValue(1.20)
						return
					}
				}
			} else if spell.Matches(ClassSpellMask_DruidWrath) && hasSunfireRune {
				wrathDamageMod.UpdateFloatValue(1)
				if druid.Sunfire.Dot(target).IsActive() {
					wrathDamageMod.Activate()
					wrathDamageMod.UpdateFloatValue(1.20)
				}
			}
		},
	}))
}

// Your Starsurge now increases the damage of your next 2 Starfires.
func (druid *Druid) applyScarletEnclaveBalance4PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneLegsStarsurge) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Druid - Balance 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			druid.StarsurgeAura.MaxStacks += 1
		},
	}))
}

// Each time your Sunfire deals periodic damage, you gain 20% increased damage to your next Wrath, stacking up to 5 times.
// Each time your Moonfire deals periodic damage, you gain 20% increased damage to your next Stafire, stacking up to 5 times.
// These bonuses do not apply to Starsurge.
func (druid *Druid) applyScarletEnclaveBalance6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Druid - Balance 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	wrathDamageMod := druid.AddDynamicMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_DruidWrath,
		Kind:      core.SpellMod_DamageDone_Flat,
	})

	wrathAura := druid.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1226105},
		Label:     "Sunsurge",
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			wrathDamageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			wrathDamageMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			wrathDamageMod.UpdateIntValue(20 * int64(newStacks))
		},
	})

	starfireDamageMod := druid.AddDynamicMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_DruidStarfire,
		Kind:      core.SpellMod_DamageDone_Flat,
	})

	starfireAura := druid.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1226106},
		Label:     "Moonsurge",
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			starfireDamageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			starfireDamageMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			starfireDamageMod.UpdateIntValue(20 * int64(newStacks))
		},
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_DruidSunfire) {
				wrathAura.Activate(sim)
				wrathAura.AddStack(sim)
			} else if spell.Matches(ClassSpellMask_DruidMoonfire) {
				starfireAura.Activate(sim)
				starfireAura.AddStack(sim)
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_DruidWrath) {
				wrathAura.Deactivate(sim)
			} else if spell.Matches(ClassSpellMask_DruidStarfire) {
				starfireAura.Deactivate(sim)
			}
		},
	}))
}

var ItemSetWaywatcherFerocity = core.NewItemSet(core.ItemSet{
	Name: "Waywatcher Ferocity",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyScarletEnclaveFeral2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyScarletEnclaveFeral4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyScarletEnclaveFeral6PBonus()
		},
	},
})

// https://www.wowhead.com/classic-ptr/spell=1226109/s03-item-scarlet-enclave-druid-feral-2p-bonus
// You gain 2 Energy each time Rake or Rip deals periodic damage.
func (druid *Druid) applyScarletEnclaveFeral2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Druid - Feral 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	actionId := core.ActionID{SpellID: 1226113}
	scentOfBloodMetrics := druid.NewEnergyMetrics(actionId)
	ripOrRakeMask := ClassSpellMask_DruidRip | ClassSpellMask_DruidRake

	// https://www.wowhead.com/classic-ptr/spell=1226113/scent-of-blood
	// Not required at all to put this in a spell, but game does it like this, so why not?
	scentOfBlood := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:    actionId,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Unit.AddEnergy(sim, 2, scentOfBloodMetrics)
		},
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226109},
		Label:    label,
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ripOrRakeMask) && druid.HasEnergyBar() {
				scentOfBlood.Cast(sim, &druid.Unit)
			}
		},
	}))
}

// Multiplies the damage bonus from Tiger's Fury by 2.
func (druid *Druid) applyScarletEnclaveFeral4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Druid - Feral 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226116},
		Label:    label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// Q: wtf is this??
			// A: We can't access and manipulate the bonus directly because it only exists in the gain/expire callbacks.
			// Without hardcoding values it's probably the easiest to just double the change made by the original handler.
			// For the rune version we need to calculate the mult needed to effectively double the originally applied mult.
			// Hacky garbage, but probably better than having hardcoded values?

			origGain := druid.TigersFuryAura.OnGain
			origExpire := druid.TigersFuryAura.OnExpire

			if druid.TigersFury.ActionID.SpellID == 417045 {
				druid.TigersFuryAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
					oldVal := druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
					origGain(aura, sim)
					origTfMult := druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] / oldVal
					multToDoubleChange := (1 + (origTfMult-1)*2) / origTfMult // 1.3/1.15 for the default 15% on the rune
					druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= multToDoubleChange
				}
				druid.TigersFuryAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
					oldVal := druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
					origExpire(aura, sim)
					origTfMult := oldVal / druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
					multToDoubleChange := (1 + (origTfMult-1)*2) / origTfMult
					druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= multToDoubleChange
				}
			} else {
				druid.TigersFuryAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
					oldVal := druid.PseudoStats.BonusPhysicalDamage
					origGain(aura, sim)
					origBonus := druid.PseudoStats.BonusPhysicalDamage - oldVal
					druid.PseudoStats.BonusPhysicalDamage += origBonus
				}
				druid.TigersFuryAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
					oldVal := druid.PseudoStats.BonusPhysicalDamage
					origExpire(aura, sim)
					origBonus := druid.PseudoStats.BonusPhysicalDamage - oldVal
					druid.PseudoStats.BonusPhysicalDamage += origBonus
				}
			}
		},
	}))
}

// https://www.wowhead.com/classic-ptr/spell=1226119/s03-item-scarlet-enclave-druid-feral-6p-bonus
// Your Finishing Moves have a 20% chance per combo point spent to trigger Clearcasting and extend the duration of your active Tiger's Fury by 6 sec.
func (druid *Druid) applyScarletEnclaveFeral6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Druid - Feral 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	finisherMask := ClassSpellMask_DruidFerociousBite | ClassSpellMask_DruidRip | ClassSpellMask_DruidSavageRoar

	// https://www.wowhead.com/classic-ptr/spell=1226118/s03-item-scarlet-enclave-druid-feral-6p-bonus-trigger
	// Not needed but this way procs can be easily tracked.
	triggerSpell := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1226118},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			druid.ClearcastingAura.Activate(sim)
			if druid.TigersFuryAura.IsActive() {
				druid.TigersFuryAura.Refresh(sim)
			}
		},
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226119},
		Label:    label,
	}))

	druid.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
		if spell.Matches(finisherMask) && sim.Proc(0.2*float64(comboPoints), label) {
			triggerSpell.Cast(sim, &druid.Unit)
		}
	})
}

var ItemSetWaywatcherGuardian = core.NewItemSet(core.ItemSet{
	Name: "Waywatcher Guardian",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyScarletEnclaveGuardian2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyScarletEnclaveGuardian4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyScarletEnclaveGuardian6PBonus()
		},
	},
})

// Your melee critical strikes in Bear Form or Dire Bear Form grant you a shield lasting until cancelled that absorbs Physical damage equal to 25% of your Attack Power the next time you take Physical damage. Stacks up to 0 times.
func (druid *Druid) applyScarletEnclaveGuardian2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Druid - Guardian 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
	}))
}

// Increases the duration of your Berserk ability by 15 sec.
func (druid *Druid) applyScarletEnclaveGuardian4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Druid - Guardian 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			druid.BerserkAura.Duration += time.Second * 15
		},
	}))
}

// You gain 30% increased attack speed for 10 sec every time you deal a critical strike.
func (druid *Druid) applyScarletEnclaveGuardian6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Druid - Guardian 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	savageFlurryAura := druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226127},
		Label:    "Savage Flurry",
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.MultiplyAttackSpeed(sim, 1.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.MultiplyAttackSpeed(sim, 1.0/1.3)
		},
	})

	core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
		Name:     label,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeCrit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			savageFlurryAura.Activate(sim)
		},
	})
}

var ItemSetWaywatcherRaiment = core.NewItemSet(core.ItemSet{
	Name: "Waywatcher Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// Each time your Lifebloom heals a target, it has a 5% chance to make your next Healing Touch, Nourish, or Regrowth within 15 sec instant cast.
		2: func(agent core.Agent) {
		},
		// Targets with your active Rejuvenation Rank 10 or Rank 11 receive 20% increased healing from your spells.
		4: func(agent core.Agent) {
		},
		// When your Regrowth Rank 8 or Rank 9 deals a non-periodic critical heal, your Rejuvenation on that target will spread to all members of the target's party within 43.5 yards not already affected by your Rejuvenation.
		6: func(agent core.Agent) {
		},
	},
})
