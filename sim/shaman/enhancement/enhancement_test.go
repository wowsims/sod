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
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassShaman,
		Level:       25,
		OtherLevels: []int32{40},

		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},

		GearSet: core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_1"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_2"),
		},

		Talents:       phase2Talents,
		Consumes:      FullConsumesRBRB,
		OtherConsumes: []core.ConsumesCombo{},

		SpecOptions: core.SpecOptionsCombo{Label: "Sync Auto", SpecOptions: PlayerOptionsSyncAuto},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Sync Delay OH", SpecOptions: PlayerOptionsSyncDelayOH},
		},

		Rotation: core.GetAplRotation("../../../ui/enhancement_shaman/apls", "phase_1"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/enhancement_shaman/apls", "phase_2"),
		},

		ItemFilter: core.ItemFilter{
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
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:          proto.Race_RaceOrc,
				Class:         proto.Class_ClassShaman,
				Level:         40,
				Equipment:     core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_1").GearSet,
				TalentsString: phase2Talents,
				Rotation:      core.GetAplRotation("../../../ui/enhancement_shaman/apls", "phase_1").Rotation,
				Consumes:      FullConsumesRBRB.Consumes,
				Spec:          PlayerOptionsSyncAuto,
				Buffs:         core.FullIndividualBuffs,
			},
			core.FullPartyBuffs,
			core.FullRaidBuffs,
			core.FullDebuffs),
		Encounter: &proto.Encounter{
			Duration: 120,
			Targets: []*proto.Target{
				core.NewDefaultTarget(),
			},
		},
		SimOptions: core.AverageDefaultSimTestOptions,
	}

	core.RaidBenchmark(b, rsr)
}

var phase1Talents = "-5005202101"
var phase2Talents = "-5005202105023051"

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
	Shield:   proto.ShamanShield_WaterShield,
	SyncType: proto.ShamanSyncType_DelayOffhandSwings,
	Totems: &proto.ShamanTotems{
		Earth: proto.EarthTotem_StrengthOfEarthTotem,
		Air:   proto.AirTotem_WindfuryTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_SearingTotem,
	},
}

var optionsSyncAuto = &proto.EnhancementShaman_Options{
	Shield:   proto.ShamanShield_LightningShield,
	SyncType: proto.ShamanSyncType_Auto,
	Totems: &proto.ShamanTotems{
		Earth: proto.EarthTotem_StrengthOfEarthTotem,
		Air:   proto.AirTotem_WindfuryTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_MagmaTotem,
	},
}

var baseConsumables = &proto.Consumes{}

var FullConsumesRBRB = core.ConsumesCombo{
	Label: "RB/RB",
	Consumes: &proto.Consumes{
		DefaultPotion:  baseConsumables.DefaultPotion,
		Food:           baseConsumables.Food,
		Flask:          baseConsumables.Flask,
		AgilityElixir:  baseConsumables.AgilityElixir,
		StrengthBuff:   baseConsumables.StrengthBuff,
		SpellPowerBuff: baseConsumables.SpellPowerBuff,
		FirePowerBuff:  baseConsumables.FirePowerBuff,

		MainHandImbue: proto.WeaponImbue_RockbiterWeapon,
		OffHandImbue:  proto.WeaponImbue_RockbiterWeapon,
	},
}

var FullConsumesWFWF = core.ConsumesCombo{
	Label: "WF/WF",
	Consumes: &proto.Consumes{
		DefaultPotion:  baseConsumables.DefaultPotion,
		Food:           baseConsumables.Food,
		Flask:          baseConsumables.Flask,
		AgilityElixir:  baseConsumables.AgilityElixir,
		StrengthBuff:   baseConsumables.StrengthBuff,
		SpellPowerBuff: baseConsumables.SpellPowerBuff,
		FirePowerBuff:  baseConsumables.FirePowerBuff,

		MainHandImbue: proto.WeaponImbue_WindfuryWeapon,
		OffHandImbue:  proto.WeaponImbue_WindfuryWeapon,
	},
}
