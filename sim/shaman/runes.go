package shaman

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (shaman *Shaman) ApplyRunes() {
	// Helm
	shaman.applyBurn()
	shaman.applyMentalDexterity()

	// Cloak
	shaman.registerFeralSpiritCD()

	// Chest
	shaman.applyDualWieldSpec()
	shaman.applyShieldMastery()
	shaman.applyTwoHandedMastery()

	// Bracers
	shaman.applyRollingThunder()

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

	// Feet
	shaman.applyAncestralAwakening()
	shaman.applySpiritOfTheAlpha()
}

var BurnFlameShockTargetCount = int32(5)
var BurnFlameShockDamageBonus = 1.0
var BurnFlameShockBonusTicks = int32(2)
var BurnSpellPowerPerLevel = int32(2)

func (shaman *Shaman) applyBurn() {
	if !shaman.HasRune(proto.ShamanRune_RuneHelmBurn) {
		return
	}

	if shaman.Consumes.MainHandImbue == proto.WeaponImbue_FlametongueWeapon || shaman.Consumes.OffHandImbue == proto.WeaponImbue_FlametongueWeapon {
		shaman.AddStat(stats.SpellDamage, float64(BurnSpellPowerPerLevel*shaman.Level))
	}

	// Other parts of burn are handled in flame_shock.go
}

func (shaman *Shaman) burnFlameShockDamageMultiplier() float64 {
	return core.TernaryFloat64(shaman.HasRune(proto.ShamanRune_RuneHelmBurn), BurnFlameShockDamageBonus, 0)
}

func (shaman *Shaman) applyMentalDexterity() {
	if !shaman.HasRune(proto.ShamanRune_RuneHelmMentalDexterity) {
		return
	}

	intToApStatDep := shaman.NewDynamicStatDependency(stats.Intellect, stats.AttackPower, .65)
	apToSpStatDep := shaman.NewDynamicStatDependency(stats.AttackPower, stats.SpellDamage, .20)

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
			if result.Landed() && (spell == shaman.LavaLash || spell == shaman.StormstrikeMH) {
				procAura.Activate(sim)
			}
		},
	})
}

func (shaman *Shaman) applyDualWieldSpec() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestDualWieldSpec) || !shaman.HasMHWeapon() || !shaman.HasOHWeapon() {
		return
	}

	shaman.AutoAttacks.OHConfig().DamageMultiplier *= 1.5

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

	defendersResolveAura := core.DefendersResolveSpellDamage(shaman.GetCharacter())

	has4PEarthfuryResolve := shaman.HasSetBonus(ItemSetEarthfuryResolve, 4)

	shaman.AddStat(stats.Block, 10)
	shaman.PseudoStats.BlockValueMultiplier = 1.15

	actionId := core.ActionID{SpellID: int32(proto.ShamanRune_RuneChestShieldMastery)}
	manaMetrics := shaman.NewManaMetrics(actionId)
	procManaReturn := 0.08
	armorPerStack := shaman.Equipment.OffHand().Stats[stats.Armor] * 0.3

	blockProcAura := shaman.RegisterAura(core.Aura{
		Label:     "Shield Mastery Block",
		ActionID:  core.ActionID{SpellID: 408525},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			shaman.AddStatDynamic(sim, stats.Armor, armorPerStack*float64(newStacks-oldStacks))
		},
	})

	affectedSpellcodes := []int32{SpellCode_ShamanEarthShock, SpellCode_ShamanFlameShock, SpellCode_ShamanFrostShock}
	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: "Shield Mastery Trigger",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeBlock) || (has4PEarthfuryResolve && (result.Outcome.Matches(core.OutcomeParry) || result.Outcome.Matches(core.OutcomeDodge))) {
				shaman.AddMana(sim, shaman.MaxMana()*procManaReturn, manaMetrics)
				blockProcAura.Activate(sim)
				blockProcAura.AddStack(sim)
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && slices.Contains(affectedSpellcodes, spell.SpellCode) {
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

	// Two-handed mastery gives +10% AP, +30% attack speed, and +10% spell hit
	attackSpeedMultiplier := 1.5
	apMultiplier := 1.1
	spellHitIncrease := core.SpellHitRatingPerHitChance * 10.0

	statDep := shaman.NewDynamicMultiplyStat(stats.AttackPower, apMultiplier)
	procAura := shaman.RegisterAura(core.Aura{
		Label:    "Two-Handed Mastery Proc",
		ActionID: core.ActionID{SpellID: procSpellId},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, attackSpeedMultiplier)
			shaman.AddStatDynamic(sim, stats.SpellHit, spellHitIncrease)
			aura.Unit.EnableDynamicStatDep(sim, statDep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyAttackSpeed(sim, 1/attackSpeedMultiplier)
			shaman.AddStatDynamic(sim, stats.SpellHit, -1*spellHitIncrease)
			aura.Unit.DisableDynamicStatDep(sim, statDep)
		},
	})

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

var RollingThunderProcChance = .50

func (shaman *Shaman) applyRollingThunder() {
	if !shaman.HasRune(proto.ShamanRune_RuneBracersRollingThunder) {
		return
	}

	impLightningShieldBonus := 1 + []float64{0, .05, .10, .15}[shaman.Talents.ImprovedLightningShield]

	// Casts handled in lightning_shield.go
	shaman.RollingThunder = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 432129},
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagShaman,

		BonusCritRating: float64(shaman.Talents.TidalMastery) * core.CritRatingPerCritChance,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if shaman.ActiveShield == nil || shaman.ActiveShield.SpellCode != SpellCode_ShamanLightningShield {
				return
			}

			rank := shaman.ActiveShield.Rank
			chargeDamage := LightningShieldBaseDamage[rank]*impLightningShieldBonus + LightningShieldSpellCoef[rank]*shaman.LightningShieldProcs[rank].BonusDamage()
			spell.CalcAndDealDamage(sim, target, chargeDamage, spell.OutcomeMagicCrit)
		},
	})
}

