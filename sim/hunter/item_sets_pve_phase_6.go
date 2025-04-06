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

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.pet != nil {
				hunter.pet.IncreaseMaxFocus(50)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.pet != nil {
				hunter.pet.DecreaseMaxFocus(50)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PeriodicDamageDone_Flat,
		ClassMask: ClassSpellMask_HunterWyvernStrike,
		IntValue:  50,
	}))
}

// Increases the Impact Damage of Mongoose Bite and all Strikes by 20%
func (hunter *Hunter) applyTAQMelee4PBonus() {
	label := "S03 - Item - TAQ - Hunter - Melee 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_ImpactDamageDone_Flat,
		ClassMask: ClassSpellMask_HunterMongooseBite | ClassSpellMask_HunterStrikes,
		IntValue:  20,
	}))

	// This also applies to the pet's Flanking Strike
	if hunter.pet != nil {
		core.MakePermanent(hunter.pet.RegisterAura(core.Aura{
			Label: label,
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_ImpactDamageDone_Flat,
			ClassMask: ClassSpellMask_HunterPetFlankingStrike,
			IntValue:  20,
		}))
	}
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

// Increases Kill Shot damage by 30% against non-player targets.
func (hunter *Hunter) applyTAQRanged2PBonus() {
	if !hunter.HasRune(proto.HunterRune_RuneLegsKillShot) {
		return
	}

	if hunter.HasAura(TAQRanged2PBonusLabel) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: TAQRanged2PBonusLabel,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterKillShot,
		IntValue:  10,
	}))
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
	clonedShotConfig.ProcMask = core.ProcMaskRangedProc | core.ProcMaskRangedDamageProc
	clonedShotConfig.Flags |= core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell
	clonedShotConfig.Flags ^= core.SpellFlagAPL
	clonedShotConfig.MinRange = core.MinRangedAttackRange
	clonedShotConfig.MaxRange = core.MaxRangedAttackRange
	clonedShotConfig.Cast.DefaultCast.GCD = 0
	clonedShotConfig.Cast.DefaultCast.Cost = 0
	clonedShotConfig.Cast.CD = core.Cooldown{}
	clonedShotConfig.ManaCost.BaseCost = 0
	clonedShotConfig.ManaCost.FlatCost = 0
	clonedShotConfig.MetricSplits = 0
	clonedShotConfig.DamageMultiplier *= 0.30

	clonedShot := hunter.RegisterSpell(clonedShotConfig)

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if !hunter.HasRune(proto.HunterRune_RuneHelmRapidKilling) {
				return
			}

			oldApplyEffects := hunter.KillShot.ApplyEffects
			hunter.KillShot.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				oldApplyEffects(sim, target, spell)

				if hunter.RapidFireAura.IsActive() {
					spell.CD.Reset()

					for i := 1; i < 4; i++ {
						core.StartDelayedAction(sim, core.DelayedActionOptions{
							DoAt: sim.CurrentTime + time.Duration(i*200)*time.Millisecond,
							OnAction: func(sim *core.Simulation) {
								// Ensure that the cloned shots get any damage amps from the main Kill Shot ability
								clonedShot.Cast(sim, target)
							},
						})
					}
				}
			}
		},
	})).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_HunterKillShot,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * 6,
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
