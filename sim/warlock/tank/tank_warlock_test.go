package tank

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterTankWarlock()
}

func TestAffliction(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GetGearSet("../../../ui/tank_warlock/gear_sets", "affi.tank"),
		Talents:     AfflictionTalents,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Affliction Warlock", SpecOptions: DefaultAfflictionWarlock},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "AffItemSwap", SpecOptions: afflictionItemSwap},
		},
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/tank_warlock/apls", "affi.tank"),
		},

		ItemFilter: ItemFilter,
	}))
}

func TestDemonology(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GetGearSet("../../../ui/tank_warlock/gear_sets", "destro.tank"),
		Talents:     DemonologyTalents,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Demonology Warlock", SpecOptions: DefaultDemonologyWarlock},
		ItemFilter:  ItemFilter,
	}))
}

func TestDestruction(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GetGearSet("../../../ui/tank_warlock/gear_sets", "destro.tank"),
		Talents:     DestructionTalents,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Destruction Warlock", SpecOptions: DefaultDestroWarlock},
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/tank_warlock/apls", "destro.tank"),
		},
		ItemFilter: ItemFilter,
	}))
}

var ItemFilter = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeDagger,
	},
	HandTypes: []proto.HandType{
		proto.HandType_HandTypeOffHand,
	},
	ArmorType: proto.ArmorType_ArmorTypeCloth,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeWand,
	},
}

var AfflictionTalents = "05002-005"
var DemonologyTalents = "05002-005"
var DestructionTalents = "05002-005"

var defaultDestroOptions = &proto.WarlockOptions{
	Armor:       proto.WarlockOptions_DemonArmor,
	Summon:      proto.WarlockOptions_Imp,
	WeaponImbue: proto.WarlockOptions_NoWeaponImbue,
}

var DefaultDestroWarlock = &proto.Player_TankWarlock{
	TankWarlock: &proto.TankWarlock{
		Options: defaultDestroOptions,
	},
}

// ---------------------------------------
var DefaultAfflictionWarlock = &proto.Player_TankWarlock{
	TankWarlock: &proto.TankWarlock{
		Options: defaultAfflictionOptions,
	},
}

var afflictionItemSwap = &proto.Player_TankWarlock{
	TankWarlock: &proto.TankWarlock{
		Options: defaultAfflictionOptions,
	},
}

var defaultAfflictionOptions = &proto.WarlockOptions{
	Armor:       proto.WarlockOptions_DemonArmor,
	Summon:      proto.WarlockOptions_Imp,
	WeaponImbue: proto.WarlockOptions_NoWeaponImbue,
}

// ---------------------------------------
var DefaultDemonologyWarlock = &proto.Player_TankWarlock{
	TankWarlock: &proto.TankWarlock{
		Options: defaultDemonologyOptions,
	},
}

var defaultDemonologyOptions = &proto.WarlockOptions{
	Armor:       proto.WarlockOptions_DemonArmor,
	Summon:      proto.WarlockOptions_Imp,
	WeaponImbue: proto.WarlockOptions_NoWeaponImbue,
}

// ---------------------------------------------------------

var FullConsumes = core.ConsumesCombo{
	Label: "Full Consumes",
	Consumes: &proto.Consumes{
		Flask:         proto.Flask_FlaskOfSupremePower,
		DefaultPotion: proto.Potions_ManaPotion,
		Food:          proto.Food_FoodBlessSunfruit,
	},
}
