package hunter

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterHunter()
}

// func TestBM(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
// 		{
// 			Class:      proto.Class_ClassHunter,
// 			Level:      40,
// 			Race:       proto.Race_RaceOrc,
// 			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

// 			Talents:     Phase1BMTalents,
// 			GearSet:     core.GetGearSet("../../ui/hunter/gear_sets", "phase2"),
// 			Rotation:    core.GetAplRotation("../../ui/hunter/apls", "melee.p2"),
// 			Buffs:       core.FullBuffsPhase2,
// 			Consumes:    Phase2Consumes,
// 			SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: Phase2PlayerOptions},

// 			ItemFilter:      ItemFilters,
// 			EPReferenceStat: proto.Stat_StatAttackPower,
// 			StatsToWeigh:    Stats,
// 		},
// 	}))
// }

// func TestMM(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
// 		{
// 			Class:      proto.Class_ClassHunter,
// 			Level:      25,
// 			Race:       proto.Race_RaceOrc,
// 			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

// 			Talents:     Phase1MMTalents,
// 			GearSet:     core.GetGearSet("../../ui/hunter/gear_sets", "phase2"),
// 			Rotation:    core.GetAplRotation("../../ui/hunter/apls", "melee.p2"),
// 			Buffs:       core.FullBuffsPhase1,
// 			Consumes:    Phase1Consumes,
// 			SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: Phase1PlayerOptions},

// 			ItemFilter:      ItemFilters,
// 			EPReferenceStat: proto.Stat_StatAttackPower,
// 			StatsToWeigh:    Stats,
// 		},
// 	}))
// }

func TestSV(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassHunter,
			Level:      25,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:     Phase1SVTalents,
			GearSet:     core.GetGearSet("../../ui/hunter/gear_sets", "phase1"),
			Rotation:    core.GetAplRotation("../../ui/hunter/apls", "melee.weave.p1"),
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
			GearSet:     core.GetGearSet("../../ui/hunter/gear_sets", "phase2"),
			Rotation:    core.GetAplRotation("../../ui/hunter/apls", "melee.p2"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: Phase2PlayerOptions},

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
					Race:          proto.Race_RaceOrc,
					Class:         proto.Class_ClassHunter,
					Level:         40,
					TalentsString: Phase2SVTalents,
					Equipment:     core.GetGearSet("../../ui/hunter/gear_sets", "phase2").GearSet,
					Rotation:      core.GetAplRotation("../../ui/hunter/apls", "melee.p2").Rotation,
					Consumes:      Phase2Consumes.Consumes,
					Spec:          Phase2PlayerOptions,
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

var Phase1BMTalents = "53000200501"
var Phase1MMTalents = "-050515"
var Phase1SVTalents = "--33502001101"

var Phase2SVTalents = "--335020051030315"

var Phase1Consumes = core.ConsumesCombo{
	Label: "Phase 1 Consumes",
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
	Label: "Phase 2 Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:  proto.AgilityElixir_ElixirOfAgility,
		DefaultPotion:  proto.Potions_ManaPotion,
		Food:           proto.Food_FoodDragonbreathChili,
		MainHandImbue:  proto.WeaponImbue_WildStrikes,
		OffHandImbue:   proto.WeaponImbue_SolidWeightstone,
		SpellPowerBuff: proto.SpellPowerBuff_LesserArcaneElixir,
		StrengthBuff:   proto.StrengthBuff_ElixirOfOgresStrength,
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
