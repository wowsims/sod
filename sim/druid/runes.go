package druid

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (druid *Druid) ApplyRunes() {
	// Chest
	druid.applyFuryOfStormRage()
	// druid.applyLivingSeed()
	// druid.applySurvivalOfTheFittest()
	druid.applyWildStrikes()

	// Hands
	// druid.applyLacerate()
	druid.applyMangle()
	druid.applySunfire()
	// druid.applyWildGrowth()

	// Belt
	// druid.applyBerserk()
	druid.applyEclipse()
	// druid.applyNourish()

	// Legs
	druid.applyStarsurge()
	druid.applySavageRoar()
	// druid.applyLifebloom()
	// druid.applySkullBash()

	// Feet
	druid.applyDreamstate()
	// druid.applyKingOfTheJungle()
	// druid.applySurvivalInstincts()
}

func (druid *Druid) applyFuryOfStormRage() {
	if !druid.HasRune(proto.DruidRune_RuneChestFuryOfStormrage) {
		return
	}

	druid.FuryOfStormrageAura = druid.RegisterAura(core.Aura{
		Label:    "Fury Of Stormrage",
		ActionID: core.ActionID{SpellID: int32(proto.DruidRune_RuneChestFuryOfStormrage)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (druid *Druid) applyEclipse() {
	if !druid.HasRune(proto.DruidRune_RuneBeltEclipse) {
		return
	}

	solarCritBonus := 30.0
	lunarCastTimeReduction := time.Second * 1

	var affectedSolarSpells []*DruidSpell
	var affectedLunarSpells []*DruidSpell

	// Solar
	druid.SolarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:     "Solar Eclipse proc",
		Duration:  time.Second * 15,
		MaxStacks: 4,
		ActionID:  core.ActionID{SpellID: 408250},
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSolarSpells = core.FilterSlice(
				core.Flatten([][]*DruidSpell{druid.Wrath, {druid.Starsurge}}),
				func(spell *DruidSpell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSolarSpells, func(spell *DruidSpell) {
				spell.BonusCritRating += solarCritBonus
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSolarSpells, func(spell *DruidSpell) {
				spell.BonusCritRating -= solarCritBonus
			})
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode != SpellCode_DruidWrath && spell.SpellCode != SpellCode_DruidStarsurge {
				return
			}

			aura.RemoveStack(sim)
		},
	})

	// Lunar
	druid.LunarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:     "Lunar Eclipse proc",
		Duration:  time.Second * 15,
		MaxStacks: 4,
		ActionID:  core.ActionID{SpellID: 408255},
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedLunarSpells = core.FilterSlice(
				core.Flatten([][]*DruidSpell{druid.Starfire}),
				func(spell *DruidSpell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedLunarSpells, func(spell *DruidSpell) {
				spell.DefaultCast.CastTime -= lunarCastTimeReduction
			})
			druid.Starsurge.BonusCritRating += 30.0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedLunarSpells, func(spell *DruidSpell) {
				spell.DefaultCast.CastTime += lunarCastTimeReduction
			})
			druid.Starsurge.BonusCritRating -= 30.0
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode != SpellCode_DruidStarfire {
				return
			}

			aura.RemoveStack(sim)
		},
	})

	druid.EclipseAura = druid.RegisterAura(core.Aura{
		Label:    "Eclipse",
		Duration: core.NeverExpires,
		ActionID: core.ActionID{SpellID: 408248},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !slices.Contains([]int32{SpellCode_DruidWrath, SpellCode_DruidStarfire, SpellCode_DruidStarsurge}, spell.SpellCode) || !result.Landed() {
				return
			}

			if spell.SpellCode == SpellCode_DruidWrath || spell.SpellCode == SpellCode_DruidStarsurge {
				druid.LunarEclipseProcAura.Activate(sim)
				// Solar gives 1 stack of lunar bonus
				druid.LunarEclipseProcAura.AddStack(sim)
			}

			if spell.SpellCode == SpellCode_DruidStarfire || spell.SpellCode == SpellCode_DruidStarsurge {
				druid.SolarEclipseProcAura.Activate(sim)
				// Lunar gives 2 staacks of solar bonus
				druid.SolarEclipseProcAura.AddStacks(sim, 2)
			}
		},
	})

}

// https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044#news-post-336044
func (druid *Druid) applySunfire() {
	if !druid.HasRune(proto.DruidRune_RuneHandsSunfire) {
		return
	}

	level := float64(druid.GetCharacter().Level)
	baseCalc := (9.183105 + 0.616405*level + 0.028608*level*level)
	baseLowDamage := baseCalc * 1.3
	baseHighDamage := baseCalc * 1.52
	spellCoeff := .15
	spellDotCoeff := .13
	baseDotDamage := baseCalc * 0.65
	ticks := int32(4)

	druid.Sunfire = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 414684},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.21,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 0,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Sunfire",
				ActionID: core.ActionID{SpellID: 414684},
			},
			NumberOfTicks: ticks,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = baseDotDamage*druid.MoonfuryDamageMultiplier() + spellDotCoeff*dot.Spell.SpellDamage()
				dot.SnapshotAttackerMultiplier = 1
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		BonusCritRating:  druid.ImprovedMoonfireCritBonus() * core.SpellCritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   druid.VengeanceCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage)*druid.MoonfuryDamageMultiplier()*druid.ImprovedMoonfireDamageMultiplier() + spellCoeff*spell.SpellDamage()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}

func (druid *Druid) applyMangle() {
	druid.applyMangleCat()
	//druid.applyMangleBear()
}

func (druid *Druid) applyWildStrikes() {
	if !druid.HasRune(proto.DruidRune_RuneChestWildStrikes) {
		return
	}

	druid.WildStrikesBuffAura = core.ApplyWildStrikes(druid.GetCharacter())
}

func (druid *Druid) applyDreamstate() {
	if !druid.HasRune(proto.DruidRune_RuneFeetDreamstate) {
		return
	}

	manaRegenDuration := time.Second * 8

	manaRegenAura := druid.RegisterAura(core.Aura{
		Label:    "Dreamstate Mana Regen",
		ActionID: core.ActionID{SpellID: int32(proto.DruidRune_RuneFeetDreamstate)},
		Duration: manaRegenDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.SpiritRegenRateCasting += .5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.SpiritRegenRateCasting -= .5
		},
	})

	druid.RegisterAura(core.Aura{
		Label:    "Dreamstate Trigger",
		ActionID: core.ActionID{SpellID: int32(proto.DruidRune_RuneFeetDreamstate)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}

			manaRegenAura.Activate(sim)
			core.DreamstateAura(result.Target).Activate(sim)
		},
	})
}
