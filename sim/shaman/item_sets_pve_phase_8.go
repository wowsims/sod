package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var ItemSetTheSoulcrushersStorm = core.NewItemSet(core.ItemSet{
	Name: "The Soulcrusher's Storm",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyScarletEnclaveElemental2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyScarletEnclaveElemental4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyScarletEnclaveElemental6PBonus()
		},
	},
})

// When your Lightning Bolt, Lava Burst, or Chain Lightning strike a target afflicted with your Flame Shock Rank 5 or Rank 6, they also deal one pulse of Flame Shock's damage.
func (shaman *Shaman) applyScarletEnclaveElemental2PBonus() {
	label := "S03 - Item - Scarlet Enclave - Shaman - Elemental 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	// These interactions are unlisted but confirmed by Zirene. All stack multiplicatively.
	// -1 per point of Concussion because the ticks double dip
	// 40% from Tier 3 2-piece
	// 60% from Storm, Earth, and Fire
	damageMultiplier := 1.0
	damageMultiplier *= 1 - 0.01*float64(shaman.Talents.Concussion)
	if shaman.HasSetBonus(ItemSetTheEarthshatterersStorm, 2) {
		damageMultiplier *= 1 + EleTier32pFlameShockDamageBonus/100
	}
	if shaman.HasRune(proto.ShamanRune_RuneCloakStormEarthAndFire) {
		damageMultiplier *= 1 + StormEarthAndFireFlameShockDamageBonus/100
	}

	flameShockCopy := shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1226972}.WithTag(1),
		ClassSpellMask: ClassSpellMask_ShamanFlameShock,
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc,
		Flags:          core.SpellFlagTreatAsPeriodic | core.SpellFlagPureDot | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Flame Shock (2pT4)",
			},

			NumberOfTicks: 1,
			TickLength:    0,
		},

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {},
	})

	flameShockSpells := []*core.Spell{}
	core.MakePermanent(shaman.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226961},
		Label:    label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if shaman.FlameShock[5] != nil {
				flameShockSpells = append(flameShockSpells, shaman.FlameShock[5])
			}
			if shaman.FlameShock[6] != nil {
				flameShockSpells = append(flameShockSpells, shaman.FlameShock[6])
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_ShamanLightningBolt|ClassSpellMask_ShamanChainLightning|ClassSpellMask_ShamanLavaBurst) && result.Landed() {
				for _, spell := range flameShockSpells {
					if dot := spell.Dot(result.Target); dot.IsActive() {
						flameShockCopy.Cast(sim, result.Target)
						flameShockCopy.CalcAndDealDamage(sim, result.Target, dot.SnapshotBaseDamage, flameShockCopy.Dot(result.Target).OutcomeTick)
						break
					}
				}
			}
		},
	}))
}

// Increases the chance to trigger your Overload by an additional 10%.
// Additionally, each time Lightning Bolt or Chain Lightning damages a target, your next Lava Burst deals 10% increased damage, stacking up to 5 times.
// Your Lava Burst deals increased damage equal to its critical strike chance.
func (shaman *Shaman) applyScarletEnclaveElemental4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Shaman - Elemental 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	classMask := ClassSpellMask_ShamanLightningBolt | ClassSpellMask_ShamanChainLightning

	damageMod := shaman.AddDynamicMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanLavaBurst,
		Kind:      core.SpellMod_DamageDone_Flat,
	})

	buffAura := shaman.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1233835},
		Label:     "Thunder and Lava",
		MaxStacks: 5,
		Duration:  time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			damageMod.UpdateIntValue(10 * int64(newStacks))
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_ShamanLavaBurst) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226977},
		Label:    label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.overloadProcChance += 0.10
			shaman.useLavaBurstCritScaling = true
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(classMask) && result.Landed() {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}

// When your Chain Lightning damages fewer than 3 targets, it deals 35% increased damage for each target less than 3.
func (shaman *Shaman) applyScarletEnclaveElemental6PBonus() {
	label := "S03 - Item - Scarlet Enclave - Shaman - Elemental 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	numTargets := shaman.Env.GetNumTargets()

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226978},
		Label:    label,
	})).AttachSpellMod(core.SpellModConfig{
		ClassMask:  ClassSpellMask_ShamanChainLightning,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 1 + (0.35 * float64(3-numTargets)),
	})
}

var ItemSetTheSoulcrushersRage = core.NewItemSet(core.ItemSet{
	Name: "The Soulcrusher's Rage",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyScarletEnclaveEnhancement2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyScarletEnclaveEnhancement4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyScarletEnclaveEnhancement6PBonus()
		},
	},
})

// While Static Shock is active, Lava Lash, Lava Burst, and Stormstrike have a 100% chance to add 1 charge to your Lightning Shield.
// If charges exceed 9, Lightning Shield will immediately deal damage to your target instead of adding charges.
func (shaman *Shaman) applyScarletEnclaveEnhancement2PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneBracersStaticShock) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Shaman - Enhancement 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	classMask := ClassSpellMask_ShamanLavaLash | ClassSpellMask_ShamanStormstrikeHit | ClassSpellMask_ShamanLavaBurst

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226984},
		Label:    label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(classMask) && shaman.ActiveShield != nil && shaman.ActiveShield.Matches(ClassSpellMask_ShamanLightningShield) {
				if shaman.ActiveShieldAura.GetStacks() == 9 {
					shaman.LightningShieldProcs[shaman.ActiveShield.Rank].Cast(sim, result.Target)
				}

				shaman.ActiveShieldAura.AddStack(sim) // Add back the charge removed
			}
		},
	}))
}

