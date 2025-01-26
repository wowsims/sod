package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerSliceAndDice() {
	hasteBonusByRank := map[int32]float64{
		25: 0.20,
		40: 0.20,
		50: 0.30,
		60: 0.30,
	}[rogue.Level]

	spellID := map[int32]int32{
		25: 5171,
		40: 5171,
		50: 6774,
		60: 6774,
	}[rogue.Level]

	actionID := core.ActionID{SpellID: spellID}

	durationMultiplier := []float64{1, 1.15, 1.3, 1.45}[rogue.Talents.ImprovedSliceAndDice]

	rogue.sliceAndDiceDurations = [6]time.Duration{
		0,
		time.Duration(float64(time.Second*9) * durationMultiplier),
		time.Duration(float64(time.Second*12) * durationMultiplier),
		time.Duration(float64(time.Second*15) * durationMultiplier),
		time.Duration(float64(time.Second*18) * durationMultiplier),
		time.Duration(float64(time.Second*21) * durationMultiplier),
	}

	hasteBonus := 1 + hasteBonusByRank
	inverseHasteBonus := 1.0 / hasteBonus

	rogue.SliceAndDiceAura = rogue.RegisterAura(core.Aura{
		Label:    "Slice and Dice",
		ActionID: actionID,
		// This will be overridden on cast, but set a non-zero default so it doesn't crash when used in APL prepull
		Duration: rogue.sliceAndDiceDurations[5],
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, hasteBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, inverseHasteBonus)
		},
	})

	rogue.SliceAndDice = rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueSliceandDice,
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		MetricSplits:   6,

		EnergyCost: core.EnergyCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.SliceAndDiceAura.Duration = rogue.sliceAndDiceDurations[rogue.ComboPoints()]
			rogue.SliceAndDiceAura.Activate(sim)
			rogue.SpendComboPoints(sim, spell)
		},
	})
	rogue.Finishers = append(rogue.Finishers, rogue.SliceAndDice)
}
