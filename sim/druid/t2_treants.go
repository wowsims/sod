package druid

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// We have limited data on these treants, so numbers are mostly made up.
// A lot of the stats and scaling were copied from Felhunter and Warlock pets with tweaks made to
// adjust to be closer to Testwerk logs.

type T2Treants struct {
	core.Pet

	Druid           *Druid
	Shadowcrawl     *core.Spell
	ShadowcrawlAura *core.Aura
}

func (druid *Druid) NewT2Treants() *T2Treants {
	// TODO: Figure out stats
	baseDamageMin := 70.0
	baseDamageMax := 97.0
	baseStats := stats.Stats{
		stats.Strength:  129,
		stats.Agility:   85,
		stats.Stamina:   234,
		stats.Intellect: 70,
		stats.Spirit:    150,
		stats.Mana:      1066,
		stats.MP5:       0,
		stats.MeleeHit:  4 * core.MeleeHitRatingPerHitChance,
		stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
		stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
	}

	treants := &T2Treants{
		Pet:   core.NewPet("Treants", &druid.Character, baseStats, druid.t2TreantsStatInheritance(), false, true),
		Druid: druid,
	}

	treants.EnableManaBarWithModifier(.77)

	treants.PseudoStats.DamageTakenMultiplier *= 0.1

	treants.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	treants.AddStat(stats.AttackPower, -20)

	// Warrior crit scaling
	treants.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[proto.Class_ClassWarrior][int(treants.Level)]*core.CritRatingPerCritChance)
	treants.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[proto.Class_ClassWarrior][int(treants.Level)]*core.SpellCritRatingPerCritChance)

	treants.EnableAutoAttacks(treants, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:     baseDamageMin,
			BaseDamageMax:     baseDamageMax,
			SwingSpeed:        0.6667, // 3 Treants attack at 2 second intervals. To avoid creating 3 pets, have the AI swing at 2/3 second intervals
			AttackPowerPerDPS: 14.0 / 6.0,
			SpellSchool:       core.SpellSchoolPhysical,
		},
		AutoSwingMelee: true,
	})

	druid.AddPet(treants)

	return treants
}

func (druid *Druid) t2TreantsStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:   ownerStats[stats.Intellect] * 0.3,
			stats.Armor:       ownerStats[stats.Armor] * 0.35,
			stats.AttackPower: ownerStats[stats.AttackPower] * 0.565,
			stats.MP5:         ownerStats[stats.Intellect] * 0.315,
			stats.SpellPower:  ownerStats[stats.SpellPower] * 0.15,
			stats.SpellDamage: ownerStats[stats.SpellDamage] * 0.15,
		}
	}
}

func (treants *T2Treants) Initialize() {
}

func (treants *T2Treants) ExecuteCustomRotation(sim *core.Simulation) {
}

func (treants *T2Treants) Reset(sim *core.Simulation) {
	treants.Disable(sim)
}

func (treants *T2Treants) OnPetDisable(sim *core.Simulation) {
}

func (treants *T2Treants) GetPet() *core.Pet {
	return &treants.Pet
}
