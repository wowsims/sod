package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) makeSuccubus() *WarlockPet {
	cfg := PetConfig{
		Name:          "Succubus",
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
			stats.Mana:      119,
			stats.MP5:       0,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin: 23,
				BaseDamageMax: 38,
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
			stats.Mana:      521,
			stats.MP5:       0,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin: 41,
				BaseDamageMax: 61,
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
				// Not updated
				BaseDamageMin: 41,
				BaseDamageMax: 61,
				SwingSpeed:    2,
			},
			AutoSwingMelee: true,
		}
	case 60:
		cfg.Stats = stats.Stats{
			stats.Strength:  129,
			stats.Agility:   85,
			stats.Stamina:   234,
			stats.Intellect: 70,
			stats.Spirit:    150,
			stats.Mana:      1066,
			stats.MP5:       0,
			stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
			stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
		}
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin: 95,
				BaseDamageMax: 131,
				SwingSpeed:    2,
			},
			AutoSwingMelee: true,
		}
	}

	return warlock.makePet(cfg, warlock.Options.Summon == proto.WarlockOptions_Succubus)
}

func (wp *WarlockPet) registerSuccubusLashOfPainSpell() {
	warlockLevel := wp.owner.Level
	// assuming max rank available
	rank := map[int32]int{25: 1, 40: 3, 50: 4, 60: 6}[warlockLevel]

	if rank == 0 {
		rank = 1
	}

	spellCoeff := [7]float64{0, .429, .429, .429, .429, .429, .429}[rank]
	baseDamage := [7]float64{0, 33, 44, 60, 73, 87, 99}[rank] * (1 + .10*float64(wp.owner.Talents.ImprovedSayaad))
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
