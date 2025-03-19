package guardians

import (
	"slices"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// https://www.wowhead.com/classic/item-set=1779/core-hounds-call
// https://www.wowhead.com/classic/spell=461267/core-hounds-call
// https://www.wowhead.com/classic/npc=229001/core-hound

type CoreHound struct {
	core.Pet
}

func NewCoreHound(character *core.Character) *CoreHound {
	coreHound := &CoreHound{
		Pet: core.NewPet("Core Hound", character, stats.Stats{}, coreHoundStatInheritance(), false, true),
	}
	// TODO: Verify
	coreHound.Level = 60

	coreHound.EnableAutoAttacks(coreHound, core.AutoAttackOptions{
		// TODO: Need Core Hound data
		MainHand: core.Weapon{
			BaseDamageMin: 1,
			BaseDamageMax: 1,
			SwingSpeed:    2.0,
			SpellSchool:   core.SpellSchoolPhysical,
			MaxRange:      core.MaxMeleeAttackRange,
		},
		AutoSwingMelee: true,
	})

	return coreHound
}

func coreHoundStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		// TODO: Needs more verification
		return stats.Stats{}
	}
}

func (hound *CoreHound) Initialize() {
}

func (hound *CoreHound) ExecuteCustomRotation(sim *core.Simulation) {
	// Run the cast check only on swings or cast completes
	if hound.AutoAttacks.NextAttackAt() != sim.CurrentTime+hound.AutoAttacks.MainhandSwingSpeed() && hound.AutoAttacks.NextAnyAttackAt()-1 > sim.CurrentTime {
		hound.WaitUntil(sim, hound.AutoAttacks.NextAttackAt()-1)
		return
	}

	hound.WaitUntil(sim, hound.AutoAttacks.NextAttackAt()-1)
}

func (hound *CoreHound) Reset(sim *core.Simulation) {
	hound.Disable(sim)
}

func (hound *CoreHound) OnPetDisable(sim *core.Simulation) {
}

func (hound *CoreHound) GetPet() *core.Pet {
	return &hound.Pet
}

func constructCoreHound(character *core.Character) {
	// Can't use the set bonus itself because of an import cycle
	hasSetBonus := slices.ContainsFunc(character.GetActiveSetBonuses(), func(bonus core.ActiveSetBonus) bool {
		return bonus.Name == "Core Hound's Call" && bonus.NumPieces >= 2
	})

	if hasSetBonus {
		character.AddPet(NewCoreHound(character))
	}
}
