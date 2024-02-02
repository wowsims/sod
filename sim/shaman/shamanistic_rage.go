package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
	"time"
)

func (shaman *Shaman) applyShamanisticRage() {
	if !shaman.HasRune(proto.ShamanRune_RuneLegsShamanisticRage) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsShamanisticRage)}
	manaMetrics := shaman.NewManaMetrics(actionID)
	srAura := shaman.RegisterAura(core.Aura{
		Label:    "Shamanistic Rage",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.8
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.8
		},
	})

	spell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 1,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			srAura.Activate(sim)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 15,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					mana := max(
						shaman.GetStat(stats.AttackPower)*0.15,
						shaman.GetStat(stats.SpellPower)*0.10,
						shaman.GetStat(stats.Healing)*0.06,
					)
					shaman.AddMana(sim, mana, manaMetrics)
				},
			})
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentManaPercent() <= 0.2
		},
	})
}
