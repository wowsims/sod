package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var ItemSetEnigmaInsight = core.NewItemSet(core.ItemSet{
	Name: "Enigma Insight",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyTAQFire2PBonus()
		},
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyTAQFire4PBonus()
		},
	},
})

// Your Fire Blast now also causes your next Fire spell to gain 50% increased critical strike chance for 10 sec.
func (mage *Mage) applyTAQFire2PBonus() {
	label := "S03 - Item - TAQ - Mage - Fire 2P Bonus"
	if mage.HasAura(label) {
		return
	}

	buffAura := mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213317},
		Label:    "Fire Blast",
		Duration: time.Second * 10,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}

			if !(spell.Matches(ClassSpellMask_MageAll) && spell.SpellSchool.Matches(core.SpellSchoolFire) && !spell.Flags.Matches(core.SpellFlagPassiveSpell)) {
				return
			}

			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + core.SpellBatchWindow,
				OnAction: func(sim *core.Simulation) {
					if aura.IsActive() {
						aura.Deactivate(sim)
					}
				},
			})
		},
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:         ClassSpellMask_MageAll,
		School:            core.SpellSchoolFire,
		SpellFlagsExclude: core.SpellFlagPassiveSpell,
		Kind:              core.SpellMod_BonusCrit_Flat,
		FloatValue:        50 * core.SpellCritRatingPerCritChance,
	})

	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: ClassSpellMask_MageFireBlast,
		Callback:       core.CallbackOnCastComplete,
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			buffAura.Activate(sim)
		},
	})
}

// Increases the damage done by your Ignite talent by 10%.
func (mage *Mage) applyTAQFire4PBonus() {
	if mage.Talents.Ignite == 0 {
		return
	}

	label := "S03 - Item - TAQ - Mage - Fire 4P Bonus"
	if mage.HasAura(label) {
		return
	}

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_MageIgnite,
		Kind:      core.SpellMod_DamageDone_Flat,
		IntValue:  10,
	}))
}

var ItemSetEnigmaMoment = core.NewItemSet(core.ItemSet{
	Name: "Enigma Moment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyTAQArcane2PBonus()
		},
		// Your Mana Shield, Fire Ward, and Frost Ward absorb 50% more damage and also place a Temporal Beacon on the target for 30 sec.
		4: func(agent core.Agent) {
			// Nothing to do
		},
	},
})

// Your Arcane Blast increases damage and healing done by an additional 10% per stack.
func (mage *Mage) applyTAQArcane2PBonus() {
	if !mage.HasRune(proto.MageRune_RuneHandsArcaneBlast) {
		return
	}

	label := "S03 - Item - TAQ - Mage - Arcane 2P Bonus"
	if mage.HasAura(label) {
		return
	}

	mage.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			mage.ArcaneBlastDamageMultiplier += 0.10
		},
	})
}

var ItemSetTrappingsOfVaultedSecrets = core.NewItemSet(core.ItemSet{
	Name: "Trappings of Vaulted Secrets",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyRAQFire3PBonus()
		},
	},
})

// Your Fireball, Frostfire Bolt, and Balefire Bolt spells gain 4% increased damage for each of your Fire effects on your target, up to a maximum increased of 12%.
func (mage *Mage) applyRAQFire3PBonus() {
	label := "S03 - Item - RAQ - Mage - Fire 3P Bonus"
	if mage.HasAura(label) {
		return
	}

	perEffectMultiplier := 0.04
	maxMultiplier := 1.12

	classSpellMasks := ClassSpellMask_MageFireball | ClassSpellMask_MageFrostfireBolt | ClassSpellMask_MageBalefireBolt
	damageMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  classSpellMasks,
		FloatValue: 1,
	})

	var dotSpells []*core.Spell
	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			dotSpells = core.FilterSlice(mage.Spellbook, func(spell *core.Spell) bool {
				return spell.Matches(ClassSpellMask_MageAll) && spell.SpellSchool.Matches(core.SpellSchoolFire) && len(spell.Dots()) > 0
			})
		},
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
			multiplier := 1.0

			for _, spell := range dotSpells {
				if spell.Dot(target).IsActive() {
					multiplier += perEffectMultiplier
				}
			}

			multiplier = min(maxMultiplier, multiplier)
			damageMod.UpdateFloatValue(multiplier)
		},
	}))
}
