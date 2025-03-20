package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetDragonstalkerProwess = core.NewItemSet(core.ItemSet{
	Name: "Dragonstalker's Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT2Melee2PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT2Melee4PBonus()
		},
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT2Melee6PBonus()
		},
	},
})

// Raptor Strike increases the damage done by your next other melee ability within 5 sec by 20%.
func (hunter *Hunter) applyT2Melee2PBonus() {
	label := "S03 - Item - T2 - Hunter - Melee 2P Bonus"
	if hunter.HasAura(label) {
		return
	}

	procMask := core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial
	affectedSpells := ClassSpellMask_HunterAll ^ ClassSpellMask_HunterRaptorStrikeHit ^ ClassSpellMask_HunterRaptorStrike ^ ClassSpellMask_HunterWingClip

	damageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ProcMask:  procMask,
		ClassMask: affectedSpells,
		IntValue:  20,
	})

	procAura := hunter.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467331},
		Label:    "Clever Strikes",
		Duration: time.Second * 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !(spell.Matches(affectedSpells) && spell.ProcMask.Matches(procMask)) {
				return
			}

			aura.Deactivate(sim)
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_HunterRaptorStrikeHit) {
				procAura.Activate(sim)
			}
		},
	}))
}

// Increases damage dealt by your main hand weapon with Raptor Strike and Wyvern Strike by 20%.
func (hunter *Hunter) applyT2Melee4PBonus() {
	label := "S03 - Item - T2 - Hunter - Melee 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterRaptorStrike | ClassSpellMask_HunterRaptorStrikeHit | ClassSpellMask_HunterWyvernStrike,
		ProcMask:  core.ProcMaskMeleeMHSpecial,
		IntValue:  20,
	}))
}

// Your periodic damage has a 5% chance to reset the cooldown on one of your Strike abilities.
// The Strike with the longest remaining cooldown is always chosen.
func (hunter *Hunter) applyT2Melee6PBonus() {
	label := "S03 - Item - T2 - Hunter - Melee 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467334}, // Tracking in APL
		Label:    label,
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(0.05, "T2 Melee 6PC Strike Reset") {
				maxSpell := hunter.RaptorStrike

				for _, strike := range hunter.Strikes {
					if strike.TimeToReady(sim) > maxSpell.TimeToReady(sim) {
						maxSpell = strike
					}
				}

				maxSpell.CD.Reset()
				aura.Activate(sim) // used for metrics
			}
		},
	}))
}

var ItemSetDragonstalkerPursuit = core.NewItemSet(core.ItemSet{
	Name: "Dragonstalker's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT2Ranged2PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT2Ranged4PBonus()
		},
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyT2Ranged6PBonus()
		},
	},
})

// Your Aimed Shot deals 20% more damage to targets afflicted by one of your trap effects.
func (hunter *Hunter) applyT2Ranged2PBonus() {
	if !hunter.Talents.AimedShot {
		return
	}

	label := "S03 - Item - T2 - Hunter - Ranged 2P Bonus"
	if hunter.HasAura(label) {
		return
	}

	damageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterAimedShot,
	})

	core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: ClassSpellMask_HunterAimedShot,
		Callback:       core.CallbackOnApplyEffects,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			hasActiveTrap := result.Target.HasActiveAuraWithTag("ImmolationTrap") || hunter.HasActiveAuraWithTag("ExplosiveTrap")
			damageMod.UpdateIntValue(core.TernaryInt64(hasActiveTrap, 20, 0))
			damageMod.Activate()
		},
	})
}

// Your damaging Shot abilities deal 10% increased damage if the previous damaging Shot used was different than the current one.
func (hunter *Hunter) applyT2Ranged4PBonus() {
	label := "S03 - Item - T2 - Hunter - Ranged 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	// Map of all possible damage mods
	damageModMap := make(map[uint64]*core.SpellMod)
	// Dynamic damage mod for the current shot
	var damageMod *core.SpellMod

	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_HunterShots) {
			damageModMap[spell.ClassSpellMask] = hunter.AddDynamicMod(core.SpellModConfig{
				Kind:      core.SpellMod_DamageDone_Flat,
				ClassMask: ClassSpellMask_HunterShots &^ spell.ClassSpellMask,
				IntValue:  10,
			})
		}
	})

	procAura := hunter.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467312},
		Label:    label + " Proc",
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
	})

	core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: ClassSpellMask_HunterShots,
		Callback:       core.CallbackOnCastComplete,
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if procAura.IsActive() {
				procAura.Deactivate(sim)
			}
			damageMod = damageModMap[spell.ClassSpellMask]
			procAura.Activate(sim)
		},
	})
}

// Your Serpent Sting damage is increased by 25% of your Attack Power over its normal duration.
func (hunter *Hunter) applyT2Ranged6PBonus() {
	label := "S03 - Item - T2 - Hunter - Ranged 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			hunter.SerpentStingAPCoeff += 0.25
		},
	})
}

var ItemSetPredatorArmor = core.NewItemSet(core.ItemSet{
	Name: "Predator's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Attack Power.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
			c.AddStat(stats.RangedAttackPower, 20)
		},
		3: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyZGBeastmaster3PBonus()
		},
		5: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyZGBeastmaster5PBonus()
		},
	},
})

// Increases the Attack Power your Beast pet gains from your attributes by 20%.
func (hunter *Hunter) applyZGBeastmaster3PBonus() {
	if hunter.pet == nil {
		return
	}

	label := "S03 - Item - ZG - Hunter - Beastmaster 3P Bonus"
	if hunter.HasAura(label) {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			oldStatInheritance := hunter.pet.GetStatInheritance()
			hunter.pet.UpdateStatInheritance(
				func(ownerStats stats.Stats) stats.Stats {
					s := oldStatInheritance(ownerStats)
					s[stats.AttackPower] *= 1.20
					return s
				},
			)
		},
	})
}

// Increases the Focus regeneration of your Beast pet by 20%.
func (hunter *Hunter) applyZGBeastmaster5PBonus() {
	if hunter.pet == nil {
		return
	}

	label := "S03 - Item - ZG - Hunter - Beastmaster 5P Bonus"
	if hunter.HasAura(label) {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			hunter.pet.AddFocusRegenMultiplier(0.20)
		},
	})
}
