package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) makeFelguard() *WarlockPet {
	cfg := PetConfig{
		Name:          "Felguard",
		PowerModifier: 0.77,
	}
	switch warlock.Level {
	case 25:
		cfg.Stats = stats.Stats{
			stats.Strength:  50,
			stats.Agility:   40,
			stats.Stamina:   87,
			stats.Intellect: 35,
			stats.Spirit:    61,
			stats.Mana:      653,
			stats.MP5:       0,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin: 24,
				BaseDamageMax: 40,
				SwingSpeed:    2,
			},
			AutoSwingMelee: true,
		}
	case 40:
		cfg.Stats = stats.Stats{
			stats.Strength:  74,
			stats.Agility:   58,
			stats.Stamina:   148,
			stats.Intellect: 49,
			stats.Spirit:    97,
			stats.Mana:      653,
			stats.MP5:       0,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin: 24,
				BaseDamageMax: 40,
				SwingSpeed:    2,
			},
			AutoSwingMelee: true,
		}
	case 50:
		cfg.Stats = stats.Stats{
			stats.Strength:  107,
			stats.Agility:   71,
			stats.Stamina:   190,
			stats.Intellect: 59,
			stats.Spirit:    123,
			stats.Mana:      912,
			stats.MP5:       0,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin: 67,
				BaseDamageMax: 101,
				SwingSpeed:    2,
			},
			AutoSwingMelee: true,
		}
	case 60:
		cfg.Stats = stats.Stats{
			stats.Strength:  129,
			stats.Agility:   85,
			stats.Stamina:   290,
			stats.Intellect: 70,
			stats.Spirit:    150,
			stats.Mana:      1066,
			stats.MP5:       0,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin: 87,
				BaseDamageMax: 128,
				SwingSpeed:    2,
			},
			AutoSwingMelee: true,
		}
	}

	pet := warlock.makePet(cfg, warlock.Options.Summon == proto.WarlockOptions_Felguard)
	// Felguard was given a ~20% damage buff on July 3rd that doesn't seem accounted for in base stats
	pet.PseudoStats.DamageDealtMultiplier *= 1.20

	return pet
}

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
				if result.Landed() {
					spell.DealDamage(sim, result)
				}
			}
		},
	})
}

func (wp *WarlockPet) registerFelguardDemonicFrenzyAura() {
	statDeps := make([]*stats.StatDependency, 11) // 10 stacks + zero condition
	for i := 1; i < 11; i++ {
		statDeps[i] = wp.NewDynamicMultiplyStat(stats.AttackPower, 1.0+.05*float64(i))
	}

	// Make a dummy copy on the Warlock for APL tracking
	ownerAura := wp.owner.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 460907},
		Label:     "Demonic Frenzy",
		Duration:  time.Second * 10,
		MaxStacks: 10,
	})

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

	core.MakePermanent(wp.RegisterAura(core.Aura{
		Label: "Demonic Frenzy Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				demonicFrenzyAura.Activate(sim)
				demonicFrenzyAura.AddStack(sim)
				ownerAura.Activate(sim)
				ownerAura.AddStack(sim)
			}
		},
	}))
}
