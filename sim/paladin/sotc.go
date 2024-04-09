package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

type SealOfTheCrusaderRankInfo struct {
	SealId      int32
	ManaCost    float64
	AttackPower float64
}

var sotcRanks = map[int32]SealOfTheCrusaderRankInfo{
	25: {
		SealId:      20305,
		ManaCost:    65,
		AttackPower: 99.1,
	},
	40: {
		SealId:      20306,
		ManaCost:    90,
		AttackPower: 161,
	},
	50: {
		SealId:      20307,
		ManaCost:    125,
		AttackPower: 238.6,
	},
	60: {
		SealId:      20308,
		ManaCost:    160,
		AttackPower: 325.2,
	},
}

func (paladin *Paladin) hasLibramOfFervor() bool {
	return paladin.Ranged().ID == LibramOfFervor
}

func makeJudgementOfTheCrusader(paladin *Paladin) *core.Spell {
	debuffs := paladin.NewEnemyAuraArray(func(u *core.Unit, i int32) *core.Aura {
		mult := paladin.GetImprovedSealOfTheCrusaderMult()
		extraBonus := core.TernaryFloat64(paladin.hasLibramOfFervor(), 33, 0)
		return core.JudgementOfTheCrusaderAura(u, i, mult, extraBonus)
	})

	return paladin.RegisterSpell(core.SpellConfig{
		ActionID:    debuffs[0].ActionID,
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeAlwaysHit)
			if result.Landed() {
				debuffs.Get(target).Activate(sim)
			}
		},
	})
}

func (paladin *Paladin) registerSealOfTheCrusader() {
	rankInfo := sotcRanks[paladin.Level]
	spellAction := core.ActionID{SpellID: rankInfo.SealId}

	apAtLevel := rankInfo.AttackPower * paladin.GetImprovedSealOfTheCrusaderMult()
	sotcMultiplier := 1.4

	if paladin.hasLibramOfFervor() {
		apAtLevel += 48
	}

	sealAura := paladin.RegisterAura(core.Aura{
		Label:    "Seal of the Crusader" + paladin.Label,
		Tag:      "Seal",
		ActionID: spellAction,
		Duration: SealDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.MultiplyMeleeSpeed(sim, sotcMultiplier)
			paladin.AutoAttacks.MHAuto().DamageMultiplier /= sotcMultiplier
			paladin.AddStatDynamic(sim, stats.AttackPower, apAtLevel)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.MultiplyMeleeSpeed(sim, 1/sotcMultiplier)
			paladin.AutoAttacks.MHAuto().DamageMultiplier *= sotcMultiplier
			paladin.AddStatDynamic(sim, stats.AttackPower, -apAtLevel)
		},
	})

	judgementSpell := makeJudgementOfTheCrusader(paladin)

	paladin.SealOfTheCrusader = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    spellAction,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   rankInfo.ManaCost - paladin.GetLibramSealCostReduction(),
			Multiplier: 1.0 - (float64(paladin.Talents.Benediction) * 0.03),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			paladin.ApplySeal(sealAura, judgementSpell, sim)
		},
	})
}
