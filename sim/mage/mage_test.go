package mage

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterMage()
}

func TestArcane(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassMage,
			Level:      25,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase1TalentsArcane,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p1_generic"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p1_arcane"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Arcane", SpecOptions: PlayerOptionsArcane},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Level:      40,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase2TalentsArcane,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p2_arcane"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p2_arcane"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Arcane", SpecOptions: PlayerOptionsArcane},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase4TalentsArcane,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p4_arcane"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p4_arcane"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Arcane", SpecOptions: PlayerOptionsArcane},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase5TalentsArcane,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p5_arcane"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p5_spellfrost"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase5Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Arcane", SpecOptions: PlayerOptionsArcane},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassMage,
			Level:      25,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase1TalentsFire,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p1_fire"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p1_fire"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Level:      40,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase2TalentsFire,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p2_fire"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p2_fire"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Level:      50,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase3TalentsFire,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p3_fire"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p3_fire"),
			Buffs:       core.FullBuffsPhase3,
			Consumes:    Phase3Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase5TalentsFire,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p4_fire"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p4_fire"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase5TalentsFire,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p5_fire"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p5_fire"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase5Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Phase:      6,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase6TalentsFire,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p6_fire"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p6_fire"),
			Buffs:       core.FullBuffsPhase6,
			Consumes:    Phase6Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassMage,
			Level:      50,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase3TalentsFrost,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p3_frost_ffb"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p3_frost"),
			Buffs:       core.FullBuffsPhase3,
			Consumes:    Phase3Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Frost", SpecOptions: PlayerOptionsFrost},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase4TalentsFrost,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p4_frost"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p4_frost"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Frost", SpecOptions: PlayerOptionsFrost},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     phase5talentsfrost,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p5_frost"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p5_spellfrost"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase5Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Frost", SpecOptions: PlayerOptionsFrost},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassMage,
			Phase:      6,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     phase6talentsfrost,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p6_frost"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p6_spellfrost"),
			Buffs:       core.FullBuffsPhase6,
			Consumes:    Phase6Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Frost", SpecOptions: PlayerOptionsFrost},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var Phase1TalentsArcane = "22500502"
var Phase1TalentsFire = "-5050020121"

var Phase2TalentsArcane = "2250050310031531"
var Phase2TalentsFire = "-5050020123033151"
var Phase2TalentsFrostfire = Phase2TalentsFire

var Phase3TalentsFire = "-0550020123033151-2035"
var Phase3TalentsFrost = "-055-20350203100351051"

var Phase4TalentsArcane = "0550050210031531-054-203500001"
var Phase4TalentsFire = "21-5052300123033151-203500031"
var Phase4TalentsFrost = "-0550320003021-2035020310035105"

var Phase5TalentsArcane = "2500550010031531--2035020310004"
var Phase5TalentsFire = "21-5052300123033151-203500031"
var phase5talentsfrost = "250025001002--05350203100351051"

var Phase6TalentsFire = "-0552323121033151-203500031"
var phase6talentsfrost = "005055001--20350203110351351"

var PlayerOptionsArcane = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: &proto.Mage_Options{
			Armor: proto.Mage_Options_MageArmor,
		},
	},
}

var PlayerOptionsFire = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: &proto.Mage_Options{
			Armor: proto.Mage_Options_MoltenArmor,
		},
	},
}

var PlayerOptionsFrost = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: &proto.Mage_Options{
			Armor: proto.Mage_Options_IceArmor,
		},
	},
}

var Phase1Consumes = core.ConsumesCombo{
	Label: "P1-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion: proto.Potions_ManaPotion,
		FirePowerBuff: proto.FirePowerBuff_ElixirOfFirepower,
		Food:          proto.Food_FoodSmokedSagefish,
		MainHandImbue: proto.WeaponImbue_BlackfathomManaOil,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "P2-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_GreaterManaPotion,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfFirepower,
		FrostPowerBuff: proto.FrostPowerBuff_ElixirOfFrostPower,
		Food:           proto.Food_FoodSagefishDelight,
		MainHandImbue:  proto.WeaponImbue_LesserWizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_LesserArcaneElixir,
	},
}

var Phase3Consumes = core.ConsumesCombo{
	Label: "P3-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_MajorManaPotion,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfGreaterFirepower,
		FrostPowerBuff: proto.FrostPowerBuff_ElixirOfFrostPower,
		Food:           proto.Food_FoodNightfinSoup,
		MainHandImbue:  proto.WeaponImbue_WizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_ArcaneElixir,
	},
}

var Phase4Consumes = core.ConsumesCombo{
	Label: "P4-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_MajorManaPotion,
		Flask:          proto.Flask_FlaskOfSupremePower,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfGreaterFirepower,
		FrostPowerBuff: proto.FrostPowerBuff_ElixirOfFrostPower,
		Food:           proto.Food_FoodRunnTumTuberSurprise,
		MainHandImbue:  proto.WeaponImbue_WizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_GreaterArcaneElixir,
	},
}

var Phase5Consumes = core.ConsumesCombo{
	Label: "P5-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_MajorManaPotion,
		Flask:          proto.Flask_FlaskOfSupremePower,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfGreaterFirepower,
		FrostPowerBuff: proto.FrostPowerBuff_ElixirOfFrostPower,
		Food:           proto.Food_FoodRunnTumTuberSurprise,
		MainHandImbue:  proto.WeaponImbue_BrilliantWizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_GreaterArcaneElixir,
	},
}

var Phase6Consumes = core.ConsumesCombo{
	Label: "P6-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_MajorManaPotion,
		Flask:          proto.Flask_FlaskOfAncientKnowledge,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfGreaterFirepower,
		FrostPowerBuff: proto.FrostPowerBuff_ElixirOfFrostPower,
		Food:           proto.Food_FoodDarkclawBisque,
		MainHandImbue:  proto.WeaponImbue_EnchantedRepellent,
		SpellPowerBuff: proto.SpellPowerBuff_ElixirOfTheMageLord,
	},
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeStaff,
	},
	ArmorType: proto.ArmorType_ArmorTypeCloth,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeWand,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatIntellect,
	proto.Stat_StatSpellPower,
	proto.Stat_StatArcanePower,
	proto.Stat_StatFirePower,
	proto.Stat_StatFrostPower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
}
