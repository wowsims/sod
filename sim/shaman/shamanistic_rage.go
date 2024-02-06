package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (shaman *Shaman) applyShamanisticRage() {
	if !shaman.HasRune(proto.ShamanRune_RuneLegsShamanisticRage) {
		return
	}

	apCoeff := .15
	spCoeff := .10
	hpCoeff := .06
	damageTakenMultiplier := .8
	duration := time.Second * 15
	cooldown := time.Minute * 1

	actionID := core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsShamanisticRage)}
	manaMetrics := shaman.NewManaMetrics(actionID)
	srAura := shaman.RegisterAura(core.Aura{
		Label:    "Shamanistic Rage",
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMultiplier
		},
	})

	spell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cooldown,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			// TODO: Raid mana regain
			srAura.Activate(sim)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 15,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					mana := max(
						shaman.GetStat(stats.AttackPower)*apCoeff,
						shaman.GetStat(stats.SpellPower)*spCoeff,
						shaman.GetStat(stats.Healing)*hpCoeff,
					)
					shaman.AddMana(sim, mana, manaMetrics)
				},
			})
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeMana,
	})
}
