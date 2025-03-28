package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var ItemSetDeathdealersThrill = core.NewItemSet(core.ItemSet{
	Name: "Deathdealer's Thrill",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyTAQDamage2PBonus()
		},
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyTAQDamage4PBonus()
		},
	},
})

// Increases Mutilate and Sinister Strike damage by 20%
func (rogue *Rogue) applyTAQDamage2PBonus() {
	if !rogue.HasRune(proto.RogueRune_RuneMutilate) && !rogue.HasRune(proto.RogueRune_RuneSaberSlash) {
		return
	}

	label := "S03 - Item - TAQ - Rogue - Damage 2P Bonus"
	if rogue.HasAura(label) {
		return
	}

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_RogueSaberSlash | ClassSpellMask_RogueMutilateHit,
		IntValue:  20,
	}))
}

// Reduces the cooldown on Adrenaline Rush by 4 minutes.
func (rogue *Rogue) applyTAQDamage4PBonus() {
	if !rogue.Talents.AdrenalineRush {
		return
	}

	label := "S03 - Item - TAQ - Rogue - Damage 4P Bonus"
	if rogue.HasAura(label) {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.AdrenalineRush.CD.Duration -= time.Minute * 4
		},
	})
}

var ItemSetDeathdealersBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Deathdealer's Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyTAQTank2PBonus()
		},
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyTAQTank4PBonus()
		},
	},
})

// Your Main Gauche now strikes 1 additional nearby target and also causes your Sinister Strike to strike 1 additional nearby target.
// These additional strikes are not duplicated by Blade Flurry.
func (rogue *Rogue) applyTAQTank2PBonus() {
	if !rogue.HasRune(proto.RogueRune_RuneMainGauche) || rogue.Env.GetNumTargets() == 1 {
		return
	}

	label := "S03 - Item - TAQ - Rogue - Tank 2P Bonus"
	if rogue.HasAura(label) {
		return
	}

	var curDmg float64

	cleaveHit := rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1213754},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, curDmg, spell.OutcomeAlwaysHit)
		},
	})

	cleaveAura := rogue.RegisterAura(core.Aura{
		Label:    "2P Cleave Buff",
		Duration: time.Second * 10,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.Matches(ClassSpellMask_RogueSinisterStrike) {
				curDmg = result.Damage / result.ResistanceMultiplier
				cleaveHit.Cast(sim, rogue.Env.NextTargetUnit(result.Target))
				cleaveHit.SpellMetrics[result.Target.UnitIndex].Casts--
			}
		},
	})

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.Matches(ClassSpellMask_RogueMainGauche) {
				cleaveAura.Activate(sim)
				curDmg = result.Damage / result.ResistanceMultiplier
				cleaveHit.Cast(sim, rogue.Env.NextTargetUnit(result.Target))
				cleaveHit.SpellMetrics[result.Target.UnitIndex].Casts--
			}
		},
	}))
}

// While active, your Main Gauche also causes you to heal for 10% of all damage done by Sinister Strike.
// Any excess healing becomes a Blood Barrier, absorbing damage up to 20% of your maximum health.
func (rogue *Rogue) applyTAQTank4PBonus() {
	if !rogue.HasRune(proto.RogueRune_RuneMainGauche) {
		return
	}

	label := "S03 - Item - TAQ - Rogue - Tank 4P Bonus"
	if rogue.HasAura(label) {
		return
	}

	healthMetrics := rogue.NewHealthMetrics(core.ActionID{SpellID: 11294})
	healAmount := 0.0
	shieldAmount := 0.0
	currentShield := 0.0

	var shieldSpell *core.Spell

	shieldSpell = rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1213761},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Shield: core.ShieldConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label:    "Blood Barrier",
				ActionID: core.ActionID{SpellID: 1213762},
				Duration: time.Second * 15,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					shieldAmount = 0.0
					currentShield = 0.0
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if currentShield <= 0 || result.Damage <= 0 {
						return
					}

					damageReduced := min(result.Damage, currentShield)
					currentShield -= damageReduced

					rogue.GainHealth(sim, damageReduced, shieldSpell.HealthMetrics(result.Target))
					if currentShield <= 0 {
						shieldSpell.SelfShield().Deactivate(sim)
					}
				},
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if currentShield < rogue.MaxHealth()*0.2 {
				shieldAmount = min(shieldAmount, rogue.MaxHealth()*0.2-currentShield)
				currentShield += shieldAmount
				spell.SelfShield().Apply(sim, shieldAmount)
			}
		},
	})

	activeAura := core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
		Name:     "Main Gauche - Blood Barrier",
		ActionID: core.ActionID{SpellID: 1213762},
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		Duration: time.Second * 15,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_RogueSinisterStrike) {
				healAmount = result.Damage * 0.15
				if rogue.CurrentHealth() < rogue.MaxHealth() {
					rogue.GainHealth(sim, healAmount, healthMetrics)
				} else {
					shieldAmount = healAmount
					shieldSpell.Cast(sim, result.Target)
				}

			}
		},
	})

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.Matches(ClassSpellMask_RogueMainGauche) {
				activeAura.Activate(sim)
			}
		},
	}))
}

var ItemSetEmblemsofVeiledShadows = core.NewItemSet(core.ItemSet{
	Name: "Emblems of Veiled Shadows",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyRAQDamage3PBonus()
		},
	},
})

// 3 pieces: Your finishing moves cost 50% less Energy.
func (rogue *Rogue) applyRAQDamage3PBonus() {
	label := "S03 - Item - RAQ - Rogue - Damage 3P Bonus"
	if rogue.HasAura(label) {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, finisher := range rogue.Finishers {
				finisher.Cost.Multiplier -= 50
			}
		},
	})
}
