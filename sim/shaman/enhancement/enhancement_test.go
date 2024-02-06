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

		GearSet:       core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_1"),
		OtherGearSets: []core.GearSetCombo{},

		Talents:     phase2Talents,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "RB", SpecOptions: PlayerOptionsRBRB},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "WF", SpecOptions: PlayerOptionsWFWF},
		},

		Rotation:       core.GetAplRotation("../../../ui/enhancement_shaman/apls", "phase_1"),
		OtherRotations: []core.RotationCombo{},

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
				Consumes:      FullConsumes,
				Spec:          PlayerOptionsRBRB,
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

var PlayerOptionsWFWF = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options: enhShamWFWF,
	},
}

var PlayerOptionsRBRB = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options: enhShamRBRB,
	},
}

var enhShamWFWF = &proto.EnhancementShaman_Options{
	Shield:   proto.ShamanShield_WaterShield,
	SyncType: proto.ShamanSyncType_DelayOffhandSwings,
	ImbueMh:  proto.ShamanImbue_WindfuryWeapon,
	ImbueOh:  proto.ShamanImbue_WindfuryWeapon,
}

var enhShamRBRB = &proto.EnhancementShaman_Options{
	Shield:   proto.ShamanShield_LightningShield,
	SyncType: proto.ShamanSyncType_Auto,
	ImbueMh:  proto.ShamanImbue_RockbiterWeapon,
	ImbueOh:  proto.ShamanImbue_RockbiterWeapon,
	Totems: &proto.ShamanTotems{
		Earth: proto.EarthTotem_StrengthOfEarthTotem,
		Air:   proto.AirTotem_WindfuryTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_SearingTotem,
	},
}

var FullConsumes = &proto.Consumes{
	// DefaultConjured: proto.Conjured_ConjuredFlameCap,
}
