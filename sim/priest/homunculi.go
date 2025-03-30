package priest

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (priest *Priest) registerHomunculiSpell() {
	if !priest.HasRune(proto.PriestRune_RuneLegsHomunculi) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.PriestRune_RuneLegsHomunculi)}
	duration := time.Minute * 2
	cooldown := time.Minute * 2

	// For timeline only
	priest.HomunculiAura = priest.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Homunculi",
		Duration: duration,
	})

	priest.Homunculi = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			core.Each(priest.HomunculiPets, func(homunculus *Homunculus) {
				homunculus.EnableWithTimeout(sim, homunculus, duration)
			})
			priest.HomunculiAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    priest.Homunculi,
		Priority: 1,
		Type:     core.CooldownTypeDPS,
	})
}

type Homunculus struct {
	core.Pet

	npcID        int32
	Priest       *Priest
	PrimarySpell *core.Spell
}

func (priest *Priest) NewHomunculus(idx int32, npcID int32) *Homunculus {
	// We can't use the legacy pet window for Homunculi so these values are determined from trying to bump damage up to match logs
	baseDamageMin := 0.0
	baseDamageMax := 0.0
	baseStats := stats.Stats{
		stats.Strength:    0,
		stats.Agility:     0,
		stats.Stamina:     0,
		stats.Intellect:   0,
		stats.Spirit:      0,
		stats.AttackPower: 0,
	}

	homunculus := &Homunculus{
		npcID:  npcID,
		Pet:    core.NewPet("Homunculi", &priest.Character, baseStats, priest.homunculusStatInheritance(), false, true),
		Priest: priest,
	}

	homunculus.AddStat(stats.AttackPower, -20)

	homunculus.EnableAutoAttacks(homunculus, core.AutoAttackOptions{
		MainHand: core.Weapon{
			// TODO: Check stats
			BaseDamageMin: baseDamageMin,
			BaseDamageMax: baseDamageMax,
			SwingSpeed:    2,
			MaxRange:      core.MaxMeleeAttackRange,
		},
		AutoSwingMelee: true,
	})

	homunculus.ApplyOnPetEnable(func(sim *core.Simulation) {
		// Priest pets only inherit the owner's cast speed
		homunculus.EnableDynamicCastSpeedInheritance(sim)
	})

	// Homunculus aren't very bright and often sit in front of targets
	homunculus.PseudoStats.InFrontOfTarget = true

	// core.ApplyPetConsumeEffects(&homunculus.Character, priest.Consumes)

	priest.AddPet(homunculus)

	return homunculus
}

// TODO: Verify any homunculus stat inheritance
func (priest *Priest) homunculusStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		// Shadowfiend seems to benefit from the owner's Spell Hit + Shadow Focus hit chance
		ownerHitChance := (ownerStats[stats.SpellHit] + 2*float64(priest.Talents.ShadowFocus)) / core.SpellHitRatingPerHitChance
		highestSchoolPower := ownerStats[stats.SpellPower] + ownerStats[stats.SpellDamage] + max(ownerStats[stats.HolyPower], ownerStats[stats.ShadowPower])

		// TODO: Needs more verification
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
			stats.SpellHit:         math.Floor(ownerHitChance / 12.0 * 17.0),
			// Shadowfiend seems to most likely use the priest's Spell Crit to scale its melee crit
			// In reality the melees are shadow damage and probably use the spell hit table but we can't configure that currently
			stats.MeleeCrit: ownerStats[stats.SpellCrit],
			stats.SpellCrit: ownerStats[stats.SpellCrit],
		}
	}
}

func (homunculus *Homunculus) newHomunculusCrippleSpell() core.SpellConfig {
	attackSpeedAuras := homunculus.NewEnemyAuraArray(core.HomunculiAttackSpeedAura)

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
			attackSpeedAuras.Get(target).Activate(sim)
		},
	}
}

func (homunculus *Homunculus) newHomunculusDegradeSpell() core.SpellConfig {
	armorAuras := homunculus.NewEnemyAuraArray(core.HomunculiArmorAura)

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
			armorAuras.Get(target).Activate(sim)
		},
	}
}

func (homunculus *Homunculus) newHomunculusDemoralizeSpell() core.SpellConfig {
	attackPowerAuras := homunculus.NewEnemyAuraArray(core.HomunculiAttackPowerAura)

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
			attackPowerAuras.Get(target).Activate(sim)
		},
	}
}

func (homunculus *Homunculus) Initialize() {
	switch homunculus.npcID {
	case 202390:
		homunculus.PrimarySpell = homunculus.GetOrRegisterSpell(homunculus.newHomunculusCrippleSpell())
	case 202392:
		homunculus.PrimarySpell = homunculus.GetOrRegisterSpell(homunculus.newHomunculusDegradeSpell())
	case 202391:
		homunculus.PrimarySpell = homunculus.GetOrRegisterSpell(homunculus.newHomunculusDemoralizeSpell())
	}
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
