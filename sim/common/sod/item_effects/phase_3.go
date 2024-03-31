package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	AtalaiBloodRitualCharm = 220634
	CobraFangClaw          = 220588
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Cloth
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Leather
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Mail
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Plate
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	core.NewItemEffect(AtalaiBloodRitualCharm, func(agent core.Agent) {
		character := agent.GetCharacter()

		bonusPerStack := stats.Stats{
			stats.SpellDamage:  8,
			stats.HealingPower: 16,
		}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Unrestrained Power",
			ActionID:  core.ActionID{SpellID: 446297},
			Duration:  time.Second * 20,
			MaxStacks: 12,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.SetStacks(sim, aura.MaxStacks)
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				bonusStats := bonusPerStack.Multiply(float64(newStacks - oldStacks))
				character.AddStatsDynamic(sim, bonusStats)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
					return
				}
				aura.RemoveStack(sim)
			},
		})

		triggerSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 446297},
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    triggerSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////
	core.NewItemEffect(CobraFangClaw, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(CobraFangClaw)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		character.RegisterAura(core.Aura{
			Label:    "Cobra Fang Claw Thrash",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}
				if !spell.ProcMask.Matches(procMask) {
					return
				}

				if ppmm.Proc(sim, procMask, "Cobra Fang Claw Extra Attack") {
					character.AutoAttacks.ExtraMHAttack(sim)
				}
			},
		})
	})

	core.AddEffectsToTest = true
}
