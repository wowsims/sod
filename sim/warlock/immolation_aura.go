package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Immolation Aura now triggers from being attacked rather than as a periodic effect. This cannot occur more than once per second, and does not require the attack to hit.
// Immolation Aura now also increases fire damage by 10%.
func (warlock *Warlock) registerImmolationAuraSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneBracerImmolationAura) {
		return
	}

	spellCoeff := 0.045
	baseDamage := warlock.baseRuneAbilityDamage() * 0.2

	immoAuraProc := warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 427725},
		ClassSpellMask: ClassSpellMask_WarlockImmolationAuraProc,
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeAlwaysHit)
			}
		},
	})

	var pa *core.PendingAction
	immoAura := warlock.RegisterAura(core.Aura{
		Label:    "Immolation Aura",
		ActionID: core.ActionID{SpellID: int32(proto.WarlockRune_RuneBracerImmolationAura)},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			pa = core.NewPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 2,
				OnAction: func(s *core.Simulation) {
					immoAuraProc.Cast(sim, warlock.CurrentTarget)
				},
			})
			// Dont proc damage in prepull
			if pa.NextActionAt < 0 {
				pa.NextActionAt = 0
			}
			sim.AddPendingAction(pa)

			for si := stats.SchoolIndexArcane; si < stats.SchoolLen; si++ {
				warlock.PseudoStats.SchoolDamageTakenMultiplier[si] *= 0.9
			}

			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.10
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			pa.Cancel(sim)

			for si := stats.SchoolIndexArcane; si < stats.SchoolLen; si++ {
				warlock.PseudoStats.SchoolDamageTakenMultiplier[si] /= 0.9
			}

			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= 1.10
		},
	})

	core.MakeProcTriggerAura(&warlock.Unit, core.ProcTrigger{
		Name:     "Immolation Aura Trigger",
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskSpellDamage,
		ICD:      time.Second * 1,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			immoAura.Activate(sim)
		},
	})

	warlock.ImmolationAura = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: int32(proto.WarlockRune_RuneBracerImmolationAura)},
		ClassSpellMask: ClassSpellMask_WarlockImmolationAura,
		SpellSchool:    core.SpellSchoolFire,
		Flags:          core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagDestruction,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if immoAura.IsActive() {
				immoAura.Deactivate(sim)
			} else {
				immoAura.Activate(sim)
			}
		},
	})
}
