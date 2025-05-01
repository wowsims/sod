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
// Glaciate now stacks to 10 and Spellfrost Bolt grants 2 stacks per hit.
func (mage *Mage) applyScarletEnclaveDamage2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Mage - Damage 2P Bonus"
	if mage.HasAura(label) {
		return
	}

	aura := core.MakePermanent(mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226423},
		Label:    label,
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
			if dot := mage.LivingBomb.Dot(result.Target); spell.Matches(ClassSpellMask_MageLivingBombExplosion) && result.Landed() &&
				// This is tricky because the Explosion is triggered after the DoT expires, but it shouldn't reapply Living Bomb to the exploding target
				// If we know the DoT isn't active and the #NextTickAt() is CurrentTTime + 1 second then it must be the same target
				!dot.IsActive() && dot.NextTickAt() != sim.CurrentTime+time.Second {
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

				aura.MaxStacks += 5
			}
		}).ApplyOnSpellHitDealt(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_MageSpellfrostBolt) && result.Landed() {
				mage.GlaciateAuras.Get(result.Target).Activate(sim)
				mage.GlaciateAuras.Get(result.Target).AddStack(sim)
			}
		})
	}
}

// Casting Deep Freeze increases the remaining duration of your Icy Veins spell by 10 sec.
// Casting Pyroblast cancels 2 stack of the effect from your Balefire Bolt.
func (mage *Mage) applyScarletEnclaveDamage4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Mage - Damage 4P Bonus"
	if mage.HasAura(label) {
		return
	}

	aura := core.MakePermanent(mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226446},
		Label:    label,
	}))

	if mage.HasRune(proto.MageRune_RuneBracersBalefireBolt) && mage.Talents.Pyroblast {
		aura.ApplyOnCastComplete(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_MagePyroblast) && mage.BalefireAura.IsActive() && mage.BalefireAura.GetStacks() > 0 {
				// These have to be separate
				mage.BalefireAura.RemoveStack(sim)
				if mage.BalefireAura.GetStacks() > 0 {
					mage.BalefireAura.RemoveStack(sim)
				}
			}
		}, false)
	}

	if mage.HasRune(proto.MageRune_RuneHelmDeepFreeze) && mage.HasRune(proto.MageRune_RuneLegsIcyVeins) {
		aura.ApplyOnCastComplete(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_MageDeepFreeze) && mage.IcyVeinsAura.IsActive() {
				mage.IcyVeinsAura.UpdateExpires(sim, mage.IcyVeinsAura.ExpiresAt()+time.Second*10)
			}
		}, false)
	}
}

// Reduces the cooldown on your Frozen Orb spell by 25 sec.
// Each time Glaciate is consumed, the cooldown on your Deep Freeze is reduced by 1.0 sec per stack consumed.
// Reduces the cooldown on Fire Blast by 5 sec and Fire Blast now refreshes the duration of your Living Bomb on the target.
func (mage *Mage) applyScarletEnclaveDamage6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Mage - Damage 6P Bonus"
	if mage.HasAura(label) {
		return
	}

	aura := core.MakePermanent(mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226432},
		Label:    label,
	}))

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
					if newStacks == 0 {
						mage.DeepFreeze.CD.ModifyRemainingCooldown(sim, -time.Second*time.Duration(oldStacks))
					}
				})
			}
		})
	}

	aura.AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_MageFireBlast,
		TimeValue: -time.Second * 5,
	})

	if mage.HasRune(proto.MageRune_RuneHandsLivingBomb) {
		aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if dot := mage.LivingBomb.Dot(result.Target); dot.IsActive() && spell.Matches(ClassSpellMask_MageFireBlast) && result.Landed() {
				dot.Rollover(sim)
			}
		}
	}
}

var ItemSetFireleafVestments = core.NewItemSet(core.ItemSet{
	Name: "Fireleaf Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyScarletEnclaveHealer2PBonus()
		},
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyScarletEnclaveHealer4PBonus()
		},
		6: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyScarletEnclaveHealer6PBonus()
		},
	},
})

// Your Arcane Blast has a 10% chance to cause Arcane Tunneling.
// Arcane Tunneling prevents your Arcane Blast effect from being consumed by the next other Arcane damage spell you cast.
// In addition, activating Arcane Power resets the cooldown on your Mass Regeneration.
func (mage *Mage) applyScarletEnclaveHealer2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Mage - Healer 2P Bonus"
	if mage.HasAura(label) {
		return
	}

	mage.ArcaneTunnelingAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Tunneling",
		ActionID: core.ActionID{SpellID: 1226406},
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolArcane) && spell.ClassSpellMask > 0 && !spell.Matches(ClassSpellMask_MageArcaneBlast) {
				aura.Deactivate(sim)
			}
		},
	})

	aura := core.MakePermanent(mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226407},
		Label:    label,
	}))

	mage.ArcaneTunnelingProcChance += 0.10

	if mage.HasRune(proto.MageRune_RuneHandsArcaneBlast) {
		aura.ApplyOnSpellHitDealt(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_MageArcaneBlast) && result.Landed() && sim.Proc(mage.ArcaneTunnelingProcChance, "Arcane Tunneling") {
				mage.ArcaneTunnelingAura.Activate(sim)
			}
		})
	}

	if mage.Talents.ArcanePower && mage.HasRune(proto.MageRune_RuneLegsMassRegeneration) {
		aura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
			mage.ArcanePower.RelatedSelfBuff.ApplyOnCastComplete(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				mage.MassRegeneration.CD.Reset()
			}, false)
		})
	}
}

// Rewind Time also reduces all damage taken by your target by 20% for 8 sec.
func (mage *Mage) applyScarletEnclaveHealer4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Mage - Healer 4P Bonus"
	if mage.HasAura(label) {
		return
	}

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: label,
	}))
}

// Reduces the cooldown of your Arcane Power by 90 sec and increases its duration by 10 sec.
// While Arcane Power is active, your chance to gain Arcane Tunneling is increased by 10% and each cast of Arcane Blast reduces the remaining cooldown on Mass Regeneration by 1.0 sec.
func (mage *Mage) applyScarletEnclaveHealer6PBonus() {
	if !mage.Talents.ArcanePower {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Mage - Healer 6P Bonus"
	if mage.HasAura(label) {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_MageArcanePower,
		TimeValue: -time.Second * 90,
	})

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			mage.ArcanePower.RelatedSelfBuff.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				mage.ArcaneTunnelingProcChance += 0.10
			}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				mage.ArcaneTunnelingProcChance -= 0.10
			})

			if mage.HasRune(proto.MageRune_RuneLegsMassRegeneration) {
				mage.ArcanePower.RelatedSelfBuff.ApplyOnCastComplete(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.Matches(ClassSpellMask_MageArcaneBlast) {
						mage.MassRegeneration.CD.ModifyRemainingCooldown(sim, -time.Second)
					}
				}, false)
			}
		},
	}))
}
