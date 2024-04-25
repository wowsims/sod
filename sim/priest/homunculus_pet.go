package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

type Homunculus struct {
	core.Pet

	npcID        int32
	Priest       *Priest
	PrimarySpell *core.Spell
}

func (priest *Priest) NewHomunculus(idx int32, npcID int32) *Homunculus {
	baseDamageMin := 0.0
	baseDamageMax := 0.0
	baseStats := stats.Stats{
		stats.Strength:    0,
		stats.Agility:     0,
		stats.Stamina:     0,
		stats.Intellect:   0,
		stats.Spirit:      0,
		stats.AttackPower: 0,
		// Across all 3 pets there seems to be somewhere around 20% crit chance on average between all 3 pets.
		// In reality the pets have vastly different chances to crit, but they're just not worth the time to dig into right now.
		stats.MeleeCrit: 20 * core.CritRatingPerCritChance,
	}
	switch priest.Level {
	case 25:
		baseDamageMin = 10
		baseDamageMax = 20
	case 40:
		baseDamageMin = 20
		baseDamageMax = 30
	case 50:
		baseDamageMin = 30
		baseDamageMax = 40
	case 60:
		baseDamageMin = 40
		baseDamageMax = 50
	}

	homunculus := &Homunculus{
		npcID:  npcID,
		Pet:    core.NewPet("Homunculi", &priest.Character, baseStats, priest.homunculusStatInheritance(), false, true),
		Priest: priest,
	}

	switch homunculus.npcID {
	case 202390:
		homunculus.PrimarySpell = homunculus.GetOrRegisterSpell(homunculus.newHomunculusCrippleSpell())
	case 202392:
		homunculus.PrimarySpell = homunculus.GetOrRegisterSpell(homunculus.newHomunculusDegradeSpell())
	case 202391:
		homunculus.PrimarySpell = homunculus.GetOrRegisterSpell(homunculus.newHomunculusDemoralizeSpell())
	}

	homunculus.EnableAutoAttacks(homunculus, core.AutoAttackOptions{
		MainHand: core.Weapon{
			// TODO: Check stats
			BaseDamageMin: baseDamageMin,
			BaseDamageMax: baseDamageMax,
			SwingSpeed:    2,
		},
		AutoSwingMelee: true,
	})

	// core.ApplyPetConsumeEffects(&homunculus.Character, priest.Consumes)

	priest.AddPet(homunculus)

	return homunculus
}

// TODO: Verify any homunculus stat inheritance
func (priest *Priest) homunculusStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{}
	}
}

func (homunculus *Homunculus) newHomunculusCrippleSpell() core.SpellConfig {
	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 402808},
		SpellSchool: core.SpellSchoolShadow,

		ProcMask: core.ProcMaskEmpty,
		Flags:    core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			core.HomunculiAttackSpeedAura(target, homunculus.Priest.Level).Activate(sim)
		},
	}
}

func (homunculus *Homunculus) newHomunculusDegradeSpell() core.SpellConfig {
	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 402818},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 15,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			core.HomunculiArmorAura(target, homunculus.Priest.Level).Activate(sim)
		},
	}
}

func (homunculus *Homunculus) newHomunculusDemoralizeSpell() core.SpellConfig {
	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 402811},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 15,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			core.HomunculiAttackPowerAura(target, homunculus.Priest.Level).Activate(sim)
		},
	}
}

func (homunculus *Homunculus) Initialize() {
}

func (homunculus *Homunculus) ExecuteCustomRotation(sim *core.Simulation) {
	homunculus.PrimarySpell.Cast(sim, nil)
}

func (homunculus *Homunculus) Reset(sim *core.Simulation) {
	homunculus.Disable(sim)
}

func (homunculus *Homunculus) OnPetDisable(sim *core.Simulation) {
}

func (homunculus *Homunculus) GetPet() *core.Pet {
	return &homunculus.Pet
}
