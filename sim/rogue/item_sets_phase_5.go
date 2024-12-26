package rogue

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetBloodfangThrill = core.NewItemSet(core.ItemSet{
	Name: "Bloodfang Thrill",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyT2Damage2PBonus()
		},
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyT2Damage4PBonus()
		},
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyT2Damage6PBonus()
		},
	},
})

// Your opening moves have a 100% chance to make your next ability cost no energy.
func (rogue *Rogue) applyT2Damage2PBonus() {
	label := "S03 - Item - T2 - Rogue - Damage 2P Bonus"
	if rogue.HasAura(label) {
		return
	}

	var affectedSpells []*core.Spell

	clearcastingAura := rogue.RegisterAura(core.Aura{
		Label:    fmt.Sprintf("Clearcasting (%s)", label),
		ActionID: core.ActionID{SpellID: 467735},
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

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.SpellCode == SpellCode_RogueAmbush || spell.SpellCode == SpellCode_RogueGarrote) && result.Landed() {
				clearcastingAura.Activate(sim)
			}
		},
	}))
}

// Increases damage dealt by your main hand weapon from combo-generating abilities by 20%
func (rogue *Rogue) applyT2Damage4PBonus() {
	label := "S03 - Item - T2 - Rogue - Damage 4P Bonus"
	if rogue.HasAura(label) {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range rogue.Spellbook {
				if spell.Flags.Matches(SpellFlagBuilder) && (spell.ProcMask.Matches(core.ProcMaskMeleeMHSpecial) || spell.SpellCode == SpellCode_RogueMainGauche || spell.SpellCode == SpellCode_RoguePoisonedKnife) {
					spell.DamageMultiplierAdditive += 0.20
				}
			}
		},
	})
}

// Reduces cooldown on vanish to 1 minute
func (rogue *Rogue) applyT2Damage6PBonus() {
	label := "S03 - Item - T2 - Rogue - Damage 6P Bonus"
	if rogue.HasAura(label) {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			//Applied after talents in sim so does not stack with elusiveness when active.
			rogue.Vanish.CD.Duration = time.Second * 60
		},
	})
}

var ItemSetBloodfangBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Bloodfang Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Rolling with the Punches now also activates every time you gain a combo point.
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			if !rogue.HasRune(proto.RogueRune_RuneRollingWithThePunches) {
				return
			}
			rogue.OnComboPointsGained(func(sim *core.Simulation) {
				rogue.RollingWithThePunchesProcAura.Activate(sim)
				rogue.RollingWithThePunchesProcAura.AddStack(sim)
			})
		},
		// Your Rolling with the Punches also grants you 20% increased Armor from items per stack.
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			if !rogue.HasRune(proto.RogueRune_RuneRollingWithThePunches) {
				return
			}
			initarmor := rogue.BaseEquipStats()[stats.Armor]

			fmt.Println(initarmor)

			rogue.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Rogue - Tank 4P Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					oldOnStacksChange := rogue.RollingWithThePunchesProcAura.OnStacksChange
					rogue.RollingWithThePunchesProcAura.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
						oldOnStacksChange(aura, sim, oldStacks, newStacks)
						rogue.AddStatDynamic(sim, stats.Armor, float64(0.2*initarmor*float64(newStacks-oldStacks)))
					}
				},
			})
		},
		// The cooldown on your Main Gauche resets every time your target Dodges or Parries.
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			if !rogue.HasRune(proto.RogueRune_RuneMainGauche) {
				return
			}

			core.MakePermanent(rogue.RegisterAura(core.Aura{
				Label: "S03 - Item - T2 - Rogue - Tank 6P Bonus",
				OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidDodge() || result.DidParry() {
						rogue.MainGauche.CD.Reset()
					}
				},
			}))
		},
	},
})

var ItemSetMadCapsOutfit = core.NewItemSet(core.ItemSet{
	Name: "Madcap's Outfit",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Attack Power
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       20,
				stats.RangedAttackPower: 20,
			})
		},
		// Increases your chance to get a critical strike with Daggers by 5%.
		3: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			switch rogue.GetProcMaskForTypes(proto.WeaponType_WeaponTypeDagger) {
			case core.ProcMaskMelee:
				rogue.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*float64(5))
			case core.ProcMaskMeleeMH:
				// the default character pane displays critical strike chance for main hand only
				rogue.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*float64(5))
				rogue.OnSpellRegistered(func(spell *core.Spell) {
					if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
						spell.BonusCritRating -= core.CritRatingPerCritChance * float64(5)
					}
				})
			case core.ProcMaskMeleeOH:
				rogue.OnSpellRegistered(func(spell *core.Spell) {
					if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
						spell.BonusCritRating += core.CritRatingPerCritChance * float64(5)
					}
				})
			}
		},
		// Increases the critical strike chance of your Ambush ability by 30%.
		5: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.OnSpellRegistered(func(spell *core.Spell) {
				if spell.SpellCode == SpellCode_RogueAmbush {
					spell.BonusCritRating += 30 * core.CritRatingPerCritChance
				}
			})
		},
	},
})
