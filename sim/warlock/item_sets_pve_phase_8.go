package warlock

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// var ItemSetHereticRaiment = core.NewItemSet(core.ItemSet{
// 	Name: "Heretic Raiment",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			warlock := agent.(WarlockAgent).GetWarlock()
// 			warlock.applyScarletEnclaveDamage2PBonus()
// 		},
// 		4: func(agent core.Agent) {
// 			warlock := agent.(WarlockAgent).GetWarlock()
// 			warlock.applyScarletEnclaveDamage4PBonus()
// 		},
// 		6: func(agent core.Agent) {
// 			warlock := agent.(WarlockAgent).GetWarlock()
// 			warlock.applyScarletEnclaveDamage6PBonus()
// 		},
// 	},
// })

// Your Shadow and Fire non-periodic critical strikes cause the target to Burn for 25% of the damage they deal over 4 sec.
func (warlock *Warlock) applyScarletEnclaveDamage2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Warlock - Damage 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	// This is the spell used for the burn proc.
	// https://www.wowhead.com/classic-ptr/spell=1227180/burn
	burnSpell := warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1227180},
		SpellSchool: core.SpellSchoolShadow | core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagIgnoreAttackerModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Burning",
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrRefresh(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
		},
	})

	core.MakeProcTriggerAura(&warlock.Unit, core.ProcTrigger{
		Name:             label,
		Callback:         core.CallbackOnSpellHitDealt,
		Outcome:          core.OutcomeCrit,
		CanProcFromProcs: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellSchool.Matches(core.SpellSchoolShadow | core.SpellSchoolFire) {
				dot := burnSpell.Dot(result.Target)
				dotDamage := result.Damage * 0.3
				if dot.IsActive() {
					dotDamage += dot.SnapshotBaseDamage * float64(dot.MaxTicksRemaining())
				}
				dot.SnapshotBaseDamage = dotDamage / float64(dot.NumberOfTicks)
				dot.SnapshotAttackerMultiplier = 1

				burnSpell.Cast(sim, result.Target)
			}
		},
	})
}

// Your Incinerate, Shadow Bolt, and Soul Fire deal 20% more damage to targets afflicted with your Corruption.
func (warlock *Warlock) applyScarletEnclaveDamage4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Warlock - Damage 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	classMask := ClassSpellMask_WarlockIncinerate | ClassSpellMask_WarlockShadowBolt
	damageMod := warlock.AddDynamicMod(core.SpellModConfig{
		ClassMask:  classMask,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 1.0,
	})

	core.MakeProcTriggerAura(&warlock.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: classMask,
		Callback:       core.CallbackOnApplyEffects,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			hasCorruption := slices.ContainsFunc(warlock.Corruption, func(spell *core.Spell) bool {
				return spell.Dot(result.Target).IsActive()
			})
			damageMod.UpdateFloatValue(core.TernaryFloat64(hasCorruption, 1.20, 1.0))
			damageMod.Activate()
		},
	})
}

// Your periodic critical strikes grant 15% spellcasting haste for 15 sec, and your Backdraft grants an additional 15% spellcasting haste.
func (warlock *Warlock) applyScarletEnclaveDamage6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Warlock - Damage 6P Bonus"
	if warlock.HasAura(label) {
		return
	}

	hasteAura := warlock.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1227200},
		Label:    "Wickedness",
		Duration: time.Second * 15,
	}).AttachMultiplyCastSpeed(&warlock.Unit, 1.15)

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if warlock.BackdraftAura != nil {
				warlock.BackdraftAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
					hasteAura.Activate(sim)
				})
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
				hasteAura.Activate(sim)
			}
		},
	}))
}

// var ItemSetHereticStitchings = core.NewItemSet(core.ItemSet{
// 	Name: "Heretic Stitchings",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			warlock := agent.(WarlockAgent).GetWarlock()
// 			warlock.applyScarletEnclaveTank2PBonus()
// 		},
// 		4: func(agent core.Agent) {
// 			warlock := agent.(WarlockAgent).GetWarlock()
// 			warlock.applyScarletEnclaveTank4PBonus()
// 		},
// 		6: func(agent core.Agent) {
// 			warlock := agent.(WarlockAgent).GetWarlock()
// 			warlock.applyScarletEnclaveTank6PBonus()
// 		},
// 	},
// })

// Your Shadowcleave now applies your Corruption Rank 7 to every target it hits but its duration is only 12 sec.
func (warlock *Warlock) applyScarletEnclaveTank2PBonus() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsMetamorphosis) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Warlock - Tank 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	core.MakeProcTriggerAura(&warlock.Unit, core.ProcTrigger{
		Name:           label,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: ClassSpellMask_WarlockShadowCleave,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dot := warlock.Corruption[6].Dot(result.Target)
			dot.ApplyOrRefresh(sim)
			dot.UpdateExpires(sim, sim.CurrentTime+time.Second*12)
		},
	})
}

// You heal for 5% of all damage done by your Corruption.
func (warlock *Warlock) applyScarletEnclaveTank4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Warlock - Tank 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	healthMetrics := warlock.NewHealthMetrics(core.ActionID{SpellID: 1227207})

	core.MakeProcTriggerAura(&warlock.Unit, core.ProcTrigger{
		Name:             label,
		Callback:         core.CallbackOnPeriodicDamageDealt,
		ClassSpellMask:   ClassSpellMask_WarlockCorruption,
		CanProcFromProcs: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warlock.GainHealth(sim, result.Damage*0.05, healthMetrics)
		},
	})
}

// Increases the damage reduction from Infernal Armor by an additional 10%.
func (warlock *Warlock) applyScarletEnclaveTank6PBonus() {
	if !warlock.HasRune(proto.WarlockRune_RuneCloakInfernalArmor) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Warlock - Tank 6P Bonus"
	if warlock.HasAura(label) {
		return
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warlock.bonusInfernalArmorMultiplier += 0.10
		},
	}))
}
