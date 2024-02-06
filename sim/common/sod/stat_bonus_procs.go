package sod

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

type ProcStatBonusEffect struct {
	Name       string
	ID         int32
	AuraID     int32
	Bonus      stats.Stats
	Duration   time.Duration
	Callback   core.AuraCallback
	ProcMask   core.ProcMask
	Outcome    core.HitOutcome
	Harmful    bool
	ProcChance float64
	PPM        float64
	ICD        time.Duration

	// For ignoring a hardcoded spell.
	IgnoreSpellID int32
}

func newProcStatBonusEffect(config ProcStatBonusEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		procID := core.ActionID{SpellID: config.AuraID}
		if procID.IsEmptyAction() {
			procID = core.ActionID{ItemID: config.ID}
		}
		procAura := character.NewTemporaryStatsAura(config.Name+" Proc", procID, config.Bonus, config.Duration)

		handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			procAura.Activate(sim)
		}
		if config.IgnoreSpellID != 0 {
			ignoreSpellID := config.IgnoreSpellID
			handler = func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if !spell.IsSpellAction(ignoreSpellID) {
					procAura.Activate(sim)
				}
			}
		}

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{ItemID: config.ID},
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: config.ProcChance,
			PPM:        config.PPM,
			ICD:        config.ICD,
			Handler:    handler,
		})
		procAura.Icd = triggerAura.Icd
	})
}
