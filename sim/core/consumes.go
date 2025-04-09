package core

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Registers all consume-related effects to the Agent.
// TODO: Classic Consumes
func applyConsumeEffects(agent Agent) {
	character := agent.GetCharacter()
	consumes := character.Consumes

	if consumes == nil {
		return
	}

	applyFlaskConsumes(character, consumes)
	applyWeaponImbueConsumes(character, consumes)
	applyFoodConsumes(character, consumes)
	applyDefensiveBuffConsumes(character, consumes)
	applyPhysicalBuffConsumes(character, consumes)
	applySpellBuffConsumes(character, consumes)
	applySealOfTheDawnBuffConsumes(character, consumes)
	applyZanzaBuffConsumes(character, consumes)
	applyMiscConsumes(character, consumes.MiscConsumes)
	applyEnchantingConsumes(character, consumes)

	registerPotionCD(agent, consumes)
	registerConjuredCD(agent, consumes)
	registerMildlyIrradiatedRejuvCD(agent, consumes)
	registerExplosivesCD(agent, consumes)
}

func ApplyPetConsumeEffects(pet *Character, ownerConsumes *proto.Consumes) {
	pet.AddStat(stats.AttackPower, []float64{0, 40}[ownerConsumes.PetAttackPowerConsumable])
	pet.AddStat(stats.Agility, []float64{0, 17, 13, 9, 5}[ownerConsumes.PetAgilityConsumable])
	pet.AddStat(stats.Strength, []float64{0, 30, 17, 13, 9, 5}[ownerConsumes.PetStrengthConsumable])
}

///////////////////////////////////////////////////////////////////////////
//                             Flasks
///////////////////////////////////////////////////////////////////////////

func applyFlaskConsumes(character *Character, consumes *proto.Consumes) {
	if consumes.Flask == proto.Flask_FlaskUnknown {
		return
	}

	switch consumes.Flask {
	case proto.Flask_FlaskOfDistilledWisdom:
		character.AddStats(stats.Stats{
			stats.Mana: 2000,
		})
	case proto.Flask_FlaskOfUnyieldingSorrow:
		character.AddStats(stats.Stats{
			stats.SpellDamage:  27,
			stats.HealingPower: 80,
			stats.MP5:          12,
		})
	case proto.Flask_FlaskOfAncientKnowledge:
		character.AddStats(stats.Stats{
			stats.SpellPower: 180,
		})
	case proto.Flask_FlaskOfTheOldGods:
		character.AddStats(stats.Stats{
			stats.Stamina: 100,
			stats.Defense: 10,
		})
	case proto.Flask_FlaskOfSupremePower:
		character.AddStats(stats.Stats{
			stats.SpellPower: 150,
		})
	case proto.Flask_FlaskOfTheTitans:
		character.AddStats(stats.Stats{
			stats.Health: 1200,
		})
	case proto.Flask_FlaskOfChromaticResistance:
		character.AddResistances(25)
	case proto.Flask_FlaskOfRestlessDreams:
		character.AddStats(stats.Stats{
			// +30 Spell Damage, +45 Healing Power, +12 MP5
			stats.SpellDamage:  30,
			stats.HealingPower: 15,
			stats.MP5:          12,
		})
	case proto.Flask_FlaskOfEverlastingNightmares:
		character.AddStats(stats.Stats{
			stats.AttackPower:       45,
			stats.RangedAttackPower: 45,
		})
	case proto.Flask_FlaskOfMadness:
		character.AddStats(stats.Stats{
			stats.AttackPower:       50,
			stats.RangedAttackPower: 50,
		})
	}
}

///////////////////////////////////////////////////////////////////////////
//                             Weapon Imbues
///////////////////////////////////////////////////////////////////////////

func applyWeaponImbueConsumes(character *Character, consumes *proto.Consumes) {
	// There must be a nicer way to do this...
	shadowOilIcd := Cooldown{
		Timer:    character.NewTimer(),
		Duration: time.Second * 10,
	}

	if character.HasMHWeapon() {
		addImbueStats(character, consumes.MainHandImbue, true, shadowOilIcd)
	}
	if character.HasOHItem() {
		addImbueStats(character, consumes.OffHandImbue, false, shadowOilIcd)
	}
}

func addImbueStats(character *Character, imbue proto.WeaponImbue, isMh bool, shadowOilIcd Cooldown) {
	if imbue != proto.WeaponImbue_WeaponImbueUnknown {
		switch imbue {
		// Wizard Oils
		case proto.WeaponImbue_MinorWizardOil:
			character.AddStats(stats.Stats{
				stats.SpellPower: 8,
			})
		case proto.WeaponImbue_LesserWizardOil:
			character.AddStats(stats.Stats{
				stats.SpellPower: 16,
			})
		case proto.WeaponImbue_WizardOil:
			character.AddStats(stats.Stats{
				stats.SpellPower: 24,
			})
		case proto.WeaponImbue_BrilliantWizardOil:
			character.AddStats(stats.Stats{
				stats.SpellPower: 36,
				stats.SpellCrit:  1 * SpellCritRatingPerCritChance,
			})
		case proto.WeaponImbue_EnchantedRepellent:
			character.AddStats(stats.Stats{
				stats.SpellPower: 45,
				stats.SpellCrit:  1 * SpellCritRatingPerCritChance,
			})
		case proto.WeaponImbue_BlessedWizardOil:
			if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
				character.PseudoStats.MobTypeSpellPower += 60
			}

		// Mana Oils
		case proto.WeaponImbue_MinorManaOil:
			character.AddStats(stats.Stats{
				stats.MP5: 4,
			})
		case proto.WeaponImbue_LesserManaOil:
			character.AddStats(stats.Stats{
				stats.MP5: 8,
			})
		case proto.WeaponImbue_BrilliantManaOil:
			character.AddStats(stats.Stats{
				stats.MP5:          12,
				stats.HealingPower: 25,
			})
		case proto.WeaponImbue_BlackfathomManaOil:
			character.AddStats(stats.Stats{
				stats.MP5:      12,
				stats.SpellHit: 2 * SpellHitRatingPerHitChance,
			})

		// Shield Oil
		case proto.WeaponImbue_ConductiveShieldCoating:
			character.AddStat(stats.SpellPower, 24)
		case proto.WeaponImbue_MagnificentTrollshine:
			character.AddStats(stats.Stats{
				stats.SpellPower: 36,
				stats.SpellCrit:  1 * CritRatingPerCritChance,
			})

		// Sharpening Stones
		case proto.WeaponImbue_SolidSharpeningStone:
			weapon := character.AutoAttacks.MH()
			if !isMh {
				weapon = character.AutoAttacks.OH()
			}
			weapon.BaseDamageMin += 6
			weapon.BaseDamageMax += 6
		case proto.WeaponImbue_DenseSharpeningStone:
			weapon := character.AutoAttacks.MH()
			if !isMh {
				weapon = character.AutoAttacks.OH()
			}
			weapon.BaseDamageMin += 8
			weapon.BaseDamageMax += 8
		case proto.WeaponImbue_ElementalSharpeningStone:
			character.AddStats(stats.Stats{
				stats.MeleeCrit: 2 * CritRatingPerCritChance,
			})
			character.AddBonusRangedCritRating(-2.0)
		case proto.WeaponImbue_BlackfathomSharpeningStone:
			character.AddStats(stats.Stats{
				stats.MeleeHit: 2 * MeleeHitRatingPerHitChance,
			})
		case proto.WeaponImbue_ConsecratedSharpeningStone:
			if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
				character.PseudoStats.MobTypeAttackPower += 100
			}
		case proto.WeaponImbue_WeightedConsecratedSharpeningStone:
			if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
				character.PseudoStats.MobTypeAttackPower += 200
			}

		// Weightstones
		case proto.WeaponImbue_SolidWeightstone:
			weapon := character.AutoAttacks.MH()
			if !isMh {
				weapon = character.AutoAttacks.OH()
			}
			weapon.BaseDamageMin += 6
			weapon.BaseDamageMax += 6
		case proto.WeaponImbue_DenseWeightstone:
			weapon := character.AutoAttacks.MH()
			if !isMh {
				weapon = character.AutoAttacks.OH()
			}
			weapon.BaseDamageMin += 8
			weapon.BaseDamageMax += 8

		// Windfury
		case proto.WeaponImbue_WildStrikes:
			//protect against double application if wild strikes is selected by a feral in sim settings
			if !character.HasRuneById(int32(proto.DruidRune_RuneChestWildStrikes)) {
				ApplyWildStrikes(character)
			}
		case proto.WeaponImbue_Windfury:
			ApplyWindfury(character)
		case proto.WeaponImbue_ShadowOil:
			registerShadowOil(character, isMh, shadowOilIcd)
		case proto.WeaponImbue_FrostOil:
			registerFrostOil(character, isMh)
		}
	}
}

