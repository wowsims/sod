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
	}))
}

func BenchmarkSimulate(b *testing.B) {
	core.Each([]*proto.RaidSimRequest{
		{
			Raid: core.SinglePlayerRaidProto(
				&proto.Player{
					Race:          proto.Race_RaceTroll,
					Class:         proto.Class_ClassShaman,
					Level:         25,
					TalentsString: Phase1Talents,
					Equipment:     core.GetGearSet("../../../ui/elemental_shaman/gear_sets", "phase_1").GearSet,
					Rotation:      core.GetAplRotation("../../../ui/elemental_shaman/apls", "phase_1").Rotation,
					Buffs:         core.FullIndividualBuffsPhase1,
					Consumes:      Phase1Consumes.Consumes,
					Spec:          PlayerOptionsAdaptive,
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
					Race:          proto.Race_RaceTroll,
					Class:         proto.Class_ClassShaman,
					Level:         40,
					TalentsString: Phase2Talents,
					Equipment:     core.GetGearSet("../../../ui/elemental_shaman/gear_sets", "phase_2").GearSet,
					Rotation:      core.GetAplRotation("../../../ui/elemental_shaman/apls", "phase_2").Rotation,
					Buffs:         core.FullIndividualBuffsPhase2,
					Consumes:      Phase2Consumes.Consumes,
					Spec:          PlayerOptionsAdaptive,
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

var Phase1Talents = "25003105"
var Phase2Talents = "550031550000151"

var PlayerOptionsAdaptive = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Options: &proto.ElementalShaman_Options{
			Shield: proto.ShamanShield_LightningShield,
		},
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
