package rogue

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

/**
Instant Poison: 20% proc chance
25: 22 +/- 3 damage, 8679 ID, 40 charges
40: 50 +/- 6 damage, 8688 ID, 70 charges
50: 76 +/- 9 damage, 11338 ID, 85 charges
60: 130 =/- 18 damage, 11340 ID, 115 charges

Deadly Poison: 30% proc chance, 5 stacks
25: 36 damage, 2823 ID, 60 charges (Deadly Brew only)
40: 52 damage, 2824 ID, 75 charges
50: 80 damage, 11355 ID, 90 charges
60: 108 damage, 11356 ID, 105 charges (Rank 4, Rank 5 is by book)

Wound Poison: 30% proc chance, 5 stacks
25: x damage, x ID (none, first rank is level 32)
40: -75 healing, 11325 ID, 75 charges (Rank 2)
50: -105 healing, 13226 ID, 90 charges (Rank 3)
60: -135 healing, 13227 ID, 105 charges (Rank 4)
*/

// TODO: Add charges to poisons (not deadly brew)

type PoisonProcSource int

const (
	NormalProc PoisonProcSource = iota
	ShivProc
	DeadlyBrewProc
)

func (rogue *Rogue) GetInstantPoisonProcChance() float64 {
	return (0.2 + rogue.improvedPoisons()) * (1 + rogue.instantPoisonProcChanceBonus)
}

func (rogue *Rogue) GetDeadlyPoisonProcChance() float64 {
	return 0.3 + rogue.improvedPoisons()
}

func (rogue *Rogue) GetWoundPoisonProcChance() float64 {
	return 0.3 + rogue.improvedPoisons()
}

func (rogue *Rogue) improvedPoisons() float64 {
	return []float64{0, 0.02, 0.04, 0.06, 0.08, 0.1}[rogue.Talents.ImprovedPoisons]
}

func (rogue *Rogue) getPoisonDamageMultiplier() float64 {
	return []float64{1, 1.04, 1.08, 1.12, 1.16, 1.2}[rogue.Talents.VilePoisons]
}

///////////////////////////////////////////////////////////////////////////
//                               Apply Poisons
///////////////////////////////////////////////////////////////////////////

func (rogue *Rogue) applyPoisons() {
	rogue.applyDeadlyPoison()
	rogue.applyInstantPoison()
	rogue.applyWoundPoison()
}

// Apply Deadly Brew Instant Poison procs
func (rogue *Rogue) applyDeadlyBrewInstant() {
	// apply IP from all weapons w/o IP, DP, or WP applied
	procMask := core.ProcMaskMelee
	procMask ^= rogue.getImbueProcMask(proto.WeaponImbue_InstantPoison)
	procMask ^= rogue.getImbueProcMask(proto.WeaponImbue_DeadlyPoison)
	procMask ^= rogue.getImbueProcMask(proto.WeaponImbue_WoundPoison)

	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Deadly Brew (Instant)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Instant Poison") < rogue.GetInstantPoisonProcChance() {
				rogue.InstantPoison[DeadlyBrewProc].Cast(sim, result.Target)
			}
		},
	})
}

// Apply Deadly Brew Deadly Poison procs
func (rogue *Rogue) applyDeadlyBrewDeadly() {
	rogue.RegisterAura(core.Aura{
		Label:    "Deadly Brew (Deadly)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.Flags.Matches(SpellFlagDeadlyBrewed) {
				return
			}
			rogue.DeadlyPoison[DeadlyBrewProc].Cast(sim, result.Target)
		},
	})
}

// Apply Instant Poison to weapon and enable procs
func (rogue *Rogue) applyInstantPoison() {
	procMask := rogue.getImbueProcMask(proto.WeaponImbue_InstantPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Instant Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			if sim.RandomFloat("Instant Poison") < rogue.GetInstantPoisonProcChance() {
				rogue.InstantPoison[NormalProc].Cast(sim, result.Target)
			}
		},
	})
}

// Apply Deadly Poison to weapon and enable procs
func (rogue *Rogue) applyDeadlyPoison() {
	procMask := rogue.getImbueProcMask(proto.WeaponImbue_DeadlyPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Deadly Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Deadly Poison") < rogue.GetDeadlyPoisonProcChance() {
				rogue.DeadlyPoison[NormalProc].Cast(sim, result.Target)
			}
		},
	})
}

// Apply Wound Poison to weapon and enable procs
func (rogue *Rogue) applyWoundPoison() {
	procMask := rogue.getImbueProcMask(proto.WeaponImbue_WoundPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Wound Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			if sim.RandomFloat("Wound Poison") < rogue.GetWoundPoisonProcChance() {
				rogue.WoundPoison[NormalProc].Cast(sim, result.Target)
			}
		},
	})
}

///////////////////////////////////////////////////////////////////////////
//                              Register Poisons
///////////////////////////////////////////////////////////////////////////

func (rogue *Rogue) registerInstantPoisonSpell() {
	rogue.InstantPoison = [3]*core.Spell{
		rogue.makeInstantPoison(NormalProc),
		rogue.makeInstantPoison(ShivProc),
		rogue.makeInstantPoison(DeadlyBrewProc),
	}
}

