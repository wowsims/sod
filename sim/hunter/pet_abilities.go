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
	DemoralizingScreech
	FuriousHowl
	LightningBreath
	ScorpidPoison
	Swipe
)

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) *core.Spell {
	switch abilityType {
	case Bite:
		return hp.newBite()
	case Claw:
		return hp.newClaw()
	// case DemoralizingScreech:
	// 	return hp.newDemoralizingScreech()
	// case FuriousHowl:
	// 	return hp.newFuriousHowl()
	case LightningBreath:
		return hp.newLightningBreath()
	case ScorpidPoison:
		return hp.newScorpidPoison()
	// case Swipe:
	// 	return hp.newSwipe()
	case Unknown:
		return nil
	default:
		panic("Invalid pet ability type")
	}
}

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

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

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
		CritMultiplier:   hp.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax) + (spell.MeleeAttackPower() * 1.5 / 14) + spell.BonusWeaponDamage()
			baseDamage *= hp.killCommandMult()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

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

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

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
		CritMultiplier:   hp.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax) + (spell.MeleeAttackPower() * 2.15 / 14) + spell.BonusWeaponDamage()
			baseDamage *= hp.killCommandMult()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

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

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,

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
		CritMultiplier:   hp.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: Figure out how LB scales as it doesnt seem to be from SP even tho the spell is listed
			// with a sp mod
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax) + (spell.MeleeAttackPower() * 2.15 / 14)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

// func (hp *HunterPet) newDemoralizingScreech() *core.Spell {
// 	//debuffs := hp.NewEnemyAuraArray(core.DemoralizingScreechAura)

// 	return hp.newSpecialAbility(PetSpecialAbilityConfig{
// 		Type:    DemoralizingScreech,
// 		Cost:    20,
// 		GCD:     PetGCD,
// 		CD:      time.Second * 10,
// 		SpellID: 55487,
// 		School:  core.SpellSchoolPhysical,
// 		MinDmg:  85,
// 		MaxDmg:  129,
// 		APRatio: 0.07,
// 		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if result.Landed() {
// 				//for _, aoeTarget := range sim.Encounter.TargetUnits {
// 				//debuffs.Get(aoeTarget).Activate(sim)
// 				//}
// 			}
// 		},
// 	})
// }

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

func (hp *HunterPet) newScorpidPoison() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55728},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,

		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ScorpidPoison",
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(100/5, 130/5) + (0.07/5)*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newSporeCloud() *core.Spell {
	//debuffs := hp.NewEnemyAuraArray(core.SporeCloudAura)
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53598},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,

		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "SporeCloud",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(22, 28) + (0.049/3)*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
			//for _, target := range spell.Unit.Env.Encounter.TargetUnits {
			//debuffs.Get(target).Activate(sim)
			//}
		},
	})
}

func (hp *HunterPet) newVenomWebSpray() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55509},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: time.Second * 40,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "VenomWebSpray",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 46 + 0.07*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}
