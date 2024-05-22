package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (druid *Druid) applyStarsurge() {
	if !druid.HasRune(proto.DruidRune_RuneLegsStarsurge) {
		return
	}

	actionID := core.ActionID{SpellID: 417157}

	baseLowDamage := druid.baseRuneAbilityDamage() * 2.48
	baseHighDamage := druid.baseRuneAbilityDamage() * 3.04
	spellCoeff := .429

	starfireAuraMultiplier := 1 + .80
	starfireAuraDuration := time.Second * 15

	starfireDamageAura := druid.RegisterAura(core.Aura{
		Label:     "Starsurge",
		ActionID:  actionID,
		Duration:  starfireAuraDuration,
		MaxStacks: 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(druid.Starfire, func(spell *DruidSpell) {
				if spell != nil {
					spell.DamageMultiplier *= starfireAuraMultiplier
				}
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(druid.Starfire, func(spell *DruidSpell) {
				if spell != nil {
					spell.DamageMultiplier /= starfireAuraMultiplier
				}
			})
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode != SpellCode_DruidStarfire {
				return
			}

			aura.Deactivate(sim)
		},
	})

	druid.Starsurge = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    actionID,
		SpellCode:   SpellCode_DruidStarsurge,
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagOmen | core.SpellFlagResetAttackSwing | core.SpellFlagBinary | core.SpellFlagAPL,

		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.01 * (1 - 0.03*float64(druid.Talents.Moonglow)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 0,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		CritDamageBonus: druid.vengeance(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) * druid.MoonfuryDamageMultiplier()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.DidCrit() && druid.NaturesGraceProcAura != nil {
				druid.NaturesGraceProcAura.Activate(sim)
			}

			// Aura applies on cast
			starfireDamageAura.Activate(sim)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
