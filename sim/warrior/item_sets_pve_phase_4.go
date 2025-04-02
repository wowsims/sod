package warrior

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetBattlegearOfValor = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Heroism",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Chance on melee attack to heal you for 88 to 132 and energize you for 10 Rage
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450587}
			healthMetrics := c.NewHealthMetrics(core.ActionID{SpellID: 450589})
			rageMetrics := c.NewRageMetrics(core.ActionID{SpellID: 450589})

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:          actionID,
				Name:              "S03 - Warrior Armor Heal Trigger - Battlegear of Valor",
				Callback:          core.CallbackOnSpellHitDealt,
				Outcome:           core.OutcomeLanded,
				ProcMask:          core.ProcMaskMelee,
				SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
				PPM:               1,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					c.GainHealth(sim, sim.Roll(88, 132), healthMetrics)
					if c.HasRageBar() {
						c.AddRage(sim, 10, rageMetrics)
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

var ItemSetUnstoppableMight = core.NewItemSet(core.ItemSet{
	Name: "Unstoppable Might",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT1Damage2PBonus()
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT1Damage4PBonus()
		},
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT1Damage6PBonus()
		},
	},
})

// After changing stances, your next offensive ability's rage cost is reduced by 10.
func (warrior *Warrior) applyT1Damage2PBonus() {
	label := "S03 - Item - T1 - Warrior - Damage 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	var affectedSpells []*core.Spell
	tacticianAura := warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 464241},
		Label:    "Tactician",
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warrior.Spellbook {
				if spell.Cost != nil && spell.Cost.CostType() == core.CostTypeRage && !spell.Flags.Matches(core.SpellFlagHelpful) {
					affectedSpells = append(affectedSpells, spell)
				}
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.Cost.FlatModifier -= 10
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.Cost.FlatModifier += 10
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if slices.Contains(affectedSpells, spell) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(StanceCodes) {
				tacticianAura.Activate(sim)
			}
		},
	}))
}

// For 15 sec after leaving a stance, you can use abilities requiring that stance as if you were still in that stance.
func (warrior *Warrior) applyT1Damage4PBonus() {
	label := "S03 - Item - T1 - Warrior - Damage 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	duration := time.Second * 15

	battleStanceAura := warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 457706},
		Label:    "Echoes of Battle Stance",
		Duration: duration,
	})
	defStanceAura := warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 457699},
		Label:    "Echoes of Defensive Stance",
		Duration: duration,
	})
	berserkStanceAura := warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 457708},
		Label:    "Echoes of Berserker Stance",
		Duration: duration,
	})
	gladStanceAura := warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 457819},
		Label:    "Echoes of Gladiator Stance",
		Duration: duration,
	})

	// We're assuming these will be exclusive but TBD
	warrior.newStanceOverrideExclusiveEffect(BattleStance, battleStanceAura)
	warrior.newStanceOverrideExclusiveEffect(DefensiveStance, defStanceAura)
	warrior.newStanceOverrideExclusiveEffect(BerserkerStance, berserkStanceAura)
	warrior.newStanceOverrideExclusiveEffect(AnyStance, gladStanceAura)

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(StanceCodes) {
				switch warrior.PreviousStance {
				case BattleStance:
					battleStanceAura.Activate(sim)
				case DefensiveStance:
					defStanceAura.Activate(sim)
				case BerserkerStance:
					berserkStanceAura.Activate(sim)
				case GladiatorStance:
					gladStanceAura.Activate(sim)
				}
			}
		},
	}))
}

// For the first 10 sec after activating a stance, you can gain an additional benefit:
// Battle Stance/Gladiator Stance: 10% increased damage done.
// Berserker Stance: 10% increased critical strike chance.
// Defensive Stance: 10% reduced Physical damage taken.
func (warrior *Warrior) applyT1Damage6PBonus() {
	label := "S03 - Item - T1 - Warrior - Damage 6P Bonus"
	if warrior.HasAura(label) {
		return
	}

	duration := time.Second * 15

	battleAura := warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 457816},
		Label:    "Battle Forecast",
		Duration: duration,
	}).AttachMultiplicativePseudoStatBuff(&warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical], 1.05)

	defenseAura := warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 457814},
		Label:    "Defense Forecast",
		Duration: duration,
	}).AttachMultiplicativePseudoStatBuff(&warrior.PseudoStats.DamageTakenMultiplier, 0.95)

	berserkAura := warrior.NewTemporaryStatsAura(
		"Berserker Forecast",
		core.ActionID{SpellID: 457817},
		stats.Stats{stats.MeleeCrit: 5 * core.CritRatingPerCritChance},
		duration,
	)

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			switch spell.ClassSpellMask {
			case ClassSpellMask_WarriorStanceBattle:
				battleAura.Activate(sim)
			case ClassSpellMask_WarriorStanceGladiator:
				battleAura.Activate(sim)
			case ClassSpellMask_WarriorStanceDefensive:
				defenseAura.Activate(sim)
			case ClassSpellMask_WarriorStanceBerserker:
				berserkAura.Activate(sim)
			}
		},
	}))
}

var ItemSetImmoveableMight = core.NewItemSet(core.ItemSet{
	Name: "Immoveable Might",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases the block value of your shield by 30.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.BlockValue, 30)
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT1Tank4PBonus()
		},
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT1Tank6PBonus()
		},
	},
})

// You gain 1 extra Rage every time you take any damage or deal auto attack damage.
func (warrior *Warrior) applyT1Tank4PBonus() {
	label := "S03 - Item - T1 - Warrior - Tank 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	warrior.RegisterAura(core.Aura{
		Label: label,
	})

	warrior.AddDamageDealtRageBonus(1)
	warrior.AddDamageTakenRageBonus(1)
}

// Increases all threat you generate in Defensive Stance by an additional 10% and increases all damage you deal in Gladiator Stance by 4%.
func (warrior *Warrior) applyT1Tank6PBonus() {
	label := "S03 - Item - T1 - Warrior - Tank 6P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.defensiveStanceThreatMultiplier, 1.10,
	).AttachMultiplicativePseudoStatBuff(
		&warrior.gladiatorStanceDamageMultiplier, 1.04,
	))
}
