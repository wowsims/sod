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

	// Solar
	solarProcMultiplier := 30.0
	druid.SolarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:     "Solar Eclipse proc",
		Duration:  time.Second * 15,
		MaxStacks: 4,
		ActionID:  core.ActionID{SpellID: 408250},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(druid.Wrath, func(spell *DruidSpell) {
				if spell == nil {
					return
				}
				spell.BonusCritRating += solarProcMultiplier
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(druid.Wrath, func(spell *DruidSpell) {
				if spell == nil {
					return
				}
				spell.BonusCritRating -= solarProcMultiplier
			})
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Assert we are  casting wrath
			if spell.SpellCode != SpellCode_DruidWrath {
				return
			}

			if !result.Landed() {
				return
			}

			aura.RemoveStack(sim)
		},
	})

	// Lunar
	lunarBonusCrit := 30.0
	druid.LunarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:     "Lunar Eclipse proc",
		Duration:  time.Second * 15,
		MaxStacks: 4,
		ActionID:  core.ActionID{SpellID: 408255},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(druid.Starfire, func(spell *DruidSpell) {
				if spell == nil {
					return
				}
				spell.BonusCritRating += lunarBonusCrit
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(druid.Starfire, func(spell *DruidSpell) {
				if spell == nil {
					return
				}
				spell.BonusCritRating -= lunarBonusCrit
			})
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Assert we are casting Starfire
			if spell.SpellCode != SpellCode_DruidStarfire {
				return
			}
			if !result.Landed() {
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
				druid.LunarEclipseProcAura.SetStacks(sim, 1)

			case SpellCode_DruidStarfire:
				// Proc Solar
				druid.SolarEclipseProcAura.Activate(sim)
				druid.SolarEclipseProcAura.SetStacks(sim, 2)

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
	baseDotDamage := baseCalc * 0.65
	ticks := int32(4)
	impMf := float64(druid.Talents.ImprovedMoonfire)
	moonfury := float64(druid.Talents.Moonfury)

	druid.Sunfire = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 414684},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,
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
				dot.SnapshotBaseDamage = baseDotDamage + 0.13*dot.Spell.SpellDamage()
				dot.SnapshotAttackerMultiplier = 1
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		BonusCritRating:  2 * impMf * core.SpellCritRatingPerCritChance,
		DamageMultiplier: 1 + 0.02*impMf + 0.02*moonfury,
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + 0.15*spell.SpellDamage()
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
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,
		SpellCode:   int32(SpellCode_DruidStarfire), // Please check if this is right - Starsurge affected by all Starfire talents/procs?
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

		BonusCritRating:  0,
		DamageMultiplier: 1 + 0.02*float64(druid.Talents.Moonfury),
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + spell.SpellDamage()
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
