package crafted

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	ShoulderplatesOfDread          = 220738
	BalefulPauldrons               = 220739
	FearmongersShoulderguards      = 220740
	ScreamingChainPauldrons        = 220741
	ShriekingSpaulders             = 220742
	CacophonousChainShoulderguards = 220743
	WailingChainMantle             = 220744
	MembraneOfDarkNeurosis         = 220745
	ParanoiaMantle                 = 220747
	ShoulderpadsOfObsession        = 220748
	MantleOfInsanity               = 220749
	FracturedMindPauldrons         = 220750
	ShoulderpadsOfTheDeranged      = 220751
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Cloth
	///////////////////////////////////////////////////////////////////////////

	// Mantle of Insanity
	core.NewItemEffect(MantleOfInsanity, echoesOfInsanityEffect)

	// Fractured Mind Pauldrons
	core.NewItemEffect(FracturedMindPauldrons, echoesOfMadnessEffect)

	// Shoulderpads of the Deranged
	core.NewItemEffect(ShoulderpadsOfTheDeranged, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.Level != 50 {
			return
		}

		procAura := character.NewTemporaryStatsAura("Echoes of the Depraved Proc", core.ActionID{SpellID: 446572}, stats.Stats{stats.SpellDamage: 30, stats.Dodge: 2}, time.Second*10)

		handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			procAura.Activate(sim)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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
	core.NewItemEffect(MembraneOfDarkNeurosis, echoesOfFearEffect)

	// Paranoia Mantle
	core.NewItemEffect(ParanoiaMantle, echoesOfDreadEffect)

	// Shoulderpads of Obsession
	core.NewItemEffect(ShoulderpadsOfObsession, echoesOfInsanityEffect)

	///////////////////////////////////////////////////////////////////////////
	//                                 Mail
	///////////////////////////////////////////////////////////////////////////

	// Screaming Chain Pauldrons
	core.NewItemEffect(ScreamingChainPauldrons, echoesOfMadnessEffect)

	// Shrieking Spaulders
	core.NewItemEffect(ShriekingSpaulders, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.Level != 50 {
			return
		}

		procAura := character.NewTemporaryStatsAura("Echoes of the Damned Proc", core.ActionID{SpellID: 446618}, stats.Stats{stats.AttackPower: 60}, time.Second*10)

		handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			procAura.Activate(sim)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Echoes of the Damned",
			Callback:          core.CallbackOnCastComplete,
			ProcMask:          core.ProcMaskMeleeOrRanged,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			ProcChance:        0.3,
			ICD:               time.Second * 40,
			Handler:           handler,
		})
	})

	// Cacophonous Chain Shoulderguards
	core.NewItemEffect(CacophonousChainShoulderguards, echoesOfDreadEffect)

	// Wailing Chain Mantle
	core.NewItemEffect(WailingChainMantle, echoesOfInsanityEffect)

	///////////////////////////////////////////////////////////////////////////
	//                                 Plate
	///////////////////////////////////////////////////////////////////////////

	// Shoulderplates of Dread
	core.NewItemEffect(ShoulderplatesOfDread, echoesOfDreadEffect)

	// Baleful Pauldrons
	core.NewItemEffect(BalefulPauldrons, echoesOfInsanityEffect)

	// Fearmonger's Shoulderguards
	core.NewItemEffect(FearmongersShoulderguards, echoesOfFearEffect)

	core.AddEffectsToTest = true
}

func echoesOfFearEffect(agent core.Agent) {
	character := agent.GetCharacter()

	if character.Level != 50 {
		return
	}

	procAura := character.NewTemporaryStatsAura("Echoes of Fear Proc", core.ActionID{SpellID: 446597}, stats.Stats{stats.SpellDamage: 50}, time.Second*10)

	handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
		if core.ProcMaskRanged.Matches(spell.ProcMask) {
			return
		}
		procAura.Activate(sim)
	}

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       "Echoes of Fear",
		Callback:   core.CallbackOnCastComplete,
		ProcMask:   core.ProcMaskDirect,
		ProcChance: 0.3,
		ICD:        time.Second * 40,
		Handler:    handler,
	})
}

func echoesOfInsanityEffect(agent core.Agent) {
	character := agent.GetCharacter()

	if character.Level != 50 {
		return
	}

	procAura := character.NewTemporaryStatsAura("Echoes of Insanity Proc", core.ActionID{SpellID: 446541}, stats.Stats{stats.HealingPower: 50}, time.Second*10)

	handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
		procAura.Activate(sim)
	}

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       "Echoes of Insanity",
		Callback:   core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		ProcChance: 0.3,
		ICD:        time.Second * 40,
		Handler:    handler,
	})
}

func echoesOfMadnessEffect(agent core.Agent) {
	character := agent.GetCharacter()

	if character.Level != 50 {
		return
	}

	procAura := character.GetOrRegisterAura(core.Aura{
		Label:    "Echoes of Madness Proc",
		ActionID: core.ActionID{SpellID: 446528},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			character.MultiplyCastSpeed(sim, 1.1)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			character.MultiplyCastSpeed(sim, 1/1.1)
		},
	})

	handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
		procAura.Activate(sim)
	}

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       "Echoes of Madness",
		Callback:   core.CallbackOnCastComplete,
		ProcMask:   core.ProcMaskSpellDamage,
		ProcChance: 0.3,
		ICD:        time.Second * 40,
		Handler:    handler,
	})
}

func echoesOfDreadEffect(agent core.Agent) {
	character := agent.GetCharacter()

	if character.Level != 50 {
		return
	}

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
		Name:              "Echoes of Dread",
		Callback:          core.CallbackOnSpellHitDealt,
		ProcMask:          core.ProcMaskMelee,
		SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
		ProcChance:        0.3,
		ICD:               time.Second * 40,
		Handler:           handler,
	})
}