func registerShadowOil(character *Character, isMh bool, icd Cooldown) {
	procChance := 0.15

	procSpell := character.GetOrRegisterSpell(SpellConfig{
		ActionID:    ActionID{SpellID: 1382},
		SpellSchool: SpellSchoolShadow,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagNoOnCastComplete | SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.56,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			damage := sim.Roll(52, 61)
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		},
	})

	label := " MH"
	procMask := ProcMaskMeleeMH
	if !isMh {
		label = " OH"
		procMask = ProcMaskMeleeOH
	}

	MakePermanent(character.GetOrRegisterAura(Aura{
		Label: "Shadow Oil" + label,
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if !result.Landed() {
				return
			}

			if !spell.ProcMask.Matches(procMask) {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Shadow Oil") < procChance {
				icd.Use(sim)
				procSpell.Cast(sim, result.Target)
			}
		},
	}))
}

func registerFrostOil(character *Character, isMh bool) {
	procChance := 0.10

	procSpell := character.GetOrRegisterSpell(SpellConfig{
		ActionID:    ActionID{SpellID: 1191},
		SpellSchool: SpellSchoolFrost,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagNoOnCastComplete | SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.269,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			damage := sim.Roll(33, 38)
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		},
	})

	label := " MH"
	procMask := ProcMaskMeleeMHAuto
	if !isMh {
		label = " OH"
		procMask = ProcMaskMeleeOHAuto
	}

	MakePermanent(character.GetOrRegisterAura(Aura{
		Label: "Frost Oil" + label,
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if !result.Landed() {
				return
			}

			if !spell.ProcMask.Matches(procMask) {
				return
			}

			if sim.RandomFloat("Frost Oil") < procChance {
				procSpell.Cast(sim, result.Target)
			}
		},
	}))
}

///////////////////////////////////////////////////////////////////////////
//                             Food
///////////////////////////////////////////////////////////////////////////

func applyFoodConsumes(character *Character, consumes *proto.Consumes) {
	if consumes.Food != proto.Food_FoodUnknown {
		switch consumes.Food {
		case proto.Food_FoodHotWolfRibs:
			character.AddStats(stats.Stats{
				stats.Stamina: 8,
				stats.Spirit:  8,
			})
		case proto.Food_FoodSmokedSagefish:
			character.AddStats(stats.Stats{
				stats.MP5: 3,
			})
		case proto.Food_FoodSagefishDelight:
			character.AddStats(stats.Stats{
				stats.MP5: 6,
			})
		case proto.Food_FoodTenderWolfSteak:
			character.AddStats(stats.Stats{
				stats.Stamina: 12,
				stats.Spirit:  12,
			})
		case proto.Food_FoodGrilledSquid:
			character.AddStats(stats.Stats{
				stats.Agility: 10,
			})
		case proto.Food_FoodSmokedDesertDumpling:
			character.AddStats(stats.Stats{
				stats.Strength: 20,
			})
		case proto.Food_FoodNightfinSoup:
			character.AddStats(stats.Stats{
				stats.MP5: 8,
			})
		case proto.Food_FoodRunnTumTuberSurprise:
			character.AddStats(stats.Stats{
				stats.Intellect: 10,
			})
		case proto.Food_FoodDirgesKickChimaerokChops:
			character.AddStats(stats.Stats{
				stats.Stamina: 25,
			})
		case proto.Food_FoodBlessedSunfruitJuice:
			character.AddStats(stats.Stats{
				stats.Spirit: 10,
			})
		case proto.Food_FoodBlessSunfruit:
			character.AddStats(stats.Stats{
				stats.Strength: 10,
			})
		case proto.Food_FoodDarkclawBisque:
			character.AddStats(stats.Stats{
				stats.SpellDamage: 12,
			})
		case proto.Food_FoodSmokedRedgill:
			character.AddStats(stats.Stats{
				stats.HealingPower: 22,
			})
		case proto.Food_FoodProwlerSteak:
			character.AddStats(stats.Stats{
				stats.Strength: 25,
				stats.Stamina:  10,
			})
		case proto.Food_FoodFiletOFlank:
			character.AddStats(stats.Stats{
				stats.Agility: 25,
				stats.Stamina: 10,
			})
		case proto.Food_FoodSunriseOmelette:
			character.AddStats(stats.Stats{
				stats.SpellPower:   29,
				stats.HealingPower: 55,
				stats.Stamina:      10,
			})
		case proto.Food_FoodSpecklefinFeast:
			character.AddStats(stats.Stats{
				stats.AttackPower:  40,
				stats.SpellPower:   23,
				stats.HealingPower: 44,
				stats.Stamina:      10,
			})
		case proto.Food_FoodGrandLobsterBanquet:
			character.AddStats(stats.Stats{
				stats.AttackPower:  40,
				stats.SpellPower:   23,
				stats.HealingPower: 44,
				stats.Stamina:      10,
			})
		}
	}

	if consumes.Alcohol != proto.Alcohol_AlcoholUnknown {
		switch consumes.Alcohol {
		case proto.Alcohol_AlcoholRumseyRumBlackLabel:
			character.AddStats(stats.Stats{
				stats.Stamina: 15,
			})
		case proto.Alcohol_AlcoholGordokGreenGrog:
			character.AddStats(stats.Stats{
				stats.Stamina: 10,
			})
		case proto.Alcohol_AlcoholRumseyRumDark:
			character.AddStats(stats.Stats{
				stats.Stamina: 10,
			})
		case proto.Alcohol_AlcoholRumseyRumLight:
			character.AddStats(stats.Stats{
				stats.Stamina: 5,
			})
		case proto.Alcohol_AlcoholKreegsStoutBeatdown:
			character.AddStats(stats.Stats{
				stats.Spirit:    25,
				stats.Intellect: -5,
			})
		}
	}

	if consumes.DragonBreathChili {
		MakePermanent(DragonBreathChiliAura(character))
	}
}

