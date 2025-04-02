package priest

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (priest *Priest) registerEyeOfTheVoidCD() {
	if !priest.HasRune(proto.PriestRune_RuneHelmEyeOfTheVoid) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.PriestRune_RuneHelmEyeOfTheVoid)}
	duration := time.Second * 30
	cooldown := time.Minute * 3

	// For timeline only
	priest.EyeOfTheVoidAura = priest.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Eye of the Void",
		Duration: duration,
	})

	priest.EyeOfTheVoid = priest.RegisterSpell(core.SpellConfig{
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
			priest.EyeOfTheVoidPet.EnableWithTimeout(sim, priest.EyeOfTheVoidPet, duration)
			priest.EyeOfTheVoidAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    priest.EyeOfTheVoid,
		Priority: 1,
		Type:     core.CooldownTypeDPS,
	})
}

type EyeOfTheVoid struct {
	core.Pet

	npcID      int32
	Priest     *Priest
	ShadowBolt *core.Spell
}

func (priest *Priest) NewEyeOfTheVoid() *EyeOfTheVoid {
	baseStats := stats.Stats{}
	// Little information about this pet is available so using Warlock Imp stats
	switch priest.Level {
	case 25:
		baseStats = stats.Stats{
			stats.Strength:  47,
			stats.Agility:   25,
			stats.Stamina:   49,
			stats.Intellect: 94,
			stats.Spirit:    95,
			stats.Mana:      149,
			stats.MP5:       0,
			stats.MeleeCrit: 3.454 * core.CritRatingPerCritChance,
			stats.SpellCrit: 0.9075 * core.CritRatingPerCritChance,
		}
	case 40:
		baseStats = stats.Stats{
			stats.Strength:  70,
			stats.Agility:   29,
			stats.Stamina:   67,
			stats.Intellect: 163,
			stats.Spirit:    163,
			stats.Mana:      318,
			stats.MP5:       0,
			stats.MeleeCrit: 3.454 * core.CritRatingPerCritChance,
			stats.SpellCrit: 0.9075 * core.CritRatingPerCritChance,
		}
	case 50:
		baseStats = stats.Stats{
			stats.Strength:  101,
			stats.Agility:   32,
			stats.Stamina:   71,
			stats.Intellect: 212,
			stats.Spirit:    211,
			stats.Mana:      476,
			stats.MP5:       0,
			stats.MeleeCrit: 3.454 * core.CritRatingPerCritChance,
			stats.SpellCrit: 0.9075 * core.CritRatingPerCritChance,
		}
	case 60:
		baseStats = stats.Stats{
			stats.Strength:  101,
			stats.Agility:   32,
			stats.Stamina:   71,
			stats.Intellect: 212,
			stats.Spirit:    211,
			stats.Mana:      476,
			stats.MP5:       0,
			stats.MeleeCrit: 3.454 * core.CritRatingPerCritChance,
			stats.SpellCrit: 0.9075 * core.CritRatingPerCritChance,
		}
	}

	eyePet := &EyeOfTheVoid{
		npcID:  202387,
		Pet:    core.NewPet("Eye of the Void", &priest.Character, baseStats, priest.eyeOfTheVoidStatInheritance(), false, true),
		Priest: priest,
	}

	eyePet.EnableManaBarWithModifier(0.33)

	// Imp gets 1mp/5 non casting regen per spirit
	eyePet.PseudoStats.SpiritRegenMultiplier = 1
	eyePet.PseudoStats.SpiritRegenRateCasting = 0
	eyePet.SpiritManaRegenPerSecond = func() float64 {
		// 1mp5 per spirit
		return eyePet.GetStat(stats.Spirit) / 5
	}

	// Mage spell crit scaling for imp
	eyePet.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[proto.Class_ClassMage][int(eyePet.Level)]*core.SpellCritRatingPerCritChance)

	eyePet.ApplyOnPetEnable(func(sim *core.Simulation) {
		// Priest pets only inherit the owner's cast speed
		eyePet.EnableDynamicCastSpeedInheritance(sim)
	})

	eyePet.ShadowBolt = eyePet.GetOrRegisterSpell(eyePet.newShadowBoltSpellConfig(priest))

	priest.AddPet(eyePet)

	return eyePet
}

// TODO: Verify any eye of the void stat inheritance
func (priest *Priest) eyeOfTheVoidStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerHitChance := ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance
		highestSchoolPower := ownerStats[stats.SpellPower] + ownerStats[stats.SpellDamage] + max(ownerStats[stats.FirePower], ownerStats[stats.ShadowPower])

		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:      highestSchoolPower * 0.565,
			stats.MP5:              ownerStats[stats.Intellect] * 0.315,
			stats.SpellPower:       ownerStats[stats.SpellPower] * 0.15,
			stats.SpellDamage:      ownerStats[stats.SpellDamage] * 0.15,
			stats.FirePower:        ownerStats[stats.FirePower] * 0.15,
			stats.ShadowPower:      ownerStats[stats.ShadowPower] * 0.15,
			stats.SpellPenetration: ownerStats[stats.SpellPenetration],
			stats.MeleeHit:         ownerHitChance * core.MeleeHitRatingPerHitChance,
			stats.SpellHit:         math.Floor(ownerStats[stats.SpellHit] / 12.0 * 17.0),
		}
	}
}

// TODO:
func (eyeOfTheVoid *EyeOfTheVoid) newShadowBoltSpellConfig(priest *Priest) core.SpellConfig {
	baseDamage := priest.baseRuneAbilityDamage() + 25
	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 402790},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			// BaseCost: .05,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}
}

func (eyeOfTheVoid *EyeOfTheVoid) Initialize() {
}

func (eyeOfTheVoid *EyeOfTheVoid) ExecuteCustomRotation(sim *core.Simulation) {
	eyeOfTheVoid.ShadowBolt.Cast(sim, nil)
}

func (eyeOfTheVoid *EyeOfTheVoid) Reset(sim *core.Simulation) {
	eyeOfTheVoid.Disable(sim)
}

func (eyeOfTheVoid *EyeOfTheVoid) OnPetDisable(sim *core.Simulation) {
}

func (eyeOfTheVoid *EyeOfTheVoid) GetPet() *core.Pet {
	return &eyeOfTheVoid.Pet
}
