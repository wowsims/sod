package crafted

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Cloth
	///////////////////////////////////////////////////////////////////////////

	// Mantle of Insanity
	core.NewItemEffect(220749, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}

		procAura := character.NewTemporaryStatsAura("Echoes of Insanity Proc", core.ActionID{SpellID: 446541}, stats.Stats{stats.HealingPower: 50}, time.Second*10)

		character.RegisterAura(core.Aura{
			Label:    "Echoes of Insanity",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskSpellHealing.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of the Depraved") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	// Fractured Mind Pauldrons
	core.NewItemEffect(220750, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}

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

		character.RegisterAura(core.Aura{
			Label:    "Echoes of Madness",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskSpellDamage.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of Madness") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	// Shoulderpads of the Deranged
	core.NewItemEffect(220751, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}

		procAura := character.NewTemporaryStatsAura("Echoes of the Depraved Proc", core.ActionID{SpellID: 446572}, stats.Stats{stats.SpellDamage: 30, stats.Dodge: 2}, time.Second*10)

		character.RegisterAura(core.Aura{
			Label:    "Echoes of the Depraved",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskSpellDamage.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of the Depraved") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Leather
	///////////////////////////////////////////////////////////////////////////

	// Membrane of Dark Neurosis
	core.NewItemEffect(220745, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}

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

		character.RegisterAura(core.Aura{
			Label:    "Echoes of Madness",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskSpellDamage.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of Madness") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	// Paranoia Mantle
	core.NewItemEffect(220747, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
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

		character.RegisterAura(core.Aura{
			Label:    "Echoes of Dread",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskMelee.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of Dread") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	// Shoulderpads of Obsession
	core.NewItemEffect(220748, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}

		procAura := character.NewTemporaryStatsAura("Echoes of Insanity Proc", core.ActionID{SpellID: 446541}, stats.Stats{stats.HealingPower: 50}, time.Second*10)

		character.RegisterAura(core.Aura{
			Label:    "Echoes of Insanity",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskSpellHealing.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of the Depraved") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Mail
	///////////////////////////////////////////////////////////////////////////

	// Shrieking Spaulders
	core.NewItemEffect(220742, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}

		procAura := character.NewTemporaryStatsAura("Echoes of the Damned Proc", core.ActionID{SpellID: 446618}, stats.Stats{stats.AttackPower: 60}, time.Second*10)

		character.RegisterAura(core.Aura{
			Label:    "Echoes of the Damned",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskMeleeOrRanged.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of the Damned") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	// Cacophonous Chain Shoulderguards
	core.NewItemEffect(220743, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
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

		character.RegisterAura(core.Aura{
			Label:    "Echoes of Dread",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskMelee.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of Dread") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	// Wailing Chain Mantle
	core.NewItemEffect(220744, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}

		procAura := character.NewTemporaryStatsAura("Echoes of Insanity Proc", core.ActionID{SpellID: 446541}, stats.Stats{stats.HealingPower: 50}, time.Second*10)

		character.RegisterAura(core.Aura{
			Label:    "Echoes of Insanity",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskSpellHealing.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of the Depraved") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Plate
	///////////////////////////////////////////////////////////////////////////

	// Shoulderplates of Dread
	core.NewItemEffect(220738, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
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

		character.RegisterAura(core.Aura{
			Label:    "Echoes of Dread",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskMelee.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of Dread") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	// Baleful Pauldrons
	core.NewItemEffect(220739, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}

		procAura := character.NewTemporaryStatsAura("Echoes of Insanity Proc", core.ActionID{SpellID: 446541}, stats.Stats{stats.HealingPower: 50}, time.Second*10)

		character.RegisterAura(core.Aura{
			Label:    "Echoes of Insanity",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) || !core.ProcMaskSpellHealing.Matches(spell.ProcMask) {
					return
				}

				if sim.RandomFloat("Echoes of the Depraved") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	// Fearmonger's Shoulderguards
	core.NewItemEffect(220740, func(agent core.Agent) {
		character := agent.GetCharacter()

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}

		procAura := character.NewTemporaryStatsAura("Echoes of Fear Proc", core.ActionID{SpellID: 446597}, stats.Stats{stats.SpellDamage: 50}, time.Second*10)

		character.RegisterAura(core.Aura{
			Label:    "Echoes of Fear",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !icd.IsReady(sim) {
					return
				}

				if !(core.ProcMaskSpellDamage.Matches(spell.ProcMask) || core.ProcMaskMelee.Matches(spell.ProcMask)) {
					return
				}

				if sim.RandomFloat("Echoes of Fear") < 0.3 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	core.AddEffectsToTest = true
}
