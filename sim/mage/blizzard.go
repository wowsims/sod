package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
)

const BlizzardRanks = 6

var BlizzardSpellId = [BlizzardRanks + 1]int32{0, 10, 6141, 8427, 10185, 10186, 10187}
var BlizzardBaseDamage = [BlizzardRanks + 1]float64{0, 200, 352, 520, 720, 936, 1192}
var BlizzardManaCost = [BlizzardRanks + 1]float64{0, 320, 520, 720, 935, 1160, 1400}
var BlizzardLevel = [BlizzardRanks + 1]int{0, 20, 28, 36, 44, 52, 60}

func (mage *Mage) registerBlizzardSpell() {
	mage.Blizzard = make([]*core.Spell, BlizzardRanks+1)

	for rank := 1; rank <= BlizzardRanks; rank++ {
		config := mage.newBlizzardSpellConfig(rank)

		if config.RequiredLevel <= int(mage.Level) {
			mage.Blizzard[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) newBlizzardSpellConfig(rank int) core.SpellConfig {
	numTicks := int32(8)
	tickLength := time.Second * 1

	spellId := BlizzardSpellId[rank]
	baseDamage := BlizzardBaseDamage[rank] / float64(numTicks)
	manaCost := BlizzardManaCost[rank]
	level := BlizzardLevel[rank]

	spellCoeff := .042

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_MageBlizzard,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: fmt.Sprintf("Blizzard (Rank %d)", rank),
			},
			NumberOfTicks:    numTicks,
			TickLength:       tickLength,
			BonusCoefficient: spellCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
				}
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	}
}
