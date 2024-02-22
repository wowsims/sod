package priest

import (
	"fmt"
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
	homunculusBaseStats := stats.Stats{}
	switch priest.Level {
	case 25:
	case 40:
		// 40 stats
		// TODO: All of the stats and stat inheritance needs to be verified
		baseDamageMin = 20
		baseDamageMax = 30
		homunculusBaseStats = stats.Stats{
			stats.Strength:    0,
			stats.Agility:     0,
			stats.Stamina:     0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.AttackPower: 0,
			// with 3% crit debuff, shadowfiend crits around 9-12% (TODO: verify and narrow down)
			stats.MeleeCrit: 8 * core.CritRatingPerCritChance,
		}
	case 50:
	case 60:
	}

	homunculus := &Homunculus{
		npcID:  npcID,
		Pet:    core.NewPet(fmt.Sprintf("Homunculi (%d)", idx), &priest.Character, homunculusBaseStats, priest.homunculusStatInheritance(), false, false),
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
			BaseDamageMin:  baseDamageMin,
			BaseDamageMax:  baseDamageMax,
			SwingSpeed:     2,
			CritMultiplier: priest.DefaultMeleeCritMultiplier(),
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
