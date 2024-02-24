package shadow

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get caster sets included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterShadowPriest()
}

func TestShadow(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassPriest,
			Level:      25,
			Race:       proto.Race_RaceUndead,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:     Phase1Talents,
			GearSet:     core.GetGearSet("../../../ui/shadow_priest/gear_sets", "phase_1"),
			Rotation:    core.GetAplRotation("../../../ui/shadow_priest/apls", "phase_1"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassPriest,
			Level:      40,
			Race:       proto.Race_RaceUndead,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},

			Talents:     Phase2Talents,
			GearSet:     core.GetGearSet("../../../ui/shadow_priest/gear_sets", "phase_1"),
			Rotation:    core.GetAplRotation("../../../ui/shadow_priest/apls", "phase_1"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var Phase1Talents = "-20535000001"
var Phase2Talents = "--5022204002501251"

var Phase1Consumes = core.ConsumesCombo{
	Label: "Phase 1 Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion: proto.Potions_ManaPotion,
		Food:          proto.Food_FoodSmokedSagefish,
		MainHandImbue: proto.WeaponImbue_BlackfathomManaOil,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "Phase 2 Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_GreaterManaPotion,
		Food:           proto.Food_FoodSagefishDelight,
		MainHandImbue:  proto.WeaponImbue_LesserWizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_LesserArcaneElixir,
	},
}

var PlayerOptionsBasic = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			Armor: proto.ShadowPriest_Options_InnerFire,
		},
	},
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeStaff,
	},
	ArmorType: proto.ArmorType_ArmorTypeCloth,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeWand,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatIntellect,
	proto.Stat_StatSpellPower,
	proto.Stat_StatShadowPower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
}
