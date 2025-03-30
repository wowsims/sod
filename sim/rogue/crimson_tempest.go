package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (rogue *Rogue) makeCrimsonTempestHitSpell() *core.Spell {
	actionID := core.ActionID{SpellID: 436611}
	procMask := core.ProcMaskMeleeMHSpecial

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagCarnage,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Crimson Tempest",
				Tag:   RogueBleedTag,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					rogue.BleedsActive[aura.Unit.UnitIndex]++
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					rogue.BleedsActive[aura.Unit.UnitIndex]--
				},
			},
			NumberOfTicks: 0,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, rogue.CrimsonTempestDamage(rogue.ComboPoints()), isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHitNoHitCounter)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.Spell = spell
				dot.NumberOfTicks = rogue.ComboPoints() + 1
				dot.Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

// TODO: Currently bugged and creates "infite loop detected" warning
func (rogue *Rogue) registerCrimsonTempestSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneCrimsonTempest) {
		return
	}

	// Must be updated to match combo points spent
	rogue.CrimsonTempestBleed = rogue.makeCrimsonTempestHitSpell()

	rogue.CrimsonTempest = rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueCrimsonTempest,
		ActionID:       core.ActionID{SpellID: 412096},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          rogue.finisherFlags(),
		MetricSplits:   6,
		MaxRange:       8,

		EnergyCost: core.EnergyCostOptions{
			Cost:   35,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				rogue.CrimsonTempestBleed.Cast(sim, aoeTarget)
			}

			rogue.SpendComboPoints(sim, spell)
		},
	})
	rogue.Finishers = append(rogue.Finishers, rogue.CrimsonTempest)
}

func (rogue *Rogue) CrimsonTempestDamage(comboPoints int32) float64 {
	tickDamageValues := []float64{0, 0.3, 0.45, 0.6, 0.75, 0.9}
	tickDamage := tickDamageValues[comboPoints] * rogue.GetStat(stats.AttackPower) / float64(comboPoints+1)
	return tickDamage
}
