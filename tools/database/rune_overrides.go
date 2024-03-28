package database

import (
	"github.com/wowsims/sod/sim/core/proto"
)

// Overrides for runes as needed
// Regen db with "go run ./tools/database/gen_db -outDir=assets -gen=db"
// And ensure db files are copied from assets/db into dist/sod/database
var RuneOverrides = []*proto.UIRune{
	// Mage
	// {Id: 428738, Name: "Engrave Helm - Advanced Warding", Icon: "spell_arcane_arcaneresilience", Type: proto.ItemType_ItemTypeHead, Class: proto.Class_ClassMage, RequiresLevel: 1},
	// {Id: 428739, Name: "Engrave Helm - Deep Freeze", Icon: "ability_mage_deepfreeze", Type: proto.ItemType_ItemTypeHead, Class: proto.Class_ClassMage, RequiresLevel: 1},
	// {Id: 428885, Name: "Engrave Helm - Temporal Anomaly", Icon: "spell_fire_blueflamering", Type: proto.ItemType_ItemTypeHead, Class: proto.Class_ClassMage, RequiresLevel: 1},

	// {Id: 428878, Name: "Engrave Bracers - Balefire Bolt", Icon: "spell_fire_firebolt", Type: proto.ItemType_ItemTypeWrist, Class: proto.Class_ClassMage, RequiresLevel: 1},
	// {Id: 428861, Name: "Engrave Bracers - Displacement", Icon: "ability_hunter_displacement", Type: proto.ItemType_ItemTypeWrist, Class: proto.Class_ClassMage, RequiresLevel: 1},
	// {Id: 428741, Name: "Engrave Bracers - Molten Armor", Icon: "ability_mage_moltenarmor", Type: proto.ItemType_ItemTypeWrist, Class: proto.Class_ClassMage, RequiresLevel: 1},

	// Rogue
	{Id: 399985, Name: "Engrave Gloves - Shadowstrike", Icon: "ability_rogue_envelopingshadows", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassRogue, RequiresLevel: 1},
	{Id: 400029, Name: "Engrave Belt - Shadowstep", Icon: "ability_rogue_shadowstep", Type: proto.ItemType_ItemTypeWaist, Class: proto.Class_ClassRogue, RequiresLevel: 30},

	// Shaman
	{Id: 415768, Name: "Engrave Helm - Burn", Icon: "spell_fire_moltenblood", Type: proto.ItemType_ItemTypeHead, Class: proto.Class_ClassShaman, RequiresLevel: 1},
	{Id: 415713, Name: "Engrave Helm - Mental Dexterity", Icon: "spell_nature_mentalquickness", Type: proto.ItemType_ItemTypeHead, Class: proto.Class_ClassShaman, RequiresLevel: 1},
	{Id: 432233, Name: "Engrave Helm - Tidal Waves", Icon: "spell_shaman_tidalwaves", Type: proto.ItemType_ItemTypeHead, Class: proto.Class_ClassShaman, RequiresLevel: 1},

	{Id: 432240, Name: "Engrave Bracers - Overcharged", Icon: "spell_fire_blueflamering", Type: proto.ItemType_ItemTypeWrist, Class: proto.Class_ClassShaman, RequiresLevel: 1},
	{Id: 409954, Name: "Engrave Bracers - Riptide", Icon: "spell_nature_riptide", Type: proto.ItemType_ItemTypeWrist, Class: proto.Class_ClassShaman, RequiresLevel: 1},
	{Id: 432235, Name: "Engrave Bracers - Rolling Thunder", Icon: "spell_nature_unrelentingstorm", Type: proto.ItemType_ItemTypeWrist, Class: proto.Class_ClassShaman, RequiresLevel: 1},
	{Id: 432237, Name: "Engrave Bracers - Static Shock", Icon: "spell_shaman_staticshock", Type: proto.ItemType_ItemTypeWrist, Class: proto.Class_ClassShaman, RequiresLevel: 1},
}
