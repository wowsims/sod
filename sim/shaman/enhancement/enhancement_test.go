package enhancement

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterEnhancementShaman()
}

func TestEnhancement(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassShaman,
			Level:      25,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     Phase1Talents,
			GearSet:     core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_1"),
			Rotation:    core.GetAplRotation("../../../ui/enhancement_shaman/apls", "phase_1"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Sync Auto", SpecOptions: PlayerOptionsSyncAuto},
			OtherSpecOptions: []core.SpecOptionsCombo{
				{Label: "Sync Delay OH", SpecOptions: PlayerOptionsSyncDelayOH},
			},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassShaman,
			Level:      40,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:  Phase2Talents,
			GearSet:  core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_2"),
			Rotation: core.GetAplRotation("../../../ui/enhancement_shaman/apls", "phase_2"),
			Buffs:    core.FullBuffsPhase2,
			Consumes: Phase2ConsumesWFWF,
			OtherConsumes: []core.ConsumesCombo{
				Phase2ConsumesWFFT,
			},
			SpecOptions: core.SpecOptionsCombo{Label: "Sync Auto", SpecOptions: PlayerOptionsSyncAuto},
			OtherSpecOptions: []core.SpecOptionsCombo{
				{Label: "Sync Delay OH", SpecOptions: PlayerOptionsSyncDelayOH},
			},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassShaman,
			Level:      50,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:  Phase3Talents,
			GearSet:  core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_3"),
			Rotation: core.GetAplRotation("../../../ui/enhancement_shaman/apls", "phase_3"),
			Buffs:    core.FullBuffsPhase3,
			Consumes: Phase3ConsumesWFWF,
			OtherConsumes: []core.ConsumesCombo{
				Phase3ConsumesWFFT,
			},
			SpecOptions: core.SpecOptionsCombo{Label: "Sync Auto", SpecOptions: PlayerOptionsSyncAuto},
			OtherSpecOptions: []core.SpecOptionsCombo{
				{Label: "Sync Delay OH", SpecOptions: PlayerOptionsSyncDelayOH},
			},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassShaman,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents: Phase4Talents,
			GearSet: core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_4_dw"),
			OtherGearSets: []core.GearSetCombo{
				core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_4_2h"),
			},
			Rotation:    core.GetAplRotation("../../../ui/enhancement_shaman/apls", "phase_4"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4ConsumesWFWF,
			SpecOptions: core.SpecOptionsCombo{Label: "Sync Auto", SpecOptions: PlayerOptionsSyncAuto},
			OtherSpecOptions: []core.SpecOptionsCombo{
				{Label: "Sync Delay OH", SpecOptions: PlayerOptionsSyncDelayOH},
			},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassShaman,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents: Phase4Talents,
			GearSet: core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_5_dw"),
			OtherGearSets: []core.GearSetCombo{
				core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_5_2h"),
			},
			Rotation:    core.GetAplRotation("../../../ui/enhancement_shaman/apls", "phase_5"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase4ConsumesWFWF,
			SpecOptions: core.SpecOptionsCombo{Label: "Sync Auto", SpecOptions: PlayerOptionsSyncAuto},
			OtherSpecOptions: []core.SpecOptionsCombo{
				{Label: "Sync Delay OH", SpecOptions: PlayerOptionsSyncDelayOH},
			},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var Phase1Talents = "-5005202101"
var Phase2Talents = "-5005202105023051"
var Phase3Talents = "05003-5005132105023051"
var Phase4Talents = "25003105003-5005032105023051"

var PlayerOptionsSyncDelayOH = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options: optionsSyncDelayOffhand,
	},
}

var PlayerOptionsSyncAuto = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options: optionsSyncAuto,
	},
}

var optionsSyncDelayOffhand = &proto.EnhancementShaman_Options{
	SyncType: proto.ShamanSyncType_DelayOffhandSwings,
}

var optionsSyncAuto = &proto.EnhancementShaman_Options{
	SyncType: proto.ShamanSyncType_Auto,
}

