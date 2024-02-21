package paladin

import "github.com/wowsims/sod/sim/core"

// Libram IDs
const (
	LibramDiscardedTenetsOfTheSilverHand = 209574
	LibramOfBenediction                  = 215435
)

func init() {
	core.NewItemEffect(LibramDiscardedTenetsOfTheSilverHand, func(agent core.Agent) {
		// character := agent.GetCharacter()
		// character.PseudoStats.
	})
}
