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

	baseDamageLow := druid.baseRuneAbilityDamage() * 0.46
	baseDamageHigh := druid.baseRuneAbilityDamage() * 0.54
	baseDamageSplash := druid.baseRuneAbilityDamage() * 0.08
	spellCoefTick := 0.3
	spellCoefSplash := 0.127

	numberOfTicks := core.TernaryInt32(druid.Env.GetNumTargets() > 1, 20, 10)
	tickLength := time.Second
	cooldown := time.Second * 90

	druid.StarfallSplash = druid.RegisterSpell(Any, core.SpellConfig{
		SpellCode:   SpellCode_DruidStarfallSplash,
		ActionID:    actionID.WithTag(2),
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoefSplash,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Apply the base spell's multipliers to pick up on effects that only affect spells with DoTs
			spell.DamageMultiplierAdditive += druid.Starfall.PeriodicDamageMultiplierAdditive - 1

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamageSplash, spell.OutcomeMagicHitAndCrit)
			}

			spell.DamageMultiplierAdditive -= druid.Starfall.PeriodicDamageMultiplierAdditive - 1
		},
	})

	druid.StarfallTick = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		SpellCode:   SpellCode_DruidStarfallTick,
		ActionID:    actionID.WithTag(1),
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage, // Shown to proc things in-game
		Flags:       core.SpellFlagBinary,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoefTick,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)

			// Apply the base spell's multipliers to pick up on effects that only affect spells with DoTs
			spell.DamageMultiplierAdditive += druid.Starfall.PeriodicDamageMultiplierAdditive - 1
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DamageMultiplierAdditive -= druid.Starfall.PeriodicDamageMultiplierAdditive - 1

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
			BaseCost: 0.39,
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

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Starfall.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
