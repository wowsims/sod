package crafted

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	echoesOfInsanityEffect := func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Echoes of Insanity Proc", core.ActionID{SpellID: 446541}, stats.Stats{stats.HealingPower: 50}, time.Second*10)

		handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			procAura.Activate(sim)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 446545},
			Name:       "Echoes of Insanity",
			Callback:   core.CallbackOnHealDealt,
			ProcMask:   core.ProcMaskSpellHealing,
			ProcChance: 0.3,
			ICD:        time.Second * 40,
			Handler:    handler,
		})
	}

	echoesOfMadnessEffect := func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Echoes of Madness Proc",
			ActionID: core.ActionID{SpellID: 446528},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyCastSpeed(1.1)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyCastSpeed(1 / 1.1)
			},
		})

		handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			procAura.Activate(sim)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 446518},
			Name:       "Echoes of Madness",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.3,
			ICD:        time.Second * 40,
			Handler:    handler,
		})
	}

	echoesOfDreadEffect := func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Echoes of Dread Proc",
			ActionID: core.ActionID{SpellID: 446577},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.05)
				character.AddStatDynamic(sim, stats.AttackPower, 50)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1/1.05)
				character.AddStatDynamic(sim, stats.AttackPower, -50)
			},
		})

		handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			procAura.Activate(sim)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 446579},
			Name:       "Echoes of Dread",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskMelee,
			ProcChance: 0.3,
			ICD:        time.Second * 40,
			Handler:    handler,
		})
	}

	///////////////////////////////////////////////////////////////////////////
	//                                 Cloth
	///////////////////////////////////////////////////////////////////////////

	// Mantle of Insanity
	core.NewItemEffect(220749, echoesOfInsanityEffect)

	// Fractured Mind Pauldrons
	core.NewItemEffect(220750, echoesOfMadnessEffect)

	// Shoulderpads of the Deranged
	core.NewItemEffect(220751, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Echoes of the Depraved Proc", core.ActionID{SpellID: 446572}, stats.Stats{stats.SpellDamage: 30, stats.Dodge: 2}, time.Second*10)

		handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			procAura.Activate(sim)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 446570},
			Name:       "Echoes of the Depraved",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.3,
			ICD:        time.Second * 40,
			Handler:    handler,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Leather
	///////////////////////////////////////////////////////////////////////////

	// Membrane of Dark Neurosis
	core.NewItemEffect(220745, echoesOfMadnessEffect)

	// Paranoia Mantle
	core.NewItemEffect(220747, echoesOfDreadEffect)

	// Shoulderpads of Obsession
	core.NewItemEffect(220748, echoesOfInsanityEffect)

	///////////////////////////////////////////////////////////////////////////
	//                                 Mail
	///////////////////////////////////////////////////////////////////////////

	// Screaming Chain Pauldrons
	core.NewItemEffect(220741, echoesOfMadnessEffect)

	// Shrieking Spaulders
	core.NewItemEffect(220742, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Echoes of the Damned Proc", core.ActionID{SpellID: 446618}, stats.Stats{stats.AttackPower: 60}, time.Second*10)

		handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			procAura.Activate(sim)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 446620},
			Name:       "Echoes of the Damned",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			ProcChance: 0.3,
			ICD:        time.Second * 40,
			Handler:    handler,
		})
	})

	// Cacophonous Chain Shoulderguards
	core.NewItemEffect(220743, echoesOfDreadEffect)

	// Wailing Chain Mantle
	core.NewItemEffect(220744, echoesOfInsanityEffect)

	///////////////////////////////////////////////////////////////////////////
	//                                 Plate
	///////////////////////////////////////////////////////////////////////////

	// Shoulderplates of Dread
	core.NewItemEffect(220738, echoesOfDreadEffect)

	// Baleful Pauldrons
	core.NewItemEffect(220739, echoesOfInsanityEffect)

	// Fearmonger's Shoulderguards
	core.NewItemEffect(220740, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Echoes of Fear Proc", core.ActionID{SpellID: 446597}, stats.Stats{stats.SpellDamage: 50}, time.Second*10)

		handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if core.ProcMaskRanged.Matches(spell.ProcMask) {
				return
			}
			procAura.Activate(sim)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{SpellID: 446592},
			Name:       "Echoes of Fear",
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskDirect,
			ProcChance: 0.3,
			ICD:        time.Second * 40,
			Handler:    handler,
		})
	})

	core.AddEffectsToTest = true
}
