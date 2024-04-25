package database

import (
	"github.com/wowsims/sod/sim/core/proto"
)

// Overrides for runes as needed
// Regen db with "go run ./tools/database/gen_db -outDir=assets -gen=db"
// And ensure db files are copied from assets/db into dist/sod/database
var RuneOverrides []*proto.UIRune

// Remove runes as you implement them.
var UnimplementedRuneOverrides = []int32{
	// Paladin
	429133, // Improved Sanctuary
	428909, // Light's Grace
}
