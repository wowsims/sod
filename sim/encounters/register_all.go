package encounters

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/encounters/naxxramas"
)

func init() {
	naxxramas.Register()
	addVaelastraszTheCorrupt("SoD")
	addLevel60("SoD")
	addSunkenTempleDragonkin("SoD")
	addLevel50("SoD")
	addGnomereganMechanical("SoD")
	addLevel40("SoD")
	addLevel25("SoD")
}

func AddSingleTargetBossEncounter(presetTarget *core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}
