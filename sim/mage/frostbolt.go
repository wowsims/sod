package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const FrostboltRanks = 11

var FrostboltSpellId = [FrostboltRanks + 1]int32{0, 116, 205, 837, 7322, 8406, 8407, 8408, 10179, 10180, 10181, 25304}
var FrostboltBaseDamage = [FrostboltRanks + 1][]float64{{0, 0}, {20, 22}, {33, 38}, {54, 61}, {78, 87}, {132, 144}, {180, 197}, {231, 251}, {301, 326}, {353, 383}, {440, 475}, {515, 555}}
var FrostboltSpellCoeff = [FrostboltRanks + 1]float64{0, .163, .269, .463, .706, .814, .814, .814, .814, .814, .814, .814}
var FrostboltCastTime = [FrostboltRanks + 1]int32{0, 1500, 1800, 2200, 2600, 3000, 3000, 3000, 3000, 3000, 3000, 3000}
var FrostboltManaCost = [FrostboltRanks + 1]float64{0, 25, 35, 50, 65, 100, 130, 160, 195, 225, 260, 290}
var FrostboltLevel = [FrostboltRanks + 1]int{0, 4, 8, 14, 20, 26, 32, 38, 44, 50, 56, 60}

func (mage *Mage) registerFrostboltSpell() {
	mage.Frostbolt = make([]*core.Spell, FrostboltRanks+1)

	for rank := 1; rank <= FrostboltRanks; rank++ {
		config := mage.getFrostboltConfig(rank)

		if config.RequiredLevel <= int(mage.Level) {
			mage.Frostbolt[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) getFrostboltConfig(rank int) core.SpellConfig {
	spellId := FrostboltSpellId[rank]
	baseDamageLow := FrostboltBaseDamage[rank][0]
	baseDamageHigh := FrostboltBaseDamage[rank][1]
	spellCoeff := FrostboltSpellCoeff[rank]
	castTime := FrostboltCastTime[rank]
	manaCost := FrostboltManaCost[rank]
	level := FrostboltLevel[rank]

	return core.SpellConfig{
		ActionID:     core.ActionID{SpellID: spellId},
		SpellCode:    SpellCode_MageFrostbolt,
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagAPL | SpellFlagMage | SpellFlagChillSpell,
		MissileSpeed: 28,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*time.Duration(castTime) - time.Millisecond*100*time.Duration(mage.Talents.ImprovedFrostbolt),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.MageCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					spell.DealDamage(sim, result)
				}
			})
		},
	}
}
