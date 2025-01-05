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

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells := core.FilterSlice(
				core.Flatten(
					[][]*DruidSpell{
						druid.Wrath,
						druid.Starfire,
						{druid.Starsurge},
					},
				),
				func(spell *DruidSpell) bool { return spell != nil },
			)

			for _, spell := range affectedSpells {
				spell.CritDamageBonus += 0.60
			}
		},
	})
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

	druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213171}, // Tracking in APL
		Label:    label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			druid.ShredPositionOverride = true
			if !druid.PseudoStats.InFrontOfTarget {
				druid.Shred.DamageMultiplierAdditive += 0.15
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

	core.MakePermanent(druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213174}, // Tracking in APL
		Label:    label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Outcome.Matches(core.OutcomeCrit) || !(spell == druid.Shred.Spell || spell == druid.MangleCat.Spell || spell == druid.FerociousBite.Spell) {
				return
			}

			dot := toothAndClawSpell.Dot(result.Target)
			dotDamage := result.Damage * 0.3
			if dot.IsActive() {
				dotDamage += dot.SnapshotBaseDamage * float64(dot.MaxTicksRemaining())
			}
			dot.SnapshotBaseDamage = dotDamage / float64(dot.NumberOfTicks)
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

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// TODO
		},
	})
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

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// TODO
		},
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
