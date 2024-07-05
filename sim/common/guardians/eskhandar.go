package guardians

import (
	"slices"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// https://www.wowhead.com/classic/item-set=1781/spirit-of-eskhandar
// https://www.wowhead.com/classic/spell=461990/call-of-eskhandar
// https://www.wowhead.com/classic/npc=14306/eskhandar

type Eskhandar struct {
	core.Pet
}

func NewEskhandar(character *core.Character) *Eskhandar {
	eskhandar := &Eskhandar{
		Pet: core.NewPet("Eskhandar", character, stats.Stats{}, eskhandarStatInheritance(), false, true),
	}

	// TODO: Verify
	eskhandar.Level = 60

	eskhandar.EnableAutoAttacks(eskhandar, core.AutoAttackOptions{
		// TODO: Need Core Hound data
		MainHand: core.Weapon{
			BaseDamageMin: 1,
			BaseDamageMax: 1,
			SwingSpeed:    2.0,
			SpellSchool:   core.SpellSchoolPhysical,
		},
		AutoSwingMelee: true,
	})

	return eskhandar
}

func eskhandarStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		// TODO: Needs more verification
		return stats.Stats{}
	}
}

func (eskhandar *Eskhandar) Initialize() {
}

func (eskhandar *Eskhandar) ExecuteCustomRotation(sim *core.Simulation) {
	// Run the cast check only on swings or cast completes
	if eskhandar.AutoAttacks.NextAttackAt() != sim.CurrentTime+eskhandar.AutoAttacks.MainhandSwingSpeed() && eskhandar.AutoAttacks.NextAnyAttackAt()-1 > sim.CurrentTime {
		eskhandar.WaitUntil(sim, eskhandar.AutoAttacks.NextAttackAt()-1)
		return
	}

	eskhandar.WaitUntil(sim, eskhandar.AutoAttacks.NextAttackAt()-1)
}

func (eskhandar *Eskhandar) Reset(sim *core.Simulation) {
	eskhandar.Disable(sim)
}

func (eskhandar *Eskhandar) OnPetDisable(sim *core.Simulation) {
}

func (eskhandar *Eskhandar) GetPet() *core.Pet {
	return &eskhandar.Pet
}

func constructEskhandar(character *core.Character) {
	// Can't use the set bonus itself because of an import cycle
	hasSetBonus := slices.ContainsFunc(character.GetActiveSetBonuses(), func(bonus core.ActiveSetBonus) bool {
		return bonus.Name == "Spirit of Eskhandar" && bonus.NumPieces == 4
	})

	if hasSetBonus {
		character.AddPet(NewEskhandar(character))
	}
}
