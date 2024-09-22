package protection

import (
	"testing"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterProtectionPaladin()
}

func TestProtection(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassPaladin,
			Phase:      4,
			Level:      60,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase4ProtTalents,
			GearSet:     core.GetGearSet("../../../ui/protection_paladin/gear_sets", "p4prot"),
			Rotation:    core.GetAplRotation("../../../ui/protection_paladin/apls", "p4prot"),
			Buffs:       core.FullBuffsPhase4,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "P4 Prot", SpecOptions: PlayerOptionsSealofMartyrdom},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var Phase4ProtTalents = "-053020335001551-0500535"

var Phase4Consumes = core.ConsumesCombo{
	Label: "P4-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:     proto.Potions_MajorManaPotion,
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		Flask:             proto.Flask_FlaskOfSupremePower,
		SpellPowerBuff:    proto.SpellPowerBuff_GreaterArcaneElixir,
		DragonBreathChili: true,
		Food:              proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_ConductiveShieldCoating,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
	},
}

var PlayerOptionsSealofCommand = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options: optionsSealOfCommand,
	},
}

var PlayerOptionsSealofMartyrdom = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options: optionsSealOfMartyrdom,
	},
}

var PlayerOptionsSealofRighteousness = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options: optionsSealOfRighteousness,
	},
}

var optionsSealOfCommand = &proto.PaladinOptions{
	PrimarySeal:   proto.PaladinSeal_Command,
	RighteousFury: true,
}

var optionsSealOfMartyrdom = &proto.PaladinOptions{
	PrimarySeal:   proto.PaladinSeal_Martyrdom,
	RighteousFury: true,
}

var optionsSealOfRighteousness = &proto.PaladinOptions{
	PrimarySeal:   proto.PaladinSeal_Righteousness,
	RighteousFury: true,
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypePolearm,
		proto.WeaponType_WeaponTypeShield,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeLibram,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatHealth,
	proto.Stat_StatMana,
	proto.Stat_StatStrength,
	proto.Stat_StatStamina,
	proto.Stat_StatAgility,
	proto.Stat_StatIntellect,
	proto.Stat_StatAttackPower,
	proto.Stat_StatMeleeHit,
	proto.Stat_StatMeleeCrit,
	proto.Stat_StatMeleeHaste,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
	proto.Stat_StatSpellPower,
	proto.Stat_StatHolyPower,
	proto.Stat_StatHealingPower,
	proto.Stat_StatArmor,
	proto.Stat_StatBonusArmor,
	proto.Stat_StatDefense,
	proto.Stat_StatDodge,
	proto.Stat_StatParry,
	proto.Stat_StatBlock,
	proto.Stat_StatBlockValue,
	proto.Stat_StatFireResistance,
	proto.Stat_StatNatureResistance,
	proto.Stat_StatShadowResistance,
	proto.Stat_StatFrostResistance,
	proto.Stat_StatArcaneResistance,
}
