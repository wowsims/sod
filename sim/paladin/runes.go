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
	paladin.registerShockAndAwe()
	paladin.registerRV()
	paladin.registerFanaticism()

	// "RuneHeadWrath" is handled in Exorcism, Holy Shock, Consecration and Holy Wrath
	paladin.registerMalleableProtection()
	paladin.registerHammerOfTheRighteous()
	// "RuneWristImprovedHammerOfWrath" is handled Hammer of Wrath
	paladin.applyPurifyingPower()
	paladin.registerAegis()
	paladin.registerAvengersShield()

	paladin.applyShoulderRuneEffect()
}

func (paladin *Paladin) registerFanaticism() {
	if paladin.hasRune(proto.PaladinRune_RuneHeadFanaticism) {
		paladin.PseudoStats.SchoolBonusCritChance[stats.SchoolIndexHoly] += 18
	}
}

func (paladin *Paladin) applyShoulderRuneEffect() {

	if paladin.Equipment.Shoulders().Rune == int32(proto.PaladinRune_PaladinRuneNone) {
		return
	}

	switch paladin.Equipment.Shoulders().Rune {
	// Prot
	case int32(proto.PaladinRune_RuneShouldersPristineBlocker):
		paladin.applyPaladinT1Prot2P()
	case int32(proto.PaladinRune_RuneShouldersLightwarden):
		paladin.applyPaladinT1Prot4P()
	case int32(proto.PaladinRune_RuneShouldersRadiantDefender):
		paladin.applyPaladinT1Prot6P()
	case int32(proto.PaladinRune_RuneShouldersShieldbearer):
		paladin.applyPaladinT2Prot2P()
	case int32(proto.PaladinRune_RuneShouldersBastion):
		paladin.applyPaladinT2Prot4P()
	case int32(proto.PaladinRune_RuneShouldersReckoner):
		paladin.applyPaladinT2Prot6P()
	case int32(proto.PaladinRune_RuneShouldersIronclad):
		paladin.applyPaladinTAQProt2P()
	case int32(proto.PaladinRune_RuneShouldersGuardian):
		paladin.applyPaladinTAQProt4P()

	// Holy
	case int32(proto.PaladinRune_RuneShouldersPeacekeeper):
		paladin.applyPaladinT1Holy2P()
	case int32(proto.PaladinRune_RuneShouldersRefinedPaladin):
		paladin.applyPaladinT1Holy4P()
	case int32(proto.PaladinRune_RuneShouldersExemplar):
		paladin.applyPaladinT1Holy6P()
	case int32(proto.PaladinRune_RuneShouldersInquisitor):
		paladin.applyPaladinT2Holy2P()
	case int32(proto.PaladinRune_RuneShouldersSovereign):
		paladin.applyPaladinT2Holy4P()
	case int32(proto.PaladinRune_RuneShouldersDominus):
		paladin.applyPaladinT2Holy6P()
	case int32(proto.PaladinRune_RuneShouldersVindicator):
		paladin.applyPaladinTAQHoly2P()
	case int32(proto.PaladinRune_RuneShouldersAltruist):
		paladin.applyPaladinTAQHoly4P()

	// Ret
	case int32(proto.PaladinRune_RuneShouldersArbiter):
		paladin.applyPaladinT1Ret2P()
	// T2 4P for ret is missing because it is the same as 4P for Holy
	case int32(proto.PaladinRune_RuneShouldersSealbearer):
		paladin.applyPaladinT1Ret6P()
	case int32(proto.PaladinRune_RuneShouldersJusticar):
		paladin.applyPaladinT2Ret2P()
	case int32(proto.PaladinRune_RuneShouldersJudicator):
		paladin.applyPaladinT2Ret4P()
	case int32(proto.PaladinRune_RuneShouldersAscendant):
		paladin.applyPaladinT2Ret6P()
	case int32(proto.PaladinRune_RuneShouldersRetributor):
		paladin.applyPaladinTAQRet2P()
	case int32(proto.PaladinRune_RuneShouldersExcommunicator):
		paladin.applyPaladinTAQRet4P()
	case int32(proto.PaladinRune_RuneShouldersTemplar):
		paladin.applyPaladinRAQ3P()

	// ZG (Shockadin)
	case int32(proto.PaladinRune_RuneShouldersLightbringer):
		paladin.applyPaladinZG3P()
	case int32(proto.PaladinRune_RuneShouldersExile):
		paladin.applyPaladinZG5P()
	}
}

