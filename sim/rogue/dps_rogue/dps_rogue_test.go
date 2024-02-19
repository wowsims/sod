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
			GearSet:     core.GetGearSet("../../../ui/rogue/gear_sets", "p1_daggers"),
			Rotation:    core.GetAplRotation("../../../ui/rogue/apls", "basic_strike"),
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
			GearSet:     core.GetGearSet("../../../ui/rogue/gear_sets", "p1_daggers"),
			Rotation:    core.GetAplRotation("../../../ui/rogue/apls", "basic_strike"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase2Consumes,
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
			Rotation:    core.GetAplRotation("../../../ui/rogue/apls", "basic_strike"),
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
			GearSet:     core.GetGearSet("../../../ui/rogue/gear_sets", "p1_daggers"),
			Rotation:    core.GetAplRotation("../../../ui/rogue/apls", "basic_strike"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "No Poisons", SpecOptions: DefaultAssassinationRogue},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	core.Each([]*proto.RaidSimRequest{
		{
			Raid: core.SinglePlayerRaidProto(
				&proto.Player{
					Race:          proto.Race_RaceHuman,
					Class:         proto.Class_ClassRogue,
					Level:         25,
					TalentsString: CombatDagger25Talents,
					Equipment:     core.GetGearSet("../../../ui/rogue/gear_sets", "p1_sword").GearSet,
					Rotation:      core.GetAplRotation("../../../ui/rogue/apls", "basic_strike").Rotation,
					Buffs:         core.FullIndividualBuffsPhase1,
					Consumes:      Phase1Consumes.Consumes,
					Spec:          DefaultCombatRogue,
				},
				core.FullPartyBuffs,
				core.FullRaidBuffsPhase1,
				core.FullDebuffsPhase1,
			),
			Encounter: &proto.Encounter{
				Duration: 120,
				Targets: []*proto.Target{
					core.NewDefaultTarget(25),
				},
			},
			SimOptions: core.AverageDefaultSimTestOptions,
		},
		{
			Raid: core.SinglePlayerRaidProto(
				&proto.Player{
					Race:          proto.Race_RaceHuman,
					Class:         proto.Class_ClassRogue,
					Level:         40,
					TalentsString: CombatDagger40Talents,
					Equipment:     core.GetGearSet("../../../ui/rogue/gear_sets", "p1_sword").GearSet,
					Rotation:      core.GetAplRotation("../../../ui/rogue/apls", "basic_strike").Rotation,
					Buffs:         core.FullIndividualBuffsPhase1,
					Consumes:      Phase2Consumes.Consumes,
					Spec:          DefaultCombatRogue,
				},
				core.FullPartyBuffs,
				core.FullRaidBuffsPhase1,
				core.FullDebuffsPhase1,
			),
			Encounter: &proto.Encounter{
				Duration: 120,
				Targets: []*proto.Target{
					core.NewDefaultTarget(40),
				},
			},
			SimOptions: core.AverageDefaultSimTestOptions,
		},
	}, func(rsr *proto.RaidSimRequest) { core.RaidBenchmark(b, rsr) })
}

var CombatDagger25Talents = "-025305000001"
var CombatDagger40Talents = "-02330500204501001-05"
var Assassination25Talents = "0053021--05"
var Assassination40Talents = "005303103551--05"

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

var DefaultDeadlyBrewOptions = &proto.RogueOptions{
	MhImbue: proto.RogueOptions_NoPoison,
	OhImbue: proto.RogueOptions_NoPoison,
}

var Phase1Consumes = core.ConsumesCombo{
	Label: "Phase 1 Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir: proto.AgilityElixir_ElixirOfLesserAgility,
		MainHandImbue: proto.WeaponImbue_WildStrikes,
		OffHandImbue:  proto.WeaponImbue_BlackfathomSharpeningStone,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "Phase 2 Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir: proto.AgilityElixir_ElixirOfAgility,
		MainHandImbue: proto.WeaponImbue_WildStrikes,
		OffHandImbue:  proto.WeaponImbue_SolidSharpeningStone,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
	},
}
