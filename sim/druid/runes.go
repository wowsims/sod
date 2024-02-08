package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (druid *Druid) ApplyRunes() {
	druid.applyEclipse()
	druid.applyFuryOfStormRage()
	druid.applySunfire()
	druid.applyStarsurge()
	druid.applyMangle()
	druid.applySavageRoar()
	druid.applyWildStrikes()
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

	// Solar
	druid.SolarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:     "Solar Eclipse proc",
		Duration:  time.Second * 15,
		MaxStacks: 4,
		ActionID:  core.ActionID{SpellID: 408250},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(
				core.FilterSlice(druid.Wrath, func(spell *DruidSpell) bool { return spell != nil }),
				func(spell *DruidSpell) {
					spell.BonusCritRating += solarCritBonus
				},
			)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(
				core.FilterSlice(druid.Wrath, func(spell *DruidSpell) bool { return spell != nil }),
				func(spell *DruidSpell) {
					spell.BonusCritRating -= solarCritBonus
				},
			)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Assert we are casting wrath
			if !result.Landed() || spell.SpellCode != SpellCode_DruidWrath {
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
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(
				core.FilterSlice(druid.Starfire, func(spell *DruidSpell) bool { return spell != nil }),
				func(spell *DruidSpell) {
					spell.DefaultCast.CastTime -= lunarCastTimeReduction
				},
			)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(
				core.FilterSlice(druid.Starfire, func(spell *DruidSpell) bool { return spell != nil }),
				func(spell *DruidSpell) {
					spell.DefaultCast.CastTime += lunarCastTimeReduction
				},
			)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Assert we are casting Starfire
			if !result.Landed() || spell.SpellCode != SpellCode_DruidStarfire {
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
			switch spell.SpellCode {
			case SpellCode_DruidWrath:
				// Proc Lunar
				druid.LunarEclipseProcAura.Activate(sim)
				// Wrath gives 1 stack of starfire bonus
				druid.LunarEclipseProcAura.AddStack(sim)
			case SpellCode_DruidStarfire:
				// Proc Solar
				druid.SolarEclipseProcAura.Activate(sim)
				// Starfire gives 2 staacks of wrath bonus
				druid.SolarEclipseProcAura.AddStack(sim)
				druid.SolarEclipseProcAura.AddStack(sim)
			default:
				return
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

// TODO: Classic verify star surge numbers
func (druid *Druid) applyStarsurge() {
	if !druid.HasRune(proto.DruidRune_RuneLegsStarsurge) {
		return
	}

	level := float64(druid.GetCharacter().Level)
	baseCalc := (9.183105 + 0.616405*level + 0.028608*level*level)
	baseLowDamage := baseCalc * 3.81
	baseHighDamage := baseCalc * 4.67

	druid.Starsurge = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 417157},
		SpellCode:   SpellCode_DruidStarsurge,
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.01 * (1 - 0.03*float64(druid.Talents.Moonglow)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 0,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   druid.VengeanceCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage)*druid.MoonfuryDamageMultiplier() + spell.SpellDamage()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
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
