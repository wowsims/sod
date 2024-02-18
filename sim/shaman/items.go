package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// Totem Item IDs
const (
	TotemCarvedDriftwoodIcon = 209575
	TotemInvigoratingFlame   = 215436
)

func init() {
	core.NewItemEffect(TotemCarvedDriftwoodIcon, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.AddStat(stats.MP5, 2)
	})
}