func (shaman *Shaman) rollRollingThunderCharge(sim *core.Simulation) {
	if shaman.ActiveShield != nil && shaman.ActiveShield.SpellCode == SpellCode_ShamanLightningShield && shaman.ActiveShieldAura.IsActive() && sim.Proc(RollingThunderProcChance, "Rolling Thunder") {
		shaman.ActiveShieldAura.AddStack(sim)
	}
}

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

	var affectedSpells []*core.Spell
	shaman.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagMaelstrom) {
			affectedSpells = append(affectedSpells, spell)
		}
	})

	shaman.MaelstromWeaponAura = shaman.RegisterAura(core.Aura{
		Label:     "MaelstromWeapon Proc",
		ActionID:  core.ActionID{SpellID: 408505},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			multDiff := 0.2 * float64(newStacks-oldStacks)
			for _, spell := range affectedSpells {
				spell.CastTimeMultiplier -= multDiff
				spell.CostMultiplier -= multDiff
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Flags.Matches(SpellFlagMaelstrom) {
				shaman.MaelstromWeaponAura.Deactivate(sim)
			}
		},
	})

	ppmm := shaman.AutoAttacks.NewPPMManager(ppm, core.ProcMaskMelee)

	// This aura is hidden, just applies stacks of the proc aura.
	shaman.RegisterAura(core.Aura{
		Label:    "MaelstromWeapon",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if ppmm.Proc(sim, spell.ProcMask, "Maelstrom Weapon") {
				shaman.MaelstromWeaponAura.Activate(sim)
				shaman.MaelstromWeaponAura.AddStack(sim)
			}
		},
	})
}

const ShamanPowerSurgeProcChance = .05

func (shaman *Shaman) applyPowerSurge() {
	// TODO: Figure out how this actually works because the 2024-02-27 tuning notes make it sound like
	// this is not just a fully passive stat boost
	if shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge) {
		shaman.AddStat(stats.MP5, shaman.GetStat(stats.Intellect)*.15)
	}

	// We want to create the power surge aura all the time because it's used by the T1 Ele 4P and can be triggered without the rune

	var affectedSpells []*core.Spell
	var affectedSpellCodes = []int32{
		SpellCode_ShamanChainLightning,
		SpellCode_ShamanChainHeal,
		SpellCode_ShamanLavaBurst,
	}

	shaman.PowerSurgeAura = shaman.RegisterAura(core.Aura{
		Label:    "Power Surge Proc",
		ActionID: core.ActionID{SpellID: 440285},
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(
				core.Flatten([][]*core.Spell{
					shaman.ChainLightning,
					shaman.ChainHeal,
					{shaman.LavaBurst},
				}), func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier -= 1
				if spell.CD.Timer != nil {
					spell.CD.Reset()
				}
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) { spell.CastTimeMultiplier += 1 })
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !slices.Contains(affectedSpellCodes, spell.SpellCode) {
				return
			}
			aura.Deactivate(sim)
		},
	})
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

	shaman.RegisterAura(core.Aura{
		Label:    "Way of Earth",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsWayOfEarth)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.EnableDynamicStatDep(sim, healthDep)
			shaman.PseudoStats.DamageTakenMultiplier *= .9
			shaman.PseudoStats.ReducedCritTakenChance += 6
			shaman.PseudoStats.ThreatMultiplier *= 1.65
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.DisableDynamicStatDep(sim, healthDep)
			shaman.PseudoStats.DamageTakenMultiplier /= .9
			shaman.PseudoStats.ReducedCritTakenChance -= 6
			shaman.PseudoStats.ThreatMultiplier /= 1.65
		},
	})
}

// https://www.wowhead.com/classic/spell=408696/spirit-of-the-alpha
func (shaman *Shaman) applySpiritOfTheAlpha() {
	if !shaman.HasRune(proto.ShamanRune_RuneFeetSpiritOfTheAlpha) {
		return
	}

	if shaman.IsTanking() {
		shaman.RegisterAura(core.Aura{
			Label:    "Spirit of the Alpha",
			ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneFeetSpiritOfTheAlpha)},
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.ThreatMultiplier *= 1.45
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.ThreatMultiplier /= 1.45
			},
		})
	} else {
		shaman.RegisterAura(core.Aura{
			Label:    "Loyal Beta",
			ActionID: core.ActionID{SpellID: 443320},
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.05
				shaman.PseudoStats.ThreatMultiplier *= .70
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.05
				shaman.PseudoStats.ThreatMultiplier /= .70
			},
		})
	}
}
