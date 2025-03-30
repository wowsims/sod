package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetGenesisEclipse = core.NewItemSet(core.ItemSet{
	Name: "Genesis Eclipse",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyTAQBalance2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyTAQBalance4PBonus()
		},
	},
})

// Your Nature's Grace talent gains 1 additional charge each time it triggers.
func (druid *Druid) applyTAQBalance2PBonus() {
	if !druid.Talents.NaturesGrace {
		return
	}

	label := "S03 - Item - TAQ - Druid - Balance 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			druid.NaturesGraceProcAura.MaxStacks += 1
		},
	})
}

// Increases the critical strike damage bonus of your Starfire, Starsurge, and Wrath by 60%.
func (druid *Druid) applyTAQBalance4PBonus() {
	label := "S03 - Item - TAQ - Druid - Balance 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CritDamageBonus_Flat,
		ClassMask:  ClassSpellMask_DruidWrath | ClassSpellMask_DruidStarfire | ClassSpellMask_DruidStarsurge,
		FloatValue: 0.60,
	}))
}

var ItemSetGenesisCunning = core.NewItemSet(core.ItemSet{
	Name: "Genesis Cunning",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyTAQFeral2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyTAQFeral4PBonus()
		},
	},
})

// Your Shred no longer has a positional requirement, but deals 15% more damage if you are behind the target.
func (druid *Druid) applyTAQFeral2PBonus() {
	label := "S03 - Item - TAQ - Druid - Feral 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	damageMod := druid.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_DruidShred,
		IntValue:  15,
	})

	druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213171}, // Tracking in APL
		Label:    label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			druid.ShredPositionOverride = true
			if !druid.PseudoStats.InFrontOfTarget {
				damageMod.Activate()
			}
		},
	})
}

// Your Mangle, Shred, and Ferocious Bite critical strikes cause your target to Bleed for 30% of the damage done over the next 4 sec sec.
func (druid *Druid) applyTAQFeral4PBonus() {
	label := "S03 - Item - TAQ - Druid - Feral 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	// This is the spell used for the bleed proc.
	// https://www.wowhead.com/classic/spell=1213176/tooth-and-claw
	toothAndClawSpell := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1213176},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Tooth and Claw",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					druid.BleedsActive[aura.Unit.UnitIndex]++
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					druid.BleedsActive[aura.Unit.UnitIndex]--
				},
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrRefresh(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
		},
	})

	catProcMasks := ClassSpellMask_DruidShred | ClassSpellMask_DruidMangleCat | ClassSpellMask_DruidFerociousBite
	bearProcMasks := ClassSpellMask_DruidMangleBear

	core.MakePermanent(druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213174}, // Tracking in APL
		Label:    label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			} else if druid.form.Matches(Cat) && !spell.Matches(catProcMasks) {
				return
			} else if druid.form.Matches(Bear) && !spell.Matches(bearProcMasks) {
				return
			}

			dot := toothAndClawSpell.Dot(result.Target)
			if dot == nil {
				return
			}

			newDamage := result.Damage * 0.3

			dot.SnapshotBaseDamage = (dot.OutstandingDmg() + newDamage) / float64(dot.NumberOfTicks)
			dot.SnapshotAttackerMultiplier = 1

			toothAndClawSpell.Cast(sim, result.Target)
		},
	}))
}

var ItemSetGenesisBounty = core.NewItemSet(core.ItemSet{
	Name: "Genesis Bounty",
	Bonuses: map[int32]core.ApplyEffect{
		// Reduces the cooldown of your Rebirth and Innervate spells by 65%.
		2: func(agent core.Agent) {
		},
		// Your critical heals with Healing Touch, Regrowth, and Nourish instantly heal the target for another 50% of the healing they dealt.
		4: func(agent core.Agent) {
		},
	},
})

var ItemSetGenesisFury = core.NewItemSet(core.ItemSet{
	Name: "Genesis Fury",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyTAQGuardian2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyTAQGuardian4PBonus()
		},
	},
})

// Each time you Dodge while in Dire Bear Form, you gain 10% increased damage on your next Mangle or Swipe, stacking up to 5 times.
func (druid *Druid) applyTAQGuardian2PBonus() {
	label := "S03 - Item - TAQ - Druid - Guardian 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	damageMod := druid.AddDynamicMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_DruidMangleBear | ClassSpellMask_DruidSwipeBear,
		Kind:      core.SpellMod_DamageDone_Flat,
	})

	buffAura := druid.RegisterAura(core.Aura{
		Label:     "Guardian 2P Bonus Proc",
		ActionID:  core.ActionID{SpellID: 1213188},
		Duration:  time.Second * 10,
		MaxStacks: 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			damageMod.UpdateIntValue(10 * int64(newStacks))
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_DruidMangleBear | ClassSpellMask_DruidSwipeBear) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if druid.form == Bear && spell.ProcMask.Matches(core.ProcMaskMelee) && result.Outcome.Matches(core.OutcomeDodge) {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}

// Reduces the cooldown on Mangle (Bear) by 1.5 sec.
func (druid *Druid) applyTAQGuardian4PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneHandsMangle) {
		return
	}
	label := "S03 - Item - TAQ - Druid - Guardian 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
	})).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_DruidMangleBear,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Millisecond * 1500,
	})
}

var ItemSetSymbolsOfUnendingLife = core.NewItemSet(core.ItemSet{
	Name: "Symbols of Unending Life",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyRAQFeral3PBonus()
		},
	},
})

// Your melee attacks have 5% less chance to be Dodged or Parried.
func (druid *Druid) applyRAQFeral3PBonus() {
	label := "S03 - Item - RAQ - Druid - Feral 3P Bonus"
	if druid.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{stats.Expertise: 5 * core.ExpertiseRatingPerExpertiseChance}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachStatsBuff(bonusStats))
}
