package database

import (
	"github.com/wowsims/sod/sim/core/proto"
)

// Overrides for runes as needed
// Regen db with "go run ./tools/database/gen_db -outDir=assets -gen=db"
// And ensure db files are copied from assets/db into dist/sod/database
var RuneOverrides = []*proto.UIRune{}

// Remove runes as you implement them.
var UnimplementedRuneOverrides = []int32{
	// Hunter
	415428, // Catlike Reflexes
	415405, // Rapid Killing
	428717, // T.N.T.

	// Paladin
	429142, // Fanaticism
	429133, // Improved Sanctuary
	428909, // Light's Grace
	429139, // Wrath
	407632, // Hammer of the Righteous
	429152, // Improved Hammer of Wrath
	429144, // Purifying power

	// Priest
	431622, // Divine Aegis
	402789, // Eye of the Void
	413251, // Pain and Suffering
	431670, // Despair
	431664, // Surge of Light
	431681, // Void Zone

	// Rogue
	432273, // Unfair Advantage
}
