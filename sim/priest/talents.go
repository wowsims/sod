package priest

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (priest *Priest) ApplyTalents() {
	// Discipline
	priest.registerInnerFocus()
	priest.applyMentalAgility()
	priest.applyForceOfWill()

	if priest.Talents.SilentResolve > 0 {
		priest.PseudoStats.ThreatMultiplier *= 1 - (.04 * float64(priest.Talents.SilentResolve))
	}

	if priest.Talents.ImprovedPowerWordFortitude > 0 {
		priest.MultiplyStat(stats.Stamina, 1.0+.15*float64(priest.Talents.ImprovedPowerWordFortitude))
	}

	priest.PseudoStats.SpiritRegenRateCasting = []float64{0.0, 0.17, 0.33, 0.5}[priest.Talents.Meditation]

	if priest.Talents.MentalStrength > 0 {
		priest.MultiplyStat(stats.Intellect, 1.0+0.02*float64(priest.Talents.MentalStrength))
	}

	// Holy
	priest.applyInspiration()
	priest.applyHolySpecialization()
	priest.applySearingLight()

	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 1 - .02*float64(priest.Talents.SpellWarding)

	if priest.Talents.SpiritualGuidance > 0 {
		priest.AddStatDependency(stats.Spirit, stats.SpellPower, 0.05*float64(priest.Talents.SpiritualGuidance))
	}

	// Shadow
	priest.registerVampiricEmbraceSpell()
	priest.registerShadowform()
	priest.applyShadowAffinity()
	priest.applyShadowFocus()
	priest.applyShadowWeaving()
}

func (priest *Priest) darknessDamageModifier() float64 {
	return 1 + .02*float64(priest.Talents.Darkness)
}

func (priest *Priest) applyMentalAgility() {
	if priest.Talents.MentalAgility == 0 {
		return
	}

	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) && spell.DefaultCast.CastTime == 0 {
			spell.CostMultiplier *= 1 - .02*float64(priest.Talents.MentalAgility)
		}
	})
}

func (priest *Priest) applyForceOfWill() {
	if priest.Talents.ForceOfWill == 0 {
		return
	}

	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) {
			spell.DamageMultiplier *= .01 * float64(priest.Talents.ForceOfWill)
			spell.BonusCritRating += 1 * float64(priest.Talents.ForceOfWill) * core.CritRatingPerCritChance
		}
	})
}

func (priest *Priest) applyHolySpecialization() {
	if priest.Talents.HolySpecialization == 0 {
		return
	}

	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) && spell.SpellSchool.Matches(core.SpellSchoolHoly) {
			spell.BonusCritRating += 1 * float64(priest.Talents.HolySpecialization) * core.CritRatingPerCritChance
		}
	})
}

func (priest *Priest) applyInspiration() {
	if priest.Talents.Inspiration == 0 {
		return
	}

	auras := make([]*core.Aura, len(priest.Env.AllUnits))
	for _, unit := range priest.Env.AllUnits {
		if !priest.IsOpponent(unit) {
			aura := core.InspirationAura(unit, priest.Talents.Inspiration)
			auras[unit.UnitIndex] = aura
		}
	}

	priest.RegisterAura(core.Aura{
		Label:    "Inspiration Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if slices.Contains([]int32{SpellCode_PriestFlashHeal, SpellCode_PriestHeal, SpellCode_PriestGreaterHeal}, spell.SpellCode) {
				auras[result.Target.UnitIndex].Activate(sim)
			}
		},
	})
}

func (priest *Priest) applySearingLight() {
	if priest.Talents.SearingLight == 0 {
		return
	}

	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellCode == SpellCode_PriestSmite || spell.SpellCode == SpellCode_PriestHolyFire {
			spell.DamageMultiplier *= 0.05 * float64(priest.Talents.SearingLight)
		}
	})
}

func (priest *Priest) applyShadowAffinity() {
	if priest.Talents.ShadowAffinity == 0 {
		return
	}

	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) || spell.SpellSchool.Matches(core.SpellSchoolShadow) {
			spell.ThreatMultiplier *= 1 - 0.08*float64(priest.Talents.ShadowAffinity)
		}
	})
}

func (priest *Priest) applyShadowFocus() {
	if priest.Talents.ShadowFocus == 0 {
		return
	}

	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) || spell.SpellSchool.Matches(core.SpellSchoolShadow) {
			spell.BonusHitRating += 2 * float64(priest.Talents.ShadowFocus) * core.SpellHitRatingPerHitChance
		}
	})
}

func (priest *Priest) applyShadowWeaving() {
	if priest.Talents.ShadowWeaving == 0 {
		return
	}

	priest.ShadowWeavingAuras = priest.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		return core.ShadowWeavingAura(unit, int(priest.Talents.ShadowWeaving))
	})

	procChance := 0.2 * float64(priest.Talents.ShadowWeaving)

	priest.ShadowWeavingProc = priest.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 15258},
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagNoMetrics,
		SpellSchool: core.SpellSchoolShadow,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if !result.Landed() {
				return
			}

			if procChance == 1.0 || sim.RollWithLabel(0, 1, "ShadowWeaving") < procChance {
				priest.ShadowWeavingAuras.Get(target).Activate(sim)
				priest.ShadowWeavingAuras.Get(target).AddStack(sim)
			}
		},
	})
}

func (priest *Priest) AddShadowWeavingStack(sim *core.Simulation, target *core.Unit) {
	if priest.ShadowWeavingProc == nil {
		return
	}

	priest.ShadowWeavingProc.Cast(sim, target)
}

func (priest *Priest) registerInnerFocus() {
	if !priest.Talents.InnerFocus {
		return
	}

	actionID := core.ActionID{SpellID: 14751}

	priest.InnerFocusAura = priest.RegisterAura(core.Aura{
		Label:    "Inner Focus",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.CostMultiplier -= 1
			for _, spell := range priest.Spellbook {
				if spell.Flags.Matches(SpellFlagPriest) {
					spell.CostMultiplier -= 1
					spell.BonusCritRating += 25 * core.SpellCritRatingPerCritChance
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.CostMultiplier += 1
			for _, spell := range priest.Spellbook {
				if spell.Flags.Matches(SpellFlagPriest) {
					spell.CostMultiplier += 1
					spell.BonusCritRating -= 25 * core.SpellCritRatingPerCritChance
				}
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Flags.Matches(SpellFlagPriest) {
				// Remove the buff and put skill on CD
				aura.Deactivate(sim)
				priest.InnerFocus.CD.Use(sim)
				priest.UpdateMajorCooldowns()
			}
		},
	})

	priest.InnerFocus = priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			priest.InnerFocusAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell: priest.InnerFocus,
		Type:  core.CooldownTypeDPS,
	})
}

func (priest *Priest) registerShadowform() {
	if !priest.Talents.Shadowform {
		return
	}

	actionID := core.ActionID{SpellID: 15473}

	priest.ShadowformAura = priest.RegisterAura(core.Aura{
		Label:    "Shadowform",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.25
			for _, spell := range priest.Spellbook {
				if spell.SpellSchool.Matches(core.SpellSchoolShadow) {
					spell.CostMultiplier *= .5
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.25
			for _, spell := range priest.Spellbook {
				if spell.SpellSchool.Matches(core.SpellSchoolShadow) {
					spell.CostMultiplier /= .5
				}
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolHoly) {
				aura.Deactivate(sim)
			}
		},
	})

	priest.Shadowform = priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			priest.ShadowformAura.Activate(sim)
		},
	})
}
