package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) registerImmolationAuraSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneBracerImmolationAura) {
		return
	}
	spellCoeff := 0.045

	level := float64(warlock.GetCharacter().Level)
	baseCalc := (6.568597 + 0.672028*level + 0.031721*level*level)
	baseDamage := baseCalc * 0.2

	immoAuraProc := warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 427725},
		SpellSchool: core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,

		BonusCritRating: float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,

		CritDamageBonus: warlock.ruin(),

		DamageMultiplierAdditive: 1, // + 0.02*float64(warlock.Talents.Emberstorm), Not affected by any talent atm
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				targetDamage := baseDamage + spellCoeff*spell.SpellDamage()

				// TODO: Right now Immolation Aura is not listed in the Lake of Fire debuff affected spells
				// if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(aoeTarget).IsActive() {
				// 	targetDamage *= 1.4
				// }

				spell.CalcAndDealDamage(sim, aoeTarget, targetDamage, spell.OutcomeAlwaysHit)
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

			warlock.PseudoStats.DamageTakenMultiplier *= 0.9
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			pa.Cancel(sim)
			warlock.PseudoStats.DamageTakenMultiplier /= 0.9
		},
	})

	warlock.ImmolationAura = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.WarlockRune_RuneBracerImmolationAura)},
		SpellSchool: core.SpellSchoolFire,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing | core.SpellFlagNoOnCastComplete,

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
