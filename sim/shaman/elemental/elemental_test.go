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
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassShaman,
		Level:       25,
		OtherLevels: []int32{40},

		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},

		GearSet: core.GetGearSet("../../../ui/elemental_shaman/gear_sets", "phase_1"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/elemental_shaman/gear_sets", "phase_2"),
		},

		Talents:          phase2Talents,
		Consumes:         FullConsumes,
		SpecOptions:      core.SpecOptionsCombo{Label: "Adaptive", SpecOptions: PlayerOptionsAdaptive},
		OtherSpecOptions: []core.SpecOptionsCombo{},

		Rotation: core.GetAplRotation("../../../ui/elemental_shaman/apls", "phase_1"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/elemental_shaman/apls", "phase_2"),
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

		EPReferenceStat: proto.Stat_StatSpellPower,
		StatsToWeigh: []proto.Stat{
			proto.Stat_StatIntellect,
			proto.Stat_StatSpellPower,
			proto.Stat_StatSpellHit,
			proto.Stat_StatSpellCrit,
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
				Equipment:     core.GetGearSet("../../../ui/elemental_shaman/gear_sets", "p2").GearSet,
				TalentsString: phase2Talents,
				Consumes:      FullConsumes,
				Spec:          PlayerOptionsAdaptive,
				Buffs:         core.FullIndividualBuffs,
				Rotation:      core.GetAplRotation("../../../ui/elemental_shaman/apls", "phase_2").Rotation,
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

var phase1Talents = "25003105"
var phase2Talents = "350031550002151"

var NoTotems = &proto.ShamanTotems{}
var BasicTotems = &proto.ShamanTotems{
	Earth: proto.EarthTotem_TremorTotem,
	Air:   proto.AirTotem_WindfuryTotem,
	Water: proto.WaterTotem_ManaSpringTotem,
	Fire:  proto.FireTotem_SearingTotem,
}

var PlayerOptionsAdaptive = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Options: &proto.ElementalShaman_Options{
			Shield:  proto.ShamanShield_WaterShield,
			ImbueMh: proto.ShamanImbue_RockbiterWeapon,
			ImbueOh: proto.ShamanImbue_RockbiterWeapon,
			Totems:  BasicTotems,
		},
	},
}

var FullConsumes = &proto.Consumes{
	// Flask:           proto.Flask_FlaskOfBlindingLight,
	// Food:            proto.Food_FoodBlackenedBasilisk,
	// DefaultPotion:   proto.Potions_SuperManaPotion,
	// DefaultConjured: proto.Conjured_ConjuredDarkRune,
}
