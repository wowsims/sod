package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetRedemptionWarplate = core.NewItemSet(core.ItemSet{
	Name: "Redemption Warplate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.applyNaxxramasRetribution2PBonus()
		},
		4: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.applyNaxxramasRetribution4PBonus()
		},
		6: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.applyNaxxramasRetribution6PBonus()
		},
	},
})

// Increases the damage done by your Divine Storm ability by 100%.
func (paladin *Paladin) applyNaxxramasRetribution2PBonus() {
	if !paladin.hasRune(proto.PaladinRune_RuneChestDivineStorm) {
		return
	}

	label := "S03 - Item - Naxxramas - Paladin - Retribution 2P Bonus"
	if paladin.HasAura(label) {
		return
	}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PaladinDivineStorm,
		IntValue:  100,
	}))
}

// Reduces the cast time of your Holy Wrath ability by 100% and reduces the cooldown and mana cost of your Holy Wrath ability by 75%.
func (paladin *Paladin) applyNaxxramasRetribution4PBonus() {
	label := "S03 - Item - Naxxramas - Paladin - Retribution 4P Bonus"
	if paladin.HasAura(label) {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range paladin.holyWrath {
				spell.CastTimeMultiplier -= 1
				spell.CD.Multiplier *= 0.25
				spell.Cost.Multiplier -= 75
			}
		},
	})
}

// Your Crusader Strike, Divine Storm, Exorcism and Holy Wrath abilities deal increased damage to Undead equal to their critical strike chance.
func (paladin *Paladin) applyNaxxramasRetribution6PBonus() {
	label := "S03 - Item - Naxxramas - Paladin - Retribution 6P Bonus"
	if paladin.HasAura(label) {
		return
	}

	classSpellMasks := ClassSpellMask_PaladinExorcism | ClassSpellMask_PaladinHolyWrath | ClassSpellMask_PaladinDivineStorm | ClassSpellMask_PaladinCrusaderStrike
	damageMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  classSpellMasks,
		FloatValue: 1,
	})

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label: label,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if !spell.Matches(classSpellMasks) {
				return
			}
			critChanceBonus := 0.0

			if target.MobType == proto.MobType_MobTypeUndead {
				if spell.Matches(ClassSpellMask_PaladinExorcism | ClassSpellMask_PaladinHolyWrath) {
					critChanceBonus = paladin.GetStat(stats.SpellCrit)/100.0 + paladin.GetSchoolBonusCritChance(spell)/100.0
				} else {
					critChanceBonus = paladin.GetStat(stats.MeleeCrit) / 100.0
				}
			}
			critChanceBonus = min(critChanceBonus, 1)
			damageMod.UpdateFloatValue(1 + critChanceBonus)
		},
	}))
}

var ItemSetRedemptionBulwark = core.NewItemSet(core.ItemSet{
	Name: "Redemption Bulwark",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.applyNaxxramasProtection2PBonus()
		},
		4: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.applyNaxxramasProtection4PBonus()
		},
		6: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.applyNaxxramasProtection6PBonus()
		},
	},
})

// Your Hand of Reckoning ability never misses, and your chance to be Dodged or Parried is reduced by 2%.
func (paladin *Paladin) applyNaxxramasProtection2PBonus() {
	label := "S03 - Item - Naxxramas - Paladin - Protection 2P Bonus"
	if paladin.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{stats.Expertise: 2 * core.ExpertiseRatingPerExpertiseChance}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachStatsBuff(bonusStats))
}

// Reduces the cooldown on your Divine Protection ability by 3 min and reduces the cooldown on your Avenging Wrath ability by 2 min.
func (paladin *Paladin) applyNaxxramasProtection4PBonus() {
	label := "S03 - Item - Naxxramas - Paladin - Protection 4P Bonus"
	if paladin.HasAura(label) {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.divineProtection != nil {
				paladin.divineProtection.CD.FlatModifier -= time.Minute * 3
			}

			if paladin.avengingWrath != nil {
				paladin.avengingWrath.CD.FlatModifier -= time.Minute * 2
			}
		},
	})
}

// When damage from an Undead enemy takes you below 35% health, the effect from Hand of Reckoning and Righteous Fury now reduces that damage by 50%.
func (paladin *Paladin) applyNaxxramasProtection6PBonus() {
	label := "S03 - Item - Naxxramas - Paladin - Protection 6P Bonus"
	if paladin.HasAura(label) {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// Implemented in righteous_fury.go
		},
	})
}

var ItemSetRedemptionArmor = core.NewItemSet(core.ItemSet{
	Name: "Redemption Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			paladin.applyNaxxramasHoly2PBonus()
		},
		// Your Flash of Light Rank 6 and Holy Light Rank 8 and Rank 9 spells have a 10% chance to imbue your target with Holy Power.
		4: func(agent core.Agent) {
		},
		// Your Beacon of Light target takes 20% reduced damage from Undead enemies.
		6: func(agent core.Agent) {
		},
	},
})

// Reduces the cooldown on your Lay on Hands ability by 35 min, and your Lay on Hands now restores you to 30% of your maximum Mana when used.
func (paladin *Paladin) applyNaxxramasHoly2PBonus() {
	label := "S03 - Item - Naxxramas - Paladin - Holy 2P Bonus"
	if paladin.HasAura(label) {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			paladin.layOnHands.CD.FlatModifier -= time.Minute * 35

			// TODO: Mana return
		},
	})
}
