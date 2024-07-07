package warlock

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) makeFelhunter() *WarlockPet {
	cfg := PetConfig{
		Name:          "Felhunter",
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
				// Not updated
				BaseDamageMin: 24,
				BaseDamageMax: 40,
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
				BaseDamageMin: 70,
				BaseDamageMax: 97,
				SwingSpeed:    2,
			},
			AutoSwingMelee: true,
		}
	}

	return warlock.makePet(cfg, warlock.Options.Summon == proto.WarlockOptions_Felhunter)
}
