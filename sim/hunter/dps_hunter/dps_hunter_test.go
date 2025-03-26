package dps_hunter

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterDPSHunter()
}

const buildsDir = "../../../ui/hunter/builds"

func TestBM(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassHunter,
			Level:      25,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:          Phase1BMTalents,
			GearSet:          core.GetGearSet("../../../ui/hunter/gear_sets", "phase1"),
			Rotation:         core.GetAplRotation("../../../ui/hunter/apls", "p1_weave"),
			Buffs:            core.FullBuffsPhase1,
			Consumes:         Phase1Consumes,
			SpecOptions:      core.SpecOptionsCombo{Label: "Basic", SpecOptions: Phase1PlayerOptions},
			StartingDistance: core.MinRangedAttackRange,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassHunter,
			Level:      40,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:     Phase2BMTalents,
			GearSet:     core.GetGearSet("../../../ui/hunter/gear_sets", "p2_melee"),
			Rotation:    core.GetAplRotation("../../../ui/hunter/apls", "p2_melee"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: Phase2PlayerOptions},

			OtherGearSets:  []core.GearSetCombo{core.GetGearSet("../../../ui/hunter/gear_sets", "p2_ranged_bm")},
			OtherRotations: []core.RotationCombo{core.GetAplRotation("../../../ui/hunter/apls", "p2_ranged_bm")},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 4, 60, buildsDir, "p4_melee_dw", ItemFilters, proto.Stat_StatAgility, Stats),
	}))
}

func TestMM(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassHunter,
			Level:      25,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:          Phase1MMTalents,
			GearSet:          core.GetGearSet("../../../ui/hunter/gear_sets", "phase1"),
			Rotation:         core.GetAplRotation("../../../ui/hunter/apls", "p1_weave"),
			Buffs:            core.FullBuffsPhase1,
			Consumes:         Phase1Consumes,
			SpecOptions:      core.SpecOptionsCombo{Label: "Basic", SpecOptions: Phase1PlayerOptions},
			StartingDistance: core.MinRangedAttackRange,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassHunter,
			Level:      40,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:          Phase2MMTalents,
			GearSet:          core.GetGearSet("../../../ui/hunter/gear_sets", "p2_ranged_mm"),
			Rotation:         core.GetAplRotation("../../../ui/hunter/apls", "p2_ranged_mm"),
			Buffs:            core.FullBuffsPhase2,
			Consumes:         Phase2Consumes,
			SpecOptions:      core.SpecOptionsCombo{Label: "Basic", SpecOptions: Phase2PlayerOptions},
			StartingDistance: core.MaxRangedAttackRange,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 4, 60, buildsDir, "p4_ranged", ItemFilters, proto.Stat_StatAgility, Stats),
	}))
}

func TestSV(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassHunter,
			Level:      25,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:     Phase1SVTalents,
			GearSet:     core.GetGearSet("../../../ui/hunter/gear_sets", "phase1"),
			Rotation:    core.GetAplRotation("../../../ui/hunter/apls", "p1_weave"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: Phase1PlayerOptions},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassHunter,
			Level:      40,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase2SVTalents,
			GearSet:     core.GetGearSet("../../../ui/hunter/gear_sets", "p2_melee"),
			Rotation:    core.GetAplRotation("../../../ui/hunter/apls", "p2_melee"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: Phase2PlayerOptions},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassHunter,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:          Phase4WeaveTalents,
			GearSet:          core.GetGearSet("../../../ui/hunter/gear_sets", "p4_weave"),
			Rotation:         core.GetAplRotation("../../../ui/hunter/apls", "p4_weave"),
			Buffs:            core.FullBuffsPhase4,
			Consumes:         Phase4Consumes,
			SpecOptions:      core.SpecOptionsCombo{Label: "Weave", SpecOptions: Phase4PlayerOptions},
			StartingDistance: core.MinRangedAttackRange,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func TestRangedHunter(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 4, 60, buildsDir, "p4_ranged", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 5, 60, buildsDir, "p5_ranged", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 6, 60, buildsDir, "p6_ranged", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 7, 60, buildsDir, "p7_ranged", ItemFilters, proto.Stat_StatAgility, Stats),
	}))
}

