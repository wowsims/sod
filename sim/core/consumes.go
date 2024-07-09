package core

import (
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
	applyZanzaBuffConsumes(character, consumes)
	applyMiscConsumes(character, consumes.MiscConsumes)
	applyEnchantingConsumes(character, consumes)

	registerPotionCD(agent, consumes)
	registerConjuredCD(agent, consumes)
	registerMildlyIrradiatedRejuvCD(agent, consumes)
	registerExplosivesCD(agent, consumes)
}

func ApplyPetConsumeEffects(pet *Character, ownerConsumes *proto.Consumes) {
	pet.AddStat(stats.Agility, []float64{0, 5, 9, 13, 17}[ownerConsumes.PetScrollOfAgility])
	pet.AddStat(stats.Strength, []float64{0, 5, 9, 13, 17}[ownerConsumes.PetScrollOfStrength])
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
	case proto.Flask_FlaskOfSupremePower:
		character.AddStats(stats.Stats{
			stats.SpellPower: 150,
		})
	case proto.Flask_FlaskOfTheTitans:
		character.AddStats(stats.Stats{
			stats.Health: 1200,
		})
	case proto.Flask_FlaskOfChromaticResistance:
		character.AddStats(stats.Stats{
			stats.ArcaneResistance: 25,
			stats.FireResistance:   25,
			stats.FrostResistance:  25,
			stats.NatureResistance: 25,
			stats.ShadowResistance: 25,
		})
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
	if character.HasOHWeapon() {
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
		case proto.WeaponImbue_BrillianWizardOil:
			character.AddStats(stats.Stats{
				stats.SpellPower: 36,
				stats.SpellCrit:  1 * SpellCritRatingPerCritChance,
			})

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
		case proto.WeaponImbue_BlackfathomSharpeningStone:
			character.AddStats(stats.Stats{
				stats.MeleeHit: 2 * MeleeHitRatingPerHitChance,
			})

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
		ProcMask:    ProcMaskSpellDamage,

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
		ProcMask:    ProcMaskSpellDamage,

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

	procSpell := character.RegisterSpell(SpellConfig{
		ActionID:    ActionID{SpellID: 15851},
		SpellSchool: SpellSchoolFire,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagNone,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
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

			if sim.RandomFloat("Dragonbreath Chili") < procChance {
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
				stats.AttackPower: 40,
			})
		case proto.AttackPowerBuff_WinterfallFirewater:
			character.AddStats(stats.Stats{
				stats.AttackPower: 35,
			})
		}
	}

	if consumes.AgilityElixir != proto.AgilityElixir_AgilityElixirUnknown {
		switch consumes.AgilityElixir {
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
		character.PseudoStats.BonusDamage += 1
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
}

///////////////////////////////////////////////////////////////////////////
//                             Enchanting Consumes
///////////////////////////////////////////////////////////////////////////

func applyEnchantingConsumes(character *Character, consumes *proto.Consumes) {
	if !character.HasProfession(proto.Profession_Enchanting) || consumes.EnchantedSigil == proto.EnchantedSigil_UnknownSigil {
		return
	}

	switch consumes.EnchantedSigil {
	case proto.EnchantedSigil_InnovationSigil:
		character.AddStats(stats.Stats{
			stats.AttackPower:       20,
			stats.RangedAttackPower: 20,
			stats.SpellPower:        20,
		})
	case proto.EnchantedSigil_LivingDreamsSigil:
		character.AddStats(stats.Stats{
			stats.AttackPower:       50,
			stats.RangedAttackPower: 50,
			stats.SpellPower:        50,
		})
	}
}

///////////////////////////////////////////////////////////////////////////
//                             Engineering Explosives
///////////////////////////////////////////////////////////////////////////

var SapperActionID = ActionID{ItemID: 10646}
var SolidDynamiteActionID = ActionID{ItemID: 10507}
var DenseDynamiteActionID = ActionID{ItemID: 18641}
var ThoriumGrenadeActionID = ActionID{ItemID: 15993}
var EzThroRadiationBombActionID = ActionID{ItemID: 215168}
var HighYieldRadiationBombActionID = ActionID{ItemID: 215127}
var GoblinLandMineActionID = ActionID{ItemID: 4395}

func registerExplosivesCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	hasFiller := consumes.FillerExplosive != proto.Explosive_ExplosiveUnknown

	if !consumes.Sapper && !hasFiller {
		return
	}
	sharedTimer := character.NewTimer()

	if consumes.Sapper && character.HasProfession(proto.Profession_Engineering) {
		character.AddMajorCooldown(MajorCooldown{
			Spell:    character.newSapperSpell(sharedTimer),
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 30,
			ShouldActivate: func(s *Simulation, c *Character) bool {
				return !character.IsShapeshifted()
			},
		})
	}

	if hasFiller {
		if consumes.FillerExplosive != proto.Explosive_ExplosiveEzThroRadiationBomb && !character.HasProfession(proto.Profession_Engineering) {
			return
		}

		var filler *Spell
		switch consumes.FillerExplosive {
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

// Creates a spell object for the common explosive case.
// TODO: create 10s delay on Goblin Landmine cast to damage
func (character *Character) newBasicExplosiveSpellConfig(sharedTimer *Timer, actionID ActionID, school SpellSchool, minDamage float64, maxDamage float64, cooldown Cooldown, selfMinDamage float64, selfMaxDamage float64) SpellConfig {
	isSapper := actionID.SameAction(SapperActionID)

	var defaultCast Cast
	if !isSapper {
		defaultCast = Cast{
			CastTime: time.Second,
		}
	}

	return SpellConfig{
		ActionID:    actionID,
		SpellSchool: school,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagCastTimeNoGCD,

		Cast: CastConfig{
			DefaultCast: defaultCast,
			CD:          cooldown,
			IgnoreHaste: true,
			SharedCD: Cooldown{
				Timer:    sharedTimer,
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

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(minDamage, maxDamage) * sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			if isSapper {
				baseDamage := sim.Roll(selfMinDamage, selfMaxDamage)
				spell.CalcAndDealDamage(sim, &character.Unit, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	}
}
func (character *Character) newSapperSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, SapperActionID, SpellSchoolFire, 450, 750, Cooldown{Timer: character.NewTimer(), Duration: time.Minute * 5}, 375, 625))
}
func (character *Character) newSolidDynamiteSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, SolidDynamiteActionID, SpellSchoolFire, 213, 287, Cooldown{}, 0, 0))
}
func (character *Character) newDenseDynamiteSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, DenseDynamiteActionID, SpellSchoolFire, 340, 460, Cooldown{}, 0, 0))
}
func (character *Character) newThoriumGrenadeSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, ThoriumGrenadeActionID, SpellSchoolFire, 300, 500, Cooldown{}, 0, 0))
}
func (character *Character) newGoblinLandMineSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, GoblinLandMineActionID, SpellSchoolFire, 394, 506, Cooldown{}, 0, 0))
}

