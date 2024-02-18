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
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class: proto.Class_ClassWarlock,
			Level: 25,
			Race:  proto.Race_RaceOrc,

			Talents:     Phase1AfflictionTalents,
			GearSet:     core.GetGearSet("../../../ui/tank_warlock/gear_sets", "affi.tank"),
			Rotation:    core.GetAplRotation("../../../ui/tank_warlock/apls", "affi.tank"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Affliction Warlock", SpecOptions: DefaultAfflictionWarlock},
			OtherSpecOptions: []core.SpecOptionsCombo{
				{Label: "AffItemSwap", SpecOptions: afflictionItemSwap},
			},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

// func TestDemonology(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
// 		{
// 			Class: proto.Class_ClassWarlock,
// 			Level: 25,
// 			Race:  proto.Race_RaceOrc,

// 			GearSet:     core.GetGearSet("../../../ui/tank_warlock/gear_sets", "destro.tank"),
// 			Talents:     Phase1DemonologyTalents,
// 			Buffs:       core.FullBuffsPhase1,
// 			Consumes:    Phase1Consumes,
// 			SpecOptions: core.SpecOptionsCombo{Label: "Demonology Warlock", SpecOptions: DefaultDemonologyWarlock},

// 			ItemFilter:      ItemFilters,
// 			EPReferenceStat: proto.Stat_StatSpellPower,
// 			StatsToWeigh:    Stats,
// 		},
// 	}))
// }

func TestDestruction(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class: proto.Class_ClassWarlock,
			Level: 25,
			Race:  proto.Race_RaceOrc,

			Talents:     Phase1DestructionTalents,
			GearSet:     core.GetGearSet("../../../ui/tank_warlock/gear_sets", "destro.tank"),
			Rotation:    core.GetAplRotation("../../../ui/tank_warlock/apls", "destro.tank"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Destruction Warlock", SpecOptions: DefaultDestroWarlock},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var Phase1AfflictionTalents = "05002-005"
var Phase1DemonologyTalents = "05002-005"
var Phase1DestructionTalents = "05002-005"

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

var Phase1Consumes = core.ConsumesCombo{
	Label: "Phase 1 Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir: proto.AgilityElixir_ElixirOfLesserAgility,
		DefaultPotion: proto.Potions_ManaPotion,
		FirePowerBuff: proto.FirePowerBuff_ElixirOfFirepower,
		Food:          proto.Food_FoodSmokedSagefish,
		MainHandImbue: proto.WeaponImbue_BlackfathomManaOil,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "Phase 2 Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_ManaPotion,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfFirepower,
		Food:           proto.Food_FoodSagefishDelight,
		MainHandImbue:  proto.WeaponImbue_LesserWizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_LesserArcaneElixir,
	},
}

var ItemFilters = core.ItemFilter{
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

var Stats = []proto.Stat{
	proto.Stat_StatIntellect,
	proto.Stat_StatSpellPower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
}
