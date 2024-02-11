package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) getConsecrationBaseConfig(rank int) core.SpellConfig {
	spellId := [4]int32{0, 26573, 20116, 20922}[rank]
	baseDamage := [4]float64{0, 8, 15, 24}[rank]
	manaCost := [4]float64{0, 135, 235, 320}[rank]
	level := [4]int{0, 20, 30, 40}[rank]

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

	maxRank := 3 // for p2, can update for p3 onwards
	for i := 1; i <= maxRank; i++ {
		config := paladin.getConsecrationBaseConfig(i)
		if config.RequiredLevel <= int(paladin.Level) {
			paladin.Consecration = paladin.GetOrRegisterSpell(config)
		}
	}
}