// Creates a spell object for the common explosive case.
func (character *Character) newRadiationBombSpellConfig(sharedTimer *Timer, actionID ActionID, minDamage float64, maxDamage float64, dotDamage float64, cooldown Cooldown) SpellConfig {
	return SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolFire,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagCastTimeNoGCD,

		Cast: CastConfig{
			DefaultCast: Cast{
				CastTime: time.Second,
			},
			IgnoreHaste: true,
			CD:          cooldown,
			SharedCD: Cooldown{
				Timer:    sharedTimer,
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
				Label: actionID.String(),
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
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(minDamage, maxDamage) * sim.Encounter.AOECapMultiplier()

				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				if result.Landed() {
					spell.Dot(aoeTarget).Apply(sim)
				}
			}
		},
	}
}
func (character *Character) newEzThroRadiationBombSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newRadiationBombSpellConfig(sharedTimer, EzThroRadiationBombActionID, 112, 188, 10, Cooldown{}))
}
func (character *Character) newHighYieldRadiationBombSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newRadiationBombSpellConfig(sharedTimer, HighYieldRadiationBombActionID, 150, 250, 25, Cooldown{}))
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
		defaultMCD.Spell.Flags |= SpellFlagCombatPotion
		character.AddMajorCooldown(defaultMCD)
	}
}

func makePotionActivation(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
	mcd := makePotionActivationInternal(potionType, character, potionCD)
	if mcd.Spell != nil {
		// Mark as 'Encounter Only' so that users are forced to select the generic Potion
		// placeholder action instead of specific potion spells, in APL prepull. This
		// prevents a mismatch between Consumes and Rotation settings.
		mcd.Spell.Flags |= SpellFlagEncounterOnly | SpellFlagPotion
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

func makeManaConsumableMCD(itemId int32, character *Character, cdTimer *Timer) MajorCooldown {
	minRoll := map[int32]float64{
		3385:  280.0,
		3827:  455.0,
		6149:  700.0,
		4381:  150.0,
		12662: 900.0,
		13443: 900.0,
		13444: 1350.0,
	}[itemId]

	maxRoll := map[int32]float64{
		3385:  360.0,
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

	return MajorCooldown{
		Type: CooldownTypeDPS,
		ShouldActivate: func(sim *Simulation, character *Character) bool {
			return true
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
				switch itemId {
				case 4623:
					lesserStoneshieldAura := character.NewTemporaryStatsAura("Lesser Stoneshield Potion", actionID, stats.Stats{stats.BonusArmor: 1000}, time.Second*90)
					lesserStoneshieldAura.Activate(sim)
				case 13455:
					greaterStoneshieldAura := character.NewTemporaryStatsAura("Greater Stoneshield Potion", actionID, stats.Stats{stats.BonusArmor: 2000}, time.Second*120)
					greaterStoneshieldAura.Activate(sim)
				}
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
	if potionType != proto.Potions_UnknownPotion {
		switch potionType {
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

		default:
			return MajorCooldown{}
		}
	} else {
		return MajorCooldown{}
	}
}

func registerConjuredCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	conjuredType := consumes.DefaultConjured

	if conjuredType == proto.Conjured_ConjuredDemonicRune || conjuredType == proto.Conjured_ConjuredMinorRecombobulator {
		itemId := map[proto.Conjured]int32{
			proto.Conjured_ConjuredDemonicRune:         12662,
			proto.Conjured_ConjuredMinorRecombobulator: 4381,
		}[conjuredType]

		character.AddMajorCooldown(makeManaConsumableMCD(itemId, character, character.GetConjuredCD()))
	}
}

func registerMildlyIrradiatedRejuvCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()

	if consumes.MildlyIrradiatedRejuvPot {
		actionID := ActionID{ItemID: 215162}
		healthMetrics := character.NewHealthMetrics(actionID)
		manaMetrics := character.NewManaMetrics(actionID)
		aura := character.NewTemporaryStatsAura("Mildly Irradiated Rejuvenation Potion", actionID, stats.Stats{stats.AttackPower: 40, stats.SpellDamage: 35}, time.Second*20)
		character.AddMajorCooldown(MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
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
