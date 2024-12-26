package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
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
		Label: label,
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
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			druid.TigersFury.CD.Multiplier -= 50
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
		Label: label,
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellSchool.Matches(core.SpellSchoolPhysical) && result.Target.MobType == proto.MobType_MobTypeUndead {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}
