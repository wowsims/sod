package mage

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (mage *Mage) ApplyRunes() {
	// Helm
	mage.registerDeepFreezeSpell()

	// CLoak
	mage.registerArcaneBarrageSpell()
	// mage.applyOverheat()
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

	// Feet
	mage.applyBrainFreeze()
	mage.applySpellPower()
}

func (mage *Mage) applyBurnout() {
	if !mage.HasRune(proto.MageRune_RuneChestBurnout) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.MageRune_RuneChestBurnout)}
	metric := mage.NewManaMetrics(actionID)

	mage.AddStat(stats.SpellCrit, 15*core.SpellCritRatingPerCritChance)

	mage.RegisterAura(core.Aura{
		Label:    "Burnout",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagMage) && !result.DidCrit() {
				return
			}
			aura.Unit.SpendMana(sim, aura.Unit.BaseMana*0.01, metric)
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
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
		},
	})

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
	bonusCrit := 10 * float64(mage.Talents.Shatter) * core.SpellCritRatingPerCritChance

	mage.FingersOfFrostAura = mage.RegisterAura(core.Aura{
		Label:     "Fingers of Frost Proc",
		ActionID:  core.ActionID{SpellID: int32(proto.MageRune_RuneChestFingersOfFrost)},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
			if mage.IceLance != nil {
				mage.IceLance.DamageMultiplier *= 4.0
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
			if mage.IceLance != nil {
				mage.IceLance.DamageMultiplier /= 4.0
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}

			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) || !spell.SpellSchool.Matches(core.SpellSchoolFrost) {
				return
			}

			if aura.GetStacks() == 1 {
				// Fingers of Frost can be batched with 2x FFBs into Deep Freeze
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + core.SpellBatchWindow,
					OnAction: func(sim *core.Simulation) {
						if aura.IsActive() {
							aura.RemoveStack(sim)
						}
					},
				})
			} else {
				aura.RemoveStack(sim)
			}
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Fingers of Frost Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Flags.Matches(SpellFlagChillSpell) && sim.RandomFloat("Fingers of Frost") < mage.FingersOfFrostProcChance {
				mage.FingersOfFrostAura.Activate(sim)
				mage.FingersOfFrostAura.SetStacks(sim, 2)
			}
		},
	})
}

func (mage *Mage) applyHotStreak() {
	if !mage.HasRune(proto.MageRune_RuneHelmHotStreak) {
		return
	}

	actionID := core.ActionID{SpellID: 48108}

	pyroblastSpells := []*core.Spell{}
	triggerSpellCodes := []int32{SpellCode_MageFireball, SpellCode_MageFrostfireBolt, SpellCode_MageBalefireBolt, SpellCode_MageFireBlast, SpellCode_MageScorch, SpellCode_MageLivingBombExplosion}

	mage.HotStreakAura = mage.RegisterAura(core.Aura{
		Label:    "Hot Streak",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			pyroblastSpells = core.FilterSlice(mage.Pyroblast, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(pyroblastSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier -= 1
				spell.Cost.Multiplier -= 100
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(pyroblastSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier += 1
				spell.Cost.Multiplier += 100
			})
		},
	})

	heatingUpAura := mage.RegisterAura(core.Aura{
		Label:    "Heating Up",
		ActionID: actionID.WithTag(1),
		Duration: time.Hour,
	})

	mage.RegisterAura(core.Aura{
		Label:    "Hot Streak Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !slices.Contains(triggerSpellCodes, spell.SpellCode) {
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
				// We can replicate this by adding a 1ms delay then checking the state of the auras again.
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
	})
}

func (mage *Mage) applyMissileBarrage() {
	if !mage.HasRune(proto.MageRune_RuneBeltMissileBarrage) {
		return
	}

	procChance := .20
	mage.ArcaneBlastMissileBarrageChance = .40
	buffDuration := time.Second * 15

	arcaneMissilesSpells := []*core.Spell{}
	affectedSpellCodes := []int32{SpellCode_MageArcaneBarrage, SpellCode_MageArcaneBlast, SpellCode_MageFireball, SpellCode_MageFrostbolt}

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
			if !slices.Contains(affectedSpellCodes, spell.SpellCode) {
				return
			}

			procChance := procChance
			if spell.SpellCode == SpellCode_MageArcaneBlast {
				procChance = mage.ArcaneBlastMissileBarrageChance
			}

			if sim.RandomFloat("Missile Barrage") < procChance {
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
	triggerSpellCodes := []int32{SpellCode_MageFireball, SpellCode_MageFrostfireBolt, SpellCode_MageSpellfrostBolt, SpellCode_MageBalefireBolt}

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

			if slices.Contains(triggerSpellCodes, spell.SpellCode) && spell.CurCast.CastTime == 0 {
				aura.Deactivate(sim)
			}
		},
	})

	units := []*core.Unit{&mage.Unit}
	// Can also proc from Frozen Orb hits
	if mage.HasRune(proto.MageRune_RuneCloakFrozenOrb) {
		units = append(units, core.MapSlice(mage.frozenOrbPets, func(orb *FrozenOrb) *core.Unit { return &orb.Unit })...)
	}

	for _, unit := range units {
		core.MakePermanent(unit.RegisterAura(core.Aura{
			Label: "Brain Freeze Trigger",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(SpellFlagChillSpell) && result.Landed() && sim.Proc(procChance, "Brain Freeze") {
					procAura.Activate(sim)
				}
			},
		}))
	}
}

func (mage *Mage) applySpellPower() {
	if !mage.HasRune(proto.MageRune_RuneFeetSpellPower) {
		return
	}

	units := []*core.Unit{&mage.Unit}
	// Can also proc from Frozen Orb hits
	if mage.HasRune(proto.MageRune_RuneCloakFrozenOrb) {
		units = append(units, core.MapSlice(mage.frozenOrbPets, func(orb *FrozenOrb) *core.Unit { return &orb.Unit })...)
	}

	for _, unit := range units {
		unit.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Flags.Matches(SpellFlagMage) {
				spell.CritDamageBonus += 0.5
			}
		})
	}
}
