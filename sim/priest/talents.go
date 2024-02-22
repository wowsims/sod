package priest

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (priest *Priest) ApplyTalents() {
	priest.applyInspiration()
	priest.applyShadowWeaving()
	priest.registerInnerFocus()

	// Meditation
	priest.PseudoStats.SpiritRegenRateCasting = []float64{0.0, 0.17, 0.33, 0.5}[priest.Talents.Meditation]

	// Spell Warding
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 1 - .02*float64(priest.Talents.SpellWarding)

	priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1 + (0.2 * float64(priest.Talents.Darkness))

	if priest.Talents.Shadowform {
		priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.15
	}

	if priest.Talents.SpiritualGuidance > 0 {
		priest.AddStatDependency(stats.Spirit, stats.SpellPower, 0.05*float64(priest.Talents.SpiritualGuidance))
	}

	if priest.Talents.MentalStrength > 0 {
		priest.MultiplyStat(stats.Intellect, 1.0+0.02*float64(priest.Talents.MentalStrength))
	}

	if priest.Talents.ImprovedPowerWordFortitude > 0 {
		priest.MultiplyStat(stats.Stamina, 1.0+.15*float64(priest.Talents.ImprovedPowerWordFortitude))
	}

	if priest.Talents.SilentResolve > 0 {
		priest.PseudoStats.ThreatMultiplier *= 1 - (.04 * float64(priest.Talents.SilentResolve))
	}
}

func (priest *Priest) forceOfWillDamageModifier() float64 {
	return 1 + .01*float64(priest.Talents.ForceOfWill)
}

func (priest *Priest) forceOfWillCritRating() float64 {
	return 1 * float64(priest.Talents.ForceOfWill) * core.CritRatingPerCritChance
}

func (priest *Priest) searingLightDamageModifier() float64 {
	return 1 + 0.05*float64(priest.Talents.SearingLight)
}
func (priest *Priest) holySpecCritRating() float64 {
	return 1 * float64(priest.Talents.HolySpecialization) * core.CritRatingPerCritChance
}

func (priest *Priest) shadowHitModifier() float64 {
	return 2 * float64(priest.Talents.ShadowFocus) * core.SpellHitRatingPerHitChance
}

func (priest *Priest) shadowThreatModifier() float64 {
	return 1 - 0.08*float64(priest.Talents.ShadowAffinity)
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

func (priest *Priest) applyShadowWeaving() {
	if priest.Talents.ShadowWeaving == 0 {
		return
	}

	priest.ShadowWeavingAuras = priest.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
		return core.ShadowWeavingAura(unit, int(priest.Talents.ShadowWeaving))
	})
}

func (priest *Priest) AddShadowWeavingStack(sim *core.Simulation, target *core.Unit) {
	if priest.ShadowWeavingAuras == nil {
		return
	}

	if sim.RollWithLabel(0, 1, "ShadowWeaving") < (0.2 * float64(priest.Talents.ShadowWeaving)) {
		priest.ShadowWeavingAuras.Get(target).Activate(sim)
		priest.ShadowWeavingAuras.Get(target).AddStack(sim)
	}
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
			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, 25*core.CritRatingPerCritChance)
			aura.Unit.PseudoStats.CostMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, -25*core.CritRatingPerCritChance)
			aura.Unit.PseudoStats.CostMultiplier += 1
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			priest.InnerFocus.CD.Use(sim)
			priest.UpdateMajorCooldowns()
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
}