var Phase1Consumes = core.ConsumesCombo{
	Label: "Phase 1 Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir: proto.AgilityElixir_ElixirOfLesserAgility,
		DefaultPotion: proto.Potions_ManaPotion,
		FirePowerBuff: proto.FirePowerBuff_ElixirOfFirepower,
		MainHandImbue: proto.WeaponImbue_RockbiterWeapon,
		OffHandImbue:  proto.WeaponImbue_RockbiterWeapon,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase2ConsumesWFWF = core.ConsumesCombo{
	Label: "Phase 2 Consumes WF/WF",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfAgility,
		DefaultPotion:     proto.Potions_ManaPotion,
		DragonBreathChili: true,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfFirepower,
		Food:              proto.Food_FoodSagefishDelight,
		MainHandImbue:     proto.WeaponImbue_WindfuryWeapon,
		OffHandImbue:      proto.WeaponImbue_WindfuryWeapon,
		SpellPowerBuff:    proto.SpellPowerBuff_LesserArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase2ConsumesWFFT = core.ConsumesCombo{
	Label: "Phase 2 Consumes WF/FT",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfAgility,
		DefaultPotion:     proto.Potions_ManaPotion,
		DragonBreathChili: true,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfFirepower,
		Food:              proto.Food_FoodSagefishDelight,
		MainHandImbue:     proto.WeaponImbue_WindfuryWeapon,
		OffHandImbue:      proto.WeaponImbue_FlametongueWeapon,
		SpellPowerBuff:    proto.SpellPowerBuff_LesserArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_ScrollOfStrength,
	},
}

var Phase3ConsumesWFWF = core.ConsumesCombo{
	Label: "Phase 3 Consumes WF/WF",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfAgility,
		DefaultPotion:     proto.Potions_ManaPotion,
		DragonBreathChili: true,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfFirepower,
		Food:              proto.Food_FoodSagefishDelight,
		MainHandImbue:     proto.WeaponImbue_WindfuryWeapon,
		OffHandImbue:      proto.WeaponImbue_WindfuryWeapon,
		SpellPowerBuff:    proto.SpellPowerBuff_LesserArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase3ConsumesWFFT = core.ConsumesCombo{
	Label: "Phase 3 Consumes WF/FT",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfAgility,
		DefaultPotion:     proto.Potions_ManaPotion,
		DragonBreathChili: true,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfFirepower,
		Food:              proto.Food_FoodSagefishDelight,
		MainHandImbue:     proto.WeaponImbue_WindfuryWeapon,
		OffHandImbue:      proto.WeaponImbue_FlametongueWeapon,
		SpellPowerBuff:    proto.SpellPowerBuff_LesserArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_ScrollOfStrength,
	},
}

var Phase4ConsumesWFWF = core.ConsumesCombo{
	Label: "Phase 4 Consumes WF/WF",
	Consumes: &proto.Consumes{
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		DefaultPotion:     proto.Potions_MajorManaPotion,
		DragonBreathChili: true,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfGreaterFirepower,
		Flask:             proto.Flask_FlaskOfSupremePower,
		Food:              proto.Food_FoodBlessSunfruit,
		MainHandImbue:     proto.WeaponImbue_WindfuryWeapon,
		OffHandImbue:      proto.WeaponImbue_WindfuryWeapon,
		SpellPowerBuff:    proto.SpellPowerBuff_GreaterArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
	},
}

var Phase4ConsumesWFFT = core.ConsumesCombo{
	Label: "Phase 4 Consumes WF/FT",
	Consumes: &proto.Consumes{
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		DefaultPotion:     proto.Potions_MajorManaPotion,
		DragonBreathChili: true,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfGreaterFirepower,
		Flask:             proto.Flask_FlaskOfSupremePower,
		Food:              proto.Food_FoodBlessSunfruit,
		MainHandImbue:     proto.WeaponImbue_WindfuryWeapon,
		OffHandImbue:      proto.WeaponImbue_FlametongueWeapon,
		SpellPowerBuff:    proto.SpellPowerBuff_GreaterArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
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
	proto.Stat_StatStrength,
	proto.Stat_StatAgility,
	proto.Stat_StatAttackPower,
	proto.Stat_StatMeleeHit,
	proto.Stat_StatMeleeCrit,
	proto.Stat_StatSpellPower,
}
