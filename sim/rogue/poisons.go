package rogue

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	// "github.com/wowsims/sod/sim/core/stats"
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
60 R4: 108 damage, 11356 ID, 105 charges
60 R5: 136 damage, 25351 ID, 120 charges

Wound Poison: 30% proc chance, 5 stacks
25: x damage, x ID (none, first rank is level 32)
40: -75 healing, 11325 ID, 75 charges (Rank 2)
50: -105 healing, 13226 ID, 90 charges (Rank 3)
60: -135 healing, 13227 ID, 105 charges (Rank 4)

Occult Poison: 30% proc chance, 5 stacks
Benefits from all Deadly Poison effects
56: 108 damage, 458821 ID, 30 minute duration (Rank 1)
60: 136 damage, 1214168 ID, 30 minute duration (Rank 2)

Sebacious Poison: 30% proc chance
60: 1700 armor for 15 sec


*/

// TODO: Add charges to poisons (not deadly brew)

type PoisonProcSource int

const (
	NormalProc PoisonProcSource = iota
	ShivProc
	DeadlyBrewProc
)

func (rogue *Rogue) GetInstantPoisonProcChance() float64 {
	return (0.2+rogue.improvedPoisonsBonusProcChance())*(1+rogue.instantPoisonProcChanceBonus) + rogue.additivePoisonBonusChance
}

// Used for all 30% proc poisons (Sebacious and others)
func (rogue *Rogue) GetDeadlyPoisonProcChance() float64 {
	return 0.3 + rogue.improvedPoisonsBonusProcChance() + rogue.additivePoisonBonusChance
}

func (rogue *Rogue) GetWoundPoisonProcChance() float64 {
	return 0.3 + rogue.improvedPoisonsBonusProcChance() + rogue.additivePoisonBonusChance
}

func (rogue *Rogue) improvedPoisonsBonusProcChance() float64 {
	return 0.02 * float64(rogue.Talents.ImprovedPoisons)
}

///////////////////////////////////////////////////////////////////////////
//                               Apply Poisons
///////////////////////////////////////////////////////////////////////////

func (rogue *Rogue) applyPoisons() {
	rogue.applyDeadlyPoison()
	rogue.applyInstantPoison()
	rogue.applyWoundPoison()
	rogue.applyOccultPoison()
	rogue.applySebaciousPoison()
	rogue.applyAtrophicPoison()
	rogue.applyNumbingPoison()

	if rogue.Options.PkSwap && rogue.HasRune(proto.RogueRune_RunePoisonedKnife) {
		rogue.RegisterAura(core.Aura{
			Label:    "Apply Sebacious on pull (PK Swap)",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				rogue.SebaciousPoison[1].Cast(sim, sim.GetTargetUnit(0))
			},
		})
	}
}

// Apply Deadly Brew Instant Poison procs
func (rogue *Rogue) applyDeadlyBrewInstant() {
	// apply IP from all weapons w/o Poisons applied
	procMask := core.ProcMaskMelee
	procMask ^= rogue.getImbueProcMask(proto.WeaponImbue_InstantPoison)
	procMask ^= rogue.getImbueProcMask(proto.WeaponImbue_DeadlyPoison)
	procMask ^= rogue.getImbueProcMask(proto.WeaponImbue_OccultPoison)
	procMask ^= rogue.getImbueProcMask(proto.WeaponImbue_WoundPoison)
	procMask ^= rogue.getImbueProcMask(proto.WeaponImbue_SebaciousPoison)
	procMask ^= rogue.getImbueProcMask(proto.WeaponImbue_AtrophicPoison)
	procMask ^= rogue.getImbueProcMask(proto.WeaponImbue_NumbingPoison)

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
	if rogue.Level == 60 {
		rogue.usingOccult = true
	} else {
		rogue.usingDeadly = true
	}

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
			if rogue.usingOccult {
				rogue.OccultPoison[DeadlyBrewProc].Cast(sim, result.Target)
			} else {
				rogue.DeadlyPoison[DeadlyBrewProc].Cast(sim, result.Target)
			}
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

	rogue.usingDeadly = true

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

// Apply Occult Poison to weapon and enable procs
func (rogue *Rogue) applyOccultPoison() {
	procMask := rogue.getImbueProcMask(proto.WeaponImbue_OccultPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.usingOccult = true

	rogue.RegisterAura(core.Aura{
		Label:    "Occult Poison Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Occult Poison") < rogue.GetDeadlyPoisonProcChance() {
				rogue.OccultPoison[NormalProc].Cast(sim, result.Target)
			}
		},
	})
}

