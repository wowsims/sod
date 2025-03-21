package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (shaman *Shaman) ApplyRunes() {
	// Helm
	shaman.applyBurn()
	shaman.applyMentalDexterity()

	// Shoulder
	shaman.applyShoulderRuneEffect()

	// Cloak
	shaman.registerFeralSpiritCD()
	shaman.applyStormEarthAndFire()

	// Chest
	shaman.applyDualWieldSpec()
	shaman.applyShieldMastery()
	shaman.applyTwoHandedMastery()
	shaman.applyOverload()

	// Bracers
	shaman.applyStaticShocks()
	shaman.registerRollingThunder()
	shaman.registerRiptideSpell()

	// Hands
	shaman.registerWaterShieldSpell()
	shaman.registerLavaBurstSpell()
	shaman.applyLavaLash()
	shaman.applyMoltenBlast()

	// Waist
	shaman.applyFireNova()
	shaman.applyMaelstromWeapon()
	shaman.applyPowerSurge()

	// Legs
	shaman.applyAncestralGuidance()
	shaman.applyWayOfEarth()
	shaman.registerEarthShieldSpell()

	// Feet
	shaman.applyAncestralAwakening()
	shaman.applySpiritOfTheAlpha()
}

func (shaman *Shaman) applyShoulderRuneEffect() {
	if shaman.Equipment.Shoulders().Rune == int32(proto.ShamanRune_RuneNone) {
		return
	}

	switch shaman.Equipment.Shoulders().Rune {
	// Elemental
	case int32(proto.ShamanRune_RuneShouldersVolcano):
		shaman.applyT1Elemental4PBonus()
	case int32(proto.ShamanRune_RuneShouldersRagingFlame):
		shaman.applyT1Elemental6PBonus()
	case int32(proto.ShamanRune_RuneShouldersElementalMaster):
		shaman.applyT2Elemental2PBonus()
	case int32(proto.ShamanRune_RuneShouldersTribesman):
		shaman.applyT2Elemental4PBonus()
	case int32(proto.ShamanRune_RuneShouldersSpiritGuide):
		shaman.applyT2Elemental6PBonus()
	case int32(proto.ShamanRune_RuneShouldersElder):
		shaman.applyTAQElemental2PBonus()
	case int32(proto.ShamanRune_RuneShouldersElements):
		shaman.applyTAQElemental4PBonus()
	case int32(proto.ShamanRune_RuneShouldersLavaSage):
		shaman.applyRAQElemental3PBonus()

	// Enhancement
	case int32(proto.ShamanRune_RuneShouldersRefinedShaman):
		shaman.applyT1Enhancement4PBonus()
	case int32(proto.ShamanRune_RuneShouldersChieftain):
		shaman.applyT1Enhancement6PBonus()
	case int32(proto.ShamanRune_RuneShouldersFurycharged):
		shaman.applyT2Enhancement2PBonus()
	case int32(proto.ShamanRune_RuneShouldersStormbreaker):
		shaman.applyT2Enhancement4PBonus()
	case int32(proto.ShamanRune_RuneShouldersTempest):
		shaman.applyT2Enhancement6PBonus()
	case int32(proto.ShamanRune_RuneShouldersSeismicSmasher):
		shaman.applyTAQEnhancement2PBonus()
	case int32(proto.ShamanRune_RuneShouldersFlamebringer):
		shaman.applyTAQEnhancement4PBonus()

	// Restoration
	case int32(proto.ShamanRune_RuneShouldersWaterWalker):
		shaman.applyT2Restoration2PBonus()
	case int32(proto.ShamanRune_RuneShouldersStormtender):
		shaman.applyT2Restoration4PBonus()
	case int32(proto.ShamanRune_RuneShouldersElementalSeer):
		shaman.applyT2Restoration6PBonus()

	// Tank
	case int32(proto.ShamanRune_RuneShouldersWindwalker):
		shaman.applyT1Tank2PBonus()
	case int32(proto.ShamanRune_RuneShouldersShieldMaster):
		shaman.applyT1Tank4PBonus()
	case int32(proto.ShamanRune_RuneShouldersTotemicProtector):
		shaman.applyT1Tank6PBonus()
	case int32(proto.ShamanRune_RuneShouldersShockAbsorber):
		shaman.applyT2Tank2PBonus()
	case int32(proto.ShamanRune_RuneShouldersSpiritualBulwark):
		shaman.applyT2Tank4PBonus()
	case int32(proto.ShamanRune_RuneShouldersMaelstrombringer):
		shaman.applyT2Tank6PBonus()
	case int32(proto.ShamanRune_RuneShouldersAncestralWarden):
		shaman.applyZGTank3PBonus()
	case int32(proto.ShamanRune_RuneShouldersCorrupt):
		shaman.applyZGTank5PBonus()
	case int32(proto.ShamanRune_RuneShouldersLavaWalker):
		shaman.applyTAQTank2PBonus()
	case int32(proto.ShamanRune_RuneShouldersTrueAlpha):
		shaman.applyTAQTank4PBonus()
	}
}

