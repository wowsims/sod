package tank

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterFeralTankDruid()
}

func TestFeralTank(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class: proto.Class_ClassDruid,
			Level: 60,
			Phase: 5,
			Race:  proto.Race_RaceTauren,

			GearSet:     core.GetGearSet("../../../ui/feral_tank_druid/gear_sets", "phase_5"),
			Talents:     StandardTalents,
			Buffs:       core.FullBuffsPhase4,
			Consumes:    FullConsumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsDefault},
			Rotation:    core.GetAplRotation("../../../ui/feral_tank_druid/apls", "phase_5"),

			ItemFilter: core.ItemFilter{
				WeaponTypes: []proto.WeaponType{
					proto.WeaponType_WeaponTypeDagger,
					proto.WeaponType_WeaponTypeMace,
					proto.WeaponType_WeaponTypeOffHand,
					proto.WeaponType_WeaponTypeStaff,
				},
				ArmorType: proto.ArmorType_ArmorTypeLeather,
				RangedWeaponTypes: []proto.RangedWeaponType{
					proto.RangedWeaponType_RangedWeaponTypeIdol,
				},
			},
		},
	}))
}

var StandardTalents = "500005001-5050321303022151-05002"

var PlayerOptionsDefault = &proto.Player_FeralTankDruid{
	FeralTankDruid: &proto.FeralTankDruid{
		Options: &proto.FeralTankDruid_Options{
			InnervateTarget: &proto.UnitReference{}, // no Innervate
			StartingRage:    20,
		},
	},
}

var FullConsumes = core.ConsumesCombo{
	Label: "Full Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		DragonBreathChili: true,
		Flask:             proto.Flask_FlaskOfTheTitans,
		Food:              proto.Food_FoodDirgesKickChimaerokChops,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
	},
}
