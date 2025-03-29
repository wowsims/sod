package rogue

import "github.com/wowsims/sod/sim/core"

// Damaging finishing moves have a 20% chance per combo point to restore 20 energy.
func (rogue *Rogue) ApplyFallenRegalityRogueBonus(aura *core.Aura) {
	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 1232184})
	aura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
		rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
			if spell.ProcMask != core.ProcMaskEmpty && sim.Proc(0.20*float64(comboPoints), "Fallen Regality Proc") {
				rogue.AddEnergy(sim, 20, energyMetrics)
			}
		})
	})
}
