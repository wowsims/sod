package shadow

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get caster sets included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterShadowPriest()
}

func TestShadow(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassPriest,
			Level:      25,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:          Phase1Talents,
			GearSet:          core.GetGearSet("../../../ui/shadow_priest/gear_sets", "phase_1"),
			Rotation:         core.GetAplRotation("../../../ui/shadow_priest/apls", "phase_1"),
			Buffs:            core.FullBuffsPhase1,
			Consumes:         Phase1Consumes,
			SpecOptions:      core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
			StartingDistance: core.MaxShortSpellRange,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPriest,
			Level:      40,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:          Phase2Talents,
			GearSet:          core.GetGearSet("../../../ui/shadow_priest/gear_sets", "phase_2"),
			Rotation:         core.GetAplRotation("../../../ui/shadow_priest/apls", "phase_2"),
			Buffs:            core.FullBuffsPhase2,
			Consumes:         Phase2Consumes,
			SpecOptions:      core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
			StartingDistance: core.MaxShortSpellRange,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPriest,
			Level:      50,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:          Phase3Talents,
			GearSet:          core.GetGearSet("../../../ui/shadow_priest/gear_sets", "phase_3"),
			Rotation:         core.GetAplRotation("../../../ui/shadow_priest/apls", "phase_3"),
			Buffs:            core.FullBuffsPhase3,
			Consumes:         Phase3Consumes,
			SpecOptions:      core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
			StartingDistance: core.MaxShortSpellRange,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPriest,
			Level:      60,
			Phase:      4,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:          Phase4Talents,
			GearSet:          core.GetGearSet("../../../ui/shadow_priest/gear_sets", "phase_4"),
			Rotation:         core.GetAplRotation("../../../ui/shadow_priest/apls", "phase_4"),
			Buffs:            core.FullBuffsPhase4,
			Consumes:         Phase4Consumes,
			SpecOptions:      core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
			StartingDistance: core.MaxShortSpellRange,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPriest,
			Level:      60,
			Phase:      5,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents: Phase4Talents,
			GearSet: core.GetGearSet("../../../ui/shadow_priest/gear_sets", "phase_5_t1"),
			OtherGearSets: []core.GearSetCombo{
				core.GetGearSet("../../../ui/shadow_priest/gear_sets", "phase_5_t2"),
			},
			Rotation:         core.GetAplRotation("../../../ui/shadow_priest/apls", "phase_5"),
			Buffs:            core.FullBuffsPhase5,
			Consumes:         Phase4Consumes,
			SpecOptions:      core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
			StartingDistance: core.MaxShortSpellRange,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPriest,
			Level:      60,
			Phase:      6,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:          Phase4Talents,
			GearSet:          core.GetGearSet("../../../ui/shadow_priest/gear_sets", "phase_6"),
			Rotation:         core.GetAplRotation("../../../ui/shadow_priest/apls", "phase_6"),
			Buffs:            core.FullBuffsPhase6,
			Consumes:         Phase6Consumes,
			SpecOptions:      core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
			StartingDistance: core.MaxShortSpellRange,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		core.GetTestBuildFromJSON(proto.Class_ClassPriest, 7, 60, "../../../ui/shadow_priest/builds", "phase_7", ItemFilters, proto.Stat_StatSpellPower, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassPriest, 8, 60, "../../../ui/shadow_priest/builds", "phase_8", ItemFilters, proto.Stat_StatSpellPower, Stats),
	}))
}

var Phase1Talents = "-20535000001"
var Phase2Talents = "--5022204002501251"
var Phase3Talents = "-0055-5022204002501251"
var Phase4Talents = "0512301302--5002504103501251"

var Phase1Consumes = core.ConsumesCombo{
	Label: "P1-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion: proto.Potions_ManaPotion,
		Food:          proto.Food_FoodSmokedSagefish,
		MainHandImbue: proto.WeaponImbue_BlackfathomManaOil,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "P2-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_GreaterManaPotion,
		Food:           proto.Food_FoodSagefishDelight,
		MainHandImbue:  proto.WeaponImbue_LesserWizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_LesserArcaneElixir,
	},
}

var Phase3Consumes = core.ConsumesCombo{
	Label: "P3-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:   proto.Potions_GreaterManaPotion,
		Food:            proto.Food_FoodNightfinSoup,
		MainHandImbue:   proto.WeaponImbue_LesserWizardOil,
		SpellPowerBuff:  proto.SpellPowerBuff_ArcaneElixir,
		ShadowPowerBuff: proto.ShadowPowerBuff_ElixirOfShadowPower,
	},
}

var Phase4Consumes = core.ConsumesCombo{
	Label: "P4-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:   proto.Potions_MajorManaPotion,
		Flask:           proto.Flask_FlaskOfSupremePower,
		Food:            proto.Food_FoodRunnTumTuberSurprise,
		MainHandImbue:   proto.WeaponImbue_WizardOil,
		SpellPowerBuff:  proto.SpellPowerBuff_GreaterArcaneElixir,
		ShadowPowerBuff: proto.ShadowPowerBuff_ElixirOfShadowPower,
	},
}

var Phase6Consumes = core.ConsumesCombo{
	Label: "P6-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:   proto.Potions_MajorManaPotion,
		Flask:           proto.Flask_FlaskOfAncientKnowledge,
		Food:            proto.Food_FoodDarkclawBisque,
		MainHandImbue:   proto.WeaponImbue_EnchantedRepellent,
		SpellPowerBuff:  proto.SpellPowerBuff_ElixirOfTheMageLord,
		ShadowPowerBuff: proto.ShadowPowerBuff_ElixirOfShadowPower,
	},
}

var PlayerOptionsBasic = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			Armor: proto.ShadowPriest_Options_InnerFire,
		},
	},
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeMace,
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
	proto.Stat_StatShadowPower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
}
