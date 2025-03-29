package warrior

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetUnstoppableWrath = core.NewItemSet(core.ItemSet{
	Name: "Unstoppable Wrath",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT2Damage2PBonus()
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT2Damage4PBonus()
		},
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT2Damage6PBonus()
		},
	},
})

// Overpower critical strikes refresh the duration of Rend on your target back to its maximum duration.
func (warrior *Warrior) applyT2Damage2PBonus() {
	label := "S03 - Item - T2 - Warrior - Damage 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_WarriorOverpower) && result.DidCrit() {
				if dot := warrior.Rend.Dot(result.Target); dot.IsActive() {
					dot.Rollover(sim)
				}
			}
		},
	}))
}

// Increases the damage of Heroic Strike, Overpower, and Slam by 25%
func (warrior *Warrior) applyT2Damage4PBonus() {
	label := "S03 - Item - T2 - Warrior - Damage 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warrior.HeroicStrike.ApplyMultiplicativeDamageBonus(1.25)
			warrior.Overpower.ApplyMultiplicativeDamageBonus(1.25)
			if warrior.SlamMH != nil {
				warrior.SlamMH.ApplyMultiplicativeDamageBonus(1.25)
			}
			if warrior.SlamOH != nil {
				warrior.SlamMH.ApplyMultiplicativeDamageBonus(1.25)
			}
			if warrior.QuickStrike != nil {
				warrior.QuickStrike.ApplyMultiplicativeDamageBonus(1.25)
			}
		},
	})
}

// Your Slam hits reset the remaining cooldown on your Mortal Strike, Bloodthirst, and Shield Slam abilities.
func (warrior *Warrior) applyT2Damage6PBonus() {
	label := "S03 - Item - T2 - Warrior - Damage 6P Bonus"
	if warrior.HasAura(label) {
		return
	}

	var affectedSpells []*core.Spell
	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = warrior.GetSpellsMatchingClassMask(ClassSpellMask_WarriorBloodthirst | ClassSpellMask_WarriorMortalStrike | ClassSpellMask_WarriorShieldSlam)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_WarriorSlamMH) && result.Landed() {
				for _, spell := range affectedSpells {
					spell.CD.Reset()
				}
			}
		},
	}))
}

var ItemSetImmoveableWrath = core.NewItemSet(core.ItemSet{
	Name: "Immoveable Wrath",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT2Protection2PBonus()
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT2Protection4PBonus()
		},
		6: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyT2Protection6PBonus()
		},
	},
})

// You gain 10 Rage every time you Parry or one of your attacks is Parried.
func (warrior *Warrior) applyT2Protection2PBonus() {
	label := "S03 - Item - T2 - Warrior - Protection 2P Bonus"
	if warrior.HasAura(label) {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 468066})

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee) && result.DidParry() {
				warrior.AddRage(sim, 10, rageMetrics)
			}
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidParry() {
				warrior.AddRage(sim, 10, rageMetrics)
			}
		},
	}))
}

// Revenge also grants you Flurry, increasing your attack speed by 30% for the next 3 swings.
func (warrior *Warrior) applyT2Protection4PBonus() {
	label := "S03 - Item - T2 - Warrior - Protection 4P Bonus"
	if warrior.HasAura(label) {
		return
	}

	flurryAura := warrior.makeFlurryAura(5)
	// The consumption trigger may not exist if the Warrior doesn't talent into Flurry
	warrior.makeFlurryConsumptionTrigger(flurryAura)

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_WarriorRevenge) {
				flurryAura.Activate(sim)
				flurryAura.SetStacks(sim, 3)
			}
		},
	}))
}

// When your target Parries an attack, you instantly Retaliate for 200% weapon damage to that target.
// Retaliate cannot be Dodged, Blocked, or Parried, but can only occur once every 30 sec per target.
func (warrior *Warrior) applyT2Protection6PBonus() {
	label := "S03 - Item - T2 - Warrior - Protection 6P Bonus"
	if warrior.HasAura(label) {
		return
	}

	retaliate := warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 468071},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial, // Retaliate and Retaliation count as normal yellow hits that can proc things
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		CritDamageBonus:  warrior.impale(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, warrior.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()), spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	icds := warrior.NewEnemyICDArray(func(u *core.Unit) *core.Cooldown {
		return &core.Cooldown{
			Timer:    warrior.NewTimer(),
			Duration: time.Second * 30,
		}
	})

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.DidParry() {
				return
			}

			if icd := icds.Get(result.Target); icd.IsReady(sim) {
				retaliate.Cast(sim, result.Target)
				icd.Use(sim)
			}
		},
	}))
}

var ItemSetVindicatorsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Vindicator's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +7.
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.AddStat(stats.Defense, 7)
		},
		3: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyZGGladiator3PBonus()
		},
		5: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			warrior.applyZGGladiator5PBonus()
		},
	},
})

// Reduces the cooldown on your Shield Slam ability by 2 sec.
func (warrior *Warrior) applyZGGladiator3PBonus() {
	if !warrior.Talents.ShieldSlam {
		return
	}

	label := "S03 - Item - ZG - Warrior - Gladiator 3P Bonus"
	if warrior.HasAura(label) {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: label,
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_WarriorShieldSlam,
		TimeValue: -time.Second * 2,
	})
}

// Reduces the cooldown on your Bloodrage ability by 30 sec while you are in Gladiator Stance.
func (warrior *Warrior) applyZGGladiator5PBonus() {
	if !warrior.HasRune(proto.WarriorRune_RuneGladiatorStance) {
		return
	}

	label := "S03 - Item - T1 - Warrior - Gladiator 5P Bonus"
	if warrior.HasAura(label) {
		return
	}

	cooldownMod := warrior.AddDynamicMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_WarriorBloodrage,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * 30,
	})

	warrior.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			idx := slices.IndexFunc(warrior.GladiatorStanceAura.ExclusiveEffects, func(ee *core.ExclusiveEffect) bool {
				return ee.Category.Name == stanceEffectCategory
			})
			ee := warrior.GladiatorStanceAura.ExclusiveEffects[idx]
			oldOnGain := ee.OnGain
			ee.OnGain = func(ee *core.ExclusiveEffect, sim *core.Simulation) {
				oldOnGain(ee, sim)
				cooldownMod.Activate()
			}

			oldOnExpire := ee.OnExpire
			ee.OnExpire = func(ee *core.ExclusiveEffect, sim *core.Simulation) {
				oldOnExpire(ee, sim)
				cooldownMod.Deactivate()
			}
		},
	})
}
