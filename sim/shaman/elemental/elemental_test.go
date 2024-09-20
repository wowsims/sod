package elemental

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterElementalShaman()
}

func TestElemental(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassShaman,
			Level:      25,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     Phase1Talents,
			GearSet:     core.GetGearSet("../../../ui/elemental_shaman/gear_sets", "phase_1"),
			Rotation:    core.GetAplRotation("../../../ui/elemental_shaman/apls", "phase_1"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Adaptive", SpecOptions: PlayerOptionsAdaptive},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassShaman,
			Level:      40,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     Phase2Talents,
			GearSet:     core.GetGearSet("../../../ui/elemental_shaman/gear_sets", "phase_2"),
			Rotation:    core.GetAplRotation("../../../ui/elemental_shaman/apls", "phase_2"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Adaptive", SpecOptions: PlayerOptionsAdaptive},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassShaman,
			Level:      50,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     Phase3Talents,
			GearSet:     core.GetGearSet("../../../ui/elemental_shaman/gear_sets", "phase_3"),
			Rotation:    core.GetAplRotation("../../../ui/elemental_shaman/apls", "phase_3"),
			Buffs:       core.FullBuffsPhase3,
			Consumes:    Phase3Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Adaptive", SpecOptions: PlayerOptionsAdaptive},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassShaman,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents: Phase4Talents,
			GearSet: core.GetGearSet("../../../ui/elemental_shaman/gear_sets", "phase_4"),
			OtherGearSets: []core.GearSetCombo{
				core.GetGearSet("../../../ui/elemental_shaman/gear_sets", "phase_5"),
			},
			Rotation: core.GetAplRotation("../../../ui/elemental_shaman/apls", "phase_4"),
			OtherRotations: []core.RotationCombo{
				core.GetAplRotation("../../../ui/elemental_shaman/apls", "phase_5"),
			},
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Adaptive", SpecOptions: PlayerOptionsAdaptive},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var Phase1Talents = "25003105"
var Phase2Talents = "550031550000151"
var Phase3Talents = "550031550000151-500203"
var Phase4Talents = "550301550000151--50205300005"

var PlayerOptionsAdaptive = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Options: &proto.ElementalShaman_Options{},
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
		Food:           proto.Food_FoodSagefishDelight,
		MainHandImbue:  proto.WeaponImbue_LesserWizardOil,
		OffHandImbue:   proto.WeaponImbue_LesserWizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_LesserArcaneElixir,
	},
}

var Phase3Consumes = core.ConsumesCombo{
	Label: "Phase 3 Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_GreaterManaPotion,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfGreaterFirepower,
		Food:           proto.Food_FoodNightfinSoup,
		MainHandImbue:  proto.WeaponImbue_FlametongueWeapon,
		OffHandImbue:   proto.WeaponImbue_LesserWizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_ArcaneElixir,
		StrengthBuff:   proto.StrengthBuff_ElixirOfGiants,
	},
}

var Phase4Consumes = core.ConsumesCombo{
	Label: "Phase 4 Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_MajorManaPotion,
		Flask:          proto.Flask_FlaskOfSupremePower,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfGreaterFirepower,
		Food:           proto.Food_FoodRunnTumTuberSurprise,
		MainHandImbue:  proto.WeaponImbue_FlametongueWeapon,
		OffHandImbue:   proto.WeaponImbue_ConductiveShieldCoating,
		SpellPowerBuff: proto.SpellPowerBuff_GreaterArcaneElixir,
	},
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeShield,
		proto.WeaponType_WeaponTypeStaff,
	},
	ArmorType: proto.ArmorType_ArmorTypeMail,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeTotem,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatIntellect,
	proto.Stat_StatSpellPower,
	proto.Stat_StatNaturePower,
	proto.Stat_StatFirePower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
}
