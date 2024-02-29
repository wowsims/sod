package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const SmiteRanks = 8

var SmiteSpellId = [SmiteRanks + 1]int32{0, 585, 591, 598, 984, 1004, 6060, 10933, 10934}
var SmiteBaseDamage = [SmiteRanks + 1][]float64{{0}, {15, 20}, {28, 34}, {58, 67}, {94, 109}, {158, 178}, {216, 244}, {296, 333}, {384, 429}}
var SmiteSpellCoef = [SmiteRanks + 1]float64{0, 0.123, 0.271, 0.554, 0.714, 0.714, 0.714, 0.714, 0.714}
var SmiteCastTime = [SmiteRanks + 1]int{0, 1500, 2000, 2500, 2500, 2500, 2500, 2500, 2500}
var SmiteManaCost = [SmiteRanks + 1]float64{0, 20, 30, 60, 95, 140, 185, 230, 280}
var SmiteLevel = [SmiteRanks + 1]int{0, 1, 6, 14, 22, 30, 38, 46, 54}

func (priest *Priest) RegisterSmiteSpell() {
	priest.Smite = make([]*core.Spell, SmiteRanks+1)

	for rank := 1; rank <= SmiteRanks; rank++ {
		config := priest.getSmiteBaseConfig(rank)

		if config.RequiredLevel <= int(priest.Level) {
			priest.Smite[rank] = priest.GetOrRegisterSpell(config)
		}
	}
}

func (priest *Priest) getSmiteBaseConfig(rank int) core.SpellConfig {
	spellId := SmiteSpellId[rank]
	baseDamageLow := SmiteBaseDamage[rank][0]
	baseDamageHigh := SmiteBaseDamage[rank][1]
	spellCoeff := SmiteSpellCoef[rank]
	castTime := SmiteCastTime[rank]
	manaCost := SmiteManaCost[rank]
	level := SmiteLevel[rank]

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*time.Duration(castTime) - time.Millisecond*100*time.Duration(priest.Talents.DivineFury),
			},
		},

		BonusCritRating:  priest.holySpecCritRating() + priest.forceOfWillCritRating(),
		DamageMultiplier: priest.searingLightDamageModifier() * priest.forceOfWillDamageModifier(),
		CritMultiplier:   priest.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}
}
