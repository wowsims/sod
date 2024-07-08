package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) getMultiShotConfig(rank int, timer *core.Timer) core.SpellConfig {
	spellId := [6]int32{0, 2643, 14288, 14289, 14290, 25294}[rank]
	baseDamage := [6]float64{0, 0, 40, 80, 120, 150}[rank]
	manaCost := [6]float64{0, 100, 140, 175, 210, 230}[rank]
	level := [6]int{0, 18, 30, 42, 54, 60}[rank]

	numHits := min(3, hunter.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	hasCobraStrikes := hunter.pet != nil && hunter.HasRune(proto.HunterRune_RuneChestCobraStrikes)
	hasSerpentSpread := hunter.HasRune(proto.HunterRune_RuneLegsSerpentSpread)

	manaCostMultiplier := 1 - 0.02*float64(hunter.Talents.Efficiency)
	if hunter.HasRune(proto.HunterRune_RuneChestMasterMarksman) {
		manaCostMultiplier -= 0.25
	}
	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolPhysical,
		DefenseType:   core.DefenseTypeRanged,
		ProcMask:      core.ProcMaskRangedSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		CastType:      proto.CastType_CastTypeRanged,
		Rank:          rank,
		RequiredLevel: level,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: manaCostMultiplier,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 500,
			},
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 10,
			},
			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.RangedSwingSpeed())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget >= core.MinRangedAttackDistance
		},

		CritDamageBonus: hunter.mortalShots(),

		DamageMultiplier: 1 + .05*float64(hunter.Talents.Barrage),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target

			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := baseDamage +
					hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target)) +
					hunter.NormalizedAmmoDamageBonus

				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeRangedHitAndCrit)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					spell.DealDamage(sim, results[hitIndex])

					if hasCobraStrikes && results[hitIndex].DidCrit() {
						hunter.CobraStrikesAura.Activate(sim)
						hunter.CobraStrikesAura.SetStacks(sim, 2)
					}

					if hasSerpentSpread {
						serpentStingAura := hunter.SerpentSting.Dot(curTarget)
						serpentStingTicks := serpentStingAura.NumberOfTicks
						if serpentStingAura.IsActive() {
							// If less then 4 ticks are left then we rollover with a 4 tick duration
							serpentStingAura.NumberOfTicks = max(4, serpentStingAura.NumberOfTicks-serpentStingAura.TickCount)
							serpentStingAura.Rollover(sim)
						} else {
							// Else we apply with a 4 tick duration
							serpentStingAura.NumberOfTicks = 4
							serpentStingAura.Apply(sim)
						}
						serpentStingAura.NumberOfTicks = serpentStingTicks
					}

					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			})

		},
	}
}

func (hunter *Hunter) registerMultiShotSpell(timer *core.Timer) {
	maxRank := 5

	for i := 1; i <= maxRank; i++ {
		config := hunter.getMultiShotConfig(i, timer)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.ArcaneShot = hunter.GetOrRegisterSpell(config)
		}
	}
}
