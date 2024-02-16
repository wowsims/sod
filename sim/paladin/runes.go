package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (paladin *Paladin) ApplyRunes() {
	paladin.registerTheArtOfWar()
	paladin.registerSheathOfLight()
	paladin.registerGuardedByTheLight()
}

func (paladin *Paladin) registerTheArtOfWar() {
	if !paladin.HasRune(proto.PaladinRune_RuneFeetTheArtOfWar) {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label:    "The Art of War",
		Duration: core.NeverExpires,
		ActionID: core.ActionID{SpellID: 426157},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			paladin.HolyShockCooldown.Reset()
			paladin.ExorcismCooldown.Reset()
		},
	})
}

func (paladin *Paladin) registerSheathOfLight() {

	if !paladin.HasRune(proto.PaladinRune_RuneWaistSheathOfLight) {
		return
	}

	dep := paladin.NewDynamicStatDependency(
		stats.AttackPower, stats.SpellDamage, 0.3,
	)

	sheathAura := paladin.RegisterAura(core.Aura{
		Label:    "Sheath of Light",
		Duration: time.Second * 60,
		ActionID: core.ActionID{SpellID: 426159},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.EnableDynamicStatDep(sim, dep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.DisableDynamicStatDep(sim, dep)
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "Sheath of Light (rune)",
		Duration: core.NeverExpires,
		ActionID: core.ActionID{SpellID: 426158},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskWhiteHit) {
				return
			}
			sheathAura.Activate(sim)
		},
	})

}

func (paladin *Paladin) registerGuardedByTheLight() {

	if !paladin.HasRune(proto.PaladinRune_RuneFeetGuardedByTheLight) {
		return
	}

	actionID := core.ActionID{SpellID: 415058}
	manaMetrics := paladin.NewManaMetrics(actionID)
	var manaPA *core.PendingAction

	guardedAura := paladin.RegisterAura(core.Aura{
		Label:    "Guarded by the Light",
		Duration: time.Second*15 + 1,
		ActionID: actionID,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			manaPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 3,
				OnAction: func(sim *core.Simulation) {
					paladin.AddMana(sim, 0.05*paladin.MaxMana(), manaMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			manaPA.Cancel(sim)
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "Guarded by the Light (rune)",
		Duration: core.NeverExpires,
		ActionID: core.ActionID{SpellID: 415755},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskWhiteHit) {
				return
			}
			guardedAura.Activate(sim)
		},
	})
}
