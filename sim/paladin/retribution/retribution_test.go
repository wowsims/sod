package retribution

import (
	"testing"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterRetributionPaladin()
}

func TestRetribution(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassPaladin,
			Level:      25,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase1RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p1-ret"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p1-ret"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P1 Seal of Command Ret", SpecOptions: PlayerOptionsSealofCommand},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPaladin,
			Level:      40,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase2RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p2-retsoc"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p2-ret"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P2 Seal of Command Ret", SpecOptions: PlayerOptionsSealofCommand},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPaladin,
			Level:      50,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase3RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p3-retsom"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p3-ret"),
			Buffs:       core.FullBuffsPhase3,
			Consumes:    Phase3Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P3 Seal of Martyrdom Ret", SpecOptions: PlayerOptionsSealofMartyrdom},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPaladin,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:  Phase456RetTalents,
			GearSet:  core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p4-twisting-6pcT1"),
			Rotation: core.GetAplRotation("../../../ui/retribution_paladin/apls", "p4-twisting-6pcT1"),

			OtherGearSets:  []core.GearSetCombo{core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p4-twist")},
			OtherRotations: []core.RotationCombo{core.GetAplRotation("../../../ui/retribution_paladin/apls", "p5p6p7-twist")},
			Buffs:          core.FullBuffsPhase5,
			Consumes:       Phase4Consumes,
			SpecOptions:    core.SpecOptionsCombo{Label: "P4 Twist", SpecOptions: PlayerOptionsSealofMartyrdomStopAttack},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPaladin,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase456RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p5-twisting"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p5p6p7-twist"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase5Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P5 Twist", SpecOptions: PlayerOptionsSealofMartyrdomStopAttack},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPaladin,
			Phase:      6,
			Level:      60,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase456RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p6-twisting"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p5p6p7-twist"),
			Buffs:       core.FullBuffsPhase6,
			Consumes:    Phase6Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P6 Twist", SpecOptions: PlayerOptionsSealofMartyrdomStopAttack},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func TestExodin(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassPaladin,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase456RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p4-exodin-6pcT1"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p4-exodin-6pcT1"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P4 Exodin", SpecOptions: PlayerOptionsSealofMartyrdom},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPaladin,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase456RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p5-exodin"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p5-exodin-6CF-2DR"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase5Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P5 Exodin", SpecOptions: PlayerOptionsSealofMartyrdom},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func TestShockadin(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassPaladin,
			Level:      40,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase2ShockadinTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p2-retsom"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p2-ret"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P2 Seal of Martyrdom Shockadin", SpecOptions: PlayerOptionsSealofMartyrdom},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPaladin,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase45ShockadinTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p5-shockadin"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p5-shockadin"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase5Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P5 Shockadin", SpecOptions: PlayerOptionsSealofRighteousness},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func TestSealStacking(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassPaladin,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase456RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p5-seal-stacking"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p5-seal-stacking-6CF-2DR"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase5Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P5 Seal Stacking", SpecOptions: PlayerOptionsSealofMartyrdom},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var Phase1RetTalents = "--05230051"
var Phase2RetTalents = "--532300512003151"
var Phase2ShockadinTalents = "55050100521151--"
var Phase3RetTalents = "500501--53230051200315"
var Phase456RetTalents = "500501-503-52230351200315"
var Phase45ShockadinTalents = "55053100501051--052303511"

var Phase1Consumes = core.ConsumesCombo{
	Label: "P1-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir: proto.AgilityElixir_ElixirOfLesserAgility,
		DefaultPotion: proto.Potions_ManaPotion,
		FirePowerBuff: proto.FirePowerBuff_ElixirOfFirepower,
		MainHandImbue: proto.WeaponImbue_WildStrikes,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "P2-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfAgility,
		DefaultPotion:     proto.Potions_ManaPotion,
		DragonBreathChili: true,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfFirepower,
		Food:              proto.Food_FoodSagefishDelight,
		MainHandImbue:     proto.WeaponImbue_WindfuryWeapon,
		SpellPowerBuff:    proto.SpellPowerBuff_LesserArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase3Consumes = core.ConsumesCombo{
	Label: "P3-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		DefaultPotion:     proto.Potions_MajorManaPotion,
		DefaultConjured:   proto.Conjured_ConjuredDemonicRune,
		DragonBreathChili: true,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfFirepower,
		Food:              proto.Food_FoodBlessSunfruit,
		MainHandImbue:     proto.WeaponImbue_WindfuryWeapon,
		SpellPowerBuff:    proto.SpellPowerBuff_GreaterArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_ElixirOfGiants,
		EnchantedSigil:    proto.EnchantedSigil_LivingDreamsSigil,
		AttackPowerBuff:   proto.AttackPowerBuff_WinterfallFirewater,
		ZanzaBuff:         proto.ZanzaBuff_AtalaiMojoOfWar,
	},
}
var Phase4Consumes = core.ConsumesCombo{
	Label: "P4-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:     proto.Potions_MajorManaPotion,
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		Flask:             proto.Flask_FlaskOfSupremePower,
		SpellPowerBuff:    proto.SpellPowerBuff_GreaterArcaneElixir,
		DragonBreathChili: true,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfFirepower,
		Food:              proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_ConductiveShieldCoating,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
		EnchantedSigil:    proto.EnchantedSigil_FlowingWatersSigil,
	},
}

