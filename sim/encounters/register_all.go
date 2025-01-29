package encounters

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/encounters/naxxramas"
)

func init() {
	// TODO: Classic encounters?
	addLevel25("SoD")
	addLevel40("SoD")
	addGnomereganMechanical("SoD")
	addLevel50("SoD")
	addSunkenTempleDragonkin("SoD")
	addLevel60("SoD")
	addVaelastraszTheCorrupt("SoD")
	naxxramas.Register()
}

func AddSingleTargetBossEncounter(presetTarget *core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}
