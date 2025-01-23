package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetBeastmasterArmor = core.NewItemSet(core.ItemSet{
	Name: "Beastmaster Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Your melee and ranged autoattacks have a 6% chance to energize you for 300 mana.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450577}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "S03 - Mana Proc on Cast - Beaststalker Armor",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskWhiteHit,
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

var ItemSetGiantstalkerProwess = core.NewItemSet(core.ItemSet{
	Name: "Giantstalker Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT1Melee2PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT1Melee4PBonus()
		},
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT1Melee6PBonus()
		},
	},
})

// Your Mongoose Bite also reduces its target's chance to Dodge by 1% and increases your chance to hit by 1% for 30 sec.
func (hunter *Hunter) applyT1Melee2PBonus() {
	label := "S03 - Item - T1 - Hunter - Melee 2P Bonus"
	if hunter.HasAura(label) {
		return
	}

	procBonus := stats.Stats{
		stats.SpellHit: 1,
		stats.MeleeHit: 1,
	}

	stalkerAura := hunter.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 458403},
		Label:    "Stalker",
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatsDynamic(sim, procBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatsDynamic(sim, procBonus.Invert())
		},
	})

	debuffAuras := hunter.NewEnemyAuraArray(core.MeleeHunterDodgeReductionAura)
	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_HunterMongooseBite) && result.Landed() {
				debuffAuras.Get(result.Target).Activate(sim)
				stalkerAura.Activate(sim)
			}
		},
	}))
}

// While tracking a creature type, you deal 3% increased damage to that creature type.
// Unsure if this stacks with the Pursuit 4p
func (hunter *Hunter) applyT1Melee4PBonus() {
	label := "S03 - Item - T1 - Hunter - Melee 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Just adding 3% damage to assume the hunter is tracking their target's type
			hunter.PseudoStats.DamageDealtMultiplier *= 1.03
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.DamageDealtMultiplier /= 1.03
		},
	}))
}

// Mongoose Bite also activates for 5 sec whenever your target Parries or Blocks or when your melee attack misses.
func (hunter *Hunter) applyT1Melee6PBonus() {
	label := "S03 - Item - T1 - Hunter - Melee 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee) && !result.Landed() {
				hunter.DefensiveState.Activate(sim)
			}
		},
	}))
}

var ItemSetGiantstalkerPursuit = core.NewItemSet(core.ItemSet{
	Name: "Giantstalker Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		// You generate 100% more threat for 8 sec after using Distracting Shot.
		2: func(agent core.Agent) {
			// Do nothing
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT1Ranged4PBonus()
		},
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT1Ranged6PBonus()
		},
	},
})

// While tracking a creature type, you deal 3% increased damage to that creature type.
// Unsure if this stacks with the Prowess 4p
func (hunter *Hunter) applyT1Ranged4PBonus() {
	label := "S03 - Item - T1 - Hunter - Ranged 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Just adding 3% damage to assume the hunter is tracking their target's type
			hunter.PseudoStats.DamageDealtMultiplier *= 1.03
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.DamageDealtMultiplier /= 1.03
		},
	}))
}

// Your next Shot ability within 12 sec after Aimed Shot deals 20% more damage.
func (hunter *Hunter) applyT1Ranged6PBonus() {
	if !hunter.Talents.AimedShot {
		return
	}

	label := "S03 - Item - T1 - Hunter - Ranged 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	damageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  ClassSpellMask_HunterShots,
		FloatValue: 0.20,
	})

	procAura := hunter.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 456382},
		Label:    "Precision",
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(ClassSpellMask_HunterShots) || (aura.RemainingDuration(sim) == aura.Duration && spell.Matches(ClassSpellMask_HunterAimedShot)) {
				return
			}

			aura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
		Name:           label,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: ClassSpellMask_HunterAimedShot,
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			procAura.Activate(sim)
		},
	})
}