func DragonBreathChiliAura(character *Character) *Aura {
	baseDamage := 60.0
	procChance := .05
	icd := Cooldown{
		Timer:    character.NewTimer(),
		Duration: time.Second * 10,
	}

	procSpell := character.RegisterSpell(SpellConfig{
		ActionID:    ActionID{SpellID: 15851},
		SpellSchool: SpellSchoolFire,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskSpellDamageProc | ProcMaskSpellProc,
		Flags:       SpellFlagNoOnCastComplete | SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			for _, aoeTarget := range sim.Environment.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	aura := character.GetOrRegisterAura(Aura{
		Label:    "Dragonbreath Chili",
		ActionID: ActionID{SpellID: 15852},
		Duration: NeverExpires,
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(ProcMaskMelee) {
				return
			}
			if icd.IsReady(sim) && sim.RandomFloat("Dragonbreath Chili") < procChance {
				icd.Use(sim)
				procSpell.Cast(sim, result.Target)
			}
		},
	})
	return aura
}

///////////////////////////////////////////////////////////////////////////
//                             Defensive Buff Consumes
///////////////////////////////////////////////////////////////////////////

func applyDefensiveBuffConsumes(character *Character, consumes *proto.Consumes) {
	if consumes.ArmorElixir != proto.ArmorElixir_ArmorElixirUnknown {
		switch consumes.ArmorElixir {
		case proto.ArmorElixir_ElixirOfTheIronside:
			character.AddStats(stats.Stats{
				stats.BonusArmor:       450,
				stats.Defense:          5,
				stats.NatureResistance: 15,
			})
		case proto.ArmorElixir_ElixirOfSuperiorDefense:
			character.AddStats(stats.Stats{
				stats.BonusArmor: 450,
			})
		case proto.ArmorElixir_ElixirOfGreaterDefense:
			character.AddStats(stats.Stats{
				stats.BonusArmor: 250,
			})
		case proto.ArmorElixir_ElixirOfDefense:
			character.AddStats(stats.Stats{
				stats.BonusArmor: 150,
			})
		case proto.ArmorElixir_ElixirOfMinorDefense:
			character.AddStats(stats.Stats{
				stats.BonusArmor: 50,
			})
		case proto.ArmorElixir_ScrollOfProtection:
			character.AddStats(BuffSpellByLevel[ScrollOfProtection][character.Level])
		}
	}

	if consumes.HealthElixir != proto.HealthElixir_HealthElixirUnknown {
		switch consumes.HealthElixir {
		case proto.HealthElixir_ElixirOfFortitude:
			character.AddStats(stats.Stats{
				stats.Health: 120,
			})
		case proto.HealthElixir_ElixirOfMinorFortitude:
			character.AddStats(stats.Stats{
				stats.Health: 27,
			})
		}
	}
}

///////////////////////////////////////////////////////////////////////////
//                             Physical Buff Consumes
///////////////////////////////////////////////////////////////////////////

func applyPhysicalBuffConsumes(character *Character, consumes *proto.Consumes) {
	if consumes.AttackPowerBuff != proto.AttackPowerBuff_AttackPowerBuffUnknown {
		switch consumes.AttackPowerBuff {
		case proto.AttackPowerBuff_JujuMight:
			character.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		case proto.AttackPowerBuff_WinterfallFirewater:
			character.AddStats(stats.Stats{
				stats.AttackPower: 35,
			})
		}
	}

	if consumes.AgilityElixir != proto.AgilityElixir_AgilityElixirUnknown {
		switch consumes.AgilityElixir {
		case proto.AgilityElixir_ElixirOfTheHoneyBadger:
			character.AddStats(stats.Stats{
				stats.Agility:          30,
				stats.MeleeCrit:        2 * CritRatingPerCritChance,
				stats.NatureResistance: 15,
			})
		case proto.AgilityElixir_ElixirOfTheMongoose:
			character.AddStats(stats.Stats{
				stats.Agility:   25,
				stats.MeleeCrit: 2 * CritRatingPerCritChance,
			})
		case proto.AgilityElixir_ElixirOfGreaterAgility:
			character.AddStats(stats.Stats{
				stats.Agility: 25,
			})
		case proto.AgilityElixir_ElixirOfAgility:
			character.AddStats(stats.Stats{
				stats.Agility: 15,
			})
		case proto.AgilityElixir_ElixirOfLesserAgility:
			character.AddStats(stats.Stats{
				stats.Agility: 8,
			})
		case proto.AgilityElixir_ScrollOfAgility:
			character.AddStats(BuffSpellByLevel[ScrollOfAgility][character.Level])
		}
	}

	if consumes.StrengthBuff != proto.StrengthBuff_StrengthBuffUnknown {
		switch consumes.StrengthBuff {
		case proto.StrengthBuff_JujuPower:
			character.AddStats(stats.Stats{
				stats.Strength: 30,
			})
		case proto.StrengthBuff_ElixirOfGiants:
			character.AddStats(stats.Stats{
				stats.Strength: 25,
			})
		case proto.StrengthBuff_ElixirOfOgresStrength:
			character.AddStats(stats.Stats{
				stats.Strength: 8,
			})
		case proto.StrengthBuff_ScrollOfStrength:
			character.AddStats(BuffSpellByLevel[ScrollOfStrength][character.Level])
		}
	}
}

///////////////////////////////////////////////////////////////////////////
//                             Spell Buff Consumes
///////////////////////////////////////////////////////////////////////////

func applySpellBuffConsumes(character *Character, consumes *proto.Consumes) {
	if consumes.SpellPowerBuff != proto.SpellPowerBuff_SpellPowerBuffUnknown {
		switch consumes.SpellPowerBuff {
		case proto.SpellPowerBuff_LesserArcaneElixir:
			character.AddStats(stats.Stats{
				stats.SpellDamage: 14,
			})
		case proto.SpellPowerBuff_ArcaneElixir:
			character.AddStats(stats.Stats{
				stats.SpellDamage: 20,
			})
		case proto.SpellPowerBuff_GreaterArcaneElixir:
			character.AddStats(stats.Stats{
				stats.SpellDamage: 35,
			})
		case proto.SpellPowerBuff_ElixirOfTheMageLord:
			character.AddStats(stats.Stats{
				stats.SpellDamage:      40,
				stats.NatureResistance: 15,
			})
		}
	}

	if consumes.FirePowerBuff != proto.FirePowerBuff_FirePowerBuffUnknown {
		switch consumes.FirePowerBuff {
		case proto.FirePowerBuff_ElixirOfFirepower:
			character.AddStats(stats.Stats{
				stats.FirePower: 10,
			})
		case proto.FirePowerBuff_ElixirOfGreaterFirepower:
			character.AddStats(stats.Stats{
				stats.FirePower: 40,
			})
		}
	}

	if consumes.ShadowPowerBuff != proto.ShadowPowerBuff_ShadowPowerBuffUnknown {
		switch consumes.ShadowPowerBuff {
		case proto.ShadowPowerBuff_ElixirOfShadowPower:
			character.AddStats(stats.Stats{
				stats.ShadowPower: 40,
			})
		}
	}

	if consumes.FrostPowerBuff != proto.FrostPowerBuff_FrostPowerBuffUnknown {
		switch consumes.FrostPowerBuff {
		case proto.FrostPowerBuff_ElixirOfFrostPower:
			character.AddStats(stats.Stats{
				stats.FrostPower: 15,
			})
		}
	}

	if consumes.ManaRegenElixir != proto.ManaRegenElixir_ManaRegenElixirUnknown {
		switch consumes.ManaRegenElixir {
		case proto.ManaRegenElixir_MagebloodPotion:
			character.AddStats(stats.Stats{
				stats.MP5: 12,
			})
		}
	}
}

///////////////////////////////////////////////////////////////////////////
//                             Seal of the Dawn Buff Consumes
///////////////////////////////////////////////////////////////////////////

func applySealOfTheDawnBuffConsumes(character *Character, consumes *proto.Consumes) {
	if consumes.SealOfTheDawn == proto.SealOfTheDawn_SealOfTheDawnUnknown {
		return
	}

	switch consumes.SealOfTheDawn {
	case proto.SealOfTheDawn_SealOfTheDawnDamageR1:
		sanctifiedDamageEffect(character, 1219539, 1.25)
	case proto.SealOfTheDawn_SealOfTheDawnDamageR2:
		sanctifiedDamageEffect(character, 1223348, 4.38)
	case proto.SealOfTheDawn_SealOfTheDawnDamageR3:
		sanctifiedDamageEffect(character, 1223349, 6.25)
	case proto.SealOfTheDawn_SealOfTheDawnDamageR4:
		sanctifiedDamageEffect(character, 1223350, 10.0)
	case proto.SealOfTheDawn_SealOfTheDawnDamageR5:
		sanctifiedDamageEffect(character, 1223351, 12.5)
	case proto.SealOfTheDawn_SealOfTheDawnDamageR6:
		sanctifiedDamageEffect(character, 1223352, 18.13)
	case proto.SealOfTheDawn_SealOfTheDawnDamageR7:
		sanctifiedDamageEffect(character, 1223353, 21.25)
	case proto.SealOfTheDawn_SealOfTheDawnDamageR8:
		sanctifiedDamageEffect(character, 1223354, 21.25)
	case proto.SealOfTheDawn_SealOfTheDawnDamageR9:
		sanctifiedDamageEffect(character, 1223355, 21.25)
	case proto.SealOfTheDawn_SealOfTheDawnDamageR10:
		sanctifiedDamageEffect(character, 1223357, 21.25)

	case proto.SealOfTheDawn_SealOfTheDawnTankR1:
		sanctifiedTankingEffect(character, 1220514, 3.13, 1.25)
	case proto.SealOfTheDawn_SealOfTheDawnTankR2:
		sanctifiedTankingEffect(character, 1223367, 4.38, 4.38)
	case proto.SealOfTheDawn_SealOfTheDawnTankR3:
		sanctifiedTankingEffect(character, 1223368, 5.0, 6.25)
	case proto.SealOfTheDawn_SealOfTheDawnTankR4:
		sanctifiedTankingEffect(character, 1223370, 5.63, 10)
	case proto.SealOfTheDawn_SealOfTheDawnTankR5:
		sanctifiedTankingEffect(character, 1223371, 6.25, 12.5)
	case proto.SealOfTheDawn_SealOfTheDawnTankR6:
		sanctifiedTankingEffect(character, 1223372, 7.5, 18.13)
	case proto.SealOfTheDawn_SealOfTheDawnTankR7:
		sanctifiedTankingEffect(character, 1223373, 8.13, 21.25)
	case proto.SealOfTheDawn_SealOfTheDawnTankR8:
		sanctifiedTankingEffect(character, 1223374, 8.13, 21.25)
	case proto.SealOfTheDawn_SealOfTheDawnTankR9:
		sanctifiedTankingEffect(character, 1223375, 8.13, 21.25)
	case proto.SealOfTheDawn_SealOfTheDawnTankR10:
		sanctifiedTankingEffect(character, 1223376, 8.13, 21.25)

	case proto.SealOfTheDawn_SealOfTheDawnHealingR1:
		sanctifiedHealingEffect(character, 1219548, 1.25, 1.25)
	case proto.SealOfTheDawn_SealOfTheDawnHealingR2:
		sanctifiedHealingEffect(character, 1223379, 3.13, 4.38)
	case proto.SealOfTheDawn_SealOfTheDawnHealingR3:
		sanctifiedHealingEffect(character, 1223380, 4.38, 6.25)
	case proto.SealOfTheDawn_SealOfTheDawnHealingR4:
		sanctifiedHealingEffect(character, 1223381, 7.5, 10.0)
	case proto.SealOfTheDawn_SealOfTheDawnHealingR5:
		sanctifiedHealingEffect(character, 1223382, 8.75, 12.5)
	case proto.SealOfTheDawn_SealOfTheDawnHealingR6:
		sanctifiedHealingEffect(character, 1223383, 13.13, 18.13)
	case proto.SealOfTheDawn_SealOfTheDawnHealingR7:
		sanctifiedHealingEffect(character, 1223384, 15.0, 21.25)
	case proto.SealOfTheDawn_SealOfTheDawnHealingR8:
		sanctifiedHealingEffect(character, 1223385, 15.0, 21.25)
	case proto.SealOfTheDawn_SealOfTheDawnHealingR9:
		sanctifiedHealingEffect(character, 1223386, 15.0, 21.25)
	case proto.SealOfTheDawn_SealOfTheDawnHealingR10:
		sanctifiedHealingEffect(character, 1223387, 15.0, 21.25)
	}
}

const MaxSanctifiedBonus = 8

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your damage by X% and your health by X% for each piece of Sanctified armor equipped.
func sanctifiedDamageEffect(character *Character, spellID int32, percentIncrease float64) {
	for _, unit := range getSanctifiedUnits(character) {
		sanctifiedBonus := int32(0)
		healthDeps := buildSanctifiedHealthDeps(unit, percentIncrease)

		unit.GetOrRegisterAura(Aura{
			Label:      "Seal of the Dawn (Damage)",
			ActionID:   ActionID{SpellID: spellID},
			BuildPhase: CharacterBuildPhaseGear,
			Duration:   NeverExpires,
			MaxStacks:  MaxSanctifiedBonus,
			OnInit: func(aura *Aura, sim *Simulation) {
				sanctifiedBonus = max(min(MaxSanctifiedBonus, character.PseudoStats.SanctifiedBonus), 0)
			},
			OnReset: func(aura *Aura, sim *Simulation) {
				aura.Activate(sim)
				aura.SetStacks(sim, sanctifiedBonus)
			},
			OnGain: func(aura *Aura, sim *Simulation) {
				aura.Unit.EnableBuildPhaseStatDep(sim, healthDeps[sanctifiedBonus])

				aura.Unit.PseudoStats.SanctifiedDamageMultiplier = 1.0 + percentIncrease/100.0*float64(sanctifiedBonus)
				aura.Unit.PseudoStats.DamageDealtMultiplier *= aura.Unit.PseudoStats.SanctifiedDamageMultiplier
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				aura.Unit.DisableBuildPhaseStatDep(sim, healthDeps[sanctifiedBonus])

				aura.Unit.PseudoStats.DamageDealtMultiplier /= aura.Unit.PseudoStats.SanctifiedDamageMultiplier
			},
		})
	}
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your healing by X% and your health by Y% for each piece of Sanctified armor equipped.
func sanctifiedHealingEffect(character *Character, spellID int32, healingPercentIncrease float64, healthPercentIncrease float64) {
	for _, unit := range getSanctifiedUnits(character) {
		sanctifiedBonus := int32(0)
		multiplier := 1.0
		healthDeps := buildSanctifiedHealthDeps(unit, healthPercentIncrease)

		unit.GetOrRegisterAura(Aura{
			Label:      "Seal of the Dawn (Healing)",
			ActionID:   ActionID{SpellID: spellID},
			BuildPhase: CharacterBuildPhaseGear,
			Duration:   NeverExpires,
			MaxStacks:  MaxSanctifiedBonus,
			OnInit: func(aura *Aura, sim *Simulation) {
				sanctifiedBonus = max(min(MaxSanctifiedBonus, character.PseudoStats.SanctifiedBonus), 0)
				multiplier = 1.0 + healingPercentIncrease/100.0*float64(sanctifiedBonus)
			},
			OnReset: func(aura *Aura, sim *Simulation) {
				aura.Activate(sim)
				aura.SetStacks(sim, sanctifiedBonus)
			},
			OnGain: func(aura *Aura, sim *Simulation) {
				aura.Unit.EnableBuildPhaseStatDep(sim, healthDeps[sanctifiedBonus])
				aura.Unit.PseudoStats.HealingDealtMultiplier *= multiplier
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				aura.Unit.DisableBuildPhaseStatDep(sim, healthDeps[sanctifiedBonus])
				aura.Unit.PseudoStats.HealingDealtMultiplier /= multiplier
			},
		})
	}
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your threat caused by X%, your damage by Y%, and your health by Y% for each piece of Sanctified armor equipped.
func sanctifiedTankingEffect(character *Character, spellID int32, threatPercentIncrease float64, damageHealthPercentIncrease float64) {
	for _, unit := range getSanctifiedUnits(character) {
		sanctifiedBonus := int32(0)
		threatMultiplier := 1.0
		healthDeps := buildSanctifiedHealthDeps(unit, damageHealthPercentIncrease)

		unit.GetOrRegisterAura(Aura{
			Label:      "Seal of the Dawn (Tanking)",
			ActionID:   ActionID{SpellID: spellID},
			BuildPhase: CharacterBuildPhaseGear,
			Duration:   NeverExpires,
			MaxStacks:  MaxSanctifiedBonus,
			OnInit: func(aura *Aura, sim *Simulation) {
				sanctifiedBonus = max(min(MaxSanctifiedBonus, character.PseudoStats.SanctifiedBonus), 0)
				threatMultiplier = 1.0 + threatPercentIncrease/100.0*float64(sanctifiedBonus)
			},
			OnReset: func(aura *Aura, sim *Simulation) {
				aura.Activate(sim)
				aura.SetStacks(sim, sanctifiedBonus)
			},
			OnGain: func(aura *Aura, sim *Simulation) {
				aura.Unit.EnableBuildPhaseStatDep(sim, healthDeps[sanctifiedBonus])

				aura.Unit.PseudoStats.SanctifiedDamageMultiplier = 1.0 + damageHealthPercentIncrease/100.0*float64(sanctifiedBonus)
				aura.Unit.PseudoStats.DamageDealtMultiplier *= aura.Unit.PseudoStats.SanctifiedDamageMultiplier
				aura.Unit.PseudoStats.ThreatMultiplier *= threatMultiplier
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				aura.Unit.DisableBuildPhaseStatDep(sim, healthDeps[sanctifiedBonus])

				aura.Unit.PseudoStats.DamageDealtMultiplier /= aura.Unit.PseudoStats.SanctifiedDamageMultiplier
				aura.Unit.PseudoStats.ThreatMultiplier /= threatMultiplier
			},
		})
	}
}

// Gets all units that the Sanctified buff should apply to.
// This includes the player and ALL pets/minions as of 2025-01-31
func getSanctifiedUnits(character *Character) []*Unit {
	units := []*Unit{&character.Unit}
	for _, pet := range character.Pets {
		units = append(units, &pet.Unit)
	}

	return units
}

func buildSanctifiedHealthDeps(unit *Unit, percentIncrease float64) []*stats.StatDependency {
	healthDeps := []*stats.StatDependency{}
	for i := 0; i < MaxSanctifiedBonus+1; i++ {
		healthDeps = append(healthDeps, unit.NewDynamicMultiplyStat(stats.Health, 1.0+percentIncrease/100.0*float64(i)))
	}

	return healthDeps
}

///////////////////////////////////////////////////////////////////////////
//                             Zanza-esque Consumes
///////////////////////////////////////////////////////////////////////////

func applyZanzaBuffConsumes(character *Character, consumes *proto.Consumes) {
	if consumes.ZanzaBuff == proto.ZanzaBuff_ZanzaBuffUnknown {
		return
	}

	switch consumes.ZanzaBuff {
	case proto.ZanzaBuff_SpiritOfZanza:
		character.AddStats(stats.Stats{
			stats.Stamina: 50,
			stats.Spirit:  50,
		})
	case proto.ZanzaBuff_ROIDS:
		character.AddStats(stats.Stats{
			stats.Strength: 25,
		})
	case proto.ZanzaBuff_GroundScorpokAssay:
		character.AddStats(stats.Stats{
			stats.Agility: 25,
		})
	case proto.ZanzaBuff_CerebralCortexCompound:
		character.AddStats(stats.Stats{
			stats.Intellect: 25,
		})
	case proto.ZanzaBuff_GizzardGum:
		character.AddStats(stats.Stats{
			stats.Spirit: 25,
		})
	case proto.ZanzaBuff_LungJuiceCocktail:
		character.AddStats(stats.Stats{
			stats.Stamina: 25,
		})
	case proto.ZanzaBuff_AtalaiMojoOfWar:
		if character.Level == 50 {
			character.AddStats(stats.Stats{
				stats.AttackPower:       48,
				stats.RangedAttackPower: 48,
			})
			ApplyAtalAiProc(character, consumes.ZanzaBuff)
		}
	case proto.ZanzaBuff_AtalaiMojoOfForbiddenMagic:
		if character.Level == 50 {
			character.AddStats(stats.Stats{
				stats.SpellPower: 40,
			})
			ApplyAtalAiProc(character, consumes.ZanzaBuff)
		}
	case proto.ZanzaBuff_AtalaiMojoOfLife:
		if character.Level == 50 {
			character.AddStats(stats.Stats{
				stats.HealingPower: 45,
				stats.MP5:          11,
			})
			ApplyAtalAiProc(character, consumes.ZanzaBuff)
		}
	}
}

func ApplyAtalAiProc(character *Character, atalaiBuff proto.ZanzaBuff) {
	icd := Cooldown{
		Timer:    character.NewTimer(),
		Duration: time.Second * 40,
	}

	switch atalaiBuff {
	case proto.ZanzaBuff_AtalaiMojoOfWar:
		procAuraStr := character.NewTemporaryStatsAura("Voodoo Frenzy Str Proc", ActionID{SpellID: 446335}, stats.Stats{stats.Strength: 35}, time.Second*10)
		procAuraAgi := character.NewTemporaryStatsAura("Voodoo Frenzy Agi Proc", ActionID{SpellID: 449409}, stats.Stats{stats.Agility: 35}, time.Second*10)
		procAuraStr.Icd = &icd
		procAuraAgi.Icd = &icd

		MakePermanent(character.RegisterAura(Aura{
			Label: "Voodoo Frenzy",
			OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(ProcMaskMeleeOrRanged) || !icd.IsReady(sim) {
					return
				}

				if sim.Proc(0.15, "Voodoo Frenzy") {
					icd.Use(sim)

					if aura.Unit.GetStat(stats.Strength) > aura.Unit.GetStat(stats.Agility) {
						procAuraStr.Activate(sim)
					} else {
						procAuraAgi.Activate(sim)
					}
				}
			},
		}))
	case proto.ZanzaBuff_AtalaiMojoOfForbiddenMagic:
		procSpell := character.RegisterSpell(SpellConfig{
			ActionID:    ActionID{SpellID: 446258},
			SpellSchool: SpellSchoolShadow,
			ProcMask:    ProcMaskEmpty,
			DefenseType: DefenseTypeMagic,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.56,

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				dmg := sim.Roll(204, 236)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicCrit) // TODO: Verify if it rolls miss? Most procs dont so we have it like this
			},
		})

		MakePermanent(character.RegisterAura(Aura{
			Label: "Forbidden Magic",
			OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(ProcMaskSpellDamage) || !icd.IsReady(sim) {
					return
				}

				if sim.Proc(0.25, "Forbidden Magic") {
					icd.Use(sim)
					procSpell.Cast(sim, character.CurrentTarget)
				}
			},
		}))
	case proto.ZanzaBuff_AtalaiMojoOfLife:
		// TODO: Your heals have a chance to restore 8 Energy, 1% Mana, or 4 Rage
		// This is also shared with the Darkmoon Card: Overgrowth trinket but unsure if they stack or not
	}
}

