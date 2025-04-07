package priest

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var ItemSetRaimentsOfRevalation = core.NewItemSet(core.ItemSet{
	Name: "Raiments of Revelation",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyScarletEnclaveShadow2PBonus()
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyScarletEnclaveShadow4PBonus()
		},
		6: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyScarletEnclaveShadow6PBonus()
		},
	},
})

// Your Mind Flay and Mind Sear no longer lose duration from taking damage during their channel.
// In addition, they deal 10% increased damage per other periodic Shadow effect you have on the target, up to a maximum increase of 30%.
func (priest *Priest) applyScarletEnclaveShadow2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Priest - Shadow 2P Bonus"
	if priest.HasAura(label) {
		return
	}

	classMask := ClassSpellMask_PriestMindFlay | ClassSpellMask_PriestMindSear

	damageMod := priest.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  classMask,
		FloatValue: 1.0,
	})

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			mindFlay := core.FilterSlice(
				core.Flatten(priest.MindFlay),
				func(spell *core.Spell) bool { return spell != nil },
			)

			for _, spell := range mindFlay {
				spell.PushbackReduction = 1
			}
		},
		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.Matches(classMask) {
				multiplier := 1.0

				for _, spell := range priest.DoTSpells {
					if spell.Dot(target).IsActive() {
						multiplier += 0.1
					}
				}

				damageMod.Activate()
				damageMod.UpdateFloatValue(min(1.3, multiplier))
			}
		},
	}))
}

// Your Mind Blast deals 50% reduced threat, and gains 20% damage increase from each stack of Mind Spike on the target.
func (priest *Priest) applyScarletEnclaveShadow4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Priest - Shadow 4P Bonus"
	if priest.HasAura(label) {
		return
	}

	priest.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_Threat_Pct,
		ClassMask:  ClassSpellMask_PriestMindBlast,
		FloatValue: -0.50,
	})

	damageMod := priest.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  ClassSpellMask_PriestMindBlast,
		FloatValue: 1.0,
	})

	aura := core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
	}))

	if priest.HasRune(proto.PriestRune_RuneWaistMindSpike) {
		aura.ApplyOnApplyEffects(func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_PriestMindBlast) {
				debuff := priest.MindSpike.RelatedAuras[0].Get(target)
				damageMod.Activate()
				damageMod.UpdateFloatValue(core.TernaryFloat64(debuff.IsActive(), 1+0.20*float64(debuff.GetStacks()), 1))
			}
		})
	}
}

// Damage done by your Mind Flay now increases the longer you channel the spell.
// Each time it deals damage, subsequent damage will increase by 70%.
// This resets on each new channel.
func (priest *Priest) applyScarletEnclaveShadow6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Priest - Shadow 6P Bonus"
	if priest.HasAura(label) {
		return
	}

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_PriestMindFlay) {
				spell.Dot(result.Target).SnapshotAttackerMultiplier *= 1.70
			}
		},
	}))
}

var ItemSetVestmentsOfRevalation = core.NewItemSet(core.ItemSet{
	Name: "Vestments of Revelation",
	Bonuses: map[int32]core.ApplyEffect{
		// Each time you cast Circle of Healing or Penance, you have a 50% chance to make your next Lesser Heal, Heal, Greater Heal, Flash Heal, Binding Heal, or Prayer of Healing within 15 sec instant cast.
		2: func(agent core.Agent) {
		},
		// Your Power Word: Shield also instantly heals the target for 20% of the absorb value.
		4: func(agent core.Agent) {
		},
		// Your Power Word: Shield has a 50% chance to not deplete when the target is damaged by an Undead enemy.
		6: func(agent core.Agent) {
		},
	},
})
