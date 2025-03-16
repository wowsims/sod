package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// var ItemSetWIPScarletCrusadeEle = core.NewItemSet(core.ItemSet{
// 	Name: "WIPScarletCrusadeEle",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			shaman := agent.(ShamanAgent).GetShaman()
// 			shaman.applyScarletEnclaveElemental2PBonus()
// 		},
// 		4: func(agent core.Agent) {
// 			shaman := agent.(ShamanAgent).GetShaman()
// 			shaman.applyScarletEnclaveElemental4PBonus()
// 		},
// 		6: func(agent core.Agent) {
// 			shaman := agent.(ShamanAgent).GetShaman()
// 			shaman.applyScarletEnclaveElemental6PBonus()
// 		},
// 	},
// })

// TODO: When your Lava Burst strikes a target afflicted with your Flame Shock Rank 5 or Rank 6, it also deals one pulse of Flame Shock's damage.
func (shaman *Shaman) applyScarletEnclaveElemental2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Shaman - Elemental 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	dotSpells := make([]*core.Spell, FlameShockRanks+1)

	for rank := 1; rank <= FlameShockRanks; rank++ {
		dotSpells[rank] = shaman.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: FlameShockSpellId[rank]}.WithTag(1),
			ClassSpellMask: ClassSpellMask_ShamanFlameShock,
			SpellSchool:    core.SpellSchoolFire,
			DefenseType:    core.DefenseTypeMagic,
			ProcMask:       core.ProcMaskSpellDamage,
			Flags:          core.SpellFlagTreatAsPeriodic | core.SpellFlagPureDot | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {},
		})
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(ClassSpellMask_ShamanLavaBurst) {
				return
			}

			for rank := 1; rank <= FlameShockRanks; rank++ {
				if dot := shaman.FlameShock[rank].Dot(result.Target); dot.IsActive() {
					copiedDoTSpell := dotSpells[rank]
					copiedDoTSpell.Cast(sim, result.Target)
					copiedDoTSpell.CalcAndDealDamage(sim, result.Target, dot.SnapshotBaseDamage, copiedDoTSpell.OutcomeAlwaysHit)
					break
				}
			}
		},
	}))
}

// Increases the chance to trigger your Rolling Thunder by an additional 10%.
func (shaman *Shaman) applyScarletEnclaveElemental4PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneBracersRollingThunder) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Shaman - Elemental 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.rollingThunderProcChance += 0.50
		},
	})
}

// Increases the chance to trigger your Overload by an additional 15%.
func (shaman *Shaman) applyScarletEnclaveElemental6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Shaman - Elemental 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.overloadProcChance += 0.40
		},
	})
}

// var ItemSetWIPScarletCrusadeEnh = core.NewItemSet(core.ItemSet{
// 	Name: "WIPScarletCrusadeEnh",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			shaman := agent.(ShamanAgent).GetShaman()
// 			shaman.applyScarletEnclaveEnhancement2PBonus()
// 		},
// 		4: func(agent core.Agent) {
// 			shaman := agent.(ShamanAgent).GetShaman()
// 			shaman.applyScarletEnclaveEnhancement4PBonus()
// 		},
// 		6: func(agent core.Agent) {
// 			shaman := agent.(ShamanAgent).GetShaman()
// 			shaman.applyScarletEnclaveEnhancement6PBonus()
// 		},
// 	},
// })

// Lava Lash and Stormstrike now have a 100% chance to add a charge to your Lightning Shield. If it exceeds 9 charges, Lightning Shield will immediately deal damage to your target instead of adding a charge.
func (shaman *Shaman) applyScarletEnclaveEnhancement2PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsLavaLash) && !shaman.Talents.Stormstrike {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Shaman - Enhancement 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	classMask := ClassSpellMask_ShamanLavaLash | ClassSpellMask_ShamanStormstrikeHit

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(classMask) && shaman.ActiveShield != nil && shaman.ActiveShield.Matches(ClassSpellMask_ShamanLightningShield) {
				if shaman.ActiveShieldAura.GetStacks() == 9 {
					shaman.LightningShieldProcs[shaman.ActiveShield.Rank].Cast(sim, result.Target)
				}

				shaman.ActiveShieldAura.AddStack(sim)
			}
		},
	}))
}

