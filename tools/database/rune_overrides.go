package database

import (
	"github.com/wowsims/sod/sim/core/proto"
)

// Overrides for runes as needed
// Regen db with "go run ./tools/database/gen_db -outDir=assets -gen=db"
// And ensure db files are copied from assets/db into dist/sod/database
var RuneOverrides = []*proto.UIRune{
	// Ring rune tooltips lack the relevant class restrictions so manually override the class allowlists
	// Ring - Arcane Specialization
	{Id: 442893, ClassAllowlist: []proto.Class{proto.Class_ClassDruid, proto.Class_ClassMage}},
	// Ring - Axe Specialization
	{Id: 442876, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassPaladin, proto.Class_ClassHunter, proto.Class_ClassShaman}},
	// Ring - Dagger Specialization
	{Id: 442887, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassHunter, proto.Class_ClassRogue, proto.Class_ClassPriest, proto.Class_ClassShaman, proto.Class_ClassMage, proto.Class_ClassWarlock, proto.Class_ClassDruid}},
	// Ring - Defense Specialization
	{Id: 459312, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassPaladin, proto.Class_ClassRogue, proto.Class_ClassShaman, proto.Class_ClassWarlock, proto.Class_ClassDruid}},
	// Ring - Fire Specialization
	{Id: 442894, ClassAllowlist: []proto.Class{proto.Class_ClassShaman, proto.Class_ClassMage, proto.Class_ClassWarlock}},
	// Ring - Fist Weapon Specialization
	{Id: 442890, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassHunter, proto.Class_ClassRogue, proto.Class_ClassShaman, proto.Class_ClassDruid}},
	// Ring - Frost Specialization
	{Id: 442895, ClassAllowlist: []proto.Class{proto.Class_ClassShaman, proto.Class_ClassMage}},
	// Ring - Holy Specialization
	{Id: 442898, ClassAllowlist: []proto.Class{proto.Class_ClassPaladin, proto.Class_ClassPriest}},
	// Ring - Mace Specialization
	{Id: 442881, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassPaladin, proto.Class_ClassRogue, proto.Class_ClassPriest, proto.Class_ClassShaman, proto.Class_ClassDruid}},
	// Ring - Nature Specialization
	{Id: 442896, ClassAllowlist: []proto.Class{proto.Class_ClassRogue, proto.Class_ClassShaman, proto.Class_ClassDruid}},
	// Ring - Pole Weapon Specialization
	{Id: 442892, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassPaladin, proto.Class_ClassHunter, proto.Class_ClassPriest, proto.Class_ClassShaman, proto.Class_ClassMage, proto.Class_ClassWarlock, proto.Class_ClassDruid}},
	// Ring - Ranged Weapon Specialization
	{Id: 442891, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassHunter, proto.Class_ClassRogue}},
	// Ring - Shadow Specialization
	{Id: 442897, ClassAllowlist: []proto.Class{proto.Class_ClassPriest, proto.Class_ClassWarlock}},
	// Ring - Sword Specialization
	{Id: 442813, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassPaladin, proto.Class_ClassHunter, proto.Class_ClassRogue, proto.Class_ClassMage, proto.Class_ClassWarlock}},

	// Hunter
	// As of 2024-06-13 Cobra Slayer is being missed by the scraper because the rune engraving ability is missing "Engrave Rune" in the name
	{Id: 458393, Name: "Engrave Chest - Cobra Slayer", Icon: "spell_nature_guardianward", Type: proto.ItemType_ItemTypeChest, ClassAllowlist: []proto.Class{proto.Class_ClassHunter}},
	// Warlock
	// TODO: These runes haven't been updated by wowhead yet but were updated on 2024-07-03
	// Cloak - Soul Siphon
	{Id: 403511, Name: "Engrave Cloak - Soul Siphon", Icon: "spell_shadow_lifedrain02", Type: proto.ItemType_ItemTypeBack, ClassAllowlist: []proto.Class{proto.Class_ClassWarlock}},

	// Bracers - Incinerate
	{Id: 412758, Name: "Engrave Bracers - Incinerate", Icon: "spell_fire_burnout", Type: proto.ItemType_ItemTypeWrist, ClassAllowlist: []proto.Class{proto.Class_ClassWarlock}},

	// Boots - Mark of Chaos
	{Id: 440892, Name: "Engrave Boots - Mark of Chaos", Icon: "spell_shadow_unstableaffliction_1", Type: proto.ItemType_ItemTypeFeet, ClassAllowlist: []proto.Class{proto.Class_ClassWarlock}},
}

// Remove runes as you implement them.
var UnimplementedRuneOverrides = []int32{
	// Druid

	// Hunter

	// Mage

	// Paladin
	440658, // Cloak - Shield of Righteousness
	440666, // Cloak - Vindicator

	// Priest

	// Shaman

	// Warlock
	440882, // Cloak - Infernal Armor

	// Warrior
}
