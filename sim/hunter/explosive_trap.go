package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) getExplosiveTrapConfig(rank int, timer *core.Timer) core.SpellConfig {
	spellId := [4]int32{0, 409532, 409534, 409535}[rank]
	dotDamage := [4]float64{0, 15, 24, 33}[rank]
	minDamage := [4]float64{0, 104, 145, 208}[rank]
	maxDamage := [4]float64{0, 135, 193, 265}[rank]
	manaCost := [4]float64{0, 275, 395, 520}[rank]
	level := [4]int{0, 34, 44, 54}[rank]

	numHits := hunter.Env.GetNumTargets()

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 15,
			},
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		DamageMultiplierAdditive: 1 + 0.15*float64(hunter.Talents.CleverTraps),
		CritMultiplier:           hunter.critMultiplier(true, hunter.CurrentTarget),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "ExplosiveTrap" + hunter.Label + strconv.Itoa(rank),
				Tag:   "ExplosiveTrap",
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = dotDamage
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					baseDamage := sim.Roll(minDamage, maxDamage)
					baseDamage *= sim.Encounter.AOECapMultiplier()
					spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
				spell.AOEDot().ApplyOrReset(sim)
			})
		},
	}
}

func (hunter *Hunter) registerExplosiveTrapSpell(timer *core.Timer) {
	if !hunter.HasRune(proto.HunterRune_RuneBootsTrapLauncher) {
		return
	}

	maxRank := 3
	for i := 1; i <= maxRank; i++ {
		config := hunter.getExplosiveTrapConfig(i, timer)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.ExplosiveTrap = hunter.GetOrRegisterSpell(config)
		}
	}
}
