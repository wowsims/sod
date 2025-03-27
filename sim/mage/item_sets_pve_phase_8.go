package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var ItemSetFireleafRegalia = core.NewItemSet(core.ItemSet{
	Name: "Fireleaf Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyScarletEnclaveDamage2PBonus()
		},
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyScarletEnclaveDamage4PBonus()
		},
		6: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyScarletEnclaveDamage6PBonus()
		},
	},
})

// Living Bomb ticks every 1 second and when it explodes it spreads Living Bomb to all targets struck that don't have an active Living Bomb.
// Glaciate now stacks to 8 and Spellfrost Bolt grants 2 stacks per hit.
func (mage *Mage) applyScarletEnclaveDamage2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Mage - Damage 2P Bonus"
	if mage.HasAura(label) {
		return
	}

	aura := core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: label,
	}))

	if mage.HasRune(proto.MageRune_RuneHandsLivingBomb) {
		ticksDelta := LivingBombBaseTickLength - time.Second

		aura.AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_DotTickLength_Flat,
			ClassMask: ClassSpellMask_MageLivingBomb,
			TimeValue: -ticksDelta,
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
			ClassMask: ClassSpellMask_MageLivingBomb,
			IntValue:  LivingBombBaseNumTicks * int64(ticksDelta.Seconds()),
		}).ApplyOnSpellHitDealt(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if dot := mage.LivingBomb.Dot(result.Target); !dot.IsActive() && spell.Matches(ClassSpellMask_MageLivingBombExplosion) && result.Landed() {
				dot.Apply(sim)
			}
		})
	}

	if mage.HasRune(proto.MageRune_RuneHandsIceLance) {
		aura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
			for _, aura := range mage.GlaciateAuras {
				if aura == nil {
					continue
				}

				aura.MaxStacks += 3
			}
		}).ApplyOnSpellHitDealt(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_MageSpellfrostBolt) && result.Landed() {
				mage.GlaciateAuras.Get(result.Target).Activate(sim)
				mage.GlaciateAuras.Get(result.Target).AddStack(sim)
			}
		})
	}
}

// Casts of Pyroblast remove a stack of Balefire Bolt.
// Casts of Deep Freeze increase the duration of your active Icy Veins by 8 seconds.
func (mage *Mage) applyScarletEnclaveDamage4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Mage - Damage 4P Bonus"
	if mage.HasAura(label) {
		return
	}

	aura := core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: label,
	}))

	if mage.HasRune(proto.MageRune_RuneBracersBalefireBolt) && mage.Talents.Pyroblast {
		aura.ApplyOnCastComplete(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_MagePyroblast) && mage.BalefireAura.IsActive() && mage.BalefireAura.GetStacks() > 0 {
				mage.BalefireAura.RemoveStack(sim)
			}
		}, false)
	}

	if mage.HasRune(proto.MageRune_RuneHelmDeepFreeze) && mage.HasRune(proto.MageRune_RuneLegsIcyVeins) {
		aura.ApplyOnCastComplete(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_MageDeepFreeze) && mage.IcyVeinsAura.IsActive() {
				mage.IcyVeinsAura.UpdateExpires(sim, mage.IcyVeinsAura.ExpiresAt()+time.Second*8)
			}
		}, false)
	}
}

// Fire Blast cooldown is reduced by 5 seconds and refreshes Living Bomb.
// Frozen Orb cooldown is reduced by 25 seconds.
// When you use consume stacks of Glaciate the cooldown of Deep Freeze is reduced by 1 second per stack consumed.
func (mage *Mage) applyScarletEnclaveDamage6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Mage - Damage 6P Bonus"
	if mage.HasAura(label) {
		return
	}

	aura := core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: label,
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_MageFireBlast,
		TimeValue: -time.Second * 5,
	})

	if mage.HasRune(proto.MageRune_RuneHandsLivingBomb) {
		aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if dot := mage.LivingBomb.Dot(result.Target); dot.IsActive() && spell.Matches(ClassSpellMask_MageFireBlast) && result.Landed() {
				dot.Refresh(sim)
			}
		}
	}

	if mage.HasRune(proto.MageRune_RuneCloakFrozenOrb) {
		aura.AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_Cooldown_Flat,
			ClassMask: ClassSpellMask_MageFrozenOrb,
			TimeValue: -time.Second * 25,
		})
	}

	if mage.HasRune(proto.MageRune_RuneHelmDeepFreeze) && mage.HasRune(proto.MageRune_RuneHandsIceLance) {
		aura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
			for _, aura := range mage.GlaciateAuras {
				if aura == nil {
					continue
				}

				aura.ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					if newStacks == 0 && !mage.DeepFreeze.CD.Cooldown.IsReady(sim) {
						mage.DeepFreeze.ModifyRemainingCooldown(sim, -time.Second*time.Duration(oldStacks))
					}
				})
			}
		})
	}
}

var ItemSetFireleafVestments = core.NewItemSet(core.ItemSet{
	Name: "Fireleaf Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		// Rewind Time reduces the damage the target takes by 20% for x seconds.
		2: func(agent core.Agent) {
		},
		// TODO: Casts of Arcane Blast have a 10% chance to cause your next Arcane spell cast to not consume your stacks of Arcane Blast
		4: func(agent core.Agent) {
		},
		// Arcane Power's duration is increased by 10 seconds, and it's cooldown is reduced by 1 minute
		6: func(agent core.Agent) {
		},
	},
})
