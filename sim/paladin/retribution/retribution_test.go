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
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p1ret"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p1ret"),
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
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p2retsoc"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p2ret"),
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
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p3retsom"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p3ret"),
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

			Talents:  Phase45RetTalents,
			GearSet:  core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p4ret-twisting-6pcT1"),
			Rotation: core.GetAplRotation("../../../ui/retribution_paladin/apls", "p4ret-twisting-6pcT1"),

			OtherGearSets:  []core.GearSetCombo{core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p4rettwist")},
			OtherRotations: []core.RotationCombo{core.GetAplRotation("../../../ui/retribution_paladin/apls", "p4ret")},
			Buffs:          core.FullBuffsPhase5,
			Consumes:       Phase4Consumes,
			SpecOptions:    core.SpecOptionsCombo{Label: "P4 Seal of Martyrdom Ret", SpecOptions: PlayerOptionsSealofMartyrdom},

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

			Talents:        Phase45RetTalents,
			GearSet:        core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p5twisting"),
			Rotation:       core.GetAplRotation("../../../ui/retribution_paladin/apls", "p5ret-twist-4DR-3.5-3.6"),
			OtherRotations: []core.RotationCombo{core.GetAplRotation("../../../ui/retribution_paladin/apls", "p5ret-twist-4DR-3.7-4.0")},
			Buffs:          core.FullBuffsPhase5,
			Consumes:       Phase4Consumes,
			SpecOptions:    core.SpecOptionsCombo{Label: "P5 Seal of Martyrdom Ret", SpecOptions: PlayerOptionsSealofMartyrdom},

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

			Talents:     Phase45RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p4ret-exodin-6pcT1"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p4ret-exodin-6pcT1"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P4 Seal of Martyrdom Ret", SpecOptions: PlayerOptionsSealofMartyrdom},

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

			Talents:     Phase45RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p5exodin"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p5ret-exodin-6CF2DR"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P5 Seal of Martyrdom Ret", SpecOptions: PlayerOptionsSealofMartyrdom},

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
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p2retsom"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p2ret"),
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
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p5shockadin"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p5Shockadin"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P5 Seal of Righteousness Shockadin", SpecOptions: PlayerOptionsSealofRighteousness},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var Phase1RetTalents = "--05230051"
var Phase2RetTalents = "--532300512003151"
var Phase2ShockadinTalents = "55050100521151--"
var Phase3RetTalents = "500501--53230051200315"
var Phase45RetTalents = "500501-503-52230351200315"
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

var PlayerOptionsSealofRighteousness = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: optionsSealOfRighteousness,
	},
}

var optionsSealOfCommand = &proto.PaladinOptions{
	PrimarySeal: proto.PaladinSeal_Command,
}

var optionsSealOfMartyrdom = &proto.PaladinOptions{
	PrimarySeal: proto.PaladinSeal_Martyrdom,
}

var optionsSealOfRighteousness = &proto.PaladinOptions{
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
