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
		ClassSpellMask: ClassSpellMask_RogueSaberSlash,
		ActionID:       core.ActionID{SpellID: int32(proto.RogueRune_RuneSaberSlash)},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          rogue.builderFlags(),
		EnergyCost: core.EnergyCostOptions{
			Cost:   45,
			Refund: 0.8,
		},
		MaxRange: 5,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Saber Slash - Bleed",
				Tag:       RogueBleedTag,
				Duration:  time.Second * 12,
				MaxStacks: 3,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					rogue.BleedsActive++
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					rogue.BleedsActive--
				},
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				if !isRollover {
					if dot.GetStacks() == 0 {
						dot.Snapshot(target, 0, isRollover)
					}
					return
				}

				// Each stack re-applies Saber Slash's stacking multiplier
				dot.SnapshotAttackerMultiplier /= rogue.saberSlashMultiplier(dot.GetStacks() - 1)
				dot.SnapshotAttackerMultiplier *= rogue.saberSlashMultiplier(dot.GetStacks())

				// each stack snapshots the AP it was applied with
				dot.SnapshotBaseDamage += 0.05 * dot.Spell.MeleeAttackPower()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			dot := spell.Dot(target)

			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
				return
			}

			rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())

			dot.ApplyOrRefresh(sim)
			if dot.GetStacks() < dot.MaxStacks {
				dot.AddStack(sim)
				// snapshotting only takes place when adding a stack
				dot.TakeSnapshot(sim, true)
			}
		},
	})

	rogue.RegisterAura(core.Aura{
		Label: "Saber Slash DoT Damage Amp",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells := rogue.GetSpellsMatchingClassMask(ClassSpellMask_RogueSinisterStrike | ClassSpellMask_RogueSaberSlash | ClassSpellMask_RoguePoisonedKnife)

			for _, spell := range affectedSpells {
				oldApplyEffects := spell.ApplyEffects
				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					multiplier := 1.0
					if dot := rogue.SaberSlash.Dot(target); dot.IsActive() {
						multiplier = rogue.saberSlashMultiplier(dot.GetStacks())
					}

					spell.ApplyMultiplicativeDamageBonus(multiplier)
					oldApplyEffects(sim, target, spell)
					spell.ApplyMultiplicativeDamageBonus(1 / multiplier)
				}
			}
		},
	})
}

func (rogue *Rogue) saberSlashMultiplier(stacks int32) float64 {
	return []float64{1, 1.33, 1.67, 2.0}[stacks]
}