var BurnFlameShockTargetCount = int32(5)
var BurnFlameShockBonusTicks = int32(2)
var BurnSpellPowerPerLevel = int32(2)

func (shaman *Shaman) applyBurn() {
	if !shaman.HasRune(proto.ShamanRune_RuneHelmBurn) {
		return
	}

	if shaman.Consumes.MainHandImbue == proto.WeaponImbue_FlametongueWeapon {
		shaman.AddStatDependency(stats.Intellect, stats.SpellDamage, 1.50)
	}

	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanFlameShock,
		Kind:      core.SpellMod_DamageDone_Flat,
		IntValue:  100,
	})

	// Other parts of burn are handled in flame_shock.go
}

func (shaman *Shaman) applyMentalDexterity() {
	if !shaman.HasRune(proto.ShamanRune_RuneHelmMentalDexterity) {
		return
	}

	intToApStatDep := shaman.NewDynamicStatDependency(stats.Intellect, stats.AttackPower, 1.50)
	apToSpStatDep := shaman.NewDynamicStatDependency(stats.AttackPower, stats.SpellDamage, .35)

	procAura := shaman.RegisterAura(core.Aura{
		Label:    "Mental Dexterity Proc",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneHelmMentalDexterity)},
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, intToApStatDep)
			aura.Unit.EnableDynamicStatDep(sim, apToSpStatDep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, intToApStatDep)
			aura.Unit.DisableDynamicStatDep(sim, apToSpStatDep)
		},
	})

	// Hidden Aura
	shaman.RegisterAura(core.Aura{
		Label:    "Mental Dexterity",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell == shaman.StormstrikeMH {
				procAura.Activate(sim)
			}
		},
	})
}

func (shaman *Shaman) applyStormEarthAndFire() {
	if !shaman.HasRune(proto.ShamanRune_RuneCloakStormEarthAndFire) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: "Storm, Earth, and Fire",
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanFlameShock,
		Kind:      core.SpellMod_PeriodicDamageDone_Flat,
		IntValue:  60,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanChainLightning,
		Kind:      core.SpellMod_Cooldown_Multi_Flat,
		IntValue:  -50,
	}))

}

func (shaman *Shaman) applyDualWieldSpec() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestDualWieldSpec) || !shaman.HasMHWeapon() || !shaman.HasOHWeapon() {
		return
	}

	shaman.AutoAttacks.OHConfig().DamageMultiplier *= 1.60

	meleeHit := float64(core.MeleeHitRatingPerHitChance * 5)
	spellHit := float64(core.SpellHitRatingPerHitChance * 5)

	shaman.AddStat(stats.MeleeHit, meleeHit)
	shaman.AddStat(stats.SpellHit, spellHit)

	dwBonusApplied := true

	shaman.RegisterAura(core.Aura{
		Label:    "DW Spec Trigger",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneChestDualWieldSpec)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		// Perform additional checks for later weapon-swapping
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
				return
			}

			if shaman.HasMHWeapon() && shaman.HasOHWeapon() {
				if dwBonusApplied {
					return
				} else {
					shaman.AddStat(stats.MeleeHit, meleeHit)
					shaman.AddStat(stats.SpellHit, spellHit)
				}
			} else {
				shaman.AddStat(stats.MeleeHit, -1*meleeHit)
				shaman.AddStat(stats.SpellHit, -1*spellHit)
				dwBonusApplied = false
			}
		},
	})
}