// Apply Sebacious Poison to weapon and enable procs
func (rogue *Rogue) applySebaciousPoison() {
	procMask := rogue.getImbueProcMask(proto.WeaponImbue_SebaciousPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Sebacious Poison Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Sebacious Poison") < rogue.GetDeadlyPoisonProcChance() {
				rogue.SebaciousPoison[NormalProc].Cast(sim, result.Target)
			}
		},
	})
}

// Apply Atrophic Poison to weapon and enable procs
func (rogue *Rogue) applyAtrophicPoison() {
	procMask := rogue.getImbueProcMask(proto.WeaponImbue_AtrophicPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Atrophic Poison Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Atrophic Poison") < rogue.GetDeadlyPoisonProcChance() {
				rogue.AtrophicPoison[NormalProc].Cast(sim, result.Target)
			}
		},
	})
}

// Apply Numbing Poison to weapon and enable procs
func (rogue *Rogue) applyNumbingPoison() {
	procMask := rogue.getImbueProcMask(proto.WeaponImbue_NumbingPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Numbing Poison Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Numbing Poison") < rogue.GetDeadlyPoisonProcChance() {
				rogue.NumbingPoison[NormalProc].Cast(sim, result.Target)
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

// TODO: Check this before PR
func (rogue *Rogue) registerDeadlyPoisonSpell() {
	baseDamageTick := map[int32]float64{
		25: 9,
		40: 13,
		50: 20,
		60: 34, //updated to Rank 5
	}[rogue.Level]
	spellID := map[int32]int32{
		25: 2823,
		40: 2824,
		50: 11355,
		60: 25351,
	}[rogue.Level]

	hasDeadlyBrew := rogue.HasRune(proto.RogueRune_RuneDeadlyBrew)

	rogue.deadlyPoisonTick = rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueDeadlyPoisonTick,
		ActionID:       core.ActionID{SpellID: spellID, Tag: 100},
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamageProc,
		Flags:          core.SpellFlagPoison | core.SpellFlagPassiveSpell | SpellFlagRoguePoison | SpellFlagCarnage,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "DeadlyPoison",
				MaxStacks: 5,
				Duration:  time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					// p8 DPS tier bonus tracking
					fmt.Println("Deadly Poison activated")
					rogue.PoisonsActive[aura.Unit.UnitIndex]++
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					// p8 DPS tier bonus tracking
					fmt.Println("Deadly Poison deactivated")
					rogue.PoisonsActive[aura.Unit.UnitIndex]--
				},
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
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
					dot.SnapshotBaseDamage = 0
				}

				// each stack snapshots the AP it was applied with
				// 3.6% per stack for all ticks, or 0.9% per stack and tick
				dot.SnapshotBaseDamage += baseDamageTick + core.TernaryFloat64(hasDeadlyBrew, 0.009*dot.Spell.MeleeAttackPower(), 0)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
	})

	rogue.DeadlyPoison = [3]*core.Spell{
		rogue.makeDeadlyPoison(NormalProc),
		rogue.makeDeadlyPoison(ShivProc),
		rogue.makeDeadlyPoison(DeadlyBrewProc),
	}
}

