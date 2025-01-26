package rogue

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetBonescytheArmor = core.NewItemSet(core.ItemSet{
	Name: "Bonescythe Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyNaxxramasDamage2PBonus()
		},
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyNaxxramasDamage4PBonus()
		},
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyNaxxramasDamage6PBonus()
		},
	},
})

// Your Ambush and Instant Poison deal 20% more damage. You heal for 5% of all damage done by your Poisons.
func (rogue *Rogue) applyNaxxramasDamage2PBonus() {
	label := "S03 - Item - Naxxramas - Rogue - Damage 2P Bonus"
	if rogue.HasAura(label) {
		return
	}

	healthMetrics := rogue.NewHealthMetrics(core.ActionID{SpellID: 1219261})

	core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: ClassSpellMask_RogueDeadlyPoisonTick | ClassSpellMask_RogueOccultPoisonTick | ClassSpellMask_RogueInstantPoison,
		Callback:       core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rogue.GainHealth(sim, result.Damage*0.05, healthMetrics)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  ClassSpellMask_RogueAmbush | ClassSpellMask_RogueInstantPoison,
		FloatValue: 0.20,
	})
}

// You have a 100% chance gain 1 Energy each time you deal periodic Nature or Bleed damage.
func (rogue *Rogue) applyNaxxramasDamage4PBonus() {
	label := "S03 - Item - Naxxramas - Rogue - Damage 4P Bonus"
	if rogue.HasAura(label) {
		return
	}

	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 1219288})

	core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
		Name:        label,
		SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolNature,
		Callback:    core.CallbackOnPeriodicDamageDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rogue.AddEnergy(sim, 1, energyMetrics)
		},
	})
}

// You gain 1% increased damage to Undead for 30 sec per Combo Point you spend, stacking up to 25 times.
func (rogue *Rogue) applyNaxxramasDamage6PBonus() {
	label := "S03 - Item - Naxxramas - Rogue - Damage 6P Bonus"
	if rogue.HasAura(label) {
		return
	}

	var undeadTargets []*core.Unit

	buffAura := rogue.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1219291},
		Label:     "Undead Slaying",
		Duration:  time.Second * 30,
		MaxStacks: 25,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			undeadTargets = core.FilterSlice(rogue.Env.Encounter.TargetUnits, func(unit *core.Unit) bool { return unit.MobType == proto.MobType_MobTypeUndead })
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			for _, unit := range undeadTargets {
				rogue.AttackTables[unit.UnitIndex][proto.CastType_CastTypeMainHand].DamageDealtMultiplier /= 1 + 0.01*float64(oldStacks)
				rogue.AttackTables[unit.UnitIndex][proto.CastType_CastTypeMainHand].DamageDealtMultiplier *= 1 + 0.01*float64(newStacks)
			}
		},
	})

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
				buffAura.Activate(sim)
				buffAura.AddStacks(sim, comboPoints)
			})
		},
	})
}

var ItemSetBonescytheLeathers = core.NewItemSet(core.ItemSet{
	Name: "Bonescythe Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyNaxxramasTank2PBonus()
		},
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyNaxxramasTank4PBonus()
		},
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyNaxxramasTank6PBonus()
		},
	},
})

// Your Tease ability never misses, and your chance to be Dodged or Parried is reduced by 2%.
func (rogue *Rogue) applyNaxxramasTank2PBonus() {
	label := "S03 - Item - Naxxramas - Rogue - Tank 2P Bonus"
	if rogue.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{stats.Expertise: 2 * core.ExpertiseRatingPerExpertiseChance}

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachStatsBuff(bonusStats))
}

// Reduces the cooldown on your Evasion ability by 2 min and reduces the cooldown on your Blade Flurry ability by 1 min.
func (rogue *Rogue) applyNaxxramasTank4PBonus() {
	label := "S03 - Item - Naxxramas - Rogue - Tank 4P Bonus"
	if rogue.HasAura(label) {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.Evasion.CD.FlatModifier -= time.Minute * 2

			if rogue.BladeFlurry != nil {
				rogue.BladeFlurry.CD.FlatModifier -= time.Minute
			}
		},
	})
}

// Any damage from an Undead attacker which would otherwise kill you will instead reduce you to 10% of your maximum health (or your current health, whichever is lower).
// In addition, all damage taken from Undead will be reduced by 90% for 3 sec. This effect cannot occur more than once per minute.
func (rogue *Rogue) applyNaxxramasTank6PBonus() {
	if !rogue.IsTanking() {
		return
	}

	label := "S03 - Item - Naxxramas - Rogue - Tank 6P Bonus"
	if rogue.HasAura(label) {
		return
	}

	actionID := core.ActionID{SpellID: 400023}
	healthMetrics := rogue.NewHealthMetrics(actionID)

	cheatDeathAura := rogue.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    fmt.Sprintf("Cheat Death (%s)", label),
		Duration: time.Second * 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageTakenMultiplier *= 0.10
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageTakenMultiplier /= 0.10
		},
	})

	cheatDeathSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.GainHealth(sim, rogue.MaxHealth()*0.10, healthMetrics)
			cheatDeathAura.Activate(sim)
		},
	})

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if rogue.CurrentHealth() == 0 && cheatDeathSpell.IsReady(sim) {
				cheatDeathSpell.Cast(sim, &rogue.Unit)
			}
		},
	}))
}
