package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (shaman *Shaman) registerWaterShieldSpell() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsWaterShield) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.ShamanRune_RuneHandsWaterShield)}

	passiveMP5Pct := .01
	onHitManaGainedPct := .04

	manaMetrics := shaman.NewManaMetrics(actionID)
	mp5StatDep := shaman.NewDynamicStatDependency(stats.Mana, stats.MP5, passiveMP5Pct)

	aura := shaman.RegisterAura(core.Aura{
		Label:    "Water Shield",
		ActionID: actionID,
		Duration: time.Minute * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, mp5StatDep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, mp5StatDep)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.ProcMask.Matches(core.ProcMaskDirect) {
				shaman.AddMana(sim, shaman.MaxMana()*onHitManaGainedPct, manaMetrics)
			}
		},
	})

	shaman.WaterShield = shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		ProcMask: core.ProcMaskEmpty,
		Flags:    SpellFlagShaman | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if shaman.ActiveShieldAura != nil {
				shaman.ActiveShieldAura.Deactivate(sim)
			}
			shaman.ActiveShield = spell
			shaman.ActiveShieldAura = aura
			shaman.ActiveShieldAura.Activate(sim)
		},
	})
}
