package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetRaimentsOfFaith = core.NewItemSet(core.ItemSet{
	Name: "Raiments of Faith",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyNaxxramasShadow2PBonus()
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyNaxxramasShadow4PBonus()
		},
		6: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyNaxxramasShadow6PBonus()
		},
	},
})

// Your Shadow Word: Pain ability deals 20% more damage.
func (priest *Priest) applyNaxxramasShadow2PBonus() {
	label := "S03 - Item - Naxxramas - Priest - Shadow 2P Bonus"
	if priest.HasAura(label) {
		return
	}

	priest.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range priest.ShadowWordPain {
				if spell == nil {
					continue
				}

				spell.DamageMultiplierAdditive += 0.20
			}
		},
	})
}

// Reduces the cooldown on your Mind Blast ability by 1.0 sec.
func (priest *Priest) applyNaxxramasShadow4PBonus() {
	label := "S03 - Item - Naxxramas - Priest - Shadow 4P Bonus"
	if priest.HasAura(label) {
		return
	}

	priest.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range priest.MindBlast {
				if spell == nil {
					continue
				}

				spell.CD.FlatModifier -= time.Second
			}
		},
	})
}

// Your Mind Flay, Mind Blast, and Mind Spike abilities deal increased damage to Undead targets equal to their critical strike chance.
func (priest *Priest) applyNaxxramasShadow6PBonus() {
	label := "S03 - Item - Naxxramas - Priest - Shadow 6P Bonus"
	if priest.HasAura(label) {
		return
	}

	priest.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells := core.Flatten([][]*core.Spell{
				priest.MindBlast,
				{priest.MindSpike},
			})

			if priest.HasRune(proto.PriestRune_RuneBracersDespair) {
				affectedSpells = append(affectedSpells, core.Flatten(priest.MindFlay)...)
			}

			for _, spell := range affectedSpells {
				if spell == nil {
					continue
				}

				oldApplyEffects := spell.ApplyEffects
				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					critChanceBonus := 0.0
					if target.MobType == proto.MobType_MobTypeUndead {
						critChanceBonus = priest.GetStat(stats.SpellCrit) / 100
					}

					spell.DamageMultiplierAdditive += critChanceBonus
					oldApplyEffects(sim, target, spell)
					spell.DamageMultiplierAdditive -= critChanceBonus
				}
			}
		},
	})
}

var ItemSetVestmentsOfFaith = core.NewItemSet(core.ItemSet{
	Name: "Vestments of Faith",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyNaxxramasHealer2PBonus()
		},
		// Your Penance, Flash Heal Rank 7, and Greater Heal Rank 4 and Rank 5 have a 9% chance to grant the target 10% increased critical strike chance for 15 sec.
		4: func(agent core.Agent) {
		},
		// Your Power Word: Shield has a 50% chance to not deplete when the target is damaged by an Undead enemy.
		6: func(agent core.Agent) {
		},
	},
})

// Reduces the cooldown on your Circle of Healing and Penance abilities by 25%.
func (priest *Priest) applyNaxxramasHealer2PBonus() {
	if !priest.HasRune(proto.PriestRune_RuneHandsPenance) {
		return
	}

	label := "S03 - Item - Naxxramas - Priest - Healer 2P Bonus"
	if priest.HasAura(label) {
		return
	}

	priest.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			priest.Penance.CD.Multiplier -= 25
			priest.PenanceHeal.CD.Multiplier -= 25
		},
	})
}
