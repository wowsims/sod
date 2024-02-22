package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

type Shadowfiend struct {
	core.Pet

	Priest          *Priest
	Shadowcrawl     *core.Spell
	ShadowcrawlAura *core.Aura
}

func (priest *Priest) NewShadowfiend() *Shadowfiend {
	baseDamageMin := 0.0
	baseDamageMax := 0.0
	shadowfiendBaseStats := stats.Stats{}
	switch priest.Level {
	case 25:
	case 40:
		// 40 stats
		// TODO: All of the stats and stat inheritance needs to be verified
		baseDamageMin = 140
		baseDamageMax = 151
		shadowfiendBaseStats = stats.Stats{
			stats.Strength:    74,
			stats.Agility:     59,
			stats.Stamina:     181,
			stats.Intellect:   49,
			stats.Spirit:      97,
			stats.AttackPower: 0,
			// with 3% crit debuff, shadowfiend crits around 9-12% (TODO: verify and narrow down)
			stats.MeleeCrit: 8 * core.CritRatingPerCritChance,
		}
	case 50:
	case 60:
	}

	shadowfiend := &Shadowfiend{
		Pet:    core.NewPet("Shadowfiend", &priest.Character, shadowfiendBaseStats, priest.shadowfiendStatInheritance(), false, false),
		Priest: priest,
	}

	manaMetric := priest.NewManaMetrics(core.ActionID{SpellID: 34433})
	_ = core.MakePermanent(shadowfiend.GetOrRegisterAura(core.Aura{
		Label: "Autoattack mana regen",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			restoreMana := priest.MaxMana() * 0.05
			priest.AddMana(sim, restoreMana, manaMetric)
		},
	}))

	actionID := core.ActionID{SpellID: 401990}
	shadowfiend.ShadowcrawlAura = shadowfiend.GetOrRegisterAura(core.Aura{
		Label:    "Shadowcrawl",
		ActionID: actionID,
		Duration: time.Second * 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shadowfiend.PseudoStats.DamageDealtMultiplier *= 1.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shadowfiend.PseudoStats.DamageDealtMultiplier /= 1.15
		},
	})

	shadowfiend.Shadowcrawl = shadowfiend.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shadowfiend.ShadowcrawlAura.Activate(sim)
		},
	})

	shadowfiend.PseudoStats.DamageTakenMultiplier *= 0.1

	shadowfiend.EnableAutoAttacks(shadowfiend, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:        baseDamageMin,
			BaseDamageMax:        baseDamageMax,
			SwingSpeed:           1.5,
			NormalizedSwingSpeed: 1.5,
			CritMultiplier:       2,
			SpellSchool:          core.SpellSchoolShadow,
		},
		AutoSwingMelee: true,
	})

	shadowfiend.AddStatDependency(stats.Strength, stats.AttackPower, 2.0)

	// core.ApplyPetConsumeEffects(&shadowfiend.Character, priest.Consumes)

	priest.AddPet(shadowfiend)

	return shadowfiend
}

func (priest *Priest) shadowfiendStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		// TODO Scaling
		return stats.Stats{}
	}
}

func (shadowfiend *Shadowfiend) Initialize() {
}

func (shadowfiend *Shadowfiend) ExecuteCustomRotation(sim *core.Simulation) {
	shadowfiend.Shadowcrawl.Cast(sim, nil)
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
