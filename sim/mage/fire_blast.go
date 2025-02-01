package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const FireBlastRanks = 7

var FireBlastSpellId = [FireBlastRanks + 1]int32{0, 2136, 2137, 2138, 8412, 8413, 10197, 10199}
var FireBlastBaseDamage = [FireBlastRanks + 1][]float64{{0}, {27, 35}, {62, 76}, {107, 132}, {177, 211}, {246, 295}, {342, 405}, {446, 524}}
var FireBlastSpellCoeff = [FireBlastRanks + 1]float64{0, .204, .332, .429, .429, .429, .429, .429}
var FireBlastManaCost = [FireBlastRanks + 1]float64{0, 40, 75, 115, 165, 220, 280, 340}
var FireBlastLevel = [FireBlastRanks + 1]int{0, 6, 14, 22, 30, 38, 46, 54}

func (mage *Mage) registerFireBlastSpell() {
	mage.FireBlast = make([]*core.Spell, FireBlastRanks+1)
	cdTimer := mage.NewTimer()

	for rank := 1; rank <= FireBlastRanks; rank++ {
		config := mage.newFireBlastSpellConfig(rank, cdTimer)

		if config.RequiredLevel <= int(mage.Level) {
			mage.FireBlast[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) newFireBlastSpellConfig(rank int, cdTimer *core.Timer) core.SpellConfig {
	spellId := FireBlastSpellId[rank]
	baseDamageLow := FireBlastBaseDamage[rank][0]
	baseDamageHigh := FireBlastBaseDamage[rank][1]
	spellCoeff := FireBlastSpellCoeff[rank]
	manaCost := FireBlastManaCost[rank]
	level := FireBlastLevel[rank]

	cooldown := time.Second * 8

	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellId},
		ClassSpellMask: ClassSpellMask_MageFireBlast,
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,

		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamageCacl := (baseDamageLow + baseDamageHigh) / 2
			return spell.CalcDamage(sim, target, baseDamageCacl, spell.OutcomeExpectedMagicHitAndCrit)
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}
}