func (paladin *Paladin) registerTheArtOfWar() {
	if !paladin.hasRune(proto.PaladinRune_RuneFeetTheArtOfWar) {
		return
	}

	actionID := core.ActionID{SpellID: 426157}
	paladin.artOfWarDelayAura = paladin.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "The Art of War Delay",
		Duration: time.Millisecond * 250,
	})

	spellQueueWindow := time.Millisecond * 400
	cdReduction := time.Second * 2
	aowSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagPassiveSpell,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			timeToReady := paladin.exorcismCooldown.TimeToReady(sim)
			actualReduction := min(timeToReady, cdReduction)
			newTimeToReady := timeToReady - actualReduction

			if sim.Log != nil {
				paladin.Log(sim, "The Art of War reduced Exorcism cooldown by %s (%s -> %s)", actualReduction, timeToReady, newTimeToReady)
			}

			if newTimeToReady <= 0 && actualReduction > spellQueueWindow && paladin.NextGCDAt()-sim.CurrentTime <= spellQueueWindow {
				paladin.artOfWarDelayAura.Activate(sim)
			}

			paladin.exorcismCooldown.Set(sim.CurrentTime + newTimeToReady)
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:       "The Art of War Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee | core.ProcMaskMeleeDamageProc,
		Outcome:    core.OutcomeCrit,
		ProcChance: 1,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if paladin.exorcismCooldown.TimeToReady(sim) <= 0 {
				return
			}

			aowSpell.Cast(sim, result.Target)
		},
	})
}

func (paladin *Paladin) registerSheathOfLight() {
	if !paladin.hasRune(proto.PaladinRune_RuneWaistSheathOfLight) {
		return
	}

	var prevSPBonus = 0.0

	sheathAura := paladin.RegisterAura(core.Aura{
		Label:    "Sheath of Light",
		Duration: time.Second * 60,
		ActionID: core.ActionID{SpellID: 426159},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			newSPBonus := paladin.GetStat(stats.AttackPower) * 0.3
			paladin.AddStatDynamic(sim, stats.SpellDamage, +newSPBonus)

			if (newSPBonus != prevSPBonus) && (sim.Log != nil) {
				paladin.Log(sim, "Sheath of Light new bonus is %d old was %d", int32(newSPBonus), int32(prevSPBonus))
			}

			prevSPBonus = newSPBonus

		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			newSPBonus := paladin.GetStat(stats.AttackPower) * 0.3
			paladin.AddStatDynamic(sim, stats.SpellDamage, newSPBonus-prevSPBonus)

			if (newSPBonus != prevSPBonus) && (sim.Log != nil) {
				paladin.Log(sim, "Sheath of Light new bonus is %d old was %d", int32(newSPBonus), int32(prevSPBonus))
			}

			prevSPBonus = newSPBonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.SpellDamage, -prevSPBonus)
			prevSPBonus = 0.0
		},
	})
	paladin.RegisterAura(core.Aura{
		Label:    "Sheath of Light (rune)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			sheathAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) registerShockAndAwe() {
	if !paladin.hasRune(proto.PaladinRune_RuneCloakShockAndAwe) {
		return
	}

	hasWrath := paladin.hasRune(proto.PaladinRune_RuneHeadWrath)

	dep := paladin.NewDynamicStatDependency(stats.Intellect, stats.SpellDamage, 2.0)

	shockAndAweAura := paladin.RegisterAura(core.Aura{
		Label:    "Shock and Awe",
		ActionID: core.ActionID{SpellID: 462832},
		Duration: time.Second * 60,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.EnableDynamicStatDep(sim, dep)
			paladin.PseudoStats.ThreatMultiplier *= 0.70
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.DisableDynamicStatDep(sim, dep)
			paladin.PseudoStats.ThreatMultiplier /= 0.70
		},
	})

	var sunlightAura *core.Aura

	if paladin.Talents.HolyShock && paladin.hasRune(proto.PaladinRune_RuneWaistInfusionOfLight) {
		sunlightSpell := paladin.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: 1239548},
			SpellSchool:    core.SpellSchoolHoly,
			DefenseType:    core.DefenseTypeMagic,
			ProcMask:       core.ProcMaskEmpty,
			ClassSpellMask: ClassSpellMask_PaladinSunlight,
			Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			BonusCoefficient: 1,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				bonusCrit := 0.0

				if hasWrath {
					bonusCrit += paladin.GetStat(stats.MeleeCrit)
				}

				spell.BonusCritRating += bonusCrit
				spell.CalcAndDealDamage(sim, target, 1, spell.OutcomeMagicCrit)
				spell.BonusCritRating -= bonusCrit
			},
		})

		sunlightAura = paladin.RegisterAura(core.Aura{
			Label:    "Sunlight",
			ActionID: core.ActionID{SpellID: 1239543},
			Duration: time.Second * 10,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(core.SpellFlagSuppressEquipProcs | core.SpellFlagSuppressWeaponProcs) {
					return
				}

				if !result.Landed() {
					return
				}

				// Does not work when RF is active
				if paladin.righteousFuryAura != nil && paladin.righteousFuryAura.IsActive() {
					return
				}

				if spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto) || spell.Matches(ClassSpellMask_PaladinCrusaderStrike|ClassSpellMask_PaladinDivineStorm) {
					sunlightSpell.Cast(sim, result.Target)
				}
			},
		})
	}

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Shock and Awe Trigger",
		ActionID:       core.ActionID{SpellID: 462834},
		Callback:       core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ClassSpellMask: ClassSpellMask_PaladinHolyShock,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shockAndAweAura.Activate(sim)

			if sunlightAura != nil {
				sunlightAura.Activate(sim)
			}
		},
	})
}

