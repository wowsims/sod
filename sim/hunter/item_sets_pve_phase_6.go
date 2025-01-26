package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var StrikersProwess = core.NewItemSet(core.ItemSet{
	Name: "Striker's Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyTAQMelee2PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyTAQMelee4PBonus()
		},
	},
})

// Increases Wyvern Strike DoT by 50% and increases your pet's maximum focus by 50.
func (hunter *Hunter) applyTAQMelee2PBonus() {
	label := "S03 - Item - TAQ - Hunter - Melee 2P Bonus"
	if hunter.HasAura(label) {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.WyvernStrike != nil {
				hunter.WyvernStrike.PeriodicDamageMultiplierAdditive += 0.50
			}

			if hunter.pet != nil {
				hunter.pet.IncreaseMaxFocus(50)
			}
		},
	})
}

// Increases the Impact Damage of Mongoose Bite and all Strikes by 15%
func (hunter *Hunter) applyTAQMelee4PBonus() {
	label := "S03 - Item - TAQ - Hunter - Melee 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range hunter.Strikes {
				spell.ImpactDamageMultiplierAdditive += 0.15
			}
			hunter.RaptorStrikeMH.ImpactDamageMultiplierAdditive += 0.15
			hunter.RaptorStrikeOH.ImpactDamageMultiplierAdditive += 0.15
			hunter.MongooseBite.ImpactDamageMultiplierAdditive += 0.15
		},
	})
}

var StrikersPursuit = core.NewItemSet(core.ItemSet{
	Name: "Striker's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyTAQRanged2PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyTAQRanged4PBonus()
		},
	},
})

const TAQRanged2PBonusLabel = "S03 - Item - TAQ - Hunter - Ranged 2P Bonus"

// Increases Kill Shot damage by 50% against non-player targets.
func (hunter *Hunter) applyTAQRanged2PBonus() {
	if !hunter.HasRune(proto.HunterRune_RuneLegsKillShot) {
		return
	}

	if hunter.HasAura(TAQRanged2PBonusLabel) {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label: TAQRanged2PBonusLabel,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			hunter.KillShot.DamageMultiplierAdditive += 0.20
		},
	})
}

// Kill Shot's cooldown is reduced by 50%.
// While Rapid Fire is active with Rapid killing engraved, Kill Shot has no cooldown and fires 3 additional Kill Shots at 30% damage, with a minimum range.
func (hunter *Hunter) applyTAQRanged4PBonus() {
	if !hunter.HasRune(proto.HunterRune_RuneLegsKillShot) {
		return
	}

	label := "S03 - Item - TAQ - Hunter - Ranged 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	clonedShotConfig := hunter.newKillShotConfig()
	clonedShotConfig.ActionID.Tag = 1
	clonedShotConfig.Flags |= core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell
	clonedShotConfig.Flags ^= core.SpellFlagAPL
	clonedShotConfig.Cast.DefaultCast.GCD = 0
	clonedShotConfig.Cast.DefaultCast.Cost = 0
	clonedShotConfig.Cast.CD = core.Cooldown{}
	clonedShotConfig.ManaCost.BaseCost = 0
	clonedShotConfig.ManaCost.FlatCost = 0
	clonedShotConfig.MetricSplits = 0
	clonedShotConfig.DamageMultiplier *= 0.30
	clonedShotConfig.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return hunter.DistanceFromTarget >= core.MinRangedAttackDistance
	}

	clonedShot := hunter.RegisterSpell(clonedShotConfig)

	hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			hunter.KillShot.CD.Multiplier *= 0.5

			if !hunter.HasRune(proto.HunterRune_RuneHelmRapidKilling) {
				return
			}

			if hunter.HasAura(TAQRanged2PBonusLabel) {
				clonedShot.DamageMultiplierAdditive += 0.20 // Add the 2p bonus 20%
			}

			oldApplyEffects := hunter.KillShot.ApplyEffects
			hunter.KillShot.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				oldApplyEffects(sim, target, spell)

				if hunter.RapidFireAura.IsActive() {
					spell.CD.Reset()

					for i := 1; i < 4; i++ {
						core.StartDelayedAction(sim, core.DelayedActionOptions{
							DoAt: sim.CurrentTime + time.Duration(i*375)*time.Millisecond,
							OnAction: func(sim *core.Simulation) {
								// Ensure that the cloned shots get any damage amps from the main Kill Shot ability
								clonedShot.DamageMultiplier *= spell.DamageMultiplier
								clonedShot.DamageMultiplierAdditive += spell.DamageMultiplierAdditive - 1
								clonedShot.Cast(sim, target)
								clonedShot.DamageMultiplier /= spell.DamageMultiplier
								clonedShot.DamageMultiplierAdditive -= spell.DamageMultiplierAdditive - 1
							},
						})
					}
				}
			}
		},
	})
}

var TrappingsOfTheUnseenPath = core.NewItemSet(core.ItemSet{
	Name: "Trappings of the Unseen Path",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyRAQBeastmastery5PBonus()
		},
	},
})

// Increases the Focus regeneration of your Beast pet by 100%.
func (hunter *Hunter) applyRAQBeastmastery5PBonus() {
	if hunter.pet == nil {
		return
	}

	label := "S03 - Item - RAQ - Hunter - Beastmastery 5P Bonus"
	if hunter.HasAura(label) {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			hunter.pet.AddFocusRegenMultiplier(1.00)
		},
	})
}
