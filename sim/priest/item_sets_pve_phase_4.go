package priest

import (
	"time"

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

	damageModifier := 0.50
	durationDivisor := time.Duration(2)

	var affectedSpells []*core.Spell

	buffAura := priest.GetOrRegisterAura(core.Aura{
		Label:    "Melting Faces",
		ActionID: core.ActionID{SpellID: 456549},
		Duration: core.NeverExpires,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(
				core.Flatten(priest.MindFlay),
				func(spell *core.Spell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += damageModifier
				for _, dot := range spell.Dots() {
					if dot != nil {
						dot.TickLength /= durationDivisor
					}
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive -= damageModifier
				for _, dot := range spell.Dots() {
					if dot != nil {
						dot.TickLength *= durationDivisor
					}
				}
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_PriestMindFlay) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(priest.GetOrRegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_PriestMindBlast) && result.DidCrit() {
				buffAura.Activate(sim)
			}
		},
	}))
}
