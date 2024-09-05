package druid

import (
	"math"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type T2Treants struct {
	core.Pet

	Druid           *Druid
	Shadowcrawl     *core.Spell
	ShadowcrawlAura *core.Aura
}

func (druid *Druid) NewT2Treants() *T2Treants {
	// TODO: Figure out stats
	baseDamageMin := 0.0
	baseDamageMax := 0.0
	baseStats := stats.Stats{
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
			SwingSpeed:        0.5,
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
		ownerHitChance := ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance
		highestSchoolPower := ownerStats[stats.SpellPower] + ownerStats[stats.SpellDamage] + max(ownerStats[stats.HolyPower], ownerStats[stats.ShadowPower])

		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * .75,
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:      highestSchoolPower * 0.57,
			stats.MP5:              ownerStats[stats.MP5] * 0.3,
			stats.SpellPower:       ownerStats[stats.SpellPower] * 0.15,
			stats.SpellDamage:      ownerStats[stats.SpellDamage] * 0.15,
			stats.ShadowPower:      ownerStats[stats.ShadowPower] * 0.15,
			stats.SpellPenetration: ownerStats[stats.SpellPenetration],
			stats.MeleeHit:         ownerHitChance * core.MeleeHitRatingPerHitChance,
			stats.SpellHit:         math.Floor(ownerStats[stats.SpellHit] / 12.0 * 17.0),
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
