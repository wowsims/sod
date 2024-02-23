package retribution

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterRetributionPaladin()
}

func TestRetribution(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassPaladin,
			Level:      25,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase1RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p1ret"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p1ret"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Seal of Command Ret", SpecOptions: PlayerOptionsSealofCommand},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPaladin,
			Level:      40,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase2RetTalents,
			GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p2ret"),
			Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p2ret"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Seal of Command Ret", SpecOptions: PlayerOptionsSealofCommand},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		// {
		// 	Class:      proto.Class_ClassPaladin,
		// 	Level:      40,
		// 	Race:       proto.Race_RaceHuman,
		// 	OtherRaces: []proto.Race{proto.Race_RaceDwarf},

		// 	Talents:     Phase2ShockadinTalents,
		// 	GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p2shockadin"),
		// 	Rotation:    core.GetAplRotation("../../../ui/retribution_paladin/apls", "p2ret"),
		// 	Buffs:       core.FullBuffsPhase2,
		// 	Consumes:    Phase2Consumes,
		// 	SpecOptions: core.SpecOptionsCombo{Label: "Seal of Martyrdom Shockadin", SpecOptions: PlayerOptionsSealofMartyrdom},

		// 	ItemFilter:      ItemFilters,
		// 	EPReferenceStat: proto.Stat_StatAttackPower,
		// 	StatsToWeigh:    Stats,
		// },
	}))
}

func BenchmarkSimulate(b *testing.B) {
	core.Each([]*proto.RaidSimRequest{
		{
			Raid: core.SinglePlayerRaidProto(
				&proto.Player{
					Race:          proto.Race_RaceHuman,
					Class:         proto.Class_ClassPaladin,
					Level:         25,
					TalentsString: Phase1RetTalents,
					Equipment:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p1ret").GearSet,
					Rotation:      core.GetAplRotation("../../../ui/retribution_paladin/apls", "p1ret").Rotation,
					Consumes:      Phase1Consumes.Consumes,
					Spec:          PlayerOptionsSealofCommand,
					Buffs:         core.FullIndividualBuffsPhase1,
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
					Class:         proto.Class_ClassPaladin,
					Level:         40,
					TalentsString: Phase2RetTalents,
					Equipment:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p2ret").GearSet,
					Rotation:      core.GetAplRotation("../../../ui/retribution_paladin/apls", "p2ret").Rotation,
					Consumes:      Phase2Consumes.Consumes,
					Spec:          PlayerOptionsSealofCommand,
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
		// {
		// 	Raid: core.SinglePlayerRaidProto(
		// 		&proto.Player{
		// 			Race:          proto.Race_RaceHuman,
		// 			Class:         proto.Class_ClassPaladin,
		// 			Level:         40,
		// 			TalentsString: Phase2RetTalents,
		// 			Equipment:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p2shockadin").GearSet,
		// 			Rotation:      core.GetAplRotation("../../../ui/retribution_paladin/apls", "p2ret").Rotation,
		// 			Consumes:      Phase2Consumes.Consumes,
		// 			Spec:          PlayerOptionsSealofMartyrdom,
		// 			Buffs:         core.FullIndividualBuffsPhase2,
		// 		},
		// 		core.FullPartyBuffs,
		// 		core.FullRaidBuffsPhase2,
		// 		core.FullDebuffsPhase2,
		// 	),
		// 	Encounter: &proto.Encounter{
		// 		Duration: 120,
		// 		Targets: []*proto.Target{
		// 			core.NewDefaultTarget(40),
		// 		},
		// 	},
		// 	SimOptions: core.AverageDefaultSimTestOptions,
		// },
	}, func(rsr *proto.RaidSimRequest) { core.RaidBenchmark(b, rsr) })
}

var Phase1RetTalents = "--05230051"
var Phase2RetTalents = "--532300512003151"
var Phase2ShockadinTalents = "55050100521151--"

var Phase1Consumes = core.ConsumesCombo{
	Label: "Phase 1 Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir: proto.AgilityElixir_ElixirOfLesserAgility,
		DefaultPotion: proto.Potions_ManaPotion,
		FirePowerBuff: proto.FirePowerBuff_ElixirOfFirepower,
		MainHandImbue: proto.WeaponImbue_WildStrikes,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "Phase 2 Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:  proto.AgilityElixir_ElixirOfAgility,
		DefaultPotion:  proto.Potions_ManaPotion,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfFirepower,
		Food:           proto.Food_FoodDragonbreathChili,
		MainHandImbue:  proto.WeaponImbue_WindfuryWeapon,
		SpellPowerBuff: proto.SpellPowerBuff_LesserArcaneElixir,
		StrengthBuff:   proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var PlayerOptionsSealofCommand = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: optionsSealOfCommand,
	},
}

var PlayerOptionsSealofMartyrdom = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: optionsSealOfMartyrdom,
	},
}

var optionsSealOfCommand = &proto.RetributionPaladin_Options{
	PrimarySeal: proto.PaladinSeal_Command,
}

var optionsSealOfMartyrdom = &proto.RetributionPaladin_Options{
	PrimarySeal: proto.PaladinSeal_Martyrdom,
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypePolearm,
		proto.WeaponType_WeaponTypeShield,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeLibram,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatStrength,
	proto.Stat_StatAgility,
	proto.Stat_StatAttackPower,
	proto.Stat_StatMeleeHit,
	proto.Stat_StatMeleeCrit,
	proto.Stat_StatSpellPower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
}
