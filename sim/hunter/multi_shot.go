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

	//numHits := min(3+hunter.MultiShotBonusTargets, hunter.Env.GetNumTargets())
	//results := make([]*core.SpellResult, numHits)

	hasSerpentSpread := hunter.HasRune(proto.HunterRune_RuneLegsSerpentSpread)

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterMultiShot,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolPhysical,
		CastType:       proto.CastType_CastTypeRanged,
		DefenseType:    core.DefenseTypeRanged,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		Rank:          rank,
		RequiredLevel: level,
		MinRange:      core.MinRangedAttackRange,
		MaxRange:      core.MaxRangedAttackRange,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 500,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
				hunter.AutoAttacks.CancelAutoSwing(sim)
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

		CritDamageBonus: hunter.mortalShots(),

		DamageMultiplier: 1 + .05*float64(hunter.Talents.Barrage),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		JumpTargets: 3,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseDamage +
				hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target, false)) +
				hunter.AmmoDamageBonus

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			hunter.AutoAttacks.EnableAutoSwing(sim)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)

				if hasSerpentSpread {
					serpentStingAura := hunter.SerpentSting.Dot(result.Target)
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
			})

		},
	}
}

func (hunter *Hunter) registerMultiShotSpell(timer *core.Timer) {
	maxRank := core.TernaryInt(core.IncludeAQ, 5, 4)
	for rank := 1; rank <= maxRank; rank++ {
		config := hunter.getMultiShotConfig(rank, timer)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.MultiShot = hunter.GetOrRegisterSpell(config)
		}
	}
}
