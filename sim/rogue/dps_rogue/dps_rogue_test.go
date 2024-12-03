package dpsrogue

import (
	"testing"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterDpsRogue()
}

func TestCombat(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassRogue,
			Level:      25,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     CombatDagger25Talents,
			GearSet:     core.GetGearSet("../../../ui/rogue/gear_sets", "p1_combat"),
			Rotation:    core.GetAplRotation("../../../ui/rogue/apls", "basic_strike_25"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "No Poisons", SpecOptions: DefaultCombatRogue},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassRogue,
			Level:      40,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     CombatDagger40Talents,
			GearSet:     core.GetGearSet("../../../ui/rogue/gear_sets", "p2_daggers"),
			Rotation:    core.GetAplRotation("../../../ui/rogue/apls", "mutilate"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "No Poisons", SpecOptions: DefaultCombatRogue},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassRogue,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     P4DeadlyPoisonTalents,
			GearSet:     core.GetGearSet("../../../ui/rogue/gear_sets", "p5_saber"),
			Rotation:    core.GetAplRotation("../../../ui/rogue/apls", "P5_Saber"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "No Poisons", SpecOptions: DefaultCombatRogue},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func TestAssassination(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassRogue,
			Level:      25,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     Assassination25Talents,
			GearSet:     core.GetGearSet("../../../ui/rogue/gear_sets", "p1_daggers"),
			Rotation:    core.GetAplRotation("../../../ui/rogue/apls", "mutilate"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "No Poisons", SpecOptions: DefaultAssassinationRogue},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassRogue,
			Level:      40,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     Assassination40Talents,
			GearSet:     core.GetGearSet("../../../ui/rogue/gear_sets", "p2_daggers"),
			Rotation:    core.GetAplRotation("../../../ui/rogue/apls", "mutilate"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "No Poisons", SpecOptions: DefaultAssassinationRogue},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassRogue,
			Phase:      5,
			Level:      60,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     P4DaggersTalents,
			GearSet:     core.GetGearSet("../../../ui/rogue/gear_sets", "p5_backstab"),
			Rotation:    core.GetAplRotation("../../../ui/rogue/apls", "P5_Assassination_Backstab"),
			Buffs:       core.FullBuffsPhase5,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "No Poisons", SpecOptions: DefaultAssassinationRogue},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var CombatDagger25Talents = "-025305000001"
var CombatDagger40Talents = "-0053052020550100201"
var Assassination25Talents = "0053021--05"
var Assassination40Talents = "005303103551--05"
var P4DaggersTalents = "005323105551051-023302-05"
var P4DeadlyPoisonTalents = "00532310155104-02330520000501"

var ItemFilters = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeLeather,
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeBow,
		proto.RangedWeaponType_RangedWeaponTypeCrossbow,
		proto.RangedWeaponType_RangedWeaponTypeGun,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatAttackPower,
	proto.Stat_StatAgility,
	proto.Stat_StatStrength,
	proto.Stat_StatMeleeHit,
	proto.Stat_StatMeleeCrit,
}

var DefaultAssassinationRogue = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options: DefaultDeadlyBrewOptions,
	},
}

var DefaultCombatRogue = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options: DefaultDeadlyBrewOptions,
	},
}

var DefaultDeadlyBrewOptions = &proto.RogueOptions{}

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
		AgilityElixir: proto.AgilityElixir_ElixirOfAgility,
		MainHandImbue: proto.WeaponImbue_WildStrikes,
		OffHandImbue:  proto.WeaponImbue_SolidSharpeningStone,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
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
	},
}

var Phase4Consumes = core.ConsumesCombo{
	Label: "P4-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		DragonBreathChili: true,
		Food:              proto.Food_FoodGrilledSquid,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_ElementalSharpeningStone,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
	},
}
