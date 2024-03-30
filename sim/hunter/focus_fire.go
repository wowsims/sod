package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) registerFocusFireSpell() {
	if hunter.pet == nil || !hunter.HasRune(proto.HunterRune_RuneBracersFocusFire) {
		return
	}

	focusFireMetrics := hunter.pet.Metrics.NewResourceMetrics(core.ActionID{SpellID: 428726}, proto.ResourceType_ResourceTypeEnergy)
	focusFireActionId := core.ActionID{SpellID: 428726}
	focusFireFrenzyActionId := core.ActionID{SpellID: 428728}

	// For tracking in timeline
	hunterFrenzyAura := hunter.RegisterAura(core.Aura{
		Label:     "Focus Fire Frenzy (Hunter)",
		ActionID:  focusFireFrenzyActionId,
		Duration:  time.Second * 10,
		MaxStacks: 5,
	})

	hunterPetFrenzyAura := hunter.pet.RegisterAura(core.Aura{
		Label:     "Focus Fire Frenzy",
		ActionID:  focusFireFrenzyActionId,
		Duration:  time.Second * 10,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1/(1+(0.06*float64(oldStacks))))
			aura.Unit.MultiplyMeleeSpeed(sim, 1+(0.06*float64(newStacks)))
			if !hunterFrenzyAura.IsActive() { hunterFrenzyAura.Activate(sim) }
			hunterFrenzyAura.SetStacks(sim, newStacks)
		},
	})

	hunterFocusFireAura := hunter.RegisterAura(core.Aura{
		Label:    "Focus Fire",
		ActionID: focusFireActionId,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			frenzyStacks := hunter.pet.GetAuraByID(focusFireFrenzyActionId).GetStacks()
			aura.Unit.MultiplyRangedSpeed(sim, 1+(0.03*float64(frenzyStacks)))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			frenzyStacks := hunter.pet.GetAuraByID(focusFireFrenzyActionId).GetStacks()
			aura.Unit.MultiplyRangedSpeed(sim, 1/(1+(0.03*float64(frenzyStacks))))
		},
	})

	hunter.pet.RegisterAura(core.Aura{
		Label:    "Focus Fire Pet",
		ActionID: focusFireActionId,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto) {
				if !hunterPetFrenzyAura.IsActive() {
					hunterPetFrenzyAura.Activate(sim)
				}

				hunterPetFrenzyAura.AddStack(sim)
				hunterPetFrenzyAura.Refresh(sim)
				hunterFrenzyAura.Refresh(sim)
			}
		},
	})

	hunter.FocusFire = hunter.RegisterSpell(core.SpellConfig{
		ActionID: focusFireActionId,

		ManaCost: core.ManaCostOptions{
			FlatCost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 15,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.pet.IsEnabled() && (hunter.pet.GetAuraByID(focusFireFrenzyActionId).GetStacks() > 0)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			frenzyStacks := hunterPetFrenzyAura.GetStacks()
			hunter.pet.AddFocus(sim, float64(4 * frenzyStacks), focusFireMetrics)
			hunterFocusFireAura.Activate(sim)
			hunterPetFrenzyAura.Deactivate(sim)
		},
	})
}
