package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const ConeOfColdRanks = 5

var ConeOfColdSpellId = [ConeOfColdRanks + 1]int32{0, 120, 8492, 10159, 10160, 10161}
var ConeOfColdBaseDamage = [ConeOfColdRanks + 1][2]float64{{0, 0}, {98, 108}, {146, 160}, {203, 223}, {264, 290}, {335, 365}}
var ConeOfColdManaCost = [ConeOfColdRanks + 1]float64{0, 210, 290, 380, 465, 555}
var ConeOfColdLevel = [ConeOfColdRanks + 1]int{0, 26, 34, 42, 50, 58}

func (mage *Mage) registerConeOfColdSpell() {
	mage.ConeOfCold = make([]*core.Spell, ConeOfColdRanks+1)

	for rank := 1; rank <= ConeOfColdRanks; rank++ {
		config := mage.newConeOfColdSpellConfig(rank)

		if config.RequiredLevel <= int(mage.Level) {
			mage.ConeOfCold[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) newConeOfColdSpellConfig(rank int) core.SpellConfig {
	spellId := ConeOfColdSpellId[rank]
	baseDamageLow := ConeOfColdBaseDamage[rank][0]
	baseDamageHigh := ConeOfColdBaseDamage[rank][1]
	manaCost := ConeOfColdManaCost[rank]
	level := ConeOfColdLevel[rank]

	spellCoeff := 0.129

	numHits := len(mage.Env.Encounter.TargetUnits)

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_MageConeOfCold,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolFrost,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | SpellFlagChillSpell,

		RequiredLevel: level,
		Rank:          rank,
		MinRange:      10,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for i := 0; i < numHits; i++ {
				spell.CalcAndDealDamage(sim, target, sim.Roll(baseDamageLow, baseDamageHigh), spell.OutcomeMagicHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}
		},
	}
}