// TODO: check this before PR
func (rogue *Rogue) registerWoundPoisonSpell() {
	woundPoisonDebuffAura := core.Aura{
		Label:     "WoundPoison-" + strconv.Itoa(int(rogue.Index)),
		ActionID:  core.ActionID{SpellID: 13219},
		MaxStacks: 5,
		Duration:  time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// all healing effects used on target reduced by x, stacks 5 times

			// p8 DPS tier bonus tracking
			fmt.Println("Wound Poison activated")
			rogue.PoisonsActive[aura.Unit.UnitIndex]++
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			// undo reduced healing effects used on targets

			// p8 DPS tier bonus tracking
			fmt.Println("Wound Poison deactivated")
			rogue.PoisonsActive[aura.Unit.UnitIndex]--
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

// TODO: check this before PR
func (rogue *Rogue) registerOccultPoisonSpell() {
	if rogue.Level < 60 {
		return
	}

	baseDamageTick := float64(34) //Updated to Rank 2
	spellID := int32(1214170)

	hasDeadlyBrew := rogue.HasRune(proto.RogueRune_RuneDeadlyBrew)

	rogue.occultPoisonTick = rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueOccultPoisonTick,
		ActionID:       core.ActionID{SpellID: spellID, Tag: 100},
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamageProc,
		Flags:          SpellFlagCarnage | core.SpellFlagPoison | SpellFlagRoguePoison,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "OccultPoison",
				MaxStacks: 5,
				Duration:  time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					// p8 DPS tier bonus tracking
					fmt.Println("Occult Poison activated")
					rogue.PoisonsActive[aura.Unit.UnitIndex]++
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					// p8 DPS tier bonus tracking
					fmt.Println("Occult Poison deactivated")
					rogue.PoisonsActive[aura.Unit.UnitIndex]--
				},
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
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
					dot.SnapshotBaseDamage = 0
				}

				// each stack snapshots the AP it was applied with
				// 3.6% per stack for all ticks, or 0.9% per stack and tick
				dot.SnapshotBaseDamage += baseDamageTick + core.TernaryFloat64(hasDeadlyBrew, 0.009*dot.Spell.MeleeAttackPower(), 0)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
	})

	rogue.OccultPoison = [3]*core.Spell{
		rogue.makeOccultPoison(NormalProc),
		rogue.makeOccultPoison(ShivProc),
		rogue.makeOccultPoison(DeadlyBrewProc),
	}
}

// TODO: Figure out how to either use the Rogue struct in debuffs since the aura is contructed in debuffs.go
// or move the aura constructor here
func (rogue *Rogue) registerSebaciousPoisonSpell() {
	if rogue.Level < 60 {
		return
	}

	rogue.sebaciousPoisonDebuffAura = rogue.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		sebaciousPoisonAura := core.SebaciousPoisonAura(unit, rogue.Talents.ImprovedExposeArmor, rogue.Level)

		sebaciousPoisonAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			// p8 DPS tier bonus tracking
			fmt.Println("Sebacious Poison activated")
			rogue.PoisonsActive[aura.Unit.UnitIndex]++
		})

		sebaciousPoisonAura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
			// p8 DPS tier bonus tracking
			fmt.Println("Sebacious Poison deactivated")
			rogue.PoisonsActive[aura.Unit.UnitIndex]--
		})

		return sebaciousPoisonAura
	})

	rogue.SebaciousPoison = [2]*core.Spell{
		rogue.makeSebaciousPoison(NormalProc),
		rogue.makeSebaciousPoison(ShivProc),
	}

}

// TODO: Figure out how to either use the Rogue struct in debuffs since the aura is contructed in debuffs.go
// or move the aura constructor here
func (rogue *Rogue) registerAtrophicPoisonSpell() {
	if rogue.Level < 60 {
		return
	}

	rogue.atrophicPoisonDebuffAura = rogue.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		atrophicPoisonAura := core.AtrophicPoisonAura(unit)

		atrophicPoisonAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			// p8 DPS tier bonus tracking
			fmt.Println("Atrophic Poison activated")
			rogue.PoisonsActive[aura.Unit.UnitIndex]++
		})

		atrophicPoisonAura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
			// p8 DPS tier bonus tracking
			fmt.Println("Atrophic Poison deactivated")
			rogue.PoisonsActive[aura.Unit.UnitIndex]--
		})

		return atrophicPoisonAura
	})

	rogue.AtrophicPoison = [2]*core.Spell{
		rogue.makeAtrophicPoison(NormalProc),
		rogue.makeAtrophicPoison(ShivProc),
	}

}

