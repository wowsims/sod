package shaman

import (
	"math"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

type SpiritWolves struct {
	core.Pet

	shamanOwner *Shaman
}

var spiritWolfBaseStats = stats.Stats{
	stats.Strength:  136,
	stats.Agility:   100,
	stats.Stamina:   265,
	stats.Intellect: 50,
	stats.Spirit:    80,
	// Base AP 265 - (136 * 2) - (100 * 0.2)
	stats.AttackPower: -27,

	// Add 1.8% because pets aren't affected by that component of crit suppression.
	stats.MeleeCrit: (1.1515 + 1.8) * core.CritRatingPerCritChance,
}

func (shaman *Shaman) NewSpiritWolves() *SpiritWolves {
	wolves := &SpiritWolves{
		Pet:         core.NewPet("Spirit Wolves", &shaman.Character, spiritWolfBaseStats, shaman.makeStatInheritance(), false, true),
		shamanOwner: shaman,
	}

	wolves.EnableAutoAttacks(wolves, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin: 55.5,
			BaseDamageMax: 71.5,
			SwingSpeed:    1.5 / 2, // Two dogs attack at 1.5s intervals, but for performance we use 1 pet so divide by 2
			MaxRange:      core.MaxMeleeAttackRange,
		},
		AutoSwingMelee: true,
	})

	// Testing found that wolves gained 2 AP per Str, and ~1 AP per 5 Agi
	// Tested using different ranks of Strength of Earth and Grace of Air Totems
	wolves.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	wolves.AddStatDependency(stats.Agility, stats.AttackPower, 0.2)

	wolves.ApplyOnPetEnable(func(sim *core.Simulation) {
		// Shaman pets inherit the highest of the owner's melee and cast speed
		wolves.EnableDynamicAttackSpeedInheritance(sim)
	})

	// Warrior crit scaling
	core.ApplyPetConsumeEffects(&wolves.Character, shaman.Consumes)

	shaman.AddPet(wolves)

	return wolves
}

func (shaman *Shaman) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerHitChance := ownerStats[stats.MeleeHit] / core.MeleeHitRatingPerHitChance
		hitRatingFromOwner := math.Floor(ownerHitChance) * core.MeleeHitRatingPerHitChance

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.3,
			stats.Armor:       ownerStats[stats.Armor] * 0.35,
			stats.AttackPower: ownerStats[stats.AttackPower] * 0.31,
			stats.MeleeCrit:   ownerStats[stats.MeleeCrit], // Logs showed high crit rates consistent with inheriting owner hit but possibly not scaling from Agi
			stats.MeleeHit:    hitRatingFromOwner,
		}
	}
}

func (wolves *SpiritWolves) Initialize() {
	// Nothing
}

func (wolves *SpiritWolves) ExecuteCustomRotation(_ *core.Simulation) {
}

func (wolves *SpiritWolves) Reset(sim *core.Simulation) {
	wolves.Disable(sim)
}

func (spiritWolf *SpiritWolves) GetPet() *core.Pet {
	return &spiritWolf.Pet
}
