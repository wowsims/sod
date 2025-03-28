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

	focusFireMetrics := hunter.pet.NewFocusMetrics(core.ActionID{SpellID: int32(proto.HunterRune_RuneBracersFocusFire)})
	focusFireActionId := core.ActionID{SpellID: int32(proto.HunterRune_RuneBracersFocusFire)}
	focusFireFrenzyActionId := core.ActionID{SpellID: 428728}

	duration := time.Second * 10
	maxStacks := int32(5)

	// Ues a dummy aura for tracking on the timeline
	hunterFrenzyAura := hunter.RegisterAura(core.Aura{
		Label:     "Focus Fire Frenzy (Hunter)",
		ActionID:  focusFireFrenzyActionId,
		Duration:  duration,
		MaxStacks: maxStacks,
	})

	hunterPetFrenzyAura := hunter.pet.RegisterAura(core.Aura{
		Label:     "Focus Fire Frenzy",
		ActionID:  focusFireFrenzyActionId,
		Duration:  duration,
		MaxStacks: maxStacks,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1/(1+(0.06*float64(oldStacks))))
			aura.Unit.MultiplyMeleeSpeed(sim, 1+(0.06*float64(newStacks)))
		},
		OnRefresh: func(aura *core.Aura, sim *core.Simulation) {
			hunterFrenzyAura.SetStacks(sim, aura.GetStacks())
		},
	}).AttachDependentAura(hunterFrenzyAura)

	hunter.FocusFireAura = hunter.RegisterAura(core.Aura{
		Label:     "Focus Fire",
		ActionID:  focusFireActionId,
		Duration:  time.Second * 20,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			aura.Unit.MultiplyRangedSpeed(sim, 1/(1+(0.03*float64(oldStacks))))
			aura.Unit.MultiplyRangedSpeed(sim, 1+(0.03*float64(newStacks)))
		},
	})

	core.MakeProcTriggerAura(&hunter.pet.Unit, core.ProcTrigger{
		Name:           "Focus Fire Trigger",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: ClassSpellMask_HunterPetBasicAttacks,
		Harmful:        true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			hunterPetFrenzyAura.Activate(sim)
			hunterPetFrenzyAura.AddStack(sim)
		},
	})

	hunter.FocusFire = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       focusFireActionId,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: ClassSpellMask_HunterFocusFire,
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.pet.IsEnabled() && hunterPetFrenzyAura.GetStacks() > 0
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			frenzyStacks := hunterPetFrenzyAura.GetStacks()
			hunter.pet.AddFocus(sim, float64(4*frenzyStacks), focusFireMetrics)

			hunter.FocusFireAura.Activate(sim)
			hunter.FocusFireAura.SetStacks(sim, frenzyStacks)
			hunterPetFrenzyAura.Deactivate(sim)
		},
	})
}
