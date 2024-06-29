package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Imp
func (wp *WarlockPet) registerImpFireboltSpell() {
	warlockLevel := wp.owner.Level
	// assuming max rank available
	rank := map[int32]int{25: 3, 40: 5, 50: 6, 60: 7}[warlockLevel]

	if rank == 0 {
		rank = 1
	}

	if wp.owner.Options.MaxFireboltRank != proto.WarlockOptions_NoMaximum {
		rank = min(rank, int(wp.owner.Options.MaxFireboltRank))
	}

	spellCoeff := [8]float64{0, .164, .314, .529, .571, .571, .571, .571}[rank]
	baseDamage := [8][]float64{{0, 0}, {7, 10}, {14, 16}, {25, 29}, {36, 41}, {52, 59}, {72, 80}, {85, 96}}[rank]
	spellId := [8]int32{0, 3110, 7799, 7800, 7801, 7802, 11762, 11763}[rank]
	manaCost := [8]float64{0, 10, 20, 35, 50, 70, 95, 115}[rank]
	level := [8]int{0, 1, 8, 18, 28, 38, 48, 58}[rank]

	improvedImp := []float64{1, 1.1, 1.2, 1.3}[wp.owner.Talents.ImprovedImp]
	baseDamage[0] *= improvedImp
	baseDamage[1] *= improvedImp

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Millisecond * 1000,
				CastTime: time.Millisecond * (2000 - time.Duration(500*wp.owner.Talents.ImprovedFirebolt)),
			},
			// Adding an artificial CD to account for real delay in imp casts in-game
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Millisecond * 200,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1])

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.DealDamage(sim, result)
		},
	})
}

// Succubus
func (wp *WarlockPet) registerSuccubusLashOfPainSpell() {
	warlockLevel := wp.owner.Level
	// assuming max rank available
	rank := map[int32]int{25: 1, 40: 3, 50: 4, 60: 6}[warlockLevel]

	if rank == 0 {
		rank = 1
	}

	spellCoeff := [7]float64{0, .429, .429, .429, .429, .429, .429}[rank]
	baseDamage := [7]float64{0, 33, 44, 60, 73, 87, 99}[rank]
	spellId := [7]int32{0, 7814, 7815, 7816, 11778, 11779, 11780}[rank]
	manaCost := [7]float64{0, 65, 80, 105, 125, 145, 160}[rank]
	level := [7]int{0, 20, 28, 36, 44, 52, 60}[rank]

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * (12 - time.Duration(3*wp.owner.Talents.ImprovedLashOfPain)),
			},
		},

		DamageMultiplier: wp.AutoAttacks.MHConfig().DamageMultiplier,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

// Felguard
func (wp *WarlockPet) registerFelguardCleaveSpell() {
	results := make([]*core.SpellResult, min(2, wp.Env.GetNumTargets()))

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 427744},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: wp.AutoAttacks.MHConfig().DamageMultiplier,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				baseDamage := 2.0 + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}
			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}

func (wp *WarlockPet) registerFelguardDemonicFrenzyAura() {
	statDeps := make([]*stats.StatDependency, 11) // 10 stacks + zero condition
	for i := 1; i < 11; i++ {
		statDeps[i] = wp.NewDynamicMultiplyStat(stats.AttackPower, 1.0+.05*float64(i))
	}

	demonicFrenzyAura := wp.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 460907},
		Label:     "Demonic Frenzy",
		Duration:  time.Second * 10,
		MaxStacks: 10,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if oldStacks != 0 {
				aura.Unit.DisableDynamicStatDep(sim, statDeps[oldStacks])
			}
			if newStacks != 0 {
				aura.Unit.EnableDynamicStatDep(sim, statDeps[newStacks])
			}
		},
	})

	wp.RegisterAura(core.Aura{
		Label:    "Demonic Frenzy Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				demonicFrenzyAura.Activate(sim)
				demonicFrenzyAura.AddStack(sim)
			}
		},
	})
}

// func (wp *WarlockPet) registerInterceptSpell() {
// 	wp.secondaryAbility = nil // not implemented
// }

// func (wp *WarlockPet) registerShadowBiteSpell() {
// 	actionID := core.ActionID{SpellID: 54053}

// 	var petManaMetrics *core.ResourceMetrics
// 	maxManaMult := 0.04 * float64(wp.owner.Talents.ImprovedFelhunter)
// 	impFelhunter := wp.owner.Talents.ImprovedFelhunter > 0
// 	if impFelhunter {
// 		petManaMetrics = wp.NewManaMetrics(actionID)
// 	}

// 	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
// 		ActionID:    actionID,
// 		SpellSchool: core.SpellSchoolShadow,
//      DefenseType: core.DefenseTypeMagic,
// 		ProcMask:    core.ProcMaskSpellDamage,

// 		ManaCost: core.ManaCostOptions{
// 			// TODO: should be 3% of BaseMana, but it's unclear what that actually refers to with pets
// 			FlatCost: 131,
// 		},
// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: core.GCDDefault,
// 			},
// 			IgnoreHaste: true,
// 			CD: core.Cooldown{
// 				Timer:    wp.NewTimer(),
// 				Duration: time.Second * (6 - time.Duration(2*wp.owner.Talents.ImprovedFelhunter)),
// 			},
// 		},

//      BonusCritDamage: wp.owner.ruin(),

// 		DamageMultiplier: 1 + 0.03*float64(wp.owner.Talents.ShadowMastery),
// 		ThreatMultiplier: 1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := sim.Roll(97+1, 97+41) + 0.429*spell.SpellDamage()

// 			w := wp.owner
// 			spells := []*core.Spell{
// 				w.UnstableAffliction,
// 				w.Immolate,
// 				w.CurseOfAgony,
// 				w.CurseOfDoom,
// 				w.Corruption,
// 				w.Conflagrate,
// 				w.Seed,
// 				w.DrainSoul,
// 				// missing: drain life, shadowflame
// 			}
// 			counter := 0
// 			for _, spell := range spells {
// 				if spell != nil && spell.Dot(target).IsActive() {
// 					counter++
// 				}
// 			}

// 			baseDamage *= 1 + 0.15*float64(counter)

// 			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
// 			if impFelhunter && result.Landed() {
// 				wp.AddMana(sim, wp.MaxMana()*maxManaMult, petManaMetrics)
// 			}
// 			spell.DealDamage(sim, result)
// 		},
// 	})
// }
