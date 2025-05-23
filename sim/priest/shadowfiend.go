package priest

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (priest *Priest) registerShadowfiendSpell() {
	actionID := core.ActionID{SpellID: 401977}
	duration := time.Second * 15
	cooldown := time.Minute * 5

	// For timeline only
	priest.ShadowfiendAura = priest.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Shadowfiend",
		Duration: duration,
	})

	priest.Shadowfiend = priest.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_PriestShadowFiend,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,

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
			priest.ShadowfiendPet.EnableWithTimeout(sim, priest.ShadowfiendPet, duration)
			priest.ShadowfiendAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    priest.Shadowfiend,
		Priority: 1,
		Type:     core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentManaPercent() <= 0.01
		},
	})
}

type Shadowfiend struct {
	core.Pet

	Priest          *Priest
	Shadowcrawl     *core.Spell
	ShadowcrawlAura *core.Aura
}

func (priest *Priest) NewShadowfiend() *Shadowfiend {
	baseDamageMin := 0.0
	baseDamageMax := 0.0
	baseStats := stats.Stats{}
	// Seems to basically be a reskinned Felhunter so using Felhunter stats
	switch priest.Level {
	case 25:
		baseStats = stats.Stats{
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
		baseDamageMin = 24
		baseDamageMax = 40
	case 40:
		// TODO: All of the stats and stat inheritance needs to be verified
		baseStats = stats.Stats{
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
		baseDamageMin = 24
		baseDamageMax = 40
	case 50:
		baseStats = stats.Stats{
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
		baseDamageMin = 24
		baseDamageMax = 40
	case 60:
		baseStats = stats.Stats{
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
		baseDamageMin = 70
		baseDamageMax = 97
	}

	shadowfiend := &Shadowfiend{
		Pet:    core.NewPet("Shadowfiend", &priest.Character, baseStats, priest.shadowfiendStatInheritance(), false, true),
		Priest: priest,
	}

	manaMetric := priest.NewManaMetrics(core.ActionID{SpellID: 34433})
	core.MakePermanent(shadowfiend.GetOrRegisterAura(core.Aura{
		Label: "Autoattack mana regen",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			restoreMana := priest.MaxMana() * 0.05
			priest.AddMana(sim, restoreMana, manaMetric)
		},
	}))

	shadowfiend.EnableManaBarWithModifier(.77)

	shadowfiend.registerShadowCrawlSpell()

	shadowfiend.PseudoStats.DamageTakenMultiplier *= 0.1

	shadowfiend.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	shadowfiend.AddStat(stats.AttackPower, -20)

	// Warrior crit scaling
	shadowfiend.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[proto.Class_ClassWarrior][int(shadowfiend.Level)]*core.CritRatingPerCritChance)
	shadowfiend.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[proto.Class_ClassWarrior][int(shadowfiend.Level)]*core.SpellCritRatingPerCritChance)

	shadowfiend.EnableAutoAttacks(shadowfiend, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:     baseDamageMin,
			BaseDamageMax:     baseDamageMax,
			SwingSpeed:        1.5,
			AttackPowerPerDPS: 14.0 / 6.0, // Observed 6 times stronger AP scaling then expected
			SpellSchool:       core.SpellSchoolShadow,
			MaxRange:          core.MaxMeleeAttackRange,
		},
		AutoSwingMelee: true,
	})

	shadowfiend.ApplyOnPetEnable(func(sim *core.Simulation) {
		// Priest pets only inherit the owner's cast speed
		shadowfiend.EnableDynamicCastSpeedInheritance(sim)
	})

	priest.AddPet(shadowfiend)

	return shadowfiend
}

func (priest *Priest) shadowfiendStatInheritance() core.PetStatInheritance {
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

func (shadowfiend *Shadowfiend) Initialize() {
}

func (shadowfiend *Shadowfiend) ExecuteCustomRotation(sim *core.Simulation) {
	shadowfiend.Shadowcrawl.Cast(sim, shadowfiend.CurrentTarget)
}

func (shadowfiend *Shadowfiend) Reset(sim *core.Simulation) {
	shadowfiend.ShadowcrawlAura.Deactivate(sim)
	shadowfiend.Disable(sim)
}

func (shadowfiend *Shadowfiend) OnPetDisable(sim *core.Simulation) {
	shadowfiend.ShadowcrawlAura.Deactivate(sim)
}

func (shadowfiend *Shadowfiend) GetPet() *core.Pet {
	return &shadowfiend.Pet
}

func (shadowfiend *Shadowfiend) registerShadowCrawlSpell() {
	actionID := core.ActionID{SpellID: 401990}
	shadowfiend.ShadowcrawlAura = shadowfiend.GetOrRegisterAura(core.Aura{
		Label:    "Shadowcrawl",
		ActionID: actionID,
		Duration: time.Second * 5,
	}).AttachMultiplicativePseudoStatBuff(&shadowfiend.PseudoStats.DamageDealtMultiplier, 1.15)

	shadowfiend.Shadowcrawl = shadowfiend.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shadowfiend.ShadowcrawlAura.Activate(sim)
		},
	})
}
