package database

import (
	"github.com/wowsims/sod/sim/core/proto"
)

// Overrides for runes as needed
// Regen db with "go run ./tools/database/gen_db -outDir=assets -gen=db"
// And ensure db files are copied from assets/db into dist/sod/database
var RuneOverrides = []*proto.UIRune{
	// {Id: 415460, Name: "Engrave Chest - Burnout", Icon: "ability_mage_burnout", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassMage, RequiresLevel: 1},
	{Id: 399985, Name: "Engrave Gloves - Shadowstrike", Icon: "ability_rogue_envelopingshadows", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassRogue, RequiresLevel: 1},
	{Id: 400029, Name: "Engrave Belt - Shadowstep", Icon: "ability_rogue_shadowstep", Type: proto.ItemType_ItemTypeWaist, Class: proto.Class_ClassRogue, RequiresLevel: 30},
}
