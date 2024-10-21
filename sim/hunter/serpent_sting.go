package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) getSerpentStingConfig(rank int) core.SpellConfig {
	spellId := [10]int32{0, 1978, 13549, 13550, 13551, 13552, 13553, 13554, 13555, 25295}[rank]
	baseDamage := [10]float64{0, 20, 40, 80, 140, 210, 290, 385, 490, 555}[rank] / 5
	spellCoeff := [10]float64{0, .4, .625, .925, 1, 1, 1, 1, 1, 1}[rank] / 5
	manaCost := [10]float64{0, 15, 30, 50, 80, 115, 150, 190, 230, 250}[rank]
	level := [10]int{0, 4, 10, 18, 26, 34, 42, 50, 58, 60}[rank]

	return core.SpellConfig{
		SpellCode:     SpellCode_HunterSerpentSting,
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolNature,
		DefenseType:   core.DefenseTypeRanged,
		ProcMask:      core.ProcMaskRangedSpecial,
		Flags:         core.SpellFlagAPL | core.SpellFlagPureDot | core.SpellFlagPoison | SpellFlagSting,
		CastType:      proto.CastType_CastTypeRanged,
		Rank:          rank,
		RequiredLevel: level,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget >= core.MinRangedAttackDistance
		},

		DamageMultiplier: 1 + 0.02*float64(hunter.Talents.ImprovedSerpentSting),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SerpentSting" + hunter.Label + strconv.Itoa(rank),
				Tag:   "SerpentSting",
			},
			NumberOfTicks:    5,
			TickLength:       time.Second * 3,
			BonusCoefficient: spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// As of phase 5 the only time serpent sting scales with AP is using the Dragonstalker's Pursuit 6P - this AP scaling doesn't benefit from target AP modifiers
				damage := baseDamage + (hunter.SerpentStingAPCoeff*dot.Spell.RangedAttackPower(target, true))/5
				dot.Snapshot(target, damage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHitNoHitCounter)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealOutcome(sim, result)

				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			})
		},
	}
}

func (hunter *Hunter) chimeraShotSerpentStingSpell(rank int) *core.Spell {
	baseDamage := [10]float64{0, 20, 40, 80, 140, 210, 290, 385, 490, 555}[rank]
	spellCoeff := [10]float64{0, .4, .625, .925, 1, 1, 1, 1, 1, 1}[rank]
	return hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 409493},
		SpellSchool: core.SpellSchoolNature,
		CastType:    proto.CastType_CastTypeRanged,
		DefenseType: core.DefenseTypeRanged,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,

		BonusCritRating: 1,

		CritDamageBonus: hunter.mortalShots(),

		DamageMultiplier:         0.48,
		DamageMultiplierAdditive: 1 + 0.02*float64(hunter.Talents.ImprovedSerpentSting),
		ThreatMultiplier:         1,
		BonusCoefficient:         spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// As of phase 5 the only time serpent sting scales with AP is using the Dragonstalker's Pursuit 6P - this AP scaling doesn't benefit from target AP modifiers
			damage := baseDamage + (hunter.SerpentStingAPCoeff * spell.RangedAttackPower(target, true))
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeRangedHitAndCrit)
		},
	})
}

func (hunter *Hunter) registerSerpentStingSpell() {
	// TODO: AQ ranks := 9
	ranks := 8
	hunter.SerpentStingAPCoeff = 0

	for i := ranks; i >= 0; i-- {
		config := hunter.getSerpentStingConfig(i)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.SerpentSting = hunter.GetOrRegisterSpell(config)
			break
		}
	}
}
