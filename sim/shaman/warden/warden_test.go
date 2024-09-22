package warden

import (
	"testing"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterWardenShaman()
}

func TestWardenShaman(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassShaman,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     Phase4Talents,
			GearSet:     core.GetGearSet("../../../ui/warden_shaman/gear_sets", "phase_4_enh_tank"),
			Rotation:    core.GetAplRotation("../../../ui/warden_shaman/apls", "phase_4_enh_tank"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsBasic},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var PlayerOptionsBasic = &proto.Player_WardenShaman{
	WardenShaman: &proto.WardenShaman{
		Options: &proto.WardenShaman_Options{},
	},
}

var Phase4Talents = "05033150003-0505032015003151"

var Phase4Consumes = core.ConsumesCombo{
	Label: "Phase 4 Consumes",
	Consumes: &proto.Consumes{
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		DefaultPotion:     proto.Potions_MajorManaPotion,
		DragonBreathChili: true,
		FirePowerBuff:     proto.FirePowerBuff_ElixirOfGreaterFirepower,
		Flask:             proto.Flask_FlaskOfTheTitans,
		Food:              proto.Food_FoodBlessSunfruit,
		MainHandImbue:     proto.WeaponImbue_WindfuryWeapon,
		OffHandImbue:      proto.WeaponImbue_ConductiveShieldCoating,
		SpellPowerBuff:    proto.SpellPowerBuff_GreaterArcaneElixir,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
	},
}

var ItemFilters = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeMail,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
		proto.WeaponType_WeaponTypeShield,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatStrength,
	proto.Stat_StatAttackPower,
	proto.Stat_StatSpellPower,
	proto.Stat_StatArmor,
	proto.Stat_StatDodge,
	proto.Stat_StatParry,
	proto.Stat_StatBlockValue,
	proto.Stat_StatDefense,
}