// TODO: Figure out how to either use the Rogue struct in debuffs since the aura is contructed in debuffs.go
// or move the aura constructor here
func (rogue *Rogue) registerNumbingPoisonSpell() {
	if rogue.Level < 60 {
		return
	}

	rogue.numbingPoisonDebuffAura = rogue.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		numbingPoisonAura := core.NumbingPoisonAura(unit)

		numbingPoisonAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			// p8 DPS tier bonus tracking
			fmt.Println("Numbing Poison activated")
			rogue.PoisonsActive[aura.Unit.UnitIndex]++
		})

		numbingPoisonAura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
			// p8 DPS tier bonus tracking
			fmt.Println("Numbing Poison deactivated")
			rogue.PoisonsActive[aura.Unit.UnitIndex]--
		})

		return numbingPoisonAura
	})

	rogue.NumbingPoison = [2]*core.Spell{
		rogue.makeNumbingPoison(NormalProc),
		rogue.makeNumbingPoison(ShivProc),
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
		ClassSpellMask: ClassSpellMask_RogueInstantPoison,
		ActionID:       core.ActionID{SpellID: spellID, Tag: int32(procSource)},
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamageProc,
		Flags:          core.SpellFlagPoison | core.SpellFlagPassiveSpell | SpellFlagDeadlyBrewed | SpellFlagCarnage | SpellFlagRoguePoison,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BonusHitRating: core.TernaryFloat64(procSource == ShivProc, 100*core.SpellHitRatingPerHitChance, 0),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageByLevel, baseDamageByLevel+damageVariance) + core.TernaryFloat64(hasDeadlyBrew, 0.05*spell.MeleeAttackPower(), 0)
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

			if rogue.Level == 60 && rogue.occultPoisonTick.Dot(target).IsActive() {
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

func (rogue *Rogue) makeOccultPoison(procSource PoisonProcSource) *core.Spell {

	rogue.occultPoisonDebuffAuras = rogue.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		return core.OccultPoisonDebuffAura(unit, rogue.Level)
	})

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: rogue.occultPoisonTick.SpellID, Tag: int32(procSource)},
		Flags:    core.Ternary(procSource == DeadlyBrewProc, core.SpellFlagNone, SpellFlagDeadlyBrewed),

		BonusHitRating: core.TernaryFloat64(procSource == ShivProc, 100*core.SpellHitRatingPerHitChance, 0),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			if !result.Landed() {
				return
			}

			rogue.deadlyPoisonTick.Dot(target).Deactivate(sim)
			rogue.occultPoisonDebuffAuras.Get(target).Activate(sim)
			rogue.occultPoisonDebuffAuras.Get(target).AddStack(sim)

			dot := rogue.occultPoisonTick.Dot(target)

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
		ProcMask:    core.ProcMaskSpellDamageProc,
		Flags:       core.SpellFlagPoison | core.SpellFlagPassiveSpell | SpellFlagDeadlyBrewed | SpellFlagRoguePoison,

		DamageMultiplier: 1,
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

func (rogue *Rogue) makeSebaciousPoison(procSource PoisonProcSource) *core.Spell {
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 439500, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamageProc,
		Flags:       core.SpellFlagPoison | core.SpellFlagPassiveSpell | SpellFlagDeadlyBrewed | SpellFlagRoguePoison,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BonusHitRating: core.TernaryFloat64(procSource == ShivProc, 100*core.SpellHitRatingPerHitChance, 0),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			if !result.Landed() {
				return
			}

			rogue.sebaciousPoisonDebuffAura.Get(target).Activate(sim)
		},
	})
}

func (rogue *Rogue) makeAtrophicPoison(procSource PoisonProcSource) *core.Spell {
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 439473, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamageProc,
		Flags:       core.SpellFlagPoison | core.SpellFlagPassiveSpell | SpellFlagDeadlyBrewed | SpellFlagRoguePoison,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BonusHitRating: core.TernaryFloat64(procSource == ShivProc, 100*core.SpellHitRatingPerHitChance, 0),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			if !result.Landed() {
				return
			}

			rogue.atrophicPoisonDebuffAura.Get(target).Activate(sim)
		},
	})
}

func (rogue *Rogue) makeNumbingPoison(procSource PoisonProcSource) *core.Spell {
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 439472, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamageProc,
		Flags:       core.SpellFlagPoison | core.SpellFlagPassiveSpell | SpellFlagDeadlyBrewed | SpellFlagRoguePoison,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BonusHitRating: core.TernaryFloat64(procSource == ShivProc, 100*core.SpellHitRatingPerHitChance, 0),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			if !result.Landed() {
				return
			}

			rogue.numbingPoisonDebuffAura.Get(target).Activate(sim)
		},
	})
}
