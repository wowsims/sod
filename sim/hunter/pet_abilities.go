package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

type PetAbilityType int

// Pet AI doesn't use abilities immediately, so model this with a 1.6s GCD.
const PetGCD = time.Millisecond * 1600

const (
	Unknown PetAbilityType = iota
	Bite
	Claw
	Screech
	FuriousHowl
	LightningBreath
	ScorpidPoison
	LavaBreath
)

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) *core.Spell {
	switch abilityType {
	case Bite:
		return hp.newBite()
	case Claw:
		return hp.newClaw()
	case Screech:
		return hp.newScreech()
	// case FuriousHowl:
	// 	return hp.newFuriousHowl()
	case LightningBreath:
		return hp.newLightningBreath()
	case ScorpidPoison:
		return hp.newScorpidPoison()
	// case Swipe:
	// 	return hp.newSwipe()
	case LavaBreath:
		return hp.newLavaBreath()
	case Unknown:
		return nil
	default:
		panic("Invalid pet ability type")
	}
}

// https://www.wowhead.com/classic/spell=3009/claw
func (hp *HunterPet) newClaw() *core.Spell {
	baseDamageMin := map[int32]float64{
		25: 16,
		40: 26,
		50: 35,
		60: 43,
	}[hp.Owner.Level]

	baseDamageMax := map[int32]float64{
		25: 22,
		40: 36,
		50: 49,
		60: 59,
	}[hp.Owner.Level]

	spellID := map[int32]int32{
		25: 16830,
		40: 16832,
		50: 3010,
		60: 3009,
	}[hp.Owner.Level]

	ApCoeff := 1.5 / 14

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		ClassSpellMask: ClassSpellMask_HunterPetClaw,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,

		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		MaxRange: 5,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax) + (spell.MeleeAttackPower() * ApCoeff)
			baseDamage *= hp.killCommandMult()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// https://www.wowhead.com/classic/spell=17261/bite
