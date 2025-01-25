package priest

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetTwilightOfTranscendence = core.NewItemSet(core.ItemSet{
	Name: "Twilight of Transcendence",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyT2Shadow2PBonus()
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyT2Shadow4PBonus()
		},
		6: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyT2Shadow6PBonus()
		},
	},
})

// Reduces the cooldown of your Shadow Word: Death spell by 6 sec.
func (priest *Priest) applyT2Shadow2PBonus() {
	if !priest.HasRune(proto.PriestRune_RuneHandsShadowWordDeath) {
		return
	}

	label := "S03 - Item - T2 - Priest - Shadow 2P Bonus"
	if priest.HasAura(label) {
		return
	}

	priest.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			priest.ShadowWordDeath.CD.FlatModifier -= time.Second * 6
		},
	})
}

// Your Shadow Word: Pain has a 2% chance per talent point in Spirit Tap to trigger your Spirit Tap talent when it deals damage,
// or a 20% chance per talent point when a target dies with your Shadow Word: Pain active.
func (priest *Priest) applyT2Shadow4PBonus() {
	if priest.Talents.SpiritTap == 0 {
		return
	}

	label := "S03 - Item - T2 - Priest - Shadow 4P Bonus"
	if priest.HasAura(label) {
		return
	}

	procChance := 0.02 * float64(priest.Talents.SpiritTap)

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if priest.Talents.InnerFocus {
				oldApplyEffects := priest.InnerFocus.ApplyEffects
				priest.InnerFocus.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					oldApplyEffects(sim, target, spell)
					priest.SpiritTapAura.Activate(sim)
				}
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_PriestShadowWordPain) && sim.Proc(procChance, "Proc Spirit Tap") {
				priest.SpiritTapAura.Activate(sim)
			}
		},
	}))
}

// While Spirit Tap is active, you deal 25% more shadow damage.
func (priest *Priest) applyT2Shadow6PBonus() {
	if priest.Talents.SpiritTap == 0 {
		return
	}

	label := "S03 - Item - T2 - Priest - Shadow 6P Bonus"
	if priest.HasAura(label) {
		return
	}

	priest.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			oldOnGain := priest.SpiritTapAura.OnGain
			priest.SpiritTapAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
				oldOnGain(aura, sim)
				priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.25
			}
			oldOnExpire := priest.SpiritTapAura.OnExpire
			priest.SpiritTapAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
				oldOnExpire(aura, sim)
				priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.25
			}
		},
	})
}

var ItemSetDawnOfTranscendence = core.NewItemSet(core.ItemSet{
	Name: "Dawn of Transcendence",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyT2Healer2PBonus()
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyT2Healer4PBonus()
		},
		// Circle of Healing and Penance also place a heal over time effect on their targets that heals for 25% as much over 15 sec.
		6: func(agent core.Agent) {
		},
	},
})

// Allows 15% of your Mana regeneration to continue while casting.
func (priest *Priest) applyT2Healer2PBonus() {
	label := "S03 - Item - T2 - Priest - Healer 2P Bonus"
	if priest.HasAura(label) {
		return
	}

	bonusRegen := 0.15

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			priest.PseudoStats.SpiritRegenRateCasting += bonusRegen
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.PseudoStats.SpiritRegenRateCasting -= bonusRegen
		},
	}))
}

// Your periodic healing has a 2% chance to make your next spell with a casting time less than 10 seconds an instant cast spell.
func (priest *Priest) applyT2Healer4PBonus() {
	label := "S03 - Item - T2 - Priest - Healer 4P Bonus"
	if priest.HasAura(label) {
		return
	}

	affectedSpells := []*core.Spell{}

	buffAura := priest.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467543},
		Label:    "Deliverance",
		Duration: core.NeverExpires, // TODO: Verify duration
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(priest.Spellbook, func(spell *core.Spell) bool {
				return spell.Matches(ClassSpellMask_PriestAll) && spell.DefaultCast.CastTime < time.Second*10
			})
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CastTimeMultiplier -= 1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CastTimeMultiplier += 1
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if slices.Contains(affectedSpells, spell) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
		OnPeriodicHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskSpellHealing) && sim.Proc(.02, "Proc Deliverance") {
				buffAura.Activate(sim)
			}
		},
	}))
}

var ItemSetConfessorsRaiment = core.NewItemSet(core.ItemSet{
	Name: "Confessor's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases healing done by up to 22 and damage done by up to 7 for all magical spells and effects.
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.AddStats(stats.Stats{
				stats.HealingPower: 22,
				stats.SpellDamage:  7,
			})
		},
		3: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyZGDiscipline3PBonus()
		},
		// Increases the damage absorbed by your Power Word: Shield spell by 20%.
		5: func(agent core.Agent) {
		},
	},
})

// Reduces the cooldown of your Penance spell by 6 sec.
func (priest *Priest) applyZGDiscipline3PBonus() {
	if !priest.HasRune(proto.PriestRune_RuneHandsPenance) {
		return
	}

	label := "S03 - Item - ZG - Priest - Discipline 3P Bonus"
	if priest.HasAura(label) {
		return
	}

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_PriestPenance,
		TimeValue: time.Second * 6,
	}))
}
