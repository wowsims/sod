package priest

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var ItemSetTwilightOfTheOracle = core.NewItemSet(core.ItemSet{
	Name: "Twilight of the Oracle",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyTAQShadow2PBonus()
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyTAQShadow4PBonus()
		},
	},
})

// Your Mind Flay no longer loses duration from taking damage and launches a free Mind Spike at the target on cast.
func (priest *Priest) applyTAQShadow2PBonus() {
	if !priest.Talents.MindFlay {
		return
	}

	label := "S03 - Item - TAQ - Priest - Shadow 2P Bonus"
	if priest.HasAura(label) {
		return
	}
	priest.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spells := range priest.MindFlay {
				for _, spell := range spells {
					if spell == nil {
						continue
					}

					spell.PushbackReduction += 1

					if priest.MindSpike == nil {
						continue
					}

					oldApplyEffects := spell.ApplyEffects
					spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						priest.MindSpike.DefaultCast.GCD = 0
						priest.MindSpike.CastTimeMultiplier -= 1
						priest.MindSpike.Cost.Multiplier -= 100
						priest.MindSpike.Cast(sim, target)
						priest.MindSpike.DefaultCast.GCD = core.GCDDefault
						priest.MindSpike.CastTimeMultiplier += 1
						priest.MindSpike.Cost.Multiplier += 100

						oldApplyEffects(sim, target, spell)
					}
				}
			}
		},
	})
}

// Your Mind Spike is now instant, deals 10% more damage, and can be cast while channeling another spell.
func (priest *Priest) applyTAQShadow4PBonus() {
	if !priest.HasRune(proto.PriestRune_RuneWaistMindSpike) {
		return
	}

	label := "S03 - Item - TAQ - Priest - Shadow 4P Bonus"
	if priest.HasAura(label) {
		return
	}

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  ClassSpellMask_PriestMindSpike,
		FloatValue: -1,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  ClassSpellMask_PriestMindSpike,
		FloatValue: 0.10,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Custom,
		ClassMask: ClassSpellMask_PriestMindSpike,
		ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
			priest.MindSpike.Flags |= core.SpellFlagCastWhileChanneling
		},
		RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
			priest.MindSpike.Flags ^= core.SpellFlagCastWhileChanneling
		},
	}))
}

var ItemSetDawnOfTheOracle = core.NewItemSet(core.ItemSet{
	Name: "Dawn of the Oracle",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Prayer of Mending gains 2 additional charges.
		2: func(agent core.Agent) {
		},
		// Your Circle of Healing now heals the most injured member of the target party for 100% more.
		4: func(agent core.Agent) {
		},
	},
})

var ItemSetFineryOfInfiniteWisdom = core.NewItemSet(core.ItemSet{
	Name: "Finery of Infinite Wisdom",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyRAQShadow3PBonus()
		},
	},
})

// Your Pain and Suffering rune can now refresh the duration of Devouring Plague.
func (priest *Priest) applyRAQShadow3PBonus() {
	if !priest.HasRune(proto.PriestRune_RuneHelmPainAndSuffering) {
		return
	}

	label := "S03 - Item - RAQ - Priest - Shadow 3P Bonus"
	if priest.HasAura(label) {
		return
	}

	priest.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			priest.PainAndSufferingDoTSpells = append(
				priest.PainAndSufferingDoTSpells,
				core.FilterSlice(priest.DevouringPlague, func(spell *core.Spell) bool { return spell != nil })...,
			)
		},
	})
}
