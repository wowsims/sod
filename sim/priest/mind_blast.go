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
	baseDamageLow := MindBlastBaseDamage[rank][0]
	baseDamageHigh := MindBlastBaseDamage[rank][1]
	spellCoeff := MindBlastSpellCoef[rank]
	castTime := time.Millisecond * 1500
	manaCost := MindBlastManaCost[rank]
	level := MindBlastLevel[rank]

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolShadow,
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

		BonusCritRating:  priest.forceOfWillCritRating(),
		BonusHitRating:   priest.shadowHitModifier(),
		DamageMultiplier: priest.forceOfWillDamageModifier(),
		CritMultiplier:   priest.DefaultSpellCritMultiplier(),
		ThreatMultiplier: priest.shadowThreatModifier(),

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			damage := ((baseDamageLow+baseDamageHigh)/2 + spellCoeff*spell.SpellDamage()) * priest.MindBlastModifier
			return spell.CalcDamage(sim, target, damage, spell.OutcomeExpectedMagicHitAndCrit)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.BonusCritRating += float64(30 * priest.MindSpikeAuras.Get(target).GetStacks() * core.CritRatingPerCritChance)
			baseDamage := (sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()) * priest.MindBlastModifier
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim, target)
				if priest.HasRune(proto.PriestRune_RuneWaistMindSpike) {
					priest.MindSpikeAuras.Get(target).Deactivate(sim)
				}
			}

			spell.DealDamage(sim, result)
		},
	}
}
