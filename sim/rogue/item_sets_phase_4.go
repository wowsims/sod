package rogue

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetDarkmantleArmor = core.NewItemSet(core.ItemSet{
	Name: "Darkmantle Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Chance on melee attack to restore 35 energy.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 27787}
			energyMetrics := c.NewEnergyMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID: actionID,
				Name:     "Rogue Armor Energize",
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
				ProcMask: core.ProcMaskMeleeWhiteHit,
				PPM:      1,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasEnergyBar() {
						c.AddEnergy(sim, 35, energyMetrics)
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

var ItemSetNightSlayerThrill = core.NewItemSet(core.ItemSet{
	Name: "Nightslayer Thrill",
	Bonuses: map[int32]core.ApplyEffect{
		// Feint also grants Avoidance for 6 sec, reducing all damage taken from area of effect attacks from non-players by 50%
		2: func(agent core.Agent) {
			// Not yet implemented
		},
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyT1Damage4PBonus()
		},
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyT1Damage6PBonus()
		},
	},
})

// Increases the critical strike damage bonus of your Poisons by 100%.
func (rogue *Rogue) applyT1Damage4PBonus() {
	label := "S03 - Item - T1 - Rogue - Damage 4P Bonus"
	if rogue.HasAura(label) {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range rogue.Spellbook {
				if spell.Flags.Matches(SpellFlagRoguePoison) {
					spell.CritDamageBonus += 1.0
				}
			}
		},
	})
}

// Your finishing moves have a 5% chance per combo point to make your next ability cost no energy.
// https://www.wowhead.com/classic/spell=457342/clearcasting
func (rogue *Rogue) applyT1Damage6PBonus() {
	label := "S03 - Item - T1 - Rogue - Damage 6P Bonus"
	if rogue.HasAura(label) {
		return
	}

	var affectedSpells []*core.Spell

	clearcastingAura := rogue.RegisterAura(core.Aura{
		Label:    fmt.Sprintf("Clearcasting (%s)", label),
		ActionID: core.ActionID{SpellID: 457342},
		Duration: time.Second * 15,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(
				rogue.Spellbook,
				func(spell *core.Spell) bool {
					return spell != nil && spell.Cost != nil && spell.Cost.CostType() == core.CostTypeEnergy
				},
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) { spell.Cost.Multiplier -= 100 })
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) { spell.Cost.Multiplier += 100 })
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if aura.RemainingDuration(sim) == aura.Duration || spell.DefaultCast.Cost == 0 {
				return
			}
			aura.Deactivate(sim)
		},
	})

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
				if sim.Proc(.05*float64(comboPoints), "Clearcasting (S03 - Item - T1 - Rogue - Damage 6P Bonus)") {
					clearcastingAura.Activate(sim)
				}
			})
		},
	})
}

var ItemSetNightSlayerBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Nightslayer Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyT1Tank2PBonus()
		},
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyT1Tank4PBonus()
		},
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyT1Tank6PBonus()
		},
	},
})

// While Just a Flesh Wound and Blade Dance are active, Crimson Tempest, Blunderbuss, and Fan of Knives cost 20 less Energy and generate 100% increased threat.
func (rogue *Rogue) applyT1Tank2PBonus() {
	if !rogue.HasRune(proto.RogueRune_RuneJustAFleshWound) || !rogue.HasRune(proto.RogueRune_RuneBladeDance) {
		return
	}

	label := "S03 - Item - T1 - Rogue - Tank 2P Bonus"
	if rogue.HasAura(label) {
		return
	}

	classSpellMasks := ClassSpellMask_RogueCrimsonTempest | ClassSpellMask_RogueCrimsonTempestHit | SpellClassMask_RogueBlunderbuss | SpellClassMask_RogueFanOfKnives
	buffAura := rogue.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 457351},
		Label:    fmt.Sprintf("Blade Dance (%s)", label),
		Duration: core.NeverExpires,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_Threat_Pct,
		ClassMask:  classSpellMasks,
		FloatValue: 2,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Flat,
		ClassMask: classSpellMasks,
		IntValue:  -20,
	})

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.BladeDanceAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				buffAura.Activate(sim)
			})

			rogue.BladeDanceAura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				buffAura.Deactivate(sim)
			})
		},
	})
}

// Vanish now reduces all Magic damage you take by 50% for its duration, but it no longer grants Stealth or breaks movement impairing effects.  - 457437
func (rogue *Rogue) applyT1Tank4PBonus() {
	label := "S03 - Item - T1 - Rogue - Tank 4P Bonus"
	if rogue.HasAura(label) {
		return
	}

	buffAura := rogue.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 457437},
		Label:    "Vanish",
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.SchoolDamageTakenMultiplier.MultiplyMagicSchools(0.5)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.SchoolDamageTakenMultiplier.MultiplyMagicSchools(1 / 0.5)
		},
	})

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// Override Vanish's Apply Effects to prevent activating stealth
			rogue.Vanish.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			}
		},
	})
}

// Your finishing moves have a 20% chance per combo point to make you take 20% less Physical damage from the next melee attack that hits you within 10 sec.
func (rogue *Rogue) applyT1Tank6PBonus() {
	label := "S03 - Item - T1 - Rogue - Tank 6P Bonus"
	if rogue.HasAura(label) {
		return
	}

	buffLabel := fmt.Sprintf("Resilient (%s)", label)
	damageTakenMultiplier := 0.8

	buffAura := rogue.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 457469},
		Label:    buffLabel,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= damageTakenMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= damageTakenMultiplier
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee) && result.Outcome.Matches(core.OutcomeLanded) {
				aura.Deactivate(sim)
			}
		},
	})

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
				if sim.Proc(0.2*float64(comboPoints), buffLabel) {
					buffAura.Activate(sim)
				}
			})
		},
	})
}