///////////////////////////////////////////////////////////////////////////
//                             Misc Consumes
///////////////////////////////////////////////////////////////////////////

func applyMiscConsumes(character *Character, miscConsumes *proto.MiscConsumes) {
	if miscConsumes == nil {
		return
	}

	if miscConsumes.BoglingRoot {
		character.PseudoStats.BonusPhysicalDamage += 1
	}

	if miscConsumes.ElixirOfCoalescedRegret {
		character.AddStats(stats.Stats{
			stats.Stamina:   1,
			stats.Agility:   1,
			stats.Strength:  1,
			stats.Intellect: 1,
			stats.Spirit:    1,
		})
	}

	if miscConsumes.RaptorPunch {
		character.AddStats(stats.Stats{
			stats.Intellect: 4,
			stats.Stamina:   -5,
		})
	}

	if miscConsumes.GreaterMarkOfTheDawn {
		character.AddStat(stats.Stamina, 30)
	}

	if miscConsumes.JujuEmber {
		character.AddStat(stats.FireResistance, 15)
	}

	if miscConsumes.JujuChill {
		character.AddStat(stats.FrostResistance, 15)
	}

	if miscConsumes.JujuFlurry {
		actionID := ActionID{SpellID: 16322}
		// In Vanilla Juju Flurry was bugged to act like Seal of the Crusader where it gave attack speed but also reduced damage done.
		jujuFlurryAura := character.RegisterAura(Aura{
			Label:    "Juju Flurry",
			ActionID: actionID,
			Duration: time.Second * 20,
			OnGain: func(aura *Aura, sim *Simulation) {
				aura.Unit.MultiplyMeleeSpeed(sim, 1.03)
				aura.Unit.AutoAttacks.MHAuto().ApplyMultiplicativeDamageBonus(1 / 1.03)
				aura.Unit.AutoAttacks.OHAuto().ApplyMultiplicativeDamageBonus(1 / 1.03)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				aura.Unit.MultiplyMeleeSpeed(sim, 1/1.03)
				aura.Unit.AutoAttacks.MHAuto().ApplyMultiplicativeDamageBonus(1.03)
				aura.Unit.AutoAttacks.OHAuto().ApplyMultiplicativeDamageBonus(1.03)
			},
		})
		jujuFlurrySpell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			ProcMask: ProcMaskEmpty,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute,
				},
				SharedCD: Cooldown{
					Timer:    character.GetAttackSpeedBuffCD(),
					Duration: time.Second * 10,
				},
			},
			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				jujuFlurryAura.Activate(sim)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell:    jujuFlurrySpell,
			Type:     CooldownTypeDPS,
			Priority: CooldownPriorityDefault,
		})
	}

	if miscConsumes.JujuEscape {
		actionID := ActionID{SpellID: 16321}
		jujuEscapeAura := character.RegisterAura(Aura{
			Label:    "Juju Escape",
			ActionID: actionID,
			Duration: time.Second * 10,
		}).AttachStatBuff(stats.Dodge, 5*DodgeRatingPerDodgeChance)

		jujuEscapeSpell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			ProcMask: ProcMaskEmpty,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute,
				},
			},
			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				jujuEscapeAura.Activate(sim)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell:    jujuEscapeSpell,
			Type:     CooldownTypeSurvival,
			Priority: CooldownPriorityDefault,
		})
	}
}

