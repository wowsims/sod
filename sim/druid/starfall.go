package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// We register two spells to apply two different dot effects and get two entries in Damage/Detailed results
func (druid *Druid) registerStarfallCD() {
	if !druid.HasRune(proto.DruidRune_RuneCloakStarfall) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.DruidRune_RuneCloakStarfall)}

	// Moonfury and Improved Moonfire only seem to be applying the crit bonus to starfall at the moment in-game
	// talentBaseMultiplier := 1 + druid.MoonfuryDamageMultiplier() + druid.ImprovedMoonfireDamageMultiplier()

	baseDamageLow := druid.baseRuneAbilityDamage() * 0.46   // * talentBaseMultiplier
	baseDamageHigh := druid.baseRuneAbilityDamage() * .54   // * talentBaseMultiplier
	baseDamageSplash := druid.baseRuneAbilityDamage() * .08 // * talentBaseMultiplier
	spellCoefTick := .3
	spellCoefSplash := .127

	numberOfTicks := core.TernaryInt32(druid.Env.GetNumTargets() > 1, 20, 10)
	tickLength := time.Second
	cooldown := time.Second * 90

	druid.StarfallSplash = druid.RegisterSpell(Any, core.SpellConfig{
		SpellCode:   SpellCode_DruidStarfallSplash,
		ActionID:    actionID.WithTag(2),
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskEmpty,

		BonusCritRating:  druid.ImprovedMoonfireCritBonus(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoefSplash,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamageSplash, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	druid.StarfallTick = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		SpellCode:   SpellCode_DruidStarfallTick,
		ActionID:    actionID.WithTag(1),
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage, // Shown to proc things in-game
		Flags:       core.SpellFlagBinary,

		BonusCritRating:  druid.ImprovedMoonfireCritBonus(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoefTick,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			druid.StarfallSplash.Cast(sim, target)
		},
	})

	druid.Starfall = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		SpellCode:   SpellCode_DruidStarfall,
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL | SpellFlagOmen,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.39,
			Multiplier: druid.MoonglowManaCostMultiplier(),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: cooldown,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Starfall",
			},
			NumberOfTicks: numberOfTicks,
			TickLength:    tickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				druid.StarfallTick.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Starfall.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
