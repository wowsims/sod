package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetDreamwalkerEclipse = core.NewItemSet(core.ItemSet{
	Name: "Dreamwalker Eclipse",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyNaxxramasBalance2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyNaxxramasBalance4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyNaxxramasBalance6PBonus()
		},
	},
})

// Your Moonfire and Sunfire deal 20% more damage.
func (druid *Druid) applyNaxxramasBalance2PBonus() {
	label := "S03 - Item - Naxxramas - Druid - Balance 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range druid.Moonfire {
				spell.DamageMultiplierAdditive += 0.20
			}

			if druid.Sunfire != nil {
				druid.Sunfire.DamageMultiplierAdditive += 0.20
			}
		},
	})
}

// The cooldown of your Starsurge spell is reduced by 1.5 sec.
func (druid *Druid) applyNaxxramasBalance4PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneLegsStarsurge) {
		return
	}

	label := "S03 - Item - Naxxramas - Druid - Balance 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			druid.Starsurge.CD.FlatModifier -= time.Millisecond * 1500
		},
	})
}

// When your Starsurge strikes an Undead target, the remaining duration on your active Starfall is reset to 10 sec.
func (druid *Druid) applyNaxxramasBalance6PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneLegsStarsurge) || !druid.HasRune(proto.DruidRune_RuneCloakStarfall) {
		return
	}

	label := "S03 - Item - Naxxramas - Druid - Balance 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if starfallDot := druid.Starfall.AOEDot(); starfallDot.IsActive() && result.Target.MobType == proto.MobType_MobTypeUndead && spell.SpellCode == SpellCode_DruidStarsurge && result.Landed() {
				starfallDot.Refresh(sim)
			}
		},
	}))
}

var ItemSetDreamwalkerFerocity = core.NewItemSet(core.ItemSet{
	Name: "Dreamwalker Ferocity",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyNaxxramasFeral2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyNaxxramasFeral4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyNaxxramasFeral6PBonus()
		},
	},
})

// Your Rake now deals its periodic damage every 1 sec, increasing its total damage over time by 200%.
func (druid *Druid) applyNaxxramasFeral2PBonus() {
	label := "S03 - Item - Naxxramas - Druid - Feral 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1218476}, // Tracking in APL
		Label:    label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, dot := range druid.Rake.Dots() {
				if dot == nil {
					continue
				}

				dot.TickLength /= 3
				dot.NumberOfTicks *= 3
			}
		},
	})
}

// Your Tiger's Fury cooldown is reduced by 50%.
func (druid *Druid) applyNaxxramasFeral4PBonus() {
	label := "S03 - Item - Naxxramas - Druid - Feral 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1218477}, // Tracking in APL
		Label:    label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			druid.TigersFury.CD.Multiplier *= 0.5
		},
	})
}

// Each time you deal Bleed damage to an Undead target, you gain 1% increased damage done to Undead for 30 sec, stacking up to 25 times.
func (druid *Druid) applyNaxxramasFeral6PBonus() {
	label := "S03 - Item - Naxxramas - Druid - Feral 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	var undeadTargets []*core.Unit

	buffAura := druid.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1218479},
		Label:     "Undead Slaying",
		Duration:  time.Second * 30,
		MaxStacks: 25,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			undeadTargets = core.FilterSlice(druid.Env.Encounter.TargetUnits, func(unit *core.Unit) bool { return unit.MobType == proto.MobType_MobTypeUndead })
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			for _, unit := range undeadTargets {
				druid.AttackTables[unit.UnitIndex][proto.CastType_CastTypeMainHand].DamageDealtMultiplier /= 1 + 0.01*float64(oldStacks)
				druid.AttackTables[unit.UnitIndex][proto.CastType_CastTypeMainHand].DamageDealtMultiplier *= 1 + 0.01*float64(newStacks)
			}
		},
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1218478}, // Tracking in APL
		Label:    label,
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellSchool.Matches(core.SpellSchoolPhysical) && result.Target.MobType == proto.MobType_MobTypeUndead {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}

var ItemSetDreamwalkerGuardian = core.NewItemSet(core.ItemSet{
	Name: "Dreamwalker Guardian",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyNaxxramasGuardian2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyNaxxramasGuardian4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyNaxxramasGuardian6PBonus()
		},
	},
})

// Your Growl ability never misses, and your chance to be Dodged or Parried is reduced by 2%.
func (druid *Druid) applyNaxxramasGuardian2PBonus() {
	label := "S03 - Item - Naxxramas - Druid - Guardian 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{stats.Expertise: 2 * core.ExpertiseRatingPerExpertiseChance}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
				aura.Unit.AddStats(bonusStats)
			} else {
				aura.Unit.AddStatsDynamic(sim, bonusStats)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
				aura.Unit.AddStats(bonusStats.Invert())
			} else {
				aura.Unit.AddStatsDynamic(sim, bonusStats.Invert())
			}
		},
	}))
}

// Reduces the cooldown on your Survival Instincts by 2 min, and reduces the cooldown on your Berserk ability by 2 min.
func (druid *Druid) applyNaxxramasGuardian4PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneFeetSurvivalInstincts) && !druid.HasRune(proto.DruidRune_RuneBeltBerserk) {
		return
	}

	label := "S03 - Item - Naxxramas - Druid - Guardian 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if druid.SurvivalInstincts != nil {
				druid.SurvivalInstincts.CD.FlatModifier -= time.Minute * 2
			}

			if druid.Berserk != nil {
				druid.Berserk.CD.FlatModifier -= time.Minute * 2
			}
		},
	})
}

// When you take damage from an Undead enemy, the remaining duration of your active Frenzied Regeneration is reset to 10 sec.
// In addition, Frenzied Regeneration will never consume your last 15 Rage to generate healing.
func (druid *Druid) applyNaxxramasGuardian6PBonus() {
	label := "S03 - Item - Naxxramas - Druid - Guardian 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	// TODO: Implement rage part when Frenzied Regeneration is implemented
	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if druid.FrenziedRegenerationAura != nil && druid.FrenziedRegenerationAura.IsActive() && spell.Unit.MobType == proto.MobType_MobTypeUndead && result.Landed() && result.Damage > 0 {
				druid.FrenziedRegenerationAura.Activate(sim)
			}
		},
	}))
}

var ItemSetDreamwalkerRaiment = core.NewItemSet(core.ItemSet{
	Name: "Dreamwalker Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Lifebloom now has a 100% chance to heal the target for one application of its dispel amount each time the target takes damage.
		2: func(agent core.Agent) {
		},
		// Each time your Rejuvenation rank 10 or rank 11 heals a target, you have a 25% chance to restore 60 mana, 8 energy, or 2 rage to your target.
		4: func(agent core.Agent) {
		},
		// Your Rejuvenation and Regrowth refresh their duration to full each time their target is damaged by an Undead enemy.
		6: func(agent core.Agent) {
		},
	},
})