///////////////////////////////////////////////////////////////////////////
//                             Enchanting Consumes
///////////////////////////////////////////////////////////////////////////

func applyEnchantingConsumes(character *Character, consumes *proto.Consumes) {
	switch consumes.EnchantedSigil {
	case proto.EnchantedSigil_InnovationSigil:
		character.AddStats(stats.Stats{
			stats.AttackPower:       20,
			stats.RangedAttackPower: 20,
			stats.SpellPower:        20,
		})
	case proto.EnchantedSigil_LivingDreamsSigil:
		character.AddStats(stats.Stats{
			stats.AttackPower:       30,
			stats.RangedAttackPower: 30,
			stats.SpellPower:        30,
		})
	case proto.EnchantedSigil_FlowingWatersSigil:
		for _, player := range character.Env.Raid.AllPlayerUnits {
			player.AddStats(stats.Stats{
				stats.AttackPower:       30,
				stats.RangedAttackPower: 30,
				stats.SpellPower:        30,
			})
		}
	case proto.EnchantedSigil_WrathOfTheStormSigil:
		for _, player := range character.Env.Raid.AllPlayerUnits {
			player.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
				stats.SpellPower:        40,
			})
		}
	}
}

///////////////////////////////////////////////////////////////////////////
//                             Engineering Explosives
///////////////////////////////////////////////////////////////////////////

var SapperActionID = ActionID{ItemID: 10646}
var FumigatorActionID = ActionID{ItemID: 233985}
var SolidDynamiteActionID = ActionID{ItemID: 10507}
var DenseDynamiteActionID = ActionID{ItemID: 18641}
var ThoriumGrenadeActionID = ActionID{ItemID: 15993}
var EzThroRadiationBombActionID = ActionID{ItemID: 215168}
var HighYieldRadiationBombActionID = ActionID{ItemID: 215127}
var GoblinLandMineActionID = ActionID{ItemID: 4395}
var ObsidianBombActionID = ActionID{ItemID: 233986}
var StratholmeHolyWaterActionID = ActionID{ItemID: 13180}

func registerExplosivesCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	hasFiller := consumes.FillerExplosive != proto.Explosive_ExplosiveUnknown
	hasSapper := consumes.SapperExplosive != proto.SapperExplosive_SapperUnknown

	if !hasSapper && !hasFiller {
		return
	}
	sharedTimer := character.NewTimer()

	if hasSapper {
		nonEngiSappers := []proto.SapperExplosive{proto.SapperExplosive_SapperFumigator}
		if !character.HasProfession(proto.Profession_Engineering) && !slices.Contains(nonEngiSappers, consumes.SapperExplosive) {
			return
		}

		var sapperSpell *Spell
		switch consumes.SapperExplosive {
		case proto.SapperExplosive_SapperGoblinSapper:
			sapperSpell = character.newSapperSpell(sharedTimer)
		case proto.SapperExplosive_SapperFumigator:
			sapperSpell = character.newFumigatorSpell(sharedTimer)
		}

		sapperSpell.ExtraCastCondition = func(sim *Simulation, target *Unit) bool {
			return character.DistanceFromTarget <= 10
		}

		character.AddMajorCooldown(MajorCooldown{
			Spell:    sapperSpell,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 30,
			ShouldActivate: func(s *Simulation, c *Character) bool {
				return !character.IsShapeshifted()
			},
		})
	}

	if hasFiller {
		// Update this list with explosives that don't require engi
		nonEngiExplosives := []proto.Explosive{proto.Explosive_ExplosiveEzThroRadiationBomb, proto.Explosive_ExplosiveObsidianBomb, proto.Explosive_ExplosiveStratholmeHolyWater}
		if !character.HasProfession(proto.Profession_Engineering) && !slices.Contains(nonEngiExplosives, consumes.FillerExplosive) {
			return
		}

		var filler *Spell
		switch consumes.FillerExplosive {
		case proto.Explosive_ExplosiveStratholmeHolyWater:
			filler = character.newStratholmeHolyWaterSpell(sharedTimer)
		case proto.Explosive_ExplosiveObsidianBomb:
			filler = character.newObisidianBombSpell(sharedTimer)
		case proto.Explosive_ExplosiveSolidDynamite:
			filler = character.newSolidDynamiteSpell(sharedTimer)
		case proto.Explosive_ExplosiveDenseDynamite:
			filler = character.newDenseDynamiteSpell(sharedTimer)
		case proto.Explosive_ExplosiveThoriumGrenade:
			filler = character.newThoriumGrenadeSpell(sharedTimer)
		case proto.Explosive_ExplosiveEzThroRadiationBomb:
			filler = character.newEzThroRadiationBombSpell(sharedTimer)
		case proto.Explosive_ExplosiveHighYieldRadiationBomb:
			filler = character.newHighYieldRadiationBombSpell(sharedTimer)
		case proto.Explosive_ExplosiveGoblinLandMine:
			filler = character.newGoblinLandMineSpell(sharedTimer)
		}

		character.AddMajorCooldown(MajorCooldown{
			Spell:    filler,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 10,
			ShouldActivate: func(s *Simulation, c *Character) bool {
				return !character.IsShapeshifted()
			},
		})
	}
}

type ExplosiveConfig struct {
	ActionID     ActionID
	SpellSchool  SpellSchool
	DefenseType  DefenseType
	MissileSpeed float64

	SharedTimer *Timer
	Cooldown    time.Duration
	CastTime    time.Duration

	// Land Mines have a 10s "arming time" before they hit
	TriggerDelay time.Duration

	OnHitAction OnHitAction

	MinDamage     float64
	MaxDamage     float64
	SelfMinDamage float64
	SelfMaxDamage float64

	BonusCoefficient float64
}

type OnHitAction func(sim *Simulation, spell *Spell, result *SpellResult)

func (character *Character) applyExplosiveDamage(sim *Simulation, spell *Spell, explosiveConfig ExplosiveConfig) {
	for _, aoeTarget := range sim.Encounter.TargetUnits {
		result := spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(explosiveConfig.MinDamage, explosiveConfig.MaxDamage), spell.OutcomeMagicHitAndCrit)

		if explosiveConfig.OnHitAction != nil {
			explosiveConfig.OnHitAction(sim, spell, result)
		}
	}

	if explosiveConfig.SelfMinDamage > 0 && explosiveConfig.SelfMaxDamage > 0 {
		spell.CalcAndDealDamage(sim, &character.Unit, sim.Roll(explosiveConfig.SelfMinDamage, explosiveConfig.SelfMaxDamage), spell.OutcomeMagicHitAndCrit)
	}
}

func (character *Character) castExplosive(sim *Simulation, spell *Spell, explosiveConfig ExplosiveConfig) {
	if explosiveConfig.MissileSpeed > 0 {
		spell.WaitTravelTime(sim, func(sim *Simulation) {
			character.applyExplosiveDamage(sim, spell, explosiveConfig)
		})
	} else {
		character.applyExplosiveDamage(sim, spell, explosiveConfig)
	}
}

// Creates a spell object for the common explosive case.
func (character *Character) newBasicExplosiveSpellConfig(explosiveConfig ExplosiveConfig) SpellConfig {
	var defaultCast Cast
	if explosiveConfig.CastTime > 0 {
		defaultCast = Cast{
			CastTime: explosiveConfig.CastTime,
		}
	}

	cooldownConfig := Cooldown{}
	if explosiveConfig.Cooldown > 0 {
		cooldownConfig = Cooldown{
			Timer:    character.NewTimer(),
			Duration: explosiveConfig.Cooldown,
		}
	}

	flags := SpellFlagCastTimeNoGCD

	if explosiveConfig.DefenseType == DefenseTypeNone {
		explosiveConfig.DefenseType = DefenseTypeMagic
	}

	return SpellConfig{
		ActionID:     explosiveConfig.ActionID,
		SpellSchool:  explosiveConfig.SpellSchool,
		DefenseType:  explosiveConfig.DefenseType,
		ProcMask:     ProcMaskEmpty,
		Flags:        flags,
		MissileSpeed: explosiveConfig.MissileSpeed,

		Cast: CastConfig{
			DefaultCast: defaultCast,
			CD:          cooldownConfig,
			IgnoreHaste: true,
			SharedCD: Cooldown{
				Timer:    explosiveConfig.SharedTimer,
				Duration: time.Minute,
			},
			ModifyCast: func(sim *Simulation, _ *Spell, _ *Cast) {
				character.CancelShapeshift(sim)
			},
		},

		// Explosives always have 1% resist chance, so just give them hit cap.
		BonusHitRating: 100 * SpellHitRatingPerHitChance,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: explosiveConfig.BonusCoefficient,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			if explosiveConfig.TriggerDelay > 0 {
				StartDelayedAction(sim, DelayedActionOptions{
					DoAt: sim.CurrentTime + explosiveConfig.TriggerDelay,
					OnAction: func(sim *Simulation) {
						character.castExplosive(sim, spell, explosiveConfig)
					},
				})
				return
			}

			character.castExplosive(sim, spell, explosiveConfig)
		},
	}
}

// Needs testing for Silithid interaction if in raid
func (character *Character) newFumigatorSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(ExplosiveConfig{
		ActionID:      FumigatorActionID,
		SpellSchool:   SpellSchoolFire,
		SharedTimer:   sharedTimer,
		Cooldown:      time.Minute * 5,
		MinDamage:     650,
		MaxDamage:     950,
		SelfMinDamage: 475,
		SelfMaxDamage: 725,
	}))
}
func (character *Character) newSapperSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(ExplosiveConfig{
		ActionID:      SapperActionID,
		SpellSchool:   SpellSchoolFire,
		SharedTimer:   sharedTimer,
		Cooldown:      time.Minute * 5,
		MinDamage:     450,
		MaxDamage:     750,
		SelfMinDamage: 375,
		SelfMaxDamage: 625,
	}))
}

func (character *Character) newStratholmeHolyWaterSpell(sharedTimer *Timer) *Spell {
	explosiveConfig := ExplosiveConfig{
		ActionID:         StratholmeHolyWaterActionID,
		SpellSchool:      SpellSchoolHoly,
		DefenseType:      DefenseTypeMelee,
		SharedTimer:      sharedTimer,
		MinDamage:        438,
		MaxDamage:        562,
		BonusCoefficient: 1,
	}
	config := character.newBasicExplosiveSpellConfig(explosiveConfig)
	var outcomeMagicHitAndBaseSpellCrit OutcomeApplier
	config.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
		for _, aoeTarget := range sim.Encounter.TargetUnits {
			if aoeTarget.MobType != proto.MobType_MobTypeUndead {
				continue
			}
			spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(explosiveConfig.MinDamage, explosiveConfig.MaxDamage), outcomeMagicHitAndBaseSpellCrit)
		}
	}
	holyWaterSpell := character.GetOrRegisterSpell(config)
	baseSpellCrit := character.GetBaseStats()[stats.SpellCrit] / 100.0
	outcomeMagicHitAndBaseSpellCrit = func(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
		isPartialResist := result.DidResist()
		if holyWaterSpell.MagicHitCheck(sim, attackTable) {
			if sim.RandomFloat("Magical Crit Roll") < baseSpellCrit {
				result.Outcome = OutcomeCrit
				result.Damage *= holyWaterSpell.CritMultiplier(attackTable)
				holyWaterSpell.SpellMetrics[result.Target.UnitIndex].Crits++
				if isPartialResist {
					holyWaterSpell.SpellMetrics[result.Target.UnitIndex].ResistedCrits++
				}
			} else {
				result.Outcome = OutcomeHit
				holyWaterSpell.SpellMetrics[result.Target.UnitIndex].Hits++
				if isPartialResist {
					holyWaterSpell.SpellMetrics[result.Target.UnitIndex].ResistedHits++
				}
			}
		} else {
			result.Outcome = OutcomeMiss
			result.Damage = 0
			holyWaterSpell.SpellMetrics[result.Target.UnitIndex].Misses++
		}
	}
	return holyWaterSpell
}

