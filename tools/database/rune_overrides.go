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
	{Id: 442893, ClassAllowlist: []proto.Class{proto.Class_ClassDruid, proto.Class_ClassMage, proto.Class_ClassHunter}},
	// Ring - Axe Specialization
	{Id: 442876, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassPaladin, proto.Class_ClassHunter, proto.Class_ClassShaman}},
	// Ring - Dagger Specialization
	{Id: 442887, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassHunter, proto.Class_ClassRogue, proto.Class_ClassPriest, proto.Class_ClassShaman, proto.Class_ClassMage, proto.Class_ClassWarlock, proto.Class_ClassDruid}},
	// Ring - Defense Specialization
	{Id: 459312, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassPaladin, proto.Class_ClassRogue, proto.Class_ClassShaman, proto.Class_ClassWarlock, proto.Class_ClassDruid}},
	// Ring - Fire Specialization
	{Id: 442894, ClassAllowlist: []proto.Class{proto.Class_ClassShaman, proto.Class_ClassMage, proto.Class_ClassWarlock, proto.Class_ClassHunter, proto.Class_ClassPriest}},
	// Ring - Fist Weapon Specialization
	{Id: 442890, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassHunter, proto.Class_ClassRogue, proto.Class_ClassShaman, proto.Class_ClassDruid}},
	// Ring - Frost Specialization
	{Id: 442895, ClassAllowlist: []proto.Class{proto.Class_ClassShaman, proto.Class_ClassMage}},
	// Ring - Healing Specialization
	{Id: 468758, ClassAllowlist: []proto.Class{proto.Class_ClassDruid, proto.Class_ClassMage, proto.Class_ClassPaladin, proto.Class_ClassPriest, proto.Class_ClassShaman}},
	// Ring - Holy Specialization
	{Id: 442898, ClassAllowlist: []proto.Class{proto.Class_ClassPaladin, proto.Class_ClassPriest}},
	// Ring - Mace Specialization
	{Id: 442881, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior, proto.Class_ClassPaladin, proto.Class_ClassRogue, proto.Class_ClassPriest, proto.Class_ClassShaman, proto.Class_ClassDruid}},
	// Ring - Meditation Specialization
	{Id: 468762, ClassAllowlist: []proto.Class{proto.Class_ClassDruid, proto.Class_ClassHunter, proto.Class_ClassMage, proto.Class_ClassPaladin, proto.Class_ClassPriest, proto.Class_ClassShaman, proto.Class_ClassWarlock}},
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
	{Id: 458393, Name: "Engrave Gloves - Cobra Slayer", Icon: "spell_nature_guardianward", Type: proto.ItemType_ItemTypeHands, ClassAllowlist: []proto.Class{proto.Class_ClassHunter}},

	// Special should pseudo-runes

	// Druid

	// Hunter

	// Mage
	// {Id: 1220158, Name: "Soul of Winter's Grasp", Icon: "spell_holy_divinespirit", Type: proto.ItemType_ItemTypeShoulder, ClassAllowlist: []proto.Class{proto.Class_ClassMage}},

	// Paladin

	// Priest
	// {Id: 1220134, Name: "Soul of the Zealot", Icon: "spell_holy_divinespirit", Type: proto.ItemType_ItemTypeShoulder, ClassAllowlist: []proto.Class{proto.Class_ClassPriest}},

	// Rogue

	// Shaman
	// {Id: 1220232, Name: "Soul of the Windwalker", Icon: "spell_holy_divinespirit", Type: proto.ItemType_ItemTypeShoulder, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	// {Id: 1220234, Name: "Soul of the Shield Master", Icon: "spell_holy_divinespirit", Type: proto.ItemType_ItemTypeShoulder, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	// {Id: 1220236, Name: "Soul of the Totemic Protector", Icon: "spell_holy_divinespirit", Type: proto.ItemType_ItemTypeShoulder, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	// {Id: 1220238, Name: "Soul of the Shock-Absorber", Icon: "spell_holy_divinespirit", Type: proto.ItemType_ItemTypeShoulder, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	// {Id: 1220240, Name: "Soul of the Spiritual Bulwark", Icon: "spell_holy_divinespirit", Type: proto.ItemType_ItemTypeShoulder, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	// {Id: 1220242, Name: "Soul of the Maelstrombringer", Icon: "spell_holy_divinespirit", Type: proto.ItemType_ItemTypeShoulder, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	// {Id: 1220244, Name: "Soul of the Lavawalker", Icon: "spell_holy_divinespirit", Type: proto.ItemType_ItemTypeShoulder, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},

	// Warlock

	// Warrior
}
