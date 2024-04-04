package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const HurricaneRanks = 3

var HurricaneSpellId = [HurricaneRanks + 1]int32{0, 16914, 17401, 17402}
var HurricaneBaseDamage = [HurricaneRanks + 1]float64{0, 72, 102, 134}
var HurricaneSpellCoef = [HurricaneRanks + 1]float64{0, .03, .03, .03}
var HurricaneManaCost = [HurricaneRanks + 1]float64{0, 880, 1180, 1495}
var HurricaneLevel = [HurricaneRanks + 1]int{0, 40, 50, 60}

func (druid *Druid) registerHurricaneSpell() {
	druid.Hurricane = make([]*DruidSpell, HurricaneRanks+1)

	cooldownTimer := druid.NewTimer()

	for rank := 1; rank <= HurricaneRanks; rank++ {
		config := druid.newHurricaneSpellConfig(rank, cooldownTimer)

		if config.RequiredLevel <= int(druid.Level) {
			druid.Hurricane[rank] = druid.RegisterSpell(Humanoid|Moonkin, config)
		}
	}
}

func (druid *Druid) newHurricaneSpellConfig(rank int, cooldownTimer *core.Timer) core.SpellConfig {
	spellId := HurricaneSpellId[rank]
	baseDamage := HurricaneBaseDamage[rank]
	spellCoeff := HurricaneSpellCoef[rank]
	manaCost := HurricaneManaCost[rank]
	level := HurricaneLevel[rank]

	damageMultiplier := 1.0
	cooldown := time.Second * 60

	if druid.HasRune(proto.DruidRune_RuneHelmGaleWinds) {
		damageMultiplier += 1.0
		cooldown = core.GCDDefault
		manaCost *= .80
	}

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagOmen | core.SpellFlagChanneled | core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1,
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

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Hurricane",
			},
			NumberOfTicks:    10,
			TickLength:       time.Second * 1,
			BonusCoefficient: spellCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, baseDamage, false)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealPeriodicDamage(sim, aoeTarget, dot.SnapshotBaseDamage, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			druid.AutoAttacks.CancelAutoSwing(sim)
			spell.AOEDot().Apply(sim)
		},
	}
}