func (paladin *Paladin) registerGuardedByTheLight() {
	if !paladin.hasRune(proto.PaladinRune_RuneFeetGuardedByTheLight) {
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

func (paladin *Paladin) applyPurifyingPower() {
	if !paladin.hasRune(proto.PaladinRune_RuneWristPurifyingPower) {
		return
	}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label: "Purifying Power",
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_PaladinExorcism | ClassSpellMask_PaladinHolyWrath,
		Kind:      core.SpellMod_Cooldown_Multi_Flat,
		IntValue:  -50,
	}))

}

func (paladin *Paladin) registerAegis() {

	if !paladin.hasRune(proto.PaladinRune_RuneChestAegis) {
		return
	}

	// The SBV bonus is additive with Shield Specialization.
	paladin.PseudoStats.BlockValueMultiplier += 0.3

	if paladin.Talents.Redoubt > 0 {
		// Redoubt now has a 10% chance to trigger on any melee or ranged attack against
		// you, and always triggers on your melee critical strikes.
		paladin.RegisterAura(core.Aura{
			Label:    "Redoubt Aegis Trigger",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) && result.Landed() {
					if sim.Proc(0.1, "Aegis Attack") {
						paladin.redoubtAura.Activate(sim)
						paladin.redoubtAura.SetStacks(sim, 5)
					}
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMelee) && result.DidCrit() {
					paladin.redoubtAura.Activate(sim)
					paladin.redoubtAura.SetStacks(sim, 5)
				}
			},
		})
	}

	if paladin.Talents.Reckoning > 0 {
		// Reckoning now also procs on any melee or ranged attack against you with (2% * talent points) chance
		procID := core.ActionID{SpellID: 20178} // reckoning proc id
		procChance := 0.02 * float64(paladin.Talents.Reckoning)

		core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
			Name:       "Reckoning Aegis Trigger",
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Outcome:    core.OutcomeLanded ^ core.OutcomeCrit,
			ProcChance: procChance,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				paladin.AutoAttacks.ExtraMHAttackProc(sim, 1, procID, spell)
			},
		})
	}
}

func (paladin *Paladin) registerMalleableProtection() {
	if !paladin.hasRune(proto.PaladinRune_RuneWaistMalleableProtection) {
		return
	}
	// Activating Holy Shield now grants 4 AP for each point of defense above paladin.Level * 5
	defendersResolveAPAura := core.DefendersResolveAttackPower(paladin.GetCharacter())
	handler := func(spell *core.Spell) {
		if !spell.Matches(ClassSpellMask_PaladinHolyShield) {
			return
		}
		oldEffects := spell.ApplyEffects
		spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			oldEffects(sim, target, spell)
			if stacks := int32(paladin.GetStat(stats.Defense)); stacks > 0 {
				defendersResolveAPAura.Activate(sim)
				if defendersResolveAPAura.GetStacks() != stacks {
					defendersResolveAPAura.SetStacks(sim, stacks)
				}
			}
		}
	}
	paladin.OnSpellRegistered(handler)

	// A prot paladin will only ever cast Divine Protection in conjunction with Malleable Protection,
	// so we register only the modified form of the spell when the rune is engraved.
	// Although there are two spell ranks, intentional downranking is never done in practice,
	// so we only register the highest spell rank available.
	if paladin.Level < 10 {
		return
	}

	isRank1 := paladin.Level < 18
	spellID := core.TernaryInt32(isRank1, 458312, 458371)
	manaCost := core.TernaryFloat64(isRank1, 15, 35)
	duration := core.TernaryDuration(isRank1, 9, 12)

	actionID := core.ActionID{SpellID: spellID}

	dpAura := paladin.RegisterAura(core.Aura{
		Label:    "Divine Protection",
		ActionID: actionID,
		Duration: time.Second * duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier *= 0.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier /= 0.5
		},
	})

	paladin.divineProtection = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_PaladinDivineProtection,
		Flags:          core.SpellFlagAPL | SpellFlag_Forbearance,
		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 5,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dpAura.Activate(sim)
		},
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell:    paladin.divineProtection,
		Priority: core.CooldownPriorityDrums, // Primary defensive cooldown
		Type:     core.CooldownTypeSurvival,
	})
}
