package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetDreadnaughtsWarplate = core.NewItemSet(core.ItemSet{
	Name: "Dreadnaught's Warplate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyNaxxramasDamage2PBonus()
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyNaxxramasDamage4PBonus()
		},
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyNaxxramasDamage6PBonus()
		},
	},
})

// Increases damage done by your Deep Wounds talent by 20%.
func (warrior *Warrior) applyNaxxramasDamage2PBonus() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	label := "S03 - Item - Naxxramas - Warrior - Damage 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warrior.DeepWounds.DamageMultiplierAdditive += 0.20
		},
	})
}

// Reduces the cooldown on your Bloodthirst, Mortal Strike, and Shield Slam abilities by 25%.
func (warrior *Warrior) applyNaxxramasDamage4PBonus() {
	if !warrior.Talents.Bloodthirst && !warrior.Talents.MortalStrike && !warrior.Talents.ShieldSlam {
		return
	}

	label := "S03 - Item - Naxxramas - Warrior - Damage 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if warrior.Bloodthirst != nil {
				warrior.Bloodthirst.CD.Multiplier *= 0.75
			}
			if warrior.MortalStrike != nil {
				warrior.MortalStrike.CD.Multiplier *= 0.75
			}
			if warrior.ShieldSlam != nil {
				warrior.ShieldSlam.CD.Multiplier *= 0.75
			}
		},
	})
}

// Your melee critical strikes against Undead enemies grant you 1% increased damage done to Undead for 30 sec, stacking up to 25 times.
func (warrior *Warrior) applyNaxxramasDamage6PBonus() {
	label := "S03 - Item - Naxxramas - Warrior - Damage 6P Bonus"
	if warrior.HasAura(label) {
		return
	}

	buffAura := warrior.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1219485},
		Label:     "Undead Slaying",
		Duration:  time.Second * 30,
		MaxStacks: 25,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			warrior.PseudoStats.DamageDealtMultiplier /= 1 + 0.01*float64(oldStacks)
			warrior.PseudoStats.DamageDealtMultiplier *= 1 + 0.01*float64(newStacks)
		},
	})

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee) && result.DidCrit() && result.Target.MobType == proto.MobType_MobTypeUndead {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}

var ItemSetDreadnaughtsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Dreadnaught's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyNaxxramasProtection2PBonus()
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyNaxxramasProtection4PBonus()
		},
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyNaxxramasProtection6PBonus()
		},
	},
})

// Your Taunt ability never misses, and your chance to be Dodged or Parried is reduced by 2%.
func (warrior *Warrior) applyNaxxramasProtection2PBonus() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	label := "S03 - Item - Naxxramas - Warrior - Protection 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{stats.Expertise: 2 * core.ExpertiseRatingPerExpertiseChance}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachBuildPhaseStatsBuff(bonusStats))
}

// Reduces the cooldown on your Shield Wall ability by 3 min and reduces the cooldown on your Recklessness ability by 3 min.
// Recklessness can now be used in any Stance and does not increase damage taken.
func (warrior *Warrior) applyNaxxramasProtection4PBonus() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	label := "S03 - Item - Naxxramas - Warrior - Protection 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if warrior.ShieldWall != nil {
				warrior.ShieldWall.CD.FlatModifier -= time.Minute * 3
			}
			if warrior.Recklessness != nil {
				warrior.Recklessness.CD.FlatModifier -= time.Minute * 3
				warrior.Recklessness.StanceMask = AnyStance
				warrior.recklessnessDamageTakenMultiplier = 1
			}
		},
	})
}

// When you take damage from an Undead enemy, the remaining duration of your active Last Stand is reset to 20 sec.
func (warrior *Warrior) applyNaxxramasProtection6PBonus() {
	if !warrior.Talents.LastStand {
		return
	}

	label := "S03 - Item - Naxxramas - Warrior - Protection 6P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if warrior.LastStandAura.IsActive() && spell.Unit.MobType == proto.MobType_MobTypeUndead && result.Landed() && result.Damage > 0 {
				warrior.LastStandAura.Activate(sim)
			}
		},
	}))
}
