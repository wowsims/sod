package mage

import (
	"github.com/wowsims/sod/sim/core"
)

const ArcaneExplosionRanks = 6

var ArcaneExplosionSpellId = [ArcaneExplosionRanks + 1]int32{0, 1449, 8437, 8438, 8439, 10201, 10202}
var ArcaneExplosionBaseDamage = [ArcaneExplosionRanks + 1][]float64{{0}, {34, 38}, {58, 65}, {101, 110}, {140, 153}, {190, 207}, {249, 270}}
var ArcaneExplosionSpellCoeff = [ArcaneExplosionRanks + 1]float64{0, .111, .143, .143, .143, .143, .143}
var ArcaneExplosionManaCost = [ArcaneExplosionRanks + 1]float64{0, 75, 120, 185, 250, 315, 390}
var ArcaneExplosionLevel = [ArcaneExplosionRanks + 1]int{0, 14, 22, 30, 38, 46, 54}

func (mage *Mage) registerArcaneExplosionSpell() {
	mage.ArcaneExplosion = make([]*core.Spell, ArcaneExplosionRanks+1)

	for rank := 1; rank <= ArcaneExplosionRanks; rank++ {
		config := mage.newArcaneExplosionSpellConfig(rank)

		if config.RequiredLevel <= int(mage.Level) {
			mage.ArcaneExplosion[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) newArcaneExplosionSpellConfig(rank int) core.SpellConfig {
	spellId := ArcaneExplosionSpellId[rank]
	baseDamageLow := ArcaneExplosionBaseDamage[rank][0]
	baseDamageHigh := ArcaneExplosionBaseDamage[rank][1]
	spellCoeff := ArcaneExplosionSpellCoeff[rank]
	manaCost := ArcaneExplosionManaCost[rank]
	level := ArcaneExplosionLevel[rank]

	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellId},
		ClassSpellMask: ClassSpellMask_MageArcaneExplosion,
		SpellSchool:    core.SpellSchoolArcane,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				damage := sim.Roll(baseDamageLow, baseDamageHigh)
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicCrit)
			}
		},
	}
}
