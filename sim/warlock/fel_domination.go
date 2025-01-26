package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warlock *Warlock) registerFelDominationCD() {
	if !warlock.Talents.FelDomination {
		return
	}

	actionID := core.ActionID{SpellID: 18708}

	aura := warlock.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Fel Domination",
		Duration: time.Second * 15,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_WarlockSummons) {
				aura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: ClassSpellMask_WarlockSummons,
		IntValue:  -50,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		ClassMask: ClassSpellMask_WarlockSummons,
		TimeValue: -time.Millisecond * 5500,
	})

	spell := warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute * 15,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeUnknown,
	})
}