func (shaman *Shaman) applyShieldMastery() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestShieldMastery) {
		return
	}

	defendersResolveAura := core.DefendersResolveSpellDamage(shaman.GetCharacter(), 4)

	shaman.AddStat(stats.Block, 10)
	shaman.PseudoStats.BlockValueMultiplier = 1.15

	actionId := core.ActionID{SpellID: int32(proto.ShamanRune_RuneChestShieldMastery)}
	manaMetrics := shaman.NewManaMetrics(actionId)
	procManaReturn := 0.08
	armorPerStack := shaman.Equipment.OffHand().Stats[stats.Armor] * 0.3

	shaman.ShieldMasteryAura = shaman.RegisterAura(core.Aura{
		Label:     "Shield Mastery Block",
		ActionID:  core.ActionID{SpellID: 408525},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnRefresh: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AddMana(sim, shaman.MaxMana()*procManaReturn, manaMetrics)
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			shaman.AddStatDynamic(sim, stats.Armor, armorPerStack*float64(newStacks-oldStacks))
		},
	})

	affectedSpellClassMasks := ClassSpellMask_ShamanEarthShock | ClassSpellMask_ShamanFlameShock | ClassSpellMask_ShamanFrostShock
	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: "Shield Mastery Trigger",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidBlock() {
				shaman.ShieldMasteryAura.Activate(sim)
				shaman.ShieldMasteryAura.AddStack(sim)
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.Matches(affectedSpellClassMasks) {
				if stacks := int32(shaman.GetStat(stats.Defense)); stacks > 0 {
					//Aura.Activate takes care of refreshing if the aura is already active
					defendersResolveAura.Activate(sim)
					if defendersResolveAura.GetStacks() != stacks {
						defendersResolveAura.SetStacks(sim, stacks)
					}
				}
			}
		},
	}))
}

func (shaman *Shaman) applyTwoHandedMastery() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestTwoHandedMastery) {
		return
	}

	procSpellId := int32(436365)

	// Two-handed mastery gives +15% AP, +30% attack speed, and +10% spell hit
	attackSpeedMultiplier := 1.5
	apMultiplier := 1.15
	spellHitIncrease := core.SpellHitRatingPerHitChance * 10.0

	statDep := shaman.NewDynamicMultiplyStat(stats.AttackPower, apMultiplier)
	procAura := shaman.RegisterAura(core.Aura{
		Label:    "Two-Handed Mastery Proc",
		ActionID: core.ActionID{SpellID: procSpellId},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, statDep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, statDep)
		},
	}).AttachStatBuff(stats.SpellHit, spellHitIncrease).AttachMultiplyAttackSpeed(&shaman.Unit, attackSpeedMultiplier)

	shaman.RegisterAura(core.Aura{
		Label:    "Two-Handed Mastery Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
				return
			}

			if shaman.MainHand().HandType == proto.HandType_HandTypeTwoHand {
				procAura.Activate(sim)
			} else {
				procAura.Deactivate(sim)
			}
		},
	})
}

