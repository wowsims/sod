package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerSaberSlashSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneSaberSlash) {
		return
	}

	rogue.SaberSlash = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.RogueRune_RuneSaberSlash)},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,
		EnergyCost: core.EnergyCostOptions{
			Cost:   []float64{45, 42, 40}[rogue.Talents.ImprovedSinisterStrike],
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		// TODO: Fix bleed so it works properly
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Saber Slash - Bleed",
				Tag:       RogueBleedTag,
				Duration:  time.Second * 12,
				MaxStacks: 3,
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				if stacks := dot.GetStacks(); stacks > 0 {
					dot.SnapshotBaseDamage = 0.05 * dot.Spell.MeleeAttackPower() * float64(stacks)
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.SnapshotCritChance = 0
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())

				dot := spell.Dot(target)
				if !dot.IsActive() {
					dot.Apply(sim)
					dot.SetStacks(sim, 1)
					dot.TakeSnapshot(sim, false)
					return
				}

				if dot.GetStacks() < 3 {
					dot.Refresh(sim)
					dot.AddStack(sim)
					dot.TakeSnapshot(sim, false)
					return
				}

				dot.Refresh(sim)
				dot.TakeSnapshot(sim, true)
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
