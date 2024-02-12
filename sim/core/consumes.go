package core

import (
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Registers all consume-related effects to the Agent.
// TODO: Classic Consumes
func applyConsumeEffects(agent Agent, partyBuffs *proto.PartyBuffs) {
	character := agent.GetCharacter()
	consumes := character.Consumes
	if consumes == nil {
		return
	}

	if consumes.Flask != proto.Flask_FlaskUnknown {
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
		}
	}

	if character.HasMHWeapon() {
		addImbueStats(character, consumes.MainHandImbue, true)
	}
	if character.HasOHWeapon() {
		addImbueStats(character, consumes.OffHandImbue, false)
	}

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

	if consumes.BoglingRoot {
		character.PseudoStats.BonusDamage += 1
	}

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

	if character.HasProfession(proto.Profession_Enchanting) && consumes.EnchantedSigil != proto.EnchantedSigil_UnknownSigil {
		switch consumes.EnchantedSigil {
		case proto.EnchantedSigil_InnovationSigil:
			character.AddStats(stats.Stats{
				stats.AttackPower:       20,
				stats.RangedAttackPower: 20,
				stats.SpellPower:        20,
			})
		}
	}

	registerPotionCD(agent, consumes)
	registerConjuredCD(agent, consumes)
	registerExplosivesCD(agent, consumes)
}

func addImbueStats(character *Character, imbue proto.WeaponImbue, isMh bool) {
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
				stats.MP5:     12,
				stats.Healing: 25,
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
		}
	}
}

var SapperActionID = ActionID{ItemID: 10646}
var SolidDynamiteActionID = ActionID{ItemID: 10507}
var DenseDynamiteActionID = ActionID{ItemID: 18641}
var ThoriumGrenadeActionID = ActionID{ItemID: 15993}
var EzThroRadiationBombActionID = ActionID{ItemID: 215168}
var HighYieldRadiationBombActionID = ActionID{ItemID: 215127}

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
		}

		character.AddMajorCooldown(MajorCooldown{
			Spell:    filler,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 10,
		})
	}
}

// Creates a spell object for the common explosive case.
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
		},

		// Explosives always have 1% resist chance, so just give them hit cap.
		BonusHitRating:   100 * SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   2,
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

// Creates a spell object for the common explosive case.
func (character *Character) newRadiationBombSpellConfig(sharedTimer *Timer, actionID ActionID, minDamage float64, maxDamage float64, dotDamage float64, cooldown Cooldown) SpellConfig {
	return SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolFire,
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
		BonusHitRating:   100 * SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

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

				dot.SnapshotBaseDamage = dotDamage

				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])

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

func registerPotionCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	defaultPotion := consumes.DefaultPotion

	potionCD := character.NewTimer()

	if defaultPotion == proto.Potions_UnknownPotion {
		return
	}

	var defaultMCD MajorCooldown
	defaultMCD = makePotionActivation(defaultPotion, character, potionCD)

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
	}[itemId]

	maxRoll := map[int32]float64{
		3385:  360.0,
		3827:  585.0,
		6149:  900.0,
		4381:  250.0,
		12662: 1500.0,
	}[itemId]

	cdDuration := map[int32]time.Duration{
		3385:  time.Minute * 2,
		3827:  time.Minute * 2,
		6149:  time.Minute * 2,
		4381:  time.Minute * 5,
		12662: time.Minute * 2,
	}[itemId]

	actionID := ActionID{ItemID: itemId}
	manaMetrics := character.NewManaMetrics(actionID)

	return MajorCooldown{
		Type: CooldownTypeMana,
		ShouldActivate: func(sim *Simulation, character *Character) bool {
			// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
			totalRegen := character.ManaRegenPerSecondWhileCasting() * 2
			return (character.MaxMana()-(character.CurrentMana()+totalRegen) >= maxRoll) && !character.PseudoStats.Shapeshifted
		},
		Spell: character.GetOrRegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    cdTimer,
					Duration: cdDuration,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {

				manaGain := sim.RollWithLabel(minRoll, maxRoll, "Mana Consumable")
				character.AddMana(sim, manaGain, manaMetrics)
			},
		}),
	}
}

func makePotionActivationInternal(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
	if potionType == proto.Potions_LesserManaPotion || potionType == proto.Potions_ManaPotion || potionType == proto.Potions_GreaterManaPotion {
		itemId := map[proto.Potions]int32{
			proto.Potions_LesserManaPotion:  3385,
			proto.Potions_ManaPotion:        3827,
			proto.Potions_GreaterManaPotion: 6149,
		}[potionType]

		return makeManaConsumableMCD(itemId, character, potionCD)
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
