package mage

import (
	"math"
	"slices"
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

	var affectedSpells []*core.Spell

	buffAura := mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213317},
		Label:    "Fire Blast",
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(mage.Spellbook, func(spell *core.Spell) bool {
				return spell.Flags.Matches(SpellFlagMage) && spell.SpellSchool.Matches(core.SpellSchoolFire) && !spell.Flags.Matches(core.SpellFlagPassiveSpell)
			})
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.BonusCritRating += 50 * core.SpellCritRatingPerCritChance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.BonusCritRating -= 50 * core.SpellCritRatingPerCritChance
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}

			if !slices.Contains(affectedSpells, spell) {
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
	})

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: label,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode == SpellCode_MageFireBlast {
				buffAura.Activate(sim)
			}
		},
	}))
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

	mage.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			mage.Ignite.PeriodicDamageMultiplierAdditive += 0.10
		},
	})
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

// Your Fireball, Frostfire Bolt, and Balefire Bolt spells gain 3% increased damage for each of your Fire effects on your target, up to a maximum increased of 9%.
func (mage *Mage) applyRAQFire3PBonus() {
	label := "S03 - Item - RAQ - Mage - Fire 3P Bonus"
	if mage.HasAura(label) {
		return
	}

	perEffectModifier := 0.03
	maxModifier := 0.09

	mage.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			dotSpells := core.FilterSlice(mage.Spellbook, func(spell *core.Spell) bool {
				return spell.Flags.Matches(SpellFlagMage) && spell.SpellSchool.Matches(core.SpellSchoolFire) && len(spell.Dots()) > 0
			})

			affectedSpells := core.FilterSlice(
				core.Flatten(
					[][]*core.Spell{
						mage.Fireball,
						{mage.FrostfireBolt},
						{mage.BalefireBolt},
					},
				), func(spell *core.Spell) bool { return spell != nil },
			)

			for _, spell := range affectedSpells {
				oldApplyEffects := spell.ApplyEffects
				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					modifier := 0.0

					for _, spell := range dotSpells {
						if spell.Dot(target).IsActive() {
							modifier += perEffectModifier
						}
					}

					modifier = math.Min(maxModifier, modifier)

					spell.DamageMultiplierAdditive += modifier
					oldApplyEffects(sim, target, spell)
					spell.DamageMultiplierAdditive -= modifier
				}
			}
		},
	})
}
