package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (mage *Mage) ApplyRunes() {
	// Helm
	mage.registerDeepFreezeSpell()

	// Shoulders
	mage.applyShoulderRuneEffect()

	// Cloak
	mage.registerArcaneBarrageSpell()
	mage.applyOverheat()
	mage.registerFrozenOrbCD()

	// Chest
	mage.applyBurnout()
	mage.applyEnlightenment()
	mage.applyFingersOfFrost()

	// Bracers
	mage.registerBalefireBoltSpell()

	// Hands
	mage.registerArcaneBlastSpell()
	mage.registerIceLanceSpell()
	mage.registerLivingBombSpell()

	// Waist
	mage.registerFrostfireBoltSpell()
	mage.applyHotStreak()
	mage.applyMissileBarrage()
	mage.registerSpellfrostBoltSpell()

	// Legs
	mage.registerArcaneSurgeSpell()
	mage.registerIcyVeinsSpell()
	mage.registerLivingFlameSpell()
	mage.registerMassRegenerationSpell()

	// Feet
	mage.applyBrainFreeze()
	mage.applySpellPower()
}

func (mage *Mage) applyShoulderRuneEffect() {
	if mage.Equipment.Shoulders().Rune == int32(proto.MageRune_MageRuneNone) {
		return
	}

	switch mage.Equipment.Shoulders().Rune {
	// Damage
	case int32(proto.MageRune_RuneShouldersElementalist):
		mage.applyT1Damage4PBonus()
	case int32(proto.MageRune_RuneShouldersMagicalArmorer):
		mage.applyT1Damage6PBonus()
	case int32(proto.MageRune_RuneShouldersKindler):
		mage.applyT2Damage2PBonus()
	case int32(proto.MageRune_RuneShouldersFieryConvergence):
		mage.applyT2Damage4PBonus()
	case int32(proto.MageRune_RuneShouldersPerpetualBlaze):
		mage.applyT2Damage6PBonus()
	case int32(proto.MageRune_RuneShouldersWintersGrasp):
		mage.applyZGFrost3PBonus()
	case int32(proto.MageRune_RuneShouldersCryomancer):
		mage.applyZGFrost5PBonus()
	case int32(proto.MageRune_RuneShouldersPyromaniac):
		mage.applyTAQFire2PBonus()
	case int32(proto.MageRune_RuneShouldersIgniter):
		mage.applyTAQFire4PBonus()
	case int32(proto.MageRune_RuneShouldersTorcher):
		mage.applyRAQFire3PBonus()

	// Healer
	case int32(proto.MageRune_RuneShouldersPrecognitive):
		mage.applyT2Healer2PBonus()
	case int32(proto.MageRune_RuneShouldersArcanist):
		mage.applyT2Healer4PBonus()
	case int32(proto.MageRune_RuneShouldersSpellbinder):
		mage.applyTAQArcane2PBonus()
	}
}

func (mage *Mage) applyOverheat() {
	if !mage.HasRune(proto.MageRune_RuneCloakOverheat) {
		return
	}

	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_MageFireBlast) {
			spell.BonusCritRating += 100 * core.SpellCritRatingPerCritChance
			spell.CD.Duration = time.Second * 15
			spell.Flags |= core.SpellFlagCastTimeNoGCD | core.SpellFlagCastWhileCasting
			spell.DefaultCast.GCD = 0
		}
	})
}

func (mage *Mage) applyBurnout() {
	if !mage.HasRune(proto.MageRune_RuneChestBurnout) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.MageRune_RuneChestBurnout)}
	metric := mage.NewManaMetrics(actionID)

	mage.AddStat(stats.SpellCrit, 15*core.SpellCritRatingPerCritChance)

	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:           "Burnout",
		Outcome:        core.OutcomeCrit,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: ClassSpellMask_MageAll,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			spell.Unit.SpendMana(sim, min(spell.Unit.BaseMana*0.01, mage.CurrentMana()), metric)
		},
	})
}