func (character *Character) newObisidianBombSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(ExplosiveConfig{
		ActionID:     ObsidianBombActionID,
		SpellSchool:  SpellSchoolFire,
		MissileSpeed: 14,
		SharedTimer:  sharedTimer,
		CastTime:     time.Second,
		MinDamage:    530,
		MaxDamage:    670,
	}))
}
func (character *Character) newSolidDynamiteSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(ExplosiveConfig{
		ActionID:     SolidDynamiteActionID,
		SpellSchool:  SpellSchoolFire,
		MissileSpeed: 14,
		SharedTimer:  sharedTimer,
		CastTime:     time.Second,
		MinDamage:    213,
		MaxDamage:    287,
	}))
}
func (character *Character) newDenseDynamiteSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(ExplosiveConfig{
		ActionID:     DenseDynamiteActionID,
		SpellSchool:  SpellSchoolFire,
		MissileSpeed: 14,
		SharedTimer:  sharedTimer,
		CastTime:     time.Second,
		MinDamage:    340,
		MaxDamage:    460,
	}))
}
func (character *Character) newThoriumGrenadeSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(ExplosiveConfig{
		ActionID:     ThoriumGrenadeActionID,
		SpellSchool:  SpellSchoolFire,
		MissileSpeed: 25,
		SharedTimer:  sharedTimer,
		CastTime:     time.Second,
		MinDamage:    300,
		MaxDamage:    500,
	}))
}
func (character *Character) newGoblinLandMineSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(ExplosiveConfig{
		ActionID:     GoblinLandMineActionID,
		SpellSchool:  SpellSchoolFire,
		SharedTimer:  sharedTimer,
		TriggerDelay: time.Second * 10,
		MinDamage:    394,
		MaxDamage:    506,
	}))
}

// Creates a spell object for the common explosive case.
func (character *Character) newRadiationBombSpellConfig(dotDamage float64, explosiveConfig ExplosiveConfig) SpellConfig {
	explosiveConfig.OnHitAction = func(sim *Simulation, spell *Spell, result *SpellResult) {
		if result.Landed() {
			spell.Dot(result.Target).Apply(sim)
		}
	}

	return SpellConfig{
		ActionID:     explosiveConfig.ActionID,
		SpellSchool:  explosiveConfig.SpellSchool,
		DefenseType:  DefenseTypeMagic,
		ProcMask:     ProcMaskEmpty,
		Flags:        SpellFlagCastTimeNoGCD,
		MissileSpeed: explosiveConfig.MissileSpeed,

		Cast: CastConfig{
			DefaultCast: Cast{
				CastTime: explosiveConfig.CastTime,
			},
			IgnoreHaste: true,
			CD:          Cooldown{},
			SharedCD: Cooldown{
				Timer:    explosiveConfig.SharedTimer,
				Duration: time.Minute,
			},
		},

		// Explosives always have 1% resist chance, so just give them hit cap.
		BonusHitRating: 100 * SpellHitRatingPerHitChance,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		// TODO: This should use another spell (443813) as the DoT
		// Also doesn't apply to bosses or something.
		Dot: DotConfig{
			Aura: Aura{
				Label: explosiveConfig.ActionID.String(),
			},

			NumberOfTicks: 5,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *Simulation, target *Unit, dot *Dot, isRollover bool) {
				// Use nature school for dot modifiers
				dot.Spell.SpellSchool = SpellSchoolNature
				dot.Spell.SchoolIndex = stats.SchoolIndexNature

				dot.Snapshot(target, dotDamage, isRollover)

				// Revert to fire school
				dot.Spell.SpellSchool = SpellSchoolFire
				dot.Spell.SchoolIndex = stats.SchoolIndexFire
			},
			OnTick: func(sim *Simulation, target *Unit, dot *Dot) {
				// Use nature school for dot ticks
				dot.Spell.SpellSchool = SpellSchoolNature
				dot.Spell.SchoolIndex = stats.SchoolIndexNature

				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)

				// Revert to fire school
				dot.Spell.SpellSchool = SpellSchoolFire
				dot.Spell.SchoolIndex = stats.SchoolIndexFire
			},
		},

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.WaitTravelTime(sim, func(simulation *Simulation) {
				character.castExplosive(sim, spell, explosiveConfig)
			})
		},
	}
}
func (character *Character) newEzThroRadiationBombSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newRadiationBombSpellConfig(10, ExplosiveConfig{
		ActionID:     EzThroRadiationBombActionID,
		SpellSchool:  SpellSchoolFire,
		MissileSpeed: 14,
		SharedTimer:  sharedTimer,
		CastTime:     time.Millisecond * 1500,
		MinDamage:    112,
		MaxDamage:    188,
	}))
}
func (character *Character) newHighYieldRadiationBombSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newRadiationBombSpellConfig(25, ExplosiveConfig{
		ActionID:     HighYieldRadiationBombActionID,
		SpellSchool:  SpellSchoolFire,
		MissileSpeed: 25,
		SharedTimer:  sharedTimer,
		CastTime:     time.Second,
		MinDamage:    150,
		MaxDamage:    250,
	}))
}

///////////////////////////////////////////////////////////////////////////
//                             Potions
///////////////////////////////////////////////////////////////////////////

func registerPotionCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	defaultPotion := consumes.DefaultPotion

	potionCD := character.NewTimer()

	if defaultPotion == proto.Potions_UnknownPotion {
		return
	}

	defaultMCD := makePotionActivation(defaultPotion, character, potionCD)

	if defaultMCD.Spell != nil {
		character.AddMajorCooldown(defaultMCD)
	}
}

func makePotionActivation(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
	mcd := makePotionActivationInternal(potionType, character, potionCD)
	if mcd.Spell != nil {
		// Mark as 'Encounter Only' so that users are forced to select the generic Potion
		// placeholder action instead of specific potion spells, in APL prepull. This
		// prevents a mismatch between Consumes and Rotation settings.
		mcd.Spell.Flags |= SpellFlagEncounterOnly | SpellFlagPotion | SpellFlagCastTimeNoGCD
		oldApplyEffects := mcd.Spell.ApplyEffects
		mcd.Spell.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			oldApplyEffects(sim, target, spell)
			if sim.CurrentTime < 0 {
				potionCD.Set(sim.CurrentTime + 2*time.Minute)
				character.UpdateMajorCooldowns()
			}
		}
	}
	return mcd
}

func makeHealthConsumableMCD(itemId int32, character *Character, cdTimer *Timer) MajorCooldown {
	// Using min values for healthstones as locks generally don't spec into improved
	minRoll := map[int32]float64{
		858:   140,
		929:   280,
		1710:  455,
		3928:  700,
		5509:  500,
		5510:  800,
		9421:  1200,
		13446: 1050,
	}[itemId]

	maxRoll := map[int32]float64{
		858:   180,
		929:   360,
		1710:  585,
		3928:  900,
		5509:  500,
		5510:  800,
		9421:  1200,
		13446: 1750,
	}[itemId]

	cdDuration := time.Minute * 2

	actionID := ActionID{ItemID: itemId}
	healthMetrics := character.NewHealthMetrics(actionID)

	return MajorCooldown{
		Type: CooldownTypeSurvival,
		ShouldActivate: func(sim *Simulation, character *Character) bool {
			// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
			return (character.MaxHealth()-(character.CurrentHealth()) >= maxRoll) && !character.IsShapeshifted()
		},
		Spell: character.GetOrRegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    cdTimer,
					Duration: cdDuration,
				},
				ModifyCast: func(sim *Simulation, _ *Spell, _ *Cast) {
					character.CancelShapeshift(sim)
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				healthGain := sim.RollWithLabel(minRoll, maxRoll, "Health Consumable")
				character.GainHealth(sim, healthGain, healthMetrics)
			},
		}),
	}
}

func makeManaConsumableMCD(itemId int32, character *Character, cdTimer *Timer) MajorCooldown {
	minRoll := map[int32]float64{
		3827:  455.0,
		6149:  700.0,
		4381:  150.0,
		12662: 900.0,
		13443: 900.0,
		13444: 1350.0,
	}[itemId]

	maxRoll := map[int32]float64{
		3827:  585.0,
		6149:  900.0,
		4381:  250.0,
		12662: 1500.0,
		13443: 1500.0,
		13444: 2250.0,
	}[itemId]

	cdDuration := time.Minute * 2

	actionID := ActionID{ItemID: itemId}
	manaMetrics := character.NewManaMetrics(actionID)

	return MajorCooldown{
		Type: CooldownTypeMana,
		ShouldActivate: func(sim *Simulation, character *Character) bool {
			// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
			totalRegen := character.ManaRegenPerSecondWhileCasting() * 2
			return (character.MaxMana()-(character.CurrentMana()+totalRegen) >= maxRoll) && !character.IsShapeshifted()
		},
		Spell: character.GetOrRegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    cdTimer,
					Duration: cdDuration,
				},
				ModifyCast: func(sim *Simulation, _ *Spell, _ *Cast) {
					character.CancelShapeshift(sim)
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				manaGain := sim.RollWithLabel(minRoll, maxRoll, "Mana Consumable")
				character.AddMana(sim, manaGain, manaMetrics)
			},
		}),
	}
}