var Phase5Consumes = core.ConsumesCombo{
	Label: "P5-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		DefaultPotion:     proto.Potions_MajorManaPotion,
		DragonBreathChili: true,
		EnchantedSigil:    proto.EnchantedSigil_FlowingWatersSigil,
		Flask:             proto.Flask_FlaskOfSupremePower,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfFirepower,
		Food:              proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_MagnificentTrollshine,
		SpellPowerBuff:    proto.SpellPowerBuff_GreaterArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
	},
}

var Phase6Consumes = core.ConsumesCombo{
	Label: "P6-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:            proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:          proto.AttackPowerBuff_JujuMight,
		DefaultConjured:          proto.Conjured_ConjuredDemonicRune,
		DefaultPotion:            proto.Potions_MajorManaPotion,
		DragonBreathChili:        true,
		EnchantedSigil:           proto.EnchantedSigil_WrathOfTheStormSigil,
		FirePowerBuff:            proto.FirePowerBuff_ElixirOfGreaterFirepower,
		Flask:                    proto.Flask_FlaskOfMadness,
		Food:                     proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:            proto.WeaponImbue_WildStrikes,
		ManaRegenElixir:          proto.ManaRegenElixir_MagebloodPotion,
		MildlyIrradiatedRejuvPot: true,
		OffHandImbue:             proto.WeaponImbue_EnchantedRepellent,
		SapperExplosive:          proto.SapperExplosive_SapperFumigator,
		SpellPowerBuff:           proto.SpellPowerBuff_ElixirOfTheMageLord,
		StrengthBuff:             proto.StrengthBuff_JujuPower,
		ZanzaBuff:                proto.ZanzaBuff_ROIDS,
	},
}

var PlayerOptionsSealofCommand = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: optionsSealOfCommand,
	},
}

var PlayerOptionsSealofMartyrdom = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: optionsSealOfMartyrdom,
	},
}

var PlayerOptionsSealofMartyrdomStopAttack = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: optionsSealOfMartyrdomStopAttack,
	},
}

var PlayerOptionsSealofRighteousness = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: optionsSealOfRighteousness,
	},
}

var optionsSealOfCommand = &proto.PaladinOptions{
	Aura:        proto.PaladinAura_SanctityAura,
	PrimarySeal: proto.PaladinSeal_Command,
}

var optionsSealOfMartyrdom = &proto.PaladinOptions{
	Aura:        proto.PaladinAura_SanctityAura,
	PrimarySeal: proto.PaladinSeal_Martyrdom,
}

var optionsSealOfMartyrdomStopAttack = &proto.PaladinOptions{
	Aura:                            proto.PaladinAura_SanctityAura,
	PrimarySeal:                     proto.PaladinSeal_Martyrdom,
	IsUsingCrusaderStrikeStopAttack: true,
	IsUsingExorcismStopAttack:       true,
	IsUsingDivineStormStopAttack:    true,
	IsUsingJudgementStopAttack:      true,
}

var optionsSealOfRighteousness = &proto.PaladinOptions{
	Aura:        proto.PaladinAura_SanctityAura,
	PrimarySeal: proto.PaladinSeal_Righteousness,
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypePolearm,
		proto.WeaponType_WeaponTypeShield,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeLibram,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatStrength,
	proto.Stat_StatAgility,
	proto.Stat_StatAttackPower,
	proto.Stat_StatMeleeHit,
	proto.Stat_StatMeleeCrit,
	proto.Stat_StatSpellPower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
}