func (shaman *Shaman) applyStaticShocks() {
	if !shaman.HasRune(proto.ShamanRune_RuneBracersStaticShock) {
		return
	}

	// DW chance base doubled by using a 2-handed weapon
	shaman.staticSHocksProcChance = .06

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: "Static Shocks",
		OnInit: func(staticShockAura *core.Aura, sim *core.Simulation) {
			for _, aura := range shaman.LightningShieldAuras {
				if aura == nil {
					continue
				}

				oldOnGain := aura.OnGain
				aura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
					oldOnGain(aura, sim)
					staticShockAura.Activate(sim)
				}

				oldOnExpire := aura.OnExpire
				aura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
					oldOnExpire(aura, sim)
					staticShockAura.Deactivate(sim)
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if shaman.ActiveShieldAura == nil || !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.Landed() {
				return
			}

			staticShockProcChance := core.TernaryFloat64(shaman.MainHand().HandType == proto.HandType_HandTypeTwoHand, shaman.staticSHocksProcChance*2, shaman.staticSHocksProcChance)
			if sim.RandomFloat("Static Shock") < staticShockProcChance {
				shaman.LightningShieldProcs[shaman.ActiveShield.Rank].Cast(sim, result.Target)
			}
		},
	}))
}

const MaelstromWeaponBaseStacks = 5
const MaelstromWeaponSplits = 11 // 0-5 stacks + 5 more if using 6pT3.5

func (shaman *Shaman) applyMaelstromWeapon() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
		return
	}

	// Chance increased by 50% while your main hand weapon is enchanted with Windfury Weapon and by another 50% if wielding a two-handed weapon.
	// Base PPM is 10
	ppm := 10.0
	if shaman.GetCharacter().Consumes.MainHandImbue == proto.WeaponImbue_WindfuryWeapon {
		ppm += 5
	}
	if shaman.MainHand().HandType == proto.HandType_HandTypeTwoHand {
		ppm += 5
	}

	shaman.MaelstromWeaponClassMask = ClassSpellMask_ShamanLightningBolt | ClassSpellMask_ShamanChainLightning | ClassSpellMask_ShamanLesserHealingWave | ClassSpellMask_ShamanLavaBurst

	castTimeMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Pct,
		ClassMask: shaman.MaelstromWeaponClassMask,
	})

	costMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: shaman.MaelstromWeaponClassMask,
	})

	shaman.MaelstromWeaponSpellMods = []*core.SpellMod{castTimeMod, costMod}

	shaman.MaelstromWeaponAura = shaman.RegisterAura(core.Aura{
		Label:     "MaelstromWeapon Proc",
		ActionID:  core.ActionID{SpellID: 408505},
		Duration:  time.Second * 30,
		MaxStacks: MaelstromWeaponBaseStacks,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
			costMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
			costMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			if newStacks > MaelstromWeaponBaseStacks {
				return
			}

			castTimeMod.UpdateFloatValue(-0.20 * float64(newStacks))
			costMod.UpdateIntValue(-20 * int64(newStacks))
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(shaman.MaelstromWeaponClassMask) {
				aura.RemoveStacks(sim, min(MaelstromWeaponBaseStacks, aura.GetStacks()))
			}
		},
	})

	shaman.maelstromWeaponPPMM = shaman.AutoAttacks.NewPPMManager(ppm, core.ProcMaskMelee)

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:              "Maelstrom Weapon Trigger",
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		ProcMask:          core.ProcMaskMelee,
		SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
		CanProcFromProcs:  true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if shaman.maelstromWeaponPPMM.Proc(sim, spell.ProcMask, "Maelstrom Weapon") {
				shaman.MaelstromWeaponAura.Activate(sim)
				shaman.MaelstromWeaponAura.AddStack(sim)
			}
		},
	})
}

