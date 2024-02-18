package druid

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// Totem Item IDs
const (
	IdolMindExpandingMushroom = 209576
	IdolOfWrath               = 216490
)

func init() {
	core.NewItemEffect(IdolMindExpandingMushroom, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.AddStat(stats.Spirit, 5)
	})
}