func makeArmorConsumableMCD(itemId int32, character *Character, cdTimer *Timer) MajorCooldown {
	actionID := ActionID{ItemID: itemId}
	cdDuration := time.Minute * 2
	lesserStoneshieldAura := character.NewTemporaryStatsAura("Lesser Stoneshield Potion", actionID, stats.Stats{stats.BonusArmor: 1000}, time.Second*90)
	greaterStoneshieldAura := character.NewTemporaryStatsAura("Greater Stoneshield Potion", actionID, stats.Stats{stats.BonusArmor: 2000}, time.Second*120)

	return MajorCooldown{
		Type: CooldownTypeSurvival,
		Spell: character.GetOrRegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    cdTimer,
					Duration: cdDuration,
				},
				ModifyCast: func(sim *Simulation, _ *Spell, _ *Cast) {
					character.CancelShapeshift(sim)
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				switch itemId {
				case 4623:
					lesserStoneshieldAura.Activate(sim)
				case 13455:

					greaterStoneshieldAura.Activate(sim)
				}
			},
		}),
	}
}

func makeMagicResistancePotionMCD(character *Character, cdTimer *Timer) MajorCooldown {
	actionID := ActionID{ItemID: 4623}
	cdDuration := time.Minute * 2

	stats := stats.Stats{
		stats.ArcaneResistance: 50,
		stats.FireResistance:   50,
		stats.FrostResistance:  50,
		stats.NatureResistance: 50,
		stats.ShadowResistance: 50,
	}

	// Since many people will keep this rolling as a substitute for capping Fire Resistance, show the stats as a baseline
	aura := character.NewTemporaryStatsAura("Magic Resistance Potion", actionID, stats, time.Minute*3)
	aura.BuildPhase = CharacterBuildPhaseConsumes

	return MajorCooldown{
		Type: CooldownTypeSurvival,
		Spell: character.GetOrRegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    cdTimer,
					Duration: cdDuration,
				},
				ModifyCast: func(sim *Simulation, _ *Spell, _ *Cast) {
					character.CancelShapeshift(sim)
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				aura.Activate(sim)
			},
		}),
	}
}

func makeRageConsumableMCD(itemId int32, character *Character, cdTimer *Timer) MajorCooldown {
	minRoll := map[int32]float64{
		5631:  20.0,
		5633:  30.0,
		13442: 45.0,
	}[itemId]

	maxRoll := map[int32]float64{
		5631:  40.0,
		5633:  60.0,
		13442: 75.0,
	}[itemId]

	cdDuration := time.Minute * 2

	actionID := ActionID{ItemID: itemId}
	rageMetrics := character.NewRageMetrics(actionID)
	aura := character.NewTemporaryStatsAura("Mighty Rage Potion", actionID, stats.Stats{stats.Strength: 60}, time.Second*20)
	return MajorCooldown{
		Type: CooldownTypeDPS,
		ShouldActivate: func(sim *Simulation, character *Character) bool {
			return !character.IsShapeshifted()
		},
		Spell: character.GetOrRegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    cdTimer,
					Duration: cdDuration,
				},
				ModifyCast: func(sim *Simulation, _ *Spell, _ *Cast) {
					character.CancelShapeshift(sim)
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				rageGain := sim.RollWithLabel(minRoll, maxRoll, "Rage Consumable")
				character.AddRage(sim, rageGain, rageMetrics)
				if itemId == 13442 {
					aura.Activate(sim)
				}
			},
		}),
	}
}

func makePotionActivationInternal(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
	if potionType == proto.Potions_UnknownPotion {
		return MajorCooldown{}
	}

	switch potionType {
	case proto.Potions_LesserHealingPotion:
		return makeHealthConsumableMCD(858, character, potionCD)
	case proto.Potions_HealingPotion:
		return makeHealthConsumableMCD(929, character, potionCD)
	case proto.Potions_GreaterHealingPotion:
		return makeHealthConsumableMCD(1710, character, potionCD)
	case proto.Potions_SuperiorHealingPotion:
		return makeHealthConsumableMCD(3928, character, potionCD)
	case proto.Potions_MajorHealingPotion:
		return makeHealthConsumableMCD(13446, character, potionCD)

	case proto.Potions_LesserManaPotion:
		return makeManaConsumableMCD(3385, character, potionCD)
	case proto.Potions_ManaPotion:
		return makeManaConsumableMCD(3827, character, potionCD)
	case proto.Potions_GreaterManaPotion:
		return makeManaConsumableMCD(6149, character, potionCD)
	case proto.Potions_SuperiorManaPotion:
		return makeManaConsumableMCD(13443, character, potionCD)
	case proto.Potions_MajorManaPotion:
		return makeManaConsumableMCD(13444, character, potionCD)

	case proto.Potions_RagePotion:
		return makeRageConsumableMCD(5631, character, potionCD)
	case proto.Potions_GreatRagePotion:
		return makeRageConsumableMCD(5633, character, potionCD)
	case proto.Potions_MightyRagePotion:
		return makeRageConsumableMCD(13442, character, potionCD)

	case proto.Potions_LesserStoneshieldPotion:
		return makeArmorConsumableMCD(4623, character, potionCD)
	case proto.Potions_GreaterStoneshieldPotion:
		return makeArmorConsumableMCD(13455, character, potionCD)

	case proto.Potions_MagicResistancePotion:
		return makeMagicResistancePotionMCD(character, potionCD)
	// case proto.Potions_GreaterArcaneProtectionPotion:
	// 	return makeSchoolProtectionConsumableMCD(13461, character, potionCD)
	// case proto.Potions_GreaterFireProtectionPotion:
	// 	return makeSchoolProtectionConsumableMCD(13457, character, potionCD)
	// case proto.Potions_GreaterFrostProtectionPotion:
	// 	return makeSchoolProtectionConsumableMCD(13456, character, potionCD)
	// case proto.Potions_GreaterFrostProtectionPotion:
	// 	return makeSchoolProtectionConsumableMCD(13460, character, potionCD)
	// case proto.Potions_GreaterFrostProtectionPotion:
	// 	return makeSchoolProtectionConsumableMCD(13458, character, potionCD)
	// case proto.Potions_GreaterFrostProtectionPotion:
	// 	return makeSchoolProtectionConsumableMCD(13459, character, potionCD)
	default:
		return MajorCooldown{}
	}
}

func registerConjuredCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	conjuredType := consumes.DefaultConjured

	if conjuredType == proto.Conjured_ConjuredUnknown {
		return
	}

	timer := character.GetConjuredCD()

	var mcd MajorCooldown
	switch conjuredType {
	case proto.Conjured_ConjuredHealthstone:
		mcd = makeHealthConsumableMCD(5509, character, timer)
	case proto.Conjured_ConjuredGreaterHealthstone:
		mcd = makeHealthConsumableMCD(5510, character, timer)
	case proto.Conjured_ConjuredMajorHealthstone:
		mcd = makeHealthConsumableMCD(9421, character, timer)
	case proto.Conjured_ConjuredDemonicRune:
		mcd = makeManaConsumableMCD(12662, character, timer)
	case proto.Conjured_ConjuredMinorRecombobulator:
		mcd = makeManaConsumableMCD(4381, character, timer)
	// Handled in the rogue package
	// case proto.Conjured_ConjuredRogueThistleTea:
	default:
		return
	}

	character.AddMajorCooldown(mcd)
}

func registerMildlyIrradiatedRejuvCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()

	if consumes.MildlyIrradiatedRejuvPot {
		actionID := ActionID{ItemID: 215162}
		healthMetrics := character.NewHealthMetrics(actionID)
		manaMetrics := character.NewManaMetrics(actionID)
		aura := character.RegisterAura(Aura{
			ActionID: actionID,
			Label:    "Mildly Irradiated Rejuvenation Potion",
			Duration: time.Second * 20,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.AddStatsDynamic(sim, stats.Stats{
					stats.AttackPower:       40,
					stats.RangedAttackPower: 40,
					stats.SpellDamage:       35,
				})
				character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 2
				character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= 2
				character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 2
				character.AddMoveSpeedModifier(&aura.ActionID, .30)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.AddStatsDynamic(sim, stats.Stats{
					stats.AttackPower:       -40,
					stats.RangedAttackPower: -40,
					stats.SpellDamage:       -35,
				})
				character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= 2
				character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= 2
				character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= 2
				character.RemoveMoveSpeedModifier(&aura.ActionID)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Type: CooldownTypeDPS,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return !character.IsShapeshifted()
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
					ModifyCast: func(sim *Simulation, _ *Spell, _ *Cast) {
						character.CancelShapeshift(sim)
					},
				},
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					healthGain := sim.RollWithLabel(340, 460, "Mildly Irradiated Rejuvenation Potion")
					manaGain := sim.RollWithLabel(262, 438, "Mildly Irradiated Rejuvenation Potion")

					character.GainHealth(sim, healthGain*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
					character.AddMana(sim, manaGain, manaMetrics)

					aura.Activate(sim)
				},
			}),
		})
	}
}
