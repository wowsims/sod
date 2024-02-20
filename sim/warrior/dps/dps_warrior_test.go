package dps

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterDpsWarrior()
}

func TestFury(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassWarrior,
			Level:      40,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:     P2FuryTalents,
			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_2_dw"),
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_2"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fury", SpecOptions: PlayerOptionsFury},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

// func TestArms(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
// 		{
// 			Class:      proto.Class_ClassWarrior,
// 			Level:      25,
// 			Race:       proto.Race_RaceOrc,
// 			OtherRaces: []proto.Race{proto.Race_RaceHuman},

// 			Talents:     P1ArmsTalents,
// 			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_1"),
// 			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_1"),
// 			Buffs:       core.FullBuffsPhase1,
// 			Consumes:    Phase1Consumes,
// 			SpecOptions: core.SpecOptionsCombo{Label: "Arms", SpecOptions: PlayerOptionsArms},

// 			ItemFilter:      ItemFilters,
// 			EPReferenceStat: proto.Stat_StatAttackPower,
// 			StatsToWeigh:    Stats,
// 		},
// 		{
// 			Class:      proto.Class_ClassWarrior,
// 			Level:      25,
// 			Race:       proto.Race_RaceOrc,
// 			OtherRaces: []proto.Race{proto.Race_RaceHuman},

// 			Talents:     P1ArmsTalents,
// 			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_1_dw"),
// 			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_1"),
// 			Buffs:       core.FullBuffsPhase1,
// 			Consumes:    Phase1Consumes,
// 			SpecOptions: core.SpecOptionsCombo{Label: "Arms", SpecOptions: PlayerOptionsArms},

// 			ItemFilter:      ItemFilters,
// 			EPReferenceStat: proto.Stat_StatAttackPower,
// 			StatsToWeigh:    Stats,
// 		},
// 	}))
// }

func BenchmarkSimulate(b *testing.B) {
	core.Each([]*proto.RaidSimRequest{
		{
			Raid: core.SinglePlayerRaidProto(
				&proto.Player{
					Race:          proto.Race_RaceOrc,
					Class:         proto.Class_ClassWarrior,
					Level:         40,
					Equipment:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_2").GearSet,
					Rotation:      core.GetAplRotation("../../../ui/warrior/apls", "phase_2").Rotation,
					Consumes:      Phase2Consumes.Consumes,
					Spec:          PlayerOptionsFury,
					TalentsString: P2FuryTalents,
					Buffs:         core.FullIndividualBuffsPhase2,
				},
				core.FullPartyBuffs,
				core.FullRaidBuffsPhase2,
				core.FullDebuffsPhase2,
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

var P1ArmsTalents = "303220203-01"

var P2FuryTalents = "-05050005405010051"

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
		Food:          proto.Food_FoodDragonbreathChili,
		MainHandImbue: proto.WeaponImbue_WildStrikes,
		OffHandImbue:  proto.WeaponImbue_SolidSharpeningStone,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
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
	StartingRage:    50,
	UseRecklessness: true,
	Shout:           proto.WarriorShout_WarriorShoutBattle,
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
