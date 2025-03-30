package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerGarrote() {
	baseDamage := map[int32]float64{
		25: 34,
		40: 59,
		50: 74,
		60: 92,
	}[rogue.Level]

	spellID := map[int32]int32{
		25: 8631,
		40: 8633,
		50: 11289,
		60: 11290,
	}[rogue.Level]

	hasCutthroatRune := rogue.HasRune(proto.RogueRune_RuneCutthroat)

	rogue.Garrote = rogue.GetOrRegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueGarrote,
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          SpellFlagBuilder | SpellFlagCarnage | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       5,

		EnergyCost: core.EnergyCostOptions{
			Cost:   50.0 - 10*float64(rogue.Talents.DirtyDeeds),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if !rogue.IsStealthed() {
				return false
			}
			return hasCutthroatRune || !rogue.PseudoStats.InFrontOfTarget
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Garrote",
				Tag:   RogueBleedTag,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					rogue.BleedsActive[aura.Unit.UnitIndex]++
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					rogue.BleedsActive[aura.Unit.UnitIndex]--
				},
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				damage := baseDamage + dot.Spell.MeleeAttackPower()*0.03
				dot.Snapshot(target, damage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCritNoHitCounter)
			if result.Landed() {
				rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())
				spell.Dot(target).Apply(sim)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
