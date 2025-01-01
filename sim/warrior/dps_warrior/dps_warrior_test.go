package dpswarrior

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterDpsWarrior()
}

func TestDualWieldWarrior(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassWarrior,
			Level:      40,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:     P2FuryTalents,
			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_2_dw"),
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_2_fury"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fury", SpecOptions: PlayerOptionsFury},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassWarrior,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:     P4FuryTalents,
			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_4_dw"),
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_4_fury"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fury", SpecOptions: PlayerOptionsFury},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassWarrior,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents: P4FuryTalents,
			GearSet: core.GetGearSet("../../../ui/warrior/gear_sets", "phase_5_dw_t1"),
			OtherGearSets: []core.GearSetCombo{
				core.GetGearSet("../../../ui/warrior/gear_sets", "phase_5_dw_t2"),
			},
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_5_dw"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fury", SpecOptions: PlayerOptionsFury},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassWarrior,
			Phase:      6,
			Level:      60,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:     P6FuryDWTalents,
			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_6_dw"),
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_6_dw"),
			Buffs:       core.FullBuffsPhase6,
			Consumes:    Phase6Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fury", SpecOptions: PlayerOptionsFury},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func TestTwoHandedWarrior(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassWarrior,
			Level:      50,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:     P3ArmsTalents,
			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_3_2h"),
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_3_arms"),
			Buffs:       core.FullBuffsPhase3,
			Consumes:    Phase3Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Arms", SpecOptions: PlayerOptionsArms},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassWarrior,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents: P4FuryTalents,
			GearSet: core.GetGearSet("../../../ui/warrior/gear_sets", "phase_5_2h_t1"),
			OtherGearSets: []core.GearSetCombo{
				core.GetGearSet("../../../ui/warrior/gear_sets", "phase_5_2h_t2"),
			},
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_5_2h"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Arms", SpecOptions: PlayerOptionsArms},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassWarrior,
			Phase:      6,
			Level:      60,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:     P6Fury2HTalents,
			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_6_2h"),
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_6_2h"),
			Buffs:       core.FullBuffsPhase6,
			Consumes:    Phase6Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Arms", SpecOptions: PlayerOptionsArms},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var P2ArmsTalents = "303050213525100001"
var P2FuryTalents = "-05050005405010051"
var P3ArmsTalents = "303050213520105001-0505"
var P4FuryTalents = "20305020302-05050005525010051"
var P6Fury2HTalents = "20305020332-05052005005012051"
var P6FuryDWTalents = "30315020302-55000005505010051"

var Phase1Consumes = core.ConsumesCombo{
	Label: "P1-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir: proto.AgilityElixir_ElixirOfLesserAgility,
		MainHandImbue: proto.WeaponImbue_WildStrikes,
		OffHandImbue:  proto.WeaponImbue_BlackfathomSharpeningStone,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "P2-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfAgility,
		DragonBreathChili: true,
		Food:              proto.Food_FoodSagefishDelight,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_SolidSharpeningStone,
		StrengthBuff:      proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase3Consumes = core.ConsumesCombo{
	Label: "P3-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		DragonBreathChili: true,
		Food:              proto.Food_FoodGrilledSquid,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_SolidSharpeningStone,
		StrengthBuff:      proto.StrengthBuff_ElixirOfOgresStrength,
		DefaultPotion:     proto.Potions_MightyRagePotion,
	},
}

var Phase4Consumes = core.ConsumesCombo{
	Label: "P4-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		DefaultPotion:     proto.Potions_MightyRagePotion,
		DragonBreathChili: true,
		Food:              proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_ElementalSharpeningStone,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
	},
}

var Phase6Consumes = core.ConsumesCombo{
	Label: "P6-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheHoneyBadger,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		DefaultPotion:     proto.Potions_MightyRagePotion,
		DragonBreathChili: true,
		Flask:             proto.Flask_FlaskOfMadness,
		Food:              proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_ElementalSharpeningStone,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
	},
}

var PlayerOptionsArms = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Options: warriorOptions,
	},
}

var PlayerOptionsFury = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Options: warriorOptions,
	},
}

var warriorOptions = &proto.Warrior_Options{
	StartingRage: 50,
	Shout:        proto.WarriorShout_WarriorShoutBattle,
}

var ItemFilters = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatStrength,
	proto.Stat_StatAgility,
	proto.Stat_StatAttackPower,
	proto.Stat_StatMeleeCrit,
	proto.Stat_StatMeleeHit,
}
