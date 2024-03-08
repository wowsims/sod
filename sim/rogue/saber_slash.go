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

		DamageMultiplier: []float64{1, 1.02, 1.04, 1.06}[rogue.Talents.Aggression],
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

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, applyStack bool) {
				if !applyStack {
					return
				}
				dot.SnapshotBaseDamage = 0.05 * dot.Spell.MeleeAttackPower() * float64(dot.GetStacks())
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType]
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())

				dot := spell.Dot(target)

				dot.ApplyOrRefresh(sim)
				if dot.GetStacks() < dot.MaxStacks {
					dot.AddStack(sim)
					// snapshotting only takes place when adding a stack
					dot.TakeSnapshot(sim, true)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
