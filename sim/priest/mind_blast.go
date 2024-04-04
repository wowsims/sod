package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const MindBlastRanks = 9

var MindBlastSpellId = [MindBlastRanks + 1]int32{0, 8092, 8102, 8103, 8104, 8105, 8106, 10945, 10946, 10947}
var MindBlastBaseDamage = [MindBlastRanks + 1][]float64{{0}, {42, 46}, {76, 83}, {115, 124}, {174, 184}, {225, 239}, {279, 297}, {354, 375}, {437, 461}, {508, 537}}
var MindBlastSpellCoef = [MindBlastRanks + 1]float64{0, .268, .364, .429, .429, .429, .429, .429, .429, .429}
var MindBlastManaCost = [MindBlastRanks + 1]float64{0, 50, 80, 110, 150, 185, 225, 265, 310, 350}
var MindBlastLevel = [MindBlastRanks + 1]int{0, 10, 16, 22, 28, 34, 40, 46, 52, 58}

func (priest *Priest) registerMindBlast() {
	priest.MindBlast = make([]*core.Spell, MindBlastRanks+1)
	cdTimer := priest.NewTimer()

	for rank := 1; rank <= MindBlastRanks; rank++ {
		config := priest.getMindBlastBaseConfig(rank, cdTimer)

		if config.RequiredLevel <= int(priest.Level) {
			priest.MindBlast[rank] = priest.GetOrRegisterSpell(config)
		}
	}
}

func (priest *Priest) getMindBlastBaseConfig(rank int, cdTimer *core.Timer) core.SpellConfig {
	spellId := MindBlastSpellId[rank]
	baseDamageLow := MindBlastBaseDamage[rank][0] * priest.darknessDamageModifier()
	baseDamageHigh := MindBlastBaseDamage[rank][1] * priest.darknessDamageModifier()
	spellCoeff := MindBlastSpellCoef[rank]
	castTime := time.Millisecond * 1500
	manaCost := MindBlastManaCost[rank]
	level := MindBlastLevel[rank]

	hasPainAndSuffering := priest.HasRune(proto.PriestRune_RuneHelmPainAndSuffering)
	hasMindSpike := priest.HasRune(proto.PriestRune_RuneWaistMindSpike)

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
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
				CastTime: castTime,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: time.Second*8 - time.Millisecond*500*time.Duration(priest.Talents.ImprovedMindBlast),
			},
		},

		BonusCritRating: priest.forceOfWillCritRating(),
		BonusHitRating:  priest.shadowHitModifier(),

		DamageMultiplier: priest.forceOfWillDamageModifier(),
		ThreatMultiplier: priest.shadowThreatModifier(),
		BonusCoefficient: spellCoeff,

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			damage := (baseDamageLow + baseDamageHigh) / 2
			spell.DamageMultiplier *= priest.MindBlastModifier
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeExpectedMagicHitAndCrit)
			spell.DamageMultiplier /= priest.MindBlastModifier
			return result
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (sim.Roll(baseDamageLow, baseDamageHigh))

			bonusCrit := 0.0
			if hasMindSpike && priest.MindSpikeAuras.Get(target).IsActive() {
				bonusCrit = float64(priest.MindSpikeAuras.Get(target).GetStacks()) * 30 * core.CritRatingPerCritChance
			}

			spell.BonusCritRating += bonusCrit
			spell.DamageMultiplier *= priest.MindBlastModifier
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.BonusCritRating -= bonusCrit
			spell.DamageMultiplier /= priest.MindBlastModifier

			if result.Landed() {
				priest.AddShadowWeavingStack(sim, target)
				if hasMindSpike {
					priest.MindSpikeAuras.Get(target).Deactivate(sim)
				}

				if hasPainAndSuffering {
					for _, spell := range priest.ShadowWordPain {
						if spell != nil && spell.Dot(target).IsActive() {
							spell.Dot(target).Rollover(sim)
						}
					}
				}
			}

			spell.DealDamage(sim, result)
		},
	}
}
