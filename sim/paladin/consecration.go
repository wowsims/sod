package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

const consecrationRanks = 5

var consecrationLevels = [consecrationRanks + 1]int{0, 20, 30, 40, 50, 60}
var consecrationSpellIDs = [consecrationRanks + 1]int32{0, 26573, 20116, 20922, 20923, 20924}
var consecrationBaseDamages = [consecrationRanks + 1]float64{0, 64 / 8, 120 / 8, 192 / 8, 280 / 8, 384 / 8}
var consecrationManaCosts = [consecrationRanks + 1]float64{0, 135, 235, 320, 435, 565}

func (paladin *Paladin) getConsecrationBaseConfig(rank int) core.SpellConfig {
	spellId := consecrationSpellIDs[rank]
	baseDamage := consecrationBaseDamages[rank]
	manaCost := consecrationManaCosts[rank]
	level := consecrationLevels[rank]

	spellCoeff := 0.042
	actionID := core.ActionID{SpellID: spellId}

	return core.SpellConfig{
		ActionID:      actionID,
		SpellSchool:   core.SpellSchoolHoly,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 8,
			},
		},
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Consecration-" + paladin.Label + strconv.Itoa(rank),
			},
			NumberOfTicks:       8,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: false,
			OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
				target := paladin.CurrentTarget
				dot.SnapshotBaseDamage = baseDamage + spellCoeff*dot.Spell.SpellDamage()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.Spell.OutcomeMagicHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	}
}

func (paladin *Paladin) registerConsecrationSpell() {
	if !paladin.Talents.Consecration {
		return
	}
	paladin.Consecration = make([]*core.Spell, consecrationRanks+1)
	for rank := 1; rank <= consecrationRanks; rank++ {
		config := paladin.getConsecrationBaseConfig(rank)
		if config.RequiredLevel <= int(paladin.Level) {
			paladin.Consecration[rank] = paladin.GetOrRegisterSpell(config)
		}
	}
}
