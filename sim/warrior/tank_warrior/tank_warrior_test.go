package tankwarrior

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterTankWarrior()
}

func TestTankWarrior(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassWarrior,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:     P4Talents,
			GearSet:     core.GetGearSet("../../../ui/tank_warrior/gear_sets", "phase_4_tanky"),
			Rotation:    core.GetAplRotation("../../../ui/tank_warrior/apls", "phase_4"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Arms", SpecOptions: PlayerOptionsBasic},
			IsTank:      true,

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var P4Talents = "20304300302-03-55200110530201051"

var PlayerOptionsBasic = &proto.Player_TankWarrior{
	TankWarrior: &proto.TankWarrior{
		Options: warriorOptions,
	},
}

var warriorOptions = &proto.TankWarrior_Options{
	Shout:        proto.WarriorShout_WarriorShoutCommanding,
	StartingRage: 0,
}

var Phase4Consumes = core.ConsumesCombo{
	Label: "P4-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		DefaultPotion:     proto.Potions_MightyRagePotion,
		DragonBreathChili: true,
		Flask:             proto.Flask_FlaskOfTheTitans,
		Food:              proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
	},
}

var ItemFilters = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

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
	proto.Stat_StatArmor,
	proto.Stat_StatDodge,
	proto.Stat_StatParry,
	proto.Stat_StatBlockValue,
	proto.Stat_StatDefense,
}
