package balance

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get caster sets included. (we use spellfire here)
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterBalanceDruid()
}

func TestBalance(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet:       core.GetGearSet("../../../ui/balance_druid/gear_sets", "phase_1"),
		OtherGearSets: []core.GearSetCombo{},
		Talents:       StandardTalents,
		Consumes:      FullConsumes,
		SpecOptions:   core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsAdaptive},
		Rotation:      core.GetAplRotation("../../../ui/balance_druid/apls", "phase_1"),

		ItemFilter: ItemFilter,
	}))
}

var StandardTalents = "50005003021"

var FullConsumes = core.ConsumesCombo{
	Label: "Full Consumes",
	Consumes: &proto.Consumes{
		Flask: proto.Flask_FlaskUnknown,
		Food:  proto.Food_FoodUnknown,
	},
}

var PlayerOptionsAdaptive = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Options: &proto.BalanceDruid_Options{
			OkfUptime: 0.2,
		},
	},
}

var ItemFilter = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeStaff,
		proto.WeaponType_WeaponTypePolearm,
	},
	ArmorType: proto.ArmorType_ArmorTypeLeather,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeIdol,
	},
}