func Test2HMeleeHunter(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 5, 60, buildsDir, "p5_melee_2h", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 6, 60, buildsDir, "p6_melee_2h", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 7, 60, buildsDir, "p7_melee_2h", ItemFilters, proto.Stat_StatAgility, Stats),
	}))
}

func TestDWMeleeHunter(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 4, 60, buildsDir, "p4_melee_dw", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 5, 60, buildsDir, "p5_melee_dw", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 6, 60, buildsDir, "p6_melee_dw", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 7, 60, buildsDir, "p7_melee_dw_bm", ItemFilters, proto.Stat_StatAgility, Stats),
	}))
}

var Phase1BMTalents = "53000200501"
var Phase1MMTalents = "-050515"
var Phase1SVTalents = "--33502001101"

var Phase2BMTalents = "5300021150501251"
var Phase2MMTalents = "-05551001503051"
var Phase2SVTalents = "--335020051030315"

var Phase4WeaveTalents = "-055500005-3305202202303051"
var Phase4RangedMMTalents = "-05451002503051-33400023023"
var Phase4RangedSVTalents = "1-054510005-334000250230305"

var Phase1Consumes = core.ConsumesCombo{
	Label: "P1-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir: proto.AgilityElixir_ElixirOfLesserAgility,
		DefaultPotion: proto.Potions_ManaPotion,
		Food:          proto.Food_FoodSmokedSagefish,
		MainHandImbue: proto.WeaponImbue_WildStrikes,
		OffHandImbue:  proto.WeaponImbue_BlackfathomSharpeningStone,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "P2-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfAgility,
		DefaultPotion:     proto.Potions_ManaPotion,
		DragonBreathChili: true,
		Food:              proto.Food_FoodSagefishDelight,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_SolidWeightstone,
		SpellPowerBuff:    proto.SpellPowerBuff_LesserArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase4Consumes = core.ConsumesCombo{
	Label: "P4-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		DefaultPotion:     proto.Potions_ManaPotion,
		DragonBreathChili: true,
		Flask:             proto.Flask_FlaskOfSupremePower,
		Food:              proto.Food_FoodSagefishDelight,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_ElementalSharpeningStone,
		SpellPowerBuff:    proto.SpellPowerBuff_GreaterArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
	},
}

var Phase1PlayerOptions = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options: &proto.Hunter_Options{
			Ammo:           proto.Hunter_Options_RazorArrow,
			PetType:        proto.Hunter_Options_Cat,
			PetUptime:      1,
			PetAttackSpeed: 2.0,
		},
	},
}

var Phase2PlayerOptions = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options: &proto.Hunter_Options{
			Ammo:           proto.Hunter_Options_JaggedArrow,
			PetType:        proto.Hunter_Options_Cat,
			PetUptime:      1,
			PetAttackSpeed: 2.0,
		},
	},
}

var Phase4PlayerOptions = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options: &proto.Hunter_Options{
			Ammo:                 proto.Hunter_Options_JaggedArrow,
			PetType:              proto.Hunter_Options_PetNone,
			PetUptime:            1,
			PetAttackSpeed:       2.0,
			SniperTrainingUptime: 1.0,
		},
	},
}

var ItemFilters = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeMail,
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypePolearm,
		proto.WeaponType_WeaponTypeStaff,
		proto.WeaponType_WeaponTypeSword,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeBow,
		proto.RangedWeaponType_RangedWeaponTypeCrossbow,
		proto.RangedWeaponType_RangedWeaponTypeGun,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatAgility,
	proto.Stat_StatAttackPower,
	proto.Stat_StatRangedAttackPower,
	proto.Stat_StatMeleeCrit,
	proto.Stat_StatMeleeHit,
}