// Reduces the cooldown on your Fire Nova Totem by 60% and your Fire Nova Totem now activates instantly on cast.
func (shaman *Shaman) applyScarletEnclaveEnhancement4PBonus() {
	if shaman.HasRune(proto.ShamanRune_RuneWaistFireNova) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Shaman - Enhancement 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
	})).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanFireNovaTotem,
		Kind:      core.SpellMod_Cooldown_Multi_Flat,
		IntValue:  -60,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanFireNovaTotem,
		Kind:      core.SpellMod_DotTickLength_Flat,
		TimeValue: -time.Second * 5,
	})
}

// Maelstrom Weapon can now stack up to 10 charges. If you have 10 charges when casting an affected spell, all charges will be used and the spell will be instantly cast twice.
func (shaman *Shaman) applyScarletEnclaveEnhancement6PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Shaman - Enhancement 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(_ *core.Aura, _ *core.Simulation) {
			shaman.MaelstromWeaponAura.MaxStacks += 5
			shaman.MaelstromWeaponAura.ApplyOnCastComplete(func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if aura.GetStacks() == 10 && spell.Matches(shaman.MaelstromWeaponClassMask) {
					if spell.CD.Duration > 0 {
						spell.CD.Reset()
					}

					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt:     sim.CurrentTime + core.SpellBatchWindow,
						Priority: core.CooldownPriorityBloodlust,
						OnAction: func(sim *core.Simulation) {
							defaultGCD := spell.DefaultCast.GCD
							spell.DefaultCast.GCD = 0
							spell.Cast(sim, shaman.CurrentTarget)
							spell.DefaultCast.GCD = defaultGCD
						},
					})
				}
			}, true)
		},
	})
}

// var ItemSetWIPScarletCrusadeTank = core.NewItemSet(core.ItemSet{
// 	Name: "WIPScarletCrusadeTank",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			shaman := agent.(ShamanAgent).GetShaman()
// 			shaman.applyScarletEnclaveTank2PBonus()
// 		},
// 		4: func(agent core.Agent) {
// 			shaman := agent.(ShamanAgent).GetShaman()
// 			shaman.applyScarletEnclaveTank4PBonus()
// 		},
// 		6: func(agent core.Agent) {
// 			shaman := agent.(ShamanAgent).GetShaman()
// 			shaman.applyScarletEnclaveTank6PBonus()
// 		},
// 	},
// })

// Your Shield Mastery effect can now stack up to 7 times.
func (shaman *Shaman) applyScarletEnclaveTank2PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestShieldMastery) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Shaman - Tank 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ShieldMasteryAura.MaxStacks += 2
		},
	})
}

// Each time your Lightning Shield deals damage, you heal for 100% of the damage it dealt, no more than once every 3 sec.
func (shaman *Shaman) applyScarletEnclaveTank4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Shaman - Tank 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	healthMetrics := shaman.NewHealthMetrics(core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsEarthShield)}) // TODO: Spell ID
	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_ShamanLightningShieldProc) {
				shaman.GainHealth(sim, result.Damage, healthMetrics)
			}
		},
	}))
}

// Your Shield Mastery stacks also reduce the cast time of your Lava Burst by 20% per stack. Lava Burst no longer consumes Maelstrom Weapon charges.
func (shaman *Shaman) applyScarletEnclaveTank6PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestShieldMastery) || !shaman.HasRune(proto.ShamanRune_RuneHandsLavaBurst) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Shaman - Tank 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	spellMod := shaman.AddDynamicMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanLavaBurst,
		Kind:      core.SpellMod_CastTime_Pct,
	})

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ShieldMasteryAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				spellMod.Activate()
			}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				spellMod.Deactivate()
			}).ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				spellMod.UpdateFloatValue(-0.20 * float64(newStacks-oldStacks))
			})

			if shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
				shaman.MaelstromWeaponClassMask ^= ClassSpellMask_ShamanLavaBurst
				for _, mod := range shaman.MaelstromWeaponSpellMods {
					mod.ClassMask ^= ClassSpellMask_ShamanLavaBurst
				}
			}
		},
	}))
}

// var ItemSetWIPScarletCrusadeResto = core.NewItemSet(core.ItemSet{
// 	Name: "WIPScarletCrusadeResto",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		// Heals from your Earth Shield have a 40% chance to make your next cast time heal instant cast.
// 		2: func(agent core.Agent) {
// 		},
// 		// Your Healing Wave Rank 9 and Rank 10 and Lesser Healing Wave Rank 6 spells have a 10% chance to imbue your target with Totemic Power.
// 		4: func(agent core.Agent) {
// 		},
// 		// The target of your Spirit of the Alpha ability takes 20% reduced damage from Undead enemies.
// 		6: func(agent core.Agent) {
// 		},
// 	},
// })
