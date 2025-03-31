package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var ItemSetConquerorsAdvance = core.NewItemSet(core.ItemSet{
	Name: "Conqueror's Advance",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyTAQDamage2PBonus()
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyTAQDamage4PBonus()
		},
	},
})

// Reduces the cooldown on your Death Wish by 50%.
func (warrior *Warrior) applyTAQDamage2PBonus() {
	if !warrior.Talents.DeathWish {
		return
	}

	label := "S03 - Item - TAQ - Warrior - Damage 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "S03 - Item - TAQ - Warrior - Damage 2P Bonus",
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Multi_Flat,
		ClassMask: ClassSpellMask_WarriorDeathWish,
		IntValue:  -50,
	})
}

// You deal 15% increased damage while any nearby enemy is afflicted with both your Rend and your Deep Wounds.
func (warrior *Warrior) applyTAQDamage4PBonus() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	label := "S03 - Item - TAQ - Warrior - Damage 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	buffAura := warrior.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1214166},
		Label:    "Bloodythirsty",
		Duration: time.Second * 3,
	}).AttachMultiplicativePseudoStatBuff(&warrior.PseudoStats.DamageDealtMultiplier, 1.15)

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.Matches(ClassSpellMask_WarriorDeepWounds) && warrior.Rend.Dot(result.Target).IsActive()) ||
				(spell.Matches(ClassSpellMask_WarriorRend) && warrior.DeepWounds.Dot(result.Target).IsActive()) {
				buffAura.Activate(sim)
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.Matches(ClassSpellMask_WarriorDeepWounds) && warrior.Rend.Dot(result.Target).IsActive()) ||
				(spell.Matches(ClassSpellMask_WarriorRend) && warrior.DeepWounds.Dot(result.Target).IsActive()) {
				buffAura.Activate(sim)
			}
		},
	}))
}

var ItemSetConquerorsBulwark = core.NewItemSet(core.ItemSet{
	Name: "Conqueror's Bulwark",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyTAQTank2PBonus()
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyTAQTank4PBonus()
		},
	},
})

// Reduces the cooldown on Thunder Clap by 100%.
func (warrior *Warrior) applyTAQTank2PBonus() {
	label := "S03 - Item - TAQ - Warrior - Tank 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warrior.ThunderClap.CD.Duration = 0
		},
	})
}

// Your Shield Slam deals 100% increased threat and its cooldown is reset if it is Dodged, Parried, or Blocked.
func (warrior *Warrior) applyTAQTank4PBonus() {
	if !warrior.Talents.ShieldSlam {
		return
	}

	label := "S03 - Item - TAQ - Warrior - Tank 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_WarriorShieldSlam) && result.Outcome.Matches(core.OutcomeDodge|core.OutcomeParry|core.OutcomeBlock) {
				spell.CD.Reset()
			}
		},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_Threat_Flat,
		ClassMask:  ClassSpellMask_WarriorShieldSlam,
		FloatValue: 1,
	})
}

var ItemSetBattlegearOfUnyieldingStrength = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Unyielding Strength",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyRAQTank3PBonus()
		},
	},
})

// Reduces the cooldown on Shockwave by 50%.
func (warrior *Warrior) applyRAQTank3PBonus() {
	if !warrior.HasRune(proto.WarriorRune_RuneShockwave) {
		return
	}

	label := "S03 - Item - RAQ - Warrior - Tank 3P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "S03 - Item - RAQ - Warrior - Tank 3P Bonus",
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Multi_Flat,
		ClassMask: ClassSpellMask_WarriorShockwave,
		IntValue:  -50,
	})
}