// Reduces the cooldown on your Fire Nova Totem by 60%, increases its damage by 200%, and reduces its mana cost by 50%.
// Additionally, your Fire Nova Totem now activates instantly on cast.
func (shaman *Shaman) applyScarletEnclaveEnhancement4PBonus() {
	if shaman.HasRune(proto.ShamanRune_RuneWaistFireNova) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Shaman - Enhancement 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226986},
		Label:    label,
	})).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanFireNovaTotem,
		Kind:      core.SpellMod_Cooldown_Multi_Flat,
		IntValue:  -60,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanFireNovaTotemAttack,
		Kind:      core.SpellMod_DamageDone_Flat,
		IntValue:  200,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanFireNovaTotem,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -50,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanFireNovaTotem,
		Kind:      core.SpellMod_DotTickLength_Flat,
		TimeValue: -time.Second * 5,
	})
}

// Maelstrom Weapon can now stack up to 10 charges. You will also gain 2 charges at a time while wielding a two-handed weapon.
// Any excess charges will increase damage or healing dealt by the affected spell by 20% per excess charge.
// If you have 10 charges when casting an affected spell, all charges will be used and the spell will be instantly cast twice for 200% of normal damage or healing.
func (shaman *Shaman) applyScarletEnclaveEnhancement6PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Shaman - Enhancement 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	twoHandedBonusAura := shaman.RegisterAura(core.Aura{
		Label:    label + " - 2h maelstrom bonus",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.maelstromWeaponProcsPerStack += 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.maelstromWeaponProcsPerStack -= 1
		},
	})

	// This may be reverted
	// if shaman.MainHand().HandType == proto.HandType_HandTypeTwoHand {
	core.MakePermanent(twoHandedBonusAura)
	// }
	// shaman.RegisterItemSwapCallback(core.AllWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
	// 	if shaman.MainHand().HandType == proto.HandType_HandTypeTwoHand {
	// 		twoHandedBonusAura.Activate(sim)
	// 	} else {
	// 		twoHandedBonusAura.Deactivate(sim)
	// 	}
	// })

	// @Lucenia: We have to use a boolean flag because otherwise the triggered cast infinitely procs this trigger
	var damageMod *core.SpellMod
	var isProcced bool

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226997},
		Label:    label,
		OnInit: func(_ *core.Aura, sim *core.Simulation) {
			shaman.MaelstromWeaponAura.MaxStacks += 5

			// We have to initialize this within the OnInit for shaman.MaelstromWeaponClassMask to be set properly
			damageMod = shaman.AddDynamicMod(core.SpellModConfig{
				ClassMask:  shaman.MaelstromWeaponClassMask,
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 1,
			})
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			isProcced = false
		},
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			damageMod.Deactivate()
		},
		OnApplyEffects: func(_ *core.Aura, _ *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if spell.Matches(shaman.MaelstromWeaponClassMask) {
				damageMod.UpdateFloatValue(1 + max(0, 0.20*float64(shaman.MaelstromWeaponAura.GetStacks()-5)))
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(shaman.MaelstromWeaponClassMask) && !spell.ProcMask.Matches(core.ProcMaskSpellProc) && shaman.MaelstromWeaponAura.GetStacks() == 10 && !isProcced {
				isProcced = true

				if spell.CD.Duration > 0 {
					spell.CD.Reset()
				}

				defaultGCD := spell.DefaultCast.GCD
				spell.DefaultCast.GCD = 0
				spell.Cast(sim, shaman.CurrentTarget)
				spell.DefaultCast.GCD = defaultGCD

				isProcced = false
			}
		},
	}))
}

var ItemSetTheSoulcrushersResolve = core.NewItemSet(core.ItemSet{
	Name: "The Soulcrusher's Resolve",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyScarletEnclaveTank2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyScarletEnclaveTank4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyScarletEnclaveTank6PBonus()
		},
	},
})

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
		ActionID: core.ActionID{SpellID: 1227153},
		Label:    label,
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

	healthMetrics := shaman.NewHealthMetrics(core.ActionID{SpellID: 1227160})

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		ActionID:       core.ActionID{SpellID: 1227159},
		Name:           label,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: ClassSpellMask_ShamanLightningShieldProc,
		Outcome:        core.OutcomeLanded,
		ICD:            time.Second * 3,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shaman.GainHealth(sim, result.Damage, healthMetrics)
		},
	})
}

// Your Shield Mastery stacks also reduce the cast time of your Lava Burst by 20% per stack.
// Lava Burst no longer consumes Maelstrom Weapon charges.
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
		ActionID: core.ActionID{SpellID: 1227164},
		Label:    label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ShieldMasteryAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				spellMod.Activate()
			}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				spellMod.Deactivate()
			}).ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				spellMod.UpdateFloatValue(-0.20 * float64(newStacks))
			})

			if shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
				shaman.MaelstromWeaponClassMask ^= ClassSpellMask_ShamanLavaBurst
				for _, spellMod := range shaman.MaelstromWeaponSpellMods {
					spellMod.RemoveSpellByClassMask(ClassSpellMask_ShamanLavaBurst)
				}
			}
		},
	}))
}

var ItemSetTheSoulcrusher = core.NewItemSet(core.ItemSet{
	Name: "The Soulcrusher",
	Bonuses: map[int32]core.ApplyEffect{
		// Heals from your Earth Shield have a 40% chance to make your next cast time heal instant cast.
		2: func(agent core.Agent) {
		},
		// Your Healing Wave Rank 9 and Rank 10 and Lesser Healing Wave Rank 6 spells have a 10% chance to imbue your target with Totemic Power.
		4: func(agent core.Agent) {
		},
		// The target of your Spirit of the Alpha ability takes 20% reduced damage from Undead enemies.
		6: func(agent core.Agent) {
		},
	},
})