func (shaman *Shaman) applyPowerSurge() {
	shaman.powerSurgeProcChance = 0.05

	// We want to create the power surge damage aura all the time because it's used by the T1 Ele 4P and can be triggered without the rune
	var affectedDamageSpells []*core.Spell
	shaman.PowerSurgeDamageAura = shaman.RegisterAura(core.Aura{
		Label:    "Power Surge Proc (Damage)",
		ActionID: core.ActionID{SpellID: 415105},
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedDamageSpells = shaman.GetSpellsMatchingClassMask(ClassSpellMask_ShamanChainLightning | ClassSpellMask_ShamanLavaBurst)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedDamageSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier -= 1
				if spell.CD.Timer != nil {
					spell.CD.Reset()
				}
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedDamageSpells, func(spell *core.Spell) { spell.CastTimeMultiplier += 1 })
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_ShamanLavaBurst|ClassSpellMask_ShamanChainLightning) && !spell.ProcMask.Matches(core.ProcMaskSpellDamageProc) {
				aura.Deactivate(sim)
			}
		},
	})

	if !shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge) {
		return
	}

	shaman.PowerSurgeHealAura = shaman.RegisterAura(core.Aura{
		Label:    "Power Surge Proc (Heal)",
		ActionID: core.ActionID{SpellID: 468526},
		Duration: time.Second * 10,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_ShamanChainHeal) && !spell.ProcMask.Matches(core.ProcMaskSpellDamageProc) {
				aura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  ClassSpellMask_ShamanChainHeal,
		FloatValue: -1,
	})

	statDep := shaman.NewDynamicStatDependency(stats.Intellect, stats.MP5, .15)
	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: "Power Surge",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_ShamanFlameShock) && sim.Proc(shaman.powerSurgeProcChance, "Power Surge Proc") {
				shaman.PowerSurgeDamageAura.Activate(sim)
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_ShamanFlameShock) && sim.Proc(shaman.powerSurgeProcChance, "Power Surge Proc") {
				shaman.PowerSurgeDamageAura.Activate(sim)
			}
		},
	}).AttachStatDependency(statDep))
}

func (shaman *Shaman) applyWayOfEarth() {
	if !shaman.HasRune(proto.ShamanRune_RuneLegsWayOfEarth) {
		return
	}

	// Way of Earth only activates if you have Rockbiter Weapon on your mainhand and a shield in your offhand
	if shaman.Consumes.MainHandImbue != proto.WeaponImbue_RockbiterWeapon || shaman.OffHand().WeaponType != proto.WeaponType_WeaponTypeShield {
		return
	}

	healthDep := shaman.NewDynamicMultiplyStat(stats.Health, 1.3)

	core.MakePermanent(
		shaman.RegisterAura(core.Aura{
			ActionID:   core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsWayOfEarth)},
			BuildPhase: core.CharacterBuildPhaseBuffs,
			Label:      "Way of Earth",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.DamageTakenMultiplier *= .9
				shaman.PseudoStats.ReducedCritTakenChance += 6
				shaman.PseudoStats.ThreatMultiplier *= 1.65
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.DamageTakenMultiplier /= .9
				shaman.PseudoStats.ReducedCritTakenChance -= 6
				shaman.PseudoStats.ThreatMultiplier /= 1.65
			},
		}).AttachStatDependency(healthDep),
	)
}

// https://www.wowhead.com/classic/spell=408696/spirit-of-the-alpha
func (shaman *Shaman) applySpiritOfTheAlpha() {
	hasSpiritOfTheAlpha := shaman.HasRune(proto.ShamanRune_RuneFeetSpiritOfTheAlpha)

	shaman.SpiritOfTheAlphaAura = core.SpiritOfTheAlphaAura(&shaman.Unit)
	if (hasSpiritOfTheAlpha && shaman.IsTanking()) || shaman.IndividualBuffs.SpiritOfTheAlpha {
		core.MakePermanent(shaman.SpiritOfTheAlphaAura)
	} else if hasSpiritOfTheAlpha {
		shaman.LoyalBetaAura = core.MakePermanent(shaman.RegisterAura(core.Aura{
			Label:    "Loyal Beta",
			ActionID: core.ActionID{SpellID: 443320},
		}).AttachMultiplicativePseudoStatBuff(
			&shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical], 1.05,
		).AttachMultiplicativePseudoStatBuff(
			&shaman.PseudoStats.ThreatMultiplier, .70,
		))
	}
}
