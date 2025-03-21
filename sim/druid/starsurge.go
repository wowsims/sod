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

	starfireAuraDuration := time.Second * 15

	damageMod := druid.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_DruidStarfire,
		IntValue:  80,
	})

	druid.StarsurgeAura = druid.RegisterAura(core.Aura{
		Label:     "Starsurge",
		ActionID:  actionID,
		Duration:  starfireAuraDuration,
		MaxStacks: 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_DruidStarfire) && aura.GetStacks() > 0 {
				aura.RemoveStack(sim)
			}
		},
	})

	druid.Starsurge = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ClassSpellMask: ClassSpellMask_DruidStarsurge,
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolArcane,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagOmen | core.SpellFlagResetAttackSwing | core.SpellFlagBinary | core.SpellFlagAPL,

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

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			// NG procs when the cast finishes
			if result.DidCrit() && druid.NaturesGraceProcAura != nil {
				druid.NaturesGraceProcAura.Activate(sim)
				druid.NaturesGraceProcAura.SetStacks(sim, druid.NaturesGraceProcAura.MaxStacks)
			}

			// Aura applies on cast
			druid.StarsurgeAura.Activate(sim)
			druid.StarsurgeAura.SetStacks(sim, druid.StarsurgeAura.MaxStacks)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
