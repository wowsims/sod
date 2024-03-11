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
	}))
}

func TestFrostFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassMage,
			Level:      40,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     Phase2TalentsFrostfire,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p2_frostfire"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "p2_frostfire"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Frostfire", SpecOptions: PlayerOptionsFire},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

// func TestFrost(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class: proto.Class_ClassMage,
// 		Race:  proto.Race_RaceTroll,

// 		GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p3_frost_alliance"),
// 		Talents:     FrostTalents,
// 		Consumes:    FullFrostConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Frost", SpecOptions: PlayerOptionsFrost},
// 		Rotation:    core.GetAplRotation("../../ui/mage/apls", "frost"),
// 		OtherRotations: []core.RotationCombo{
// 			core.GetAplRotation("../../ui/mage/apls", "frost_aoe"),
// 		},

// 		ItemFilter: ItemFilter,
// 	}))
// }

var Phase1TalentsArcane = "22500502"
var Phase1TalentsFire = "-5050020121"

var Phase2TalentsArcane = "2250050310031531"
var Phase2TalentsFire = "-5050020123033151"
var Phase2TalentsFrostfire = "-0550020122-203500031"

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
			Armor: proto.Mage_Options_MageArmor,
		},
	},
}

var PlayerOptionsFrost = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: &proto.Mage_Options{
			Armor: proto.Mage_Options_MageArmor,
		},
	},
}

var Phase1Consumes = core.ConsumesCombo{
	Label: "Phase 1 Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion: proto.Potions_ManaPotion,
		FirePowerBuff: proto.FirePowerBuff_ElixirOfFirepower,
		Food:          proto.Food_FoodSmokedSagefish,
		MainHandImbue: proto.WeaponImbue_BlackfathomManaOil,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "Phase 2 Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_GreaterManaPotion,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfFirepower,
		FrostPowerBuff: proto.FrostPowerBuff_ElixirOfFrostPower,
		Food:           proto.Food_FoodSagefishDelight,
		MainHandImbue:  proto.WeaponImbue_LesserWizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_LesserArcaneElixir,
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
