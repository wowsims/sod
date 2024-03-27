package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warrior *Warrior) registerRampage() {
	if !warrior.HasRune(proto.WarriorRune_RuneRampage) {
		return
	}
	actionID := core.ActionID{SpellID: 426940}
	auraActionID := core.ActionID{SpellID: 426942}

	warrior.RampageAura = warrior.RegisterAura(core.Aura{
		Label:     "Rampage",
		ActionID:  auraActionID,
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			bonusPerStack := warrior.ApplyStatDependencies(stats.Stats{stats.AttackPower: aura.Unit.GetStat(stats.AttackPower) * 0.02})
			warrior.AddStatsDynamic(sim, bonusPerStack.Multiply(float64(newStacks-oldStacks)))
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && sim.RandomFloat("Rampage") < 0.8 {
				warrior.RampageAura.AddStack(sim)
			}
		},
	})

	warrior.rampageValidAura = warrior.RegisterAura(core.Aura{
		Label:    "Rampage Valid Aura",
		Duration: time.Second * 5,
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Rampage Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeCrit) {
				warrior.rampageValidAura.Activate(sim)
			}
		},
	})

	warrior.Rampage = warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 20,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.rampageValidAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.rampageValidAura.Deactivate(sim)
			warrior.RampageAura.Activate(sim)
			warrior.RampageAura.AddStack(sim)
		},
	})
}
