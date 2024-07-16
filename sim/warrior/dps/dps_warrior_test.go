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
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_2_fury"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fury", SpecOptions: PlayerOptionsFury},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassWarrior,
			Level:      60,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:     P4FuryTalents,
			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_4_dw"),
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_4_fury"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fury", SpecOptions: PlayerOptionsFury},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func TestArms(t *testing.T) {
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
	}))
}

func BenchmarkSimulate(b *testing.B) {
	core.Each([]*proto.RaidSimRequest{
		{
			Raid: core.SinglePlayerRaidProto(
				&proto.Player{
					Race:               proto.Race_RaceOrc,
					Class:              proto.Class_ClassWarrior,
					Level:              40,
					Equipment:          core.GetGearSet("../../../ui/warrior/gear_sets", "phase_2_2h").GearSet,
					Rotation:           core.GetAplRotation("../../../ui/warrior/apls", "phase_2_arms").Rotation,
					Consumes:           Phase2Consumes.Consumes,
					Spec:               PlayerOptionsFury,
					TalentsString:      P2FuryTalents,
					Buffs:              core.FullIndividualBuffsPhase2,
					DistanceFromTarget: 0,
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

var P2ArmsTalents = "303050213525100001"
var P2FuryTalents = "-05050005405010051"
var P3ArmsTalents = "303050213520105001-0505"
var P4FuryTalents = "20305020302-05050005525010051"

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
		AgilityElixir:     proto.AgilityElixir_ElixirOfAgility,
		DragonBreathChili: true,
		Food:              proto.Food_FoodSagefishDelight,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_SolidSharpeningStone,
		StrengthBuff:      proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase3Consumes = core.ConsumesCombo{
	Label: "Phase 3 Consumes",
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
	Label: "Phase 4 Consumes",
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
