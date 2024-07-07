package guardians

import "github.com/wowsims/sod/sim/core"

func ConstructGuardians(character *core.Character) {
	constructEmeralDragonWhelps(character)
	constructEskhandar(character)
	constructCoreHound(character)
}
