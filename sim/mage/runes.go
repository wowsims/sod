package mage

import (
	"fmt"
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (mage *Mage) ApplyRunes() {
	// Chest
	mage.applyBurnout()
	mage.applyEnlightenment()
	mage.applyFingersOfFrost()

	// Hands
	mage.registerArcaneBlastSpell()
	mage.registerIceLanceSpell()
	mage.registerLivingBombSpell()

	// Waist
	mage.registerFrostfireBoltSpell()
	mage.applyHotStreak()
	mage.applyMissileBarrage()
	mage.registerSpellfrostBolt()

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
			aura.Unit.PseudoStats.SpiritRegenRateCasting *= 1.1
			mage.UpdateManaRegenRates()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpiritRegenRateCasting /= 1.1
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

	procChance := 0.15
	bonusCrit := 10 * float64(mage.Talents.Shatter) * core.SpellCritRatingPerCritChance

	fmt.Println(mage.Talents.Shatter)
	var proccedAt time.Duration

	procAura := mage.RegisterAura(core.Aura{
		Label:     "Fingers of Frost Proc",
		ActionID:  core.ActionID{SpellID: int32(proto.MageRune_RuneChestFingersOfFrost)},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if proccedAt != sim.CurrentTime {
				aura.RemoveStack(sim)
			}
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Fingers of Frost Rune",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Flags.Matches(SpellFlagChillSpell) && sim.RandomFloat("Fingers of Frost") < procChance {
				procAura.Activate(sim)
				procAura.SetStacks(sim, 2)
				proccedAt = sim.CurrentTime
			}
		},
	})
}

func (mage *Mage) applyHotStreak() {
	if !mage.HasRune(proto.MageRune_RuneBeltHotStreak) {
		return
	}

	actionID := core.ActionID{SpellID: 48108}

	pyroblastSpells := []*core.Spell{}
	triggerSpellCodes := []int32{SpellCode_MageFireball, SpellCode_MageFireBlast, SpellCode_MageScorch, SpellCode_MageLivingBomb}

	mage.HotStreakAura = mage.RegisterAura(core.Aura{
		Label:    "Hot Streak",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			pyroblastSpells = core.FilterSlice(mage.Pyroblast, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(pyroblastSpells, func(spell *core.Spell) { spell.CastTimeMultiplier -= 1 })
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(pyroblastSpells, func(spell *core.Spell) { spell.CastTimeMultiplier += 1 })
		},
	})

	procAura := mage.RegisterAura(core.Aura{
		Label:     "Heating Up",
		ActionID:  actionID.WithTag(1),
		MaxStacks: 2,
		Duration:  time.Hour,
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
				procAura.Deactivate(sim)
				return
			}

			if procAura.GetStacks() == 1 {
				procAura.Deactivate(sim)
				mage.HotStreakAura.Activate(sim)
			} else {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			}
		},
	})
}

func (mage *Mage) applyMissileBarrage() {
	if !mage.HasRune(proto.MageRune_RuneBeltMissileBarrage) {
		return
	}

	procChanceArcaneBlast := .40
	procChanceFireballFrostbolt := .20
	buffDuration := time.Second * 15

	arcaneMissilesSpells := []*core.Spell{}
	affectedSpellCodes := []int32{SpellCode_MageArcaneBlast, SpellCode_MageFireball, SpellCode_MageFrostbolt}

	mage.MissileBarrageAura = mage.RegisterAura(core.Aura{
		Label:    "Missile Barrage",
		ActionID: core.ActionID{SpellID: 400589},
		Duration: buffDuration,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			arcaneMissilesSpells = core.FilterSlice(mage.ArcaneMissiles, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(arcaneMissilesSpells, func(spell *core.Spell) {
				spell.CostMultiplier -= 100
				for _, target := range sim.Encounter.TargetUnits {
					spell.Dot(target).TickLength /= 2
				}
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(arcaneMissilesSpells, func(spell *core.Spell) {
				spell.CostMultiplier += 100
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

			procChance := procChanceFireballFrostbolt
			if spell.SpellCode == SpellCode_MageArcaneBlast {
				procChance = procChanceArcaneBlast
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

	procChance := .15
	buffDuration := time.Second * 15

	affectedSpells := []*core.Spell{}
	triggerSpellCodes := []int32{SpellCode_MageFireball, SpellCode_MageFrostfireBolt, SpellCode_MageSpellfrostBolt}

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
				}),
				func(spell *core.Spell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier -= 1
				spell.CostMultiplier -= 1
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier += 1
				spell.CostMultiplier += 1
			})
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !slices.Contains(triggerSpellCodes, spell.SpellCode) {
				return
			}
			aura.Deactivate(sim)
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Brain Freeze Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.Flags.Matches(SpellFlagChillSpell) {
				return
			}

			if sim.RandomFloat("Brain Freeze") < procChance {
				procAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) applySpellPower() {
	if !mage.HasRune(proto.MageRune_RuneFeetSpellPower) {
		return
	}

	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagMage) {
			spell.CritDamageBonus += 0.5
		}
	})
}
