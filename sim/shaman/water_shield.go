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

	shaman.WaterShieldAura = shaman.RegisterAura(core.Aura{
		Label:     "Water Shield",
		ActionID:  actionID,
		Duration:  time.Minute * 10,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, mp5StatDep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, mp5StatDep)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.ProcMask.Matches(core.ProcMaskDirect) {
				shaman.WaterShieldRestore.Cast(sim, aura.Unit)
				aura.RemoveStack(sim)

				if aura.GetStacks() == 0 {
					aura.Deactivate(sim)
				}
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
			shaman.ActiveShieldAura = shaman.WaterShieldAura
			shaman.WaterShieldAura.Activate(sim)
			shaman.WaterShieldAura.SetStacks(sim, shaman.WaterShieldAura.MaxStacks)
		},
	})

	shaman.WaterShieldRestore = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 408511},
		SpellSchool: core.SpellSchoolNature,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoMetrics,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			shaman.AddMana(sim, shaman.MaxMana()*onHitManaGainedPct, manaMetrics)
		},
	})
}