func (rogue *Rogue) registerDeadlyPoisonSpell() {
	baseDamageTick := map[int32]float64{
		25: 9,
		40: 13,
		50: 20,
		60: 27,
	}[rogue.Level]
	spellID := map[int32]int32{
		25: 2823,
		40: 2824,
		50: 11355,
		60: 11356,
	}[rogue.Level]

	hasDeadlyBrew := rogue.HasRune(proto.RogueRune_RuneDeadlyBrew)

	rogue.deadlyPoisonTick = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID, Tag: 100},
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskWeaponProc,
		Flags:       SpellFlagCarnage | core.SpellFlagPoison | SpellFlagRoguePoison,

		DamageMultiplier: rogue.getPoisonDamageMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "DeadlyPoison",
				MaxStacks: 5,
				Duration:  time.Second * 12,
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, applyStack bool) {
				if !applyStack {
					return
				}

				// only the first stack snapshots the multiplier
				if dot.GetStacks() == 1 {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType]
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
					dot.SnapshotBaseDamage = 0
				}

				// each stack snapshots the AP it was applied with
				// 3.6% per stack for all ticks, or 0.9% per stack and tick
				dot.SnapshotBaseDamage += baseDamageTick + core.TernaryFloat64(hasDeadlyBrew, 0.009*dot.Spell.MeleeAttackPower(), 0)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
			},
		},
	})

	rogue.DeadlyPoison = [3]*core.Spell{
		rogue.makeDeadlyPoison(NormalProc),
		rogue.makeDeadlyPoison(ShivProc),
		rogue.makeDeadlyPoison(DeadlyBrewProc),
	}
}

func (rogue *Rogue) registerWoundPoisonSpell() {
	woundPoisonDebuffAura := core.Aura{
		Label:     "WoundPoison-" + strconv.Itoa(int(rogue.Index)),
		ActionID:  core.ActionID{SpellID: 13219},
		MaxStacks: 5,
		Duration:  time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// all healing effects used on target reduced by x, stacks 5 times
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			// undo reduced healing effects used on targets
		},
	}

	rogue.woundPoisonDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return target.RegisterAura(woundPoisonDebuffAura)
	})
	rogue.WoundPoison = [2]*core.Spell{
		rogue.makeWoundPoison(NormalProc),
		rogue.makeWoundPoison(ShivProc),
	}
}

///////////////////////////////////////////////////////////////////////////
//                              Make Poisons
///////////////////////////////////////////////////////////////////////////

// Make a source based variant of Instant Poison
func (rogue *Rogue) makeInstantPoison(procSource PoisonProcSource) *core.Spell {
	baseDamageByLevel := map[int32]float64{
		25: 19,
		40: 44,
		50: 67,
		60: 112,
	}[rogue.Level]

	damageVariance := map[int32]float64{
		25: 6,
		40: 12,
		50: 18,
		60: 36,
	}[rogue.Level]

	spellID := map[int32]int32{
		25: 8679,
		40: 8688,
		50: 11338,
		60: 11340,
	}[rogue.Level]

	hasDeadlyBrew := rogue.HasRune(proto.RogueRune_RuneDeadlyBrew)

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskWeaponProc,
		Flags:       SpellFlagDeadlyBrewed | SpellFlagCarnage | core.SpellFlagPoison | SpellFlagRoguePoison,

		DamageMultiplier: rogue.getPoisonDamageMultiplier(),
		ThreatMultiplier: 1,

		BonusHitRating: core.TernaryFloat64(procSource == ShivProc, 100*core.SpellHitRatingPerHitChance, 0),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageByLevel, baseDamageByLevel+damageVariance) + core.TernaryFloat64(hasDeadlyBrew, 0.03*spell.MeleeAttackPower(), 0)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (rogue *Rogue) makeDeadlyPoison(procSource PoisonProcSource) *core.Spell {
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: rogue.deadlyPoisonTick.SpellID, Tag: int32(procSource)},
		Flags:    core.Ternary(procSource == DeadlyBrewProc, core.SpellFlagNone, SpellFlagDeadlyBrewed),

		BonusHitRating: core.TernaryFloat64(procSource == ShivProc, 100*core.SpellHitRatingPerHitChance, 0),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			if !result.Landed() {
				return
			}

			dot := rogue.deadlyPoisonTick.Dot(target)

			dot.ApplyOrRefresh(sim)
			if dot.GetStacks() < dot.MaxStacks {
				dot.AddStack(sim)
				// snapshotting only takes place when adding a stack
				dot.TakeSnapshot(sim, true)
			}
		},
	})
}

// Make a source based variant of Wound Poison
func (rogue *Rogue) makeWoundPoison(procSource PoisonProcSource) *core.Spell {
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 13219, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskWeaponProc,
		Flags:       SpellFlagDeadlyBrewed | core.SpellFlagPoison | SpellFlagRoguePoison,

		DamageMultiplier: rogue.getPoisonDamageMultiplier(),
		ThreatMultiplier: 1,

		BonusHitRating: core.TernaryFloat64(procSource == ShivProc, 100*core.SpellHitRatingPerHitChance, 0),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			if !result.Landed() {
				return
			}

			aura := rogue.woundPoisonDebuffAuras.Get(target)
			if !aura.IsActive() {
				aura.Activate(sim)
				aura.SetStacks(sim, 1)
				return
			}

			if aura.GetStacks() < 5 {
				aura.Refresh(sim)
				aura.AddStack(sim)
				return
			}
			aura.Refresh(sim)
		},
	})
}