func (mage *Mage) applyEnlightenment() {
	if !mage.HasRune(proto.MageRune_RuneChestEnlightenment) {
		return
	}

	damageAuraThreshold := .70
	manaAuraThreshold := .30

	// https://www.wowhead.com/classic/spell=412326/enlightenment
	damageAura := mage.RegisterAura(core.Aura{
		Label:    "Enlightenment (Damage)",
		ActionID: core.ActionID{SpellID: 412326},
		Duration: core.NeverExpires,
	}).AttachMultiplicativePseudoStatBuff(&mage.PseudoStats.DamageDealtMultiplier, 1.1)

	// https://www.wowhead.com/classic/spell=412325/enlightenment
	manaAura := mage.RegisterAura(core.Aura{
		Label:    "Enlightenment (Mana)",
		ActionID: core.ActionID{SpellID: 412325},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpiritRegenRateCasting += 0.10
			mage.UpdateManaRegenRates()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpiritRegenRateCasting -= .10
			mage.UpdateManaRegenRates()
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Enlightenment",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
			damageAura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			percentMana := aura.Unit.CurrentManaPercent()

			if percentMana > damageAuraThreshold && !damageAura.IsActive() {
				damageAura.Activate(sim)
			} else if percentMana <= damageAuraThreshold {
				damageAura.Deactivate(sim)
			}

			if percentMana < manaAuraThreshold && !manaAura.IsActive() {
				manaAura.Activate(sim)
			} else if percentMana >= manaAuraThreshold {
				manaAura.Deactivate(sim)
			}
		},
	})
}

func (mage *Mage) applyFingersOfFrost() {
	if !mage.HasRune(proto.MageRune_RuneChestFingersOfFrost) {
		return
	}

	mage.FingersOfFrostProcChance += 0.25

	mage.FingersOfFrostAura = mage.RegisterAura(core.Aura{
		Label:     "Fingers of Frost Proc",
		ActionID:  core.ActionID{SpellID: int32(proto.MageRune_RuneChestFingersOfFrost)},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, aura := range mage.FrozenAuras {
				aura.Activate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, aura := range mage.FrozenAuras {
				aura.Deactivate(sim)
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}

			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}

			if aura.GetStacks() == 1 {
				// Fingers of Frost can be batched with a casted spell into an instant
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + core.SpellBatchWindow,
					OnAction: func(sim *core.Simulation) {
						if aura.IsActive() && aura.GetStacks() == 1 {
							aura.RemoveStack(sim)
						}
					},
				})
			} else {
				aura.RemoveStack(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:           "Fingers of Frost Trigger - Direct",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: ClassSpellMask_MageAll ^ ClassSpellMask_MageFrozenOrbTick, // Blizzard seems to have intentionally made Frozen Orb's chill not proc Fingers of Frost but a Frostbite proc from the orb still can
		SpellFlags:     SpellFlagChillSpell,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(mage.FingersOfFrostProcChance, "Fingers of Frost") {
				mage.FingersOfFrostAura.Activate(sim)
				mage.FingersOfFrostAura.SetStacks(sim, mage.FingersOfFrostAura.MaxStacks)
			}
		},
	})

	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:           "Fingers of Frost Trigger - Periodic",
		Callback:       core.CallbackOnPeriodicDamageDealt,
		ClassSpellMask: ClassSpellMask_MageBlizzard, // Only procs from Blizzard and only with Improved Blizzard for the chill effect
		SpellFlags:     SpellFlagChillSpell,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(mage.FingersOfFrostProcChance, "Fingers of Frost") {
				mage.FingersOfFrostAura.Activate(sim)
				mage.FingersOfFrostAura.SetStacks(sim, mage.FingersOfFrostAura.MaxStacks)
			}
		},
	})
}

func (mage *Mage) applyHotStreak() {
	if !mage.HasRune(proto.MageRune_RuneHelmHotStreak) {
		return
	}

	actionID := core.ActionID{SpellID: 48108}

	triggerSpellClassMasks := ClassSpellMask_MageFireball |
		ClassSpellMask_MageFrostfireBolt |
		ClassSpellMask_MageBalefireBolt |
		ClassSpellMask_MageFireBlast |
		ClassSpellMask_MageScorch |
		ClassSpellMask_MageLivingBombExplosion

	castTimeMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  ClassSpellMask_MagePyroblast,
		FloatValue: -1,
	})
	costMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: ClassSpellMask_MagePyroblast,
		IntValue:  -100,
	})

	mage.HotStreakAura = mage.RegisterAura(core.Aura{
		Label:    "Hot Streak",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
			costMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
			costMod.Deactivate()
		},
	})

	heatingUpAura := mage.RegisterAura(core.Aura{
		Label:    "Heating Up",
		ActionID: actionID.WithTag(1),
		Duration: time.Hour,
	})

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Hot Streak Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(triggerSpellClassMasks) {
				return
			}

			if !result.DidCrit() {
				if heatingUpAura.IsActive() {
					heatingUpAura.Deactivate(sim)
				}

				return
			}

			if heatingUpAura.IsActive() {
				heatingUpAura.Deactivate(sim)
				mage.HotStreakAura.Activate(sim)
			} else if mage.HotStreakAura.IsActive() {
				// When batching a Scorch crit into an instant Pyro, the Pyro consumes Hot Streak before the Scorch hits, so the Scorch re-applies Heating Up
				// We can replicate this by adding a batch delay then checking the state of the auras again.
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + core.SpellBatchWindow,
					OnAction: func(sim *core.Simulation) {
						if heatingUpAura.IsActive() {
							heatingUpAura.Deactivate(sim)
							mage.HotStreakAura.Activate(sim)
						} else if !mage.HotStreakAura.IsActive() {
							heatingUpAura.Activate(sim)
						}
					},
				})
			} else {
				heatingUpAura.Activate(sim)
			}
		},
	}))
}

