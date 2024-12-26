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

	rogue.saberSlashTick = rogue.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_RogueSaberSlashDoT,
		ActionID:    core.ActionID{SpellID: int32(proto.RogueRune_RuneSaberSlash), Tag: 100},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

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

				// only the first stack snapshots the multiplier
				if dot.GetStacks() == 1 {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType]
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
					dot.SnapshotBaseDamage = 0
				}

				dot.SnapshotAttackerMultiplier /= rogue.saberSlashMultiplier(dot.GetStacks() - 1)
				dot.SnapshotAttackerMultiplier *= rogue.saberSlashMultiplier(dot.GetStacks())

				// each stack snapshots the AP it was applied with
				dot.SnapshotBaseDamage += 0.05 * dot.Spell.MeleeAttackPower()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
	})

	rogue.SaberSlash = rogue.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_RogueSaberSlash,
		ActionID:    core.ActionID{SpellID: int32(proto.RogueRune_RuneSaberSlash)},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       rogue.builderFlags(),
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

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: []float64{1, 1.02, 1.04, 1.06}[rogue.Talents.Aggression],
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			dot := rogue.saberSlashTick.Dot(target)
			oldMultiplier := spell.DamageMultiplier
			if dot.IsActive() {
				spell.DamageMultiplier *= rogue.saberSlashMultiplier(dot.GetStacks())
			}

			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.DamageMultiplier = oldMultiplier

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())

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

func (rogue *Rogue) saberSlashMultiplier(stacks int32) float64 {
	return []float64{1, 1.33, 1.67, 2.0}[stacks]
}
