package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const BlastWaveRanks = 5

var BlastWaveSpellId = [BlastWaveRanks + 1]int32{0, 11113, 13018, 13019, 13020, 13021}
var BlastWaveBaseDamage = [BlastWaveRanks + 1][]float64{{0}, {160, 192}, {205, 246}, {285, 338}, {374, 443}, {462, 544}}
var BlastWaveManaCost = [BlastWaveRanks + 1]float64{0, 215, 270, 355, 450, 545}
var BlastWaveLevel = [BlastWaveRanks + 1]int{0, 30, 36, 44, 52, 60}

func (mage *Mage) registerBlastWaveSpell() {
	if !mage.Talents.BlastWave {
		return
	}

	mage.BlastWave = make([]*core.Spell, BlastWaveRanks+1)
	cdTimer := mage.NewTimer()

	for rank := 1; rank <= BlastWaveRanks; rank++ {
		config := mage.newBlastWaveSpellConfig(rank, cdTimer)

		if config.RequiredLevel <= int(mage.Level) {
			mage.BlastWave[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) newBlastWaveSpellConfig(rank int, cooldownTimer *core.Timer) core.SpellConfig {
	spellId := BlastWaveSpellId[rank]
	baseDamageLow := BlastWaveBaseDamage[rank][0]
	baseDamageHigh := BlastWaveBaseDamage[rank][1]
	manaCost := BlastWaveManaCost[rank]
	level := BlastWaveLevel[rank]

	spellCoeff := .129
	cooldown := time.Second * 45

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagBinary | core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    cooldownTimer,
				Duration: cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicCrit)
			}
		},
	}
}