func (mage *Mage) applyMissileBarrage() {
	if !mage.HasRune(proto.MageRune_RuneBeltMissileBarrage) {
		return
	}

	fireballFrostboltMissileBarrageChance := 0.20
	mage.ArcaneBlastMissileBarrageChance += 0.40
	buffDuration := time.Second * 15

	arcaneMissilesSpells := []*core.Spell{}
	affectedSpellClassMasks := ClassSpellMask_MageArcaneBlast | ClassSpellMask_MageFireball | ClassSpellMask_MageFrostbolt

	mage.MissileBarrageAura = mage.RegisterAura(core.Aura{
		Label:    "Missile Barrage",
		ActionID: core.ActionID{SpellID: 400589},
		Duration: buffDuration,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			arcaneMissilesSpells = core.FilterSlice(mage.ArcaneMissiles, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(arcaneMissilesSpells, func(spell *core.Spell) {
				spell.Cost.Multiplier -= 10000
				for _, target := range sim.Encounter.TargetUnits {
					spell.Dot(target).TickLength /= 2
				}
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(arcaneMissilesSpells, func(spell *core.Spell) {
				spell.Cost.Multiplier += 10000
				for _, target := range sim.Encounter.TargetUnits {
					spell.Dot(target).TickLength *= 2
				}
			})
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Missile Barrage Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(affectedSpellClassMasks) {
				return
			}

			procChance := core.TernaryFloat64(spell.Matches(ClassSpellMask_MageArcaneBlast), mage.ArcaneBlastMissileBarrageChance, fireballFrostboltMissileBarrageChance)
			if sim.Proc(procChance, "Missile Barrage") {
				mage.MissileBarrageAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) applyBrainFreeze() {
	if !mage.HasRune(proto.MageRune_RuneFeetBrainFreeze) {
		return
	}

	procChance := .20
	buffDuration := time.Second * 15

	affectedSpells := []*core.Spell{}
	triggerSpellClassMasks := ClassSpellMask_MageFireball | ClassSpellMask_MageFrostfireBolt | ClassSpellMask_MageSpellfrostBolt | ClassSpellMask_MageBalefireBolt

	procAura := mage.RegisterAura(core.Aura{
		Label:    "Brain Freeze",
		ActionID: core.ActionID{SpellID: 400730},
		Duration: buffDuration,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(
				core.Flatten([][]*core.Spell{
					mage.Fireball,
					{mage.FrostfireBolt},
					{mage.SpellfrostBolt},
					{mage.BalefireBolt},
				}),
				func(spell *core.Spell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier -= 1
				spell.Cost.Multiplier -= 100
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier += 1
				spell.Cost.Multiplier += 100
			})
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}

			if spell.Matches(triggerSpellClassMasks) && spell.CurCast.CastTime == 0 {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Brain Freeze Trigger",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Only Blizzard ticks proc
			if spell.Matches(ClassSpellMask_MageBlizzard) && spell.Flags.Matches(SpellFlagChillSpell) && sim.Proc(procChance, "Brain Freeze") {
				procAura.Activate(sim)
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Flags.Matches(SpellFlagChillSpell) && result.Landed() && sim.Proc(procChance, "Brain Freeze") {
				procAura.Activate(sim)
			}
		},
	}))
}

func (mage *Mage) applySpellPower() {
	if !mage.HasRune(proto.MageRune_RuneFeetSpellPower) {
		return
	}

	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_MageAll) {
			spell.CritDamageBonus += 0.5
		}
	})
}
