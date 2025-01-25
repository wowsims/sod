package priest

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetVestmentsOfTheVirtuous = core.NewItemSet(core.ItemSet{
	Name: "Vestments of the Virtuous",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Your spellcasts have a 6% chance to energize you for 300 mana.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450576}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Mana Proc on Cast - Vestments of the Devout",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.06,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasManaBar() {
						c.AddMana(sim, 300, manaMetrics)
					}
				},
			})
		},
		// +8 All Resistances.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// +200 Armor.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
	},
})

var ItemSetDawnProphecy = core.NewItemSet(core.ItemSet{
	Name: "Dawn Prophecy",
	Bonuses: map[int32]core.ApplyEffect{
		// -0.1 sec to the casting time of Flash Heal and -0.1 sec to the casting time of Greater Heal.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// Increases your critical strike chance with spells and attacks by 2%.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.MeleeCrit: 2 * core.CritRatingPerCritChance,
				stats.SpellCrit: 2 * core.CritRatingPerCritChance,
			})
		},
		// Increases your critical strike chance with Prayer of Healing and Circle of Healing by 25%.
		6: func(agent core.Agent) {
			// Nothing to do
		},
	},
})

var ItemSetTwilightProphecy = core.NewItemSet(core.ItemSet{
	Name: "Twilight Prophecy",
	Bonuses: map[int32]core.ApplyEffect{
		// You may cast Flash Heal while in Shadowform.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyT1Shadow4PBonus()
		},
		6: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyT1Shadow6PBonus()
		},
	},
})

// Increases your critical strike chance with spells and attacks by 2%.
func (priest *Priest) applyT1Shadow4PBonus() {
	label := "S03 - Item - T1 - Priest - Shadow 4P Bonus"
	if priest.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{
		stats.MeleeCrit: 2 * core.CritRatingPerCritChance,
		stats.SpellCrit: 2 * core.SpellCritRatingPerCritChance,
	}

	core.MakePermanent(priest.RegisterAura(core.Aura{
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

// Mind Blast critical strikes reduce the duration of your next Mind Flay by 50% while increasing its total damage by 50%.
func (priest *Priest) applyT1Shadow6PBonus() {
	if !priest.Talents.MindFlay {
		return
	}

	label := "S03 - Item - T1 - Priest - Shadow 6P Bonus"
	if priest.HasAura(label) {
		return
	}

	damageMod := priest.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  ClassSpellMask_PriestMindFlay,
		FloatValue: 0.50,
	})

	dotLengthMod := priest.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DotTickLength_Pct,
		ClassMask:  ClassSpellMask_PriestMindFlay,
		FloatValue: 0.50,
	})

	buffAura := priest.GetOrRegisterAura(core.Aura{
		Label:    "Melting Faces",
		ActionID: core.ActionID{SpellID: 456549},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
			dotLengthMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
			dotLengthMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_PriestMindFlay) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&priest.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: ClassSpellMask_PriestMindBlast,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeCrit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buffAura.Activate(sim)
		},
	})
}