func (hp *HunterPet) newBite() *core.Spell {
	baseDamageMin := map[int32]float64{
		25: 31,
		40: 49,
		50: 66,
		60: 81,
	}[hp.Owner.Level]

	baseDamageMax := map[int32]float64{
		25: 37,
		40: 59,
		50: 80,
		60: 91,
	}[hp.Owner.Level]

	spellID := map[int32]int32{
		25: 17257,
		40: 17259,
		50: 17260,
		60: 17261,
	}[hp.Owner.Level]

	ApCoeff := 2.15 / 14

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		ClassSpellMask: ClassSpellMask_HunterPetBite,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,

		FocusCost: core.FocusCostOptions{
			Cost: 35,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: 10 * time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		MaxRange: 5,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax) + (spell.MeleeAttackPower() * ApCoeff)
			baseDamage *= hp.killCommandMult()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// https://www.wowhead.com/classic/spell=25012/lightning-breath
func (hp *HunterPet) newLightningBreath() *core.Spell {
	baseDamageMin := map[int32]float64{
		25: 36,
		40: 36,
		50: 78,
		60: 99,
	}[hp.Owner.Level]

	baseDamageMax := map[int32]float64{
		25: 41,
		40: 41,
		50: 91,
		60: 113,
	}[hp.Owner.Level]

	spellID := map[int32]int32{
		25: 25009,
		40: 25009, // rank 4 not available in SoD Phase 2
		50: 25011,
		60: 25012,
	}[hp.Owner.Level]

	ApCoeff := 3.65 / 14
	SpCoeff := 0.429

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		ClassSpellMask: ClassSpellMask_HunterPetLightningBreath,
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,

		FocusCost: core.FocusCostOptions{
			Cost: 50,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: SpCoeff,

		MaxRange: 20,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax) + (spell.MeleeAttackPower() * ApCoeff)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

// https://www.wowhead.com/classic/spell=24579/screech
func (hp *HunterPet) newScreech() *core.Spell {
	baseDamageMin := map[int32]float64{
		25: 12,
		40: 12,
		50: 19,
		60: 26,
	}[hp.Owner.Level]

	baseDamageMax := map[int32]float64{
		25: 16,
		40: 16,
		50: 25,
		60: 46,
	}[hp.Owner.Level]

	spellID := map[int32]int32{
		15: 24580,
		40: 24580,
		50: 24581,
		60: 24582,
	}[hp.Owner.Level]

	ApCoeff := 1.15 / 14

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		ClassSpellMask: ClassSpellMask_HunterPetScreech,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics,

		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		MaxRange: 5,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax) + (spell.MeleeAttackPower() * ApCoeff)
			// This ability also applies a melee attack power reduction similar to demoralizing shout - left it out for now
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// func (hp *HunterPet) newFuriousHowl() *core.Spell {
// 	actionID := core.ActionID{SpellID: 64495}

// 	petAura := hp.NewTemporaryStatsAura("FuriousHowl", actionID, stats.Stats{stats.AttackPower: 320, stats.RangedAttackPower: 320}, time.Second*20)
// 	ownerAura := hp.hunterOwner.NewTemporaryStatsAura("FuriousHowl", actionID, stats.Stats{stats.AttackPower: 320, stats.RangedAttackPower: 320}, time.Second*20)

// 	howlSpell := hp.RegisterSpell(core.SpellConfig{
// 		ActionID: actionID,

// 		FocusCost: core.FocusCostOptions{
// 			Cost: 20,
// 		},
// 		Cast: core.CastConfig{
// 			CD: core.Cooldown{
// 				Timer:    hp.NewTimer(),
// 				Duration: time.Second * 40,
// 			},
// 		},
// 		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
// 			return hp.IsEnabled()
// 		},
// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
// 			petAura.Activate(sim)
// 			ownerAura.Activate(sim)
// 		},
// 	})

// 	hp.hunterOwner.RegisterSpell(core.SpellConfig{
// 		ActionID: actionID,
// 		Flags:    core.SpellFlagAPL | core.SpellFlagMCD,
// 		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
// 			return howlSpell.CanCast(sim, target)
// 		},
// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
// 			howlSpell.Cast(sim, target)
// 		},
// 	})

// 	hp.hunterOwner.AddMajorCooldown(core.MajorCooldown{
// 		Spell: howlSpell,
// 		Type:  core.CooldownTypeDPS,
// 	})

// 	return nil
// }

// https://www.wowhead.com/classic/spell=24587/scorpid-poison
func (hp *HunterPet) newScorpidPoison() *core.Spell {
	baseDamageTick := map[int32]float64{
		25: 3,
		40: 6,
		50: 6,
		60: 8,
	}[hp.Owner.Level]
	spellID := map[int32]int32{
		25: 24583,
		40: 24586,
		50: 24586,
		60: 24587,
	}[hp.Owner.Level]

	ApCoeff := 0.07 / 5

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		ClassSpellMask: ClassSpellMask_HunterPetScorpidPoison,
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagPassiveSpell | core.SpellFlagPoison,

		FocusCost: core.FocusCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: time.Second * 4,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		MaxRange: 5,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "ScorpidPoison",
				MaxStacks: 5,
				Duration:  time.Second * 10,
			},
			NumberOfTicks: 5,
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

				dot.SnapshotBaseDamage += baseDamageTick + ApCoeff*dot.Spell.MeleeAttackPower()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if !result.Landed() {
				return
			}

			dot := spell.Dot(target)
			dot.ApplyOrRefresh(sim)
			if dot.GetStacks() < dot.MaxStacks {
				dot.AddStack(sim)
				dot.TakeSnapshot(sim, true)
			}
		},
	})
}

// https://www.wowhead.com/classic/spell=444681/lava-breath
func (hp *HunterPet) newLavaBreath() *core.Spell {
	baseDamageMin := map[int32]float64{
		50: 78,
		60: 101,
	}[hp.Owner.Level]
	baseDamageMax := map[int32]float64{
		50: 91,
		60: 116,
	}[hp.Owner.Level]
	spellID := map[int32]int32{
		50: 444678,
		60: 444681,
	}[hp.Owner.Level]

	ApCoeff := 3.65 / 14
	SpCoeff := 0.429

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		ClassSpellMask: ClassSpellMask_HunterPetLavaBreath,
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagBinary,

		FocusCost: core.FocusCostOptions{
			Cost: 50,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: SpCoeff,

		MaxRange: 5,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamageMin, baseDamageMax) + ApCoeff*spell.MeleeAttackPower()
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			if sim.Environment.GetNumTargets() > 1 {
				target = sim.Environment.NextTargetUnit(target)
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			}

		},
	})
}
